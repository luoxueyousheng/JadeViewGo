//go:build windows || (linux && amd64 && !jadeview_cgo)

package jadeview

import (
	"runtime"
	"unsafe"
)

type cWindowOptions struct {
	structSize         uint32
	abiVersion         uint32
	title              *byte
	width              int32
	height             int32
	resizable          int32
	noframe            int32
	titleBarStyle      *byte
	transparent        int32
	backgroundColor    *byte
	alwaysOnTop        int32
	theme              *byte
	maximized          int32
	maximizable        int32
	minimizable        int32
	x                  int32
	y                  int32
	minWidth           int32
	minHeight          int32
	maxWidth           int32
	maxHeight          int32
	fullscreen         int32
	focus              int32
	hideWindow         int32
	usePageIcon        int32
	skipTaskbar        int32
	noActivate         int32
	backgroundMaterial *byte
	titlebarBack       int32
	titlebarRefresh    int32
	titlebarIcon       int32
}
type cWebViewSettingsFields struct {
	structSize                uint32
	abiVersion                uint32
	autoplay                  int32
	backgroundThrottling      int32
	allowRightClick           int32
	userAgent                 *byte
	preloadJS                 *byte
	allowFullscreen           int32
	postmessageAllowedOrigins *byte
	corsAllowedOrigins        *byte
	passwordAutofill          int32
	formAutofill              int32
	incognito                 int32
	disableClipboard          int32
	proxyURL                  *byte
	focused                   int32
	profileName               *byte
	devtools                  int32
	zoomFactor                float64
	backgroundColor           *byte
}

func bi32(b bool) int32 {
	if b {
		return 1
	}
	return 0
}

var (
	procWindowOptionsDefault   = np("window_options_default")
	procWebViewSettingsDefault = np("webview_settings_default")
	procCreateEmptyWindow      = np("create_window")
)

func (o *WindowOptions) toC(pool *cstrs) *cWindowOptions {
	if o == nil {
		return nil
	}
	c := new(cWindowOptions)
	r, _, _ := procWindowOptionsDefault.Call(uintptr(unsafe.Pointer(c)), unsafe.Sizeof(*c))
	if i32(r) != 0 {
		panic("JadeView: WindowOptionsDefault failed")
	}
	c.title = pool.p(o.Title)
	c.width = int32(o.Width)
	c.height = int32(o.Height)
	c.resizable = bi32(o.Resizable)
	c.noframe = bi32(o.NoFrame)
	c.titleBarStyle = pool.p(o.TitleBarStyle)
	c.transparent = bi32(o.Transparent)
	c.backgroundColor = pool.p(o.BackgroundColor)
	c.alwaysOnTop = bi32(o.AlwaysOnTop)
	c.theme = pool.p(o.Theme)
	c.maximized = bi32(o.Maximized)
	c.maximizable = bi32(o.Maximizable)
	c.minimizable = bi32(o.Minimizable)
	c.x = int32(o.X)
	c.y = int32(o.Y)
	c.minWidth = int32(o.MinWidth)
	c.minHeight = int32(o.MinHeight)
	c.maxWidth = int32(o.MaxWidth)
	c.maxHeight = int32(o.MaxHeight)
	c.fullscreen = bi32(o.Fullscreen)
	c.focus = bi32(o.Focus)
	c.hideWindow = bi32(o.HideWindow)
	c.usePageIcon = bi32(o.UsePageIcon)
	c.skipTaskbar = bi32(o.SkipTaskbar)
	c.noActivate = bi32(o.NoActivate)
	c.backgroundMaterial = pool.p(o.BackgroundMaterial)
	c.titlebarBack = bi32(!o.HideTitlebarBack)
	c.titlebarRefresh = bi32(!o.HideTitlebarRefresh)
	c.titlebarIcon = bi32(!o.HideTitlebarIcon)
	return c
}
func (s *WebViewSettings) toC(pool *cstrs) *cWebViewSettings {
	if s == nil {
		return nil
	}
	c := new(cWebViewSettings)
	r, _, _ := procWebViewSettingsDefault.Call(uintptr(unsafe.Pointer(c)), unsafe.Sizeof(*c))
	if i32(r) != 0 {
		panic("JadeView: WebViewSettingsDefault failed")
	}
	c.autoplay = bi32(s.Autoplay)
	c.backgroundThrottling = bi32(s.BackgroundThrottling)
	c.allowRightClick = bi32(s.AllowRightClick)
	c.userAgent = pool.p(s.UserAgent)
	c.preloadJS = pool.p(s.PreloadJS)
	c.allowFullscreen = bi32(s.AllowFullscreen)
	c.postmessageAllowedOrigins = pool.p(s.PostMessageAllowedOrigins)
	c.corsAllowedOrigins = pool.p(s.CORSAllowedOrigins)
	c.passwordAutofill = bi32(s.PasswordAutofill)
	c.formAutofill = bi32(s.FormAutofill)
	c.incognito = bi32(s.Incognito)
	c.disableClipboard = bi32(s.DisableClipboard)
	c.proxyURL = pool.p(s.ProxyURL)
	c.focused = bi32(s.Focused)
	c.profileName = pool.p(s.ProfileName)
	c.devtools = int32(s.DevTools)
	c.zoomFactor = s.ZoomFactor
	c.backgroundColor = pool.p(s.BackgroundColor)
	return c
}

