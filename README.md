# JadeView Go v3 预览版

JadeView WebView 桌面库的 Go 封装。当前配套上游 **3.0.0-beta.4.26K01**，原生 ABI 为 `0x00030003`，需要 Go 1.23+。

| 平台 | 架构 | 构建与运行 |
| --- | --- | --- |
| Windows | amd64 / 386 / arm64 | 纯 Go，无需 cgo。DLL 内嵌 exe，目标机需要 WebView2 Runtime。 |
| Linux | amd64 | 默认纯 Go，支持 `CGO_ENABLED=0` 和跨编译。运行时加载 `libJadeView.so`，需要兼容的 glibc、GTK 3、WebKitGTK 4.1 等系统库。 |
| Linux | arm64 | 本预览包未交付该架构，不使用 v2 库混编。 |

Go 导入路径为 `github.com/luoxueyousheng/JadeViewGo/v3`。安装此预览版：

```sh
go get github.com/luoxueyousheng/JadeViewGo/v3@v3.0.0-beta.4
```

Go 模块标签为 `v3.0.0-beta.4`，配套的原生运行库完整版本为 `3.0.0-beta.4.26K01`。

## 运行示例

Windows：

```powershell
go run ./example

# 发布构建：缩小体积，并使用 GUI 子系统（不弹出命令提示符）
go build -trimpath -ldflags '-s -w -H=windowsgui' -o ./jadeview-example.exe ./example
```

Windows 默认将内嵌 DLL 释放到 系统临时目录下的 `jadeview/` 目录中的内容寻址子目录。若 exe 同目录有 `JadeView.dll`，优先用它；`Preload()` 会检查 ABI、位数与配置布局。未声明 Common Controls v6 的宿主会在加载时自动使用兼容激活上下文。Windows 宿主建议在 exe 清单中声明 Windows 8 或更新版本的兼容性，以支持标题栏的分层子窗口。

Linux x64（只需 Go 工具链即可构建）：

```sh
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags '-s -w' -o ./demo ./example
JADEVIEW_LIBRARY=./lib/linux_amd64/libJadeView.so ./demo
```

Linux 默认使用 purego，不需要 C 编译器和开发头文件。加载顺序为 `JADEVIEW_LIBRARY` 指定的文件、可执行文件同目录的 `libJadeView.so`、系统动态库搜索路径。分发时可将 SO 放在可执行文件旁；没有 SO 或它的链接依赖时，`Preload()` 会报告加载错误。建窗还需要完整的 WebKitGTK 运行包（含辅助进程）与 GIO TLS 后端（Ubuntu 的 `glib-networking`），不能只复制 WebKitGTK 的 SO。

需要回退旧实现时使用 `CGO_ENABLED=1 go build -tags jadeview_cgo ./example`，它仍需要 C 编译器和可用的原生链接依赖。托盘依赖桌面环境；使用 JAPK 挂载时还需 WinFsp 2.x（Windows）或 FUSE3（Linux）。

示例默认通过回环地址提供内置网页，演示窗口、IPC、异步对话框、通知、YAML、剪贴板与托盘菜单，无需签名包即可启动。

## 快速开始

```go
package main

import (
    "encoding/json"
    "log"
    "runtime"

    jadeview "github.com/luoxueyousheng/JadeViewGo/v3"
)

func main() {
    runtime.LockOSThread()
    if err := jadeview.Preload(); err != nil {
        log.Fatal(err)
    }
    jadeview.On(jadeview.EventAppReady, func(_ uint32, data string) string {
        var result struct { OK bool `json:"ok"` }
        if err := json.Unmarshal([]byte(data), &result); err != nil || !result.OK {
            log.Print(data)
            jadeview.Exit()
            return ""
        }
        if jadeview.CreateWindow("about:blank", 0, nil, nil) == 0 {
            jadeview.Exit()
        }
        return ""
    })
    jadeview.On(jadeview.EventWindowAllClosed, func(uint32, string) string {
        jadeview.Exit()
        return ""
    })
    if !jadeview.Init(false, "", "", "my-app", "com.example.myapp", false) {
        log.Fatal("JadeView 初始化未受理")
    }
    jadeview.RunMessageLoop()
}
```

## 从 v2 迁移

- `app-ready` 的来源 ID 是 `0`；解析 JSON 的 `ok:true`，不再判断来源 ID 是否为 `1`。建窗返回非零只代表已受理，操作页面前要等加载完成。
- `CreateWindow` 返回 **Window ID**，用于尺寸、标题、窗口关闭与对话框；`DefaultWebView(windowID)` 返回 **WebView ID**，用于导航、历史、缩放、脚本、IPC、页面打印、DevTools 与视图资源。IPC 回调参数为 WebView ID，调用窗口 API 前用 `WebViewParentWindow` 转换。
- `FrameStyle` 拆成创建期的 `NoFrame` 和 `TitleBarStyle`（`default`、`hidden`、`BarOverlay`）；运行中可用 `SetWindowFrame` 切换标题栏，不改变边框。标题栏配色迁移至 `SetTitlebarStyle`。
- WebView 配置改用 `PostMessageAllowedOrigins`、`CORSAllowedOrigins`、`PasswordAutofill`、`FormAutofill`；创建接口将 `nil` 原样传给原生库，采用原生默认值。
- `ExecuteJavaScript` 返回 64 位 `RequestID`，不再发出 `javascript-result`。用 `GetRequestResult` 查询，完成后调用 `ReleaseRequest`。
- 托盘右键监听 `EventTrayEvent` 后用 `PopupMenu` 弹出菜单，在 `EventMenuCommand` 中读取字符串 `item_id`；网页右键模板用 `SetContextMenu`。
- 普通命名事件仍用 `On` / `Off`；已进入的回调可在 `Off` 返回后继续执行，Go 在原生确认回调结束后才复用槽位。

本次升级迁移现有绑定及核心 v3 功能，不代表头文件中的每项新增接口都已提供 Go 封装。完整声明见 [JadeView.h](include/JadeView.h)，版本变更见 [CHANGELOG](CHANGELOG.md)。
