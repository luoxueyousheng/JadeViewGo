# JadeView Go 封装

[JadeView](https://jade.run) WebView 桌面库的 Go 封装 —— 用 Go + HTML/CSS/JS 写 Windows 桌面应用。窗口、事件、双向 IPC、托盘、对话框、通知、YAML 持久化、NTP 授时一应俱全,头文件 129 个导出函数中 128 个已封装(仅 `yaml_get_str` 因跨平台内存管理差异未封装)。

> ⚠️ **v2.4.0 破坏性变更**:JAPK 包从此**只能加载带签名的资源包**(v3 签名协议,平台根证书链严格离线验签),混淆包与公钥注入机制(`SetPublicKey`)均已移除。详见[已知问题](#已知问题--注意事项)与 [CHANGELOG](CHANGELOG.md)。

当前对应上游 **v2.4.3 (Build 26I01)**;要求 **Go 1.23+**。

## 目录

- [支持平台](#支持平台)
- [安装](#安装)
- [前置条件与上手](#前置条件与上手)
- [快速开始](#快速开始)
- [示例](#示例example)
- [API 总览](#api-总览)
- [目录结构](#目录结构) · [实现原理](#纯-go-实现原理) · [升级上游库](#升级上游库) · [已知问题](#已知问题--注意事项)

## 支持平台

| 平台 | 架构 | 实现 | 构建依赖 | 运行时依赖 | 分发形态 |
|------|------|------|----------|------------|----------|
| **Windows** | amd64 / 386 / arm64 | 纯 Go（syscall 直调内置 DLL） | **仅 Go 工具链** | WebView2 Runtime（Win11 自带） | 单 exe 自包含 |

## 安装

作为依赖引入你的项目:

```bash
go get github.com/luoxueyousheng/JadeViewGo@latest    # 最新正式版
go get github.com/luoxueyousheng/JadeViewGo@v2.4.3    # 锁定指定版本
```

只想先跑一眼内置示例(示例是模块子包,可直接运行):

```bash
go run github.com/luoxueyousheng/JadeViewGo/example@v2.4.3
```

> 若 `@latest` 一时解析不到刚发布的 tag(官方 proxy 索引有几分钟延迟),改用精确版本号
> `@v2.4.3`,或加 `GOPROXY=https://proxy.golang.org,direct` 显式拉取。

## 前置条件与上手

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
    // Init(开发模式, 日志路径, 数据目录, 应用名, 应用签名, 单实例)
    // 签名≥6字符,建议反域名格式——JAPK 模式下它就是 JADE:// URL 的主机名
    jadeview.Init(true, "", "", "my-app", "com.example.myapp", false)
    jadeview.RunMessageLoop() // 阻塞直到退出
}
```

## 示例（example/）

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
- **托盘**:右键菜单显示/隐藏窗口、退出。
- **WebView 设置**:`DefaultWebViewSettings` 起步,用 `PreloadJS` 在页面脚本运行前注入
  平台信息(`window.__JV_ENV`),前端同步读取做平台适配(标题栏/材质),`env` IPC 通道兜底。

**前端站点加载(三选一,见 `onAppReady` 里的 `plan`)**:

| plan | 方式 | 适用场景 |
|------|------|----------|
| 0(默认) | **JAPK 资源包**:`go:embed` 的 `app.japk` 经 `LoadFromBytes` 内存加载,再以**空路径**调 `SetProtocolServicePath("", false)` 切到内存 JAPK 模式,URL 形如 `JADE://<app_signature>` | 加密分发(⚠️ v2.4.0 起仅限签名包) |
| 1 | **协议服务挂源码目录**:`SetProtocolServicePath("example/site", hotReload=true)`,改动站点文件页面即时刷新 | 开发调试 |
| 2 | **进程内回环 HTTP**:127.0.0.1 随机端口直出 `embed.FS`,磁盘零前端文件 | 单 exe 分发(不想用 JAPK 时) |

三种方式返回/构造的 URL 都直接建窗导航。注意:

- **协议服务的站点目录不要与 `Init` 的 data_directory 相同或嵌套**——库持续写数据会触发
  「写→热载刷新」死循环。
- **跨域与 IPC**:用在线或本地 http(s) URL(如 plan 2 的回环 HTTP、任意线上页面)而非
  jade:// 资源方式(plan 0/1)加载页面时,页面源与库不同源,IPC 通讯(`jade.invoke`/
  postMessage)可能被跨域限制拦截。`WebViewSettings.CORSWhitelist`/`PostMessageWhitelist`
  可设白名单,但库**没有接口能查到它注册的临时域**,无法精确加白——示例用
  `PostMessageWhitelist: "*"` 兜底;分发场景优先选 jade:// 同源方案(plan 0/1)。
- **JAPK 签名强制(v2.4.0)**:`LoadFromBytes` 只接受带 JadeTweak 平台根证明与叶子签名的 **v3 签名包**;混淆包(JPKBIN02)不再支持,`app_name`/`app_signature`
  必须与 `Init` 及平台证明一致。加载错误详情经 `japk-load-failed` 事件回报;`data:` URL 方案实测不可行(WebView2 拒绝 data: 顶层导航)。
- 托盘图标走内存 API(`TraySetIconFromData`,.ico 仅 Windows 传)。

## API 总览

共享类型在 `types.go`,实现是 `*_windows.go`(纯 Go)。

| 模块 | 主要函数 |
|------|----------|
| 生命周期 | `Init` / `Version` / `RunMessageLoop` / `Exit` / `ExitWait`（等待后台线程完全退出,便于可靠卸载）/ `Preload`（提前加载 DLL 并拿到错误） |
| 窗口创建 | `CreateWindow`（`WindowOptions`/`WebViewSettings`,默认值用 `DefaultWindowOptions`/`DefaultWebViewSettings`;v2.4.0 新增 `ProfileName`,多窗口 Cookie/存储/缓存隔离）、`CreateBorderlessWindow`、`Navigate`、`GoBack`/`GoForward`、`CanGoBack`/`CanGoForward`、`ExecuteJavaScript`、`SetTitle/SetSize/SetPosition/...` |
| 窗口扩展 | 状态查询 `Is*`、`GetWindowBounds`、`GetWindowHWND`⇄`GetWindowID`、层级/背景/全屏/主题/缩放、DevTools、`SendIPCMessage`、任务栏进度/闪烁 |
| 事件桥 | `On` / `Off` / `RegisterIPCHandler`（槽位池,上限 `MaxEventHandlers`=64）、统一网页权限处理器 `SetWebviewPermissionHandler` / `ClearWebviewPermissionHandler`（摄像头/麦克风/录屏/文件访问等） |
| 对话框/菜单 | `ShowNotification`、`ShowOpenDialog`/`ShowSaveDialog`/`ShowMessageBox`/`ShowErrorBox`、右键菜单 `MenuItemCreate`/`SetContextMenuItems` |
| 异步对话框 | `ShowOpenDialogAsync`/`ShowSaveDialogAsync`/`ShowMessageBoxAsync`（上限 `MaxAsyncDialogs`=16） |
| 托盘 | `TrayCreate`/`TraySetMenu`（扁平表）/`TraySetIconFromFile`/`TraySetIconFromData` |
| YAML 存储 | `YAMLSet`/`YAMLGet`/`YAMLGetAll`/`YAMLKeys`/`YAMLHas`/`YAMLDelete`/`YAMLLen`/`YAMLClear`/`YAMLDeleteFile` |
| 系统工具 | 剪贴板、`GetPath`/`GetLocale`/`GetDisplaysInfo`、打印、全局热键、开机自启、URL 协议/文件关联、安全资源、`GetFileIcon`、`SmartConvertEncoding`、`NTPNow` |
| JAPK 资源包 | `LoadFromBytes`/`IsLoaded`/`GetAppSignature`/`GetSignatureInfo`/`Unload`（v2.4.0 起仅限签名包,`SetPublicKey` 已随上游移除） |

有说明的几点:`cleanup_all_windows` 与 `JadeView_set_public_key` 已被上游 v2.4.0 **整体移除**,本封装不再提供;导航历史 API
已封装为 `GoBack` / `GoForward` / `CanGoBack` / `CanGoForward`(对应 `webview_go_back` 等 4 个导出);`yaml_get_str`
不封装(要求 `CoTaskMemFree` 释放,不可移植,用缓冲区版 `YAMLGet` 替代)。

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
JadeViewGo/
├── include/JadeView.h            # C 头文件（上游官方 2.4.0 版,API 参考）
├── lib/windows_{amd64,386,arm64}/JadeView.dll   # MSVC 版（自含 WebView2Loader）,被 go:embed 内置
├── doc.go / types.go             # 包文档 + 共享类型
├── enums.go / events_names.go    # 参数枚举 + 事件名常量
├── *_windows.go                  # 纯 Go 实现,共 10 个：
│                                 #   dll(核心+地址表) / window(生命周期+窗口) /
│                                 #   events(事件桥) / dialog(对话框+托盘) /
│                                 #   system(系统+YAML+JAPK) / fltcall×2 / embed×3
└── example/                      # 可交互示例
```

> 仓库内另有历史遗留的 Linux cgo 实现与 `lib/linux_*` 旧库（上游已停止更新,暂不维护）,不随 Windows 构建参与编译。

## 纯 Go 实现原理

1. **内置与释放**:`dll_embed_windows_*.go` 按架构 go:embed 对应的 `JadeView.dll`;
   **首次调用任一 API(或 `Preload`)时**才释放到 `%TEMP%\jadeview\<架构>-<内容哈希前8位>\`
   (内容寻址:换版本换目录,已存在文件按完整 sha256 校验、不符重写,多进程多版本并存安全;
   仅 import 本包无任何磁盘副作用)。exe 同目录的 `JadeView.dll` 优先。
2. **加载与调用**:`syscall.NewLazyDLL` 按绝对路径惰性加载;头文件 129 个导出函数中已封装的
   128 个经惰性代理 `jvProc.Call` 直调(`dll_windows.go` 内含完整地址表)。**加载失败时首次 API
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

把新版三个架构的 `JadeView.dll` 覆盖到 `lib/windows_*/` 即可,重新构建自动生效
(go:embed 重新打包,运行时按新哈希释放新目录)。**不再需要 dlltool/objdump 重做导入库**。
上游若新增/修改 API,在 `dll_windows.go` 的地址表加条目并补包装函数;改动结构体时需重新做布局比对。

升级后建议:

1. 核对新头文件与 `dll_windows.go` 的函数清单差异(上游自动生成的头文件出过
   `i64` 这类非 C 类型笔误);
2. 跑一遍 `go run ./example` 冒烟验证。

## 已知问题 / 注意事项

- **JAPK 仅限签名包(v2.4.0 破坏性变更)**:`LoadFromBytes` 只接受 v3 签名包(平台根证书链严格离线验签),
  混淆包(JPKBIN02)与运行时公钥注入(`SetPublicKey`)均已废弃。分发前端资源可临时改用协议服务目录模式(plan 1)
  或回环 HTTP(plan 2)。
- **`app-ready` 之后再调持久化 API**:YAML 等依赖 `Init` 的 `data_directory` 就绪。
- **`app_signature` 至少 6 个字符**,过短 `Init` 返回失败且不启动 GUI 线程;建议反域名格式
  (如 `com.example.myapp`)——JAPK 模式下它会作为 `JADE://` URL 的主机名。
- **Windows 杀软告警**:个别杀软对「释放 DLL 并加载」或浮点跳板的可执行内存分配有
  启发式告警,属正常;可引导用户加白。若 DLL 释放/加载被拦截,首次 API 调用会 panic——
  用 `Preload()` 在启动早期探测并优雅提示。
- **事件槽位上限**:`On` 与 `RegisterIPCHandler` **共享** `MaxEventHandlers`=64 个槽位;
  IPC handler 无注销 API(上游头文件亦无),注册后**永久占用**一个槽位,规划通道数量时留意。
- **http(s) 页面的 IPC 跨域风险**:页面以在线/本地 http(s) URL 加载(而非 jade:// 资源方式)时
  与库不同源,IPC 可能被跨域拦截;`WebViewSettings.CORSWhitelist` 可设白名单,但库无接口查询
  其注册的临时域,无法精确加白。规避:优先用协议服务/JAPK(jade:// 同源),或
  `PostMessageWhitelist: "*"` 兜底。
- **`jade-region-drag` 拖动区为 Windows 特性**(上游文档标注)。
- **`lib/` 下的二进制**是 JadeView 作者的第三方产物,**不在本项目 MIT 许可范围内**。

## 许可证

本 Go 封装层代码以 [MIT](LICENSE) 许可证发布。
