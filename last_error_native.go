//go:build windows || (linux && amd64 && !jadeview_cgo)

package jadeview

var (
	procLastErrorCode    = np("jade_last_error_code")
	procLastErrorMessage = np("jade_last_error_message")
)

// LastErrorCode 查询当前 OS 线程的错误码。
func LastErrorCode() int32 { r, _, _ := procLastErrorCode.Call(); return i32(r) }

// LastErrorMessage 查询当前 OS 线程的错误消息，并释放原生字符串。
func LastErrorMessage() string { r, _, _ := procLastErrorMessage.Call(); return goStringFree(r) }
