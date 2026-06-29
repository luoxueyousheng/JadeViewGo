# JadeView Go 封装（跨平台：Windows + Linux）

JadeView WebView 库的 Go(cgo) 封装。

| 平台 | 架构 | 链接方式 | 库文件 | 运行时还需 |
|------|------|----------|--------|-----------|
| Linux | amd64 | 静态 | `lib/linux_amd64/libJadeView.a` | 系统 GTK3/WebKit2GTK |
| Linux | arm64 | 静态 | `lib/linux_arm64/libJadeView.a` | 系统 GTK3/WebKit2GTK |
| Windows | amd64 (x64) | **延迟链 DLL**(MSVC) | `lib/windows_amd64/JadeView.dll` | 仅 WebView2 Runtime（系统级） |
| Windows | 386 (x86) | **延迟链 DLL**(MSVC) | `lib/windows_386/JadeView.dll` | 仅 WebView2 Runtime（需 32 位 MinGW 构建） |

> Windows 用官方 **MSVC 版 `JadeView.dll`**（自包含 WebView2Loader）。通过 **延迟加载 + go:embed**
> 把 DLL 内置进二进制、运行时自动释放到临时目录——**真正的单 exe 分发，无需携带任何 DLL**
> （详见「Windows 自包含」）。唯一外部依赖是系统级 Edge WebView2 Runtime（Win11 自带，任何方案都去不掉）。
>
> （早期曾用 GNU 静态库方案，因作者下版本不再发 GNU 产物，已弃用，改为对 MSVC DLL 做同样的内置+释放。）

## 安装

```bash
go get github.com/luoxueyousheng/JadeViewGo
```

cgo 项目，**构建需 C 编译器**（Windows 装 MinGW-w64，Linux 装 gcc + GTK3/WebKit2GTK，见「前置依赖」）。

## 最小用法

```go
package main

import jadeview "github.com/luoxueyousheng/JadeViewGo"

func main() {
    // 必须在 Init 之前注册 app-ready，并在其回调里建窗
    jadeview.On("app-ready", func(windowID uint32, data string) string {
        if windowID == 1 {
            opts := jadeview.DefaultWindowOptions()
            opts.Title = "Hello JadeView"
            jadeview.CreateWindow("https://example.com", 0, &opts, nil)
        }
        return ""
    })
    jadeview.On("window-all-closed", func(uint32, string) string {
        jadeview.Exit()
        return ""
    })
    jadeview.Init(true, "", "", "my-app", "my-app-signature", false)
    jadeview.RunMessageLoop()
}
```

完整可交互示例见 [`example/`](example/)。

## Windows 自包含（DLL 内置 + 运行时释放）

DLL 默认是**加载期导入**——进程启动前 Windows 就要它存在，没法在进程内自解压。
本封装通过两步把 `JadeView.dll` 变成「内置 + 自动释放」：

1. **延迟加载**：用 `dlltool` 从 DLL 导出表生成延迟导入库 `lib/windows_amd64/libJadeView_delay.a`
   替代普通导入（`jadeview_windows_amd64.go` 链接 `-lJadeView_delay -ldelayimp`），
   使 DLL 改为「首次调用其函数时才 LoadLibrary」。这是纯动态链接，不碰 MSVC 静态库的 CRT 符号问题。
