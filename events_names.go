package jadeview

// 事件名常量，与 include/JadeView.h 的 JADEVIEW_EVENT_* 宏一一对应，
// 供 On/Off 使用，避免裸写字符串拼错（Windows/Linux 两侧共用）。
//
//	jadeview.On(jadeview.EventAppReady, handler)
const (
	EventAppReady                 = "app-ready"                    // 应用初始化完成（回调里判断 windowID==1 才算成功）
	EventContextMenu              = "context-menu"                 // 右键菜单（配合 SetContextMenuItems）
	EventCrash                    = "crash"                        // 程序崩溃（event_data 为 Crash* 错误代码）
	EventDragDrop                 = "drag-drop"                    // 拖拽事件（enter/over/drop/leave）
	EventGlobalHotkey             = "global-hotkey"                // 全局热键触发
	EventJapkLoadFailed           = "japk-load-failed"             // JAPK 资源包加载失败
	EventJavascriptResult         = "javascript-result"            // ExecuteJavaScript 的执行结果
	EventMenuItemClicked          = "menu-item-clicked"            // 菜单项点击
	EventNotificationAction       = "notification-action"          // 通知操作被点击
	EventNotificationDismissed    = "notification-dismissed"       // 通知已关闭
	EventNotificationFailed       = "notification-failed"          // 通知显示失败
	EventNotificationShown        = "notification-shown"           // 通知已显示
	EventPostMessageReceived      = "postmessage-received"         // 收到前端 postMessage
	EventSecondInstance           = "second-instance"              // 第二个实例启动（单实例模式）
	EventThemeChanged             = "theme-changed"                // 系统主题变化
	EventTrayEvent                = "tray-event"                   // 托盘图标事件
	EventTrayMenuCommand          = "tray-menu-command"            // 托盘菜单命令
	EventUpdateWindowIcon         = "update-window-icon"           // 更新窗口图标
	EventWebViewDidFinishLoad     = "webview-did-finish-load"      // WebView 加载完成
	EventWebViewDidStartLoading   = "webview-did-start-loading"    // WebView 开始加载
	EventWebViewDownloadCompleted = "webview-download-completed"   // WebView 下载完成
	EventWebViewFaviconUpdated    = "webview-page-favicon-updated" // WebView 页面图标更新
	EventWebViewTitleUpdated      = "webview-page-title-updated"   // WebView 页面标题更新
	EventWindowAllClosed          = "window-all-closed"            // 所有窗口已关闭
	EventWindowBlurred            = "window-blurred"               // 窗口失去焦点
	EventWindowBounds             = "window-bounds"                // 窗口边界变化
	EventWindowClosed             = "window-closed"                // 窗口关闭
	EventWindowCreated            = "window-created"               // 窗口创建完成
	EventWindowDestroyed          = "window-destroyed"             // 窗口销毁
	EventWindowFocused            = "window-focused"               // 窗口获得焦点
	EventWindowFullscreen         = "window-fullscreen"            // 窗口全屏状态变化
	EventWindowMoved              = "window-moved"                 // 窗口移动
	EventWindowResized            = "window-resized"               // 窗口大小改变
	EventWindowStateChanged       = "window-state-changed"         // 窗口状态变化（最大化/还原等）
)

// 崩溃错误代码（EventCrash 事件的 event_data），对应头文件 JADEVIEW_CRASH_* 宏。
const (
	CrashSEHAccessViolation    = "SEH_ACCESS_VIOLATION"      // SEH：内存访问违规
	CrashSEHStackOverflow      = "SEH_STACK_OVERFLOW"        // SEH：栈溢出
	CrashSEHIllegalInstruction = "SEH_ILLEGAL_INSTRUCTION"   // SEH：非法指令
	CrashSEHInvalidHandle      = "SEH_INVALID_HANDLE"        // SEH：无效句柄
	CrashSEHUnknown            = "SEH_UNKNOWN"               // SEH：未知
	CrashRuntimePanic          = "RUNTIME_PANIC"             // 运行时 Panic
	CrashWV2BrowserExited      = "WV2_BROWSER_EXITED"        // WebView2：浏览器进程已退出
	CrashWV2RenderExited       = "WV2_RENDER_EXITED"         // WebView2：渲染进程已退出
	CrashWV2RenderUnresponsive = "WV2_RENDER_UNRESPONSIVE"   // WebView2：渲染进程无响应
	CrashWV2FrameRenderExited  = "WV2_FRAME_RENDER_EXITED"   // WebView2：框架渲染进程已退出
	CrashWV2UtilityExited      = "WV2_UTILITY_EXITED"        // WebView2：工具进程已退出
	CrashWV2SandboxExited      = "WV2_SANDBOX_HELPER_EXITED" // WebView2：沙箱辅助进程已退出
	CrashWV2GPUExited          = "WV2_GPU_EXITED"            // WebView2：GPU 进程已退出
	CrashWV2PPAPIPluginExited  = "WV2_PPAPI_PLUGIN_EXITED"   // WebView2：PPAPI 插件进程已退出
	CrashWV2PPAPIBrokerExited  = "WV2_PPAPI_BROKER_EXITED"   // WebView2：PPAPI 代理进程已退出
	CrashWV2UnknownExited      = "WV2_UNKNOWN_EXITED"        // WebView2：未知进程已退出
)
