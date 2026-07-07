//go:build windows

package jadeview

// Windows 纯 Go 实现的核心：DLL 定位/释放/加载与调用辅助函数。
//
// 不使用 cgo——go:embed 内置 MSVC 版 JadeView.dll，运行时释放到临时目录，
// syscall.NewLazyDLL 按绝对路径加载，各 API 经 LazyProc.Call 直调。
// 构建只需 Go 工具链，amd64/386/arm64 三架构均为单 exe 自包含分发。
//
// 结构体布局、回调约定与 C 头文件的对应关系见各 *_windows.go；
// 布局已用 C(offsetof) 与 Go(unsafe.Offsetof) 双端比对验证（amd64/386 逐字段一致，
// arm64 与 amd64 对齐规则相同）。

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
	"unsafe"
)

// jadeDLL 惰性加载 JadeView.dll：首次调用任一 API 时才真正 LoadLibrary。
// 路径解析优先级：exe 同目录的 JadeView.dll（便于手动覆盖/调试）→
// 释放内置副本到 %TEMP%\jadeview\<arch>-<内容哈希前8位>\（内容寻址，
// 不同版本/架构各占一目录，互不覆盖，也无需按字节比对旧文件）。
var jadeDLL = syscall.NewLazyDLL(resolveDLLPath())

func resolveDLLPath() string {
	if exe, err := os.Executable(); err == nil {
		p := filepath.Join(filepath.Dir(exe), "JadeView.dll")
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			return p
		}
	}
	sum := sha256.Sum256(embeddedJadeViewDLL)
	dir := filepath.Join(os.TempDir(), "jadeview", runtime.GOARCH+"-"+hex.EncodeToString(sum[:4]))
	dst := filepath.Join(dir, "JadeView.dll")
	if st, err := os.Stat(dst); err != nil || st.Size() != int64(len(embeddedJadeViewDLL)) {
		_ = os.MkdirAll(dir, 0o755)
		// 写失败不致命（可能另一进程刚写完并占用）；加载失败会在首次调用时报错
		_ = os.WriteFile(dst, embeddedJadeViewDLL, 0o644)
	}
	return dst
}

// i32 截断 C int32 返回值（x64 下 RAX 高 32 位未定义，必须先截断再判断）。
func i32(r uintptr) int32 { return int32(uint32(r)) }

// b2u 把 Go bool 转成库约定的 int32(0/1) 参数。
func b2u(b bool) uintptr {
	if b {
		return 1
	}
	return 0
}

// cstrs 管理一批调用期间存活的 NUL 结尾 UTF-8 字符串副本。
// 调用后需 runtime.KeepAlive(pool) 保证 GC 不提前回收。
type cstrs struct {
	hold [][]byte
}

// p 返回字符串的 C 指针；空字符串映射为 NULL（库会使用对应默认值）。
func (c *cstrs) p(s string) *byte {
	if s == "" {
		return nil
	}
	return c.pAlways(s)
}

// pAlways 同 p，但空字符串返回指向 "\0" 的有效指针（对应 cgo C.CString("")）。
func (c *cstrs) pAlways(s string) *byte {
	b := append([]byte(s), 0)
	c.hold = append(c.hold, b)
	return &b[0]
}

// foreignPtr 把指向「非 Go 堆内存」（DLL 内部分配/VirtualAlloc）的地址转成
// unsafe.Pointer。该内存不受 Go GC 管理，转换安全；经指针变量中转以规避
// vet 对 uintptr→Pointer 直转的误报。
func foreignPtr(p uintptr) unsafe.Pointer {
	return *(*unsafe.Pointer)(unsafe.Pointer(&p))
}

// goString 把 NUL 结尾的 C 字符串指针转成 Go string（0 → ""）。
func goString(p uintptr) string {
	if p == 0 {
		return ""
	}
	base := foreignPtr(p)
	n := 0
	for *(*byte)(unsafe.Add(base, n)) != 0 {
		n++
	}
	if n == 0 {
		return ""
	}
	return string(unsafe.Slice((*byte)(base), n))
}

