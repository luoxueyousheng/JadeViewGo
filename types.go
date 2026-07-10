package jadeview

// 本文件是 Windows(纯 Go)与 Linux(cgo)两套实现共享的公共类型与常量，
// 不得引入 cgo 或平台相关依赖。

// WindowOptions 对应 C 的 WebViewWindowOptions。
type WindowOptions struct {
	Title             string
	Width             int
	Height            int
	Resizable         bool
	FrameStyle        string // 见 FrameStyle 枚举（enums.go）
	Transparent       bool
	BackgroundColor   string // "#RRGGBBAA"
	AlwaysOnTop       bool
	Theme             string // 见 Theme 枚举（enums.go）
	Maximized         bool
	Maximizable       bool
	Minimizable       bool
	X                 int // -1 = 居中
	Y                 int // -1 = 居中
	MinWidth          int
	MinHeight         int
	MaxWidth          int
	MaxHeight         int
	Fullscreen        bool
	Focus             bool
	HideWindow        bool
	UsePageIcon       bool
	ContentProtection bool
	AutoSaveState     bool
	SkipTaskbar       bool
	NoActivate        bool
}

// DefaultWindowOptions 返回一组常用默认值。
func DefaultWindowOptions() WindowOptions {
	return WindowOptions{
		Title:       "JadeView",
		Width:       1024,
		Height:      768,
		Resizable:   true,
		FrameStyle:  FrameStyle.Normal,
		Maximizable: true,
		Minimizable: true,
		X:           -1,
		Y:           -1,
		Focus:       true,
	}
}

// WebViewSettings 对应 C 的 WebViewSettings。
type WebViewSettings struct {
	Autoplay               bool   // 允许媒体自动播放
	BackgroundThrottling   bool   // true=禁用背景限速（false=库默认限速策略）
	AllowRightClick        bool   // 允许页面右键菜单
	UserAgent              string // 空=默认 UA
	PreloadJS              string // 页面加载前注入的 JS，空=不注入
	AllowFullscreen        bool   // 允许页面全屏
	PostMessageWhitelist   string // postMessage 白名单（单个域名），空=默认
	CORSWhitelist          string // CORS 来源白名单（逗号分隔），空=默认
	Autofill               bool   // 账号/密码自动填充
	GeneralAutofillEnabled bool   // 通用表单自动填充
	Incognito              bool   // 无痕模式
	DisableClipboard       bool   // 禁用剪贴板读写权限
	ProxyURL               string // 代理，如 "http://host:port"/"socks5://host:port"，空=不使用
	Focused                bool   // 创建后 WebView 自动获取焦点
}

// DefaultWebViewSettings 返回一组桌面应用常用默认值（与 DefaultWindowOptions 对称）：
// 允许自动播放/右键/全屏/自动填充、创建即聚焦，其余零值（默认 UA、无预载 JS、无代理）。
// 注意：CreateWindow 的 settings 传 nil 表示交给库用内部默认值，二者不保证逐项一致；
// 需要在库默认基础上只改个别项时从本函数出发即可。
func DefaultWebViewSettings() WebViewSettings {
	return WebViewSettings{
		Autoplay:               true,
		AllowRightClick:        true,
		AllowFullscreen:        true,
		Autofill:               true,
		GeneralAutofillEnabled: true,
		Focused:                true,
	}
}

// NotificationParams 对应 C 的 NotificationParams。
type NotificationParams struct {
	Summary string // 标题/摘要
	Body    string // 正文
	Icon    string // 图标路径
	Timeout int    // 毫秒，-1=系统默认
	Button1 string
	Button2 string
	Text3   string // 按钮3文本
	Action  string // 默认操作动作
}

// FileDialogParams 对应 C 的 FileDialogParams。
type FileDialogParams struct {
	WindowID    uint32
	Title       string
	DefaultPath string
	ButtonLabel string
	Filters     string // JSON 格式过滤器
	Properties  string // JSON 格式属性
}

// MessageBoxParams 对应 C 的 MessageBoxParams。
type MessageBoxParams struct {
	WindowID  uint32
	Title     string
	Message   string
	Detail    string
	Buttons   string // JSON 数组
	DefaultID int
	CancelID  int
	Type      string // none / info / warning / error
}

// Deprecated: 请使用 TrayItem.Normal 等二级枚举（enums.go）。
const (
	TrayItemNormal  = 0 // Deprecated: 用 TrayItem.Normal
	TrayItemSubmenu = 1 // Deprecated: 用 TrayItem.Submenu
	TrayItemDivider = 2 // Deprecated: 用 TrayItem.Divider
	TrayItemGroup   = 3 // Deprecated: 用 TrayItem.Group
)

// TrayMenuItem 对应 C 的 TrayMenuItemDesc（扁平表，用 ParentKey 指向父项的 Key）。
type TrayMenuItem struct {
	Type      int    // 见 TrayItem 枚举（enums.go）
	Key       string // 全表唯一、非空（分隔线也需唯一 key）
	Label     string
	ParentKey string // 空=根下子项；否则须等于某行的 Key 且该行为 Submenu/Group
	Disabled  bool
	Dangerous bool
}

// Deprecated: 请使用 MenuKind.Normal 等二级枚举（enums.go）。
const (
	MenuKindNormal    = 0 // Deprecated: 用 MenuKind.Normal
	MenuKindSeparator = 1 // Deprecated: 用 MenuKind.Separator
	MenuKindCheckbox  = 2 // Deprecated: 用 MenuKind.Checkbox
	MenuKindRadio     = 3 // Deprecated: 用 MenuKind.Radio
	MenuKindSubmenu   = 4 // Deprecated: 用 MenuKind.Submenu
)

// EventHandler 是事件处理函数。返回非空字符串会作为响应回传给库（多数事件可返回 ""）。
type EventHandler func(windowID uint32, data string) string

// MaxEventHandlers 是可同时注册的事件处理器上限（回调跳板槽位数）。
const MaxEventHandlers = 64

// DialogResultHandler 异步对话框结果回调，result 为结果 JSON（取消时通常为空/null）。
type DialogResultHandler func(result string)

// MaxAsyncDialogs 同时在途的异步对话框上限（回调跳板槽位数）。
const MaxAsyncDialogs = 16

// cBufToString 把以 NUL 结尾的缓冲区转成 Go string。
func cBufToString(buf []byte) string {
	for i, b := range buf {
		if b == 0 {
			return string(buf[:i])
		}
	}
	return string(buf)
}
