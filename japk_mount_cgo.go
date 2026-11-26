//go:build linux && amd64 && jadeview_cgo && cgo

package jadeview

// #include "JadeView.h"
import "C"

func RegisterJAPK(path string) JAPKArchiveID {
	pool := &cstrPool{}
	defer pool.free()
	return JAPKArchiveID(C.register_japk(pool.s(path)))
}
func RegisterLoadedJAPK() JAPKArchiveID     { return JAPKArchiveID(C.register_loaded_japk()) }
func UnregisterJAPK(id JAPKArchiveID) int32 { return int32(C.unregister_japk(C.uint64_t(id))) }
func MountJAPK(id JAPKArchiveID, root string) JAPKMountID {
	pool := &cstrPool{}
	defer pool.free()
	return JAPKMountID(C.mount_japk(C.uint64_t(id), pool.s(root)))
}
func UnmountJAPK(id JAPKMountID) int32 { return int32(C.unmount_japk(C.uint64_t(id))) }
func GetJAPKMountRoot(id JAPKMountID) (string, bool) {
	s := C.get_japk_mount_root(C.uint64_t(id))
	return goStringFree(s), s != nil
}
func GetJAPKMountState(id JAPKMountID) (JAPKMountState, int32) {
	var out C.JapkMountState
	r := C.get_japk_mount_state(C.uint64_t(id), &out)
	return JAPKMountState{out.healthy != 0, uintptr(out.active_handles)}, int32(r)
}