// CreateWindow 创建 Window 和默认 WebView，返回 Window ID。nil 配置使用原生默认值。
func CreateWindow(url string, parentID uint32, opts *WindowOptions, settings *WebViewSettings) uint32 {
	pool := &cstrs{}
	co, cs, f := opts.toC(pool), settings.toC(pool), opts.features(pool)
	r, _, _ := procCreateWebviewWindow.Call(uintptr(unsafe.Pointer(pool.p(url))), uintptr(parentID), uintptr(unsafe.Pointer(co)), uintptr(unsafe.Pointer(cs)), uintptr(unsafe.Pointer(f)))
	runtime.KeepAlive(pool)
	runtime.KeepAlive(opts)
	return uint32(r)
}

// CreateEmptyWindow 创建不带 WebView 的原生 Window。
func CreateEmptyWindow(parentID uint32, opts *WindowOptions) uint32 {
	pool := &cstrs{}
	co, f := opts.toC(pool), opts.features(pool)
	r, _, _ := procCreateEmptyWindow.Call(uintptr(parentID), uintptr(unsafe.Pointer(co)), uintptr(unsafe.Pointer(f)))
	runtime.KeepAlive(pool)
	runtime.KeepAlive(opts)
	return uint32(r)
}

func (o *WindowOptions) features(pool *cstrs) *cWindowFeatures {
	if o == nil || (!o.ContentProtection && !o.AutoSaveState && o.Splash == nil) {
		return nil
	}
	f := &cWindowFeatures{size: uint32(unsafe.Sizeof(cWindowFeatures{})), abi: ABIVersion, contentProtection: bi32(o.ContentProtection), autoSaveState: bi32(o.AutoSaveState)}
	if s := o.Splash; s != nil {
		f.splash = &cSplashOptions{size: uint32(unsafe.Sizeof(cSplashOptions{})), abi: ABIVersion, flags: uint32(b2u(s.Disabled)),
			path: pool.p(s.LogoPath), data: unsafe.SliceData(s.LogoData), dataSize: uintptr(len(s.LogoData)), background: s.BackgroundColor, accent: s.AccentColor, minimumMS: s.MinimumMS}
	}
	return f
}

type cWindowFeatures struct {
	size, abi                        uint32
	contentProtection, autoSaveState int32
	splash                           *cSplashOptions
}

type cSplashOptions struct {
	size, abi, flags              uint32
	path, data                    *byte
	dataSize                      uintptr
	background, accent, minimumMS uint32
}

