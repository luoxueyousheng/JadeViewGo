//go:build linux

// Linux(cgo) 实现：生命周期函数（初始化、版本号、消息循环、退出）。
// 静态链接 libJadeView.a，链接参数见 jadeview_linux_*.go。
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

// Preload 与 Windows 版对齐的占位实现：Linux 侧 libJadeView.a 为静态链接，
// 无运行时加载步骤，恒返回 nil。跨平台代码可无条件调用。
func Preload() error { return nil }

// RunMessageLoop 运行消息循环（阻塞，直到窗口全部关闭/退出）。
func RunMessageLoop() {
	C.run_message_loop()
}

// Exit 清理所有窗口并结束消息循环。
func Exit() {
	C.jadeview_exit()
}
