//go:build windows

package jadeview

import (
	"runtime"
	"sync"
	"syscall"
	"unsafe"
)

// --- 通知 ---

// cNotificationParams 是 C NotificationParams 的逐字段镜像。
type cNotificationParams struct {
	summary *byte
	body    *byte
	icon    *byte
	timeout int32
	button1 *byte
	button2 *byte
	text3   *byte
	action  *byte
}

// ShowNotification 显示系统通知。
func ShowNotification(p NotificationParams) bool {
	pool := &cstrs{}
	c := cNotificationParams{
		summary: pool.p(p.Summary),
		body:    pool.p(p.Body),
		icon:    pool.p(p.Icon),
		timeout: int32(p.Timeout),
		button1: pool.p(p.Button1),
		button2: pool.p(p.Button2),
		text3:   pool.p(p.Text3),
		action:  pool.p(p.Action),
	}
	r, _, _ := procShowNotification.Call(uintptr(unsafe.Pointer(&c)))
	runtime.KeepAlive(pool)
	runtime.KeepAlive(&c)
	return i32(r) == 1
}

// --- 文件对话框 ---

// cFileDialogParams 是 C FileDialogParams 的逐字段镜像。
type cFileDialogParams struct {
	windowID    uint32
	title       *byte
	defaultPath *byte
	buttonLabel *byte
	filters     *byte
	properties  *byte
}

func (p *FileDialogParams) toC(pool *cstrs) cFileDialogParams {
	return cFileDialogParams{
		windowID:    p.WindowID,
		title:       pool.p(p.Title),
		defaultPath: pool.p(p.DefaultPath),
		buttonLabel: pool.p(p.ButtonLabel),
		filters:     pool.p(p.Filters),
		properties:  pool.p(p.Properties),
	}
}

// ShowOpenDialog 显示打开文件对话框，返回结果 JSON（取消时通常为空/null）。
func ShowOpenDialog(p FileDialogParams) string {
	pool := &cstrs{}
	c := p.toC(pool)
	r, _, _ := procShowOpenDialog.Call(uintptr(unsafe.Pointer(&c)))
	runtime.KeepAlive(pool)
	runtime.KeepAlive(&c)
	return goStringFree(r)
}

// ShowSaveDialog 显示保存文件对话框，返回结果 JSON。
func ShowSaveDialog(p FileDialogParams) string {
	pool := &cstrs{}
	c := p.toC(pool)
	r, _, _ := procShowSaveDialog.Call(uintptr(unsafe.Pointer(&c)))
	runtime.KeepAlive(pool)
	runtime.KeepAlive(&c)
	return goStringFree(r)
}

// --- 消息框 ---

// cMessageBoxParams 是 C MessageBoxParams 的逐字段镜像。
type cMessageBoxParams struct {
	windowID  uint32
	title     *byte
	message   *byte
	detail    *byte
	buttons   *byte
	defaultID int32
	cancelID  int32
	typ       *byte
}

func (p *MessageBoxParams) toC(pool *cstrs) cMessageBoxParams {
	return cMessageBoxParams{
		windowID:  p.WindowID,
		title:     pool.p(p.Title),
		message:   pool.p(p.Message),
		detail:    pool.p(p.Detail),
		buttons:   pool.p(p.Buttons),
		defaultID: int32(p.DefaultID),
		cancelID:  int32(p.CancelID),
		typ:       pool.p(p.Type),
	}
}

// ShowMessageBox 显示消息框，返回结果 JSON（含点击的按钮索引）。
func ShowMessageBox(p MessageBoxParams) string {
	pool := &cstrs{}
	c := p.toC(pool)
	r, _, _ := procShowMessageBox.Call(uintptr(unsafe.Pointer(&c)))
	runtime.KeepAlive(pool)
	runtime.KeepAlive(&c)
	return goStringFree(r)
}

// ShowErrorBox 显示错误框。
func ShowErrorBox(windowID uint32, title, content string) bool {
	pool := &cstrs{}
	r, _, _ := procShowErrorBox.Call(uintptr(windowID),
		uintptr(unsafe.Pointer(pool.p(title))), uintptr(unsafe.Pointer(pool.p(content))))
	runtime.KeepAlive(pool)
	return i32(r) == 1
}

// --- 右键/上下文菜单项 ---

