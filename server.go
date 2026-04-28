package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	// 在线用户列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	// 消息广播channel
	Message chan string
}

// 创建一个实例 (构造函数)
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}

	return server
}

// 监听 Message 广播消息 channel 的goroutine，一旦有消息就发送给全部在线User
func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message

		// 将 msg 发送给全部在线 User
		this.mapLock.RLock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.RUnlock()
	}
}

// 广播消息的方法
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	this.Message <- sendMsg
}

func (this *Server) Handler(conn net.Conn) {
	// ...当前链接的业务
	// fmt.Println("链接建立成功")

	user := NewUser(conn, this)

	user.Online()

	// 监听用户是否活跃的 channel
	isLive := make(chan bool, 1)

	// 监听用户是否已经断开连接
	isQuit := make(chan bool, 1)
	// 定义一个函数，用于发送退出信号
	notifyQuit := func() {
		select {
		case isQuit <- true: // 如果可以发送（channel未满），执行这个
		default: // 如果无法发送（channel已满），立即执行这个
		}
	}

	// 接受客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			// 从连接中读取数据到 buf
			// n: 实际读取到的字节数（0 到 len(buf) 之间）
			// err: 错误信息，如果没有错误则为 nil
			n, err := conn.Read(buf)

			if err != nil {
				if err == io.EOF {
					fmt.Println("客户端正常断开")
				} else {
					fmt.Println("读取错误: ", err)
				}
				user.offline()
				notifyQuit()
				return
			}

			if n == 0 {
				// 理论上不应该在没有错误的情况下 n==0
				// 但为了安全，也当作断开处理
				user.offline()
				notifyQuit()
				return
			}

			// 提取用户消息（去除'\n'）
			msg := string(buf[:n-1])

			user.DoMessage(msg)

			// 用户的任意消息，代表当前用户是一个活跃的
			select {
			case isLive <- true:
			default:
			}
		}
	}()

	// 当前 handler 阻塞
	for {
		select {
		case <-isLive:
			// 说明当前用户是活跃的，重置定时器
			// 不做任何事情，激活 select 更新下面的定时器
		case <-isQuit:
			// 当前用户已经断开连接，退出 Handler
			return
		case <-time.After(time.Second * 300):
			// 已经超时
			// 将当前 User 强制关闭
			user.SendMsg("超时了...你被踢了!")
			// 销毁资源并广播下线
			user.offline()
			// 退出当前 Handle
			return
		}
	}
}

// 启动服务器的接口
func (this *Server) Start() {
	// socket listen
	// fmt.Sprintf：这是一个字符串拼接工具，把 this.Ip 和 this.Port 拼凑成网络需要的格式（如 "127.0.0.1:8080"）。
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	// close listen socket
	defer listener.Close()

	// 启动监听 Message 的 goroutine ———— 一个专门负责广播的 goroutine
	go this.ListenMessage()

	for {
		// accept
		// conn 是一个双向管道，它同时包含了客户端和服务器的信息，通过不同的方法获取两端的信息。
		// - 本地：127.0.0.1:8888（服务器）
		// - 远程：192.168.1.100:54321（客户端）
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}

		// do handler
		go this.Handler(conn)
	}

}
