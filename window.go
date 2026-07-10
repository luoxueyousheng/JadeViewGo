//go:build linux

package jadeview

/*
#include <stdlib.h>
#include "JadeView.h"
*/
import "C"

import (
	"unsafe"
)

// cstrPool 管理一批临时 C 字符串，统一在调用结束后释放。
// 空字符串映射为 NULL 指针（库会使用对应默认值）。
type cstrPool struct {
	ptrs []unsafe.Pointer
}

func (p *cstrPool) s(str string) *C.char {
	if str == "" {
		return nil
	}
	return p.sAlways(str)
}

// sAlways 同 s，但空字符串也返回指向 "\0" 的有效指针。个别 API 区分 NULL 与空串
// （如 set_protocol_service_path：空串=内存 JAPK 模式，NULL 被拒绝）。
func (p *cstrPool) sAlways(str string) *C.char {
	c := C.CString(str)
	p.ptrs = append(p.ptrs, unsafe.Pointer(c))
	return c
}

func (p *cstrPool) free() {
	for _, x := range p.ptrs {
		C.free(x)
	}
	p.ptrs = nil
}

func b2i(b bool) C.int32_t {
	if b {
		return 1
	}
	return 0
}

func (o *WindowOptions) toC(pool *cstrPool) C.WebViewWindowOptions {
	return C.WebViewWindowOptions{
		title:              pool.s(o.Title),
		width:              C.int32_t(o.Width),
		height:             C.int32_t(o.Height),
		resizable:          b2i(o.Resizable),
		frame_style:        pool.s(o.FrameStyle),
		transparent:        b2i(o.Transparent),
		background_color:   pool.s(o.BackgroundColor),
		always_on_top:      b2i(o.AlwaysOnTop),
		theme:              pool.s(o.Theme),
		maximized:          b2i(o.Maximized),
		maximizable:        b2i(o.Maximizable),
		minimizable:        b2i(o.Minimizable),
		x:                  C.int32_t(o.X),
		y:                  C.int32_t(o.Y),
		min_width:          C.int32_t(o.MinWidth),
		min_height:         C.int32_t(o.MinHeight),
		max_width:          C.int32_t(o.MaxWidth),
		max_height:         C.int32_t(o.MaxHeight),
		fullscreen:         b2i(o.Fullscreen),
		focus:              b2i(o.Focus),
		hide_window:        b2i(o.HideWindow),
		use_page_icon:      b2i(o.UsePageIcon),
		content_protection: b2i(o.ContentProtection),
		auto_save_state:    b2i(o.AutoSaveState),
		skip_taskbar:       b2i(o.SkipTaskbar),
		no_activate:        b2i(o.NoActivate),
	}
}

func (s *WebViewSettings) toC(pool *cstrPool) C.WebViewSettings {
	return C.WebViewSettings{
		autoplay:                 b2i(s.Autoplay),
		background_throttling:    b2i(s.BackgroundThrottling),
		allow_right_click:        b2i(s.AllowRightClick),
		ua:                       pool.s(s.UserAgent),
		preload_js:               pool.s(s.PreloadJS),
		allow_fullscreen:         b2i(s.AllowFullscreen),
		postmessage_whitelist:    pool.s(s.PostMessageWhitelist),
		cors_whitelist:           pool.s(s.CORSWhitelist),
		autofill:                 b2i(s.Autofill),
		general_autofill_enabled: b2i(s.GeneralAutofillEnabled),
		incognito:                b2i(s.Incognito),
		disable_clipboard:        b2i(s.DisableClipboard),
		proxy_url:                pool.s(s.ProxyURL),
		focused:                  b2i(s.Focused),
	}
}

// CreateWindow 创建一个 WebView 窗口，返回 window_id（0 表示失败）。
//
//   - url      : 初始地址（http(s):// 或自定义协议）
//   - parentID : 父窗口 ID，0 表示无父窗口
//   - opts     : 窗口选项，nil 则用 DefaultWindowOptions()
//   - settings : WebView 高级设置，nil 则全用默认
func CreateWindow(url string, parentID uint32, opts *WindowOptions, settings *WebViewSettings) uint32 {
	pool := &cstrPool{}
	defer pool.free()

	if opts == nil {
		d := DefaultWindowOptions()
		opts = &d
	}
	if settings == nil {
		settings = &WebViewSettings{}
	}
	copts := opts.toC(pool)
	csettings := settings.toC(pool)

	id := C.create_webview_window(pool.s(url), C.uint32_t(parentID), &copts, &csettings)
	return uint32(id)
}

// CreateBorderlessWindow 创建独立无边框 WebView 窗口，返回 window_id。
// 仅此类窗口可通过 GetWindowHWND 取原生句柄。
func CreateBorderlessWindow(url string, settings *WebViewSettings) uint32 {
	pool := &cstrPool{}
	defer pool.free()
	if settings == nil {
		settings = &WebViewSettings{}
	}
	csettings := settings.toC(pool)
	id := C.create_borderless_webview_window(pool.s(url), &csettings)
	return uint32(id)
}

// --- 常用窗口操作（薄封装） ---

func Navigate(windowID uint32, url, headersJSON string) bool {
	pool := &cstrPool{}
	defer pool.free()
	return C.navigate_to_url(C.uint32_t(windowID), pool.s(url), pool.s(headersJSON)) == 1
}

func Reload(windowID uint32) bool {
	return C.reload_webview_window(C.uint32_t(windowID)) == 1
}

// ExecuteJavaScript 执行 JS，返回一个唯一 id；结果通过 "javascript-result" 事件异步返回。
func ExecuteJavaScript(windowID uint32, script string) int32 {
	pool := &cstrPool{}
	defer pool.free()
	return int32(C.execute_javascript(C.uint32_t(windowID), pool.s(script)))
}

func SetTitle(windowID uint32, title string) bool {
	pool := &cstrPool{}
	defer pool.free()
	return C.set_window_title(C.uint32_t(windowID), pool.s(title)) == 1
}

func SetSize(windowID uint32, width, height int) bool {
	return C.set_window_size(C.uint32_t(windowID), C.int32_t(width), C.int32_t(height)) == 1
}

func SetPosition(windowID uint32, x, y int) bool {
	return C.set_window_position(C.uint32_t(windowID), C.int32_t(x), C.int32_t(y)) == 1
}

func SetVisible(windowID uint32, visible bool) bool {
	return C.set_window_visible(C.uint32_t(windowID), b2i(visible)) == 1
}

func SetFocus(windowID uint32) bool {
	return C.set_window_focus(C.uint32_t(windowID)) == 1
}

func SetAlwaysOnTop(windowID uint32, on bool) bool {
	return C.set_window_always_on_top(C.uint32_t(windowID), b2i(on)) == 1
}

func Close(windowID uint32) bool {
	return C.close_window(C.uint32_t(windowID)) == 1
}

func Minimize(windowID uint32) bool {
	return C.minimize_window(C.uint32_t(windowID)) == 1
}

func ToggleMaximize(windowID uint32) bool {
	return C.toggle_maximize_window(C.uint32_t(windowID)) == 1
}

func IsMaximized(windowID uint32) bool {
	return C.is_window_maximized(C.uint32_t(windowID)) == 1
}

// WindowCount 返回当前窗口数量。
func WindowCount() uint32 {
	return uint32(C.get_window_count())
}
