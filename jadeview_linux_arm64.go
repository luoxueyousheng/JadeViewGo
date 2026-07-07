//go:build linux && arm64

package jadeview

// ARM64(aarch64)：静态链接 lib/arm64/libJadeView.a。
// 说明同 amd64 版本；交叉编译到 arm64 时需配置对应的交叉工具链与系统库。

// #cgo LDFLAGS: -L${SRCDIR}/lib/linux_arm64 -l:libJadeView.a
// #cgo pkg-config: gtk+-3.0 webkit2gtk-4.1
// #cgo LDFLAGS: -lpthread -ldl -lm
import "C"
