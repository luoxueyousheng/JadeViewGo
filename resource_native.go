//go:build windows || (linux && amd64 && !jadeview_cgo)

package jadeview

import (
	"runtime"
	"unsafe"
)

var (
	procResourceRegisterEx   = np("register_resource_ex")
	procResourceMode         = np("get_resource_delivery_mode")
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

// RegisterResourceEx 按指定模式登记文件；nil 使用受控模式，不自动切换到原生模式。
func RegisterResourceEx(path string, view, ttl uint32, options *ResourceOptions) (string, bool) {
	pool := &cstrs{}
	name := pool.p(path)
	var mode uint32
	var opt *uint32
	if options != nil {
		mode = uint32(options.Mode)
		opt = &mode
	}
	url, ok := bufCallSize(2048, func(buffer unsafe.Pointer, n uintptr) int32 {
		r, _, _ := procResourceRegisterEx.Call(uintptr(unsafe.Pointer(name)), uintptr(view), uintptr(ttl), uintptr(unsafe.Pointer(opt)), uintptr(buffer), n)
		return i32(r)
	})
	runtime.KeepAlive(pool)
	return url, ok
}

// GetResourceDeliveryMode 查询当前资源的实际交付模式。使用注册时返回的原始地址和所属 WebView ID。
func GetResourceDeliveryMode(url string, view uint32) (ResourceMode, int32) {
	pool := &cstrs{}
	var mode uint32
	r, _, _ := procResourceMode.Call(uintptr(unsafe.Pointer(pool.p(url))), uintptr(view), uintptr(unsafe.Pointer(&mode)))
	runtime.KeepAlive(pool)
	return ResourceMode(mode), i32(r)
}
