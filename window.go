//go:build linux && amd64 && jadeview_cgo && cgo

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

func (o *WindowOptions) toC(pool *cstrPool) *C.WebViewWindowOptions {
	if o == nil {
		return nil
	}
	c := new(C.WebViewWindowOptions)
	if err := Preload(); err != nil {
		panic(err)
	}
	if C.window_options_default(c, C.sizeof_WebViewWindowOptions) != 0 {
		panic("JadeView: WindowOptionsDefault failed")
	}
	c.title = pool.s(o.Title)
	c.width = C.int32_t(o.Width)
	c.height = C.int32_t(o.Height)
	c.resizable = b2i(o.Resizable)
	c.noframe = b2i(o.NoFrame)
	c.title_bar_style = pool.s(o.TitleBarStyle)
	c.transparent = b2i(o.Transparent)
	c.background_color = pool.s(o.BackgroundColor)
	c.always_on_top = b2i(o.AlwaysOnTop)
	c.theme = pool.s(o.Theme)
	c.maximized = b2i(o.Maximized)
	c.maximizable = b2i(o.Maximizable)
	c.minimizable = b2i(o.Minimizable)
	c.x = C.int32_t(o.X)
	c.y = C.int32_t(o.Y)
	c.min_width = C.int32_t(o.MinWidth)
	c.min_height = C.int32_t(o.MinHeight)
	c.max_width = C.int32_t(o.MaxWidth)
	c.max_height = C.int32_t(o.MaxHeight)
	c.fullscreen = b2i(o.Fullscreen)
	c.focus = b2i(o.Focus)
	c.hide_window = b2i(o.HideWindow)
	c.use_page_icon = b2i(o.UsePageIcon)
	c.skip_taskbar = b2i(o.SkipTaskbar)
	c.no_activate = b2i(o.NoActivate)
	c.background_material = pool.s(o.BackgroundMaterial)
	c.titlebar_back = b2i(!o.HideTitlebarBack)
	c.titlebar_refresh = b2i(!o.HideTitlebarRefresh)
	c.titlebar_icon = b2i(!o.HideTitlebarIcon)
	return c
}
func (s *WebViewSettings) toC(pool *cstrPool) *C.WebViewSettings {
	if s == nil {
		return nil
	}
	c := new(C.WebViewSettings)
	if err := Preload(); err != nil {
		panic(err)
	}
	if C.webview_settings_default(c, C.sizeof_WebViewSettings) != 0 {
		panic("JadeView: WebViewSettingsDefault failed")
	}
	c.autoplay = b2i(s.Autoplay)
	c.background_throttling = b2i(s.BackgroundThrottling)
	c.allow_right_click = b2i(s.AllowRightClick)
	c.user_agent = pool.s(s.UserAgent)
	c.preload_js = pool.s(s.PreloadJS)
	c.allow_fullscreen = b2i(s.AllowFullscreen)
	c.postmessage_allowed_origins = pool.s(s.PostMessageAllowedOrigins)
	c.cors_allowed_origins = pool.s(s.CORSAllowedOrigins)
	c.password_autofill = b2i(s.PasswordAutofill)
	c.form_autofill = b2i(s.FormAutofill)
	c.incognito = b2i(s.Incognito)
	c.disable_clipboard = b2i(s.DisableClipboard)
	c.proxy_url = pool.s(s.ProxyURL)
	c.focused = b2i(s.Focused)
	c.profile_name = pool.s(s.ProfileName)
	c.devtools = C.int32_t(s.DevTools)
	c.zoom_factor = C.double(s.ZoomFactor)
	c.background_color = pool.s(s.BackgroundColor)
	return c
}

// CreateWindow 创建 Window 和默认 WebView，返回 Window ID。nil 配置使用原生默认值。
func CreateWindow(url string, parentID uint32, opts *WindowOptions, settings *WebViewSettings) uint32 {
	pool := &cstrPool{}
	defer pool.free()
	if err := Preload(); err != nil {
		return 0
	}
	return uint32(C.create_webview_window(pool.s(url), C.uint32_t(parentID), opts.toC(pool), settings.toC(pool), opts.features(pool)))
}

// CreateEmptyWindow 创建不带 WebView 的原生 Window。
func CreateEmptyWindow(parentID uint32, opts *WindowOptions) uint32 {
	pool := &cstrPool{}
	defer pool.free()
	if err := Preload(); err != nil {
		return 0
	}
	return uint32(C.create_window(C.uint32_t(parentID), opts.toC(pool), opts.features(pool)))
}

func (o *WindowOptions) features(pool *cstrPool) *C.JadeWindowFeatures {
	if o == nil || (!o.ContentProtection && !o.AutoSaveState && o.Splash == nil) {
		return nil
	}
	f := &C.JadeWindowFeatures{struct_size: C.sizeof_JadeWindowFeatures, abi_version: C.uint32_t(ABIVersion), content_protection: b2i(o.ContentProtection), auto_save_state: b2i(o.AutoSaveState)}
	if s := o.Splash; s != nil {
		// cgo 不允许 C 结构内嵌 Go 指针。两层原生结构以及 LOGO 都放到 C 内存。
		raw := C.calloc(1, C.sizeof_JadeSplashOptions)
		if raw == nil {
			return nil
		}
		pool.ptrs = append(pool.ptrs, raw)
		dst := (*C.JadeSplashOptions)(raw)
		*dst = C.JadeSplashOptions{struct_size: C.sizeof_JadeSplashOptions, abi_version: C.uint32_t(ABIVersion),
			flags: C.uint32_t(b2i(s.Disabled)), logo_path: pool.s(s.LogoPath), background_color: C.uint32_t(s.BackgroundColor),
			accent_color: C.uint32_t(s.AccentColor), minimum_ms: C.uint32_t(s.MinimumMS)}
		if len(s.LogoData) > 0 {
			buf := C.CBytes(s.LogoData)
			pool.ptrs = append(pool.ptrs, buf)
			dst.logo_data = (*C.uint8_t)(buf)
			dst.logo_data_size = C.size_t(len(s.LogoData))
		}
		f.splash = dst
	}
	return f
}

// --- 常用窗口操作（薄封装） ---

func Navigate(webviewID uint32, url, headersJSON string) bool {
	pool := &cstrPool{}
	defer pool.free()
	return C.navigate_to_url(C.uint32_t(webviewID), pool.s(url), pool.s(headersJSON)) == 1
}

func Reload(webviewID uint32) bool {
	return C.reload_webview_window(C.uint32_t(webviewID)) == 1
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

// CreateBorderlessWindow 创建无边框窗口，返回 Window ID。
func CreateBorderlessWindow(url string, settings *WebViewSettings) uint32 {
	opts := DefaultWindowOptions()
	opts.NoFrame = true
	return CreateWindow(url, 0, &opts, settings)
}
