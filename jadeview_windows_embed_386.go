//go:build windows && 386

package jadeview

// 自包含运行时（x86）：把 MSVC 版 JadeView.dll 用 go:embed 编进二进制，启动时释放到临时目录
// 并预加载。配合 jadeview_windows_386.go 的延迟加载，单 exe 分发、运行时自动释放。
// 原理同 amd64（见 jadeview_windows_embed.go）。

import (
	_ "embed"
	"os"
	"path/filepath"
	"syscall"
)

//go:embed lib/windows_386/JadeView.dll
var jadeViewDLL []byte

func init() {
	dir := filepath.Join(os.TempDir(), "jadeview")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return
	}
	dst := filepath.Join(dir, "JadeView.dll")
	if fi, err := os.Stat(dst); err != nil || fi.Size() != int64(len(jadeViewDLL)) {
		_ = os.WriteFile(dst, jadeViewDLL, 0o644)
	}
	_, _ = syscall.LoadLibrary(dst)
}