2. **go:embed + 释放**：`jadeview_windows_embed.go` 把 `JadeView.dll` 编进二进制，`init()` 启动时
   释放到 `%TEMP%\jadeview\` 并按全路径预加载，之后延迟加载命中它。

> **换新版 DLL** 后需重新生成延迟库（需 MinGW 的 `objdump`/`dlltool` 在 PATH）：
> 把新的 MSVC `JadeView.dll` 覆盖到 `lib/windows_amd64/`，然后从导出表生成 `.def` 再做延迟库——
> ```bash
> cd lib/windows_amd64
> # 1) 从 DLL 导出表抽取函数名，写成 JadeView.def（LIBRARY 行 + EXPORTS + 每行一个函数名）
> objdump -p JadeView.dll | grep '+base\[' | awk '{print $NF}' | sort -u > exports.txt
> { echo "LIBRARY JadeView.dll"; echo "EXPORTS"; cat exports.txt; } > JadeView.def
> # 2) 生成延迟导入库
> dlltool --input-def JadeView.def --output-delaylib libJadeView_delay.a --dllname JadeView.dll
> ```

## 目录结构

```
JadeView/
├── include/JadeView.h            # C 头文件（作者官方版，跨编译器一致）
├── lib/
│   ├── linux_amd64/libJadeView.{a,so}
│   ├── linux_arm64/libJadeView.{a,so}
│   ├── windows_amd64/            # x64（已启用、已验证）
│   │   ├── JadeView.dll          # MSVC 版（自包含 WebView2Loader），被 go:embed 内置
│   │   └── libJadeView_delay.a   # dlltool 生成的延迟导入库（链接用）
│   └── windows_386/              # x86（已启用，需 32 位 MinGW 构建）
│       ├── JadeView.dll
│       └── libJadeView_delay.a   # 32 位延迟库（@N 修饰 + --kill-at）
├── beta/                         # 官方版本/API 文档
├── go.mod
├── jadeview.go                   # 公共 API + C 前导（所有平台）
├── jadeview_linux_amd64.go       # Linux x64 链接配置（静态）
├── jadeview_linux_arm64.go       # Linux arm64 链接配置（静态）
├── jadeview_windows_amd64.go     # Windows x64 链接配置（延迟链 DLL）
├── jadeview_windows_embed.go     # Windows x64 内置 DLL + 运行时释放
├── example/                      # 可交互完整示例
│   ├── main.go                   #   窗口+事件+IPC+异步对话框+托盘+YAML
│   ├── index.html                #   前端页面（jade.invoke / jade.on，go:embed 内置）
│   └── icon.ico                  #   托盘图标（go:embed 内置）
└── README.md
```

## 示例（example/）

`go run ./example` 启动一个可交互窗口（Windows 需先装好 MinGW，见下）。它把封装的各模块串了起来，
便于手动验证 GUI 交互——尤其是**异步对话框的回调桥**：

- 页面按钮 → `jade.invoke(channel)` → Go 的 `RegisterIPCHandler` 处理 → 返回值回显到页面日志区。
- **打开/保存文件、消息框**走异步对话框：弹框后立即返回，用户操作完的结果经 `SendIPCMessage`
  推回页面的 `jade.on('dialog-result')`——点一下就能确认回调桥工作正常。
- YAML 写入/读取、置顶切换、最小化、DevTools、显示器信息、系统通知等均可点按验证。
- **托盘**：右键托盘图标可显示/隐藏窗口、退出（验证 `tray-menu-command` 事件）。

页面与图标用 `go:embed` 内置，运行时释放到 `%TEMP%\jadeview-demo\`，并通过自定义协议服务加载
（URL 形如 `JADE://<app_signature>/index.html`）。

## 前置依赖：必须有 C 编译器（cgo）

cgo 在任何平台都需要 C 编译器。当前这台机器 **没有 gcc**，`go build` 会报
`cgo: C compiler "gcc" not found`。