// MenuItemCreate 创建菜单项，返回 menu_id（0=失败）。
//   - kind: 见 MenuKind* 常量
//   - parentMenuID: 0=顶级，>0=加入指定子菜单
//   - itemID: 用户自定义 ID，点击时通过 "menu-item-clicked" 事件回传
func MenuItemCreate(label string, kind int, parentMenuID uint32, itemID int) uint32 {
	pool := &cstrs{}
	r, _, _ := procMenuItemCreate.Call(uintptr(unsafe.Pointer(pool.p(label))),
		uintptr(uint32(int32(kind))), uintptr(parentMenuID), uintptr(uint32(int32(itemID))))
	runtime.KeepAlive(pool)
	return uint32(r)
}

func MenuItemSetEnabled(menuID uint32, enabled bool) bool {
	r, _, _ := procMenuItemSetEnabled.Call(uintptr(menuID), b2u(enabled))
	return i32(r) == 1
}

func MenuItemSetChecked(menuID uint32, checked bool) bool {
	r, _, _ := procMenuItemSetChecked.Call(uintptr(menuID), b2u(checked))
	return i32(r) == 1
}

func MenuItemDestroy(menuID uint32) bool {
	r, _, _ := procMenuItemDestroy.Call(uintptr(menuID))
	return i32(r) == 1
}

// SetContextMenuItems 设置当前右键菜单要显示的顶级菜单项（在 "context-menu" 事件回调中调用）。
func SetContextMenuItems(windowID uint32, menuIDs []uint32) bool {
	if len(menuIDs) == 0 {
		r, _, _ := procSetContextMenuItems.Call(uintptr(windowID), 0, 0)
		return i32(r) == 1
	}
	r, _, _ := procSetContextMenuItems.Call(uintptr(windowID),
		uintptr(unsafe.Pointer(&menuIDs[0])), uintptr(uint32(int32(len(menuIDs)))))
	runtime.KeepAlive(menuIDs)
	return i32(r) == 1
}

// --- 异步对话框回调桥 ---
// 完成后经 void(const char*) 回调回传结果 JSON。固定槽位池 + 惰性 NewCallback 跳板。
// 该回调类型无 JADEVIEW_CALL 修饰，386 下为 cdecl，用 NewCallbackCDecl（64 位与 NewCallback 等价）。

var (
	dlgMu     sync.Mutex
	dlgSlots  [MaxAsyncDialogs]DialogResultHandler
	dlgTramps [MaxAsyncDialogs]uintptr
)

// dlgTramp 返回槽位跳板地址（惰性创建，调用方须持有 dlgMu）。
func dlgTramp(slot int) uintptr {
	if dlgTramps[slot] == 0 {
		s := slot
		dlgTramps[slot] = syscall.NewCallbackCDecl(func(result uintptr) uintptr {
			goDialogDispatch(s, result)
			return 0
		})
	}
	return dlgTramps[slot]
}

func goDialogDispatch(slot int, result uintptr) {
	dlgMu.Lock()
	h := dlgSlots[slot]
	dlgSlots[slot] = nil // 一次性，回调后释放槽位
	dlgMu.Unlock()
	if h != nil {
		h(goString(result))
	}
}

// dlgAcquire 占用一个空闲槽位，返回 (槽位跳板地址, 槽位号, ok)。
func dlgAcquire(h DialogResultHandler) (uintptr, int, bool) {
	dlgMu.Lock()
	defer dlgMu.Unlock()
	for i := range dlgSlots {
		if dlgSlots[i] == nil {
			dlgSlots[i] = h
			return dlgTramp(i), i, true
		}
	}
	return 0, 0, false
}

func dlgRelease(slot int) {
	dlgMu.Lock()
	dlgSlots[slot] = nil
	dlgMu.Unlock()
}

// ShowOpenDialogAsync 异步显示打开文件对话框，结果通过 handler 回传。
// 返回 false 表示槽位已满或库调用失败。
func ShowOpenDialogAsync(p FileDialogParams, handler DialogResultHandler) bool {
	tramp, slot, ok := dlgAcquire(handler)
	if !ok {
		return false
	}
	pool := &cstrs{}
	c := p.toC(pool)
	r, _, _ := procShowOpenDialogAsync.Call(uintptr(unsafe.Pointer(&c)), tramp)
	runtime.KeepAlive(pool)
	runtime.KeepAlive(&c)
	if i32(r) != 1 {
		dlgRelease(slot)
		return false
	}
	return true
}

