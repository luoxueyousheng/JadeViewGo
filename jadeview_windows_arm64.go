//go:build windows && arm64

package jadeview

// Windows ARM64：常规加载期链接 MSVC 版 JadeView.dll（直接对 DLL 链接，无需导入库）。
//
// 与 amd64/386 的区别：**没有延迟加载 + go:embed 自包含**。
// binutils 的 dlltool 至今没有 aarch64-PE 的延迟库支持（-m arm64 生成 delaylib
// 会静默失败），llvm-dlltool 也不支持 --output-delaylib。因此 arm64 走普通
// 加载期导入：分发时 exe 旁边必须携带 lib/windows_arm64/JadeView.dll。
//
// 构建需要 llvm-mingw 工具链（https://github.com/mstorsjo/llvm-mingw，
// 下载 x86_64 hosted 版即可在 x64 机器上交叉编译）：
//
//	$env:PATH = "C:\llvm-mingw\bin;" + $env:PATH
//	$env:CGO_ENABLED="1"; $env:GOARCH="arm64"
//	$env:CC="aarch64-w64-mingw32-clang"
//	go build -ldflags "-H windowsgui" .
//
// -lJadeView 在只有 JadeView.dll 的目录下会由链接器（lld/GNU ld 均支持）
// 直接解析 DLL 导出表完成链接；lib/windows_arm64/JadeView.def 是导出表备份，
// 如需生成传统导入库可用 llvm-dlltool：
//
//	llvm-dlltool -m arm64 -d JadeView.def -D JadeView.dll -l libJadeView.dll.a

// #cgo CFLAGS: -I${SRCDIR}/include
// #cgo LDFLAGS: -L${SRCDIR}/lib/windows_arm64 -lJadeView
import "C"
