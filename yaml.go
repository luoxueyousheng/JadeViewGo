package jadeview

/*
#include <stdlib.h>
#include "JadeView.h"
*/
import "C"

// YAML 持久化存储（存于 JadeView_init 设置的 data_directory 下）。
// 返回 int32 状态码：1=成功，0=路径/文件不存在，-1=IO，-2=类型不匹配，-3=已存在，-4=解析失败。
// getter 用缓冲区两阶段查询（避免 yaml_get_str 的 CoTaskMemFree 跨平台问题）。

// YAMLSet 设置路径值（自动解析 JSON/YAML/纯文本）。
func YAMLSet(fileName, keyPath, value string) int32 {
	pool := &cstrPool{}
	defer pool.free()
	return int32(C.yaml_set(pool.s(fileName), pool.s(keyPath), pool.s(value)))
}

// YAMLSetStr 强制按字符串存储（不解析 JSON/YAML）。
func YAMLSetStr(fileName, keyPath, value string) int32 {
	pool := &cstrPool{}
	defer pool.free()
	return int32(C.yaml_set_str(pool.s(fileName), pool.s(keyPath), pool.s(value)))
}

// YAMLGet 获取路径值，返回 JSON 字符串。第二返回值 true 表示成功。
func YAMLGet(fileName, keyPath string) (string, bool) {
	pool := &cstrPool{}
	defer pool.free()
	cf, ck := pool.s(fileName), pool.s(keyPath)
	s, rc := yamlGet(func(buf *C.char, n C.size_t) C.int32_t {
		return C.yaml_get(cf, ck, buf, n)
	})
	return s, rc == 1
}

// YAMLGetAll 读取整个文件，返回 JSON 字符串。
func YAMLGetAll(fileName string) (string, bool) {
	pool := &cstrPool{}
	defer pool.free()
	cf := pool.s(fileName)
	s, rc := yamlGet(func(buf *C.char, n C.size_t) C.int32_t {
		return C.yaml_get_all(cf, buf, n)
	})
	return s, rc == 1
}

// YAMLKeys 列出路径下所有 key，返回 JSON 数组字符串。keyPath 为空查询根节点。
func YAMLKeys(fileName, keyPath string) (string, bool) {
	pool := &cstrPool{}
	defer pool.free()
	cf, ck := pool.s(fileName), pool.s(keyPath)
	s, rc := yamlGet(func(buf *C.char, n C.size_t) C.int32_t {
		return C.yaml_keys(cf, ck, buf, n)
	})
	return s, rc == 1
}

// YAMLHas 检查路径是否存在。
func YAMLHas(fileName, keyPath string) bool {
	pool := &cstrPool{}
	defer pool.free()
	return C.yaml_has(pool.s(fileName), pool.s(keyPath)) == 1
}

// YAMLDelete 删除指定路径。
func YAMLDelete(fileName, keyPath string) int32 {
	pool := &cstrPool{}
	defer pool.free()
	return int32(C.yaml_delete(pool.s(fileName), pool.s(keyPath)))
}

// YAMLLen 返回数组长度 / 对象 key 数。≥0=长度，-2=非映射非序列。keyPath 为空查询根节点。
func YAMLLen(fileName, keyPath string) int32 {
	pool := &cstrPool{}
	defer pool.free()
	return int32(C.yaml_len(pool.s(fileName), pool.s(keyPath)))
}

// YAMLClear 清空文件为 {}。
func YAMLClear(fileName string) bool {
	pool := &cstrPool{}
	defer pool.free()
	return C.yaml_clear(pool.s(fileName)) == 1
}

// YAMLDeleteFile 删除文件（含锁文件/临时文件）。
func YAMLDeleteFile(fileName string) int32 {
	pool := &cstrPool{}
	defer pool.free()
	return int32(C.yaml_delete_file(pool.s(fileName)))
}
