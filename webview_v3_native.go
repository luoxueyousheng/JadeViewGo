//go:build windows || (linux && amd64 && !jadeview_cgo)

package jadeview

import (
	"runtime"
	"unsafe"
)

type cWebViewBounds struct{ x, y, width, height *byte }

func (b *WebViewBounds) toC(pool *cstrs) *cWebViewBounds {
	if b == nil {
		return nil
	}
	return &cWebViewBounds{pool.p(b.X), pool.p(b.Y), pool.p(b.Width), pool.p(b.Height)}
}

type cRequestResult struct {
	state     uint32
	errorCode int32
	objectID  uint32
	message   *byte
}

type cWebViewIDArray struct {
	ids   *uint32
	count uintptr
}

type cWebViewLayout struct {
	bounds   cWebViewBounds
	logical  LogicalRect
	physical PhysicalRect
}

// DefaultWebView 查找 Window 的默认 WebView；没有默认视图时返回 0。
func DefaultWebView(windowID uint32) uint32 {
	lo, _, _ := procDefaultWebView.Call(uintptr(windowID))
	return uint32(lo)
}

// SetDefaultWebView 选择 Window 的默认 WebView。
func SetDefaultWebView(windowID, webviewID uint32) bool {
	lo, _, _ := procSetDefaultWebView.Call(uintptr(windowID), uintptr(webviewID))
	return i32(lo) == 0
}

// WebViewParentWindow 返回 WebView 所属 Window 的 ID。
func WebViewParentWindow(webviewID uint32) uint32 {
	lo, _, _ := procWebViewParentWindow.Call(uintptr(webviewID))
	return uint32(lo)
}

// CreateWebView 创建尚未挂载到 Window 的逻辑 WebView。
func CreateWebView(url string, settings *WebViewSettings) uint32 {
	pool := &cstrs{}
	lo, _, _ := procCreateWebView.Call(uintptr(unsafe.Pointer(pool.p(url))), uintptr(unsafe.Pointer(settings.toC(pool))))
	runtime.KeepAlive(pool)
	return uint32(lo)
}

// AddWebView 创建并挂载一个 WebView，返回 WebView ID。
func AddWebView(windowID uint32, url string, settings *WebViewSettings, bounds *WebViewBounds) uint32 {
	pool := &cstrs{}
	lo, _, _ := procAddWebView.Call(uintptr(windowID), uintptr(unsafe.Pointer(pool.p(url))), uintptr(unsafe.Pointer(settings.toC(pool))), uintptr(unsafe.Pointer(bounds.toC(pool))))
	runtime.KeepAlive(pool)
	return uint32(lo)
}

// AttachWebView 把已存在的逻辑 WebView 挂载到 Window。
func AttachWebView(windowID, webviewID uint32, bounds *WebViewBounds) RequestID {
	pool := &cstrs{}
	lo, hi, _ := procAttachWebView.Call(uintptr(windowID), uintptr(webviewID), uintptr(unsafe.Pointer(bounds.toC(pool))))
	runtime.KeepAlive(pool)
	return requestID(lo, hi)
}

// SetWebViewBounds 更新 WebView 在 Window 内容区内的布局。
func SetWebViewBounds(webviewID uint32, bounds *WebViewBounds) RequestID {
	pool := &cstrs{}
	lo, hi, _ := procSetWebViewBounds.Call(uintptr(webviewID), uintptr(unsafe.Pointer(bounds.toC(pool))))
	runtime.KeepAlive(pool)
	return requestID(lo, hi)
}

// NavigateWebView 导航并返回可查询的请求 ID。
func NavigateWebView(webviewID uint32, url string) RequestID {
	pool := &cstrs{}
	lo, hi, _ := procNavigateWebView.Call(uintptr(webviewID), uintptr(unsafe.Pointer(pool.p(url))))
	runtime.KeepAlive(pool)
	return requestID(lo, hi)
}

// ReloadWebView 刷新并返回可查询的请求 ID。
func ReloadWebView(webviewID uint32) RequestID {
	lo, hi, _ := procReloadWebView.Call(uintptr(webviewID))
	return requestID(lo, hi)
}

// ExecuteJavaScript 提交脚本；用 GetRequestResult 查询结果，完成后调用 ReleaseRequest。
func ExecuteJavaScript(webviewID uint32, script string) RequestID {
	pool := &cstrs{}
	lo, hi, _ := procExecuteJavaScript.Call(uintptr(webviewID), uintptr(unsafe.Pointer(pool.p(script))))
	runtime.KeepAlive(pool)
	return requestID(lo, hi)
}

// SetWebViewVisible 独立设置 WebView 的可见性。
func SetWebViewVisible(webviewID uint32, visible bool) RequestID {
	lo, hi, _ := procSetWebViewVisible.Call(uintptr(webviewID), b2u(visible))
	return requestID(lo, hi)
}

// FocusWebView 聚焦 WebView。
func FocusWebView(webviewID uint32) RequestID {
	lo, hi, _ := procFocusWebView.Call(uintptr(webviewID))
	return requestID(lo, hi)
}

// DestroyWebView 销毁 WebView，而不关闭其 Window。
func DestroyWebView(webviewID uint32) RequestID {
	lo, hi, _ := procDestroyWebView.Call(uintptr(webviewID))
	return requestID(lo, hi)
}

