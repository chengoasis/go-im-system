package main

import (
	"net"
	"strings"
	"sync"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	// 整个程序里通常只有一个 Server 对象，就是在 main.go (line 3) 里 NewServer(...) 创建出来的那个。
	// 每个 User 如果加上 server *Server 字段，保存的不是一份新的服务器副本，而只是“指向这个公共服务器的地址”。
	server *Server
	// 避免客户端断开、超时、重复错误处理时多次广播下线、多次关闭 channel/conn
	offlineOnce sync.Once
}

// 创建一个用户 API
func NewUser(conn net.Conn, server *Server) *User {
	// RemoteAddr() 是网络连接 conn 自带的方法，用来获取客户端的 IP 和端口
	// 客户端的 IP:Port，如 "192.168.1.100:54321"
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,

		server: server,
	}

	// 启动监听当前 user channel 消息的 goroutine
	go user.ListenMessage()

	return user
}

// 用户上线业务
func (this *User) Online() {
	// 用户上线， 将用户加入到 server 的 onlineMap 中
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	// 广播当前用户的上线消息
	this.server.BroadCast(this, "已上线")
}

// 用户下线业务
func (this *User) offline() {
	// 这里的代码只会执行一次
	// 即使多个 goroutine 同时调用 offline()
	// 也只有一个能进入这里执行
	this.offlineOnce.Do(func() {
		// 用户下线，将用户从 onlineMap 中删除
		this.server.mapLock.Lock()
		delete(this.server.OnlineMap, this.Name)
		this.server.mapLock.Unlock()

		// 广播当前用户下线
		this.server.BroadCast(this, "下线")

		close(this.C)
		this.conn.Close()
	})
}

// 当前 User 对应的客户端发送消息
func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}

// 用户处理消息业务
func (this *User) DoMessage(msg string) {
	if msg == "who" {
		// 查询当前用户
		this.server.mapLock.RLock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "在线...\n"
			this.SendMsg(onlineMsg)
		}
		this.server.mapLock.RUnlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		// 消息格式：rename|程
		newName := strings.SplitN(msg, "|", 2)[1]

		this.server.mapLock.Lock()
		if _, ok := this.server.OnlineMap[newName]; ok {
			this.server.mapLock.Unlock()
			this.SendMsg("当前用户名 " + newName + " 已经被使用, 不可修改\n")
		} else {
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[newName] = this
			// 更新用户名
			this.Name = newName
			this.server.mapLock.Unlock()
			this.SendMsg("您已经成功更新用户名:" + this.Name + "\n")
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		// 私聊: 消息格式 to|cheng|消息内容
		msgSlice := strings.SplitN(msg, "|", 3)
		if len(msgSlice) != 3 {
			this.SendMsg("消息格式不正确，请使用 \"to|cheng|你好呀!!!\"格式。\n")
			return
		}

		// 1. 获取对方用户名
		remoteName := msgSlice[1]
		if remoteName == "" {
			this.SendMsg("消息格式不正确，请使用 \"to|cheng|你好呀!!!\"格式。\n")
			return
		}
		// 2. 根据对方用户名得到对方 User 对象
		this.server.mapLock.RLock()
		remoteUser, ok := this.server.OnlineMap[remoteName]
		this.server.mapLock.RUnlock()
		if !ok {
			this.SendMsg("该用户名不存在\n")
			return
		}
		// 3. 获取消息内容，通知对方的 User 对象将消息发送过去
		content := msgSlice[2]
		if content == "" {
			this.SendMsg("消息内容为空, 请重发\n")
			return
		}
		remoteUser.SendMsg(this.Name + "对你说：" + content)

	} else {
		this.server.BroadCast(this, msg)
	}
}

// 监听当前 user channel 的方法，一旦有消息发送给客户端
func (this *User) ListenMessage() {
	for msg := range this.C {
		// []byte(msg) 把文字强制转换成网络能传输的字节形式，通过 conn 发送给客户端
		// 这里的“发送给客户端”，指的是发送给【当前这个 User 本身】所在的那个终端屏幕，而不是发送给其他用户。
		this.conn.Write([]byte(msg + "\n"))
	}
}
