//go:build windows && amd64

package jadeview

// Windows x64：延迟链接 MSVC 版 JadeView.dll，配合 go:embed 做到单 exe 自包含。
//
// 方案（弃用了早期 GNU 静态库）：
//   - 链接用 dlltool 从 DLL 导出表生成的「延迟导入库」libJadeView_delay.a，
//     使 JadeView.dll 变为延迟加载（首次调用某函数时才 LoadLibrary）。
//   - 因此进程启动时无需该 DLL 存在，jadeview_windows_embed.go 的 init() 便有机会
//     先把内置的 DLL 释放到临时目录并预加载。
//   - 这是纯动态链接，不碰 MSVC 静态库的 CRT 符号问题；MSVC 版 JadeView.dll 还
//     自包含了 WebView2Loader（无需额外携带 loader）。
//
// 换新版 DLL：把新的 MSVC JadeView.dll 覆盖到 lib/windows_amd64/，再按
// README「Windows 自包含」一节的 objdump+dlltool 命令重新生成 libJadeView_delay.a。
//
// 运行时仅需目标机安装 Edge WebView2 Runtime（系统级，Win11 自带）。
// GUI 程序（无控制台黑窗）构建加 -ldflags "-H windowsgui"。

// #cgo CFLAGS: -I${SRCDIR}/include
// #cgo LDFLAGS: -L${SRCDIR}/lib/windows_amd64 -lJadeView_delay -ldelayimp
import "C"
