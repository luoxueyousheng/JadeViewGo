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

## beta.4 API 覆盖

本节描述当前源码的覆盖范围；已发布的 `v3.0.0-beta.4` 标签不含以下新增封装。Go 封装会在边界上复制原生数据并释放原生内存，不需要宿主手动释放原生指针。

| 功能 | Go 接口 |
| --- | --- |
| ABI 查询 | `GetABIInfo` |
| 结构化事件 | `OnAppEvent`、`OnWindowEvent`、`OnWebViewEvent`、`UnsubscribeEvent` |
| 事件快照 | 回调收到 Go 拥有的 `TypedEvent`；可保留或调用 `Clone`，无需 `event_clone/event_free` |
| JAPK 注册与挂载 | `RegisterJAPK`、`RegisterLoadedJAPK`、`MountJAPK`、`GetJAPKMountCapabilities`、`GetJAPKMountRoot`、`GetJAPKMountState`、`UnmountJAPK`、`UnregisterJAPK` |
| 资源交付 | `RegisterResourceEx`、`GetResourceCapabilities`、`GetResourceDeliveryMode`、`GetResourceDeliveryInfo` |
| 错误详情 | `LastError`、`LastErrorCode`、`LastErrorMessage`；原生错误按 OS 线程保存，失败调用和读取错误应在同一个 `runtime.LockOSThread` 区域 |
| 菜单与存储 | `ClearContextMenu`、`YAMLGetStr`。Linux 的 `YAMLGetStr` 使用等价缓冲区接口，避免依赖 Windows COM 分配器 |

结构化订阅的注销返回 `RequestID`，调用者用 `GetRequestResult` 查询并用 `ReleaseRequest` 释放；不要在回调内部等待自身注销。命名事件优先用 `Off`；`OffImmediate` 无完成请求，为免迟到的原生回调误投，隔离原槽位直到进程退出。

原生接口有平台限制：原生文件交付模式仅 Windows 提供；真正的 JAPK 系统挂载仍需 WinFsp 或 FUSE3。完整原生声明见 [JadeView.h](include/JadeView.h)，版本变更见 [CHANGELOG](CHANGELOG.md)。
