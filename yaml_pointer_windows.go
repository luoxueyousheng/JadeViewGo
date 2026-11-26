package jadeview

import (
	"runtime"
	"syscall"
	"unsafe"
)

var (
	procYAMLGetStr    = np("yaml_get_str")
	procCoTaskMemFree = syscall.NewLazyDLL("ole32.dll").NewProc("CoTaskMemFree")
)

// YAMLGetStr 返回路径值的 JSON 文本（不限定 YAML 值类型）。不存在时返回 false。
func YAMLGetStr(fileName, keyPath string) (string, bool) {
	pool := &cstrs{}
	r, _, _ := procYAMLGetStr.Call(uintptr(unsafe.Pointer(pool.p(fileName))), uintptr(unsafe.Pointer(pool.p(keyPath))))
	runtime.KeepAlive(pool)
	if r == 0 {
		return "", false
	}
	defer procCoTaskMemFree.Call(r)
	return goString(r), true
}
