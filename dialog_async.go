//go:build linux

package jadeview

// 异步对话框回调桥：异步对话框完成后通过 void(const char*) 回调回传结果 JSON。
// 与事件桥同理，为每个并发槽位生成一个 C 跳板，槽位号编译期固定，回调时传回 Go 查表。

/*
#include "JadeView.h"

extern void goDialogDispatch(int slot, char* result);

static void jv_dlg_tr0(const char* r){ goDialogDispatch(0, (char*)r); }
static void jv_dlg_tr1(const char* r){ goDialogDispatch(1, (char*)r); }
static void jv_dlg_tr2(const char* r){ goDialogDispatch(2, (char*)r); }
static void jv_dlg_tr3(const char* r){ goDialogDispatch(3, (char*)r); }
static void jv_dlg_tr4(const char* r){ goDialogDispatch(4, (char*)r); }
static void jv_dlg_tr5(const char* r){ goDialogDispatch(5, (char*)r); }
static void jv_dlg_tr6(const char* r){ goDialogDispatch(6, (char*)r); }
static void jv_dlg_tr7(const char* r){ goDialogDispatch(7, (char*)r); }
static void jv_dlg_tr8(const char* r){ goDialogDispatch(8, (char*)r); }
static void jv_dlg_tr9(const char* r){ goDialogDispatch(9, (char*)r); }
static void jv_dlg_tr10(const char* r){ goDialogDispatch(10, (char*)r); }
static void jv_dlg_tr11(const char* r){ goDialogDispatch(11, (char*)r); }
static void jv_dlg_tr12(const char* r){ goDialogDispatch(12, (char*)r); }
static void jv_dlg_tr13(const char* r){ goDialogDispatch(13, (char*)r); }
static void jv_dlg_tr14(const char* r){ goDialogDispatch(14, (char*)r); }
static void jv_dlg_tr15(const char* r){ goDialogDispatch(15, (char*)r); }
typedef void (*JvDlgCb)(const char*);
static JvDlgCb jv_dlg_tramps[16] = {jv_dlg_tr0, jv_dlg_tr1, jv_dlg_tr2, jv_dlg_tr3, jv_dlg_tr4, jv_dlg_tr5, jv_dlg_tr6, jv_dlg_tr7, jv_dlg_tr8, jv_dlg_tr9, jv_dlg_tr10, jv_dlg_tr11, jv_dlg_tr12, jv_dlg_tr13, jv_dlg_tr14, jv_dlg_tr15};
static JvDlgCb jv_dlg_get_tramp(int i){ return jv_dlg_tramps[i]; }
*/
import "C"

import "sync"

var (
	dlgMu    sync.Mutex
	dlgSlots [MaxAsyncDialogs]DialogResultHandler
)

//export goDialogDispatch
func goDialogDispatch(slot C.int, result *C.char) {
	dlgMu.Lock()
	h := dlgSlots[int(slot)]
	dlgSlots[int(slot)] = nil // 一次性，回调后释放槽位
	dlgMu.Unlock()
	if h != nil {
		h(C.GoString(result))
	}
}

// dlgAcquire 占用一个空闲槽位，返回 (slot, ok)。
func dlgAcquire(h DialogResultHandler) (C.int, bool) {
	dlgMu.Lock()
	defer dlgMu.Unlock()
	for i := range dlgSlots {
		if dlgSlots[i] == nil {
			dlgSlots[i] = h
			return C.int(i), true
		}
	}
	return 0, false
}

func dlgRelease(slot C.int) {
	dlgMu.Lock()
	dlgSlots[int(slot)] = nil
	dlgMu.Unlock()
}

// ShowOpenDialogAsync 异步显示打开文件对话框，结果通过 handler 回传。
// 返回 false 表示槽位已满或库调用失败。
func ShowOpenDialogAsync(p FileDialogParams, handler DialogResultHandler) bool {
	slot, ok := dlgAcquire(handler)
	if !ok {
		return false
	}
	pool := &cstrPool{}
	defer pool.free()
	c := p.toC(pool)
	if C.jade_dialog_show_open_dialog_async(&c, C.jv_dlg_get_tramp(slot)) != 1 {
		dlgRelease(slot)
		return false
	}
	return true
}

// ShowSaveDialogAsync 异步显示保存文件对话框。
func ShowSaveDialogAsync(p FileDialogParams, handler DialogResultHandler) bool {
	slot, ok := dlgAcquire(handler)
	if !ok {
		return false
	}
	pool := &cstrPool{}
	defer pool.free()
	c := p.toC(pool)
	if C.jade_dialog_show_save_dialog_async(&c, C.jv_dlg_get_tramp(slot)) != 1 {
		dlgRelease(slot)
		return false
	}
	return true
}

// ShowMessageBoxAsync 异步显示消息框。
func ShowMessageBoxAsync(p MessageBoxParams, handler DialogResultHandler) bool {
	slot, ok := dlgAcquire(handler)
	if !ok {
		return false
	}
	pool := &cstrPool{}
	defer pool.free()
	c := C.MessageBoxParams{
		window_id:  C.uint32_t(p.WindowID),
		title:      pool.s(p.Title),
		message:    pool.s(p.Message),
		detail:     pool.s(p.Detail),
		buttons:    pool.s(p.Buttons),
		default_id: C.int32_t(p.DefaultID),
		cancel_id:  C.int32_t(p.CancelID),
		type_:      pool.s(p.Type),
	}
	if C.jade_dialog_show_message_box_async(&c, C.jv_dlg_get_tramp(slot)) != 1 {
		dlgRelease(slot)
		return false
	}
	return true
}
