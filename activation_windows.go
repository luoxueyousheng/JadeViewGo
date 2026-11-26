package jadeview

import (
	"fmt"
	"os"
	"runtime"
	"syscall"
	"unsafe"
)

const commonControlsManifest = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly xmlns="urn:schemas-microsoft-com:asm.v1" manifestVersion="1.0">
 <assemblyIdentity type="win32" name="JadeViewGo.NativeLoader" version="1.0.0.0" processorArchitecture="*"/>
 <dependency><dependentAssembly><assemblyIdentity type="win32" name="Microsoft.Windows.Common-Controls"
  version="6.0.0.0" processorArchitecture="*" publicKeyToken="6595b64144ccf1df" language="*"/></dependentAssembly></dependency>
</assembly>`

type activationContext struct {
	size, flags                                  uint32
	source                                       *uint16
	architecture, language                       uint16
	assemblyDirectory, resourceName, application *uint16
	module                                       uintptr
}

func loadWithCommonControls(dll *syscall.LazyDLL) error {
	f, err := os.CreateTemp("", "jadeview-controls-*.manifest")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	_, writeErr := f.WriteString(commonControlsManifest)
	closeErr := f.Close()
	if writeErr != nil {
		return writeErr
	}
	if closeErr != nil {
		return closeErr
	}
	path, err := syscall.UTF16PtrFromString(f.Name())
	if err != nil {
		return err
	}

	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	ctx := activationContext{source: path}
	ctx.size = uint32(unsafe.Sizeof(ctx))
	handle, _, err := kernel32.NewProc("CreateActCtxW").Call(uintptr(unsafe.Pointer(&ctx)))
	if handle == ^uintptr(0) {
		return fmt.Errorf("CreateActCtxW: %w", err)
	}
	defer kernel32.NewProc("ReleaseActCtx").Call(handle)

	// 激活上下文属于当前 OS 线程，必须在同一线程加载并恢复。
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	var cookie uintptr
	ok, _, err := kernel32.NewProc("ActivateActCtx").Call(handle, uintptr(unsafe.Pointer(&cookie)))
	if ok == 0 {
		return fmt.Errorf("ActivateActCtx: %w", err)
	}
	defer kernel32.NewProc("DeactivateActCtx").Call(0, cookie)
	return dll.Load()
}
