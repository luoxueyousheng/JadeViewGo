//go:build linux && amd64 && jadeview_cgo && cgo

// Linux(cgo) 实现：生命周期函数（初始化、版本号、消息循环、退出）。
// 动态链接 libJadeView.so。
package jadeview

/*
#cgo CFLAGS: -I${SRCDIR}/include
#include <stdlib.h>
#include "JadeView.h"
*/
import "C"

import (
	"fmt"
	"sync"
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
//   - appSignature : 应用签名，trim 后至少 6 个 Unicode 字符；建议反域名格式
//     （如 com.example.myapp）——它会成为 JAPK 模式下 JADE:// URL 的主机名
//   - singleInstance: 是否启用单实例
//
// 返回 true 表示成功。注意：宿主仍应在 jade_on("app-ready", ...) 回调里
// 解析 JSON 的 ok 字段；回调来源 ID 为 0。
func Init(enableDevmod bool, logPath, dataDir, appName, appSignature string, singleInstance bool) bool {
	if Preload() != nil {
		return false
	}
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

var abiOnce sync.Once
var abiErr error

// Preload 校验动态链接的原生库 ABI。系统动态加载器在进程启动时解析共享库。
func Preload() error {
	abiOnce.Do(func() {
		var info C.JadeAbiInfo
		if code := C.jade_get_abi_info(&info, C.sizeof_JadeAbiInfo); code != 0 {
			abiErr = fmt.Errorf("jade_get_abi_info: %d", code)
			return
		}
		loadedABI = abiInfo{uint32(info.abi_version), uint32(info.pointer_bits), uint32(info.window_options_size), uint32(info.webview_settings_size)}
		abiErr = checkABI(loadedABI, C.sizeof_WebViewWindowOptions, C.sizeof_WebViewSettings)
	})
	return abiErr
}

// RunMessageLoop 运行消息循环（阻塞，直到窗口全部关闭/退出）。
func RunMessageLoop() {
	C.run_message_loop()
}

// Exit 清理所有窗口并结束消息循环。
func Exit() {
	C.jadeview_exit()
}

// ExitWait 发起退出并等待原生运行时收尾，返回库状态码。
func ExitWait(timeoutMs uint32) int32 {
	return int32(C.jadeview_exit_wait(C.uint32_t(timeoutMs)))
}
