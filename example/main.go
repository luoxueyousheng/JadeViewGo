//go:build linux || windows

// JadeView Go 封装 · 跨平台可交互示例（Windows / Linux 通用，单一代码库）。
//
// 演示内容：窗口、事件桥、IPC 双向通信、异步对话框、系统通知、托盘菜单、
// YAML 持久化、剪贴板、NTP 网络时间、HWND/窗口ID 互查、显示器/系统信息。
// 平台差异统一用 runtime.GOOS 分支处理，两端均可直接构建运行。
//
// 运行（详见 README「各平台使用方式」）：
//
//	Windows: go run ./example            （纯 Go，无需 C 编译器，单 exe 自包含 DLL）
//	Linux  : go run ./example            （cgo，需 gcc + GTK3/WebKit2GTK，静态链接）
//
// 关键时序（官方要求）：
//  1. app-ready 必须在 Init 之前注册；
//  2. 建窗、协议服务、托盘等都放在 app-ready 回调里；
//  3. 回调里必须判断 windowID==1 才算初始化成功。
package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	jadeview "github.com/luoxueyousheng/JadeViewGo"
)

// 前端站点整目录内置（HTML/CSS/JS，含子目录）。运行时**不落盘**：
// 库的协议服务(SetProtocolServicePath)只认磁盘路径、无法直接挂 embed.FS，
// 故进程内起一个 127.0.0.1 回环 HTTP 服务，直接以 embed.FS 为文件系统对外服务——
// 多文件/相对路径/子目录全部可用，新增文件无需改代码。
// 托盘图标同样走内存 API（TraySetIconFromData）。
// `all:` 前缀确保以 . 和 _ 开头的文件也被包含。
//
//go:embed all:site
var siteFS embed.FS

//go:embed icon.ico
var iconICO []byte

// titlebarHeight 标题栏覆盖层高度（像素），须与前端 fluent.css 里
// #app 的网格行高（.title-bar 所在行）保持一致，否则内置控制按钮与自绘内容错位。
const titlebarHeight = 40

var (
	mainWindowID uint32
	trayID       uint32
	alwaysOnTop  bool
	dataDir      string // JadeView 数据目录（库运行时数据，WebView 内核要求必须在磁盘上）
)

func main() {
	fmt.Printf("平台: %s/%s\n", runtime.GOOS, runtime.GOARCH)

	// 前端零落盘：页面走 data: URL、托盘图标走内存 API，见 buildDataURL/setupTray。
	// dataDir 是库自身的运行时数据目录（WebView2 用户数据、YAML 存储），无法免除。
	dataDir = filepath.Join(os.TempDir(), "jadeview-demo", "data")
	_ = os.MkdirAll(dataDir, 0o755)

	// 1) app-ready 之前注册事件（事件名用库提供的 Event* 常量，避免拼错）
	jadeview.On(jadeview.EventAppReady, onAppReady)
	jadeview.On(jadeview.EventWindowCreated, func(id uint32, _ string) string {
		fmt.Printf("[事件] window-created id=%d\n", id)
		return ""
	})
	jadeview.On(jadeview.EventWindowClosed, func(id uint32, _ string) string {
		fmt.Printf("[事件] window-closed id=%d\n", id)
		return ""
	})
	jadeview.On(jadeview.EventWindowAllClosed, func(_ uint32, _ string) string {
		fmt.Println("[事件] 所有窗口已关闭，退出")
		jadeview.Exit()
		return ""
	})
	jadeview.On(jadeview.EventThemeChanged, func(_ uint32, data string) string {
		fmt.Printf("[事件] 系统主题变化: %s\n", data)
		// 转发给前端，System 模式下页面据此重算明暗
		if mainWindowID != 0 {
			jadeview.SendIPCMessage(mainWindowID, "theme-changed", data)
		}
		return ""
	})
	jadeview.On(jadeview.EventTrayMenuCommand, onTrayMenuCommand)
	jadeview.On(jadeview.EventTrayEvent, func(_ uint32, data string) string {
		fmt.Printf("[事件] tray-event: %s\n", data)
		return ""
	})
	jadeview.On(jadeview.EventCrash, func(_ uint32, data string) string {
		fmt.Printf("[事件] 崩溃报告: %s\n", data) // data 为 Crash* 错误代码
		return ""
	})

	// 2) IPC 处理器（前端 jade.invoke 的目标）
	registerIPCHandlers()

	// 3) 初始化（版本号在 Init 之后才可取）
	if !jadeview.Init(true, "", dataDir, "jadeview-go-demo", "jadeview-go-demo-signature", true) {
		fmt.Println("Init 失败")
		return
	}
	fmt.Println("JadeView 版本:", jadeview.Version())

	// 4) 阻塞运行消息循环
	jadeview.RunMessageLoop()
	fmt.Println("消息循环结束")
}

