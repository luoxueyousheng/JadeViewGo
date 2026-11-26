//go:build linux && amd64 && jadeview_cgo && cgo

package jadeview

/*
#include <stdlib.h>
#include "JadeView.h"
*/
import "C"

import "unsafe"

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
