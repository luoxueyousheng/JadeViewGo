//go:build !jadeview_cgo

package jadeview

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
)

const linuxLibrary = "libJadeView.so"

var (
	linuxLoadOnce sync.Once
	linuxHandle   uintptr
	linuxLoadErr  error
)

// Preload 加载并校验 JadeView ABI。优先使用 JADEVIEW_LIBRARY 指定的文件，
// 然后查找可执行文件旁的 libJadeView.so，最后交给系统动态加载器。
// 加载成功的库保留至进程退出，确保原生线程与回调入口仍有效。
func Preload() error {
	linuxLoadOnce.Do(func() {
		path, err := linuxLibraryPath()
		if err != nil {
			linuxLoadErr = err
			return
		}
		linuxHandle, loadedABI, linuxLoadErr = loadLinuxLibrary(path)
	})
	return linuxLoadErr
}

func linuxLibraryPath() (string, error) {
	if override := os.Getenv("JADEVIEW_LIBRARY"); override != "" {
		absolute, err := filepath.Abs(override)
		if err != nil {
			return "", fmt.Errorf("JADEVIEW_LIBRARY: %w", err)
		}
		return absolute, nil
	}
	if exe, err := os.Executable(); err == nil {
		path := filepath.Join(filepath.Dir(exe), linuxLibrary)
		if stat, err := os.Stat(path); err == nil && !stat.IsDir() {
			return path, nil
		} else if err != nil && !os.IsNotExist(err) {
			return "", err
		}
	}
	return linuxLibrary, nil
}

func loadLinuxLibrary(path string) (uintptr, abiInfo, error) {
	handle, err := purego.Dlopen(path, purego.RTLD_NOW|purego.RTLD_LOCAL)
	if err != nil {
		return 0, abiInfo{}, fmt.Errorf("加载 %s: %w", path, err)
	}
	info, err := readLinuxABI(handle)
	if err != nil {
		_ = purego.Dlclose(handle)
		return 0, abiInfo{}, fmt.Errorf("%s: %w", path, err)
	}
	return handle, info, nil
}

func readLinuxABI(handle uintptr) (abiInfo, error) {
	addr, err := purego.Dlsym(handle, "jade_get_abi_info")
	if err != nil {
		return abiInfo{}, fmt.Errorf("缺少 jade_get_abi_info: %w", err)
	}
	var getABI func(*abiInfo, uintptr) int32
	purego.RegisterFunc(&getABI, addr)
	var info abiInfo
	if code := getABI(&info, unsafe.Sizeof(info)); code != 0 {
		return abiInfo{}, fmt.Errorf("jade_get_abi_info: %d", code)
	}
	return info, checkABI(info, unsafe.Sizeof(cWindowOptions{}), unsafe.Sizeof(cWebViewSettings{}))
}

type jvProc struct {
	name string
	once sync.Once
	addr uintptr
	err  error
}

func np(name string) *jvProc { return &jvProc{name: name} }

func (p *jvProc) Find() error {
	if err := Preload(); err != nil {
		return err
	}
	p.once.Do(func() { p.addr, p.err = purego.Dlsym(linuxHandle, p.name) })
	return p.err
}

func (p *jvProc) Addr() uintptr {
	if err := p.Find(); err != nil {
		panic(err)
	}
	return p.addr
}

//go:uintptrescapes
func (p *jvProc) Call(args ...uintptr) (uintptr, uintptr, error) {
	// 导出 API 使用各自返回值表达错误，不能按 libc errno 判断。
	r1, r2, _ := purego.SyscallN(p.Addr(), args...)
	return r1, r2, nil
}

func newNativeCallback(f any) uintptr      { return purego.NewCallback(f) }
func newNativeCDeclCallback(f any) uintptr { return purego.NewCallback(f) }

type cWebViewSettings = cWebViewSettingsFields

var (
	linuxZoomOnce sync.Once
	linuxZoom     func(uint32, float64) int32
)

func callF64I32(p *jvProc, id uint32, value float64) int32 {
	linuxZoomOnce.Do(func() {
		purego.RegisterFunc(&linuxZoom, p.Addr())
	})
	return linuxZoom(id, value)
}