- **Windows**：装 MinGW-w64（推荐用 [MSYS2](https://www.msys2.org/) 的
  `mingw-w64-x86_64-gcc`，或 TDM-GCC），把 `bin` 目录加入 `PATH`，使 `gcc` 可用。
- **Linux**：`sudo apt install build-essential pkg-config libgtk-3-dev libwebkit2gtk-4.1-dev`
  （JadeView 静态库依赖 GTK3 / WebKit2GTK；若只有 webkit 4.0，把
  `jadeview_linux_*.go` 里的 `webkit2gtk-4.1` 改成 `webkit2gtk-4.0`）。

## 构建 / 运行

```bash
cd JadeView
CGO_ENABLED=1 go build ./...
go run ./example
```

**Windows 运行时**：`JadeView.dll` 已内置进 exe，**单文件分发，无需携带任何 DLL**
（首次运行自动释放到 `%TEMP%\jadeview\`）。目标机仅需 Edge WebView2 Runtime（Win11 自带）。
GUI 程序（无控制台黑窗）构建加 `-ldflags "-H windowsgui"`。
（Linux 走静态链接，无额外 DLL，依赖系统 GTK3/WebKit2GTK。）

## 当前进度

已封装并在 Windows 验证通过：

- **生命周期**：`Init` / `Version` / `RunMessageLoop` / `Exit`
- **窗口管理**：`CreateWindow`（含 `WindowOptions`/`WebViewSettings` 结构体）、
  `CreateBorderlessWindow`、`Navigate` / `Reload` / `ExecuteJavaScript` /
  `SetTitle` / `SetSize` / `SetPosition` / `SetVisible` / `SetFocus` /
  `SetAlwaysOnTop` / `Close` / `Minimize` / `ToggleMaximize` / `IsMaximized` /
  `WindowCount`
- **事件桥接**：`On` / `Off` / `RegisterIPCHandler`（C 跳板槽位池，上限
  `MaxEventHandlers`=64）

- **对话框/通知**（dialog.go）：`ShowNotification`、`ShowOpenDialog`/`ShowSaveDialog`、
  `ShowMessageBox`、`ShowErrorBox`，以及右键菜单 `MenuItemCreate`/`SetContextMenuItems` 等
- **异步对话框**（dialog_async.go）：`ShowOpenDialogAsync`/`ShowSaveDialogAsync`/
  `ShowMessageBoxAsync`（回调桥，上限 `MaxAsyncDialogs`=16）
- **托盘**（tray.go）：`TrayCreate`/`TraySetVisible`/`TraySetTooltip`/`TraySetIconFromFile`/
  `TraySetIconFromData`/`TraySetMenu`（扁平表 `TrayMenuItem`）
- **YAML 存储**（yaml.go）：`YAMLSet`/`YAMLGet`/`YAMLGetAll`/`YAMLKeys`/`YAMLHas`/
  `YAMLDelete`/`YAMLLen`/`YAMLClear` 等（getter 用缓冲区两阶段，避开 CoTaskMemFree）
- **JAPK 资源包**（japk.go）：`SetPublicKey`/`LoadFromBytes`/`IsLoaded`/`GetAppSignature`/
  `GetSignatureInfo`/`Unload`
- **系统工具**（system.go）：剪贴板、`GetPath`/`GetLocale`/`GetDisplaysInfo`、打印、全局热键、
  开机自启、URL 协议/文件关联、安全资源、`SmartConvertEncoding`、`NTPNow` 等
- **窗口扩展**（window_ext.go）：状态查询、devtools、IPC 发送、层级/背景/全屏、任务栏进度等

> **覆盖度**：头文件 123 个导出函数已全部封装，仅 2 个有意未封——`cleanup_all_windows`（已废弃，
> 用 `Exit`）、`yaml_get_str`（用 `CoTaskMemFree`，跨平台不可移植，已用缓冲区版 `YAMLGet` 替代）。
> YAML、剪贴板、显示器、版本、JAPK 等已在 Windows 实跑验证（YAML 需在 app-ready 之后调用）。

## 已知问题 / 注意事项

- **当前全部为 beta.8.26F13**：Windows DLL 与 Linux 库版本已统一。
- **Windows 支持 amd64(x64) 与 386(x86)**：均用 MSVC 版 `JadeView.dll`（官方发行物，自包含 WebView2Loader）。
- **DLL 已内置**：通过延迟加载 + go:embed 编进二进制、运行时自动释放，无需随 exe 携带任何 DLL。
- **换 DLL 后重生成延迟库**：覆盖 `lib/windows_amd64/JadeView.dll` 后，按「Windows 自包含」一节里的
  `objdump`+`dlltool` 两条命令重新生成 `libJadeView_delay.a`；若链接报 `undefined reference to <某导出>`，多半是没重生成。
- **临时释放注意**：首次运行写 `%TEMP%\jadeview\JadeView.dll`（已做大小校验，不重复写）；
  个别杀软可能对「释放 DLL 并加载」有启发式告警，属正常。
- **Windows x86(386) 构建需 32 位工具链**：装 MSYS2 的 `mingw-w64-i686-toolchain`（同一个 MSYS2 里
  `pacman -S mingw-w64-i686-toolchain`，装到 `C:\msys64\mingw32`），构建时：
  ```powershell
  $env:PATH = "C:\msys64\mingw32\bin;" + $env:PATH
  $env:CGO_ENABLED="1"; $env:GOARCH="386"; $env:CC="C:\msys64\mingw32\bin\gcc.exe"
  go build ./example
  ```
  32 位 `__stdcall` 有 `@N` 名字修饰：延迟库的 `.def` 用 `name@N`（N=参数字节数，从头文件算）+
  `dlltool --kill-at`，使链接器符号 `_name@N` 与 DLL 无修饰导出对上。换 DLL 后需照此重生成。
- **Windows arm64 未启用**：官方有 MSVC DLL，但需 `aarch64-w64-mingw32` 工具链（本机 binutils 不认
  ARM64 PE），暂未纳入；按 x86 同法可扩。

## 待办（封装完整 API）

- [x] 窗口管理（`create_webview_window` 等，含结构体映射）
- [x] 事件回调桥接（`jade_on` / `register_ipc_handler`，C 跳板槽位池）
- [x] 对话框、通知、托盘菜单、右键菜单
- [x] YAML store、system 工具函数、窗口扩展
- [x] 返回 `char*` 的函数用 `goStringFree`（jade_text_free）释放
- [x] 异步对话框（`*_async` 回调版）、JAPK 资源包、打印
- [x] **全部 123 个导出函数已封装**（仅 2 个有意未封，见上）
- [ ] 单元测试 / Linux 端实跑验证 / 异步对话框回调的 GUI 实测
```

## 许可证

本 Go 封装层代码以 [MIT](LICENSE) 许可证发布。

> 注意：`lib/` 下的 JadeView 库二进制（`.dll`/`.a`/`.so`）是 JadeView 作者的第三方产物，
> **不在本 MIT 许可范围内**，其使用/再分发受 JadeView 自身条款约束。
