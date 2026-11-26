//go:build windows || (linux && amd64 && !jadeview_cgo)

package jadeview

import (
	"runtime"
	"unsafe"
)

var (
	procResourceCaps         = np("get_resource_capabilities")
	procResourceDelivery     = np("get_resource_delivery_info")
	procResourceDeliveryFree = np("free_resource_delivery_info")
	procJAPKMountCaps        = np("get_japk_mount_capabilities")
)

type cResourceDelivery struct {
	mode, reason uint32
	code         int32
	message      *byte
}

// GetResourceCapabilities 返回指定 WebView 的交付能力与原生状态码。
func GetResourceCapabilities(view uint32) (ResourceCapabilities, int32) {
	var caps ResourceCapabilities
	r, _, _ := procResourceCaps.Call(uintptr(view), uintptr(unsafe.Pointer(&caps)))
	return caps, i32(r)
}

// GetResourceDeliveryInfo 查询注册原因；url 须为 RegisterResource 返回的原始地址。
func GetResourceDeliveryInfo(url string, view uint32) (ResourceDeliveryInfo, int32) {
	var out cResourceDelivery
	pool := &cstrs{}
	r, _, _ := procResourceDelivery.Call(uintptr(unsafe.Pointer(pool.p(url))), uintptr(view), uintptr(unsafe.Pointer(&out)), unsafe.Sizeof(out))
	runtime.KeepAlive(pool)
	if i32(r) != 0 {
		return ResourceDeliveryInfo{}, i32(r)
	}
	defer procResourceDeliveryFree.Call(uintptr(unsafe.Pointer(&out)))
	return ResourceDeliveryInfo{ResourceMode(out.mode), out.reason, out.code, goString(uintptr(unsafe.Pointer(out.message)))}, 0
}

// GetJAPKMountCapabilities 返回本机原生挂载能力。
func GetJAPKMountCapabilities() (JAPKMountCapabilities, int32) {
	var caps JAPKMountCapabilities
	r, _, _ := procJAPKMountCaps.Call(uintptr(unsafe.Pointer(&caps)))
	return caps, i32(r)
}
