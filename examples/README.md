# JadeViewGo 基础示例

此目录是独立的 Go 模块，依赖已推送的 JadeViewGo v3 SDK 版本。不需要本地 `replace`，示例源码和资源也不会随 SDK 模块下载。

`basic/` 使用内置网页演示窗口、IPC、异步对话框、通知、YAML、剪贴板和托盘菜单。默认使用进程内回环服务，不需要签名 JAPK 即可启动。

以下命令均在 **JadeViewGo 仓库根目录** 执行。

## Windows

目标机需要 WebView2 Runtime。对应架构的原生 DLL 由 SDK 内嵌。

```powershell
go -C examples run ./basic
go -C examples build -trimpath -ldflags '-s -w -H=windowsgui' -o ./jadeview-example.exe ./basic
.\examples\jadeview-example.exe
```

去掉 `-H=windowsgui` 可以查看运行日志。

## Linux x64

Go 1.23+ 即可无 cgo 构建。运行时需要配套 JadeView SO、GTK 3、WebKitGTK 4.1 和 GIO TLS 后端。

```sh
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go -C examples build -trimpath -ldflags '-s -w' -o ./demo ./basic
JADEVIEW_LIBRARY=./lib/linux_amd64/libJadeView.so ./examples/demo
```

如果使用旧 cgo 后端，在有 C 工具链和链接依赖的环境下执行 `CGO_ENABLED=1 go -C examples build -tags jadeview_cgo -o ./demo ./basic`。分发时也可以把 `libJadeView.so` 放到可执行文件同目录。

## 综合 Playground

完整交互示例和前端放在独立的 [JadeView Playground 仓库](https://github.com/luoxueyousheng/jadeview-playground)，不属于此基础示例模块。
