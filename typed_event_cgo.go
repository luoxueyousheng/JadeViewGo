//go:build linux && amd64 && jadeview_cgo && cgo

package jadeview

/*
#include "JadeView.h"
extern void goTypedDispatch(struct JadeEvent*, uintptr_t);
static void typed_dispatch(const struct JadeEvent* e, void* user) {
 goTypedDispatch((struct JadeEvent*)e, (uintptr_t)user);
}
static uint64_t typed_subscribe(uint8_t scope,uint32_t id,uint32_t event_type,uintptr_t token) {
 switch (scope) {
  case 0: return on_app_event(event_type,typed_dispatch,(void*)token);
  case 1: return on_window_event(id,event_type,typed_dispatch,(void*)token);
  case 2: return on_webview_event(id,event_type,typed_dispatch,(void*)token);
  default: return 0;
 }
}
*/
import "C"

import "unsafe"

//export goTypedDispatch
func goTypedDispatch(e *C.JadeEvent, token C.uintptr_t) {
	h := typedHandler(uintptr(token))
	if h == nil || e == nil {
		return
	}
	h(TypedEvent{EventType(e.event_type), uint32(e.source_kind), uint64(e.source_id), uint32(e.window_id), uint32(e.webview_id),
		uint64(e.sequence), RequestID(e.request_id), uint32(e.flags), uint32(e.state), int32(e.error_code), C.GoString(e.name),
		eventDataCopy(unsafe.Pointer(e.data), uintptr(e.data_size))})
}

func registerTypedEvent(scope subscriptionScope, objectID uint32, eventType EventType, token uintptr) SubscriptionID {
	return SubscriptionID(C.typed_subscribe(C.uint8_t(scope), C.uint32_t(objectID), C.uint32_t(eventType), C.uintptr_t(token)))
}
func unregisterTypedEvent(id SubscriptionID) RequestID {
	return RequestID(C.unsubscribe_event(C.uint64_t(id)))
}
