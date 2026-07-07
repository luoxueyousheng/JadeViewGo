//go:build windows && amd64

package jadeview

import _ "embed"

// x64 版 MSVC JadeView.dll（自包含 WebView2Loader），编进二进制、运行时释放。
//
//go:embed lib/windows_amd64/JadeView.dll
var embeddedJadeViewDLL []byte
