//go:build windows

package jadeview

import (
	"runtime"
	"unsafe"
)

// --- 剪贴板 ---

func ClipboardReadText() (string, bool) {
	return bufCallInt(4096, func(buf unsafe.Pointer, size int32) int32 {
		r, _, _ := procClipboardReadText.Call(uintptr(buf), uintptr(uint32(size)))
		return i32(r)
	})
}

func ClipboardWriteText(text string) bool {
	pool := &cstrs{}
	r, _, _ := procClipboardWriteText.Call(uintptr(unsafe.Pointer(pool.p(text))))
	runtime.KeepAlive(pool)
	return i32(r) == 1
}

// --- 路径 / 系统信息 ---

// GetPath 获取系统路径（name 如 home/appData/temp 等，具体见库文档）。
func GetPath(name string) (string, bool) {
	pool := &cstrs{}
	cn := pool.p(name)
	s, ok := bufCallSize(1024, func(buf unsafe.Pointer, size uintptr) int32 {
		r, _, _ := procGetPath.Call(uintptr(unsafe.Pointer(cn)), uintptr(buf), size)
		return i32(r)
	})
	runtime.KeepAlive(pool)
	return s, ok
}

// GetLocale 获取系统语言（BCP 47，如 zh-CN）。
func GetLocale() (string, bool) {
	return bufCallSize(64, func(buf unsafe.Pointer, size uintptr) int32 {
		r, _, _ := procGetLocale.Call(uintptr(buf), size)
		return i32(r)
	})
}

// GetDisplaysInfo 返回显示器信息 JSON 数组（含 bounds、work_area、scale_factor、dpi、is_primary）。
func GetDisplaysInfo() (string, bool) {
	return bufCallSize(8192, func(buf unsafe.Pointer, size uintptr) int32 {
		r, _, _ := procGetDisplaysInfo.Call(uintptr(buf), size)
		return i32(r)
	})
}

// GetCursorPosition 返回光标位置 JSON。
func GetCursorPosition() (string, bool) {
	return bufCallInt(128, func(buf unsafe.Pointer, size int32) int32 {
		r, _, _ := procGetCursorPosition.Call(uintptr(buf), uintptr(uint32(size)))
		return i32(r)
	})
}

// GetWebViewVersion 获取 WebView 内核版本号。
func GetWebViewVersion() (string, bool) {
	return bufCallSize(256, func(buf unsafe.Pointer, size uintptr) int32 {
		r, _, _ := procGetWebviewVersion.Call(uintptr(buf), size)
		return i32(r)
	})
}

// IsWindows11 检查当前系统是否为 Windows 11。
func IsWindows11() bool {
	r, _, _ := procIsWindows11.Call()
	return i32(r) == 1
}

// GetPrinterList 返回打印机名称 JSON 数组和打印机数量。
// 注意：jade_get_printer_list 返回的是打印机数量（非 1=成功），不能走 bufCallInt；
// 头文件未定义缓冲不足时的行为（16KB 对打印机名列表实际足够）。
func GetPrinterList() (string, int32, bool) {
	buf := make([]byte, 16384)
	r, _, _ := procGetPrinterList.Call(uintptr(unsafe.Pointer(&buf[0])), uintptr(uint32(len(buf))))
	rc := i32(r)
	if rc < 0 {
		return "", rc, false
	}
	return cBufToString(buf), rc, true
}

// --- 打印 ---

// Print 打开 WebView2 内置打印对话框。
func Print(windowID uint32) bool {
	r, _, _ := procJadePrint.Call(uintptr(windowID))
	return i32(r) == 1
}

// PrintFile 用系统关联程序打印文件（Windows=ShellExecute "print"，Linux=CUPS lp）。
func PrintFile(filePath string) bool {
	pool := &cstrs{}
	r, _, _ := procJadePrintDialog.Call(uintptr(unsafe.Pointer(pool.p(filePath))))
	runtime.KeepAlive(pool)
	return i32(r) == 1
}

// --- NTP 网络时间 ---

// NTPNow 获取网络时间戳（UTC 毫秒，北京时间需 +8 小时）。
// server 为空时使用内置服务器列表逐个尝试。失败返回 -1。
func NTPNow(server string) int64 {
	pool := &cstrs{}
	r1, r2, _ := procJadeNtpNow.Call(uintptr(unsafe.Pointer(pool.p(server))))
	runtime.KeepAlive(pool)
	if unsafe.Sizeof(uintptr(0)) == 4 {
		// 386：int64 返回值经 EDX:EAX 寄存器对（r2:r1）
		return int64(uint64(uint32(r1)) | uint64(uint32(r2))<<32)
	}
	return int64(r1)
}

// --- 全局热键 ---

