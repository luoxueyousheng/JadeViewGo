# JadeView Go 封装

[JadeView](https://jadeview.com) WebView 桌面库的 Go(cgo) 封装。用 Go + HTML/CSS/JS 写跨平台桌面应用：窗口、事件、双向 IPC、托盘、对话框、通知、YAML 持久化、NTP 授时等一应俱全。

当前对应上游 **v2.3.0-beta.10**，头文件 124 个导出函数已全部封装。

## 支持平台

| 平台 | 架构 | 链接方式 | 库文件 | 目标机运行时依赖 |
|------|------|----------|--------|-----------------|
| Windows | amd64 (x64) | 延迟链 DLL（MSVC）+ go:embed | `lib/windows_amd64/` | 仅 Edge WebView2 Runtime（Win11 自带） |
| Windows | 386 (x86) | 延迟链 DLL（MSVC）+ go:embed | `lib/windows_386/` | 同上 |
| Windows | arm64 | 常规链 DLL（需随 exe 携带） | `lib/windows_arm64/` | WebView2 Runtime + `JadeView.dll` |
| Linux | amd64 | 静态链接 `.a` | `lib/linux_amd64/` | 系统 GTK3 / WebKit2GTK |
| Linux | arm64 | 静态链接 `.a` | `lib/linux_arm64/` | 系统 GTK3 / WebKit2GTK |

> **Windows x64/x86 单 exe 分发**：`JadeView.dll` 通过 go:embed 编进二进制、运行时自动释放到
> `%TEMP%\jadeview\` 并预加载（延迟加载机制，详见「Windows 自包含原理」）。
> **无需随程序携带任何 DLL**。Windows arm64 例外：工具链不支持 arm64 延迟库，
> 需随 exe 携带 DLL（见其小节）。
>
> **Linux 无额外文件**：静态链接进二进制，只依赖系统包管理器安装的 GTK3/WebKit2GTK。

## 安装

```bash
go get github.com/luoxueyousheng/JadeViewGo
```

这是 cgo 项目，**构建机器必须有 C 编译器**，各平台前置依赖见下一节。

## 各平台使用方式

### Windows x64

**前置依赖**（仅构建机需要）：

1. 安装 [MSYS2](https://www.msys2.org/)，然后装 64 位 GCC 工具链：
   ```bash
   pacman -S mingw-w64-x86_64-toolchain
   ```
2. 把 `C:\msys64\mingw64\bin` 加入 `PATH`（使 `gcc` 可用）。

**构建 / 运行**（PowerShell）：

```powershell
$env:PATH = "C:\msys64\mingw64\bin;" + $env:PATH
$env:CGO_ENABLED = "1"
go run ./example                # 直接运行示例
go build -o myapp.exe .         # 构建你的应用
go build -ldflags "-H windowsgui" -o myapp.exe .   # GUI 程序（无控制台黑窗）
```

产物是**单个 exe**：DLL 已内置，首次运行自动释放到 `%TEMP%\jadeview\`（按字节校验，内容一致不重复写）。目标机只需系统装有 Edge WebView2 Runtime（Win11 自带；Win10 可用微软的 Evergreen Bootstrapper 安装）。

### Windows x86（32 位）

**前置依赖**：同一个 MSYS2 里再装 32 位工具链：

```bash
pacman -S mingw-w64-i686-toolchain
```

**构建**（PowerShell，注意 `GOARCH=386` 和 32 位 gcc）：

```powershell
$env:PATH = "C:\msys64\mingw32\bin;" + $env:PATH
$env:CGO_ENABLED = "1"
$env:GOARCH = "386"
$env:CC = "C:\msys64\mingw32\bin\gcc.exe"
go build -o myapp_x86.exe ./example
```

> 32 位 `__stdcall` 有 `@N` 名字修饰，延迟库已按此生成，构建时无需额外处理；
> 仅当**换新版 DLL** 时需重新生成延迟库（见「升级上游库」）。

### Windows arm64

**与 x64/x86 的区别**：binutils 的 `dlltool` 没有 aarch64-PE 延迟库支持（MSYS2 的
`dlltool -m arm64` 生成延迟库会**静默失败**），llvm-dlltool 也不支持 delaylib。
因此 arm64 走**常规加载期链接**——没有单 exe 自包含，**分发时 exe 旁边必须携带
`lib/windows_arm64/JadeView.dll`**。

**前置依赖**：[llvm-mingw](https://github.com/mstorsjo/llvm-mingw/releases)
（下载 `llvm-mingw-*-ucrt-x86_64.zip`，解压即用，可在 x64 机器上交叉编译 arm64）。
MSYS2 没有 x64→arm64 的交叉 GCC，MinGW-GCC 编不了。

**构建**（PowerShell，x64 机器交叉编译）：

```powershell
$env:PATH = "C:\llvm-mingw\bin;" + $env:PATH
$env:CGO_ENABLED = "1"
$env:GOARCH = "arm64"
$env:CC = "aarch64-w64-mingw32-clang"
go build -ldflags "-H windowsgui" -o myapp_arm64.exe .
```

链接器（lld）直接解析 `JadeView.dll` 导出表完成链接，无需导入库。
`lib/windows_arm64/JadeView.def` 是导出表备份，如需传统导入库：
`llvm-dlltool -m arm64 -d JadeView.def -D JadeView.dll -l libJadeView.dll.a`。

### Linux amd64

**前置依赖**（Debian/Ubuntu 系）：

```bash
sudo apt install build-essential pkg-config libgtk-3-dev libwebkit2gtk-4.1-dev
```

> 老发行版若只有 WebKit2GTK 4.0（没有 4.1 包），把
> `jadeview_linux_amd64.go` / `jadeview_linux_arm64.go` 里的
> `webkit2gtk-4.1` 改成 `webkit2gtk-4.0` 即可。

**构建 / 运行**：

```bash
CGO_ENABLED=1 go build ./...
go run ./example
```

`libJadeView.a` 静态链接进二进制；GTK3/WebKit2GTK 走系统动态库，目标机用包管理器装运行时包即可：

```bash
# 目标机（运行时，非 -dev 包）
sudo apt install libgtk-3-0 libwebkit2gtk-4.1-0
```

### Linux arm64

与 amd64 完全相同，库文件在 `lib/linux_arm64/`。**推荐在 arm64 机器上原生构建**（树莓派、云 ARM 实例等）：

```bash
sudo apt install build-essential pkg-config libgtk-3-dev libwebkit2gtk-4.1-dev
CGO_ENABLED=1 go build ./example
```

交叉编译需要 `aarch64-linux-gnu-gcc` 工具链**加上 arm64 版 GTK/WebKit 的 sysroot**，配置繁琐，一般不值得——原生构建更省事。

## 快速开始

```go
package main

import jadeview "github.com/luoxueyousheng/JadeViewGo"

func main() {
    // 关键时序：app-ready 必须在 Init 之前注册，并在回调里判断 windowID==1 再建窗
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
    // Init(开发模式, 日志路径, 数据目录, 应用名, 应用签名≥6字符, 单实例)
    jadeview.Init(true, "", "", "my-app", "my-app-signature", false)
    jadeview.RunMessageLoop() // 阻塞直到退出
}
```

## 示例（example/）

一份代码，Windows / Linux 都能跑（平台差异用 `runtime.GOOS` 分支）：

```bash
go run ./example        # Linux 直接跑；Windows 先按上面配好 MinGW
```

演示内容：

- **IPC 双向通信**：页面按钮 → `jade.invoke(channel)` → Go `RegisterIPCHandler` → 返回值回显；Go 侧 `SendIPCMessage` → 页面 `jade.on('dialog-result')`。
- **异步对话框**：打开/保存/消息框，弹框立即返回、结果回调推给页面。
- **托盘**：右键菜单显示/隐藏窗口、退出（Linux 依赖桌面环境托盘支持，失败自动跳过）。
- **YAML 持久化**：写入/读取/全量读取（存于 `Init` 的数据目录）。
- **系统能力**：剪贴板读写、NTP 网络时间、显示器信息、系统通知、窗口边界、HWND⇄窗口ID 互查、任务栏闪烁、DevTools。
- **`jade-region-drag` 拖动区**（Windows）：HTML 属性即可拖动窗口/双击最大化，无右键系统菜单。

页面与图标用 `go:embed` 内置，运行时释放到 `%TEMP%/jadeview-demo/`，经自定义协议服务加载（URL 形如 `JADE://<app_signature>/index.html`）。

## API 总览

| 模块 | 文件 | 主要函数 |
|------|------|----------|
| 生命周期 | `jadeview.go` | `Init` / `Version` / `RunMessageLoop` / `Exit` |
| 窗口创建 | `window.go` | `CreateWindow`（`WindowOptions`/`WebViewSettings`）、`CreateBorderlessWindow`、`Navigate`、`ExecuteJavaScript`、`SetTitle/SetSize/SetPosition/...` |
| 窗口扩展 | `window_ext.go` | 状态查询 `Is*`、`GetWindowBounds`、`GetWindowHWND`⇄`GetWindowID`、层级/背景/全屏/主题、DevTools、`SendIPCMessage`、任务栏进度/闪烁 |
| 事件桥 | `events.go` | `On` / `Off` / `RegisterIPCHandler`（跳板槽位池，上限 `MaxEventHandlers`=64） |
| 对话框/菜单 | `dialog.go` | `ShowNotification`、`ShowOpenDialog`/`ShowSaveDialog`/`ShowMessageBox`/`ShowErrorBox`、右键菜单 `MenuItemCreate`/`SetContextMenuItems` |
| 异步对话框 | `dialog_async.go` | `ShowOpenDialogAsync`/`ShowSaveDialogAsync`/`ShowMessageBoxAsync`（上限 `MaxAsyncDialogs`=16） |
| 托盘 | `tray.go` | `TrayCreate`/`TraySetMenu`（扁平表）/`TraySetIconFromFile`/`TraySetIconFromData` |
| YAML 存储 | `yaml.go` | `YAMLSet`/`YAMLGet`/`YAMLGetAll`/`YAMLKeys`/`YAMLHas`/`YAMLDelete`/`YAMLLen`/`YAMLClear`/`YAMLDeleteFile` |
| 系统工具 | `system.go` | 剪贴板、`GetPath`/`GetLocale`/`GetDisplaysInfo`、打印、全局热键、开机自启、URL 协议/文件关联、安全资源 `RegisterResource`、`GetFileIcon`、`SmartConvertEncoding`、`NTPNow` |
| JAPK 资源包 | `japk.go` | `SetPublicKey`/`LoadFromBytes`/`IsLoaded`/`GetAppSignature`/`GetSignatureInfo`/`Unload` |

有意不封装的 2 个：`cleanup_all_windows`（上游已废弃，用 `Exit`）、`yaml_get_str`（要求 `CoTaskMemFree` 释放，跨平台不可移植，用缓冲区版 `YAMLGet` 替代）。

### 事件系统要点

- `On(event, handler)` 注册、`Off(event, cbID)` 注销；事件名常量见 `include/JadeView.h`（`app-ready`、`window-closed`、`theme-changed`、`tray-menu-command`、`crash` 等 30+ 个）。
- **`app-ready` 必须在 `Init` 之前注册**，且回调里要判断 `windowID == 1`（0 = 初始化失败，`data` 为错误描述）。
- handler 返回非空字符串会作为响应回传给库；多数事件返回 `""` 即可。

## 目录结构

```
JadeView/
├── include/JadeView.h            # C 头文件（上游官方版）
├── lib/
│   ├── linux_amd64/libJadeView.{a,so}
│   ├── linux_arm64/libJadeView.{a,so}
│   ├── windows_amd64/
│   │   ├── JadeView.dll          # MSVC 版（自包含 WebView2Loader），被 go:embed 内置
│   │   ├── libJadeView_delay.a   # dlltool 生成的延迟导入库（链接用）
│   │   └── JadeView.def          # 导出表记录（重生成延迟库用）
│   └── windows_386/              # 同上（32 位，def 含 @N 修饰表）
├── beta/                         # 上游版本/API 文档
├── jadeview.go                   # 生命周期 + cgo 公共前导
├── jadeview_linux_{amd64,arm64}.go   # Linux 链接配置（静态）
├── jadeview_windows_{amd64,386}.go   # Windows 链接配置（延迟链 DLL）
├── jadeview_windows_embed*.go        # Windows 内置 DLL + 运行时释放
├── window.go / window_ext.go / events.go / dialog.go / dialog_async.go
├── tray.go / yaml.go / system.go / japk.go / helpers.go
└── example/                      # 跨平台可交互示例
```

## Windows 自包含原理

DLL 默认是加载期导入——进程启动前 Windows 就要求它存在，无法进程内自解压。本封装分两步解决：

1. **延迟加载**：用 `dlltool` 从 DLL 导出表生成延迟导入库 `libJadeView_delay.a`（链接参数 `-lJadeView_delay -ldelayimp`），DLL 变为「首次调用其函数时才 `LoadLibrary`」。
2. **go:embed + 预加载**：`jadeview_windows_embed*.go` 把 DLL 编进二进制，`init()` 启动时释放到 `%TEMP%\jadeview\` 并按全路径预加载，之后延迟加载命中该同名已加载模块。

纯动态链接，不碰 MSVC 静态库的 CRT 符号问题；MSVC 版 `JadeView.dll` 自带 WebView2Loader，无需额外 loader。

## 升级上游库

上游发新版后按平台替换，并注意以下步骤：

**Linux**：直接用新版 `libJadeView.a`/`libJadeView.so` 覆盖 `lib/linux_*/`，重新构建即可。

**Windows**：覆盖 `lib/windows_*/JadeView.dll` 后，**必须重新生成延迟导入库**（需 MinGW 的 `objdump`/`dlltool`）：

```bash
# x64（无名字修饰，直接从导出表生成）
cd lib/windows_amd64
objdump -p JadeView.dll | grep '+base\[' | awk '{print $NF}' | grep -v '^RVA$' | sort -u > exports.txt
{ echo "LIBRARY JadeView.dll"; echo "EXPORTS"; cat exports.txt; } > JadeView.def
dlltool --input-def JadeView.def --output-delaylib libJadeView_delay.a --dllname JadeView.dll
```

```bash
# x86（stdcall 有 @N 修饰：.def 用 name@N + --kill-at）
# 在现有 JadeView.def 基础上增补新函数（@N = 参数字节数，按头文件签名计算，
# 指针/int32 各 4 字节；不要从导出表重建，会丢失 @N 信息）
cd lib/windows_386
C:/msys64/mingw32/bin/dlltool --kill-at --input-def JadeView.def \
  --output-delaylib libJadeView_delay.a --dllname JadeView.dll
```

升级后建议全流程核对：

1. `gcc` 试编译核对新头文件（上游自动生成，出过 `i64` 这类非 C 类型笔误）；
2. 对比 DLL 导出表与延迟库符号，确认无缺漏（链接报 `undefined reference to <某导出>` 多半是没重生成）;
3. 检查新增 API 是否已有 Go 封装。

## 已知问题 / 注意事项

- **上游版本**：当前全部为 v2.3.0-beta.10，Windows DLL 与 Linux 库已统一。
- **YAML 等持久化 API 须在 `app-ready` 之后调用**（依赖 `Init` 的 `data_directory` 就绪）。
- **`app_signature` 至少 6 个字符**，过短 `Init` 返回失败且不启动 GUI 线程。
- **Windows 临时释放**：个别杀软可能对「释放 DLL 并加载」有启发式告警，属正常；可引导用户加白。
- **Linux 托盘**：依赖桌面环境（appindicator 等）支持，无托盘协议的环境创建会失败，代码应容错。
- **`jade-region-drag` 拖动区为 Windows 平台特性**（上游文档标注），Linux 请用 CSS `-webkit-app-region: drag`。
- **Windows arm64 无自包含**：链接配置已启用（`jadeview_windows_arm64.go`），但受限于工具链
  无延迟库支持，需随 exe 携带 DLL；构建需 llvm-mingw（见其小节）。
- **重生成延迟库后务必确认文件已产出**：binutils `dlltool` 在不支持的目标上会**退出码 0 但不写文件**
  （arm64 踩过），脚本里要检查 `.a` 的存在与时间戳。
- `lib/` 下的二进制是 JadeView 作者的第三方产物，**不在本项目 MIT 许可范围内**。

## 许可证

本 Go 封装层代码以 [MIT](LICENSE) 许可证发布。