// goStringFree 转换库返回的 char*（约定用 jade_text_free 释放）并释放。
func goStringFree(p uintptr) string {
	if p == 0 {
		return ""
	}
	s := goString(p)
	procJadeTextFree.Call(p)
	return s
}

// bufCallSize 调用「填充缓冲区、size_t 长度、返回 1=成功」类函数。
// 缓冲区不足时（库返回 0）自动放大重试一次。语义与 Linux 版 helpers 一致。
func bufCallSize(initial int, fn func(buf unsafe.Pointer, size uintptr) int32) (string, bool) {
	for _, size := range []int{initial, initial * 16} {
		buf := make([]byte, size)
		if fn(unsafe.Pointer(&buf[0]), uintptr(size)) == 1 {
			return cBufToString(buf), true
		}
	}
	return "", false
}

// bufCallInt 同 bufCallSize，但长度参数为 int32。
func bufCallInt(initial int, fn func(buf unsafe.Pointer, size int32) int32) (string, bool) {
	for _, size := range []int{initial, initial * 16} {
		buf := make([]byte, size)
		if fn(unsafe.Pointer(&buf[0]), int32(size)) == 1 {
			return cBufToString(buf), true
		}
	}
	return "", false
}

// yamlGet 执行 YAML 两阶段查询：先取所需字节数再分配读取。
// 返回 (值, 状态码)；状态码 1=成功，其余见文档（0=不存在，-1 IO，-2 类型，-4 解析）。
func yamlGet(fn func(buf unsafe.Pointer, size uintptr) int32) (string, int32) {
	n := fn(nil, 0)
	if n < 2 {
		if n == 1 {
			return "", 1
		}
		return "", n
	}
	buf := make([]byte, int(n))
	rc := fn(unsafe.Pointer(&buf[0]), uintptr(len(buf)))
	if rc == 1 {
		return cBufToString(buf), 1
	}
	return "", rc
}

// JadeView.dll 全部导出函数的惰性地址表（GetProcAddress 按需解析）。
// 函数名与 include/JadeView.h 一一对应；DLL 导出为无修饰名，三架构一致。

