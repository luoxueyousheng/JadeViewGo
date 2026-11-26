//go:build !jadeview_cgo

package jadeview

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/ebitengine/purego"
)

func TestLinuxLibraryOverride(t *testing.T) {
	t.Setenv("JADEVIEW_LIBRARY", "./custom/libJadeView.so")
	got, err := linuxLibraryPath()
	want, _ := filepath.Abs("./custom/libJadeView.so")
	if err != nil || got != want {
		t.Fatalf("library path = %q, %v; want %q", got, err, want)
	}
}

func TestLinuxRejectsIncompatibleLibrary(t *testing.T) {
	lib, _, err := loadLinuxLibrary("libc.so.6")
	if lib != 0 || err == nil || !strings.Contains(err.Error(), "jade_get_abi_info") {
		t.Fatalf("loadLinuxLibrary: handle=%#x, err=%v", lib, err)
	}
}

func TestLinuxCallbackOnForeignThread(t *testing.T) {
	lib, err := purego.Dlopen("libc.so.6", purego.RTLD_NOW|purego.RTLD_LOCAL)
	if err != nil {
		t.Fatal(err)
	}
	defer purego.Dlclose(lib)

	var create func(*uintptr, uintptr, uintptr, uintptr) int32
	var join func(uintptr, *uintptr) int32
	purego.RegisterLibFunc(&create, lib, "pthread_create")
	purego.RegisterLibFunc(&join, lib, "pthread_join")
	callback := newNativeCallback(func(arg uintptr) uintptr { return arg + 1 })
	var thread, value uintptr
	if rc := create(&thread, 0, callback, 41); rc != 0 {
		t.Fatalf("pthread_create: %d", rc)
	}
	if rc := join(thread, &value); rc != 0 || value != 42 {
		t.Fatalf("pthread_join: %d, callback result: %d", rc, value)
	}
}
