//go:build windows

package jadeview

// 事件回调桥接（Windows 纯 Go 版）：库的 IpcCallback 签名为 (window_id, data)，
// 不携带事件名，因此为每个注册槽位建一个 syscall.NewCallback 跳板（闭包捕获槽位号），
// 回调时查表找到对应的 Go 处理函数。NewCallback 创建的跳板进程内不可释放，
// 故用固定槽位池复用（上限 MaxEventHandlers），与 Linux cgo 版的 C 跳板设计一致。
//
// IpcCallback 是 __stdcall（386 下有意义），对应 syscall.NewCallback。

import (
	"runtime"
	"sync"
	"syscall"
	"unsafe"
)

type eventReg struct {
	event   string
	handler EventHandler
	cbID    uint32
}

var (
	evMu     sync.Mutex
	evSlots  [MaxEventHandlers]*eventReg
	evTramps [MaxEventHandlers]uintptr
)

// evTramp 返回槽位的跳板地址（惰性创建，调用方须持有 evMu）。
func evTramp(slot int) uintptr {
	if evTramps[slot] == 0 {
		s := slot
		evTramps[slot] = syscall.NewCallback(func(windowID, data uintptr) uintptr {
			return goEventDispatch(s, windowID, data)
		})
	}
	return evTramps[slot]
}

func goEventDispatch(slot int, windowID, data uintptr) uintptr {
	evMu.Lock()
	reg := evSlots[slot]
	evMu.Unlock()
	if reg == nil || reg.handler == nil {
		return 0
	}
	resp := reg.handler(uint32(windowID), goString(data))
	if resp == "" {
		return 0
	}
	// 用库自带的 jade_text_create 生成一份由库管理的副本回传。
	pool := &cstrs{}
	r, _, _ := procJadeTextCreate.Call(uintptr(unsafe.Pointer(pool.pAlways(resp))))
	runtime.KeepAlive(pool)
	return r
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
	tramp := evTramp(slot)
	evMu.Unlock()

	pool := &cstrs{}
	r, _, _ := procJadeOn.Call(uintptr(unsafe.Pointer(pool.p(event))), tramp)
	runtime.KeepAlive(pool)
	cbID := uint32(r)
	evMu.Lock()
	reg.cbID = cbID // 锁内赋值，避免与 Off 的读取竞态
	evMu.Unlock()
	return cbID, true
}

// Off 注销事件处理器。注销成功后才释放对应槽位。
func Off(event string, cbID uint32) bool {
	pool := &cstrs{}
	r, _, _ := procJadeOff.Call(uintptr(unsafe.Pointer(pool.p(event))), uintptr(cbID))
	runtime.KeepAlive(pool)
	if i32(r) != 1 {
		return false
	}
	evMu.Lock()
	for i := range evSlots {
		if reg := evSlots[i]; reg != nil && reg.event == event && reg.cbID == cbID {
			evSlots[i] = nil
			break
		}
	}
	evMu.Unlock()
	return true
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
	tramp := evTramp(slot)
	evMu.Unlock()

	pool := &cstrs{}
	r, _, _ := procRegisterIPCHandler.Call(uintptr(unsafe.Pointer(pool.p(channel))), tramp)
	runtime.KeepAlive(pool)
	if i32(r) != 1 {
		// 注册失败，释放槽位，避免永久泄漏
		evMu.Lock()
		evSlots[slot] = nil
		evMu.Unlock()
		return false
	}
	return true
}
