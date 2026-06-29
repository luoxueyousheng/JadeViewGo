//go:build windows && 386

package jadeview

// Windows x86(32位)：延迟链接 MSVC 版 JadeView.dll，配合 go:embed 单 exe 自包含。
// 方案同 amd64（见 jadeview_windows_amd64.go），库文件在 lib/windows_386/。
//
// 构建需 32 位 MinGW（MSYS2 的 mingw-w64-i686-toolchain，C:\msys64\mingw32\bin）：
//   GOARCH=386  CGO_ENABLED=1  CC=i686-w64-mingw32-gcc  go build ...

// #cgo CFLAGS: -I${SRCDIR}/include
// #cgo LDFLAGS: -L${SRCDIR}/lib/windows_386 -lJadeView_delay -ldelayimp
import "C"
