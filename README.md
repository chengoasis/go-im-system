# Go IM System - Golang TCP 即时通讯系统

[![Go](https://img.shields.io/badge/Go-1.22%2B-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Status](https://img.shields.io/badge/status-learning%20project-F59E0B)](#)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

一个基于 Go 标准库实现的轻量级即时通讯系统。项目不依赖第三方库，核心代码围绕 TCP 连接、goroutine、channel、锁和命令行交互展开，适合 Go 网络编程入门练习。

## 项目出处

本项目学习自 B 站课程：[8小时转职Golang工程师](https://www.bilibili.com/video/BV1gf4y1r79E/?spm_id_from=333.337.search-card.all.click&vd_source=e2695651ecce1cb03c26637e203cd31f)。

在课程代码基础上，本仓库补充和整理了：

- 更完整的新手向 README
- 架构说明和 Mermaid 流程图
- 部分边界情况和并发稳定性处理

## 项目简介

你可以把它理解成一个命令行聊天室：

- 启动一个服务端
- 打开多个客户端连接服务端
- 多个用户之间可以群聊、私聊、改名、查看在线用户
- 用户长时间不发消息会被服务端自动踢下线

用尽量少的代码把 IM 系统的基础流程跑通。

## 功能列表

| 功能 | 说明 |
| --- | --- |
| TCP 服务端 | 服务端默认监听 `127.0.0.1:8888` |
| 多客户端连接 | 每个客户端连接都会交给独立 goroutine 处理 |
| 群聊 | 普通文本会广播给所有在线用户 |
| 私聊 | 使用 `to|用户名|消息内容` 给指定用户发消息 |
| 查看在线用户 | 使用 `who` 查看当前在线用户 |
| 修改用户名 | 使用 `rename|新用户名` 修改自己的名字 |
| 上线/下线广播 | 用户进入和离开聊天室时会广播通知 |
| 空闲超时 | 300 秒不发送消息会被服务端踢下线 |
| 并发保护 | 使用 `sync.RWMutex` 保护在线用户列表 |
| 下线清理 | 使用 `sync.Once` 避免重复下线、重复关闭资源 |

## 项目结构

```text
.
├── main.go                  # 服务端入口：创建 Server 并启动
├── server.go                # 服务端核心：监听连接、广播消息、连接生命周期、超时控制
├── user.go                  # 用户逻辑：上线、下线、群聊、私聊、改名、消息写回
├── client.go                # 命令行客户端：连接服务端并提供交互菜单
├── docs/
│   └── architecture.md      # 架构说明和核心流程图
├── .gitignore               # Git 忽略规则
├── LICENSE                  # 开源许可证
└── README.md                # 项目首页文档
```

## 环境准备

你需要先安装 Go。

安装完成后，在终端执行：

```bash
go version
```

如果能看到类似下面的输出，说明 Go 已经安装成功：

```text
go version go1.22.x linux/amd64
```

如果提示：

```text
go: command not found
```

说明 Go 还没有安装，或者环境变量没有配置好。可以先去 [Go 官方下载页](https://go.dev/dl/) 安装。

## 快速开始

### 1. 进入项目目录

如果你已经把项目放在本地或虚拟机里，直接进入项目目录：

```bash
cd IM-System
```

如果是从 GitHub 克隆：

```bash
git clone <your-repo-url>
cd IM-System
```

### 2. 启动服务端

打开第一个终端，执行：

```bash
go run main.go server.go user.go
```

服务端默认监听：

```text
127.0.0.1:8888
```

当前版本服务端启动成功后不会主动打印“启动成功”，这是正常的。保持这个终端不要关闭。

### 3. 启动第一个客户端

打开第二个终端，执行：

```bash
go run client.go
```

看到菜单就表示连接成功：

```text
>>>>>>> 链接服务器成功!!!
1. 公聊模式
2. 私聊模式
3. 更新用户名
0. 退出
```

### 4. 再启动一个客户端

打开第三个终端，再执行一次：

```bash
go run client.go
```

这样就有两个用户在线，可以开始测试群聊、私聊和改名。

## 客户端使用方式

客户端启动后会显示菜单：

```text
1. 公聊模式
2. 私聊模式
3. 更新用户名
0. 退出
```

### 修改用户名

输入：

```text
3
```

然后输入新用户名：

```text
cheng
```

修改成功后，后续聊天会显示新用户名。

### 群聊

输入：

```text
1
```

然后输入聊天内容：

```text
hello
```

其他在线用户会收到类似消息：

```text
[127.0.0.1:xxxxx]cheng:hello
```

在公聊模式里输入：

```text
exit
```

会退出公聊模式，回到主菜单。

### 私聊

建议先把两个客户端分别改名，例如：

```text
cheng
llo
```

在 `cheng` 的客户端菜单里输入：

```text
2
```

客户端会先查询在线用户，然后提示你输入聊天对象。输入：

```text
llo
```

再输入消息内容：

```text
hello
```

`llo` 会收到来自 `cheng` 的私聊消息。

在私聊消息输入处输入：

```text
exit
```

会退出私聊模式，回到主菜单。

### 退出客户端

在主菜单输入：

```text
0
```

## 使用 nc 测试

如果不想使用项目自带的 `client.go`，也可以使用 `nc` 直接连接服务端。

先启动服务端：

```bash
go run main.go server.go user.go
```

再打开两个新终端，分别执行：

```bash
nc 127.0.0.1 8888
```

连接后可以直接输入协议命令：

```text
who
rename|cheng
to|llo|hello
大家好
```

## 编译成可执行文件

因为服务端和客户端都有自己的 `main()` 函数，所以要分别指定文件编译。

编译服务端：

```bash
go build -o im-server main.go server.go user.go
```

编译客户端：

```bash
go build -o im-client client.go
```

运行：

```bash
./im-server
./im-client
```

Windows 可以编译成：

```bash
go build -o im-server.exe main.go server.go user.go
go build -o im-client.exe client.go
```

## 虚拟机运行说明

如果服务端和客户端都在同一台虚拟机里运行，默认的 `127.0.0.1:8888` 可以直接使用。

如果服务端运行在虚拟机里，客户端运行在宿主机上，需要注意：

- `127.0.0.1` 只代表当前机器自己
- 宿主机访问虚拟机服务时，通常不能直接连虚拟机里的 `127.0.0.1`
- 可以把服务端监听地址改成虚拟机网卡 IP 或 `0.0.0.0`
- 还需要确认虚拟机网络模式、端口转发和防火墙规则

服务端监听地址在 `main.go`：

```go
server := NewServer("127.0.0.1", 8888)
```

## 新手常见问题

### 为什么不能直接执行 `go run .`？

因为当前目录下有两个入口函数：

- `main.go` 里的 `main()` 是服务端入口
- `client.go` 里的 `main()` 是客户端入口

如果直接执行：

```bash
go run .
```

Go 会同时编译当前目录所有 `.go` 文件，就会出现 `main redeclared` 之类的错误。

正确方式是分别运行：

```bash
go run main.go server.go user.go
go run client.go
```

## 核心设计思路

这个项目可以分成几个模块理解：

| 模块 | 作用 |
| --- | --- |
| 服务端启动 | 使用 `net.Listen` 监听 TCP 端口 |
| 连接接入 | 使用 `Accept` 接收客户端连接 |
| 连接处理 | 每个客户端连接使用一个 goroutine 单独处理 |
| 用户管理 | 使用 `OnlineMap` 保存当前在线用户 |
| 并发保护 | 使用 `sync.RWMutex` 保护 `OnlineMap` |
| 消息广播 | 使用 `Server.Message` channel 分发群聊消息 |
| 用户收件箱 | 每个 `User` 都有自己的 `C` channel |
| 下线清理 | 使用 `sync.Once` 保证资源只清理一次 |

更详细的架构图和流程拆解可以看：[docs/architecture.md](docs/architecture.md)。

## 建议阅读顺序

如果你是第一次看这个项目，推荐按下面的顺序读代码：

1. `main.go`：看服务端是怎么启动的
2. `server.go`：看服务端如何监听连接、处理连接、广播消息、控制超时
3. `user.go`：看用户上线、下线、改名、私聊和群聊逻辑
4. `client.go`：看客户端如何连接服务端并发送命令
5. `docs/architecture.md`：结合流程图再整体复盘一遍

## 当前版本已处理的细节

- 私聊命令使用 `strings.SplitN` 解析，避免格式不完整时直接 panic
- `OnlineMap` 的读写使用锁保护，减少并发读写 map 的风险
- 用户下线使用 `sync.Once`，避免重复广播下线和重复关闭资源
- 用户自己的消息 channel 关闭后，`User.ListenMessage()` 会自然退出
- 客户端断开后，`Handler` 会通过 `isQuit` 及时退出
- 活跃通知 `isLive` 使用带缓冲 channel，避免读消息 goroutine 被活跃通知卡住

## 后续可以扩展的方向

- 把 IP、端口和超时时间做成配置项
- 给服务端增加启动成功日志
- 将服务端和客户端拆到不同目录，避免 `go run .` 报错
- 增加登录、注册、密码校验
- 增加离线消息
- 增加聊天室或群组
- 增加消息持久化，比如保存到 MySQL 或 Redis
- 增加单元测试

## License

本项目使用 [MIT License](LICENSE)。
