# 变更日志

格式参考 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.1.0/),版本号跟随上游 JadeView。

## [v2.4.3] - 2026-09-27

对应上游 JadeView **v2.4.3 (Build 26I01)**,刷新 Windows 三架构 DLL。已核对导出表:129 个公开 API 与 v2.4.0 官方头文件完全一致(另有 4 个内部辅助导出),API 无变化,Go 侧无需改动;头文件维持 2.4.0 版。Linux 库维持 v2.3.x(上游已停止 Linux 更新)。

### 修复

- **模块路径升级为 `github.com/luoxueyousheng/JadeViewGo/v2`**:此前 go.mod 缺少 `/v2` 后缀,导致全部 v2.x tag 按 Go modules 规则不可安装
  (`go get` 报 "module path must match major version"),`@latest` 也会错误解析到 v0.x。现导入路径改为
  `import jadeview "github.com/luoxueyousheng/JadeViewGo/v2"`,示例子包为 `.../JadeViewGo/v2/example`。

## [v2.4.0] - 2026-08-27

对应上游 JadeView **v2.4.0 (Build 26H03)**。Windows 三架构 DLL 与 2.4.0 官方头文件同步更新;Linux 库维持 v2.3.x(上游已停止 Linux 更新)。

### 破坏性变更

- **JAPK 只能加载带签名的资源包**:上游 JAPK 升级为平台根证书链 + v3 签名协议,`LoadFromBytes` 执行严格离线验签;
  混淆包(JPKBIN02)不再支持。
- **移除 `SetPublicKey`**:上游 v2.4.0 删除了 `JadeView_set_public_key`,公钥注入机制一并废弃(Ed25519 公钥现已内置长度校验,
  拒绝截断/拼接的非法公钥)。依赖 `SetPublicKey` 的宿主需改用签名分发流程。
- 上游同时移除的废弃 API:`cleanup_all_windows`(请用 `Exit`)。

### 新增

- **统一网页权限处理器**:`SetWebviewPermissionHandler` / `ClearWebviewPermissionHandler`(`PermissionHandler` 回调类型),
  支持摄像头、麦克风、录屏、文件访问等权限的允许/拒绝决策。
- **`ExitWait(timeoutMs)`**:清理所有窗口并等待事件循环与后台线程(GUI/回调/日志/托盘/单实例)完全退出,便于可靠卸载 DLL
  与重复初始化;`Exit` 内部亦改为等待运行时完成退出。
- **`WebViewSettings.ProfileName`**:Windows WebView2 Profile 名称,非空时多窗口的 Cookie、存储与缓存相互隔离。
- **导航历史 API**:`GoBack` / `GoForward` / `CanGoBack` / `CanGoForward`,
  对应头文件 `webview_go_back` / `webview_go_forward` / `webview_can_go_back` / `webview_can_go_forward`。
  至此头文件 129 个导出函数中已封装 128 个,仅 `yaml_get_str` 未封装(要求 `CoTaskMemFree` 释放,不可移植;
  `YAMLGet` 已用缓冲区两阶段查询等价实现)。

### 变更

- Windows 大消息 IPC 使用 WebView2 SharedBuffer 通道,替代旧分片拉取路径;`jade.invoke` 响应更快。
- IPC 协议在执行前严格校验 Origin 白名单;`postMessage` 白名单改为完整来源精确匹配;修复回调重绑定后的野指针崩溃、
  WebView 销毁重入与异常 URL 导致的 panic。

### 平台说明

- **Windows**:amd64 / 386 / arm64 DLL 已刷新为 2.4.0,构建与运行全部验证通过。
- **Linux**:上游已停止发布新版 Linux 库,`lib/linux_*` 维持 v2.3.x;本版本起一并移除了 Linux cgo 中引用已删符号的 `SetPublicKey`,Linux 构建暂不维护。

### 未封装说明

头文件 129 个导出函数中已封装 128 个。未封装:`yaml_get_str`(要求 `CoTaskMemFree` 释放,不可移植)。
DLL 中另有 4 个内部辅助导出(`gbk_to_utf8` 等)不属于公开 API。

## [v2.3.2] - 此前

对应上游 JadeView v2.3.2 (Build 26H01),刷新 Windows 三架构 DLL。详见该版本之前的提交历史。
