//go:build linux && amd64 && jadeview_cgo && cgo

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

func (p *MessageBoxParams) toC(pool *cstrPool) C.MessageBoxParams {
	return C.MessageBoxParams{
		window_id:  C.uint32_t(p.WindowID),
		title:      pool.s(p.Title),
		message:    pool.s(p.Message),
		detail:     pool.s(p.Detail),
		buttons:    pool.s(p.Buttons),
		default_id: C.int32_t(p.DefaultID),
		cancel_id:  C.int32_t(p.CancelID),
		type_:      pool.s(p.Type),
	}
}

// ShowMessageBox 显示消息框，返回结果 JSON（含点击的按钮索引）。
func ShowMessageBox(p MessageBoxParams) string {
	pool := &cstrPool{}
	defer pool.free()
	c := p.toC(pool)
	return goStringFree(C.jade_dialog_show_message_box(&c))
}

// ShowErrorBox 显示错误框。
func ShowErrorBox(windowID uint32, title, content string) bool {
	pool := &cstrPool{}
	defer pool.free()
	return C.jade_dialog_show_error_box(C.uint32_t(windowID), pool.s(title), pool.s(content)) == 1
}