// RecreateWebView 重建 WebView 并返回请求 ID。
func RecreateWebView(webviewID uint32, url string, settings *WebViewSettings) RequestID {
	pool := &cstrs{}
	lo, hi, _ := procRecreateWebView.Call(uintptr(webviewID), uintptr(unsafe.Pointer(pool.p(url))), uintptr(unsafe.Pointer(settings.toC(pool))))
	runtime.KeepAlive(pool)
	return requestID(lo, hi)
}

var (
	procDefaultWebView      = np("get_default_webview")
	procSetDefaultWebView   = np("set_default_webview")
	procWebViewParentWindow = np("get_webview_parent_window_id")
	procCreateWebView       = np("create_webview")
	procAddWebView          = np("add_webview")
	procAttachWebView       = np("add_webview_to_window")
	procSetWebViewBounds    = np("set_webview_bounds")
	procNavigateWebView     = np("navigate_webview")
	procReloadWebView       = np("reload_webview")
	procExecuteJavaScript   = np("execute_webview_script")
	procSetWebViewVisible   = np("set_webview_visible")
	procFocusWebView        = np("focus_webview")
	procDestroyWebView      = np("destroy_webview")
	procRecreateWebView     = np("recreate_webview")
)

var (
	procGetRequestResult   = np("get_request_result")
	procFreeRequestResult  = np("free_request_result")
	procReleaseRequest     = np("release_request")
	procGetWindowState     = np("get_window_state")
	procGetWebViewState    = np("get_webview_state")
	procGetWindowWebViews  = np("get_window_webviews")
	procFreeWebViewIDArray = np("free_webview_id_array")
	procGetWebViewBounds   = np("get_webview_bounds")
	procFreeWebViewBounds  = np("free_webview_bounds_result")
)

// GetWebViewBounds 返回已应用的布局快照，原生状态码 0 表示成功。
func GetWebViewBounds(webviewID uint32) (WebViewLayout, int32) {
	var out cWebViewLayout
	r, _, _ := procGetWebViewBounds.Call(uintptr(webviewID), uintptr(unsafe.Pointer(&out)))
	if i32(r) != 0 {
		return WebViewLayout{}, i32(r)
	}
	defer procFreeWebViewBounds.Call(uintptr(unsafe.Pointer(&out)))
	s := func(p *byte) string { return goString(uintptr(unsafe.Pointer(p))) }
	b := out.bounds
	return WebViewLayout{WebViewBounds{s(b.x), s(b.y), s(b.width), s(b.height)}, out.logical, out.physical}, 0
}

// GetRequestResult 返回请求快照及原生状态码（0 表示查询成功）。
func GetRequestResult(id RequestID) (RequestResult, int32) {
	var out cRequestResult
	r := callRequestOut(procGetRequestResult, id, &out)
	if r != 0 {
		return RequestResult{}, r
	}
	defer procFreeRequestResult.Call(uintptr(unsafe.Pointer(&out)))
	return RequestResult{RequestState(out.state), out.errorCode, out.objectID, goString(uintptr(unsafe.Pointer(out.message)))}, 0
}

func ReleaseRequest(id RequestID) int32 {
	if unsafe.Sizeof(uintptr(0)) == 4 {
		r, _, _ := procReleaseRequest.Call(uintptr(uint32(id)), uintptr(id>>32))
		return i32(r)
	}
	r, _, _ := procReleaseRequest.Call(uintptr(id))
	return i32(r)
}

func callRequestOut(p *jvProc, id RequestID, out *cRequestResult) int32 {
	if unsafe.Sizeof(uintptr(0)) == 4 {
		r, _, _ := p.Call(uintptr(uint32(id)), uintptr(id>>32), uintptr(unsafe.Pointer(out)))
		return i32(r)
	}
	r, _, _ := p.Call(uintptr(id), uintptr(unsafe.Pointer(out)))
	return i32(r)
}

func requestID(lo, hi uintptr) RequestID {
	if unsafe.Sizeof(uintptr(0)) == 4 {
		return RequestID(uint64(lo) | uint64(hi)<<32)
	}
	return RequestID(lo)
}

func GetWindowState(windowID uint32) (ObjectState, int32) {
	return getObjectState(procGetWindowState, windowID)
}

func GetWebViewState(webviewID uint32) (ObjectState, int32) {
	return getObjectState(procGetWebViewState, webviewID)
}

func getObjectState(p *jvProc, id uint32) (ObjectState, int32) {
	var state ObjectState
	r, _, _ := p.Call(uintptr(id), uintptr(unsafe.Pointer(&state)))
	return state, i32(r)
}

func GetWindowWebViews(windowID uint32) ([]uint32, int32) {
	var a cWebViewIDArray
	r, _, _ := procGetWindowWebViews.Call(uintptr(windowID), uintptr(unsafe.Pointer(&a)))
	if i32(r) != 0 {
		return nil, i32(r)
	}
	defer procFreeWebViewIDArray.Call(uintptr(unsafe.Pointer(&a)))
	return append([]uint32(nil), unsafe.Slice(a.ids, a.count)...), 0
}
