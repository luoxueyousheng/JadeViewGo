//go:build windows || (linux && amd64 && !jadeview_cgo)

package jadeview

import (
	"sync"
	"unsafe"
)

type cTypedEvent struct {
	eventType, sourceKind uint32
	sourceID              uint64
	windowID, viewID      uint32
	sequence, requestID   uint64
	flags, state          uint32
	errorCode             int32
	name, data            *byte
	size                  uintptr
}

var (
	procOnAppEvent       = np("on_app_event")
	procOnWindowEvent    = np("on_window_event")
	procOnWebViewEvent   = np("on_webview_event")
	procUnsubscribeEvent = np("unsubscribe_event")
	typedTrampolineOnce  sync.Once
	typedTrampoline      uintptr
)

func typedCallback(raw, token uintptr) uintptr {
	h := typedHandler(token)
	if h == nil || raw == 0 {
		return 0
	}
	e := (*cTypedEvent)(foreignPtr(raw))
	h(TypedEvent{EventType(e.eventType), e.sourceKind, e.sourceID, e.windowID, e.viewID, e.sequence, RequestID(e.requestID),
		e.flags, e.state, e.errorCode, goString(uintptr(unsafe.Pointer(e.name))), eventDataCopy(unsafe.Pointer(e.data), e.size)})
	return 0
}

func registerTypedEvent(scope subscriptionScope, objectID uint32, eventType EventType, token uintptr) SubscriptionID {
	typedTrampolineOnce.Do(func() { typedTrampoline = newNativeCallback(typedCallback) })
	var lo, hi uintptr
	switch scope {
	case subscriptionApp:
		lo, hi, _ = procOnAppEvent.Call(uintptr(eventType), typedTrampoline, token)
	case subscriptionWindow:
		lo, hi, _ = procOnWindowEvent.Call(uintptr(objectID), uintptr(eventType), typedTrampoline, token)
	case subscriptionWebView:
		lo, hi, _ = procOnWebViewEvent.Call(uintptr(objectID), uintptr(eventType), typedTrampoline, token)
	}
	return SubscriptionID(requestID(lo, hi))
}

func unregisterTypedEvent(id SubscriptionID) RequestID {
	lo, hi, _ := callU64First(procUnsubscribeEvent, uint64(id))
	return requestID(lo, hi)
}
