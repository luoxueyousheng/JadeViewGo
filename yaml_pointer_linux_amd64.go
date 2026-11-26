//go:build !jadeview_cgo || cgo

package jadeview

// YAMLGetStr 返回路径值的 JSON 文本，不存在时返回 false。
// Linux 使用功能等价的缓冲区接口，避免将 Windows COM 分配器规则用于 Linux。
func YAMLGetStr(fileName, keyPath string) (string, bool) { return YAMLGet(fileName, keyPath) }