// --- 常用窗口操作（薄封装） ---

func Navigate(webviewID uint32, url, headersJSON string) bool {
	pool := &cstrs{}
	r, _, _ := procNavigateToURL.Call(uintptr(webviewID),
		uintptr(unsafe.Pointer(pool.p(url))), uintptr(unsafe.Pointer(pool.p(headersJSON))))
	runtime.KeepAlive(pool)
	return i32(r) == 1
}

func Reload(webviewID uint32) bool {
	r, _, _ := procReloadWebviewWindow.Call(uintptr(webviewID))
	return i32(r) == 1
}

// GoBack 导航历史后退（原生后退）。
func GoBack(webviewID uint32) bool {
	r, _, _ := procWebviewGoBack.Call(uintptr(webviewID))
	return i32(r) == 1
}

// GoForward 导航历史前进（原生前进）。
func GoForward(webviewID uint32) bool {
	r, _, _ := procWebviewGoForward.Call(uintptr(webviewID))
	return i32(r) == 1
}

// CanGoBack 查询当前页面能否后退。
func CanGoBack(webviewID uint32) bool {
	r, _, _ := procWebviewCanGoBack.Call(uintptr(webviewID))
	return i32(r) == 1
}

// CanGoForward 查询当前页面能否前进。
func CanGoForward(webviewID uint32) bool {
	r, _, _ := procWebviewCanGoForward.Call(uintptr(webviewID))
	return i32(r) == 1
}

func SetTitle(windowID uint32, title string) bool {
	pool := &cstrs{}
	r, _, _ := procSetWindowTitle.Call(uintptr(windowID), uintptr(unsafe.Pointer(pool.p(title))))
	runtime.KeepAlive(pool)
	return i32(r) == 1
}

func SetSize(windowID uint32, width, height int) bool {
	r, _, _ := procSetWindowSize.Call(uintptr(windowID), uintptr(uint32(int32(width))), uintptr(uint32(int32(height))))
	return i32(r) == 1
}

func SetPosition(windowID uint32, x, y int) bool {
	r, _, _ := procSetWindowPosition.Call(uintptr(windowID), uintptr(uint32(int32(x))), uintptr(uint32(int32(y))))
	return i32(r) == 1
}

func SetVisible(windowID uint32, visible bool) bool {
	r, _, _ := procSetWindowVisible.Call(uintptr(windowID), b2u(visible))
	return i32(r) == 1
}

func SetFocus(windowID uint32) bool {
	r, _, _ := procSetWindowFocus.Call(uintptr(windowID))
	return i32(r) == 1
}

func SetAlwaysOnTop(windowID uint32, on bool) bool {
	r, _, _ := procSetWindowAlwaysOnTop.Call(uintptr(windowID), b2u(on))
	return i32(r) == 1
}

func Close(windowID uint32) bool {
	r, _, _ := procCloseWindow.Call(uintptr(windowID))
	return i32(r) == 1
}

func Minimize(windowID uint32) bool {
	r, _, _ := procMinimizeWindow.Call(uintptr(windowID))
	return i32(r) == 1
}

func ToggleMaximize(windowID uint32) bool {
	r, _, _ := procToggleMaximizeWindow.Call(uintptr(windowID))
	return i32(r) == 1
}

func IsMaximized(windowID uint32) bool {
	r, _, _ := procIsWindowMaximized.Call(uintptr(windowID))
	return i32(r) == 1
}

// WindowCount 返回当前窗口数量。
func WindowCount() uint32 {
	r, _, _ := procGetWindowCount.Call()
	return uint32(r)
}

// --- 生命周期 ---

