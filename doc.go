// Package jadeview 是 JadeView 3.0.0-beta.4.26K01 的 Go 封装。
//
// Windows（amd64、386、arm64）通过纯 Go 加载内置 DLL；Linux（amd64）
// 默认通过纯 Go 动态加载 libJadeView.so，也可显式选择旧 cgo 后端。
// 宿主可在启动时调用 Preload 校验 ABI。
//
// CreateWindow 返回 Window ID；调用 DefaultWebView 获取用于页面操作和 IPC 的 WebView ID。
package jadeview
