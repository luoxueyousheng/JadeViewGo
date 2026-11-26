package jadeview

// 原生状态码。仅用于约定 0 表示成功的函数；其他函数可能使用 1/0 表示布尔值。
const (
	StatusOK               int32 = 0
	StatusInvalidArgument  int32 = -1
	StatusInvalidID        int32 = -2
	StatusWrongObjectType  int32 = -3
	StatusInvalidState     int32 = -4
	StatusQueueFull        int32 = -5
	StatusNotRunning       int32 = -6
	StatusPlatformError    int32 = -7
	StatusBusy             int32 = -8
	StatusInstanceReplaced int32 = -9
	StatusParentDestroyed  int32 = -10
	StatusIDExhausted      int32 = -11
	StatusWouldDeadlock    int32 = -12
	StatusDocumentReplaced int32 = -13
	StatusCancelled        int32 = -14
	StatusTimedOut         int32 = -15
)

const (
	ObjectReserved uint32 = iota + 1
	ObjectConfigured
	ObjectCreating
	ObjectReady
	ObjectClosing
	ObjectDestroying
	ObjectDestroyed
	ObjectFailed
	ObjectCrashed
)

const (
	ResourceCapControlled uint32 = 1 << iota
	ResourceCapNative
)

const (
	ResourceReasonExplicitControlled uint32 = iota
	ResourceReasonNativeSelected
	ResourceReasonFiniteTTL
	ResourceReasonGlobal
	ResourceReasonPlatform
	ResourceReasonNativeRejected
)

// ResourceOptions 指定交付模式：零值为受控交付。原生模式仅在支持的 Windows 环境可用。
type ResourceOptions struct{ Mode ResourceMode }

// ResourcePermanent 是不自动过期的资源 TTL。
const ResourcePermanent uint32 = ^uint32(0)

// JAPKArchiveID 标识已注册的签名磁盘包或内存包；0 表示注册失败。
type JAPKArchiveID uint64

// JAPKMountID 标识一个系统挂载；0 表示挂载失败。
type JAPKMountID uint64

// JAPKMountState 是挂载提供器的状态快照。
type JAPKMountState struct {
	Healthy       bool
	ActiveHandles uintptr
}