func onAppReady(windowID uint32, data string) string {
	if windowID != 1 {
		fmt.Println("[app-ready] 初始化失败:", data)
		jadeview.Exit()
		return ""
	}
	fmt.Println("[app-ready] 初始化成功")

	// 纯内存前端：回环 HTTP 服务直出 embed.FS，运行期不向磁盘释放任何前端文件
	url, err := serveSite()
	if err != nil {
		fmt.Println("启动内置站点服务失败:", err)
		jadeview.Exit()
		return ""
	}
	fmt.Println("[app-ready] 站点 URL:", url)

	// 按 DESIGN.md（Fluent 2）建窗：title-overlay 内置窗口控制按钮 + 透明窗口 + Mica 材质
	opts := jadeview.DefaultWindowOptions()
	opts.Title = "JadeView Go Demo (" + runtime.GOOS + "/" + runtime.GOARCH + ")"
	opts.Width = 1000
	opts.Height = 720
	opts.MinWidth = 640
	opts.MinHeight = 480
	opts.FrameStyle = jadeview.FrameStyle.TitleOverlay // 保留边框 + 无标题栏 + 内置控制按钮
	opts.Transparent = true                            // 配合 backdrop 材质
	opts.Theme = jadeview.Theme.System
	opts.AutoSaveState = true
	mainWindowID = jadeview.CreateWindow(url, 0, &opts, nil)
	if mainWindowID == 0 {
		fmt.Println("创建窗口失败")
		jadeview.Exit()
		return ""
	}
	fmt.Printf("[app-ready] 窗口创建成功 id=%d\n", mainWindowID)

	// 主背景 Mica；标题栏覆盖层高 40，与页面 .title-bar / #app 网格行高保持一致。
	// 图标色初始按浅色主题给，前端探测到实际明暗后经 apply-titlebar 通道再同步。
	jadeview.SetBackdrop(mainWindowID, jadeview.Backdrop.Mica)
	jadeview.SetTitlebarOverlayStyle(mainWindowID, titlebarHeight, "#1A1A1A", "#E5E5E5")

	setupTray()
	return ""
}

// serveSite 在 127.0.0.1 随机端口上起 HTTP 服务，直接以 go:embed 的 site/ 为根，
// 返回入口 URL。前端文件全程只存在于 exe 内，磁盘上没有任何副本。
// 仅监听回环地址；生产环境若担心本机其它进程访问，可在 Handler 外再校验随机 token。
func serveSite() (string, error) {
	sub, err := fs.Sub(siteFS, "site")
	if err != nil {
		return "", err
	}
	ln, err := net.Listen("tcp", "127.0.0.1:0") // 随机空闲端口
	if err != nil {
		return "", err
	}
	go func() {
		srv := &http.Server{Handler: http.FileServer(http.FS(sub))}
		_ = srv.Serve(ln) // 进程退出即销毁，无需优雅关闭
	}()
	return "http://" + ln.Addr().String() + "/index.html", nil
}

// jsonStr / jsonNum 从 invoke 的 payload JSON 里取字段（容错：非 JSON/缺字段返回零值）。
func jsonStr(data, key string) string {
	var m map[string]any
	_ = json.Unmarshal([]byte(data), &m)
	v, _ := m[key].(string)
	return v
}

func jsonNum(data, key string) float64 {
	var m map[string]any
	_ = json.Unmarshal([]byte(data), &m)
	v, _ := m[key].(float64)
	return v
}

