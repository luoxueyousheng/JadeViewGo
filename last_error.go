//go:build windows || (linux && amd64 && (!jadeview_cgo || cgo))

package jadeview

import "runtime"

// LastError 返回当前 OS 线程上的原生错误快照。
// 在调用可能失败的 API 之前先锁定 OS 线程，并保持到获取快照之后；
// 否则 goroutine 可能在不同线程调用 API 与读取错误。
func LastError() (code int32, message string) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	code = LastErrorCode()
	message = LastErrorMessage()
	return
}
