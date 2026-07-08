//go:build linux && amd64

package jadeview

// x86-64：静态链接 lib/x64/libJadeView.a。
//
// libJadeView.a 是 Rust staticlib，本身仍依赖一批系统动态库（GTK3 / WebKit2GTK
// 以及 pthread/dl/m 等）。下面用 pkg-config 拉取 GTK/WebKit 的链接参数；若你的发行版
// webkit 包名不是 webkit2gtk-4.1（老系统可能是 4.0），改这里即可。
//
// 注意：目录下同时存在 libJadeView.so，-lJadeView 会被 ld 优先解析成动态链接，
// 因此用 -l:libJadeView.a 显式指定静态库。静态库须排在系统库之前，保证符号正确解析。

// libJadeView.a 还引用了 libxdo(xdotool，X11 键盘/窗口自动化)的符号
// （xdo_new / xdo_send_keysequence_window / xdo_free 等），故补 -lxdo；
// 构建需装 libxdo-dev，运行需 libxdo3。
//
// #cgo LDFLAGS: -L${SRCDIR}/lib/linux_amd64 -l:libJadeView.a
// #cgo pkg-config: gtk+-3.0 webkit2gtk-4.1
// #cgo LDFLAGS: -lpthread -ldl -lm -lxdo
import "C"
