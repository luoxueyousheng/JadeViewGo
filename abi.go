package jadeview

import (
	"fmt"
	"strconv"
)

// ABIVersion 是此 SDK 使用的 JadeView 3.0.0-beta.4.26K01 原生 ABI。
const ABIVersion uint32 = 0x00030003

type abiInfo struct {
	version, pointerBits, windowOptionsSize, webviewSettingsSize uint32
}

func checkABI(info abiInfo, windowSize, webviewSize uintptr) error {
	if info.version != ABIVersion || info.pointerBits != strconv.IntSize ||
		info.windowOptionsSize != uint32(windowSize) || info.webviewSettingsSize != uint32(webviewSize) {
		return fmt.Errorf("JadeView ABI 不匹配：库 %#08x/%d 位，布局 %d/%d；Go 需要 %#08x/%d 位，布局 %d/%d",
			info.version, info.pointerBits, info.windowOptionsSize, info.webviewSettingsSize,
			ABIVersion, strconv.IntSize, windowSize, webviewSize)
	}
	return nil
}
