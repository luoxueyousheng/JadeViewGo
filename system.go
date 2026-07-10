//go:build linux

package jadeview

/*
#include <stdlib.h>
#include "JadeView.h"
*/
import "C"

import "unsafe"

// --- 剪贴板 ---

func ClipboardReadText() (string, bool) {
	return bufCallInt(4096, func(buf *C.char, n C.int32_t) C.int32_t {
		return C.clipboard_read_text(buf, n)
	})
}

func ClipboardWriteText(text string) bool {
	pool := &cstrPool{}
	defer pool.free()
	return C.clipboard_write_text(pool.s(text)) == 1
}

// --- 路径 / 系统信息 ---

// GetPath 获取系统路径（name 如 home/appData/temp 等，具体见库文档）。
func GetPath(name string) (string, bool) {
	pool := &cstrPool{}
	defer pool.free()
	cn := pool.s(name)
	return bufCallSize(1024, func(buf *C.char, n C.size_t) C.int32_t {
		return C.getPath(cn, buf, n)
	})
}

// GetLocale 获取系统语言（BCP 47，如 zh-CN）。
func GetLocale() (string, bool) {
	return bufCallSize(64, func(buf *C.char, n C.size_t) C.int32_t {
		return C.getLocale(buf, n)
	})
}

// GetDisplaysInfo 返回显示器信息 JSON 数组（含 bounds、work_area、scale_factor、dpi、is_primary）。
func GetDisplaysInfo() (string, bool) {
	return bufCallSize(8192, func(buf *C.char, n C.size_t) C.int32_t {
		return C.get_displays_info(buf, n)
	})
}

// GetCursorPosition 返回光标位置 JSON。
func GetCursorPosition() (string, bool) {
	return bufCallInt(128, func(buf *C.char, n C.int32_t) C.int32_t {
		return C.get_cursor_position(buf, n)
	})
}

// GetWebViewVersion 获取 WebView 内核版本号。
func GetWebViewVersion() (string, bool) {
	return bufCallSize(256, func(buf *C.char, n C.size_t) C.int32_t {
		return C.get_webview_version(buf, n)
	})
}

// IsWindows11 检查当前系统是否为 Windows 11。
func IsWindows11() bool {
	return C.is_windows_11() == 1
}

// GetPrinterList 返回打印机名称 JSON 数组和打印机数量。
// 注意：jade_get_printer_list 返回的是打印机数量（非 1=成功），不能走 bufCallInt；
// 头文件未定义缓冲不足时的行为（16KB 对打印机名列表实际足够）。
func GetPrinterList() (string, int32, bool) {
	buf := make([]byte, 16384)
	rc := C.jade_get_printer_list((*C.char)(unsafe.Pointer(&buf[0])), C.int32_t(len(buf)))
	if rc < 0 {
		return "", int32(rc), false
	}
	return cBufToString(buf), int32(rc), true
}

// --- 打印 ---

// Print 打开 WebView2 内置打印对话框。
func Print(windowID uint32) bool {
	return C.jade_print(C.uint32_t(windowID)) == 1
}

// PrintFile 用系统关联程序打印文件（Windows=ShellExecute "print"，Linux=CUPS lp）。
func PrintFile(filePath string) bool {
	pool := &cstrPool{}
	defer pool.free()
	return C.jade_print_dialog(pool.s(filePath)) == 1
}

// --- NTP 网络时间 ---

// NTPNow 获取网络时间戳（UTC 毫秒，北京时间需 +8 小时）。
// server 为空时使用内置服务器列表逐个尝试。失败返回 -1。
func NTPNow(server string) int64 {
	pool := &cstrPool{}
	defer pool.free()
	return int64(C.jade_ntp_now(pool.s(server)))
}

// --- 全局热键 ---

// RegisterGlobalHotkey 注册全局热键，返回 hotkey_id（0=失败）。触发时通过 "global-hotkey" 事件回报。
func RegisterGlobalHotkey(modifiers, vk uint32) uint32 {
	return uint32(C.register_global_hotkey(C.uint32_t(modifiers), C.uint32_t(vk)))
}

func UnregisterGlobalHotkey(hotkeyID uint32) bool {
	return C.unregister_global_hotkey(C.uint32_t(hotkeyID)) == 1
}

// --- 开机自启 ---

// SetLoginAutostart 启用/取消开机自启。args 为追加到可执行路径后的启动参数（可空）。
func SetLoginAutostart(enable bool, args string) bool {
	pool := &cstrPool{}
	defer pool.free()
	return C.set_login_autostart(b2i(enable), pool.s(args)) == 1
}

func GetLoginAutostart() bool {
	return C.get_login_autostart() == 1
}

