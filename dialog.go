//go:build linux

package jadeview

/*
#include <stdlib.h>
#include "JadeView.h"
*/
import "C"

// --- 通知 ---

// ShowNotification 显示系统通知。
func ShowNotification(p NotificationParams) bool {
	pool := &cstrPool{}
	defer pool.free()
	c := C.NotificationParams{
		summary: pool.s(p.Summary),
		body:    pool.s(p.Body),
		icon:    pool.s(p.Icon),
		timeout: C.int32_t(p.Timeout),
		button1: pool.s(p.Button1),
		button2: pool.s(p.Button2),
		text3:   pool.s(p.Text3),
		action:  pool.s(p.Action),
	}
	return C.show_notification(&c) == 1
}

// --- 文件对话框 ---

func (p *FileDialogParams) toC(pool *cstrPool) C.FileDialogParams {
	return C.FileDialogParams{
		window_id:    C.uint32_t(p.WindowID),
		title:        pool.s(p.Title),
		default_path: pool.s(p.DefaultPath),
		button_label: pool.s(p.ButtonLabel),
		filters:      pool.s(p.Filters),
		properties:   pool.s(p.Properties),
	}
}

// ShowOpenDialog 显示打开文件对话框，返回结果 JSON（取消时通常为空/null）。
func ShowOpenDialog(p FileDialogParams) string {
	pool := &cstrPool{}
	defer pool.free()
	c := p.toC(pool)
	return goStringFree(C.jade_dialog_show_open_dialog(&c))
}

// ShowSaveDialog 显示保存文件对话框，返回结果 JSON。
func ShowSaveDialog(p FileDialogParams) string {
	pool := &cstrPool{}
	defer pool.free()
	c := p.toC(pool)
	return goStringFree(C.jade_dialog_show_save_dialog(&c))
}

// --- 消息框 ---

// ShowMessageBox 显示消息框，返回结果 JSON（含点击的按钮索引）。
func ShowMessageBox(p MessageBoxParams) string {
	pool := &cstrPool{}
	defer pool.free()
	c := C.MessageBoxParams{
		window_id:  C.uint32_t(p.WindowID),
		title:      pool.s(p.Title),
		message:    pool.s(p.Message),
		detail:     pool.s(p.Detail),
		buttons:    pool.s(p.Buttons),
		default_id: C.int32_t(p.DefaultID),
		cancel_id:  C.int32_t(p.CancelID),
		type_:      pool.s(p.Type),
	}
	return goStringFree(C.jade_dialog_show_message_box(&c))
}

// ShowErrorBox 显示错误框。
func ShowErrorBox(windowID uint32, title, content string) bool {
	pool := &cstrPool{}
	defer pool.free()
	return C.jade_dialog_show_error_box(C.uint32_t(windowID), pool.s(title), pool.s(content)) == 1
}

// --- 右键/上下文菜单项 ---

// MenuItemCreate 创建菜单项，返回 menu_id（0=失败）。
//   - kind: 见 MenuKind* 常量
//   - parentMenuID: 0=顶级，>0=加入指定子菜单
//   - itemID: 用户自定义 ID，点击时通过 "menu-item-clicked" 事件回传
func MenuItemCreate(label string, kind int, parentMenuID uint32, itemID int) uint32 {
	pool := &cstrPool{}
	defer pool.free()
	return uint32(C.jade_menu_item_create(pool.s(label), C.int32_t(kind), C.uint32_t(parentMenuID), C.int32_t(itemID)))
}

func MenuItemSetEnabled(menuID uint32, enabled bool) bool {
	return C.jade_menu_item_set_enabled(C.uint32_t(menuID), b2i(enabled)) == 1
}

func MenuItemSetChecked(menuID uint32, checked bool) bool {
	return C.jade_menu_item_set_checked(C.uint32_t(menuID), b2i(checked)) == 1
}

func MenuItemDestroy(menuID uint32) bool {
	return C.jade_menu_item_destroy(C.uint32_t(menuID)) == 1
}

// SetContextMenuItems 设置当前右键菜单要显示的顶级菜单项（在 "context-menu" 事件回调中调用）。
func SetContextMenuItems(windowID uint32, menuIDs []uint32) bool {
	if len(menuIDs) == 0 {
		return C.jade_set_context_menu_items(C.uint32_t(windowID), nil, 0) == 1
	}
	ids := make([]C.uint32_t, len(menuIDs))
	for i, v := range menuIDs {
		ids[i] = C.uint32_t(v)
	}
	return C.jade_set_context_menu_items(C.uint32_t(windowID), &ids[0], C.int32_t(len(ids))) == 1
}
