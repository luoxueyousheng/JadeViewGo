package jadeview

// RequestID 是 v3 原生异步请求 ID。0 表示请求未被接受。
type RequestID uint64

type RequestState uint32

const (
	RequestAccepted RequestState = iota + 1
	RequestWaiting
	RequestSucceeded
	RequestFailed
	RequestCancelled
)

// RequestResult 是可独立于原生缓冲区保留的请求结果。
type RequestResult struct {
	State     RequestState
	ErrorCode int32
	ObjectID  uint32
	Message   string
}

// WebViewBounds 使用 CSS 表达式描述位置和尺寸，例如 "50%-20px"。
type WebViewBounds struct{ X, Y, Width, Height string }

// LogicalRect 是以逻辑像素表示的 WebView 矩形。
type LogicalRect struct{ X, Y, Width, Height float64 }

// PhysicalRect 是以实际像素表示的 WebView 矩形。
type PhysicalRect struct {
	X, Y          int32
	Width, Height uint32
}

// WebViewLayout 包含已应用的布局表达式及计算后的矩形。
type WebViewLayout struct {
	Bounds   WebViewBounds
	Logical  LogicalRect
	Physical PhysicalRect
}

// ObjectState 是 Window 或 WebView 的原生状态快照。
type ObjectState struct {
	Kind, State, ParentWindowID, DefaultWebViewID uint32
	ErrorCode                                     int32
}
