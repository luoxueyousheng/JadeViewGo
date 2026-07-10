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
	"sync"
	"syscall"
	"unsafe"
)

// jadeDLL 惰性加载 JadeView.dll：首次调用任一 API（或 Preload）时才释放文件并
// LoadLibrary，仅 import 本包不产生任何磁盘副作用。
// 路径解析优先级：exe 同目录的 JadeView.dll（便于手动覆盖/调试）→
// 释放内置副本到 %TEMP%\jadeview\<arch>-<内容哈希前8位>\（内容寻址，
// 不同版本/架构各占一目录；已存在文件按完整 sha256 校验，不符则重写）。
var (
	dllOnce sync.Once
	jadeDLL *syscall.LazyDLL
)

func loadedDLL() *syscall.LazyDLL {
	dllOnce.Do(func() { jadeDLL = syscall.NewLazyDLL(resolveDLLPath()) })
	return jadeDLL
}

// Preload 提前释放并加载内置 JadeView.dll，返回加载错误（nil=成功）。
// 可选调用：不调用时首次 API 调用会自动加载，但那时加载失败会 panic
// （syscall.LazyProc 语义）；需要优雅处理加载失败（临时目录禁止执行、
// 杀软拦截等）的宿主应在启动早期调用本函数并检查错误。
func Preload() error { return loadedDLL().Load() }

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
	// 目录名已含内容哈希，但仍校验磁盘文件的完整 sha256：同大小篡改/写坏残留不会被加载
	if b, err := os.ReadFile(dst); err != nil || sha256.Sum256(b) != sum {
		_ = os.MkdirAll(dir, 0o755)
		// 写失败不致命（可能另一进程刚写完并占用）；加载失败会在 Preload/首次调用时报错
		_ = os.WriteFile(dst, embeddedJadeViewDLL, 0o644)
	}
	return dst
}

// jvProc 惰性解析的导出函数句柄：首次 Call/Find/Addr 时才触发 DLL 释放与加载
// （见 loadedDLL），方法集与 syscall.LazyProc 的使用面一致。
type jvProc struct {
	name string
	once sync.Once
	proc *syscall.LazyProc
}

// np 声明一个导出函数（仅记录名字，不触发加载）。
func np(name string) *jvProc { return &jvProc{name: name} }

func (p *jvProc) resolve() *syscall.LazyProc {
	p.once.Do(func() { p.proc = loadedDLL().NewProc(p.name) })
	return p.proc
}

func (p *jvProc) Call(a ...uintptr) (uintptr, uintptr, error) { return p.resolve().Call(a...) }
func (p *jvProc) Find() error                                 { return p.resolve().Find() }
func (p *jvProc) Addr() uintptr                               { return p.resolve().Addr() }

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
// 缓冲区不足时（库返回 0）按 16×/256× 放大重试两次（覆盖超长 URL 等场景）。
// 语义与 Linux 版 helpers 一致。
func bufCallSize(initial int, fn func(buf unsafe.Pointer, size uintptr) int32) (string, bool) {
	for _, size := range []int{initial, initial * 16, initial * 256} {
		buf := make([]byte, size)
		if fn(unsafe.Pointer(&buf[0]), uintptr(size)) == 1 {
			return cBufToString(buf), true
		}
	}
	return "", false
}

