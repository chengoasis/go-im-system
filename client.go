package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int // 当前 client 的模式
}

func NewClient(serverIp string, serverPort int) *Client {
	// 创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}

	// 链接 server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}

	client.conn = conn

	return client
}

// 处理 server 回应的消息，直接显示到标准输出即可
func (client *Client) DealResponse() {
	// 一旦 client.conn 有数据，就直接 copy 到stdout 标准输出上，永久阻塞监听
	io.Copy(os.Stdout, client.conn)
}

func (client *Client) menu() bool {
	var flag int

	fmt.Println("1. 公聊模式")
	fmt.Println("2. 私聊模式")
	fmt.Println("3. 更新用户名")
	fmt.Println("0. 退出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>>> 请输入合法范围内的数字 <<<<<<")
		return false
	}
}

// 查询在线用户
func (client *Client) SelectUser() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn write err:", err)
		return
	}
}

func (client *Client) PrivateChat() {
	var remoteName string
	var chatMsg string

	client.SelectUser()
	fmt.Println(">>>>>>>请输入聊天对象[用户名],exit退出")
	fmt.Scanln(&remoteName)

	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			sendMsg := "to|" + remoteName + "|" + chatMsg + "\n\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn Write err:", err)
				return
			}
		}

		chatMsg = ""
		fmt.Println(">>>>>>>请输入消息内容,exit退出")
		fmt.Scanln(&chatMsg)
	}
	return
}

func (client *Client) PublicChat() {
	// 提示用户输入消息
	var chatMsg string

	fmt.Println(">>>>>>>> 请输入聊天内容, exit退出.")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		// 发给服务器

		// 消息不为空则发送
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn Write err:", err)
				break
			}
		}
		chatMsg = ""
		// 提示用户输入消息
		fmt.Println(">>>>>>>> 请输入聊天内容, exit退出.")
		fmt.Scanln(&chatMsg)
	}
}

func (client *Client) UpdateName() bool {
	fmt.Println(">>>>>>请输入用户名:")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}
	return true
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {

		}

		// 根据不同的模式处理不同的业务
		switch client.flag {
		case 1:
			client.PublicChat()
		case 2:
			client.PrivateChat()
		case 3:
			client.UpdateName()
			break
		}
	}
}

var serverIp string
var serverPort int

// ./client -ip 127.0.0.1 -port 8888
func init() {
	// 注册一个字符串类型的命令行参数
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器 ip 地址(默认是127.0.0.1)")
	// 注册一个整数类型的命令行参数
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口(默认是8888)")
}

func main() {
	// 解析命令行参数
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>>>> 链接服务器失败...")
		return
	}

	// 单独开启一个 goroutine 去处理 server 到回执消息
	go client.DealResponse()

	fmt.Println(">>>>>>> 链接服务器成功!!!")

	// 启动客户端的业务
	client.Run()
}