// ShowSaveDialogAsync 异步显示保存文件对话框。
func ShowSaveDialogAsync(p FileDialogParams, handler DialogResultHandler) bool {
	tramp, slot, ok := dlgAcquire(handler)
	if !ok {
		return false
	}
	pool := &cstrs{}
	c := p.toC(pool)
	r, _, _ := procShowSaveDialogAsync.Call(uintptr(unsafe.Pointer(&c)), tramp)
	runtime.KeepAlive(pool)
	runtime.KeepAlive(&c)
	if i32(r) != 1 {
		dlgRelease(slot)
		return false
	}
	return true
}

// ShowMessageBoxAsync 异步显示消息框。
func ShowMessageBoxAsync(p MessageBoxParams, handler DialogResultHandler) bool {
	tramp, slot, ok := dlgAcquire(handler)
	if !ok {
		return false
	}
	pool := &cstrs{}
	c := p.toC(pool)
	r, _, _ := procShowMessageBoxAsync.Call(uintptr(unsafe.Pointer(&c)), tramp)
	runtime.KeepAlive(pool)
	runtime.KeepAlive(&c)
	if i32(r) != 1 {
		dlgRelease(slot)
		return false
	}
	return true
}

// --- 托盘 ---

// cTrayMenuItemDesc 是 C TrayMenuItemDesc 的逐字段镜像。
type cTrayMenuItemDesc struct {
	itemType  int32
	key       *byte
	label     *byte
	parentKey *byte
	disabled  int32
	dangerous int32
}

// TrayCreate 创建托盘图标，返回 tray_id（0=失败）。
func TrayCreate() uint32 {
	r, _, _ := procTrayCreate.Call()
	return uint32(r)
}

func TrayDestroy(trayID uint32) bool {
	r, _, _ := procTrayDestroy.Call(uintptr(trayID))
	return i32(r) == 1
}

func TraySetVisible(trayID uint32, visible bool) bool {
	r, _, _ := procTraySetVisible.Call(uintptr(trayID), b2u(visible))
	return i32(r) == 1
}

func TraySetTooltip(trayID uint32, tooltip string) bool {
	pool := &cstrs{}
	r, _, _ := procTraySetTooltip.Call(uintptr(trayID), uintptr(unsafe.Pointer(pool.p(tooltip))))
	runtime.KeepAlive(pool)
	return i32(r) == 1
}

func TraySetIconFromFile(trayID uint32, iconPath string) bool {
	pool := &cstrs{}
	r, _, _ := procTraySetIconFromFile.Call(uintptr(trayID), uintptr(unsafe.Pointer(pool.p(iconPath))))
	runtime.KeepAlive(pool)
	return i32(r) == 1
}

// TraySetIconFromData 用内存中的图标数据设置托盘图标。
func TraySetIconFromData(trayID uint32, data []byte) bool {
	if len(data) == 0 {
		return false
	}
	r, _, _ := procTraySetIconFromData.Call(uintptr(trayID),
		uintptr(unsafe.Pointer(&data[0])), uintptr(uint32(len(data))))
	runtime.KeepAlive(data)
	return i32(r) == 1
}

// TraySetMenu 用扁平表设置托盘根菜单（库内深拷贝字符串）。空切片表示清除菜单。
func TraySetMenu(trayID uint32, items []TrayMenuItem) bool {
	if len(items) == 0 {
		r, _, _ := procTraySetMenuItems.Call(uintptr(trayID), 0, 0)
		return i32(r) == 1
	}
	pool := &cstrs{}
	arr := make([]cTrayMenuItemDesc, len(items))
	for i, it := range items {
		arr[i] = cTrayMenuItemDesc{
			itemType:  int32(it.Type),
			key:       pool.p(it.Key),
			label:     pool.p(it.Label),
			parentKey: pool.p(it.ParentKey),
			disabled:  bi32(it.Disabled),
			dangerous: bi32(it.Dangerous),
		}
	}
	r, _, _ := procTraySetMenuItems.Call(uintptr(trayID),
		uintptr(unsafe.Pointer(&arr[0])), uintptr(uint32(len(arr))))
	runtime.KeepAlive(pool)
	runtime.KeepAlive(arr)
	return i32(r) == 1
}
