//go:build windows || (linux && amd64 && !jadeview_cgo)

package jadeview

import (
	"testing"
	"unsafe"
)

func TestNativeConfigLayout(t *testing.T) {
	type layout struct {
		windowSize, backOff, iconOff, settingsSize, profileOff, zoomOff, featuresSize, menuSize, menuIconOff, menuIconSize uintptr
	}
	want := layout{160, 144, 152, 128, 96, 112, 24, 72, 40, 32}
	if unsafe.Sizeof(uintptr(0)) == 4 {
		want = layout{124, 112, 120, 88, 64, 72, 20, 40, 24, 16}
	}
	var w cWindowOptions
	var s cWebViewSettingsFields
	got := layout{unsafe.Sizeof(w), unsafe.Offsetof(w.titlebarBack), unsafe.Offsetof(w.titlebarIcon),
		unsafe.Sizeof(cWebViewSettings{}), unsafe.Offsetof(s.profileName), unsafe.Offsetof(s.zoomFactor),
		unsafe.Sizeof(cWindowFeatures{}), unsafe.Sizeof(cMenuItem{}), unsafe.Offsetof(cMenuItem{}.icon), unsafe.Sizeof(cMenuIcon{})}
	if got != want {
		t.Fatalf("layout %+v, want %+v", got, want)
	}

	wp := &cstrs{}
	o := DefaultWindowOptions()
	o.HideTitlebarRefresh = true
	o.NoFrame = true
	co := o.toC(wp)
	cs := (&WebViewSettings{ZoomFactor: 1.25, ProfileName: "abi-test"}).toC(wp)
	if co.structSize != uint32(want.windowSize) || co.abiVersion != ABIVersion ||
		co.noframe != 1 || co.titlebarBack != 1 || co.titlebarRefresh != 0 || co.titlebarIcon != 1 ||
		cs.structSize != uint32(want.settingsSize) || cs.abiVersion != ABIVersion ||
		cs.zoomFactor != 1.25 || cs.profileName == nil {
		t.Fatal("unexpected native config translation")
	}
}