// bufCallInt 同 bufCallSize，但长度参数为 int32。
func bufCallInt(initial int, fn func(buf unsafe.Pointer, size int32) int32) (string, bool) {
	for _, size := range []int{initial, initial * 16, initial * 256} {
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
	procJadeViewInit   = np("JadeView_init")
	procRunMessageLoop = np("run_message_loop")
	procJadeviewExit   = np("jadeview_exit")

	// --- japk ---
	procSetPublicKey    = np("JadeView_set_public_key")
	procLoadFromBytes   = np("JadeView_load_from_bytes")
	procIsLoaded        = np("JadeView_is_loaded")
	procGetAppSignature = np("JadeView_get_app_signature")
	procGetSignatureInf = np("JadeView_get_signature_info")
	procUnload          = np("JadeView_unload")

	// --- window ---
	procCreateWebviewWindow    = np("create_webview_window")
	procCreateBorderlessWindow = np("create_borderless_webview_window")
	procGetWindowHwnd          = np("get_window_hwnd")
	procGetWindowID            = np("get_window_id")
	procNavigateToURL          = np("navigate_to_url")
	procReloadWebviewWindow    = np("reload_webview_window")
	procExecuteJavascript      = np("execute_javascript")
	procSetWindowTitle         = np("set_window_title")
	procSetWindowSize          = np("set_window_size")
	procSetWindowPosition      = np("set_window_position")
	procSetWindowVisible       = np("set_window_visible")
	procSetWindowFocus         = np("set_window_focus")
	procSetWindowAlwaysOnTop   = np("set_window_always_on_top")
	procSetWindowSkipTaskbar   = np("set_window_skip_taskbar")
	procSetWindowNoActivate    = np("set_window_no_activate")
	procSetWindowLevel         = np("set_window_level")
	procCloseWindow            = np("close_window")
	procMinimizeWindow         = np("minimize_window")
	procToggleMaximizeWindow   = np("toggle_maximize_window")
	procIsWindowMaximized      = np("is_window_maximized")
	procSetContentProtection   = np("set_content_protection")
	procSetWebviewZoom         = np("set_webview_zoom")
	procSetWindowFrameStyle    = np("set_window_frame_style")
	procGetWindowCount         = np("get_window_count")
	procJadeOn                 = np("jade_on")
	procJadeOff                = np("jade_off")
	procRegisterIPCHandler     = np("register_ipc_handler")
	procSetWindowTheme         = np("set_window_theme")
	procSetTitlebarOverlay     = np("set_titlebar_overlay_style")
	procSetWindowEnabled       = np("set_window_enabled")
	procGetWindowTheme         = np("get_window_theme")
	procRequestRedraw          = np("request_redraw")
	procSendIPCMessage         = np("send_ipc_message")
	procSetWindowBackdrop      = np("set_window_backdrop")
	procSetWindowBgColor       = np("set_window_background_color")
	procSetWindowFullscreen    = np("set_window_fullscreen")
	procGetWindowBounds        = np("get_window_bounds")
	procIsWindowMinimized      = np("is_window_minimized")
	procIsWindowVisible        = np("is_window_visible")
	procIsWindowFocused        = np("is_window_focused")
	procIsWindowFullscreen     = np("is_window_fullscreen")
	procSetWindowMinSize       = np("set_window_min_size")
	procSetWindowMaxSize       = np("set_window_max_size")
	procSetWindowResizable     = np("set_window_resizable")
	procSetIgnoreCursorEvents  = np("set_window_ignore_cursor_events")
	procGetWebviewURL          = np("get_webview_url")
	procOpenDevtools           = np("open_devtools")
	procCloseDevtools          = np("close_devtools")
	procIsDevtoolsOpen         = np("is_devtools_open")
	procClearBrowsingData      = np("clear_browsing_data")
	procSetWindowProgress      = np("set_window_progress")
	procFlashWindow            = np("flash_window")
	procShowAboutDialog        = np("show_about_dialog")

	// --- tray ---
	procTrayCreate          = np("tray_create")
	procTrayDestroy         = np("tray_destroy")
	procTraySetVisible      = np("tray_set_visible")
	procTraySetTooltip      = np("tray_set_tooltip")
	procTraySetIconFromFile = np("tray_set_icon_from_file")
	procTraySetMenuItems    = np("tray_set_menu_items")
	procTraySetIconFromData = np("set_tray_icon_from_data")

	// --- dialog / menu ---
	procShowNotification    = np("show_notification")
	procShowOpenDialog      = np("jade_dialog_show_open_dialog")
	procShowSaveDialog      = np("jade_dialog_show_save_dialog")
	procShowMessageBox      = np("jade_dialog_show_message_box")
	procShowErrorBox        = np("jade_dialog_show_error_box")
	procShowOpenDialogAsync = np("jade_dialog_show_open_dialog_async")
	procShowSaveDialogAsync = np("jade_dialog_show_save_dialog_async")
	procShowMessageBoxAsync = np("jade_dialog_show_message_box_async")
	procMenuItemCreate      = np("jade_menu_item_create")
	procMenuItemSetEnabled  = np("jade_menu_item_set_enabled")
	procMenuItemSetChecked  = np("jade_menu_item_set_checked")
	procSetContextMenuItems = np("jade_set_context_menu_items")
	procMenuItemDestroy     = np("jade_menu_item_destroy")

	// --- yaml ---
	procYamlSet        = np("yaml_set")
	procYamlGet        = np("yaml_get")
	procYamlSetStr     = np("yaml_set_str")
	procYamlGetAll     = np("yaml_get_all")
	procYamlHas        = np("yaml_has")
	procYamlDelete     = np("yaml_delete")
	procYamlClear      = np("yaml_clear")
	procYamlDeleteFile = np("yaml_delete_file")
	procYamlKeys       = np("yaml_keys")
	procYamlLen        = np("yaml_len")

	// --- system ---
	procJadePrint           = np("jade_print")
	procJadePrintDialog     = np("jade_print_dialog")
	procGetPrinterList      = np("jade_get_printer_list")
	procSmartConvertEnc     = np("smart_convert_encoding")
	procSetProtocolSvcPath  = np("set_protocol_service_path")
	procGetPath             = np("getPath")
	procGetDisplaysInfo     = np("get_displays_info")
	procGetLocale           = np("getLocale")
	procClearDataDirectory  = np("clear_data_directory")
	procJadeviewVersion     = np("jadeview_version")
	procRegisterURLScheme   = np("register_url_scheme")
	procUnregisterURLScheme = np("unregister_url_scheme")
	procRegisterFileAssoc   = np("register_file_association")
	procUnregisterFileAssoc = np("unregister_file_association")
	procRegisterGlobalHK    = np("register_global_hotkey")
	procUnregisterGlobalHK  = np("unregister_global_hotkey")
	procSetLoginAutostart   = np("set_login_autostart")
	procGetLoginAutostart   = np("get_login_autostart")
	procGetFileIcon         = np("get_file_icon")
	procGetWebviewVersion   = np("get_webview_version")
	procIsWindows11         = np("is_windows_11")
	procJadeTextFree        = np("jade_text_free")
	procJadeTextCreate      = np("jade_text_create")
	procRegisterResource    = np("register_resource")
	procUnregisterResource  = np("unregister_resource")
	procClearWindowRes      = np("clear_window_resources")
	procClipboardReadText   = np("clipboard_read_text")
	procClipboardWriteText  = np("clipboard_write_text")
	procGetCursorPosition   = np("get_cursor_position")
	procJadeNtpNow          = np("jade_ntp_now")
)
