//go:build linux && amd64

package jadeview

// x86-64：静态链接 lib/x64/libJadeView.a。
//
// libJadeView.a 是 Rust staticlib，本身仍依赖一批系统动态库（GTK3 / WebKit2GTK
// 以及 pthread/dl/m 等）。下面用 pkg-config 拉取 GTK/WebKit 的链接参数；若你的发行版
// webkit 包名不是 webkit2gtk-4.1（老系统可能是 4.0），改这里即可。
//
// -lJadeView 必须排在系统库之前，保证静态符号正确解析。

// #cgo LDFLAGS: -L${SRCDIR}/lib/linux_amd64 -lJadeView
// #cgo pkg-config: gtk+-3.0 webkit2gtk-4.1
// #cgo LDFLAGS: -lpthread -ldl -lm
import "C"
