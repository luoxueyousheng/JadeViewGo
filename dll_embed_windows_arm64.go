//go:build windows && arm64

package jadeview

import _ "embed"

// ARM64 版 MSVC JadeView.dll，编进二进制、运行时释放。
// 纯 Go 方案下 arm64 与 x64/x86 一样支持单 exe 自包含（旧 cgo 延迟库方案做不到）。
//
//go:embed lib/windows_arm64/JadeView.dll
var embeddedJadeViewDLL []byte
