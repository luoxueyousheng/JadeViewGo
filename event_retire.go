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
