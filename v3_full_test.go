//go:build windows || (linux && amd64 && (!jadeview_cgo || cgo))

package jadeview

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
	"time"
)

func TestABIInfoAndEventClone(t *testing.T) {
	info, err := GetABIInfo()
	if err != nil || info.Version != ABIVersion || info.PointerBits != strconv.IntSize {
		t.Fatalf("%+v: %v", info, err)
	}
	event := TypedEvent{Name: "example", Data: []byte("value")}
	cloned := event.Clone()
	cloned.Data[0] = 'V'
	if string(event.Data) != "value" || string(cloned.Data) != "Value" {
		t.Fatal("event clone shares storage")
	}
	if OnAppEvent(EventTypeAll, nil) != 0 {
		t.Fatal("nil handler accepted")
	}
}

func TestV3DesktopIntegration(t *testing.T) {
	if os.Getenv("JADEVIEW_TEST_DESKTOP") != "1" {
		t.Skip("set JADEVIEW_TEST_DESKTOP=1 to open a hidden native WebView")
	}
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	wait := func(poll func() bool, desc string) {
		t.Helper()
		end := time.Now().Add(8 * time.Second)
		for time.Now().Before(end) {
			if poll() {
				return
			}
			time.Sleep(15 * time.Millisecond)
		}
		t.Fatalf("timeout: %s", desc)
	}
	received := make(chan TypedEvent, 8)
	enqueue := func(e TypedEvent) {
		select {
		case received <- e:
		default:
		}
	}
	app := OnAppEvent(EventTypeAll, func(e TypedEvent) {
		if e.Name == "app-ready" {
			enqueue(e)
		}
	})
	if app == 0 {
		t.Fatal("OnAppEvent returned 0")
	}
	dir := t.TempDir()
	if !Init(false, filepath.Join(dir, "jade.log"), dir, "JadeViewGo integration", "com.jadeviewgo.integration", false) {
		t.Fatal("Init rejected")
	}
	defer ExitWait(5000)
	var appEvent TypedEvent
	wait(func() bool {
		select {
		case appEvent = <-received:
			return true
		default:
			return false
		}
	}, "typed app-ready")
	var ready struct{ OK bool }
	if appEvent.SourceID != 0 || json.Unmarshal(appEvent.Data, &ready) != nil || !ready.OK {
		t.Fatalf("unexpected typed app-ready: %+v", appEvent)
	}
	opts := DefaultWindowOptions()
	opts.HideWindow = true
	window := CreateWindow("about:blank", 0, &opts, nil)
	view := DefaultWebView(window)
	if window == 0 || view == 0 {
		t.Fatal("failed to reserve window or WebView")
	}
	defer Close(window)
	wait(func() bool { s, rc := GetWebViewState(view); return rc == 0 && s.State == ObjectReady }, "WebView ready")
	wsub := OnWindowEvent(window, EventTypeCreated, enqueue)
	vsub := OnWebViewEvent(view, EventTypeCreated, enqueue)
	if wsub == 0 || vsub == 0 {
		t.Fatalf("invalid subscriptions %d %d", wsub, vsub)
	}
	gotWindow, gotView := false, false
	wait(func() bool {
		select {
		case e := <-received:
			if e.Type != EventTypeCreated || e.Flags&EventFlagReplay == 0 {
				return false
			}
			gotWindow = gotWindow || (e.WindowID == window && e.SourceID == uint64(window))
			gotView = gotView || (e.WebViewID == view && e.SourceID == uint64(view))
		default:
		}
		return gotWindow && gotView
	}, "native creation replay for both Window and WebView")
	for _, s := range []SubscriptionID{app, wsub, vsub} {
		r := UnsubscribeEvent(s)
		if r == 0 {
			t.Fatal("unsubscribe did not return a request")
		}
		wait(func() bool {
			result, rc := GetRequestResult(r)
			if rc != 0 || result.State == RequestFailed || result.State == RequestCancelled {
				t.Fatalf("unsubscribe %d: %+v (%d)", s, result, rc)
			}
			return result.State == RequestSucceeded
		}, "unsubscribe completion")
		if rc := ReleaseRequest(r); rc != 0 {
			t.Fatalf("ReleaseRequest: %d", rc)
		}
	}
	named, ok := On(EventWindowFocused, func(uint32, string) string { return "" })
	if !ok || !OffImmediate(EventWindowFocused, named) || OffImmediate(EventWindowFocused, named) {
		t.Fatal("named immediate off failed")
	}
	if !SetContextMenu(view, []ContextMenuItem{{MenuItem: MenuItem{ID: "test-item", Label: "Test"}}}) || !ClearContextMenu(view) {
		t.Fatal("setting or clearing view context menu failed")
	}
	file := filepath.Join(dir, "data.txt")
	if err := os.WriteFile(file, []byte("JadeView v3 API"), 0600); err != nil {
		t.Fatal(err)
	}
	url, ok := RegisterResourceEx(file, view, 30, &ResourceOptions{Mode: ResourceControlled})
	if !ok {
		c, m := LastError()
		t.Fatalf("register_resource_ex: %d %s", c, m)
	}
	if mode, rc := GetResourceDeliveryMode(url, view); rc != 0 || mode != ResourceControlled {
		t.Fatalf("mode: %d (%d)", mode, rc)
	}
	if info, rc := GetResourceDeliveryInfo(url, view); rc != 0 || info.Mode != ResourceControlled {
		t.Fatalf("delivery info: %+v (%d)", info, rc)
	}
	if !UnregisterResource(url) {
		t.Fatal("failed to unregister resource")
	}
	if _, ok := RegisterResourceEx(file, view, 30, &ResourceOptions{Mode: ResourceMode(99)}); ok {
		t.Fatal("accepted an invalid resource mode")
	}
	if code, message := LastError(); code != StatusInvalidArgument || message == "" {
		t.Fatalf("thread-local LastError: %d %q", code, message)
	}
	if code := YAMLSetStr("api_test", "key", "123"); code != 1 {
		t.Fatalf("YAMLSetStr: %d", code)
	}
	if s, ok := YAMLGetStr("api_test", "key"); !ok || s != `"123"` {
		t.Fatalf("YAMLGetStr: %q (%t)", s, ok)
	}
	if _, ok := YAMLGetStr("api_test", "absent"); ok {
		t.Fatal("absent key reported valid")
	}

	// Unknown high-bit IDs also exercise x86's low/high stack word convention.
	badArchive := JAPKArchiveID(1 << 44)
	badMount := JAPKMountID(1 << 44)
	if RegisterJAPK(filepath.Join(dir, "missing.japk")) != 0 || RegisterLoadedJAPK() != 0 {
		t.Fatal("registered nonexistent package")
	}
	if MountJAPK(badArchive, "") != 0 || UnregisterJAPK(badArchive) == 0 || UnmountJAPK(badMount) == 0 {
		t.Fatal("accepted invalid archive or mount ID")
	}
	if _, ok := GetJAPKMountRoot(badMount); ok {
		t.Fatal("returned path for invalid mount")
	}
	if _, rc := GetJAPKMountState(badMount); rc == 0 {
		t.Fatal("returned state for invalid mount")
	}
	if caps, rc := GetJAPKMountCapabilities(); rc != 0 {
		t.Fatalf("get mount capabilities: %d", rc)
	} else {
		t.Logf("mount capabilities: %+v", caps)
	}
	t.Log("typed notifications + replay + unsubscribe, resource mode, context menu, YAML, and JAPK ABI passed")
}