// Init 初始化 JadeView。
//
//   - enableDevmod : 是否启用开发模式
//   - logPath      : 日志路径（可空字符串）
//   - dataDir      : 数据目录（可空字符串）
//   - appName      : 应用名，必填、非纯空白
//   - appSignature : 应用签名，trim 后至少 6 个 Unicode 字符；建议反域名格式
//     （如 com.example.myapp）——它会成为 JAPK 模式下 JADE:// URL 的主机名
//   - singleInstance: 是否启用单实例
//
// 返回 true 表示成功。注意：宿主仍应在 jade_on("app-ready", ...) 回调里
// 解析 JSON 的 ok 字段；回调来源 ID 为 0。
func Init(enableDevmod bool, logPath, dataDir, appName, appSignature string, singleInstance bool) bool {
	pool := &cstrs{}
	r, _, _ := procJadeViewInit.Call(
		b2u(enableDevmod),
		uintptr(unsafe.Pointer(pool.pAlways(logPath))),
		uintptr(unsafe.Pointer(pool.pAlways(dataDir))),
		uintptr(unsafe.Pointer(pool.pAlways(appName))),
		uintptr(unsafe.Pointer(pool.pAlways(appSignature))),
		b2u(singleInstance),
	)
	runtime.KeepAlive(pool)
	return i32(r) == 1
}

// Version 返回 JadeView 版本字符串。
func Version() string {
	s, ok := bufCallSize(256, func(buf unsafe.Pointer, size uintptr) int32 {
		r, _, _ := procJadeviewVersion.Call(uintptr(buf), size)
		return i32(r)
	})
	if !ok {
		return ""
	}
	return s
}

// RunMessageLoop 运行消息循环（阻塞，直到窗口全部关闭/退出）。
func RunMessageLoop() {
	procRunMessageLoop.Call()
}

// Exit 清理所有窗口并结束消息循环。
func Exit() {
	procJadeviewExit.Call()
}

// ExitWait 发起退出并等待原生运行时收尾，返回库状态码。
func ExitWait(timeoutMs uint32) int32 {
	r, _, _ := procJadeviewExitWait.Call(uintptr(timeoutMs))
	return i32(r)
}

// --- 窗口状态查询 ---

func IsMinimized(windowID uint32) bool {
	r, _, _ := procIsWindowMinimized.Call(uintptr(windowID))
	return i32(r) == 1
}

func IsVisible(windowID uint32) bool {
	r, _, _ := procIsWindowVisible.Call(uintptr(windowID))
	return i32(r) == 1
}

func IsFocused(windowID uint32) bool {
	r, _, _ := procIsWindowFocused.Call(uintptr(windowID))
	return i32(r) == 1
}

func IsFullscreen(windowID uint32) bool {
	r, _, _ := procIsWindowFullscreen.Call(uintptr(windowID))
	return i32(r) == 1
}

// GetWindowBounds 返回窗口边界 JSON。
func GetWindowBounds(windowID uint32) (string, bool) {
	return bufCallInt(256, func(buf unsafe.Pointer, size int32) int32 {
		r, _, _ := procGetWindowBounds.Call(uintptr(windowID), uintptr(buf), uintptr(uint32(size)))
		return i32(r)
	})
}

// GetWebViewURL 返回当前 WebView 地址。
func GetWebViewURL(webviewID uint32) (string, bool) {
	return bufCallInt(2048, func(buf unsafe.Pointer, size int32) int32 {
		r, _, _ := procGetWebviewURL.Call(uintptr(webviewID), uintptr(buf), uintptr(uint32(size)))
		return i32(r)
	})
}

// GetWindowHWND 返回窗口原生句柄。
func GetWindowHWND(windowID uint32) uintptr {
	r, _, _ := procGetWindowHwnd.Call(uintptr(windowID))
	return r
}

// GetWindowID 根据 HWND 反查窗口 ID，0 表示未找到。
func GetWindowID(hwnd uintptr) uint32 {
	r, _, _ := procGetWindowID.Call(hwnd)
	return uint32(r)
}

// --- 尺寸 / 可调整 ---

func SetMinSize(windowID uint32, width, height int) bool {
	r, _, _ := procSetWindowMinSize.Call(uintptr(windowID), uintptr(uint32(int32(width))), uintptr(uint32(int32(height))))
	return i32(r) == 1
}

