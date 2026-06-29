//go:build windows && amd64

package jadeview

// 自包含运行时：把 MSVC 版 JadeView.dll 用 go:embed 编进二进制，启动时释放到临时目录
// 并预加载。配合 jadeview_windows_amd64.go 的延迟加载，即可单 exe 分发、运行时自动释放，
// 无需随程序携带任何 DLL（MSVC 版 JadeView.dll 已自包含 WebView2Loader）。
//
// 原理：JadeView.dll 已改为延迟加载（首次调用其函数时才 LoadLibrary），故进程启动时
// 不需要它存在；本 init() 先把内置副本写到临时目录并按全路径预加载，之后延迟加载触发的
// LoadLibrary("JadeView.dll") 会命中这个同名已加载模块。
//
// 注意：目标机仍需安装 Edge WebView2 Runtime（系统级组件，任何方案都去不掉）。

import (
	_ "embed"
	"os"
	"path/filepath"
	"syscall"
)

//go:embed lib/windows_amd64/JadeView.dll
var jadeViewDLL []byte

func init() {
	dir := filepath.Join(os.TempDir(), "jadeview")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return
	}
	dst := filepath.Join(dir, "JadeView.dll")

	// 仅当不存在或大小不一致时才写，避免每次启动重复写 / 多进程争用已加载的文件。
	if fi, err := os.Stat(dst); err != nil || fi.Size() != int64(len(jadeViewDLL)) {
		_ = os.WriteFile(dst, jadeViewDLL, 0o644)
	}

	// 按全路径预加载；后续 delay-load 的 LoadLibrary("JadeView.dll") 会命中此同名模块。
	// 失败不致命（可能同目录/系统已有），留给实际调用处报错。
	_, _ = syscall.LoadLibrary(dst)
}