// RegisterGlobalHotkey 注册全局热键，返回 hotkey_id（0=失败）。触发时通过 "global-hotkey" 事件回报。
func RegisterGlobalHotkey(modifiers, vk uint32) uint32 {
	r, _, _ := procRegisterGlobalHK.Call(uintptr(modifiers), uintptr(vk))
	return uint32(r)
}

func UnregisterGlobalHotkey(hotkeyID uint32) bool {
	r, _, _ := procUnregisterGlobalHK.Call(uintptr(hotkeyID))
	return i32(r) == 1
}

// --- 开机自启 ---

// SetLoginAutostart 启用/取消开机自启。args 为追加到可执行路径后的启动参数（可空）。
func SetLoginAutostart(enable bool, args string) bool {
	pool := &cstrs{}
	r, _, _ := procSetLoginAutostart.Call(b2u(enable), uintptr(unsafe.Pointer(pool.p(args))))
	runtime.KeepAlive(pool)
	return i32(r) == 1
}

func GetLoginAutostart() bool {
	r, _, _ := procGetLoginAutostart.Call()
	return i32(r) == 1
}

// --- URL 协议 / 文件关联 ---

func RegisterURLScheme(scheme string) bool {
	pool := &cstrs{}
	r, _, _ := procRegisterURLScheme.Call(uintptr(unsafe.Pointer(pool.p(scheme))))
	runtime.KeepAlive(pool)
	return i32(r) == 1
}

func UnregisterURLScheme(scheme string) bool {
	pool := &cstrs{}
	r, _, _ := procUnregisterURLScheme.Call(uintptr(unsafe.Pointer(pool.p(scheme))))
	runtime.KeepAlive(pool)
	return i32(r) == 1
}

func RegisterFileAssociation(extension, friendlyName string) bool {
	pool := &cstrs{}
	r, _, _ := procRegisterFileAssoc.Call(uintptr(unsafe.Pointer(pool.p(extension))),
		uintptr(unsafe.Pointer(pool.p(friendlyName))))
	runtime.KeepAlive(pool)
	return i32(r) == 1
}

func UnregisterFileAssociation(extension string) bool {
	pool := &cstrs{}
	r, _, _ := procUnregisterFileAssoc.Call(uintptr(unsafe.Pointer(pool.p(extension))))
	runtime.KeepAlive(pool)
	return i32(r) == 1
}

// --- 安全资源 / 协议服务 ---

// RegisterResource 注册本地文件为安全资源，返回 jade:// URL。
// windowID=0 表示全局；ttlSeconds=0 表示永不过期。
func RegisterResource(path string, windowID, ttlSeconds uint32) (string, bool) {
	pool := &cstrs{}
	cp := pool.p(path)
	s, ok := bufCallSize(512, func(buf unsafe.Pointer, size uintptr) int32 {
		r, _, _ := procRegisterResource.Call(uintptr(unsafe.Pointer(cp)),
			uintptr(windowID), uintptr(ttlSeconds), uintptr(buf), size)
		return i32(r)
	})
	runtime.KeepAlive(pool)
	return s, ok
}

func UnregisterResource(tokenOrURL string) bool {
	pool := &cstrs{}
	r, _, _ := procUnregisterResource.Call(uintptr(unsafe.Pointer(pool.p(tokenOrURL))))
	runtime.KeepAlive(pool)
	return i32(r) == 1
}

// ClearWindowResources 清理指定窗口的所有已注册资源，返回清理数量。
func ClearWindowResources(windowID uint32) int32 {
	r, _, _ := procClearWindowRes.Call(uintptr(windowID))
	return i32(r)
}

// GetFileIcon 提取任意路径的图标为 PNG 注册为 jade:// 资源，返回 URL。
// size 目标边长（16/32/48/64/128/256，<=0 取 48）；windowID=0 全局；ttlSeconds=0 默认。
func GetFileIcon(path string, size int, windowID, ttlSeconds uint32) (string, bool) {
	pool := &cstrs{}
	cp := pool.p(path)
	s, ok := bufCallSize(512, func(buf unsafe.Pointer, bufSize uintptr) int32 {
		r, _, _ := procGetFileIcon.Call(uintptr(unsafe.Pointer(cp)), uintptr(uint32(int32(size))),
			uintptr(windowID), uintptr(ttlSeconds), uintptr(buf), bufSize)
		return i32(r)
	})
	runtime.KeepAlive(pool)
	return s, ok
}

