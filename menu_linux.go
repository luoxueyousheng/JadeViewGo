//go:build linux && amd64 && jadeview_cgo && cgo

package jadeview

/*
#include "JadeView.h"
*/
import "C"

import "unsafe"

func (icon MenuIcon) toC(pool *cstrPool) C.JadeMenuIconDesc {
	var data unsafe.Pointer
	if len(icon.Data) > 0 {
		data = C.CBytes(icon.Data)
		pool.ptrs = append(pool.ptrs, data)
	}
	return C.JadeMenuIconDesc{icon_path: pool.s(icon.Path), icon_data: (*C.uint8_t)(data), icon_data_size: C.size_t(len(icon.Data)), icon_codepoint: C.uint32_t(icon.Codepoint)}
}

func menuItems(pool *cstrPool, items []MenuItem) []C.JadeMenuItemDesc {
	if len(items) == 0 {
		return nil
	}
	a := make([]C.JadeMenuItemDesc, len(items))
	for i, m := range items {
		a[i] = C.JadeMenuItemDesc{id: pool.s(m.ID), parent_id: pool.s(m.ParentID), label: pool.s(m.Label), accelerator: pool.s(m.Accelerator), kind: C.uint32_t(m.Kind), flags: C.uint32_t(m.flags()), icon: m.Icon.toC(pool)}
	}
	return a
}

// PopupMenu 弹出菜单。windowID=0 时不绑定窗口；position=nil 时用鼠标位置。
func PopupMenu(windowID uint32, items []MenuItem, position *MenuPosition) bool {
	pool := &cstrPool{}
	defer pool.free()
	var pos *C.JadeMenuPosition
	if position != nil {
		pos = &C.JadeMenuPosition{x: C.int32_t(position.X), y: C.int32_t(position.Y)}
	}
	return C.jade_menu_popup(C.uint32_t(windowID), unsafe.SliceData(menuItems(pool, items)), C.uint32_t(len(items)), pos) == 1
}

// SetApplicationMenu 替换应用菜单，空表移除菜单。
func SetApplicationMenu(items []MenuItem) bool {
	pool := &cstrPool{}
	defer pool.free()
	return C.jade_set_application_menu(unsafe.SliceData(menuItems(pool, items)), C.uint32_t(len(items))) == 1
}

// SetWindowMenu 替换指定 Window 的菜单，空表移除菜单。
func SetWindowMenu(windowID uint32, items []MenuItem) bool {
	pool := &cstrPool{}
	defer pool.free()
	return C.jade_window_set_menu(C.uint32_t(windowID), unsafe.SliceData(menuItems(pool, items)), C.uint32_t(len(items))) == 1
}

// SetContextMenu 替换指定 WebView 的右键菜单模板，空表清除模板。
func SetContextMenu(webviewID uint32, items []ContextMenuItem) bool {
	pool := &cstrPool{}
	defer pool.free()
	var a []C.JadeContextMenuItemDesc
	if len(items) > 0 {
		a = make([]C.JadeContextMenuItemDesc, len(items))
	}
	for i, m := range items {
		a[i] = C.JadeContextMenuItemDesc{struct_size: C.sizeof_JadeContextMenuItemDesc, id: pool.s(m.ID), parent_id: pool.s(m.ParentID), label: pool.s(m.Label), accelerator: pool.s(m.Accelerator), role: pool.s(m.Role), kind: C.uint32_t(m.Kind), flags: C.uint32_t(m.flags()), icon: m.Icon.toC(pool)}
	}
	return C.jade_webview_set_context_menu(C.uint32_t(webviewID), unsafe.SliceData(a), C.uint32_t(len(a))) == 1
}

func ClearContextMenu(view uint32) bool {
	return C.jade_webview_clear_context_menu(C.uint32_t(view)) == 1
}