func SetMaxSize(windowID uint32, width, height int) bool {
	r, _, _ := procSetWindowMaxSize.Call(uintptr(windowID), uintptr(uint32(int32(width))), uintptr(uint32(int32(height))))
	return i32(r) == 1
}

func SetResizable(windowID uint32, resizable bool) bool {
	r, _, _ := procSetWindowResizable.Call(uintptr(windowID), b2u(resizable))
	return i32(r) == 1
}

func SetFullscreen(windowID uint32, fullscreen bool) bool {
	r, _, _ := procSetWindowFullscreen.Call(uintptr(windowID), b2u(fullscreen))
	return i32(r) == 1
}

// SetIgnoreCursorEvents 设置窗口是否穿透鼠标事件。
func SetIgnoreCursorEvents(windowID uint32, ignore bool) bool {
	r, _, _ := procSetIgnoreCursorEvents.Call(uintptr(windowID), b2u(ignore))
	return i32(r) == 1
}

// --- 外观 / 层级 ---

func SetSkipTaskbar(windowID uint32, skip bool) bool {
	r, _, _ := procSetWindowSkipTaskbar.Call(uintptr(windowID), b2u(skip))
	return i32(r) == 1
}

func SetNoActivate(windowID uint32, noActivate bool) bool {
	r, _, _ := procSetWindowNoActivate.Call(uintptr(windowID), b2u(noActivate))
	return i32(r) == 1
}

// SetLevel 设置窗口层级：topmost | normal | bottom | desktop。
func SetLevel(windowID uint32, level string) bool {
	pool := &cstrs{}
	r, _, _ := procSetWindowLevel.Call(uintptr(windowID), uintptr(unsafe.Pointer(pool.p(level))))
	runtime.KeepAlive(pool)
	return i32(r) == 1
}

func SetContentProtection(windowID uint32, on bool) bool {
	r, _, _ := procSetContentProtection.Call(uintptr(windowID), b2u(on))
	return i32(r) == 1
}

// SetZoom 设置 WebView 缩放级别。
// double 参数无法经整数寄存器传递（x64/arm64 走浮点寄存器），见 fltcall_windows_*.go。
func SetZoom(webviewID uint32, level float64) bool {
	return callF64I32(procSetWebviewZoom, webviewID, level) == 1
}

func SetWindowFrame(windowID uint32, titleBarStyle string) bool {
	pool := &cstrs{}
	r, _, _ := procSetWindowFrame.Call(uintptr(windowID), uintptr(unsafe.Pointer(pool.p(titleBarStyle))))
	runtime.KeepAlive(pool)
	return i32(r) == 1
}

func SetTheme(windowID uint32, theme string) bool {
	pool := &cstrs{}
	r, _, _ := procSetWindowTheme.Call(uintptr(windowID), uintptr(unsafe.Pointer(pool.p(theme))))
	runtime.KeepAlive(pool)
	return i32(r) == 1
}

// GetTheme 返回窗口主题（int32，含义见库文档）。
func GetTheme(windowID uint32) int32 {
	r, _, _ := procGetWindowTheme.Call(uintptr(windowID))
	return i32(r)
}

func SetBackdrop(windowID uint32, backdropType string) bool {
	pool := &cstrs{}
	r, _, _ := procSetWindowBackdrop.Call(uintptr(windowID), uintptr(unsafe.Pointer(pool.p(backdropType))))
	runtime.KeepAlive(pool)
	return i32(r) == 1
}

// SetBackgroundColor 设置窗口纯色底：#RRGGBBAA。
func SetBackgroundColor(windowID uint32, colorHex string) bool {
	pool := &cstrs{}
	r, _, _ := procSetWindowBgColor.Call(uintptr(windowID), uintptr(unsafe.Pointer(pool.p(colorHex))))
	runtime.KeepAlive(pool)
	return i32(r) == 1
}

