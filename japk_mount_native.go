//go:build windows || (linux && amd64 && !jadeview_cgo)

package jadeview

import (
	"runtime"
	"unsafe"
)

var (
	procRegisterJAPK       = np("register_japk")
	procRegisterLoadedJAPK = np("register_loaded_japk")
	procUnregisterJAPK     = np("unregister_japk")
	procMountJAPK          = np("mount_japk")
	procUnmountJAPK        = np("unmount_japk")
	procJAPKMountRoot      = np("get_japk_mount_root")
	procJAPKMountState     = np("get_japk_mount_state")
)

// RegisterJAPK 校验签名磁盘包并登记；返回 0 表示失败。
func RegisterJAPK(path string) JAPKArchiveID {
	pool := &cstrs{}
	lo, hi, _ := procRegisterJAPK.Call(uintptr(unsafe.Pointer(pool.p(path))))
	runtime.KeepAlive(pool)
	return JAPKArchiveID(requestID(lo, hi))
}

// RegisterLoadedJAPK 登记已用 LoadFromBytes 成功载入的当前内存包。
func RegisterLoadedJAPK() JAPKArchiveID {
	lo, hi, _ := procRegisterLoadedJAPK.Call()
	return JAPKArchiveID(requestID(lo, hi))
}

// UnregisterJAPK 停止使用包接受新资源；不会关闭已有文件句柄或代替卸载。
func UnregisterJAPK(id JAPKArchiveID) int32 {
	r, _, _ := callU64First(procUnregisterJAPK, uint64(id))
	return i32(r)
}

// MountJAPK 创建系统挂载。磁盘包可传空 root 使用默认路径；内存包必须指定 root。
func MountJAPK(id JAPKArchiveID, root string) JAPKMountID {
	pool := &cstrs{}
	lo, hi, _ := callU64First(procMountJAPK, uint64(id), uintptr(unsafe.Pointer(pool.p(root))))
	runtime.KeepAlive(pool)
	return JAPKMountID(requestID(lo, hi))
}

// UnmountJAPK 尝试卸载。返回 StatusBusy 时应先关闭文件、目录迭代器和映射再重试。
func UnmountJAPK(id JAPKMountID) int32 {
	r, _, _ := callU64First(procUnmountJAPK, uint64(id))
	return i32(r)
}

// GetJAPKMountRoot 返回实际挂载路径，第二个值表示是否成功。
func GetJAPKMountRoot(id JAPKMountID) (string, bool) {
	r, _, _ := callU64First(procJAPKMountRoot, uint64(id))
	return goStringFree(r), r != 0
}

type cJAPKMountState struct {
	healthy int32
	active  uintptr
}

// GetJAPKMountState 返回挂载快照和原生状态码。
func GetJAPKMountState(id JAPKMountID) (JAPKMountState, int32) {
	var out cJAPKMountState
	r, _, _ := callU64First(procJAPKMountState, uint64(id), uintptr(unsafe.Pointer(&out)))
	return JAPKMountState{out.healthy != 0, out.active}, i32(r)
}
