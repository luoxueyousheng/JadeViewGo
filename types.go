package jadeview

// 本文件是 Windows(纯 Go)与 Linux(cgo)两套实现共享的公共类型与常量，
// 不得引入 cgo 或平台相关依赖。

// WindowOptions 对应 C 的 WebViewWindowOptions。
type WindowOptions struct {
	// Splash 可选；预览运行库在 Windows 绘制启动画面。
	Splash              *SplashOptions
	Title               string
	Width               int
	Height              int
	Resizable           bool
	NoFrame             bool
	TitleBarStyle       string // default / hidden / BarOverlay
	BackgroundMaterial  string // default / mica / micaAlt / acrylic
	HideTitlebarBack    bool
	HideTitlebarRefresh bool
	HideTitlebarIcon    bool
	Transparent         bool
	BackgroundColor     string // "#RRGGBBAA"
	AlwaysOnTop         bool
	Theme               string // 见 Theme 枚举（enums.go）
	Maximized           bool
	Maximizable         bool
	Minimizable         bool
	X                   int // -1 = 居中
	Y                   int // -1 = 居中
	MinWidth            int
	MinHeight           int
	MaxWidth            int
	MaxHeight           int
	Fullscreen          bool
	Focus               bool
	HideWindow          bool
	UsePageIcon         bool
	ContentProtection   bool
	AutoSaveState       bool
	SkipTaskbar         bool
	NoActivate          bool
}

// DefaultWindowOptions 返回一组常用默认值。
func DefaultWindowOptions() WindowOptions {
	return WindowOptions{
		Title:         "JadeView",
		Width:         1024,
		Height:        768,
		Resizable:     true,
		TitleBarStyle: TitleBarStyle.Default,
		Maximizable:   true,
		Minimizable:   true,
		X:             -1,
		Y:             -1,
		Focus:         true,
	}
}

// WebViewSettings 对应 C 的 WebViewSettings。
type WebViewSettings struct {
	Autoplay                  bool    // 允许媒体自动播放
	BackgroundThrottling      bool    // true=禁用背景限速（false=库默认限速策略）
	AllowRightClick           bool    // 允许页面右键菜单
	UserAgent                 string  // 空=默认 UA
	PreloadJS                 string  // 页面加载前注入的 JS，空=不注入
	AllowFullscreen           bool    // 允许页面全屏
	PostMessageAllowedOrigins string  // postMessage 白名单（单个域名），空=默认
	CORSAllowedOrigins        string  // CORS 来源白名单（逗号分隔），空=默认
	PasswordAutofill          bool    // 账号/密码自动填充
	FormAutofill              bool    // 通用表单自动填充
	Incognito                 bool    // 无痕模式
	DisableClipboard          bool    // 禁用剪贴板读写权限
	ProxyURL                  string  // 代理，如 "http://host:port"/"socks5://host:port"，空=不使用
	Focused                   bool    // 创建后 WebView 自动获取焦点
	ProfileName               string  // Windows WebView2 / Linux 命名 Profile
	DevTools                  int32   // 0=跟随全局，1=启用，2=禁用
	ZoomFactor                float64 // 0=引擎默认
	BackgroundColor           string  // 空=继承窗口背景
}

// DefaultWebViewSettings 返回一组桌面应用常用默认值（与 DefaultWindowOptions 对称）：
// 允许自动播放/右键/全屏/自动填充、创建即聚焦，其余零值（默认 UA、无预载 JS、无代理）。
// 注意：CreateWindow 的 settings 传 nil 表示交给库用内部默认值，二者不保证逐项一致；
// 需要在库默认基础上只改个别项时从本函数出发即可。
func DefaultWebViewSettings() WebViewSettings {
	return WebViewSettings{
		Autoplay:         true,
		AllowRightClick:  true,
		AllowFullscreen:  true,
		PasswordAutofill: true,
		FormAutofill:     true,
		Focused:          true,
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

// EventHandler 是事件处理函数。返回非空字符串会作为响应回传给库（多数事件可返回 ""）。
type EventHandler func(sourceID uint32, data string) string

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

// TitlebarStyleOptions 是 Windows 标题栏样式。Height=-1 保持，0 恢复默认；空颜色恢复主题色。
type TitlebarStyleOptions struct {
	Height                                                  int
	BackgroundColor, ForegroundColor, ButtonBackgroundColor string
	ButtonHoverBackgroundColor, ButtonHoverForegroundColor  string
}

// SplashOptions 描述窗口启动画面。颜色编码为 0xAARRGGBB。
// LogoData 优先于 LogoPath。若两者都为空，则仅显示加载动画。
type SplashOptions struct {
	LogoPath                     string
	LogoData                     []byte
	BackgroundColor, AccentColor uint32
	MinimumMS                    uint32
	Disabled                     bool
}

// PermissionDecision 是浏览器的同步权限决策。
type PermissionDecision int32

const (
	PermissionDeny    PermissionDecision = -1
	PermissionDefault PermissionDecision = 0
	PermissionAllow   PermissionDecision = 1
)

// PermissionHandler 返回 true 允许，false 拒绝；需要浏览器默认行为时使用三态处理器。
type PermissionHandler func(webviewID uint32, data string) bool

// PermissionDecisionHandler 返回三态决策。回调可能在 GUI 线程执行，不能阻塞。
type PermissionDecisionHandler func(webviewID uint32, data string) PermissionDecision

func normalizePermission(d PermissionDecision) PermissionDecision {
	if d == PermissionDefault || d == PermissionAllow {
		return d
	}
	return PermissionDeny
}

// ResourceCapabilities 是指定 WebView 的资源交付能力与计数。
type ResourceCapabilities struct {
	SupportedModes                 uint32 `json:"supported_modes"`
	NativeRequiresPrivateView      uint32 `json:"native_requires_private_view"`
	NativeRequiresPermanent        uint32 `json:"native_requires_permanent"`
	NativeRequiresDedicatedProfile uint32 `json:"native_requires_dedicated_profile"`
	NativeRegistered               uint32 `json:"native_registered"`
	NativePendingCleanup           uint32 `json:"native_pending_cleanup"`
}

type ResourceMode uint32

const (
	ResourceControlled ResourceMode = iota
	ResourceNative
)

// ResourceDeliveryInfo 中的字符串已复制到 Go 内存，不需要调用者释放。
type ResourceDeliveryInfo struct {
	Mode            ResourceMode
	Reason          uint32
	NativeErrorCode int32
	Message         string
}

type JAPKMountCapabilities struct {
	Backend               uint32 `json:"backend"`
	Available             int32  `json:"available"`
	ReadOnly              int32  `json:"read_only"`
	CachesPlaintextToDisk int32  `json:"caches_plaintext_to_disk"`
	OriginalArchivePrefix int32  `json:"original_archive_prefix"`
}
