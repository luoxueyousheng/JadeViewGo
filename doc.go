// Package jadeview 是 JadeView WebView 库的 Go 封装。
//
// 公共 API 跨平台一致，内部实现按平台分开：
//
//   - Windows（amd64/386/arm64）：纯 Go 实现（*_windows.go），syscall 直调
//     go:embed 内置的 JadeView.dll，构建无需 C 编译器/cgo；
//   - Linux（amd64/arm64）：cgo 静态链接 libJadeView.a（原方式，构建需
//     gcc + GTK3/WebKit2GTK 开发包）。
//
// API 按模块拆分：生命周期(jadeview*.go)、窗口(window*.go)、事件(events*.go)、
// 对话框(dialog*.go)、托盘(tray*.go)、YAML(yaml*.go)、系统工具(system*.go)、
// JAPK(japk*.go)；共享类型见 types.go。
package jadeview
