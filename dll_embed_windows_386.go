//go:build windows && 386

package jadeview

import _ "embed"

// x86(32位) 版 MSVC JadeView.dll，编进二进制、运行时释放。
//
//go:embed lib/windows_386/JadeView.dll
var embeddedJadeViewDLL []byte