// SetProtocolServicePath 设置自定义协议服务路径，返回可访问的 URL（直接用于建窗导航）。
// rootPath 三种取值（beta.10 实测）：
//   - 磁盘目录路径：文件系统模式，服务该目录（hotReload 仅此模式有效，改文件即时刷新页面）；
//   - .japk 文件路径：挂载磁盘上的 JAPK 资源包；
//   - 特殊值 "japk"：内存模式，服务 LoadFromBytes 已加载的资源包
//     （须先加载成功，返回 URL 形如 JADE://<app_signature>）。
func SetProtocolServicePath(rootPath string, hotReload bool) (string, bool) {
	pool := &cstrs{}
	cr := pool.p(rootPath)
	s, ok := bufCallSize(512, func(buf unsafe.Pointer, size uintptr) int32 {
		r, _, _ := procSetProtocolSvcPath.Call(uintptr(unsafe.Pointer(cr)), uintptr(buf), size, b2u(hotReload))
		return i32(r)
	})
	runtime.KeepAlive(pool)
	return s, ok
}

// ClearDataDirectory 清空数据目录。confirm_token 必须等于 "I_UNDERSTAND_CLEAR_DATA"。
func ClearDataDirectory(confirmToken string) bool {
	pool := &cstrs{}
	r, _, _ := procClearDataDirectory.Call(uintptr(unsafe.Pointer(pool.p(confirmToken))))
	runtime.KeepAlive(pool)
	return i32(r) == 1
}

// --- 编码转换 ---

// SmartConvertEncoding 智能检测输入编码并转换为目标编码。
// 返回 (转换结果, 检测到的源编码, 是否成功)。targetEncoding 见库文档（utf-8/gbk/big5 等）。
func SmartConvertEncoding(input []byte, targetEncoding string) (string, string, bool) {
	pool := &cstrs{}
	ct := pool.p(targetEncoding)
	var inPtr unsafe.Pointer
	if len(input) > 0 {
		inPtr = unsafe.Pointer(&input[0])
	}

	outSize := len(input)*4 + 64
	detSize := 64
	for attempt := 0; attempt < 2; attempt++ {
		out := make([]byte, outSize)
		det := make([]byte, detSize)
		r, _, _ := procSmartConvertEnc.Call(
			uintptr(inPtr), uintptr(uint32(int32(len(input)))), uintptr(unsafe.Pointer(ct)),
			uintptr(unsafe.Pointer(&out[0])), uintptr(uint32(int32(outSize))),
			uintptr(unsafe.Pointer(&det[0])), uintptr(uint32(int32(detSize))),
		)
		rc := i32(r)
		if rc > 0 {
			runtime.KeepAlive(pool)
			runtime.KeepAlive(input)
			return string(out[:int(rc)]), cBufToString(det), true
		}
		if rc < 0 {
			// 缓冲区不足，绝对值为所需大小，放大重试
			outSize = int(-rc) + 16
			continue
		}
		break // 0 = 失败
	}
	runtime.KeepAlive(pool)
	runtime.KeepAlive(input)
	return "", "", false
}

// --- YAML 持久化存储（存于 JadeView_init 设置的 data_directory 下） ---
// 返回 int32 状态码：1=成功，0=路径/文件不存在，-1=IO，-2=类型不匹配，-3=已存在，-4=解析失败。
// getter 用缓冲区两阶段查询（避免 yaml_get_str 的 CoTaskMemFree 跨平台问题）。

// YAMLSet 设置路径值（自动解析 JSON/YAML/纯文本）。
func YAMLSet(fileName, keyPath, value string) int32 {
	pool := &cstrs{}
	r, _, _ := procYamlSet.Call(uintptr(unsafe.Pointer(pool.p(fileName))),
		uintptr(unsafe.Pointer(pool.p(keyPath))), uintptr(unsafe.Pointer(pool.p(value))))
	runtime.KeepAlive(pool)
	return i32(r)
}

// YAMLSetStr 强制按字符串存储（不解析 JSON/YAML）。
func YAMLSetStr(fileName, keyPath, value string) int32 {
	pool := &cstrs{}
	r, _, _ := procYamlSetStr.Call(uintptr(unsafe.Pointer(pool.p(fileName))),
		uintptr(unsafe.Pointer(pool.p(keyPath))), uintptr(unsafe.Pointer(pool.p(value))))
	runtime.KeepAlive(pool)
	return i32(r)
}

// YAMLGet 获取路径值，返回 JSON 字符串。第二返回值 true 表示成功。
func YAMLGet(fileName, keyPath string) (string, bool) {
	pool := &cstrs{}
	cf, ck := pool.p(fileName), pool.p(keyPath)
	s, rc := yamlGet(func(buf unsafe.Pointer, size uintptr) int32 {
		r, _, _ := procYamlGet.Call(uintptr(unsafe.Pointer(cf)), uintptr(unsafe.Pointer(ck)),
			uintptr(buf), size)
		return i32(r)
	})
	runtime.KeepAlive(pool)
	return s, rc == 1
}

