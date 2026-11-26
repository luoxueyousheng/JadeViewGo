//go:build windows || (linux && amd64 && !jadeview_cgo)

package jadeview

import (
	"runtime"
	"unsafe"
)

type cMenuIcon struct {
	path      *byte
	data      *byte
	size      uintptr
	codepoint uint32
}

type cMenuItem struct {
	id, parentID, label, accelerator *byte
	kind, flags                      uint32
	icon                             cMenuIcon
}

type cContextMenuItem struct {
	size                                   uint32
	id, parentID, label, accelerator, role *byte
	kind, flags                            uint32
	icon                                   cMenuIcon
}

func (icon MenuIcon) toC(pool *cstrs) cMenuIcon {
	return cMenuIcon{pool.p(icon.Path), unsafe.SliceData(icon.Data), uintptr(len(icon.Data)), icon.Codepoint}
}

func menuItems(pool *cstrs, items []MenuItem) []cMenuItem {
	if len(items) == 0 {
		return nil
	}
	out := make([]cMenuItem, len(items))
	for i, m := range items {
		out[i] = cMenuItem{pool.p(m.ID), pool.p(m.ParentID), pool.p(m.Label), pool.p(m.Accelerator), uint32(m.Kind), m.flags(), m.Icon.toC(pool)}
	}
	return out
}

var (
	procPopupMenu          = np("jade_menu_popup")
	procSetApplicationMenu = np("jade_set_application_menu")
	procSetWindowMenu      = np("jade_window_set_menu")
	procSetContextMenu     = np("jade_webview_set_context_menu")
)

// PopupMenu 弹出菜单。windowID=0 时不绑定窗口；position=nil 时用鼠标位置。
func PopupMenu(windowID uint32, items []MenuItem, position *MenuPosition) bool {
	pool := &cstrs{}
	a := menuItems(pool, items)
	r, _, _ := procPopupMenu.Call(uintptr(windowID), uintptr(unsafe.Pointer(unsafe.SliceData(a))), uintptr(len(a)), uintptr(unsafe.Pointer(position)))
	runtime.KeepAlive(pool)
	return i32(r) == 1
}

// SetApplicationMenu 替换应用菜单，空表移除菜单。
func SetApplicationMenu(items []MenuItem) bool {
	pool := &cstrs{}
	a := menuItems(pool, items)
	r, _, _ := procSetApplicationMenu.Call(uintptr(unsafe.Pointer(unsafe.SliceData(a))), uintptr(len(a)))
	runtime.KeepAlive(pool)
	return i32(r) == 1
}

// SetWindowMenu 替换指定 Window 的菜单，空表移除菜单。
func SetWindowMenu(windowID uint32, items []MenuItem) bool {
	pool := &cstrs{}
	a := menuItems(pool, items)
	r, _, _ := procSetWindowMenu.Call(uintptr(windowID), uintptr(unsafe.Pointer(unsafe.SliceData(a))), uintptr(len(a)))
	runtime.KeepAlive(pool)
	return i32(r) == 1
}

// SetContextMenu 替换指定 WebView 的右键菜单模板，空表清除模板。
func SetContextMenu(webviewID uint32, items []ContextMenuItem) bool {
	pool := &cstrs{}
	var a []cContextMenuItem
	if len(items) > 0 {
		a = make([]cContextMenuItem, len(items))
	}
	for i, m := range items {
		a[i] = cContextMenuItem{uint32(unsafe.Sizeof(cContextMenuItem{})), pool.p(m.ID), pool.p(m.ParentID), pool.p(m.Label), pool.p(m.Accelerator), pool.p(m.Role), uint32(m.Kind), m.flags(), m.Icon.toC(pool)}
	}
	r, _, _ := procSetContextMenu.Call(uintptr(webviewID), uintptr(unsafe.Pointer(unsafe.SliceData(a))), uintptr(len(a)))
	runtime.KeepAlive(pool)
	return i32(r) == 1
}
