//go:build windows || (linux && amd64 && !jadeview_cgo)

package jadeview

import "unsafe"

// Windows x86 在 stdcall 栈上以低字、高字顺序传递 uint64。
//
//go:uintptrescapes
func callU64First(p *jvProc, id uint64, tail ...uintptr) (uintptr, uintptr, error) {
	if unsafe.Sizeof(uintptr(0)) == 4 {
		return p.Call(append([]uintptr{uintptr(uint32(id)), uintptr(uint32(id >> 32))}, tail...)...)
	}
	return p.Call(append([]uintptr{uintptr(id)}, tail...)...)
}
