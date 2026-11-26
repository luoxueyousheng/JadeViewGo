//go:build windows || (linux && amd64 && (!jadeview_cgo || cgo))

package jadeview

// ABIInfo 是原生库报告的 ABI 信息与关键结构尺寸。
type ABIInfo struct {
	Version, PointerBits, WindowOptionsSize, WebViewSettingsSize uint32
}

var loadedABI abiInfo // 由各平台的 Preload 在 sync.Once 中写入。

// GetABIInfo 加载库并返回经过校验的原生 ABI 信息。
func GetABIInfo() (ABIInfo, error) {
	if err := Preload(); err != nil {
		return ABIInfo{}, err
	}
	info := loadedABI
	return ABIInfo{info.version, info.pointerBits, info.windowOptionsSize, info.webviewSettingsSize}, nil
}
