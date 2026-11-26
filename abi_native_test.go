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

func TestSplashConfigLayout(t *testing.T) {
	wantSize, wantPath, wantDuration := uintptr(56), uintptr(16), uintptr(48)
	if unsafe.Sizeof(uintptr(0)) == 4 {
		wantSize, wantPath, wantDuration = 36, 12, 32
	}
	var s cSplashOptions
	if unsafe.Sizeof(s) != wantSize || unsafe.Offsetof(s.path) != wantPath || unsafe.Offsetof(s.minimumMS) != wantDuration {
		t.Fatalf("Splash layout: %d/%d/%d", unsafe.Sizeof(s), unsafe.Offsetof(s.path), unsafe.Offsetof(s.minimumMS))
	}
	pool := &cstrs{}
	o := DefaultWindowOptions()
	o.Splash = &SplashOptions{LogoData: []byte{137, 80, 78, 71}, MinimumMS: 123, BackgroundColor: 0xff000000}
	f := o.features(pool)
	if f == nil || f.splash == nil || f.splash.size != uint32(wantSize) || f.splash.abi != ABIVersion ||
		f.splash.minimumMS != 123 || f.splash.dataSize != 4 || f.splash.data == nil || *f.splash.data != 137 {
		t.Fatalf("bad splash conversion: %+v", f)
	}
}

func TestResourceAndJAPKLayouts(t *testing.T) {
	if got := unsafe.Sizeof(ResourceCapabilities{}); got != 24 {
		t.Errorf("resource caps layout: %d", got)
	}
	if got := unsafe.Sizeof(JAPKMountCapabilities{}); got != 20 {
		t.Errorf("mount caps layout: %d", got)
	}
	want := uintptr(24)
	if unsafe.Sizeof(uintptr(0)) == 4 {
		want = 16
	}
	if got := unsafe.Sizeof(cResourceDelivery{}); got != want {
		t.Errorf("delivery layout: %d", got)
	}
}
