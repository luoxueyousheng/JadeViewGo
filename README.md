# JadeView Go 封装

[JadeView](https://jade.run) WebView 桌面库的 Go 封装 —— 用 Go + HTML/CSS/JS 写跨平台桌面应用。窗口、事件、双向 IPC、托盘、对话框、通知、YAML 持久化、NTP 授时一应俱全,头文件 124 个导出函数全部封装。

当前对应上游 **v2.3.0-beta.10**;要求 **Go 1.23+**。

## 目录

- [支持平台](#支持平台)
- [安装](#安装)
- [平台前置条件与上手](#平台前置条件与上手) — [Windows](#windows) · [Linux](#linux)
- [快速开始](#快速开始)
- [示例](#示例example)
- [API 总览](#api-总览)
- [目录结构](#目录结构) · [实现原理](#windows-纯-go-实现原理) · [升级上游库](#升级上游库) · [已知问题](#已知问题--注意事项)

## 支持平台

| 平台 | 架构 | 实现 | 构建依赖 | 运行时依赖 | 分发形态 |
|------|------|------|----------|------------|----------|
| **Windows** | amd64 / 386 / arm64 | 纯 Go（syscall 直调内置 DLL） | **仅 Go 工具链** | WebView2 Runtime（Win11 自带） | 单 exe 自包含 |
| **Linux** | amd64 / arm64 | cgo 静态链接 `libJadeView.a` | Go + gcc + GTK3 / WebKit2GTK / xdo 开发包 | GTK3 / WebKit2GTK / libxdo3 + 图形桌面 | 单二进制（依赖系统库） |

两侧公共 API 完全一致,同一份业务代码可直接跨平台编译;差异只在**构建前置条件**——详见下方[平台前置条件与上手](#平台前置条件与上手)。

## 安装

作为依赖引入你的项目:

```bash
go get github.com/luoxueyousheng/JadeViewGo@latest    # 最新正式版
go get github.com/luoxueyousheng/JadeViewGo@v0.2.1    # 锁定指定版本
```

只想先跑一眼内置示例(示例是模块子包,可直接运行):

```bash
go run github.com/luoxueyousheng/JadeViewGo/example@v0.2.1
```

> - 依赖拉下来后,**构建仍需满足对应平台的前置条件**(下一节);Linux 尤其别漏系统开发包。
> - 若 `@latest` 一时解析不到刚发布的 tag(官方 proxy 索引有几分钟延迟),改用精确版本号
>   `@v0.2.1`,或加 `GOPROXY=https://proxy.golang.org,direct` 显式拉取。

## 平台前置条件与上手

> 挑你的目标系统看对应小节即可,每节都是「前置条件 → 构建 → 运行」自包含流程。

### Windows

**前置条件**

| 环节 | 要求 |
|------|------|
| 构建机 | 只需 **Go 工具链**。无需 MinGW / MSYS2 / llvm-mingw,无需设置 `CC` / `CGO_ENABLED`。 |
| 目标机 | Edge **WebView2 Runtime**(Win11 自带;Win10 用微软 Evergreen Bootstrapper 安装)。 |

对应架构的 `JadeView.dll` 已用 go:embed 编进二进制,**无需随程序分发 DLL**。

**构建**

```powershell
go build -o myapp.exe .                             # 当前架构,控制台版(能看日志)
go build -ldflags "-H windowsgui" -o myapp.exe .    # GUI 版(无 cmd 黑窗)

# 交叉编译其它架构(纯 Go,任何机器上都行)
$env:GOARCH="386";   go build -ldflags "-H windowsgui" -o myapp_x86.exe .
$env:GOARCH="arm64"; go build -ldflags "-H windowsgui" -o myapp_arm64.exe .
```

**运行 / 分发**

产物是**单个 exe**。首次运行时对应架构的 `JadeView.dll` 自动释放到
`%TEMP%\jadeview\<架构>-<内容哈希>\`(内容寻址,多版本/多架构并存互不覆盖)。
若 exe 同目录放了 `JadeView.dll`,则优先用它(便于调试或临时换库)。

> 调试看日志用控制台版;`-H windowsgui` 无控制台,`fmt.Printf` 看不到输出。

### Linux

Linux 侧走 cgo,静态链接 `libJadeView.a`;该库仍依赖系统的 GTK3 / WebKit2GTK / libxdo 动态库。构建机与目标机的依赖要**分别**装。

**前置条件 · 构建机**(Debian / Ubuntu 系)

```bash
sudo apt install build-essential pkg-config \
    libgtk-3-dev libwebkit2gtk-4.1-dev libxdo-dev
```

| 包 | 作用 | 漏了会报 |
|----|------|----------|
| `build-essential` | gcc + libc 头文件(cgo 必需) | `stdlib.h: No such file or directory` |
| `pkg-config` | 拉取 GTK/WebKit 链接参数 | 找不到库 |
| `libgtk-3-dev` · `libwebkit2gtk-4.1-dev` | GTK3 / WebKit2GTK 开发包 | 编译或链接失败 |
| `libxdo-dev` | libxdo(xdotool)——`libJadeView.a` 硬引用它(托盘菜单发按键序列) | `cannot find -lxdo` |

> **老发行版只有 WebKit2GTK 4.0**(无 4.1 包):把 `jadeview_linux_amd64.go` /
> `jadeview_linux_arm64.go` 里的 `webkit2gtk-4.1` 改成 `webkit2gtk-4.0`。先用
> `pkg-config --exists webkit2gtk-4.1 && echo 有4.1 || echo 用4.0` 确认。

**前置条件 · 目标机(运行时)**

```bash
sudo apt install libgtk-3-0 libwebkit2gtk-4.1-0 libxdo3   # 运行时库(非 -dev)
```

外加 **X11 / Wayland 图形桌面**——GUI 需要显示环境。

**构建 / 运行**

```bash
CGO_ENABLED=1 go build ./...     # 验证编译 + 链接
go run ./example                 # 需要图形桌面
```

**无桌面 / 远程 X11 / 无 GPU 环境**(headless 服务器、NAS、SSH X11 转发)

WebKit 默认走 GPU 合成,拿不到 EGL/DRI(`/dev/dri` 缺失或无权限)时会直接崩溃。强制软件渲染:

```bash
WEBKIT_DISABLE_DMABUF_RENDERER=1 WEBKIT_DISABLE_COMPOSITING_MODE=1 \
LIBGL_ALWAYS_SOFTWARE=1 go run ./example
```

- **SSH X11 转发**:用**原登录用户**跑;`su` / `sudo` 切到 root 会丢失 X 授权 cookie,报 `No authorisation provided`。
- **纯 headless**(无任何显示):GUI 起不来,只能 `go build ./...` 验证编译链接是否通过。

**arm64 补充**

库文件在 `lib/linux_arm64/`,用法与 amd64 完全一致。**推荐在 arm64 机器上原生构建**
(树莓派、云 ARM 实例、ARM NAS 等);交叉编译需 `aarch64-linux-gnu-gcc` 加 arm64 版
GTK/WebKit sysroot,配置繁琐,一般不值得。

## 快速开始

```go
package main

import jadeview "github.com/luoxueyousheng/JadeViewGo"

func main() {
    // 关键时序：app-ready 必须在 Init 之前注册，并在回调里判断 windowID==1 再建窗
    jadeview.On(jadeview.EventAppReady, func(windowID uint32, data string) string {
        if windowID == 1 {
            opts := jadeview.DefaultWindowOptions()
            opts.Title = "Hello JadeView"
            jadeview.CreateWindow("https://example.com", 0, &opts, nil)
        }
        return ""
    })
    jadeview.On(jadeview.EventWindowAllClosed, func(uint32, string) string {
        jadeview.Exit()
        return ""
    })
    // Init(开发模式, 日志路径, 数据目录, 应用名, 应用签名≥6字符, 单实例)
    jadeview.Init(true, "", "", "my-app", "my-app-signature", false)
    jadeview.RunMessageLoop() // 阻塞直到退出
}
```

## 示例（example/）

一份代码,Windows / Linux 都能跑(平台差异用 `runtime.GOOS` 分支)。装好对应平台前置条件后:

```bash
go run ./example
```

演示内容:

- **外观**:颜色模式切换(浅色/深色/跟随系统,联动 `SetTheme` + 标题栏图标色)、
  窗口材质切换(Mica / Mica Alt / Acrylic / 纯色,`SetBackdrop`/`SetBackgroundColor`)、
  页面缩放(`SetZoom`)。
- **IPC 测试**:任意通道 + JSON payload 的 `jade.invoke` 回声(显示往返耗时)、
  宿主连发推送(`SendIPCMessage` → `jade.on`)、四级 Toast 契约演示、通信日志。
- **窗口**:置顶开关、最小化、全屏、任务栏闪烁、边界查询、HWND⇄窗口ID 互查、DevTools。
- **系统**:异步对话框(打开/保存/消息框)、系统通知、剪贴板读写、NTP 网络时间。
- **存储**:YAML 写入/读取/全量读取(存于 `Init` 的数据目录)。
- **托盘**:右键菜单显示/隐藏窗口、退出(Linux 先探测 D-Bus 托盘协议,无支持则跳过,见「已知问题」)。
- **WebView 设置**:`DefaultWebViewSettings` 起步,用 `PreloadJS` 在页面脚本运行前注入
  平台信息(`window.__JV_ENV`),前端同步读取做平台适配(标题栏/材质),`env` IPC 通道兜底。

**前端整目录内置,运行时零落盘**:`example/site/`(index.html + fluent.css + app.js)用
`//go:embed all:site` 打成 `embed.FS` 编进 exe;运行时在 127.0.0.1 随机端口起进程内 HTTP
服务直接以 `embed.FS` 为根对外服务(`http.FileServer(http.FS(...))`),窗口导航到
`http://127.0.0.1:<port>/index.html`——**磁盘上不出现任何前端文件**,多文件/相对路径/子目录
全部可用,加文件无需改 Go 代码。托盘图标同样走内存 API(`TraySetIconFromData`)。

> 为什么不用库的协议服务?`SetProtocolServicePath` 只接受磁盘目录(库自己读文件),无法直接
> 挂 `embed.FS`;要用它就得先把 embed 内容释放到临时目录(此方式的代码见 git 历史/文档,
> 附带热载能力)。纯内存的官方替代是 JAPK 资源包(`LoadFromBytes`,需上游打包工具)。
> `data:` URL 方案实测不可行(WebView2 拒绝 data: 顶层导航,窗口直接关闭)。

## API 总览

公共 API 跨平台一致,共享类型在 `types.go`;Windows 实现是 `*_windows.go`(纯 Go),
Linux 实现是不带后缀的 cgo 文件(`//go:build linux`)。

| 模块 | 主要函数 |
|------|----------|
| 生命周期 | `Init` / `Version` / `RunMessageLoop` / `Exit` / `Preload`（Windows 提前加载 DLL 并拿到错误;Linux 恒 nil） |
| 窗口创建 | `CreateWindow`（`WindowOptions`/`WebViewSettings`,默认值用 `DefaultWindowOptions`/`DefaultWebViewSettings`）、`CreateBorderlessWindow`、`Navigate`、`ExecuteJavaScript`、`SetTitle/SetSize/SetPosition/...` |
| 窗口扩展 | 状态查询 `Is*`、`GetWindowBounds`、`GetWindowHWND`⇄`GetWindowID`、层级/背景/全屏/主题/缩放、DevTools、`SendIPCMessage`、任务栏进度/闪烁 |
| 事件桥 | `On` / `Off` / `RegisterIPCHandler`（槽位池,上限 `MaxEventHandlers`=64） |
| 对话框/菜单 | `ShowNotification`、`ShowOpenDialog`/`ShowSaveDialog`/`ShowMessageBox`/`ShowErrorBox`、右键菜单 `MenuItemCreate`/`SetContextMenuItems` |
| 异步对话框 | `ShowOpenDialogAsync`/`ShowSaveDialogAsync`/`ShowMessageBoxAsync`（上限 `MaxAsyncDialogs`=16） |
| 托盘 | `TrayCreate`/`TraySetMenu`（扁平表）/`TraySetIconFromFile`/`TraySetIconFromData` |
| YAML 存储 | `YAMLSet`/`YAMLGet`/`YAMLGetAll`/`YAMLKeys`/`YAMLHas`/`YAMLDelete`/`YAMLLen`/`YAMLClear`/`YAMLDeleteFile` |
| 系统工具 | 剪贴板、`GetPath`/`GetLocale`/`GetDisplaysInfo`、打印、全局热键、开机自启、URL 协议/文件关联、安全资源、`GetFileIcon`、`SmartConvertEncoding`、`NTPNow` |
| JAPK 资源包 | `SetPublicKey`/`LoadFromBytes`/`IsLoaded`/`GetAppSignature`/`GetSignatureInfo`/`Unload` |

有意不封装的 2 个:`cleanup_all_windows`(上游已废弃,用 `Exit`)、`yaml_get_str`
(要求 `CoTaskMemFree` 释放,跨平台不可移植,用缓冲区版 `YAMLGet` 替代)。

**枚举**:固定取值的参数都有二级命名空间枚举(`enums.go`),不必裸写字符串/数字——
`Theme.Dark`、`FrameStyle.TitleOverlay`、`WindowLevel.Topmost`、`Backdrop.Mica`、
`MsgBoxType.Warning`、`ProgressState.Indeterminate`、`TrayItem.Divider`、`MenuKind.Checkbox`、
`DialogProp.MultiSelections`、`Encoding.GBK`;事件名见 `Event*` 常量(`events_names.go`)。

### 事件系统要点

- `On(event, handler)` 注册、`Off(event, cbID)` 注销;**事件名用库提供的 Go 常量**(`events_names.go`,
  与头文件 `JADEVIEW_EVENT_*` 宏一一对应):`jadeview.EventAppReady`、`EventWindowClosed`、
  `EventThemeChanged`、`EventTrayMenuCommand` 等 34 个;`EventCrash` 的 `data` 取值见 `Crash*` 常量。
- **`app-ready` 必须在 `Init` 之前注册**,且回调里要判断 `windowID == 1`(0 = 初始化失败,`data` 为错误描述)。
- handler 返回非空字符串会作为响应回传给库;多数事件返回 `""` 即可。

## 目录结构

```
JadeView/
├── include/JadeView.h            # C 头文件（上游官方版，Linux cgo 用；Windows 仅作 API 参考）
├── lib/
│   ├── linux_amd64/libJadeView.{a,so}
│   ├── linux_arm64/libJadeView.{a,so}
│   ├── windows_amd64/JadeView.dll    # MSVC 版（自含 WebView2Loader），被 go:embed 内置
│   ├── windows_386/JadeView.dll
│   └── windows_arm64/JadeView.dll
├── beta/                         # 上游版本/API 文档
├── doc.go / types.go            # 包文档 + 跨平台共享类型
├── enums.go / events_names.go    # 参数枚举 + 事件名常量
├── *_windows.go                  # Windows 纯 Go 实现，共 10 个：
│                                 #   dll(核心+地址表) / window(生命周期+窗口) /
│                                 #   events(事件桥) / dialog(对话框+托盘) /
│                                 #   system(系统+YAML+JAPK) / fltcall×2 / embed×3
├── jadeview.go window.go ...     # Linux cgo 实现（//go:build linux）
├── jadeview_linux_{amd64,arm64}.go   # Linux 链接配置（静态 + -lxdo）
└── example/                      # 跨平台可交互示例
```

## Windows 纯 Go 实现原理

1. **内置与释放**:`dll_embed_windows_*.go` 按架构 go:embed 对应的 `JadeView.dll`;
   **首次调用任一 API(或 `Preload`)时**才释放到 `%TEMP%\jadeview\<架构>-<内容哈希前8位>\`
   (内容寻址:换版本换目录,已存在文件按完整 sha256 校验、不符重写,多进程多版本并存安全;
   仅 import 本包无任何磁盘副作用)。exe 同目录的 `JadeView.dll` 优先。
2. **加载与调用**:`syscall.NewLazyDLL` 按绝对路径惰性加载;全部 124 个导出函数经
   惰性代理 `jvProc.Call` 直调(`dll_windows.go` 内含完整地址表)。**加载失败时首次 API
   调用会 panic**(`syscall.LazyProc` 语义)——需优雅降级的宿主在启动早期调 `Preload()`
   检查错误即可。
3. **结构体传参**:`WebViewWindowOptions` 等 6 个 C 结构体在 Go 侧逐字段镜像
   (`window_windows.go` 等),布局已用 C `offsetof` 与 Go `unsafe.Offsetof`
   双端逐字段比对验证(amd64/386;arm64 与 amd64 对齐规则相同)。
4. **回调**:事件桥用 `syscall.NewCallback`(stdcall)+ 固定槽位池;异步对话框回调
   在 386 下是 cdecl,用 `NewCallbackCDecl`(64 位两者等价)。
5. **double 参数**:`set_webview_zoom` 的 double 在 x64/arm64 走浮点寄存器,syscall
   无法直传,经一段运行时生成的 8 字节跳板装入 XMM1/D1 后跳转(`fltcall_windows_*.go`,
   已用测试 DLL 端到端验证);386 的 double 走栈,直接拆两个字传递。

> 与 cgo 方案的取舍:失去了编译期对头文件的类型校验——升级上游 API 时须人工核对
> `dll_windows.go` 的地址表与各镜像结构体(布局比对方法见上),参数错误只会在运行时暴露。

## 升级上游库

**Windows**:把新版三个架构的 `JadeView.dll` 覆盖到 `lib/windows_*/` 即可,重新构建
自动生效(go:embed 重新打包,运行时按新哈希释放新目录)。**不再需要 dlltool/objdump
重做导入库**。上游若新增/修改 API,在 `dll_windows.go` 的地址表加条目并补包装函数;
改动结构体时需重新做布局比对。

**Linux**:直接用新版 `libJadeView.a`/`libJadeView.so` 覆盖 `lib/linux_*/`,重新构建。
若上游新增了对其它系统库的依赖,记得在 `jadeview_linux_*.go` 的 `#cgo LDFLAGS` 补 `-l`。

升级后建议:

1. 核对新头文件与 `dll_windows.go` 的函数清单差异(上游自动生成的头文件出过
   `i64` 这类非 C 类型笔误,Linux cgo 侧会直接编译失败);
2. Windows 与 Linux 各跑一遍 `go run ./example` 冒烟验证。

## 已知问题 / 注意事项

- **上游版本**:当前全部为 v2.3.0-beta.10,Windows DLL 与 Linux 库已统一。
- **`app-ready` 之后再调持久化 API**:YAML 等依赖 `Init` 的 `data_directory` 就绪。
- **`app_signature` 至少 6 个字符**,过短 `Init` 返回失败且不启动 GUI 线程。
- **Windows 杀软告警**:个别杀软对「释放 DLL 并加载」或浮点跳板的可执行内存分配有
  启发式告警,属正常;可引导用户加白。若 DLL 释放/加载被拦截,首次 API 调用会 panic——
  用 `Preload()` 在启动早期探测并优雅提示。
- **事件槽位上限**:`On` 与 `RegisterIPCHandler` **共享** `MaxEventHandlers`=64 个槽位;
  IPC handler 无注销 API(上游头文件亦无),注册后**永久占用**一个槽位,规划通道数量时留意。
- **Linux 托盘会崩、须先探测**:beta.10 的 `tray_create` 在**没有 StatusNotifier 托盘协议**的
  桌面(如 Debian/GNOME 默认桌面,需另装 AppIndicator 扩展)上不是返回 0,而是让库 GUI 线程
  RUNTIME_PANIC 直接 abort(已反馈上游)。调用前先探测会话 D-Bus 上有无
  `org.kde.StatusNotifierWatcher`,没有就跳过托盘——参考 `example/main.go` 的
  `hasStatusNotifierWatcher`(`dbus-send` 查 `NameHasOwner`)。
- **`jade-region-drag` 拖动区为 Windows 特性**(上游文档标注),Linux 请用 CSS `-webkit-app-region: drag`。
- **`lib/` 下的二进制**是 JadeView 作者的第三方产物,**不在本项目 MIT 许可范围内**。

## 许可证

本 Go 封装层代码以 [MIT](LICENSE) 许可证发布。
