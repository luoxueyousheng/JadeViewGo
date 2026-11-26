//go:build windows

package jadeview

// Windows 纯 Go 实现的核心：DLL 定位/释放/加载与调用辅助函数。
//
// 不使用 cgo——go:embed 内置 MSVC 版 JadeView.dll，运行时释放到临时目录，
// syscall.NewLazyDLL 按绝对路径加载，各 API 经 LazyProc.Call 直调。
// 构建只需 Go 工具链，amd64/386/arm64 三架构均为单 exe 自包含分发。
//
// 结构体布局、回调约定与 C 头文件的对应关系见共用的 *_native.go；
// 布局已用 C(offsetof) 与 Go(unsafe.Offsetof) 双端比对验证（amd64/386 逐字段一致，
// arm64 与 amd64 对齐规则相同）。

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"
	"unsafe"
)

// jadeDLL 惰性加载 JadeView.dll：首次调用任一 API（或 Preload）时才释放文件并
// LoadLibrary，仅 import 本包不产生任何磁盘副作用。
// 路径解析优先级：exe 同目录的 JadeView.dll（便于手动覆盖/调试）→
// 释放内置副本到 %TEMP%\jadeview\<arch>-<内容哈希前8位>\（内容寻址，
// 不同版本/架构各占一目录；已存在文件按完整 sha256 校验，不符则重写）。
var (
	dllOnce sync.Once
	jadeDLL *syscall.LazyDLL
	abiOnce sync.Once
	abiErr  error
)

func loadedDLL() *syscall.LazyDLL {
	dllOnce.Do(func() { jadeDLL = syscall.NewLazyDLL(resolveDLLPath()) })
	return jadeDLL
}

// Preload 提前释放并加载内置 JadeView.dll，返回加载错误（nil=成功）。
// 可选调用：不调用时首次 API 调用会自动加载，但那时加载失败会 panic
// （syscall.LazyProc 语义）；需要优雅处理加载失败（临时目录禁止执行、
// 杀软拦截等）的宿主应在启动早期调用本函数并检查错误。
func Preload() error {
	abiOnce.Do(func() {
		d := loadedDLL()
		if err := d.Load(); err != nil {
			// 本库导入 GetWindowSubclass，未声明 Common Controls v6 的
			// Go 宿主可能默认得到了没有此导出的 comctl32 v5。
			if retryErr := loadWithCommonControls(d); retryErr != nil {
				abiErr = fmt.Errorf("%w；启用 Common Controls v6 后重试: %v", err, retryErr)
				return
			}
		}
		p := d.NewProc("jade_get_abi_info")
		if err := p.Find(); err != nil {
			abiErr = fmt.Errorf("JadeView v3 ABI (%s): %w", d.Name, err)
			return
		}
		var info abiInfo
		r, _, _ := p.Call(uintptr(unsafe.Pointer(&info)), unsafe.Sizeof(info))
		if i32(r) != 0 {
			abiErr = fmt.Errorf("jade_get_abi_info (%s): %d", d.Name, i32(r))
			return
		}
		abiErr = checkABI(info, unsafe.Sizeof(cWindowOptions{}), unsafe.Sizeof(cWebViewSettings{}))
	})
	return abiErr
}

func resolveDLLPath() string {
	if exe, err := os.Executable(); err == nil {
		p := filepath.Join(filepath.Dir(exe), "JadeView.dll")
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			return p
		}
	}
	sum := sha256.Sum256(embeddedJadeViewDLL)
	dir := filepath.Join(os.TempDir(), "jadeview", runtime.GOARCH+"-"+hex.EncodeToString(sum[:4]))
	dst := filepath.Join(dir, "JadeView.dll")
	// 目录名已含内容哈希，但仍校验磁盘文件的完整 sha256：同大小篡改/写坏残留不会被加载
	if b, err := os.ReadFile(dst); err != nil || sha256.Sum256(b) != sum {
		_ = os.MkdirAll(dir, 0o755)
		// 写失败不致命（可能另一进程刚写完并占用）；加载失败会在 Preload/首次调用时报错
		_ = os.WriteFile(dst, embeddedJadeViewDLL, 0o644)
	}
	return dst
}

// jvProc 惰性解析的导出函数句柄：首次 Call/Find/Addr 时才触发 DLL 释放与加载
// （见 loadedDLL），方法集与 syscall.LazyProc 的使用面一致。
type jvProc struct {
	name string
	once sync.Once
	proc *syscall.LazyProc
}

// np 声明一个导出函数（仅记录名字，不触发加载）。
func np(name string) *jvProc { return &jvProc{name: name} }

func (p *jvProc) resolve() *syscall.LazyProc {
	if err := Preload(); err != nil {
		panic(err)
	}
	p.once.Do(func() { p.proc = loadedDLL().NewProc(p.name) })
	return p.proc
}

//go:uintptrescapes
func (p *jvProc) Call(a ...uintptr) (uintptr, uintptr, error) { return p.resolve().Call(a...) }
func (p *jvProc) Find() error                                 { return p.resolve().Find() }
func (p *jvProc) Addr() uintptr                               { return p.resolve().Addr() }

func newNativeCallback(f any) uintptr      { return syscall.NewCallback(f) }
func newNativeCDeclCallback(f any) uintptr { return syscall.NewCallbackCDecl(f) }