// setupTray 创建托盘图标与右键菜单。
// Windows 直接用 .ico；Linux 托盘依赖桌面环境的 appindicator 支持，
// 创建失败（如无托盘协议的桌面）时跳过，不影响其余功能。
func setupTray() {
	trayID = jadeview.TrayCreate()
	if trayID == 0 {
		fmt.Println("[托盘] 创建失败（当前桌面环境可能不支持，跳过）")
		return
	}
	jadeview.TraySetIconFromData(trayID, iconICO) // 内存图标，不落盘
	jadeview.TraySetTooltip(trayID, "JadeView Go Demo")
	jadeview.TraySetMenu(trayID, []jadeview.TrayMenuItem{
		{Type: jadeview.TrayItem.Normal, Key: "show", Label: "显示窗口"},
		{Type: jadeview.TrayItem.Normal, Key: "hide", Label: "隐藏窗口"},
		{Type: jadeview.TrayItem.Divider, Key: "sep1"},
		{Type: jadeview.TrayItem.Normal, Key: "quit", Label: "退出", Dangerous: true},
	})
	jadeview.TraySetVisible(trayID, true)
	fmt.Printf("[托盘] 已创建 id=%d\n", trayID)
}

func onTrayMenuCommand(_ uint32, data string) string {
	fmt.Printf("[事件] tray-menu-command: %s\n", data)
	// 载荷为 JSON，优先取被点菜单项的 key 精确匹配；
	// 库未按 {"key":...} 回报时退回对原始载荷的子串匹配兜底。
	key := jsonStr(data, "key")
	if key == "" {
		key = data
	}
	switch {
	case strings.Contains(key, "show"):
		jadeview.SetVisible(mainWindowID, true)
		jadeview.SetFocus(mainWindowID)
	case strings.Contains(key, "hide"):
		jadeview.SetVisible(mainWindowID, false)
	case strings.Contains(key, "quit"):
		jadeview.Exit()
	}
	return ""
}