// --- URL 协议 / 文件关联 ---

func RegisterURLScheme(scheme string) bool {
	pool := &cstrPool{}
	defer pool.free()
	return C.register_url_scheme(pool.s(scheme)) == 1
}

func UnregisterURLScheme(scheme string) bool {
	pool := &cstrPool{}
	defer pool.free()
	return C.unregister_url_scheme(pool.s(scheme)) == 1
}

func RegisterFileAssociation(extension, friendlyName string) bool {
	pool := &cstrPool{}
	defer pool.free()
	return C.register_file_association(pool.s(extension), pool.s(friendlyName)) == 1
}

func UnregisterFileAssociation(extension string) bool {
	pool := &cstrPool{}
	defer pool.free()
	return C.unregister_file_association(pool.s(extension)) == 1
}

// --- 安全资源 / 协议服务 ---

// RegisterResource 注册本地文件为安全资源，返回 jade:// URL。
// windowID=0 表示全局；ttlSeconds=0 表示永不过期。
func RegisterResource(path string, windowID, ttlSeconds uint32) (string, bool) {
	pool := &cstrPool{}
	defer pool.free()
	cp := pool.s(path)
	return bufCallSize(512, func(buf *C.char, n C.size_t) C.int32_t {
		return C.register_resource(cp, C.uint32_t(windowID), C.uint32_t(ttlSeconds), buf, n)
	})
}

func UnregisterResource(tokenOrURL string) bool {
	pool := &cstrPool{}
	defer pool.free()
	return C.unregister_resource(pool.s(tokenOrURL)) == 1
}

// ClearWindowResources 清理指定窗口的所有已注册资源，返回清理数量。
func ClearWindowResources(windowID uint32) int32 {
	return int32(C.clear_window_resources(C.uint32_t(windowID)))
}

// GetFileIcon 提取任意路径的图标为 PNG 注册为 jade:// 资源，返回 URL。
// size 目标边长（16/32/48/64/128/256，<=0 取 48）；windowID=0 全局；ttlSeconds=0 默认。
func GetFileIcon(path string, size int, windowID, ttlSeconds uint32) (string, bool) {
	pool := &cstrPool{}
	defer pool.free()
	cp := pool.s(path)
	return bufCallSize(512, func(buf *C.char, n C.size_t) C.int32_t {
		return C.get_file_icon(cp, C.int32_t(size), C.uint32_t(windowID), C.uint32_t(ttlSeconds), buf, n)
	})
}

// SetProtocolServicePath 设置自定义协议服务路径，返回可访问的 URL。
// hotReload 仅文件系统模式有效。
func SetProtocolServicePath(rootPath string, hotReload bool) (string, bool) {
	pool := &cstrPool{}
	defer pool.free()
	cr := pool.s(rootPath)
	return bufCallSize(512, func(buf *C.char, n C.size_t) C.int32_t {
		return C.set_protocol_service_path(cr, buf, n, b2i(hotReload))
	})
}

// ClearDataDirectory 清空数据目录。confirm_token 必须等于 "I_UNDERSTAND_CLEAR_DATA"。
func ClearDataDirectory(confirmToken string) bool {
	pool := &cstrPool{}
	defer pool.free()
	return C.clear_data_directory(pool.s(confirmToken)) == 1
}

// --- 编码转换 ---

// SmartConvertEncoding 智能检测输入编码并转换为目标编码。
// 返回 (转换结果, 检测到的源编码, 是否成功)。targetEncoding 见库文档（utf-8/gbk/big5 等）。
func SmartConvertEncoding(input []byte, targetEncoding string) (string, string, bool) {
	pool := &cstrPool{}
	defer pool.free()
	ct := pool.s(targetEncoding)
	var inPtr *C.uint8_t
	if len(input) > 0 {
		inPtr = (*C.uint8_t)(unsafe.Pointer(&input[0]))
	}

	outSize := len(input)*4 + 64
	detSize := 64
	for attempt := 0; attempt < 2; attempt++ {
		out := make([]byte, outSize)
		det := make([]byte, detSize)
		rc := C.smart_convert_encoding(
			inPtr, C.int32_t(len(input)), ct,
			(*C.char)(unsafe.Pointer(&out[0])), C.int32_t(outSize),
			(*C.char)(unsafe.Pointer(&det[0])), C.int32_t(detSize),
		)
		if rc > 0 {
			return string(out[:int(rc)]), cBufToString(det), true
		}
		if rc < 0 {
			// 缓冲区不足，绝对值为所需大小，放大重试
			outSize = int(-rc) + 16
			continue
		}
		break // 0 = 失败
	}
	return "", "", false
}
