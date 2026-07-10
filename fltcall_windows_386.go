//go:build windows && 386

package jadeview

import (
	"math"
)

// callF64I32 调用签名为 int32(uint32, double) 的导出函数。
// 386 stdcall 的 double 直接走栈（8 字节拆两个字压栈），无需跳板。
func callF64I32(p *jvProc, id uint32, f float64) int32 {
	b := math.Float64bits(f)
	r, _, _ := p.Call(uintptr(id), uintptr(uint32(b)), uintptr(uint32(b>>32)))
	return i32(r)
}
