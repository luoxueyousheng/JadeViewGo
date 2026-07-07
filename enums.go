package jadeview

// 固定取值的输入参数枚举，按用途分组为二级命名空间（Windows/Linux 两侧共用）：
//
//	jadeview.Theme.Dark            // 窗口主题
//	jadeview.FrameStyle.TitleOverlay
//	jadeview.Backdrop.Mica
//
// 取值与 include/JadeView.h 注释及官方文档一致；事件名见 events_names.go。

// Theme 窗口主题（WindowOptions.Theme / SetTheme）。
var Theme = struct {
	Light  string // 强制浅色
	Dark   string // 强制深色
	System string // 跟随 Windows 个性化设置（推荐默认）
}{"Light", "Dark", "System"}

// FrameStyle 窗口边框样式（WindowOptions.FrameStyle / SetFrameStyle）。
var FrameStyle = struct {
	Normal       string // 系统边框 + 标题栏
	NoTitlebar   string // 保留边框，去掉标题栏（需自绘标题栏与控制按钮）
	Borderless   string // 完全无边框（启动页/异形窗口）
	TitleOverlay string // 保留边框 + 无标题栏 + 库内置右上角控制按钮（Windows，推荐）
}{"normal", "no-titlebar", "borderless", "title-overlay"}

// WindowLevel 窗口层级（SetLevel）。
var WindowLevel = struct {
	Topmost string // 置顶
	Normal  string // 常规
	Bottom  string // 置底
	Desktop string // 贴桌面壁纸层（Windows 寄生 WorkerW / Linux desktop hint）
}{"topmost", "normal", "bottom", "desktop"}

// Backdrop 窗口材质（SetBackdrop，Windows 11 DWM）。
// 使用材质需建窗时 Transparent=true 且页面背景透明；纯色底用 SetBackgroundColor。
var Backdrop = struct {
	Mica    string // 云母：从壁纸取色，窗口主背景（低开销，推荐）
	MicaAlt string // 云母变体：色调更深，适合侧边栏/分区
	Acrylic string // 亚克力：实时高斯模糊，适合弹层/Flyout（中开销）
}{"mica", "micaAlt", "acrylic"}

// MsgBoxType 消息框类型（MessageBoxParams.Type），与状态色对应。
var MsgBoxType = struct {
	None     string
	Info     string // 信息（Accent）
	Warning  string // 警告（Caution）
	Error    string // 错误（Critical）
	Question string // 询问
}{"none", "info", "warning", "error", "question"}

// ProgressState 任务栏进度状态（SetWindowProgress 的 state 参数）。
var ProgressState = struct {
	None          int // 清除进度
	Normal        int // 正常（绿色）
	Paused        int // 暂停（黄色）
	Error         int // 出错（红色）
	Indeterminate int // 不确定进度（跑马灯）
}{0, 1, 2, 3, 4}

// TrayItem 托盘菜单项类型（TrayMenuItem.Type）。
var TrayItem = struct {
	Normal  int // 普通项
	Submenu int // 子菜单
	Divider int // 分隔线
	Group   int // 分组
}{0, 1, 2, 3}

// MenuKind 右键/上下文菜单项类型（MenuItemCreate 的 kind 参数）。
var MenuKind = struct {
	Normal    int // 普通命令
	Separator int // 分隔线
	Checkbox  int // 复选框
	Radio     int // 单选
	Submenu   int // 子菜单
}{0, 1, 2, 3, 4}

// DialogProp 文件对话框属性（FileDialogParams.Properties 的 JSON 数组元素）。
//
//	Properties: `["openFile","multiSelections"]`
var DialogProp = struct {
	OpenFile        string // 选择文件
	OpenDirectory   string // 选择目录
	MultiSelections string // 允许多选
	ShowHiddenFiles string // 显示隐藏文件
	PromptToCreate  string // 路径不存在时提示创建（Windows）
}{"openFile", "openDirectory", "multiSelections", "showHiddenFiles", "promptToCreate"}

// Encoding 常用目标编码（SmartConvertEncoding 的 targetEncoding，大小写不敏感）。
// 完整别名表见头文件注释与 https://encoding.spec.whatwg.org/。
var Encoding = struct {
	UTF8     string
	GBK      string
	GB18030  string
	Big5     string
	ShiftJIS string
	EUCKR    string
	Latin1   string // windows-1252（含 ascii/iso-8859-1 别名）
}{"utf-8", "gbk", "gb18030", "big5", "shift_jis", "euc-kr", "windows-1252"}
