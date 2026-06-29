package jadeview

/*
#include <stdlib.h>
#include "JadeView.h"
*/
import "C"

import "unsafe"

// 托盘菜单项类型（TrayMenuItemDesc.item_type）。
const (
	TrayItemNormal  = 0 // 普通项
	TrayItemSubmenu = 1 // 子菜单
	TrayItemDivider = 2 // 分隔线
	TrayItemGroup   = 3 // 分组
)

// TrayMenuItem 对应 C 的 TrayMenuItemDesc（扁平表，用 ParentKey 指向父项的 Key）。
type TrayMenuItem struct {
	Type      int    // 见 TrayItem* 常量
	Key       string // 全表唯一、非空（分隔线也需唯一 key）
	Label     string
	ParentKey string // 空=根下子项；否则须等于某行的 Key 且该行为 Submenu/Group
	Disabled  bool
	Dangerous bool
}

// TrayCreate 创建托盘图标，返回 tray_id（0=失败）。
func TrayCreate() uint32 {
	return uint32(C.tray_create())
}

func TrayDestroy(trayID uint32) bool {
	return C.tray_destroy(C.uint32_t(trayID)) == 1
}

func TraySetVisible(trayID uint32, visible bool) bool {
	return C.tray_set_visible(C.uint32_t(trayID), b2i(visible)) == 1
}

func TraySetTooltip(trayID uint32, tooltip string) bool {
	pool := &cstrPool{}
	defer pool.free()
	return C.tray_set_tooltip(C.uint32_t(trayID), pool.s(tooltip)) == 1
}

func TraySetIconFromFile(trayID uint32, iconPath string) bool {
	pool := &cstrPool{}
	defer pool.free()
	return C.tray_set_icon_from_file(C.uint32_t(trayID), pool.s(iconPath)) == 1
}

// TraySetIconFromData 用内存中的图标数据设置托盘图标。
func TraySetIconFromData(trayID uint32, data []byte) bool {
	if len(data) == 0 {
		return false
	}
	return C.set_tray_icon_from_data(C.uint32_t(trayID), (*C.uint8_t)(unsafe.Pointer(&data[0])), C.uint32_t(len(data))) == 1
}

// TraySetMenu 用扁平表设置托盘根菜单（库内深拷贝字符串）。空切片表示清除菜单。
func TraySetMenu(trayID uint32, items []TrayMenuItem) bool {
	if len(items) == 0 {
		return C.tray_set_menu_items(C.uint32_t(trayID), nil, 0) == 1
	}
	pool := &cstrPool{}
	defer pool.free()
	arr := make([]C.TrayMenuItemDesc, len(items))
	for i, it := range items {
		arr[i] = C.TrayMenuItemDesc{
			item_type:  C.int32_t(it.Type),
			key:        pool.s(it.Key),
			label:      pool.s(it.Label),
			parent_key: pool.s(it.ParentKey),
			disabled:   b2i(it.Disabled),
			dangerous:  b2i(it.Dangerous),
		}
	}
	return C.tray_set_menu_items(C.uint32_t(trayID), &arr[0], C.uint32_t(len(arr))) == 1
}