func SetEnabled(windowID uint32, enabled bool) bool {
	r, _, _ := procSetWindowEnabled.Call(uintptr(windowID), b2u(enabled))
	return i32(r) == 1
}

func RequestRedraw(windowID uint32) bool {
	r, _, _ := procRequestRedraw.Call(uintptr(windowID))
	return i32(r) == 1
}

var procSetTitlebarStyle = np("set_titlebar_style")
var procSetTitlebarComponents = np("set_titlebar_components")

// SetTitlebarStyle 更新 Windows 标题栏颜色和高度。
func SetTitlebarStyle(windowID uint32, style TitlebarStyleOptions) bool {
	pool := &cstrs{}
	r, _, _ := procSetTitlebarStyle.Call(uintptr(windowID), uintptr(uint32(int32(style.Height))),
		uintptr(unsafe.Pointer(pool.p(style.BackgroundColor))), uintptr(unsafe.Pointer(pool.p(style.ForegroundColor))),
		uintptr(unsafe.Pointer(pool.p(style.ButtonBackgroundColor))), uintptr(unsafe.Pointer(pool.p(style.ButtonHoverBackgroundColor))),
		uintptr(unsafe.Pointer(pool.p(style.ButtonHoverForegroundColor))))
	runtime.KeepAlive(pool)
	return i32(r) == 1
}

// SetTitlebarComponents 一次更新五个 Windows 标题栏开关。
func SetTitlebarComponents(windowID uint32, back, refresh, icon, minimize, maximize bool) bool {
	r, _, _ := procSetTitlebarComponents.Call(uintptr(windowID), b2u(back), b2u(refresh), b2u(icon), b2u(minimize), b2u(maximize))
	return i32(r) == 1
}

// --- IPC / DevTools / 其它 ---

// SendIPCMessage 向前端发送 IPC 消息。
func SendIPCMessage(webviewID uint32, messageType, messageContent string) bool {
	pool := &cstrs{}
	r, _, _ := procSendIPCMessage.Call(uintptr(webviewID),
		uintptr(unsafe.Pointer(pool.p(messageType))), uintptr(unsafe.Pointer(pool.p(messageContent))))
	runtime.KeepAlive(pool)
	return i32(r) == 1
}

func OpenDevtools(webviewID uint32) bool {
	r, _, _ := procOpenDevtools.Call(uintptr(webviewID))
	return i32(r) == 1
}

func CloseDevtools(webviewID uint32) bool {
	r, _, _ := procCloseDevtools.Call(uintptr(webviewID))
	return i32(r) == 1
}

func IsDevtoolsOpen(webviewID uint32) bool {
	r, _, _ := procIsDevtoolsOpen.Call(uintptr(webviewID))
	return i32(r) == 1
}

func ClearBrowsingData(webviewID uint32) bool {
	r, _, _ := procClearBrowsingData.Call(uintptr(webviewID))
	return i32(r) == 1
}

// SetWindowProgress 设置任务栏进度。state 含义见库文档。
func SetWindowProgress(windowID uint32, progress, state int) bool {
	r, _, _ := procSetWindowProgress.Call(uintptr(windowID), uintptr(uint32(int32(progress))), uintptr(uint32(int32(state))))
	return i32(r) == 1
}

// FlashWindow 闪烁窗口 count 次。
func FlashWindow(windowID uint32, count uint32) bool {
	r, _, _ := procFlashWindow.Call(uintptr(windowID), uintptr(count))
	return i32(r) == 1
}

func ShowAboutDialog(windowID uint32) bool {
	r, _, _ := procShowAboutDialog.Call(uintptr(windowID))
	return i32(r) == 1
}

// CreateBorderlessWindow 创建无边框窗口，返回 Window ID。
func CreateBorderlessWindow(url string, settings *WebViewSettings) uint32 {
	opts := DefaultWindowOptions()
	opts.NoFrame = true
	return CreateWindow(url, 0, &opts, settings)
}
