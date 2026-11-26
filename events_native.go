//go:build windows || (linux && amd64 && !jadeview_cgo)

package jadeview

// 固定槽位为命名事件和 IPC 提供原生回调入口。
// Windows x86 使用 stdcall；Linux x64 使用 System V ABI。

import (
	"runtime"
	"sync"
	"unsafe"
)

type eventReg struct {
	event    string
	handler  EventHandler
	cbID     uint32
	retiring bool
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
		evTramps[slot] = newNativeCallback(func(sourceID, data uintptr) uintptr {
			return goEventDispatch(s, sourceID, data)
		})
	}
	return evTramps[slot]
}

func goEventDispatch(slot int, sourceID, data uintptr) uintptr {
	evMu.Lock()
	reg := evSlots[slot]
	evMu.Unlock()
	if reg == nil || reg.handler == nil {
		return 0
	}
	resp := reg.handler(uint32(sourceID), goString(data))
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
// 第二个返回值为 false 表示槽位已满（超过 MaxEventHandlers）或库注册失败。
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
	if cbID == 0 {
		// 库注册失败（callback_id 与 window_id/tray_id 同约定：0=失败），回收槽位
		evSlots[slot] = nil
		evMu.Unlock()
		return 0, false
	}
	reg.cbID = cbID // 锁内赋值，避免与 Off 的读取竞态
	evMu.Unlock()
	return cbID, true
}

// Off 异步注销事件处理器；只有原生确认在途回调已结束后才复用槽位。
func Off(event string, cbID uint32) bool {
	evMu.Lock()
	slot := -1
	for i, r := range evSlots {
		if r != nil && !r.retiring && r.event == event && r.cbID == cbID && cbID != 0 {
			slot = i
			r.retiring = true
			break
		}
	}
	evMu.Unlock()
	if slot < 0 {
		return false
	}

	pool := &cstrs{}
	lo, hi, _ := procJadeOffWait.Call(uintptr(unsafe.Pointer(pool.p(event))), uintptr(cbID))
	runtime.KeepAlive(pool)
	id := requestID(lo, hi)
	if id == 0 {
		evMu.Lock()
		evSlots[slot].retiring = false
		evMu.Unlock()
		return false
	}
	go drainEventRequest(id, func() {
		evMu.Lock()
		evSlots[slot] = nil
		evMu.Unlock()
	})
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

// --- 网页权限处理器 ---

var permissionCallbackOnce sync.Once
var permissionCallback uintptr

func setPermissionCallback(enable bool) bool {
	if !enable {
		r, _, _ := procClearWebviewPermissionHandler.Call()
		return i32(r) == 1
	}
	permissionCallbackOnce.Do(func() {
		permissionCallback = newNativeCallback(func(view, data uintptr) uintptr {
			decision := permissionDecision(uint32(view), goString(data))
			return uintptr(uint32(decision))
		})
	})
	r, _, _ := procSetWebviewPermissionHandler.Call(permissionCallback)
	return i32(r) == 1
}
