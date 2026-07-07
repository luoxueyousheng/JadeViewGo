package jadeview

/*
#include <stdlib.h>
#include "JadeView.h"
*/
import "C"

// --- 窗口状态查询 ---

func IsMinimized(windowID uint32) bool {
	return C.is_window_minimized(C.uint32_t(windowID)) == 1
}

func IsVisible(windowID uint32) bool {
	return C.is_window_visible(C.uint32_t(windowID)) == 1
}

func IsFocused(windowID uint32) bool {
	return C.is_window_focused(C.uint32_t(windowID)) == 1
}

func IsFullscreen(windowID uint32) bool {
	return C.is_window_fullscreen(C.uint32_t(windowID)) == 1
}

// GetWindowBounds 返回窗口边界 JSON。
func GetWindowBounds(windowID uint32) (string, bool) {
	return bufCallInt(256, func(buf *C.char, n C.int32_t) C.int32_t {
		return C.get_window_bounds(C.uint32_t(windowID), buf, n)
	})
}

// GetWebViewURL 返回当前 WebView 地址。
func GetWebViewURL(windowID uint32) (string, bool) {
	return bufCallInt(2048, func(buf *C.char, n C.int32_t) C.int32_t {
		return C.get_webview_url(C.uint32_t(windowID), buf, n)
	})
}

// GetWindowHWND 返回窗口原生句柄（beta.9 起支持所有方法创建的窗口）。
func GetWindowHWND(windowID uint32) uintptr {
	return uintptr(C.get_window_hwnd(C.uint32_t(windowID)))
}

// GetWindowID 根据 HWND 反查窗口 ID（beta.9 新增），0 表示未找到。
func GetWindowID(hwnd uintptr) uint32 {
	return uint32(C.get_window_id(C.int32_t(hwnd)))
}

// --- 尺寸 / 可调整 ---

func SetMinSize(windowID uint32, width, height int) bool {
	return C.set_window_min_size(C.uint32_t(windowID), C.int32_t(width), C.int32_t(height)) == 1
}

func SetMaxSize(windowID uint32, width, height int) bool {
	return C.set_window_max_size(C.uint32_t(windowID), C.int32_t(width), C.int32_t(height)) == 1
}

func SetResizable(windowID uint32, resizable bool) bool {
	return C.set_window_resizable(C.uint32_t(windowID), b2i(resizable)) == 1
}

func SetFullscreen(windowID uint32, fullscreen bool) bool {
	return C.set_window_fullscreen(C.uint32_t(windowID), b2i(fullscreen)) == 1
}

// SetIgnoreCursorEvents 设置窗口是否穿透鼠标事件。
func SetIgnoreCursorEvents(windowID uint32, ignore bool) bool {
	return C.set_window_ignore_cursor_events(C.uint32_t(windowID), b2i(ignore)) == 1
}

// --- 外观 / 层级 ---

func SetSkipTaskbar(windowID uint32, skip bool) bool {
	return C.set_window_skip_taskbar(C.uint32_t(windowID), b2i(skip)) == 1
}

func SetNoActivate(windowID uint32, noActivate bool) bool {
	return C.set_window_no_activate(C.uint32_t(windowID), b2i(noActivate)) == 1
}

// SetLevel 设置窗口层级：topmost | normal | bottom | desktop。
func SetLevel(windowID uint32, level string) bool {
	pool := &cstrPool{}
	defer pool.free()
	return C.set_window_level(C.uint32_t(windowID), pool.s(level)) == 1
}

func SetContentProtection(windowID uint32, on bool) bool {
	return C.set_content_protection(C.uint32_t(windowID), b2i(on)) == 1
}

func SetZoom(windowID uint32, level float64) bool {
	return C.set_webview_zoom(C.uint32_t(windowID), C.double(level)) == 1
}

func SetFrameStyle(windowID uint32, frameStyle string) bool {
	pool := &cstrPool{}
	defer pool.free()
	return C.set_window_frame_style(C.uint32_t(windowID), pool.s(frameStyle)) == 1
}

func SetTheme(windowID uint32, theme string) bool {
	pool := &cstrPool{}
	defer pool.free()
	return C.set_window_theme(C.uint32_t(windowID), pool.s(theme)) == 1
}

// GetTheme 返回窗口主题（int32，含义见库文档）。
func GetTheme(windowID uint32) int32 {
	return int32(C.get_window_theme(C.uint32_t(windowID)))
}

func SetBackdrop(windowID uint32, backdropType string) bool {
	pool := &cstrPool{}
	defer pool.free()
	return C.set_window_backdrop(C.uint32_t(windowID), pool.s(backdropType)) == 1
}

// SetBackgroundColor 设置窗口纯色底：#RRGGBBAA。
func SetBackgroundColor(windowID uint32, colorHex string) bool {
	pool := &cstrPool{}
	defer pool.free()
	return C.set_window_background_color(C.uint32_t(windowID), pool.s(colorHex)) == 1
}

func SetEnabled(windowID uint32, enabled bool) bool {
	return C.set_window_enabled(C.uint32_t(windowID), b2i(enabled)) == 1
}

func RequestRedraw(windowID uint32) bool {
	return C.request_redraw(C.uint32_t(windowID)) == 1
}

// SetTitlebarOverlayStyle 设置标题栏覆盖层样式（Windows only）。height<=0 不修改高度。
func SetTitlebarOverlayStyle(windowID uint32, height int, iconColorHex, hoverBgHex string) bool {
	pool := &cstrPool{}
	defer pool.free()
	return C.set_titlebar_overlay_style(C.uint32_t(windowID), C.int32_t(height), pool.s(iconColorHex), pool.s(hoverBgHex)) == 1
}

// --- IPC / DevTools / 其它 ---

// SendIPCMessage 向前端发送 IPC 消息。
func SendIPCMessage(windowID uint32, messageType, messageContent string) bool {
	pool := &cstrPool{}
	defer pool.free()
	return C.send_ipc_message(C.uint32_t(windowID), pool.s(messageType), pool.s(messageContent)) == 1
}

func OpenDevtools(windowID uint32) bool  { return C.open_devtools(C.uint32_t(windowID)) == 1 }
func CloseDevtools(windowID uint32) bool { return C.close_devtools(C.uint32_t(windowID)) == 1 }
func IsDevtoolsOpen(windowID uint32) bool {
	return C.is_devtools_open(C.uint32_t(windowID)) == 1
}

func ClearBrowsingData(windowID uint32) bool {
	return C.clear_browsing_data(C.uint32_t(windowID)) == 1
}

// SetWindowProgress 设置任务栏进度。state 含义见库文档。
func SetWindowProgress(windowID uint32, progress, state int) bool {
	return C.set_window_progress(C.uint32_t(windowID), C.int32_t(progress), C.int32_t(state)) == 1
}

// FlashWindow 闪烁窗口 count 次。
func FlashWindow(windowID uint32, count uint32) bool {
	return C.flash_window(C.uint32_t(windowID), C.uint32_t(count)) == 1
}

func ShowAboutDialog(windowID uint32) bool {
	return C.show_about_dialog(C.uint32_t(windowID)) == 1
}
