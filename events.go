package jadeview

// 事件回调桥接：库的 IpcCallback 签名为 (window_id, data)，不携带事件名，
// 因此为每个注册槽位生成一个独立的 C 跳板函数，槽位号编译期固定，
// 回调时把槽位号传回 Go，再查表找到对应的 Go 处理函数。
//
// 槽位总数 = 同时可注册的事件处理器上限（见下方 N）。

/*
#include "JadeView.h"

extern char* goEventDispatch(int slot, uint32_t window_id, char* data);

static const char* JADEVIEW_CALL jv_tr0(uint32_t w, const char* d){ return goEventDispatch(0, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr1(uint32_t w, const char* d){ return goEventDispatch(1, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr2(uint32_t w, const char* d){ return goEventDispatch(2, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr3(uint32_t w, const char* d){ return goEventDispatch(3, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr4(uint32_t w, const char* d){ return goEventDispatch(4, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr5(uint32_t w, const char* d){ return goEventDispatch(5, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr6(uint32_t w, const char* d){ return goEventDispatch(6, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr7(uint32_t w, const char* d){ return goEventDispatch(7, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr8(uint32_t w, const char* d){ return goEventDispatch(8, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr9(uint32_t w, const char* d){ return goEventDispatch(9, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr10(uint32_t w, const char* d){ return goEventDispatch(10, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr11(uint32_t w, const char* d){ return goEventDispatch(11, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr12(uint32_t w, const char* d){ return goEventDispatch(12, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr13(uint32_t w, const char* d){ return goEventDispatch(13, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr14(uint32_t w, const char* d){ return goEventDispatch(14, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr15(uint32_t w, const char* d){ return goEventDispatch(15, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr16(uint32_t w, const char* d){ return goEventDispatch(16, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr17(uint32_t w, const char* d){ return goEventDispatch(17, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr18(uint32_t w, const char* d){ return goEventDispatch(18, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr19(uint32_t w, const char* d){ return goEventDispatch(19, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr20(uint32_t w, const char* d){ return goEventDispatch(20, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr21(uint32_t w, const char* d){ return goEventDispatch(21, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr22(uint32_t w, const char* d){ return goEventDispatch(22, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr23(uint32_t w, const char* d){ return goEventDispatch(23, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr24(uint32_t w, const char* d){ return goEventDispatch(24, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr25(uint32_t w, const char* d){ return goEventDispatch(25, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr26(uint32_t w, const char* d){ return goEventDispatch(26, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr27(uint32_t w, const char* d){ return goEventDispatch(27, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr28(uint32_t w, const char* d){ return goEventDispatch(28, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr29(uint32_t w, const char* d){ return goEventDispatch(29, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr30(uint32_t w, const char* d){ return goEventDispatch(30, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr31(uint32_t w, const char* d){ return goEventDispatch(31, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr32(uint32_t w, const char* d){ return goEventDispatch(32, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr33(uint32_t w, const char* d){ return goEventDispatch(33, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr34(uint32_t w, const char* d){ return goEventDispatch(34, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr35(uint32_t w, const char* d){ return goEventDispatch(35, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr36(uint32_t w, const char* d){ return goEventDispatch(36, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr37(uint32_t w, const char* d){ return goEventDispatch(37, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr38(uint32_t w, const char* d){ return goEventDispatch(38, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr39(uint32_t w, const char* d){ return goEventDispatch(39, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr40(uint32_t w, const char* d){ return goEventDispatch(40, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr41(uint32_t w, const char* d){ return goEventDispatch(41, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr42(uint32_t w, const char* d){ return goEventDispatch(42, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr43(uint32_t w, const char* d){ return goEventDispatch(43, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr44(uint32_t w, const char* d){ return goEventDispatch(44, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr45(uint32_t w, const char* d){ return goEventDispatch(45, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr46(uint32_t w, const char* d){ return goEventDispatch(46, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr47(uint32_t w, const char* d){ return goEventDispatch(47, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr48(uint32_t w, const char* d){ return goEventDispatch(48, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr49(uint32_t w, const char* d){ return goEventDispatch(49, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr50(uint32_t w, const char* d){ return goEventDispatch(50, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr51(uint32_t w, const char* d){ return goEventDispatch(51, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr52(uint32_t w, const char* d){ return goEventDispatch(52, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr53(uint32_t w, const char* d){ return goEventDispatch(53, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr54(uint32_t w, const char* d){ return goEventDispatch(54, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr55(uint32_t w, const char* d){ return goEventDispatch(55, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr56(uint32_t w, const char* d){ return goEventDispatch(56, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr57(uint32_t w, const char* d){ return goEventDispatch(57, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr58(uint32_t w, const char* d){ return goEventDispatch(58, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr59(uint32_t w, const char* d){ return goEventDispatch(59, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr60(uint32_t w, const char* d){ return goEventDispatch(60, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr61(uint32_t w, const char* d){ return goEventDispatch(61, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr62(uint32_t w, const char* d){ return goEventDispatch(62, w, (char*)d); }
static const char* JADEVIEW_CALL jv_tr63(uint32_t w, const char* d){ return goEventDispatch(63, w, (char*)d); }
static IpcCallback jv_tramps[64] = {jv_tr0, jv_tr1, jv_tr2, jv_tr3, jv_tr4, jv_tr5, jv_tr6, jv_tr7, jv_tr8, jv_tr9, jv_tr10, jv_tr11, jv_tr12, jv_tr13, jv_tr14, jv_tr15, jv_tr16, jv_tr17, jv_tr18, jv_tr19, jv_tr20, jv_tr21, jv_tr22, jv_tr23, jv_tr24, jv_tr25, jv_tr26, jv_tr27, jv_tr28, jv_tr29, jv_tr30, jv_tr31, jv_tr32, jv_tr33, jv_tr34, jv_tr35, jv_tr36, jv_tr37, jv_tr38, jv_tr39, jv_tr40, jv_tr41, jv_tr42, jv_tr43, jv_tr44, jv_tr45, jv_tr46, jv_tr47, jv_tr48, jv_tr49, jv_tr50, jv_tr51, jv_tr52, jv_tr53, jv_tr54, jv_tr55, jv_tr56, jv_tr57, jv_tr58, jv_tr59, jv_tr60, jv_tr61, jv_tr62, jv_tr63};
static IpcCallback jv_get_tramp(int i){ return jv_tramps[i]; }
*/
import "C"

