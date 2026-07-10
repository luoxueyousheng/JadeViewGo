//go:build windows

package jadeview

import (
	"runtime"
	"unsafe"
)

// cWindowOptions 是 C WebViewWindowOptions 的逐字段镜像（布局已双端比对验证）。
type cWindowOptions struct {
	title             *byte
	width             int32
	height            int32
	resizable         int32
	frameStyle        *byte
	transparent       int32
	backgroundColor   *byte
	alwaysOnTop       int32
	theme             *byte
	maximized         int32
	maximizable       int32
	minimizable       int32
	x                 int32
	y                 int32
	minWidth          int32
	minHeight         int32
	maxWidth          int32
	maxHeight         int32
	fullscreen        int32
	focus             int32
	hideWindow        int32
	usePageIcon       int32
	contentProtection int32
	autoSaveState     int32
	skipTaskbar       int32
	noActivate        int32
}

// cWebViewSettings 是 C WebViewSettings 的逐字段镜像。
type cWebViewSettings struct {
	autoplay               int32
	backgroundThrottling   int32
	allowRightClick        int32
	ua                     *byte
	preloadJS              *byte
	allowFullscreen        int32
	postmessageWhitelist   *byte
	corsWhitelist          *byte
	autofill               int32
	generalAutofillEnabled int32
	incognito              int32
	disableClipboard       int32
	proxyURL               *byte
	focused                int32
}

func bi32(b bool) int32 {
	if b {
		return 1
	}
	return 0
}

func (o *WindowOptions) toC(pool *cstrs) cWindowOptions {
	return cWindowOptions{
		title:             pool.p(o.Title),
		width:             int32(o.Width),
		height:            int32(o.Height),
		resizable:         bi32(o.Resizable),
		frameStyle:        pool.p(o.FrameStyle),
		transparent:       bi32(o.Transparent),
		backgroundColor:   pool.p(o.BackgroundColor),
		alwaysOnTop:       bi32(o.AlwaysOnTop),
		theme:             pool.p(o.Theme),
		maximized:         bi32(o.Maximized),
		maximizable:       bi32(o.Maximizable),
		minimizable:       bi32(o.Minimizable),
		x:                 int32(o.X),
		y:                 int32(o.Y),
		minWidth:          int32(o.MinWidth),
		minHeight:         int32(o.MinHeight),
		maxWidth:          int32(o.MaxWidth),
		maxHeight:         int32(o.MaxHeight),
		fullscreen:        bi32(o.Fullscreen),
		focus:             bi32(o.Focus),
		hideWindow:        bi32(o.HideWindow),
		usePageIcon:       bi32(o.UsePageIcon),
		contentProtection: bi32(o.ContentProtection),
		autoSaveState:     bi32(o.AutoSaveState),
		skipTaskbar:       bi32(o.SkipTaskbar),
		noActivate:        bi32(o.NoActivate),
	}
}

func (s *WebViewSettings) toC(pool *cstrs) cWebViewSettings {
	return cWebViewSettings{
		autoplay:               bi32(s.Autoplay),
		backgroundThrottling:   bi32(s.BackgroundThrottling),
		allowRightClick:        bi32(s.AllowRightClick),
		ua:                     pool.p(s.UserAgent),
		preloadJS:              pool.p(s.PreloadJS),
		allowFullscreen:        bi32(s.AllowFullscreen),
		postmessageWhitelist:   pool.p(s.PostMessageWhitelist),
		corsWhitelist:          pool.p(s.CORSWhitelist),
		autofill:               bi32(s.Autofill),
		generalAutofillEnabled: bi32(s.GeneralAutofillEnabled),
		incognito:              bi32(s.Incognito),
		disableClipboard:       bi32(s.DisableClipboard),
		proxyURL:               pool.p(s.ProxyURL),
		focused:                bi32(s.Focused),
	}
}

// CreateWindow 创建一个 WebView 窗口，返回 window_id（0 表示失败）。
//
//   - url      : 初始地址（http(s):// 或自定义协议）
//   - parentID : 父窗口 ID，0 表示无父窗口
//   - opts     : 窗口选项，nil 则用 DefaultWindowOptions()
//   - settings : WebView 高级设置，nil 则全用默认
func CreateWindow(url string, parentID uint32, opts *WindowOptions, settings *WebViewSettings) uint32 {
	pool := &cstrs{}
	if opts == nil {
		d := DefaultWindowOptions()
		opts = &d
	}
	if settings == nil {
		settings = &WebViewSettings{}
	}
	copts := opts.toC(pool)
	csettings := settings.toC(pool)

	r, _, _ := procCreateWebviewWindow.Call(
		uintptr(unsafe.Pointer(pool.p(url))),
		uintptr(parentID),
		uintptr(unsafe.Pointer(&copts)),
		uintptr(unsafe.Pointer(&csettings)),
	)
	runtime.KeepAlive(pool)
	runtime.KeepAlive(&copts)
	runtime.KeepAlive(&csettings)
	return uint32(r)
}

