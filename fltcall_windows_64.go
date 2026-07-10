//go:build windows && (amd64 || arm64)

package jadeview

// 64 位下 double 参数走浮点寄存器（x64 XMM1 / arm64 D1），syscall 只装填整数寄存器，
// 故经一段运行时生成的 8 字节跳板转发：跳板把第 2 个整数寄存器里的位模式装入
// 对应浮点寄存器后跳转目标函数。W^X：先 RW 写入，再改 RX 并刷新指令缓存。
// （386 的 double 直接走栈，见 fltcall_windows_386.go。）

import (
	"math"
	"runtime"
	"sync"
	"syscall"
	"unsafe"
)

var (
	// kernel32 进程启动即已加载，按名加载无搜索路径劫持风险
	kernel32              = syscall.NewLazyDLL("kernel32.dll")
	procVirtualAlloc      = kernel32.NewProc("VirtualAlloc")
	procVirtualProtect    = kernel32.NewProc("VirtualProtect")
	procFlushInstrCache   = kernel32.NewProc("FlushInstructionCache")
	procGetCurrentProcess = kernel32.NewProc("GetCurrentProcess")
)

const (
	memCommitReserve = 0x3000 // MEM_COMMIT | MEM_RESERVE
	pageReadwrite    = 0x04
	pageExecuteRead  = 0x20
)

// fltStubCode 返回本架构的跳板机器码。
//   - amd64: movq xmm1, rdx; jmp r8   —— 目标函数看到 RCX=id、XMM1=double
//   - arm64: fmov d1, x1;    br x2    —— 目标函数看到 W0=id、D1=double（小端字节序）
func fltStubCode() []byte {
	if runtime.GOARCH == "arm64" {
		return []byte{
			0x21, 0x00, 0x67, 0x9E, // fmov d1, x1
			0x40, 0x00, 0x1F, 0xD6, // br   x2
		}
	}
	return []byte{
		0x66, 0x48, 0x0F, 0x6E, 0xCA, // movq xmm1, rdx
		0x41, 0xFF, 0xE0, //             jmp  r8
	}
}

var fltStubOnce = sync.OnceValues(func() (uintptr, error) {
	code := fltStubCode()
	addr, _, err := procVirtualAlloc.Call(0, uintptr(len(code)), memCommitReserve, pageReadwrite)
	if addr == 0 {
		return 0, err
	}
	copy(unsafe.Slice((*byte)(foreignPtr(addr)), len(code)), code)
	var old uint32
	r, _, err := procVirtualProtect.Call(addr, uintptr(len(code)), pageExecuteRead, uintptr(unsafe.Pointer(&old)))
	if r == 0 {
		return 0, err
	}
	hproc, _, _ := procGetCurrentProcess.Call()
	procFlushInstrCache.Call(hproc, addr, uintptr(len(code)))
	return addr, nil
})

// callF64I32 调用签名为 int32(uint32, double) 的导出函数。
// 跳板约定：arg1=id（原位透传）、arg2=double 位模式、arg3=目标函数地址。
func callF64I32(p *jvProc, id uint32, f float64) int32 {
	stub, err := fltStubOnce()
	if err != nil {
		return 0
	}
	if p.Find() != nil {
		return 0
	}
	r, _, _ := syscall.SyscallN(stub, uintptr(id), uintptr(math.Float64bits(f)), p.Addr())
	return i32(r)
}