// YAMLGetAll 读取整个文件，返回 JSON 字符串。
func YAMLGetAll(fileName string) (string, bool) {
	pool := &cstrs{}
	cf := pool.p(fileName)
	s, rc := yamlGet(func(buf unsafe.Pointer, size uintptr) int32 {
		r, _, _ := procYamlGetAll.Call(uintptr(unsafe.Pointer(cf)), uintptr(buf), size)
		return i32(r)
	})
	runtime.KeepAlive(pool)
	return s, rc == 1
}

// YAMLKeys 列出路径下所有 key，返回 JSON 数组字符串。keyPath 为空查询根节点。
func YAMLKeys(fileName, keyPath string) (string, bool) {
	pool := &cstrs{}
	cf, ck := pool.p(fileName), pool.p(keyPath)
	s, rc := yamlGet(func(buf unsafe.Pointer, size uintptr) int32 {
		r, _, _ := procYamlKeys.Call(uintptr(unsafe.Pointer(cf)), uintptr(unsafe.Pointer(ck)),
			uintptr(buf), size)
		return i32(r)
	})
	runtime.KeepAlive(pool)
	return s, rc == 1
}

// YAMLHas 检查路径是否存在。
func YAMLHas(fileName, keyPath string) bool {
	pool := &cstrs{}
	r, _, _ := procYamlHas.Call(uintptr(unsafe.Pointer(pool.p(fileName))),
		uintptr(unsafe.Pointer(pool.p(keyPath))))
	runtime.KeepAlive(pool)
	return i32(r) == 1
}

// YAMLDelete 删除指定路径。
func YAMLDelete(fileName, keyPath string) int32 {
	pool := &cstrs{}
	r, _, _ := procYamlDelete.Call(uintptr(unsafe.Pointer(pool.p(fileName))),
		uintptr(unsafe.Pointer(pool.p(keyPath))))
	runtime.KeepAlive(pool)
	return i32(r)
}

// YAMLLen 返回数组长度 / 对象 key 数。≥0=长度，-2=非映射非序列。keyPath 为空查询根节点。
func YAMLLen(fileName, keyPath string) int32 {
	pool := &cstrs{}
	r, _, _ := procYamlLen.Call(uintptr(unsafe.Pointer(pool.p(fileName))),
		uintptr(unsafe.Pointer(pool.p(keyPath))))
	runtime.KeepAlive(pool)
	return i32(r)
}

// YAMLClear 清空文件为 {}。
func YAMLClear(fileName string) bool {
	pool := &cstrs{}
	r, _, _ := procYamlClear.Call(uintptr(unsafe.Pointer(pool.p(fileName))))
	runtime.KeepAlive(pool)
	return i32(r) == 1
}

// YAMLDeleteFile 删除文件（含锁文件/临时文件）。
func YAMLDeleteFile(fileName string) int32 {
	pool := &cstrs{}
	r, _, _ := procYamlDeleteFile.Call(uintptr(unsafe.Pointer(pool.p(fileName))))
	runtime.KeepAlive(pool)
	return i32(r)
}

// --- JAPK 资源包：签名校验与加载 ---

// SetPublicKey 设置 Ed25519 公钥（Base64，44 字符），必须在加载 JAPK 之前调用。
// 返回 0=成功，负数=错误码。
func SetPublicKey(publicKey string) int32 {
	pool := &cstrs{}
	r, _, _ := procSetPublicKey.Call(uintptr(unsafe.Pointer(pool.p(publicKey))))
	runtime.KeepAlive(pool)
	return i32(r)
}

// LoadFromBytes 从内存加载 JAPK 文件（支持 v2 签名包与混淆包）。
//   - 已设公钥：必须是签名包，验签后加载；
//   - 未设公钥：仅支持混淆包（JPKBIN02）。
//
// app_name / app_signature 须与 Init 时一致。返回 0=成功，负数=错误码。
// 错误信息也会通过事件异步通知。
func LoadFromBytes(data []byte) int32 {
	if len(data) == 0 {
		return -1
	}
	r, _, _ := procLoadFromBytes.Call(uintptr(unsafe.Pointer(&data[0])), uintptr(len(data)))
	runtime.KeepAlive(data)
	return i32(r)
}

// IsLoaded 返回 JAPK 是否已加载。
func IsLoaded() bool {
	r, _, _ := procIsLoaded.Call()
	return i32(r) == 1
}

// GetAppSignature 返回当前 app_signature。
func GetAppSignature() string {
	r, _, _ := procGetAppSignature.Call()
	return goStringFree(r)
}

// GetSignatureInfo 返回签名信息 JSON。
func GetSignatureInfo() string {
	r, _, _ := procGetSignatureInf.Call()
	return goStringFree(r)
}

// Unload 清除 JAPK 加载状态。返回 0=成功。
func Unload() int32 {
	r, _, _ := procUnload.Call()
	return i32(r)
}
