//go:build linux && amd64 && jadeview_cgo && cgo

package jadeview

// #include "JadeView.h"
import "C"

func GetResourceCapabilities(view uint32) (ResourceCapabilities, int32) {
	var a C.ResourceCapabilities
	r := C.get_resource_capabilities(C.uint32_t(view), &a)
	return ResourceCapabilities{uint32(a.supported_modes), uint32(a.native_requires_private_view), uint32(a.native_requires_permanent), uint32(a.native_requires_dedicated_profile), uint32(a.native_registered), uint32(a.native_pending_cleanup)}, int32(r)
}

func GetResourceDeliveryInfo(url string, view uint32) (ResourceDeliveryInfo, int32) {
	pool := &cstrPool{}
	defer pool.free()
	var a C.ResourceDeliveryInfo
	r := C.get_resource_delivery_info(pool.s(url), C.uint32_t(view), &a, C.sizeof_ResourceDeliveryInfo)
	if r != 0 {
		return ResourceDeliveryInfo{}, int32(r)
	}
	defer C.free_resource_delivery_info(&a)
	return ResourceDeliveryInfo{ResourceMode(a.mode), uint32(a.reason), int32(a.native_error_code), C.GoString(a.message)}, 0
}

func GetJAPKMountCapabilities() (JAPKMountCapabilities, int32) {
	var a C.JapkMountCapabilities
	r := C.get_japk_mount_capabilities(&a)
	return JAPKMountCapabilities{uint32(a.backend), int32(a.available), int32(a.read_only), int32(a.caches_plaintext_to_disk), int32(a.original_archive_prefix)}, int32(r)
}

func RegisterResourceEx(path string, view, ttl uint32, options *ResourceOptions) (string, bool) {
	pool := &cstrPool{}
	defer pool.free()
	name := pool.s(path)
	var opt *C.ResourceOptions
	if options != nil {
		opt = &C.ResourceOptions{mode: C.uint32_t(options.Mode)}
	}
	return bufCallSize(2048, func(buffer *C.char, n C.size_t) C.int32_t {
		return C.register_resource_ex(name, C.uint32_t(view), C.uint32_t(ttl), opt, buffer, n)
	})
}

func GetResourceDeliveryMode(url string, view uint32) (ResourceMode, int32) {
	pool := &cstrPool{}
	defer pool.free()
	var mode C.uint32_t
	r := C.get_resource_delivery_mode(pool.s(url), C.uint32_t(view), &mode)
	return ResourceMode(mode), int32(r)
}
