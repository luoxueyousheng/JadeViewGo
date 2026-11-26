//go:build windows || (linux && amd64 && (!jadeview_cgo || cgo))

package jadeview

import "time"

// 等待期间不占用原生回调线程，允许从自身的处理器中调用 Off。
func drainEventRequest(id RequestID, releaseSlot func()) {
	ticker := time.NewTicker(20 * time.Millisecond)
	defer ticker.Stop()
	timer := time.NewTimer(30 * time.Second)
	defer timer.Stop()
	defer ReleaseRequest(id)
	for {
		r, code := GetRequestResult(id)
		if code != 0 || r.State == RequestFailed || r.State == RequestCancelled {
			return
		}
		if r.State == RequestSucceeded {
			releaseSlot()
			return
		}
		select {
		case <-ticker.C:
		case <-timer.C:
			return
		}
	}
}

// OffImmediate 映射原生 jade_off：提交注销后立即返回。已进入的回调可能继续执行。
// 此接口没有完成请求，无法安全复用原槽位；槽位留空隔离到进程退出。
// 普通调用优先使用 Off，它会等待原生确认后自动回收槽位。
func OffImmediate(event string, id uint32) bool {
	if id == 0 {
		return false
	}
	evMu.Lock()
	slot := -1
	for i, r := range evSlots {
		if r != nil && !r.retiring && r.event == event && r.cbID == id {
			slot = i
			r.retiring = true
			break
		}
	}
	evMu.Unlock()
	if slot < 0 {
		return false
	}

	ok := nativeNamedOff(event, id)
	evMu.Lock()
	if ok {
		// 保留跳板占位，不保留用户 handler；已取到旧记录的回调持有独立引用。
		evSlots[slot] = &eventReg{retiring: true}
	} else {
		evSlots[slot].retiring = false
	}
	evMu.Unlock()
	return ok
}