func registerIPCHandlers() {
	// --- 外观：主题 / 标题栏 / 材质 / 缩放 ---

	// 设置窗口主题：payload {"mode":"Light"|"Dark"|"System"}
	jadeview.RegisterIPCHandler("set-theme", func(wid uint32, data string) string {
		mode := jsonStr(data, "mode")
		if mode == "" {
			return "缺少 mode"
		}
		ok := jadeview.SetTheme(wid, mode)
		return fmt.Sprintf("set_window_theme(%s): %v", mode, ok)
	})

	// 同步标题栏覆盖层图标色：payload {"dark":true|false}（仅样式，按钮功能库内置）
	jadeview.RegisterIPCHandler("apply-titlebar", func(wid uint32, data string) string {
		if strings.Contains(data, "true") {
			jadeview.SetTitlebarOverlayStyle(wid, titlebarHeight, "#FFFFFF", "#3A3A3A")
			return "标题栏: 深色图标方案"
		}
		jadeview.SetTitlebarOverlayStyle(wid, titlebarHeight, "#1A1A1A", "#E5E5E5")
		return "标题栏: 浅色图标方案"
	})

	// 切换窗口材质：payload {"type":"mica"|"micaAlt"|"acrylic"|"none","color":"#RRGGBBAA"}
	jadeview.RegisterIPCHandler("set-backdrop", func(wid uint32, data string) string {
		t := jsonStr(data, "type")
		if t == "none" {
			color := jsonStr(data, "color")
			if color == "" {
				color = "#F3F3F3FF"
			}
			ok := jadeview.SetBackgroundColor(wid, color)
			return fmt.Sprintf("纯色背景 %s: %v", color, ok)
		}
		ok := jadeview.SetBackdrop(wid, t)
		return fmt.Sprintf("set_window_backdrop(%s): %v", t, ok)
	})

	// WebView 缩放：payload {"level":1.25}
	jadeview.RegisterIPCHandler("zoom", func(wid uint32, data string) string {
		level := jsonNum(data, "level")
		if level <= 0 {
			level = 1.0
		}
		ok := jadeview.SetZoom(wid, level)
		return fmt.Sprintf("set_webview_zoom(%.2f): %v", level, ok)
	})

	// --- IPC 测试 ---

	// 回声：原样返回 payload + 服务端时间戳
	jadeview.RegisterIPCHandler("ipc-echo", func(_ uint32, data string) string {
		return fmt.Sprintf(`{"echo":%s,"serverTime":%q}`, data, time.Now().Format("15:04:05.000"))
	})

	// 后端主动推送：异步连发 3 条 push-demo 消息（验证宿主→前端方向）
	jadeview.RegisterIPCHandler("ipc-push", func(wid uint32, _ string) string {
		go func() {
			for i := 1; i <= 3; i++ {
				time.Sleep(400 * time.Millisecond)
				jadeview.SendIPCMessage(wid, "push-demo",
					fmt.Sprintf(`{"seq":%d,"time":%q}`, i, time.Now().Format("15:04:05.000")))
			}
		}()
		return "已排队 3 条推送（间隔 400ms）"
	})

	// Toast 契约演示：按 DESIGN.md §12 的 payload 规范经 "toast" 事件推送
	jadeview.RegisterIPCHandler("demo-toast", func(wid uint32, data string) string {
		level := jsonStr(data, "level")
		if level == "" {
			level = "info"
		}
		payload := map[string]any{
			"level":   level,
			"title":   map[string]string{"info": "提示", "success": "已完成", "warning": "请注意", "error": "出错了"}[level],
			"message": fmt.Sprintf("这是一条来自 Go 宿主的 %s 通知（%s）。", level, time.Now().Format("15:04:05")),
		}
		if level == "error" {
			payload["duration"] = 0 // 错误不自动消失
		}
		b, _ := json.Marshal(payload)
		ok := jadeview.SendIPCMessage(wid, "toast", string(b))
		return fmt.Sprintf("toast 已推送: %v", ok)
	})

	// --- 系统信息 ---
	jadeview.RegisterIPCHandler("sysinfo", func(_ uint32, _ string) string {
		locale, _ := jadeview.GetLocale()
		wv, _ := jadeview.GetWebViewVersion()
		info := fmt.Sprintf("平台 %s/%s | JadeView %s | WebView %s | 语言 %s",
			runtime.GOOS, runtime.GOARCH, jadeview.Version(), wv, locale)
		if runtime.GOOS == "windows" {
			info += fmt.Sprintf(" | Win11: %v", jadeview.IsWindows11())
		}
		return info
	})
	jadeview.RegisterIPCHandler("displays", func(_ uint32, _ string) string {
		v, _ := jadeview.GetDisplaysInfo()
		if len(v) > 160 {
			v = v[:160] + "..."
		}
		return v
	})

	// --- 异步对话框：弹框后立即返回，结果经 SendIPCMessage 推回页面 ---
	jadeview.RegisterIPCHandler("open-file", func(wid uint32, _ string) string {
		jadeview.ShowOpenDialogAsync(jadeview.FileDialogParams{WindowID: wid, Title: "选择文件"},
			func(result string) { jadeview.SendIPCMessage(wid, "dialog-result", result) })
		return "打开对话框已弹出"
	})
	jadeview.RegisterIPCHandler("save-file", func(wid uint32, _ string) string {
		jadeview.ShowSaveDialogAsync(jadeview.FileDialogParams{WindowID: wid, Title: "保存为"},
			func(result string) { jadeview.SendIPCMessage(wid, "dialog-result", result) })
		return "保存对话框已弹出"
	})
	jadeview.RegisterIPCHandler("msgbox", func(wid uint32, _ string) string {
		jadeview.ShowMessageBoxAsync(jadeview.MessageBoxParams{
			WindowID: wid, Title: "提示", Message: "这是一个来自 Go 的消息框",
			Type: jadeview.MsgBoxType.Info, Buttons: `["确定","取消"]`,
		}, func(result string) { jadeview.SendIPCMessage(wid, "dialog-result", result) })
		return "消息框已弹出"
	})
	jadeview.RegisterIPCHandler("notify", func(_ uint32, _ string) string {
		ok := jadeview.ShowNotification(jadeview.NotificationParams{
			Summary: "JadeView Go", Body: "这是一条系统通知", Timeout: -1,
		})
		return fmt.Sprintf("通知发送: %v", ok)
	})

	// --- YAML 持久化（存于 Init 设置的 data_directory，须在 app-ready 后调用） ---
	jadeview.RegisterIPCHandler("yaml-set", func(_ uint32, _ string) string {
		rc := jadeview.YAMLSet("demo", "app.name", "JadeView-Go")
		rc2 := jadeview.YAMLSet("demo", "app.updated", time.Now().Format("2006-01-02 15:04:05"))
		return fmt.Sprintf("写入 rc=%d,%d", rc, rc2)
	})
	jadeview.RegisterIPCHandler("yaml-get", func(_ uint32, _ string) string {
		v, _ := jadeview.YAMLGet("demo", "app.name")
		return "app.name = " + v
	})
	jadeview.RegisterIPCHandler("yaml-all", func(_ uint32, _ string) string {
		v, _ := jadeview.YAMLGetAll("demo")
		return v
	})

	// --- 剪贴板 / NTP 网络时间 ---
	jadeview.RegisterIPCHandler("clip-write", func(_ uint32, _ string) string {
		text := "JadeView " + time.Now().Format("15:04:05")
		ok := jadeview.ClipboardWriteText(text)
		return fmt.Sprintf("写入剪贴板 %q: %v", text, ok)
	})
	jadeview.RegisterIPCHandler("clip-read", func(_ uint32, _ string) string {
		v, ok := jadeview.ClipboardReadText()
		if !ok {
			return "剪贴板读取失败"
		}
		return "剪贴板: " + v
	})
	jadeview.RegisterIPCHandler("ntp", func(_ uint32, _ string) string {
		// 空字符串 = 使用内置服务器列表逐个尝试；返回 UTC 毫秒，-1 失败
		ms := jadeview.NTPNow("")
		if ms < 0 {
			return "NTP 获取失败（检查 UDP/123 出站网络）"
		}
		utc := time.UnixMilli(ms).UTC()
		return fmt.Sprintf("NTP UTC: %s（北京时间 %s）",
			utc.Format("15:04:05"), utc.Add(8*time.Hour).Format("15:04:05"))
	})

	// --- 窗口操作 ---
	jadeview.RegisterIPCHandler("toggle-top", func(wid uint32, _ string) string {
		alwaysOnTop = !alwaysOnTop
		jadeview.SetAlwaysOnTop(wid, alwaysOnTop)
		return fmt.Sprintf("置顶: %v", alwaysOnTop)
	})
	jadeview.RegisterIPCHandler("minimize", func(wid uint32, _ string) string {
		jadeview.Minimize(wid)
		return "已最小化"
	})
	jadeview.RegisterIPCHandler("bounds", func(wid uint32, _ string) string {
		v, _ := jadeview.GetWindowBounds(wid)
		return "窗口边界: " + v
	})
	jadeview.RegisterIPCHandler("hwnd", func(wid uint32, _ string) string {
		// beta.9 起所有窗口可取句柄；GetWindowID 反查验证一致性。
		// Linux 下句柄语义由库决定（可能为 0），此处仅作演示。
		h := jadeview.GetWindowHWND(wid)
		back := jadeview.GetWindowID(h)
		return fmt.Sprintf("HWND=0x%X → 反查窗口ID=%d（当前=%d）", h, back, wid)
	})
	jadeview.RegisterIPCHandler("fullscreen", func(wid uint32, _ string) string {
		to := !jadeview.IsFullscreen(wid)
		jadeview.SetFullscreen(wid, to)
		return fmt.Sprintf("全屏: %v", to)
	})
	jadeview.RegisterIPCHandler("devtools", func(wid uint32, _ string) string {
		if jadeview.IsDevtoolsOpen(wid) {
			jadeview.CloseDevtools(wid)
			return "DevTools 已关闭"
		}
		jadeview.OpenDevtools(wid)
		return "DevTools 已打开"
	})
	jadeview.RegisterIPCHandler("flash", func(wid uint32, _ string) string {
		// Windows 闪烁任务栏；Linux 对应 urgency 提示，依桌面环境而定
		return fmt.Sprintf("闪烁窗口: %v", jadeview.FlashWindow(wid, 5))
	})
}
