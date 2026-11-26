//go:build windows || (linux && amd64 && (!jadeview_cgo || cgo))

package jadeview

import (
	"strconv"
	"strings"
	"testing"
)

func TestABIValidation(t *testing.T) {
	good := abiInfo{ABIVersion, uint32(strconv.IntSize), 160, 128}
	if err := checkABI(good, 160, 128); err != nil {
		t.Fatal(err)
	}
	for _, bad := range []abiInfo{
		{ABIVersion - 1, good.pointerBits, 160, 128},
		{ABIVersion, 0, 160, 128},
		{ABIVersion, good.pointerBits, 124, 128},
		{ABIVersion, good.pointerBits, 160, 88},
	} {
		if err := checkABI(bad, 160, 128); err == nil {
			t.Errorf("accepted invalid ABI %+v", bad)
		}
	}
}

func TestNativeABI(t *testing.T) {
	if err := Preload(); err != nil {
		t.Fatal(err)
	}
	version := Version()
	t.Logf("native JadeView %s / ABI %#08x", version, ABIVersion)
	if !strings.Contains(version, "3.0.0-beta.4") {
		t.Fatalf("unexpected native version %q", version)
	}
	if _, rc := GetRequestResult(RequestID(1) << 40); rc >= 0 {
		t.Fatalf("unknown request ID returned status %d", rc)
	}
	if _, rc := GetWebViewState(0); rc >= 0 {
		t.Fatalf("invalid WebView ID returned status %d", rc)
	}
}
