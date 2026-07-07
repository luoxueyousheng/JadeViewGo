//go:build linux

package jadeview

/*
#include <stdlib.h>
#include "JadeView.h"
*/
import "C"

import "unsafe"

// JAPK 资源包：加密/签名的前端资源包，加载后通过 jade:// 协议访问。
// 返回值约定与其它模块不同：0=成功，负数=错误码。

// SetPublicKey 设置 Base64 编码的 Ed25519 公钥（44 字符），必须在 LoadFromBytes 之前调用。
// 返回 0=成功，负数=错误码。
func SetPublicKey(publicKey string) int32 {
	pool := &cstrPool{}
	defer pool.free()
	return int32(C.JadeView_set_public_key(pool.s(publicKey)))
}

// LoadFromBytes 从内存加载 JAPK 文件（支持 v2 签名包与混淆包）。
//   - 已设公钥：必须是签名包，验签后加载；
//   - 未设公钥：仅支持混淆包（JPKBIN02）。
//
// app_name / app_signature 须与 Init 时一致。返回 0=成功，负数=错误码。
// 错误信息也会通过事件异步通知。
func LoadFromBytes(data []byte) int32 {
	if len(data) == 0 {
		return -1
	}
	return int32(C.JadeView_load_from_bytes(
		(*C.uint8_t)(unsafe.Pointer(&data[0])),
		C.size_t(len(data)),
	))
}

// IsLoaded 返回 JAPK 是否已加载。
func IsLoaded() bool {
	return C.JadeView_is_loaded() == 1
}

// GetAppSignature 返回当前 app_signature。
func GetAppSignature() string {
	return goStringFree(C.JadeView_get_app_signature())
}

// GetSignatureInfo 返回签名信息 JSON。
func GetSignatureInfo() string {
	return goStringFree(C.JadeView_get_signature_info())
}

// Unload 清除 JAPK 加载状态。返回 0=成功。
func Unload() int32 {
	return int32(C.JadeView_unload())
}
