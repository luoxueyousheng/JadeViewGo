// Package jadeview 是 JadeView WebView 库的 Go(cgo)封装。
//
// 本文件为最小骨架：仅封装初始化、版本号、消息循环、退出等少量函数，
// 用于先跑通编译与链接。完整 API（窗口/托盘/对话框/YAML/system 等）待后续补全。
//
// 链接方式：静态链接 libJadeView.a（见 jadeview_linux_amd64.go / _arm64.go）。
package jadeview

/*
#cgo CFLAGS: -I${SRCDIR}/include
#include <stdlib.h>
#include "JadeView.h"
*/
import "C"

import (
	"unsafe"
)

// boolToC 把 Go bool 转成库约定的 int32(0/1)。
func boolToC(b bool) C.int32_t {
	if b {
		return 1
	}
	return 0
}

// Init 初始化 JadeView。
//
//   - enableDevmod : 是否启用开发模式
//   - logPath      : 日志路径（可空字符串）
//   - dataDir      : 数据目录（可空字符串）
//   - appName      : 应用名，必填、非纯空白
//   - appSignature : 应用签名，trim 后至少 6 个 Unicode 字符
//   - singleInstance: 是否启用单实例
//
// 返回 true 表示成功。注意：宿主仍应在 jade_on("app-ready", ...) 回调里
// 判断 window_id 以确认初始化结果（详见头文件说明）。
func Init(enableDevmod bool, logPath, dataDir, appName, appSignature string, singleInstance bool) bool {
	cLog := C.CString(logPath)
	cData := C.CString(dataDir)
	cName := C.CString(appName)
	cSig := C.CString(appSignature)
	defer C.free(unsafe.Pointer(cLog))
	defer C.free(unsafe.Pointer(cData))
	defer C.free(unsafe.Pointer(cName))
	defer C.free(unsafe.Pointer(cSig))

	ret := C.JadeView_init(
		boolToC(enableDevmod),
		cLog, cData, cName, cSig,
		boolToC(singleInstance),
	)
	return ret == 1
}

// Version 返回 JadeView 版本字符串。
func Version() string {
	const bufSize = 256
	buf := make([]byte, bufSize)
	ret := C.jadeview_version((*C.char)(unsafe.Pointer(&buf[0])), C.size_t(bufSize))
	if ret != 1 {
		return ""
	}
	return cBufToString(buf)
}

// RunMessageLoop 运行消息循环（阻塞，直到窗口全部关闭/退出）。
func RunMessageLoop() {
	C.run_message_loop()
}

// Exit 清理所有窗口并结束消息循环。
func Exit() {
	C.jadeview_exit()
}

// cBufToString 把以 NUL 结尾的 C 缓冲区转成 Go string。
func cBufToString(buf []byte) string {
	for i, b := range buf {
		if b == 0 {
			return string(buf[:i])
		}
	}
	return string(buf)
}
