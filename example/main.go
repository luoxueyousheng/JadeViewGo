//go:build linux || windows

// JadeView Go 封装 · 跨平台可交互示例（Windows / Linux 通用，单一代码库）。
//
// 演示内容：窗口、事件桥、IPC 双向通信、异步对话框、系统通知、托盘菜单、
// YAML 持久化、剪贴板、NTP 网络时间、HWND/窗口ID 互查、显示器/系统信息。
// 平台差异统一用 runtime.GOOS 分支处理，两端均可直接构建运行。
//
// 运行（详见 README「各平台使用方式」）：
//
//	Windows: go run ./example            （需 MinGW-w64，产物单 exe 自包含 DLL）
//	Linux  : go run ./example            （需 gcc + GTK3/WebKit2GTK，静态链接）
//
// 关键时序（官方要求）：
//  1. app-ready 必须在 Init 之前注册；
//  2. 建窗、协议服务、托盘等都放在 app-ready 回调里；
//  3. 回调里必须判断 windowID==1 才算初始化成功。
package main

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	jadeview "github.com/luoxueyousheng/JadeViewGo"
)

//go:embed index.html
var indexHTML []byte

//go:embed icon.ico
var iconICO []byte

var (
	mainWindowID uint32
	trayID       uint32
	alwaysOnTop  bool
	assetDir     string
)

func main() {
	fmt.Printf("平台: %s/%s\n", runtime.GOOS, runtime.GOARCH)

	// 内置页面与图标释放到临时目录，供协议服务/托盘使用。
	assetDir = filepath.Join(os.TempDir(), "jadeview-demo")
	_ = os.MkdirAll(assetDir, 0o755)
	_ = os.WriteFile(filepath.Join(assetDir, "index.html"), indexHTML, 0o644)
	_ = os.WriteFile(filepath.Join(assetDir, "icon.ico"), iconICO, 0o644)

	// 1) app-ready 之前注册事件
	jadeview.On("app-ready", onAppReady)
	jadeview.On("window-created", func(id uint32, _ string) string {
		fmt.Printf("[事件] window-created id=%d\n", id)
		return ""
	})
	jadeview.On("window-closed", func(id uint32, _ string) string {
		fmt.Printf("[事件] window-closed id=%d\n", id)
		return ""
	})
	jadeview.On("window-all-closed", func(_ uint32, _ string) string {
		fmt.Println("[事件] 所有窗口已关闭，退出")
		jadeview.Exit()
		return ""
	})
	jadeview.On("theme-changed", func(_ uint32, data string) string {
		fmt.Printf("[事件] 系统主题变化: %s\n", data)
		return ""
	})
	jadeview.On("tray-menu-command", onTrayMenuCommand)
	jadeview.On("tray-event", func(_ uint32, data string) string {
		fmt.Printf("[事件] tray-event: %s\n", data)
		return ""
	})

	// 2) IPC 处理器（前端 jade.invoke 的目标）
	registerIPCHandlers()

	// 3) 初始化（版本号在 Init 之后才可取）
	if !jadeview.Init(true, "", assetDir, "jadeview-go-demo", "jadeview-go-demo-signature", false) {
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

	// 协议服务：把临时目录作为站点根，返回自定义协议 URL
	base, ok := jadeview.SetProtocolServicePath(assetDir, true)
	if !ok {
		fmt.Println("协议服务设置失败")
		jadeview.Exit()
		return ""
	}
	url := strings.TrimRight(base, "/") + "/index.html"
	fmt.Println("[app-ready] 站点 URL:", url)

	opts := jadeview.DefaultWindowOptions()
	opts.Title = "JadeView Go Demo (" + runtime.GOOS + "/" + runtime.GOARCH + ")"
	opts.Width = 760
	opts.Height = 680
	mainWindowID = jadeview.CreateWindow(url, 0, &opts, nil)
	if mainWindowID == 0 {
		fmt.Println("创建窗口失败")
		jadeview.Exit()
		return ""
	}
	fmt.Printf("[app-ready] 窗口创建成功 id=%d\n", mainWindowID)

	setupTray()
	return ""
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
	jadeview.TraySetIconFromFile(trayID, filepath.Join(assetDir, "icon.ico"))
	jadeview.TraySetTooltip(trayID, "JadeView Go Demo")
	jadeview.TraySetMenu(trayID, []jadeview.TrayMenuItem{
		{Type: jadeview.TrayItemNormal, Key: "show", Label: "显示窗口"},
		{Type: jadeview.TrayItemNormal, Key: "hide", Label: "隐藏窗口"},
		{Type: jadeview.TrayItemDivider, Key: "sep1"},
		{Type: jadeview.TrayItemNormal, Key: "quit", Label: "退出", Dangerous: true},
	})
	jadeview.TraySetVisible(trayID, true)
	fmt.Printf("[托盘] 已创建 id=%d\n", trayID)
}

func onTrayMenuCommand(_ uint32, data string) string {
	fmt.Printf("[事件] tray-menu-command: %s\n", data)
	switch {
	case strings.Contains(data, "show"):
		jadeview.SetVisible(mainWindowID, true)
		jadeview.SetFocus(mainWindowID)
	case strings.Contains(data, "hide"):
		jadeview.SetVisible(mainWindowID, false)
	case strings.Contains(data, "quit"):
		jadeview.Exit()
	}
	return ""
}

func registerIPCHandlers() {
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
			Type: "info", Buttons: `["确定","取消"]`,
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