// CreateBorderlessWindow 创建独立无边框 WebView 窗口，返回 window_id。
func CreateBorderlessWindow(url string, settings *WebViewSettings) uint32 {
	pool := &cstrs{}
	if settings == nil {
		settings = &WebViewSettings{}
	}
	csettings := settings.toC(pool)
	r, _, _ := procCreateBorderlessWindow.Call(
		uintptr(unsafe.Pointer(pool.p(url))),
		uintptr(unsafe.Pointer(&csettings)),
	)
	runtime.KeepAlive(pool)
	runtime.KeepAlive(&csettings)
	return uint32(r)
}

// --- 常用窗口操作（薄封装） ---

func Navigate(windowID uint32, url, headersJSON string) bool {
	pool := &cstrs{}
	r, _, _ := procNavigateToURL.Call(uintptr(windowID),
		uintptr(unsafe.Pointer(pool.p(url))), uintptr(unsafe.Pointer(pool.p(headersJSON))))
	runtime.KeepAlive(pool)
	return i32(r) == 1
}

func Reload(windowID uint32) bool {
	r, _, _ := procReloadWebviewWindow.Call(uintptr(windowID))
	return i32(r) == 1
}

// ExecuteJavaScript 执行 JS，返回一个唯一 id；结果通过 "javascript-result" 事件异步返回。
func ExecuteJavaScript(windowID uint32, script string) int32 {
	pool := &cstrs{}
	r, _, _ := procExecuteJavascript.Call(uintptr(windowID), uintptr(unsafe.Pointer(pool.p(script))))
	runtime.KeepAlive(pool)
	return i32(r)
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
// 判断 window_id 以确认初始化结果（详见头文件说明）。
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
func GetWebViewURL(windowID uint32) (string, bool) {
	return bufCallInt(2048, func(buf unsafe.Pointer, size int32) int32 {
		r, _, _ := procGetWebviewURL.Call(uintptr(windowID), uintptr(buf), uintptr(uint32(size)))
		return i32(r)
	})
}

// GetWindowHWND 返回窗口原生句柄（beta.9 起支持所有方法创建的窗口）。
func GetWindowHWND(windowID uint32) uintptr {
	r, _, _ := procGetWindowHwnd.Call(uintptr(windowID))
	return r
}

// GetWindowID 根据 HWND 反查窗口 ID（beta.9 新增），0 表示未找到。
func GetWindowID(hwnd uintptr) uint32 {
	r, _, _ := procGetWindowID.Call(uintptr(uint32(hwnd)))
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
func SetZoom(windowID uint32, level float64) bool {
	return callF64I32(procSetWebviewZoom, windowID, level) == 1
}

func SetFrameStyle(windowID uint32, frameStyle string) bool {
	pool := &cstrs{}
	r, _, _ := procSetWindowFrameStyle.Call(uintptr(windowID), uintptr(unsafe.Pointer(pool.p(frameStyle))))
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

// SetTitlebarOverlayStyle 设置标题栏覆盖层样式（Windows only）。height<=0 不修改高度。
func SetTitlebarOverlayStyle(windowID uint32, height int, iconColorHex, hoverBgHex string) bool {
	pool := &cstrs{}
	r, _, _ := procSetTitlebarOverlay.Call(uintptr(windowID), uintptr(uint32(int32(height))),
		uintptr(unsafe.Pointer(pool.p(iconColorHex))), uintptr(unsafe.Pointer(pool.p(hoverBgHex))))
	runtime.KeepAlive(pool)
	return i32(r) == 1
}

// --- IPC / DevTools / 其它 ---

// SendIPCMessage 向前端发送 IPC 消息。
func SendIPCMessage(windowID uint32, messageType, messageContent string) bool {
	pool := &cstrs{}
	r, _, _ := procSendIPCMessage.Call(uintptr(windowID),
		uintptr(unsafe.Pointer(pool.p(messageType))), uintptr(unsafe.Pointer(pool.p(messageContent))))
	runtime.KeepAlive(pool)
	return i32(r) == 1
}

func OpenDevtools(windowID uint32) bool {
	r, _, _ := procOpenDevtools.Call(uintptr(windowID))
	return i32(r) == 1
}

func CloseDevtools(windowID uint32) bool {
	r, _, _ := procCloseDevtools.Call(uintptr(windowID))
	return i32(r) == 1
}

func IsDevtoolsOpen(windowID uint32) bool {
	r, _, _ := procIsDevtoolsOpen.Call(uintptr(windowID))
	return i32(r) == 1
}

func ClearBrowsingData(windowID uint32) bool {
	r, _, _ := procClearBrowsingData.Call(uintptr(windowID))
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
