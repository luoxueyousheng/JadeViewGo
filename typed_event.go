//go:build windows || (linux && amd64 && (!jadeview_cgo || cgo))

package jadeview

import (
	"slices"
	"sync"
	"unsafe"
)

// EventType 是结构化事件类型；EventTypeAll 允许接收所有类型。
type EventType uint32

const (
	EventTypeAll EventType = iota
	EventTypeCreated
	EventTypeCreating
	EventTypeCreateFailed
	EventTypeDestroyed
	EventTypeLoadStarted
	EventTypeLoadFinished
	EventTypeCrashed
	EventTypeOperationCompleted
	EventTypeOperationFailed
	EventTypeOther EventType = 100
)

// EventFlagReplay 标识订阅时的当前状态回放。
const EventFlagReplay uint32 = 1

// TypedEvent 是结构化事件的 Go 快照。Name 和 Data 在回调返回后仍然有效；
// 它们不引用原生内存，无需调用原生 event_free。
type TypedEvent struct {
	Type                EventType
	SourceKind          uint32
	SourceID            uint64
	WindowID, WebViewID uint32
	Sequence            uint64
	RequestID           RequestID
	Flags, State        uint32
	ErrorCode           int32
	Name                string
	Data                []byte
}

// Clone 返回独立快照。需要修改 Data 而不影响现有快照时使用。
func (e TypedEvent) Clone() TypedEvent { e.Data = slices.Clone(e.Data); return e }

type TypedEventHandler func(event TypedEvent)

// SubscriptionID 是结构化订阅 ID，不同于 On 的命名回调 ID。
type SubscriptionID uint64

type subscriptionScope uint8

const (
	subscriptionApp subscriptionScope = iota
	subscriptionWindow
	subscriptionWebView
)

// 单个进程级原生回调跳板。传给原生的是永不复用的标量 token，不是 Go 指针。
// 注销后立即删除 token；已取到 Go handler 的回调仍可安全完成。
var typedSubscriptions struct {
	sync.Mutex
	next     uintptr
	handlers map[uintptr]TypedEventHandler
	byID     map[SubscriptionID]uintptr
}

func subscribeTyped(scope subscriptionScope, objectID uint32, eventType EventType, h TypedEventHandler) SubscriptionID {
	if h == nil {
		return 0
	}
	typedSubscriptions.Lock()
	if typedSubscriptions.next == ^uintptr(0) {
		typedSubscriptions.Unlock()
		return 0
	}
	if typedSubscriptions.handlers == nil {
		typedSubscriptions.handlers = make(map[uintptr]TypedEventHandler)
		typedSubscriptions.byID = make(map[SubscriptionID]uintptr)
	}
	typedSubscriptions.next++
	token := typedSubscriptions.next
	typedSubscriptions.handlers[token] = h
	typedSubscriptions.Unlock()

	id := registerTypedEvent(scope, objectID, eventType, token)
	typedSubscriptions.Lock()
	if id == 0 {
		delete(typedSubscriptions.handlers, token)
	} else {
		typedSubscriptions.byID[id] = token
	}
	typedSubscriptions.Unlock()
	return id
}

// OnAppEvent 订阅应用事件；应在 Init 前订阅 app-ready。
func OnAppEvent(eventType EventType, h TypedEventHandler) SubscriptionID {
	return subscribeTyped(subscriptionApp, 0, eventType, h)
}

// OnWindowEvent 订阅指定窗口的结构化事件。
func OnWindowEvent(windowID uint32, eventType EventType, h TypedEventHandler) SubscriptionID {
	return subscribeTyped(subscriptionWindow, windowID, eventType, h)
}

// OnWebViewEvent 订阅指定 WebView 的结构化事件。
func OnWebViewEvent(viewID uint32, eventType EventType, h TypedEventHandler) SubscriptionID {
	return subscribeTyped(subscriptionWebView, viewID, eventType, h)
}

// UnsubscribeEvent 停止订阅并返回异步请求 ID。用 GetRequestResult 查询，完成后
// 用 ReleaseRequest 释放；不要在事件回调内等待自己的注销请求。
func UnsubscribeEvent(id SubscriptionID) RequestID {
	if id == 0 {
		return 0
	}
	request := unregisterTypedEvent(id)
	if request != 0 {
		typedSubscriptions.Lock()
		if token, ok := typedSubscriptions.byID[id]; ok {
			delete(typedSubscriptions.handlers, token)
			delete(typedSubscriptions.byID, id)
		}
		typedSubscriptions.Unlock()
	}
	return request
}

func typedHandler(token uintptr) TypedEventHandler {
	typedSubscriptions.Lock()
	h := typedSubscriptions.handlers[token]
	typedSubscriptions.Unlock()
	return h
}

func eventDataCopy(data unsafe.Pointer, size uintptr) []byte {
	if data == nil || size == 0 || size > ^uintptr(0)>>1 {
		return nil
	}
	return slices.Clone(unsafe.Slice((*byte)(data), int(size)))
}
