# GitHub 发布前检查清单

这份清单用于把项目发布到 GitHub 前做最后整理。你可以按顺序检查一遍，避免把临时文件、编译产物或不清楚的说明一起提交上去。

## 1. 确认项目能运行

先启动服务端：

```bash
go run main.go server.go user.go
```

再打开另一个终端启动客户端：

```bash
go run client.go
```

如果本机没有安装 Go，会提示 `go: command not found`。这时需要先安装 Go。

## 2. 确认不要上传编译产物

下面这些文件不建议提交到 GitHub：

- `server`
- `client`
- `im-server`
- `im-client`
- `*.exe`
- `.DS_Store`

仓库已经提供 `.gitignore`，正常执行 `git add .` 时会忽略它们。

## 3. 初始化 Git 仓库

如果当前项目还不是 Git 仓库，可以执行：

```bash
git init
git add .
git commit -m "Initial commit"
```

如果已经是 Git 仓库，可以跳过这一步。

## 4. 在 GitHub 创建远程仓库

在 GitHub 上新建一个仓库，比如：

```text
IM-System
```

建议选择：

- Public：适合作为学习项目展示
- 不勾选自动生成 README：本项目已经有 README
- 不勾选自动生成 License：本项目已经有 LICENSE

## 5. 关联远程仓库

把下面的地址替换成你自己的 GitHub 仓库地址：

```bash
git remote add origin https://github.com/<your-name>/IM-System.git
git branch -M main
git push -u origin main
```

如果 GitHub 提示需要登录，可以按终端提示完成认证。

## 6. 发布后检查

推送完成后，打开 GitHub 仓库页面，重点检查：

- README 首页是否正常显示
- Mermaid 流程图是否正常渲染
- `docs/architecture.md` 链接是否能打开
- `docs/publish-checklist.md` 链接是否能打开
- 仓库里没有上传本地二进制文件和 `.DS_Store`

## 7. 推荐的第一次发布说明

可以在 GitHub Releases 或提交说明里写：

```text
v0.1.0

- 实现 TCP 服务端
- 实现命令行客户端
- 支持群聊、私聊、改名、查看在线用户
- 支持用户上线/下线广播
- 支持 300 秒空闲超时踢出
```

## 8. 后续可以继续完善

- 增加服务端启动日志
- 将服务端和客户端拆分到不同目录
- 增加 `go.mod`
- 增加单元测试
- 增加 GitHub Actions
- 增加截图或终端演示 GIF
- 增加英文 README
