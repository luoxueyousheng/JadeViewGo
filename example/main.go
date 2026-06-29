//go:build linux || windows

// 完整示例：窗口 + 事件 + IPC + 异步对话框 + 托盘 + YAML 串起来的可交互 demo。
//
// 运行：
//
//	cd JadeView && go run ./example
//
// 交互方式：
//   - 页面上的按钮 → jade.invoke(channel) → 下面 RegisterIPCHandler 处理 → 结果回传页面。
//   - 异步对话框（打开/保存/消息框）的结果通过 SendIPCMessage 推回页面的 jade.on('dialog-result')。
//   - 托盘右键菜单可显示/隐藏窗口、退出。
//
// 关键时序（官方）：app-ready 必须在 Init 之前注册；窗口、协议服务、托盘都在 app-ready 回调里建。
package main

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	// 把内置的页面与图标释放到临时目录，供协议服务/托盘使用。
	assetDir = filepath.Join(os.TempDir(), "jadeview-demo")
	_ = os.MkdirAll(assetDir, 0o755)
	_ = os.WriteFile(filepath.Join(assetDir, "index.html"), indexHTML, 0o644)
	iconPath := filepath.Join(assetDir, "icon.ico")
	_ = os.WriteFile(iconPath, iconICO, 0o644)

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
	jadeview.On("tray-menu-command", onTrayMenuCommand)
	jadeview.On("tray-event", func(_ uint32, data string) string {
		fmt.Printf("[事件] tray-event: %s\n", data)
		return ""
	})

	// 2) IPC 处理器（前端 jade.invoke 的目标）
	registerIPCHandlers()

	// 3) 初始化
	if !jadeview.Init(true, "", assetDir, "jadeview-go-demo", "jadeview-go-demo-signature", false) {
		fmt.Println("Init 失败")
		return
	}

	// 4) 阻塞运行
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

	// 用协议服务把临时目录作为站点根，返回可访问 URL
	base, ok := jadeview.SetProtocolServicePath(assetDir, true)
	if !ok {
		fmt.Println("协议服务设置失败，回退到 https://jade.run/")
		base = "https://jade.run/"
	}
	url := strings.TrimRight(base, "/") + "/index.html"
	fmt.Println("[app-ready] 站点 URL:", url)

	opts := jadeview.DefaultWindowOptions()
	opts.Title = "JadeView Go Demo"
	opts.Width = 720
	opts.Height = 640
	mainWindowID = jadeview.CreateWindow(url, 0, &opts, nil)
	if mainWindowID == 0 {
		fmt.Println("创建窗口失败")
		jadeview.Exit()
		return ""
	}
	fmt.Printf("[app-ready] 窗口创建成功 id=%d\n", mainWindowID)

	// 托盘
	setupTray()
	return ""
}

func setupTray() {
	trayID = jadeview.TrayCreate()
	if trayID == 0 {
		fmt.Println("托盘创建失败（跳过）")
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
	// 异步对话框：弹框后立即返回，结果通过 SendIPCMessage 回推页面
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

	// YAML
	jadeview.RegisterIPCHandler("yaml-set", func(_ uint32, _ string) string {
		rc := jadeview.YAMLSet("demo", "app.name", "JadeView-Go")
		rc2 := jadeview.YAMLSet("demo", "ui.theme", "dark")
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

	// 窗口 / 系统
	jadeview.RegisterIPCHandler("toggle-top", func(wid uint32, _ string) string {
		alwaysOnTop = !alwaysOnTop
		jadeview.SetAlwaysOnTop(wid, alwaysOnTop)
		return fmt.Sprintf("置顶: %v", alwaysOnTop)
	})
	jadeview.RegisterIPCHandler("minimize", func(wid uint32, _ string) string {
		jadeview.Minimize(wid)
		return "已最小化"
	})
	jadeview.RegisterIPCHandler("devtools", func(wid uint32, _ string) string {
		if jadeview.IsDevtoolsOpen(wid) {
			jadeview.CloseDevtools(wid)
			return "DevTools 已关闭"
		}
		jadeview.OpenDevtools(wid)
		return "DevTools 已打开"
	})
	jadeview.RegisterIPCHandler("displays", func(_ uint32, _ string) string {
		v, _ := jadeview.GetDisplaysInfo()
		if len(v) > 120 {
			v = v[:120] + "..."
		}
		return v
	})
	jadeview.RegisterIPCHandler("webview-ver", func(_ uint32, _ string) string {
		v, _ := jadeview.GetWebViewVersion()
		return "WebView " + v
	})
}
