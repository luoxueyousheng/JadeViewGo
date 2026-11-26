//go:build linux && amd64 && jadeview_cgo && cgo

package jadeview

/*
#include "JadeView.h"
*/
import "C"

import "unsafe"

func (b *WebViewBounds) toC(pool *cstrPool) *C.WebViewBounds {
	if b == nil {
		return nil
	}
	return &C.WebViewBounds{x: pool.s(b.X), y: pool.s(b.Y), width: pool.s(b.Width), height: pool.s(b.Height)}
}

// GetWebViewBounds 返回已应用的布局快照，原生状态码 0 表示成功。
func GetWebViewBounds(webviewID uint32) (WebViewLayout, int32) {
	var out C.WebViewBoundsResult
	r := int32(C.get_webview_bounds(C.uint32_t(webviewID), &out))
	if r != 0 {
		return WebViewLayout{}, r
	}
	defer C.free_webview_bounds_result(&out)
	b, l, p := out.bounds, out.logical, out.physical
	return WebViewLayout{
		WebViewBounds{C.GoString(b.x), C.GoString(b.y), C.GoString(b.width), C.GoString(b.height)},
		LogicalRect{float64(l.x), float64(l.y), float64(l.width), float64(l.height)},
		PhysicalRect{int32(p.x), int32(p.y), uint32(p.width), uint32(p.height)},
	}, 0
}

// DefaultWebView 查找 Window 的默认 WebView；没有默认视图时返回 0。
func DefaultWebView(windowID uint32) uint32 {
	return uint32(C.get_default_webview(C.uint32_t(windowID)))
}

// SetDefaultWebView 选择 Window 的默认 WebView。
func SetDefaultWebView(windowID, webviewID uint32) bool {
	return C.set_default_webview(C.uint32_t(windowID), C.uint32_t(webviewID)) == 0
}

// WebViewParentWindow 返回 WebView 所属 Window 的 ID。
func WebViewParentWindow(webviewID uint32) uint32 {
	return uint32(C.get_webview_parent_window_id(C.uint32_t(webviewID)))
}

// CreateWebView 创建尚未挂载到 Window 的逻辑 WebView。
func CreateWebView(url string, settings *WebViewSettings) uint32 {
	pool := &cstrPool{}
	defer pool.free()
	return uint32(C.create_webview(pool.s(url), settings.toC(pool)))
}

// AddWebView 创建并挂载一个 WebView，返回 WebView ID。
func AddWebView(windowID uint32, url string, settings *WebViewSettings, bounds *WebViewBounds) uint32 {
	pool := &cstrPool{}
	defer pool.free()
	return uint32(C.add_webview(C.uint32_t(windowID), pool.s(url), settings.toC(pool), bounds.toC(pool)))
}

// AttachWebView 把已存在的逻辑 WebView 挂载到 Window。
func AttachWebView(windowID, webviewID uint32, bounds *WebViewBounds) RequestID {
	pool := &cstrPool{}
	defer pool.free()
	return RequestID(C.add_webview_to_window(C.uint32_t(windowID), C.uint32_t(webviewID), bounds.toC(pool)))
}

// SetWebViewBounds 更新 WebView 在 Window 内容区内的布局。
func SetWebViewBounds(webviewID uint32, bounds *WebViewBounds) RequestID {
	pool := &cstrPool{}
	defer pool.free()
	return RequestID(C.set_webview_bounds(C.uint32_t(webviewID), bounds.toC(pool)))
}

// NavigateWebView 导航并返回可查询的请求 ID。
func NavigateWebView(webviewID uint32, url string) RequestID {
	pool := &cstrPool{}
	defer pool.free()
	return RequestID(C.navigate_webview(C.uint32_t(webviewID), pool.s(url)))
}

// ReloadWebView 刷新并返回可查询的请求 ID。
func ReloadWebView(webviewID uint32) RequestID {
	return RequestID(C.reload_webview(C.uint32_t(webviewID)))
}

// ExecuteJavaScript 提交脚本；用 GetRequestResult 查询结果，完成后调用 ReleaseRequest。
func ExecuteJavaScript(webviewID uint32, script string) RequestID {
	pool := &cstrPool{}
	defer pool.free()
	return RequestID(C.execute_webview_script(C.uint32_t(webviewID), pool.s(script)))
}

// SetWebViewVisible 独立设置 WebView 的可见性。
func SetWebViewVisible(webviewID uint32, visible bool) RequestID {
	return RequestID(C.set_webview_visible(C.uint32_t(webviewID), b2i(visible)))
}

// FocusWebView 聚焦 WebView。
func FocusWebView(webviewID uint32) RequestID {
	return RequestID(C.focus_webview(C.uint32_t(webviewID)))
}

// DestroyWebView 销毁 WebView，而不关闭其 Window。
func DestroyWebView(webviewID uint32) RequestID {
	return RequestID(C.destroy_webview(C.uint32_t(webviewID)))
}

// RecreateWebView 重建 WebView 并返回请求 ID。
func RecreateWebView(webviewID uint32, url string, settings *WebViewSettings) RequestID {
	pool := &cstrPool{}
	defer pool.free()
	return RequestID(C.recreate_webview(C.uint32_t(webviewID), pool.s(url), settings.toC(pool)))
}

// GetRequestResult 返回请求快照及原生状态码（0 表示查询成功）。
func GetRequestResult(id RequestID) (RequestResult, int32) {
	var out C.JadeRequestResult
	r := int32(C.get_request_result(C.uint64_t(id), &out))
	if r != 0 {
		return RequestResult{}, r
	}
	defer C.free_request_result(&out)
	return RequestResult{RequestState(out.state), int32(out.error_code), uint32(out.object_id), C.GoString(out.message)}, 0
}

func ReleaseRequest(id RequestID) int32 { return int32(C.release_request(C.uint64_t(id))) }

func GetWindowState(windowID uint32) (ObjectState, int32) {
	var out C.JadeObjectState
	r := int32(C.get_window_state(C.uint32_t(windowID), &out))
	return objectState(out), r
}

func GetWebViewState(webviewID uint32) (ObjectState, int32) {
	var out C.JadeObjectState
	r := int32(C.get_webview_state(C.uint32_t(webviewID), &out))
	return objectState(out), r
}

func objectState(out C.JadeObjectState) ObjectState {
	return ObjectState{uint32(out.kind), uint32(out.state), uint32(out.parent_window_id), uint32(out.default_webview_id), int32(out.error_code)}
}

func GetWindowWebViews(windowID uint32) ([]uint32, int32) {
	var a C.WebViewIdArray
	r := int32(C.get_window_webviews(C.uint32_t(windowID), &a))
	if r != 0 {
		return nil, r
	}
	defer C.free_webview_id_array(&a)
	return append([]uint32(nil), unsafe.Slice((*uint32)(unsafe.Pointer(a.ids)), int(a.count))...), 0
}
