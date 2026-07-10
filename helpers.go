//go:build linux

package jadeview

/*
#include <stdlib.h>
#include "JadeView.h"
*/
import "C"

import "unsafe"

// goStringFree 把库返回的 char*（约定用 jade_text_free 释放）转成 Go string 并释放。
func goStringFree(c *C.char) string {
	if c == nil {
		return ""
	}
	s := C.GoString(c)
	C.jade_text_free(c)
	return s
}

// bufCallSize 调用「填充缓冲区、size_t 长度、返回 1=成功」类函数。
// 缓冲区不足时（库返回 0）按 16×/256× 放大重试两次（覆盖超长 URL 等场景）。
func bufCallSize(initial int, fn func(buf *C.char, n C.size_t) C.int32_t) (string, bool) {
	for _, size := range []int{initial, initial * 16, initial * 256} {
		buf := make([]byte, size)
		if fn((*C.char)(unsafe.Pointer(&buf[0])), C.size_t(size)) == 1 {
			return cBufToString(buf), true
		}
	}
	return "", false
}

// bufCallInt 同 bufCallSize，但长度参数为 int32。
func bufCallInt(initial int, fn func(buf *C.char, n C.int32_t) C.int32_t) (string, bool) {
	for _, size := range []int{initial, initial * 16, initial * 256} {
		buf := make([]byte, size)
		if fn((*C.char)(unsafe.Pointer(&buf[0])), C.int32_t(size)) == 1 {
			return cBufToString(buf), true
		}
	}
	return "", false
}

// yamlGet 执行 YAML 两阶段查询：先取所需字节数再分配读取。
// 返回 (值, 状态码)；状态码 1=成功，其余见文档（0=不存在，-1 IO，-2 类型，-4 解析）。
func yamlGet(fn func(buf *C.char, n C.size_t) C.int32_t) (string, int32) {
	n := fn(nil, 0)
	if n < 2 {
		// 0/负数：不存在或错误；1：理论上空值成功
		if n == 1 {
			return "", 1
		}
		return "", int32(n)
	}
	buf := make([]byte, int(n))
	rc := fn((*C.char)(unsafe.Pointer(&buf[0])), C.size_t(len(buf)))
	if rc == 1 {
		return cBufToString(buf), 1
	}
	return "", int32(rc)
}