var (
	// --- lifecycle ---
	procJadeViewInit   = jadeDLL.NewProc("JadeView_init")
	procRunMessageLoop = jadeDLL.NewProc("run_message_loop")
	procJadeviewExit   = jadeDLL.NewProc("jadeview_exit")

	// --- japk ---
	procSetPublicKey    = jadeDLL.NewProc("JadeView_set_public_key")
	procLoadFromBytes   = jadeDLL.NewProc("JadeView_load_from_bytes")
	procIsLoaded        = jadeDLL.NewProc("JadeView_is_loaded")
	procGetAppSignature = jadeDLL.NewProc("JadeView_get_app_signature")
	procGetSignatureInf = jadeDLL.NewProc("JadeView_get_signature_info")
	procUnload          = jadeDLL.NewProc("JadeView_unload")

	// --- window ---
	procCreateWebviewWindow    = jadeDLL.NewProc("create_webview_window")
	procCreateBorderlessWindow = jadeDLL.NewProc("create_borderless_webview_window")
	procGetWindowHwnd          = jadeDLL.NewProc("get_window_hwnd")
	procGetWindowID            = jadeDLL.NewProc("get_window_id")
	procNavigateToURL          = jadeDLL.NewProc("navigate_to_url")
	procReloadWebviewWindow    = jadeDLL.NewProc("reload_webview_window")
	procExecuteJavascript      = jadeDLL.NewProc("execute_javascript")
	procSetWindowTitle         = jadeDLL.NewProc("set_window_title")
	procSetWindowSize          = jadeDLL.NewProc("set_window_size")
	procSetWindowPosition      = jadeDLL.NewProc("set_window_position")
	procSetWindowVisible       = jadeDLL.NewProc("set_window_visible")
	procSetWindowFocus         = jadeDLL.NewProc("set_window_focus")
	procSetWindowAlwaysOnTop   = jadeDLL.NewProc("set_window_always_on_top")
	procSetWindowSkipTaskbar   = jadeDLL.NewProc("set_window_skip_taskbar")
	procSetWindowNoActivate    = jadeDLL.NewProc("set_window_no_activate")
	procSetWindowLevel         = jadeDLL.NewProc("set_window_level")
	procCloseWindow            = jadeDLL.NewProc("close_window")
	procMinimizeWindow         = jadeDLL.NewProc("minimize_window")
	procToggleMaximizeWindow   = jadeDLL.NewProc("toggle_maximize_window")
	procIsWindowMaximized      = jadeDLL.NewProc("is_window_maximized")
	procSetContentProtection   = jadeDLL.NewProc("set_content_protection")
	procSetWebviewZoom         = jadeDLL.NewProc("set_webview_zoom")
	procSetWindowFrameStyle    = jadeDLL.NewProc("set_window_frame_style")
	procGetWindowCount         = jadeDLL.NewProc("get_window_count")
	procJadeOn                 = jadeDLL.NewProc("jade_on")
	procJadeOff                = jadeDLL.NewProc("jade_off")
	procRegisterIPCHandler     = jadeDLL.NewProc("register_ipc_handler")
	procSetWindowTheme         = jadeDLL.NewProc("set_window_theme")
	procSetTitlebarOverlay     = jadeDLL.NewProc("set_titlebar_overlay_style")
	procSetWindowEnabled       = jadeDLL.NewProc("set_window_enabled")
	procGetWindowTheme         = jadeDLL.NewProc("get_window_theme")
	procRequestRedraw          = jadeDLL.NewProc("request_redraw")
	procSendIPCMessage         = jadeDLL.NewProc("send_ipc_message")
	procSetWindowBackdrop      = jadeDLL.NewProc("set_window_backdrop")
	procSetWindowBgColor       = jadeDLL.NewProc("set_window_background_color")
	procSetWindowFullscreen    = jadeDLL.NewProc("set_window_fullscreen")
	procGetWindowBounds        = jadeDLL.NewProc("get_window_bounds")
	procIsWindowMinimized      = jadeDLL.NewProc("is_window_minimized")
	procIsWindowVisible        = jadeDLL.NewProc("is_window_visible")
	procIsWindowFocused        = jadeDLL.NewProc("is_window_focused")
	procIsWindowFullscreen     = jadeDLL.NewProc("is_window_fullscreen")
	procSetWindowMinSize       = jadeDLL.NewProc("set_window_min_size")
	procSetWindowMaxSize       = jadeDLL.NewProc("set_window_max_size")
	procSetWindowResizable     = jadeDLL.NewProc("set_window_resizable")
	procSetIgnoreCursorEvents  = jadeDLL.NewProc("set_window_ignore_cursor_events")
	procGetWebviewURL          = jadeDLL.NewProc("get_webview_url")
	procOpenDevtools           = jadeDLL.NewProc("open_devtools")
	procCloseDevtools          = jadeDLL.NewProc("close_devtools")
	procIsDevtoolsOpen         = jadeDLL.NewProc("is_devtools_open")
	procClearBrowsingData      = jadeDLL.NewProc("clear_browsing_data")
	procSetWindowProgress      = jadeDLL.NewProc("set_window_progress")
	procFlashWindow            = jadeDLL.NewProc("flash_window")
	procShowAboutDialog        = jadeDLL.NewProc("show_about_dialog")

	// --- tray ---
	procTrayCreate          = jadeDLL.NewProc("tray_create")
	procTrayDestroy         = jadeDLL.NewProc("tray_destroy")
	procTraySetVisible      = jadeDLL.NewProc("tray_set_visible")
	procTraySetTooltip      = jadeDLL.NewProc("tray_set_tooltip")
	procTraySetIconFromFile = jadeDLL.NewProc("tray_set_icon_from_file")
	procTraySetMenuItems    = jadeDLL.NewProc("tray_set_menu_items")
	procTraySetIconFromData = jadeDLL.NewProc("set_tray_icon_from_data")

	// --- dialog / menu ---
	procShowNotification    = jadeDLL.NewProc("show_notification")
	procShowOpenDialog      = jadeDLL.NewProc("jade_dialog_show_open_dialog")
	procShowSaveDialog      = jadeDLL.NewProc("jade_dialog_show_save_dialog")
	procShowMessageBox      = jadeDLL.NewProc("jade_dialog_show_message_box")
	procShowErrorBox        = jadeDLL.NewProc("jade_dialog_show_error_box")
	procShowOpenDialogAsync = jadeDLL.NewProc("jade_dialog_show_open_dialog_async")
	procShowSaveDialogAsync = jadeDLL.NewProc("jade_dialog_show_save_dialog_async")
	procShowMessageBoxAsync = jadeDLL.NewProc("jade_dialog_show_message_box_async")
	procMenuItemCreate      = jadeDLL.NewProc("jade_menu_item_create")
	procMenuItemSetEnabled  = jadeDLL.NewProc("jade_menu_item_set_enabled")
	procMenuItemSetChecked  = jadeDLL.NewProc("jade_menu_item_set_checked")
	procSetContextMenuItems = jadeDLL.NewProc("jade_set_context_menu_items")
	procMenuItemDestroy     = jadeDLL.NewProc("jade_menu_item_destroy")

	// --- yaml ---
	procYamlSet        = jadeDLL.NewProc("yaml_set")
	procYamlGet        = jadeDLL.NewProc("yaml_get")
	procYamlSetStr     = jadeDLL.NewProc("yaml_set_str")
	procYamlGetAll     = jadeDLL.NewProc("yaml_get_all")
	procYamlHas        = jadeDLL.NewProc("yaml_has")
	procYamlDelete     = jadeDLL.NewProc("yaml_delete")
	procYamlClear      = jadeDLL.NewProc("yaml_clear")
	procYamlDeleteFile = jadeDLL.NewProc("yaml_delete_file")
	procYamlKeys       = jadeDLL.NewProc("yaml_keys")
	procYamlLen        = jadeDLL.NewProc("yaml_len")

	// --- system ---
	procJadePrint           = jadeDLL.NewProc("jade_print")
	procJadePrintDialog     = jadeDLL.NewProc("jade_print_dialog")
	procGetPrinterList      = jadeDLL.NewProc("jade_get_printer_list")
	procSmartConvertEnc     = jadeDLL.NewProc("smart_convert_encoding")
	procSetProtocolSvcPath  = jadeDLL.NewProc("set_protocol_service_path")
	procGetPath             = jadeDLL.NewProc("getPath")
	procGetDisplaysInfo     = jadeDLL.NewProc("get_displays_info")
	procGetLocale           = jadeDLL.NewProc("getLocale")
	procClearDataDirectory  = jadeDLL.NewProc("clear_data_directory")
	procJadeviewVersion     = jadeDLL.NewProc("jadeview_version")
	procRegisterURLScheme   = jadeDLL.NewProc("register_url_scheme")
	procUnregisterURLScheme = jadeDLL.NewProc("unregister_url_scheme")
	procRegisterFileAssoc   = jadeDLL.NewProc("register_file_association")
	procUnregisterFileAssoc = jadeDLL.NewProc("unregister_file_association")
	procRegisterGlobalHK    = jadeDLL.NewProc("register_global_hotkey")
	procUnregisterGlobalHK  = jadeDLL.NewProc("unregister_global_hotkey")
	procSetLoginAutostart   = jadeDLL.NewProc("set_login_autostart")
	procGetLoginAutostart   = jadeDLL.NewProc("get_login_autostart")
	procGetFileIcon         = jadeDLL.NewProc("get_file_icon")
	procGetWebviewVersion   = jadeDLL.NewProc("get_webview_version")
	procIsWindows11         = jadeDLL.NewProc("is_windows_11")
	procJadeTextFree        = jadeDLL.NewProc("jade_text_free")
	procJadeTextCreate      = jadeDLL.NewProc("jade_text_create")
	procRegisterResource    = jadeDLL.NewProc("register_resource")
	procUnregisterResource  = jadeDLL.NewProc("unregister_resource")
	procClearWindowRes      = jadeDLL.NewProc("clear_window_resources")
	procClipboardReadText   = jadeDLL.NewProc("clipboard_read_text")
	procClipboardWriteText  = jadeDLL.NewProc("clipboard_write_text")
	procGetCursorPosition   = jadeDLL.NewProc("get_cursor_position")
	procJadeNtpNow          = jadeDLL.NewProc("jade_ntp_now")
)