import (
	"sync"
	"unsafe"
)

// EventHandler 是事件处理函数。返回非空字符串会作为响应回传给库（多数事件可返回 ""）。
type EventHandler func(windowID uint32, data string) string

// MaxEventHandlers 是可同时注册的事件处理器上限（C 跳板槽位数）。
const MaxEventHandlers = 64

type eventReg struct {
	event   string
	handler EventHandler
	cbID    uint32
}

var (
	evMu    sync.Mutex
	evSlots [MaxEventHandlers]*eventReg
)

//export goEventDispatch
func goEventDispatch(slot C.int, windowID C.uint32_t, data *C.char) *C.char {
	evMu.Lock()
	reg := evSlots[int(slot)]
	evMu.Unlock()
	if reg == nil || reg.handler == nil {
		return nil
	}
	resp := reg.handler(uint32(windowID), C.GoString(data))
	if resp == "" {
		return nil
	}
	// 用库自带的 jade_text_create 生成一份由库管理的副本回传。
	tmp := C.CString(resp)
	out := C.jade_text_create(tmp)
	C.free(unsafe.Pointer(tmp))
	return out
}

// On 注册事件处理器，返回库分配的 callback_id（配合 Off 注销）。
// 第二个返回值为 false 表示槽位已满（超过 MaxEventHandlers）。
func On(event string, handler EventHandler) (uint32, bool) {
	evMu.Lock()
	slot := -1
	for i := range evSlots {
		if evSlots[i] == nil {
			slot = i
			break
		}
	}
	if slot == -1 {
		evMu.Unlock()
		return 0, false
	}
	reg := &eventReg{event: event, handler: handler}
	evSlots[slot] = reg
	evMu.Unlock()

	cev := C.CString(event)
	defer C.free(unsafe.Pointer(cev))
	cbID := uint32(C.jade_on(cev, C.jv_get_tramp(C.int(slot))))
	reg.cbID = cbID
	return cbID, true
}

// Off 注销事件处理器。
func Off(event string, cbID uint32) bool {
	evMu.Lock()
	for i := range evSlots {
		if r := evSlots[i]; r != nil && r.event == event && r.cbID == cbID {
			evSlots[i] = nil
			break
		}
	}
	evMu.Unlock()

	cev := C.CString(event)
	defer C.free(unsafe.Pointer(cev))
	return C.jade_off(cev, C.uint32_t(cbID)) == 1
}

// RegisterIPCHandler 订阅前端通过指定 channel 发送的 IPC 消息。
// 处理函数的返回值作为响应回传前端。
func RegisterIPCHandler(channel string, handler EventHandler) bool {
	evMu.Lock()
	slot := -1
	for i := range evSlots {
		if evSlots[i] == nil {
			slot = i
			break
		}
	}
	if slot == -1 {
		evMu.Unlock()
		return false
	}
	evSlots[slot] = &eventReg{event: "ipc:" + channel, handler: handler}
	evMu.Unlock()

	cch := C.CString(channel)
	defer C.free(unsafe.Pointer(cch))
	return C.register_ipc_handler(cch, C.jv_get_tramp(C.int(slot))) == 1
}
