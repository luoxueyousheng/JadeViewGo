#ifndef JADEVIEW_H
#define JADEVIEW_H

#include <stdarg.h>
#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
#include <stdlib.h>

/**
 * WebView设置选项结构体，C-compatible
 * 用于设置webview的高级选项
 */
typedef struct WebViewSettings {
  /**
   * Allocated structure size in bytes; initialize with the sized default function.
   */
  uint32_t struct_size;
  /**
   * JADE_ABI_VERSION; checked before reading any configuration fields.
   */
  uint32_t abi_version;
  /**
   * 是否允许自动播放媒体 (0=false, 1=true)
   */
  int32_t autoplay;
  /**
   * 后台定时器：0/2=引擎默认，1=禁用后台定时器节流；Windows 使用启动参数，Linux 不支持非零设置
   */
  int32_t background_throttling;
  /**
   * 是否允许页面右键菜单 (0=禁用, 1=允许)
   */
  int32_t allow_right_click;
  /**
   * 自定义User-Agent字符串指针，NULL表示使用默认UA
   */
  const char *user_agent;
  /**
   * 预载JavaScript代码字符串指针，NULL表示不预载JS
   */
  const char *preload_js;
  /**
   * 是否允许页面全屏 (0=false, 1=true)
   */
  int32_t allow_fullscreen;
  /**
   * PostMessage 允许来源，逗号或分号分隔；* 允许所有来源
   */
  const char *postmessage_allowed_origins;
  /**
   * CORS 允许来源，逗号或分号分隔；* 允许所有来源
   */
  const char *cors_allowed_origins;
  /**
   * 是否启用账号/密码自动填充 (0=禁用, 1=启用)
   */
  int32_t password_autofill;
  /**
   * 是否启用通用表单自动填充 (0=禁用, 1=启用)
   */
  int32_t form_autofill;
  /**
   * 无痕模式 (0=否, 1=是)
   */
  int32_t incognito;
  /**
   * 禁用剪贴板读写权限 (0=启用, 1=禁用)
   */
  int32_t disable_clipboard;
  /**
   * 代理URL字符串指针(NULL=不使用代理)，格式如 "http://host:port" 或 "socks5://host:port"
   */
  const char *proxy_url;
  /**
   * 初始是否自动获取焦点 (0=否, 1=是)
   */
  int32_t focused;
  /**
   * Profile 名称；Linux 命名 Profile 独占视图并保持固定 HTTPS 来源，未命名使用临时隔离。
   */
  const char *profile_name;
  /**
   * 开发者工具：0=跟随全局 devmod，1=启用，2=禁用。
   */
  int32_t devtools;
  /**
   * 页面初始缩放倍数；0 采用默认 1.0。
   */
  double zoom_factor;
  /**
   * WebView 背景色 #RRGGBB 或 #RRGGBBAA；NULL 继承窗口背景。
   */
  const char *background_color;
} WebViewSettings;

/**
 * WebView窗口选项结构体，C-compatible
 * 宿主绑定必须与当前头文件的完整字段布局一致。
 */
typedef struct WebViewWindowOptions {
  /**
   * Allocated structure size in bytes; initialize with the sized default function.
   */
  uint32_t struct_size;
  /**
   * JADE_ABI_VERSION; checked before reading any configuration fields.
   */
  uint32_t abi_version;
  /**
   * 标题
   */
  const char *title;
  /**
   * 宽度
   */
  int32_t width;
  /**
   * 高度
   */
  int32_t height;
  /**
   * 可调整大小边框
   */
  int32_t resizable;
  /**
   * 无边框：1=使用现有无边框实现，0=保留系统边框；默认 0。
   */
  int32_t noframe;
  /**
   * default=系统标题栏，hidden=隐藏标题栏，BarOverlay=标题栏按钮覆盖层；NULL 等同 default。
   */
  const char *title_bar_style;
  /**
   * 透明背景
   */
  int32_t transparent;
  /**
   * 背景颜色："#RRGGBBAA" 十六进制字符串，例如 "#ffffff00"
   */
  const char *background_color;
  /**
   * 置顶窗口
   */
  int32_t always_on_top;
  /**
   * 窗口主题
   */
  const char *theme;
  /**
   * 最大化
   */
  int32_t maximized;
  /**
   * 最大化按钮
   */
  int32_t maximizable;
  /**
   * 最小化按钮
   */
  int32_t minimizable;
  /**
   * X坐标（-1 表示居中）
   */
  int32_t x;
  /**
   * Y坐标（-1 表示居中）
   */
  int32_t y;
  /**
   * 最小宽度
   */
  int32_t min_width;
  /**
   * 最小高度
   */
  int32_t min_height;
  /**
   * 最大宽度
   */
  int32_t max_width;
  /**
   * 最大高度
   */
  int32_t max_height;
  /**
   * 全屏
   */
  int32_t fullscreen;
  /**
   * 焦点
   */
  int32_t focus;
  /**
   * 隐藏窗口
   */
  int32_t hide_window;
  /**
   * 使用页面图标
   */
  int32_t use_page_icon;
  /**
   * 不进任务栏和 Alt-Tab（0=否，1=是）。
   */
  int32_t skip_taskbar;
  /**
   * 点击或显示时不激活窗口（0=否，1=是）。
   */
  int32_t no_activate;
  /**
   * 创建时背景材质：default、mica、micaAlt 或 acrylic；NULL 为默认背景。
   */
  const char *background_material;
  /**
   * Windows 默认标题栏与 BarOverlay 共用的返回按钮开关；0=关闭，1=开启，默认 1。
   */
  int32_t titlebar_back;
  /**
   * Windows 默认标题栏与 BarOverlay 共用的刷新按钮开关；0=关闭，1=开启，默认 1。
   */
  int32_t titlebar_refresh;
  /**
   * Windows 默认标题栏的窗口图标显示开关；0=隐藏，1=显示，默认 1。
   */
  int32_t titlebar_icon;
} WebViewWindowOptions;

/**
 * 通知参数结构体
 */
typedef struct NotificationParams {
  /**
   * 通知摘要/标题
   */
  const char *summary;
  /**
   * 通知正文
   */
  const char *body;
  /**
   * 通知图标路径
   */
  const char *icon;
  /**
   * 通知超时时间（毫秒，-1=系统默认）
   */
  int32_t timeout;
  /**
   * 按钮1文本
   */
  const char *button1;
  /**
   * 按钮2文本
   */
  const char *button2;
  /**
   * 按钮3文本
   */
  const char *text3;
  /**
   * 默认操作动作
   */
  const char *action;
} NotificationParams;

/**
 * 文件对话框参数结构体
 */
typedef struct FileDialogParams {
  /**
   * 关联窗口ID
   */
  uint32_t window_id;
  /**
   * 对话框标题
   */
  const char *title;
  /**
   * 默认路径
   */
  const char *default_path;
  /**
   * 确认按钮文本
   */
  const char *button_label;
  /**
   * 文件过滤器，JSON格式
   */
  const char *filters;
  /**
   * 对话框属性，JSON格式
   */
  const char *properties;
} FileDialogParams;

/**
 * 消息框参数结构体
 */
typedef struct MessageBoxParams {
  /**
   * 关联窗口ID
   */
  uint32_t window_id;
  /**
   * 消息框标题
   */
  const char *title;
  /**
   * 消息正文
   */
  const char *message;
  /**
   * 详细信息
   */
  const char *detail;
  /**
   * 按钮列表，JSON格式
   */
  const char *buttons;
  /**
   * 默认按钮索引
   */
  int32_t default_id;
  /**
   * 取消按钮索引
   */
  int32_t cancel_id;
  /**
   * 消息框类型：none / info / warning / error
   */
  const char *type_;
} MessageBoxParams;

typedef struct WebViewBounds {
  const char *x;
  const char *y;
  const char *width;
  const char *height;
} WebViewBounds;

typedef struct WebViewIdArray {
  const uint32_t *ids;
  size_t count;
} WebViewIdArray;

typedef struct JadeObjectState {
  uint32_t kind;
  uint32_t state;
  uint32_t parent_window_id;
  uint32_t default_webview_id;
  int32_t error_code;
} JadeObjectState;

typedef struct JadeRequestResult {
  uint32_t state;
  int32_t error_code;
  uint32_t object_id;
  char *message;
} JadeRequestResult;

typedef struct LogicalRect {
  double x;
  double y;
  double width;
  double height;
} LogicalRect;

typedef struct PhysicalRect {
  int32_t x;
  int32_t y;
  uint32_t width;
  uint32_t height;
} PhysicalRect;

typedef struct WebViewBoundsResult {
  struct WebViewBounds bounds;
  struct LogicalRect logical;
  struct PhysicalRect physical;
} WebViewBoundsResult;

typedef struct JapkMountCapabilities {
  uint32_t backend;
  int32_t available;
  int32_t read_only;
  int32_t caches_plaintext_to_disk;
  int32_t original_archive_prefix;
} JapkMountCapabilities;

/**
 * register_resource_ex options. mode: 0=controlled (default), 1=native file.
 */
typedef struct ResourceOptions {
  uint32_t mode;
} ResourceOptions;

/**
 * Per-WebView engine/policy capabilities; individual files can still be refused.
 */
typedef struct ResourceCapabilities {
  /**
   * Bit 1=controlled; bit 2=native currently available for this Ready WebView.
   */
  uint32_t supported_modes;
  uint32_t native_requires_private_view;
  uint32_t native_requires_permanent;
  uint32_t native_requires_dedicated_profile;
  /**
   * Registrations belonging to this WebView.
   */
  uint32_t native_registered;
  /**
   * Global retired filesystem backings awaiting browser handle release.
   */
  uint32_t native_pending_cleanup;
} ResourceCapabilities;

/**
 * Frozen output layout. message is library-owned; release with free_resource_delivery_info.
 */
typedef struct ResourceDeliveryInfo {
  uint32_t mode;
  uint32_t reason;
  int32_t native_error_code;
  char *message;
} ResourceDeliveryInfo;

typedef struct JadeAbiInfo {
  uint32_t abi_version;
  uint32_t pointer_bits;
  uint32_t window_options_size;
  uint32_t webview_settings_size;
} JadeAbiInfo;

typedef struct JapkMountState {
  int32_t healthy;
  size_t active_handles;
} JapkMountState;

typedef struct JadeEvent {
  uint32_t event_type;
  uint32_t source_kind;
  uint64_t source_id;
  uint32_t window_id;
  uint32_t webview_id;
  uint64_t sequence;
  uint64_t request_id;
  uint32_t flags;
  uint32_t state;
  int32_t error_code;
  const char *name;
  const uint8_t *data;
  size_t data_size;
} JadeEvent;

typedef struct JadeMenuIconDesc {
  /**
   * PNG/ICO path, or TTF/OTF/TTC/OTC path with a nonzero `icon_codepoint`.
   * Mounted JAPK paths work for both.
   */
  const char *icon_path;
  /**
   * PNG/ICO bytes. Copied during the call and used before the image path.
   */
  const uint8_t *icon_data;
  /**
   * Byte length of `icon_data`; at most 8 MiB.
   */
  size_t icon_data_size;
  /**
   * Glyph codepoint in the supplied font, or the default icon font if no font path is supplied.
   * Nonzero takes precedence over image data.
   * Icon list: https://learn.microsoft.com/zh-cn/windows/apps/design/iconography/segoe-fluent-icons-font
   */
  uint32_t icon_codepoint;
} JadeMenuIconDesc;

typedef struct JadeMenuItemDesc {
  const char *id;
  const char *parent_id;
  const char *label;
  const char *accelerator;
  /**
   * 0 command, 1 submenu, 2 separator, 3 checkbox, 4 radio.
   */
  uint32_t kind;
  /**
   * 1 disables the item, 2 marks a checkbox or radio item as checked.
   */
  uint32_t flags;
  /**
   * Optional item icon; zero-initialize for no icon.
   */
  struct JadeMenuIconDesc icon;
} JadeMenuItemDesc;

/**
 * WebView page context-menu template item. `kind` keeps the same meaning as
 * `JadeMenuItemDesc`; `role` optionally binds the item to a unified action.
 */
typedef struct JadeContextMenuItemDesc {
  uint32_t struct_size;
  const char *id;
  const char *parent_id;
  const char *label;
  const char *accelerator;
  /**
   * Unified role (e.g. "copy", "paste", "print", "saveImageAs").
   */
  const char *role;
  /**
   * 0 command, 1 submenu, 2 separator, 3 checkbox, 4 radio.
   */
  uint32_t kind;
  /**
   * 1 disables the item, 2 marks a checkbox or radio item as checked.
   */
  uint32_t flags;
  struct JadeMenuIconDesc icon;
} JadeContextMenuItemDesc;

typedef struct JadeMenuPosition {
  int32_t x;
  int32_t y;
} JadeMenuPosition;

typedef struct JadeSplashOptions {
  uint32_t struct_size;
  uint32_t abi_version;
  /**
   * Bit 0 disables the splash.
   */
  uint32_t flags;
  const char *logo_path;
  const uint8_t *logo_data;
  size_t logo_data_size;
  /**
   * 0xAARRGGBB.
   */
  uint32_t background_color;
  uint32_t accent_color;
  uint32_t minimum_ms;
} JadeSplashOptions;

/**
 * Optional window features. Zero disables each feature.
 */
typedef struct JadeWindowFeatures {
  uint32_t struct_size;
  uint32_t abi_version;
  /**
   * Enable window content protection.
   */
  int32_t content_protection;
  /**
   * Save and restore the window position.
   */
  int32_t auto_save_state;
  /**
   * Optional startup splash for this window.
   */
  const struct JadeSplashOptions *splash;
} JadeWindowFeatures;


// -----------------------
// 事件名常量（自动从源码提取）
// -----------------------
// 应用初始化完成
#define JADEVIEW_EVENT_APP_READY "app-ready"
// 程序崩溃
#define JADEVIEW_EVENT_CRASH "crash"
// 拖拽事件（enter/over/drop/leave）
#define JADEVIEW_EVENT_DRAG_DROP "drag-drop"
// 全局热键触发
#define JADEVIEW_EVENT_GLOBAL_HOTKEY "global-hotkey"
// JAPK 资源包加载失败
#define JADEVIEW_EVENT_JAPK_LOAD_FAILED "japk-load-failed"
// 菜单命令
#define JADEVIEW_EVENT_MENU_COMMAND "menu-command"
// 菜单显示错误
#define JADEVIEW_EVENT_MENU_ERROR "menu-error"
// 通知操作被点击
#define JADEVIEW_EVENT_NOTIFICATION_ACTION "notification-action"
// 通知已关闭
#define JADEVIEW_EVENT_NOTIFICATION_DISMISSED "notification-dismissed"
// 通知显示失败
#define JADEVIEW_EVENT_NOTIFICATION_FAILED "notification-failed"
// 通知已显示
#define JADEVIEW_EVENT_NOTIFICATION_SHOWN "notification-shown"
// 收到前端 postMessage
#define JADEVIEW_EVENT_POSTMESSAGE_RECEIVED "postmessage-received"
// 第二个实例启动
#define JADEVIEW_EVENT_SECOND_INSTANCE "second-instance"
// 系统主题变化
#define JADEVIEW_EVENT_THEME_CHANGED "theme-changed"
// 托盘图标事件
#define JADEVIEW_EVENT_TRAY_EVENT "tray-event"
// 更新窗口图标
#define JADEVIEW_EVENT_UPDATE_WINDOW_ICON "update-window-icon"
// WebView 加载完成
#define JADEVIEW_EVENT_WEBVIEW_DID_FINISH_LOAD "webview-did-finish-load"
// WebView 开始加载
#define JADEVIEW_EVENT_WEBVIEW_DID_START_LOADING "webview-did-start-loading"
#define JADEVIEW_EVENT_WEBVIEW_DOWNLOAD_COMPLETED "webview-download-completed"
// WebView 页面图标更新
#define JADEVIEW_EVENT_WEBVIEW_PAGE_FAVICON_UPDATED "webview-page-favicon-updated"
// WebView 页面标题更新
#define JADEVIEW_EVENT_WEBVIEW_PAGE_TITLE_UPDATED "webview-page-title-updated"
// 所有窗口已关闭
#define JADEVIEW_EVENT_WINDOW_ALL_CLOSED "window-all-closed"
// 窗口失去焦点
#define JADEVIEW_EVENT_WINDOW_BLURRED "window-blurred"
#define JADEVIEW_EVENT_WINDOW_BOUNDS "window-bounds"
// 窗口关闭
#define JADEVIEW_EVENT_WINDOW_CLOSED "window-closed"
// 窗口销毁
#define JADEVIEW_EVENT_WINDOW_DESTROYED "window-destroyed"
// 窗口获得焦点
#define JADEVIEW_EVENT_WINDOW_FOCUSED "window-focused"
// 窗口全屏状态变化
#define JADEVIEW_EVENT_WINDOW_FULLSCREEN "window-fullscreen"
// 窗口移动
#define JADEVIEW_EVENT_WINDOW_MOVED "window-moved"
// 窗口大小改变
#define JADEVIEW_EVENT_WINDOW_RESIZED "window-resized"
// 窗口状态变化（最大化/还原等）
#define JADEVIEW_EVENT_WINDOW_STATE_CHANGED "window-state-changed"


// -----------------------
// 崩溃错误代码（crash 事件 event_data）
// -----------------------
// SEH 异常：内存访问违规
#define JADEVIEW_CRASH_SEH_ACCESS_VIOLATION "SEH_ACCESS_VIOLATION"
// SEH 异常：栈溢出
#define JADEVIEW_CRASH_SEH_STACK_OVERFLOW "SEH_STACK_OVERFLOW"
// SEH 异常：非法指令
#define JADEVIEW_CRASH_SEH_ILLEGAL_INSTRUCTION "SEH_ILLEGAL_INSTRUCTION"
// SEH 异常：无效句柄
#define JADEVIEW_CRASH_SEH_INVALID_HANDLE "SEH_INVALID_HANDLE"
// SEH 异常：未知
#define JADEVIEW_CRASH_SEH_UNKNOWN "SEH_UNKNOWN"
// 运行时异常：Panic
#define JADEVIEW_CRASH_RUNTIME_PANIC "RUNTIME_PANIC"
// WebView2：浏览器进程已退出
#define JADEVIEW_CRASH_WV2_BROWSER_EXITED "WV2_BROWSER_EXITED"
// WebView2：渲染进程已退出
#define JADEVIEW_CRASH_WV2_RENDER_EXITED "WV2_RENDER_EXITED"
// WebView2：渲染进程无响应
#define JADEVIEW_CRASH_WV2_RENDER_UNRESPONSIVE "WV2_RENDER_UNRESPONSIVE"
// WebView2：框架渲染进程已退出
#define JADEVIEW_CRASH_WV2_FRAME_RENDER_EXITED "WV2_FRAME_RENDER_EXITED"
// WebView2：工具进程已退出
#define JADEVIEW_CRASH_WV2_UTILITY_EXITED "WV2_UTILITY_EXITED"
// WebView2：沙箱辅助进程已退出
#define JADEVIEW_CRASH_WV2_SANDBOX_HELPER_EXITED "WV2_SANDBOX_HELPER_EXITED"
// WebView2：GPU进程已退出
#define JADEVIEW_CRASH_WV2_GPU_EXITED "WV2_GPU_EXITED"
// WebView2：PPAPI插件进程已退出
#define JADEVIEW_CRASH_WV2_PPAPI_PLUGIN_EXITED "WV2_PPAPI_PLUGIN_EXITED"
// WebView2：PPAPI代理进程已退出
#define JADEVIEW_CRASH_WV2_PPAPI_BROKER_EXITED "WV2_PPAPI_BROKER_EXITED"
// WebView2：未知进程已退出
#define JADEVIEW_CRASH_WV2_UNKNOWN_EXITED "WV2_UNKNOWN_EXITED"



// -----------------------
// C API 函数声明（自动从 src/ffi/ 生成）
// -----------------------

#ifndef JADEVIEW_CALL
#ifdef _WIN32
#define JADEVIEW_CALL __stdcall
#else
#define JADEVIEW_CALL
#endif
#endif

// 回调函数类型定义
typedef const char* (JADEVIEW_CALL *IpcCallback)(uint32_t, const char*);
typedef int32_t (JADEVIEW_CALL *PermissionCallback)(uint32_t, const char*);
typedef void (JADEVIEW_CALL *JadeEventCallback)(const struct JadeEvent*, void*);

#define JADE_OK 0
#define JADE_INVALID_ARGUMENT -1
#define JADE_INVALID_ID -2
#define JADE_WRONG_OBJECT_TYPE -3
#define JADE_INVALID_STATE -4
#define JADE_QUEUE_FULL -5
#define JADE_NOT_RUNNING -6
#define JADE_PLATFORM_ERROR -7
#define JADE_BUSY -8
#define JADE_INSTANCE_REPLACED -9
#define JADE_PARENT_DESTROYED -10
#define JADE_ID_EXHAUSTED -11
#define JADE_WOULD_DEADLOCK -12
#define JADE_DOCUMENT_REPLACED -13
#define JADE_CANCELLED -14
#define JADE_TIMED_OUT -15
#define JADE_ABI_VERSION 0x00030003
#define JADE_SPLASH_DISABLED 1
#define JADE_RESOURCE_REASON_EXPLICIT_CONTROLLED 0
#define JADE_RESOURCE_REASON_NATIVE_SELECTED 1
#define JADE_RESOURCE_REASON_FINITE_TTL 2
#define JADE_RESOURCE_REASON_GLOBAL 3
#define JADE_RESOURCE_REASON_PLATFORM 4
#define JADE_RESOURCE_REASON_NATIVE_REJECTED 5
#define JADE_RESOURCE_CONTROLLED 0
#define JADE_RESOURCE_NATIVE 1
#define JADE_RESOURCE_CAP_CONTROLLED 1
#define JADE_RESOURCE_CAP_NATIVE 2
#define JADE_EVENT_ALL 0
#define JADE_EVENT_CREATED 1
#define JADE_EVENT_CREATING 2
#define JADE_EVENT_CREATE_FAILED 3
#define JADE_EVENT_DESTROYED 4
#define JADE_EVENT_LOAD_STARTED 5
#define JADE_EVENT_LOAD_FINISHED 6
#define JADE_EVENT_CRASHED 7
#define JADE_EVENT_OPERATION_COMPLETED 8
#define JADE_EVENT_OPERATION_FAILED 9
#define JADE_EVENT_OTHER 100
#define JADE_EVENT_REPLAY 1
#define JADE_STATE_RESERVED 1
#define JADE_STATE_CONFIGURED 2
#define JADE_STATE_CREATING 3
#define JADE_STATE_READY 4
#define JADE_STATE_CLOSING 5
#define JADE_STATE_DESTROYING 6
#define JADE_STATE_DESTROYED 7
#define JADE_STATE_FAILED 8
#define JADE_STATE_CRASHED 9

#ifdef __cplusplus
extern "C" {
#endif

// --- from abi.rs ---
int32_t JADEVIEW_CALL jade_get_abi_info(struct JadeAbiInfo* out, size_t capacity);

// --- from resource_diagnostics.rs ---
// Query a live registration's original decision; output is untouched on failure.
// Pass an empty output, and free it before reuse. Returns JADE_OK or an error code.
int32_t JADEVIEW_CALL get_resource_delivery_info(const char* url, uint32_t webview_id, struct ResourceDeliveryInfo* out, size_t capacity);
void JADEVIEW_CALL free_resource_delivery_info(struct ResourceDeliveryInfo* out);

// --- from japk_api.rs ---
// 从内存加载 JAPK 文件（仅支持 JAPK v2 签名包）
// # 参数
// - `japk_data`: JAPK 文件数据指针
// - `data_size`: 数据大小
// # 返回
// - 0=成功, 负数=错误码
// # 加载逻辑
// - 仅接受带 JadeTweak 平台根证明和叶子签名的 JAPK v2 包
// - 不接受运行时注入的任意公钥，也不回退到旧签名协议或无签名格式
// # 说明
// - app_name、app_signature 和当前进程 EXE 名必须与平台证明一致
// - 证书到期或吊销不影响已签包的离线密码学验证
// - 错误信息通过 jade_on 事件异步通知
int32_t JADEVIEW_CALL JadeView_load_from_bytes(const uint8_t* japk_data, size_t data_size);
// 获取加载状态
// # 返回
// - true=已加载, false=未加载
int32_t JADEVIEW_CALL JadeView_is_loaded(void);
// 获取当前 app_signature
// # 返回
// - app_signature 字符串指针 (调用方需使用 jade_text_free 释放)
char* JADEVIEW_CALL JadeView_get_app_signature(void);
// 获取签名信息 JSON
// # 返回
// - 签名信息 JSON 字符串指针 (调用方需使用 jade_text_free 释放)
char* JADEVIEW_CALL JadeView_get_signature_info(void);
// 清除 JAPK 加载状态
// # 返回
// - 0=成功
int32_t JADEVIEW_CALL JadeView_unload(void);

// --- from webview.rs ---
// Creates a native Window without a WebView. Nonzero means accepted.
// The optional splash is copied now and used on the first WebView attachment.
uint32_t JADEVIEW_CALL create_window(uint32_t parent_window_id, const struct WebViewWindowOptions* options, const struct JadeWindowFeatures* features);
// Create a logical WebView; the native engine is allocated when attached.
uint32_t JADEVIEW_CALL create_webview(const char* url, const struct WebViewSettings* settings);
// Create and attach a second (or first) WebView. NULL settings/bounds use defaults.
uint32_t JADEVIEW_CALL add_webview(uint32_t window_id, const char* url, const struct WebViewSettings* settings, const struct WebViewBounds* bounds);
uint64_t JADEVIEW_CALL add_webview_to_window(uint32_t window_id, uint32_t webview_id, const struct WebViewBounds* bounds);
uint64_t JADEVIEW_CALL set_webview_bounds(uint32_t webview_id, const struct WebViewBounds* bounds);
uint64_t JADEVIEW_CALL navigate_webview(uint32_t webview_id, const char* url);
uint64_t JADEVIEW_CALL reload_webview(uint32_t webview_id);
uint64_t JADEVIEW_CALL execute_webview_script(uint32_t webview_id, const char* script);
uint64_t JADEVIEW_CALL set_webview_visible(uint32_t webview_id, int32_t visible);
uint64_t JADEVIEW_CALL focus_webview(uint32_t webview_id);
uint64_t JADEVIEW_CALL destroy_webview(uint32_t webview_id);
uint64_t JADEVIEW_CALL recreate_webview(uint32_t webview_id, const char* url, const struct WebViewSettings* settings);
uint32_t JADEVIEW_CALL get_default_webview(uint32_t window_id);
int32_t JADEVIEW_CALL set_default_webview(uint32_t window_id, uint32_t webview_id);
int32_t JADEVIEW_CALL get_window_webviews(uint32_t window_id, struct WebViewIdArray* out_array);
void JADEVIEW_CALL free_webview_id_array(struct WebViewIdArray* array);
uint32_t JADEVIEW_CALL get_webview_parent_window_id(uint32_t webview_id);
int32_t JADEVIEW_CALL get_webview_state(uint32_t webview_id, struct JadeObjectState* out);
int32_t JADEVIEW_CALL get_window_state(uint32_t window_id, struct JadeObjectState* out);
// Result state: 1 accepted, 2 waiting, 3 succeeded, 4 failed, 5 cancelled.
// Initialize out->message to NULL. Free with free_request_result, then release_request.
int32_t JADEVIEW_CALL get_request_result(uint64_t request_id, struct JadeRequestResult* out);
void JADEVIEW_CALL free_request_result(struct JadeRequestResult* result);
int32_t JADEVIEW_CALL release_request(uint64_t request_id);
int32_t JADEVIEW_CALL jade_last_error_code(void);
// Copy the thread-local message. Release the returned string with jade_text_free.
char* JADEVIEW_CALL jade_last_error_message(void);
// Write defaults only within capacity. Returns JADE_OK or INVALID_ARGUMENT.
int32_t JADEVIEW_CALL window_options_default(struct WebViewWindowOptions* out, size_t capacity);
// Write defaults only within capacity. Returns JADE_OK or INVALID_ARGUMENT.
int32_t JADEVIEW_CALL webview_settings_default(struct WebViewSettings* out, size_t capacity);
// Applied bounds snapshot. Initialize all string pointers to NULL; release with
// free_webview_bounds_result before reuse. Does not wait for the UI thread.
int32_t JADEVIEW_CALL get_webview_bounds(uint32_t webview_id, struct WebViewBoundsResult* out);
void JADEVIEW_CALL free_webview_bounds_result(struct WebViewBoundsResult* out);

// --- from host_fs.rs ---
// Verify/decrypt a JAPK once, then retain its shared immutable index. 0 is failure.
uint64_t JADEVIEW_CALL register_japk(const char* archive_path);
// Retire a package. Existing native file handles and WebView resource references
// remain valid; mounted filesystems stop accepting new opens. Unmount separately.
int32_t JADEVIEW_CALL unregister_japk(uint64_t archive_id);
// NULL mount_root uses the archive's path without .japk (app.japk -> app/).
// Does not install drivers, enable components, extract files or intercept APIs.
uint64_t JADEVIEW_CALL mount_japk(uint64_t archive_id, const char* mount_root);
// BUSY leaves the provider and existing handles alive; new opens are denied.
// Release files, directory iterators and mappings, then retry.
int32_t JADEVIEW_CALL unmount_japk(uint64_t mount_id);
// Library-allocated UTF-8 string; release with jade_text_free.
char* JADEVIEW_CALL get_japk_mount_root(uint64_t mount_id);
int32_t JADEVIEW_CALL get_japk_mount_capabilities(struct JapkMountCapabilities* out);
int32_t JADEVIEW_CALL get_japk_mount_state(uint64_t mount_id, struct JapkMountState* out);
// Register the already-verified default in-memory package without decrypting it again.
// Memory packages require an explicit mount_root in mount_japk.
uint64_t JADEVIEW_CALL register_loaded_japk(void);

// --- from events.rs ---
uint64_t JADEVIEW_CALL on_app_event(uint32_t event_type, JadeEventCallback callback, void* user_data);
uint64_t JADEVIEW_CALL on_window_event(uint32_t window_id, uint32_t event_type, JadeEventCallback callback, void* user_data);
uint64_t JADEVIEW_CALL on_webview_event(uint32_t webview_id, uint32_t event_type, JadeEventCallback callback, void* user_data);
uint64_t JADEVIEW_CALL unsubscribe_event(uint64_t subscription_id);
// Deep-copy a callback-scoped event. Free the result with event_free.
struct JadeEvent* JADEVIEW_CALL event_clone(const struct JadeEvent* event);
void JADEVIEW_CALL event_free(struct JadeEvent* event);

// --- from lifecycle.rs ---
// DLL初始化
// 备注 · app_signature 长度限制：
// `app_signature` 必填，且 **trim 后至少 `MIN_APP_SIGNATURE_CHARS`（=6）个 Unicode 字符**；
// `app_name` 必填、非纯空白。过短/缺失会被拒绝（例如 `"jvh"` 仅 3 字符 → 拒绝）。
// 返回值：成功返回 `1`；失败返回 `0`（参数校验未过，或启用单实例时本进程为第二实例）。
// 参数校验失败时不启动 GUI 线程，`app-ready` 正文包含错误描述。
// 宿主必须检查初始化返回值及 app-ready 正文，不能仅收到同名通知就继续建窗。
int32_t JADEVIEW_CALL JadeView_init(int32_t enable_devmod, const char* log_path, const char* data_directory, const char* app_name, const char* app_signature, int32_t single_instance);
// 运行消息循环
int32_t JADEVIEW_CALL run_message_loop(void);
// 清理所有窗口并结束消息循环；默认最多等待 30 秒。
int32_t JADEVIEW_CALL jadeview_exit(void);
// 清理所有窗口并等待事件循环和 JadeView 后台线程退出。
int32_t JADEVIEW_CALL jadeview_exit_wait(uint32_t timeout_ms);

// --- from window.rs ---
// Creates a Window and its default WebView. Returns a Window ID.
// Optional features, including the splash and logo, are copied during the call.
uint32_t JADEVIEW_CALL create_webview_window(const char* url, uint32_t parent_window_id, const struct WebViewWindowOptions* options, const struct WebViewSettings* webview_settings, const struct JadeWindowFeatures* features);
// 获取 Windows 窗口的 HWND。窗口尚未就绪或平台不提供 HWND 时返回 0。
// window_id is a Window ID, not a WebView ID.
size_t JADEVIEW_CALL get_window_hwnd(uint32_t window_id);
// 根据 HWND 获取窗口 ID
uint32_t JADEVIEW_CALL get_window_id(size_t hwnd);
// 导航到URL
// webview_id is a WebView ID; use get_default_webview(window_id) after creating a window.
int32_t JADEVIEW_CALL navigate_to_url(uint32_t webview_id, const char* url, const char* headers_json);
// 刷新webview页面
// webview_id is a WebView ID; use get_default_webview(window_id) after creating a window.
int32_t JADEVIEW_CALL reload_webview_window(uint32_t webview_id);
// 原生后退
// webview_id is a WebView ID; use get_default_webview(window_id) after creating a window.
int32_t JADEVIEW_CALL webview_go_back(uint32_t webview_id);
// 原生前进
// webview_id is a WebView ID; use get_default_webview(window_id) after creating a window.
int32_t JADEVIEW_CALL webview_go_forward(uint32_t webview_id);
// 查询是否可后退
// webview_id is a WebView ID; use get_default_webview(window_id) after creating a window.
int32_t JADEVIEW_CALL webview_can_go_back(uint32_t webview_id);
// 查询是否可前进
// webview_id is a WebView ID; use get_default_webview(window_id) after creating a window.
int32_t JADEVIEW_CALL webview_can_go_forward(uint32_t webview_id);
// 设置窗口标题
// window_id is a Window ID, not a WebView ID.
int32_t JADEVIEW_CALL set_window_title(uint32_t window_id, const char* title);
// 设置窗口大小
// window_id is a Window ID, not a WebView ID.
int32_t JADEVIEW_CALL set_window_size(uint32_t window_id, int32_t width, int32_t height);
// 设置窗口位置
// window_id is a Window ID, not a WebView ID.
int32_t JADEVIEW_CALL set_window_position(uint32_t window_id, int32_t x, int32_t y);
// 设置窗口可见性
// window_id is a Window ID, not a WebView ID.
int32_t JADEVIEW_CALL set_window_visible(uint32_t window_id, int32_t visible);
// 设置窗口焦点
// window_id is a Window ID, not a WebView ID.
int32_t JADEVIEW_CALL set_window_focus(uint32_t window_id);
// 设置窗口是否置顶
// window_id is a Window ID, not a WebView ID.
int32_t JADEVIEW_CALL set_window_always_on_top(uint32_t window_id, int32_t always_on_top);
// 不进任务栏/Alt-Tab（跨平台：Windows WS_EX_TOOLWINDOW、Linux GTK skip-taskbar-hint）
// window_id is a Window ID, not a WebView ID.
int32_t JADEVIEW_CALL set_window_skip_taskbar(uint32_t window_id, int32_t skip);
// 不抢焦点：点击/显示不激活窗口（Windows WS_EX_NOACTIVATE、Linux gtk accept-focus=false）
// window_id is a Window ID, not a WebView ID.
int32_t JADEVIEW_CALL set_window_no_activate(uint32_t window_id, int32_t no_activate);
// 设置窗口层级：topmost | normal | bottom | desktop
// topmost/normal/bottom 跨平台（tao set_always_on_top/bottom）；
// desktop=贴桌面壁纸层（Windows 寄生 WorkerW、Linux GTK desktop type hint）。
// window_id is a Window ID, not a WebView ID.
int32_t JADEVIEW_CALL set_window_level(uint32_t window_id, const char* level);
// 关闭窗口
// window_id is a Window ID, not a WebView ID.
int32_t JADEVIEW_CALL close_window(uint32_t window_id);
// 最小化窗口
// window_id is a Window ID, not a WebView ID.
int32_t JADEVIEW_CALL minimize_window(uint32_t window_id);
// 切换窗口最大化/还原状态
// window_id is a Window ID, not a WebView ID.
int32_t JADEVIEW_CALL toggle_maximize_window(uint32_t window_id);
// 获取窗口是否最大化
// window_id is a Window ID, not a WebView ID.
int32_t JADEVIEW_CALL is_window_maximized(uint32_t window_id);
// 设置内容保护
// window_id is a Window ID, not a WebView ID.
int32_t JADEVIEW_CALL set_content_protection(uint32_t window_id, int32_t content_protection);
// 设置WebView缩放级别
// webview_id is a WebView ID; use get_default_webview(window_id) after creating a window.
int32_t JADEVIEW_CALL set_webview_zoom(uint32_t webview_id, double level);
// 修改标题栏样式，窗口边框保持创建时的配置。
int32_t JADEVIEW_CALL set_window_frame(uint32_t window_id, const char* title_bar_style);
// 获取窗口数量
uint32_t JADEVIEW_CALL get_window_count(void);
// 注册事件处理器
// Callback argument 1 is the event source ID: App=0, Window=Window ID, WebView=WebView ID. app-ready success is JSON with ok:true, not source ID 1.
uint32_t JADEVIEW_CALL jade_on(const char* event_name, IpcCallback callback);
// 注销事件处理器
// 返回时只保证回调不再接收后续事件，不等待已经开始执行的回调；注销后要立即释放上下文时使用 jade_off_wait。
int32_t JADEVIEW_CALL jade_off(const char* event_name, uint32_t callback_id);
// Unregisters a named-event callback and returns a request that completes only
// after every invocation that already started has returned. Wait for the
// request before freeing the callback context. Unknown or inactive callback
// IDs complete the request with JADE_INVALID_ID; an invalid event name uses
// JADE_INVALID_ARGUMENT; 0 means no request could be created.
// `jade_off` remains the immediate-return variant and does not wait.
uint64_t JADEVIEW_CALL jade_off_wait(const char* event_name, uint32_t callback_id);
// 订阅前端发送的IPC消息
// Invoke callback argument 1 is the originating WebView ID, not its parent Window ID.
int32_t JADEVIEW_CALL register_ipc_handler(const char* channel, IpcCallback ipc_cb);
// 设置窗口主题
// window_id is a Window ID, not a WebView ID.
int32_t JADEVIEW_CALL set_window_theme(uint32_t window_id, const char* theme);
// Windows: configure navigation buttons, caption icon visibility and minimize/maximize capability.
// Icon visibility applies to the default caption and never changes the taskbar icon.
// Back appears only when history exists. Refresh is disabled and spins while loading.
// Returns 1 when queued, 0 on rejection. Every flag must be 0 or 1.
int32_t JADEVIEW_CALL set_titlebar_components(uint32_t window_id, int32_t back, int32_t refresh, int32_t icon, int32_t minimize, int32_t maximize);
// Windows: style default and BarOverlay captions. Window ID, not WebView ID.
// height: -1 keeps the current height, 0 resets to the mode default, 20..=100 sets DIP.
// NULL and empty colors reset to the current theme's default.
// background_color applies to the default caption. foreground_color sets the
// title text and normal button glyph color in both caption modes where drawn.
int32_t JADEVIEW_CALL set_titlebar_style(uint32_t window_id, int32_t height, const char* background_color, const char* foreground_color, const char* button_background_color, const char* button_hover_background_color, const char* button_hover_foreground_color);
// 启用/禁用窗口
// window_id is a Window ID, not a WebView ID.
int32_t JADEVIEW_CALL set_window_enabled(uint32_t window_id, int32_t enabled);
// 获取窗口主题
// window_id is a Window ID, not a WebView ID.
int32_t JADEVIEW_CALL get_window_theme(uint32_t window_id);
// 请求窗口重绘
// window_id is a Window ID, not a WebView ID.
int32_t JADEVIEW_CALL request_redraw(uint32_t window_id);
// 发送IPC消息
// webview_id is a WebView ID; use get_default_webview(window_id) after creating a window.
// Returns 1 when queued, 0 on failure; inspect jade_last_error_code/message on this thread.
int32_t JADEVIEW_CALL send_ipc_message(uint32_t webview_id, const char* message_type, const char* message_content);
// 设置窗口背景材料
// window_id is a Window ID, not a WebView ID.
int32_t JADEVIEW_CALL set_window_backdrop(uint32_t window_id, const char* backdrop_type);
// 设置窗口背景色（纯色底）：#RRGGBBAA
// window_id is a Window ID, not a WebView ID.
int32_t JADEVIEW_CALL set_window_background_color(uint32_t window_id, const char* background_color_hex);
// 设置窗口全屏
// window_id is a Window ID, not a WebView ID.
int32_t JADEVIEW_CALL set_window_fullscreen(uint32_t window_id, int32_t fullscreen);
// window_id is a Window ID, not a WebView ID.
int32_t JADEVIEW_CALL get_window_bounds(uint32_t window_id, char* buffer, int32_t buffer_size);
// window_id is a Window ID, not a WebView ID.
int32_t JADEVIEW_CALL is_window_minimized(uint32_t window_id);
// window_id is a Window ID, not a WebView ID.
int32_t JADEVIEW_CALL is_window_visible(uint32_t window_id);
// window_id is a Window ID, not a WebView ID.
int32_t JADEVIEW_CALL is_window_focused(uint32_t window_id);
// window_id is a Window ID, not a WebView ID.
int32_t JADEVIEW_CALL is_window_fullscreen(uint32_t window_id);
// window_id is a Window ID, not a WebView ID.
int32_t JADEVIEW_CALL set_window_min_size(uint32_t window_id, int32_t width, int32_t height);
// window_id is a Window ID, not a WebView ID.
int32_t JADEVIEW_CALL set_window_max_size(uint32_t window_id, int32_t width, int32_t height);
// window_id is a Window ID, not a WebView ID.
int32_t JADEVIEW_CALL set_window_resizable(uint32_t window_id, int32_t resizable);
// window_id is a Window ID, not a WebView ID.
int32_t JADEVIEW_CALL set_window_ignore_cursor_events(uint32_t window_id, int32_t ignore);
// webview_id is a WebView ID; use get_default_webview(window_id) after creating a window.
int32_t JADEVIEW_CALL get_webview_url(uint32_t webview_id, char* buffer, int32_t buffer_size);
// 设置统一网页权限处理器
int32_t JADEVIEW_CALL set_webview_permission_handler(PermissionCallback callback);
// 清除统一网页权限处理器
int32_t JADEVIEW_CALL clear_webview_permission_handler(void);
// webview_id is a WebView ID; use get_default_webview(window_id) after creating a window.
int32_t JADEVIEW_CALL open_devtools(uint32_t webview_id);
// webview_id is a WebView ID; use get_default_webview(window_id) after creating a window.
int32_t JADEVIEW_CALL close_devtools(uint32_t webview_id);
// webview_id is a WebView ID; use get_default_webview(window_id) after creating a window.
int32_t JADEVIEW_CALL is_devtools_open(uint32_t webview_id);
// webview_id is a WebView ID; use get_default_webview(window_id) after creating a window.
int32_t JADEVIEW_CALL clear_browsing_data(uint32_t webview_id);
// window_id is a Window ID, not a WebView ID.
int32_t JADEVIEW_CALL set_window_progress(uint32_t window_id, int32_t progress, int32_t state);
// window_id is a Window ID, not a WebView ID.
int32_t JADEVIEW_CALL flash_window(uint32_t window_id, uint32_t count);
// window_id is a Window ID, not a WebView ID.
int32_t JADEVIEW_CALL show_about_dialog(uint32_t window_id);

// --- from tray.rs ---
uint32_t JADEVIEW_CALL tray_create(void);
int32_t JADEVIEW_CALL tray_destroy(uint32_t tray_id);
int32_t JADEVIEW_CALL tray_set_visible(uint32_t tray_id, int32_t visible);
int32_t JADEVIEW_CALL tray_set_tooltip(uint32_t tray_id, const char* tooltip_utf8);
int32_t JADEVIEW_CALL tray_set_icon_from_file(uint32_t tray_id, const char* icon_path_utf8);
int32_t JADEVIEW_CALL set_tray_icon_from_data(uint32_t tray_id, const uint8_t* icon_data, uint32_t data_len);

// --- from menu.rs ---
// Shows an independent popup. Window ID 0 uses the screen as its coordinate space.
// Position is in client logical pixels for a window, or desktop coordinates without one.
// NULL position uses the cursor. Entries and position are copied before return.
// Item icons are part of JadeMenuItemDesc.
// Returns 1 when queued. A selection emits menu-command with item_id.
int32_t JADEVIEW_CALL jade_menu_popup(uint32_t window_id, const struct JadeMenuItemDesc* items, uint32_t item_count, const struct JadeMenuPosition* position);
// Replaces the application menu, or hides it with a zero item count.
int32_t JADEVIEW_CALL jade_set_application_menu(const struct JadeMenuItemDesc* items, uint32_t item_count);
// Replaces the menu for a Window, or hides it with a zero item count.
int32_t JADEVIEW_CALL jade_window_set_menu(uint32_t window_id, const struct JadeMenuItemDesc* items, uint32_t item_count);
// Sets the page context-menu template for one WebView. The rows and strings
// are copied during the call. Items without a role are custom items and
// report `context-menu-click`; role items use unified action names.
// Pass count=0 to clear.
int32_t JADEVIEW_CALL jade_webview_set_context_menu(uint32_t webview_id, const struct JadeContextMenuItemDesc* items, uint32_t item_count);
// Removes the page context-menu template for one WebView.
int32_t JADEVIEW_CALL jade_webview_clear_context_menu(uint32_t webview_id);

// --- from dialog.rs ---
// 显示通知
int32_t JADEVIEW_CALL show_notification(const struct NotificationParams* params);
// 显示打开文件对话框
// 参数:
// - params: 对话框参数结构体指针
// 返回:
// - 1: 成功
// - 0: 失败
char* JADEVIEW_CALL jade_dialog_show_open_dialog(const struct FileDialogParams* params);
// 显示保存文件对话框
// 参数:
// - params: 对话框参数结构体指针
// 返回:
// - 1: 成功
// - 0: 失败
char* JADEVIEW_CALL jade_dialog_show_save_dialog(const struct FileDialogParams* params);
// 显示消息框
// 参数:
// - params: 消息框参数结构体指针
// 返回:
// - 1: 成功
// - 0: 失败
char* JADEVIEW_CALL jade_dialog_show_message_box(const struct MessageBoxParams* params);
// 显示错误框
// 参数:
// - window_id: 窗口ID
// - title: 错误标题
// - content: 错误内容
// 返回:
// - 1: 成功
// - 0: 失败
int32_t JADEVIEW_CALL jade_dialog_show_error_box(uint32_t window_id, const char* title, const char* content);
int32_t JADEVIEW_CALL jade_dialog_show_open_dialog_async(const struct FileDialogParams* params, void (*callback)(const char*));
int32_t JADEVIEW_CALL jade_dialog_show_save_dialog_async(const struct FileDialogParams* params, void (*callback)(const char*));
int32_t JADEVIEW_CALL jade_dialog_show_message_box_async(const struct MessageBoxParams* params, void (*callback)(const char*));

// --- from yaml_store.rs ---
// YAML: set value by key path (a.b.c). `value` 可为 JSON、YAML 片段或任意纯文本（见 parse_yaml_set_payload）。
// 返回值：1=成功，0=路径不存在/空操作，-1=IO错误，-2=类型不匹配，-4=格式解析失败
int32_t JADEVIEW_CALL yaml_set(const char* file_name, const char* key_path, const char* value);
// YAML: get value by key path; returns JSON string in buffer.
// 返回值：1=成功，≥2=两阶段查询所需字节数(含NUL)，0=路径不存在，-1=IO错误，-4=格式解析失败
int32_t JADEVIEW_CALL yaml_get(const char* file_name, const char* key_path, char* buffer, size_t buffer_size);
// YAML: 强制字符串存储（不尝试 JSON/YAML 解析）
// 返回值：1=成功，0=空操作，-1=IO错误，-2=类型不匹配，-4=格式解析失败
int32_t JADEVIEW_CALL yaml_set_str(const char* file_name, const char* key_path, const char* value);
// YAML: 内部 malloc 返回字符串指针，调用方 CoTaskMemFree 释放
// 返回值：非空指针=成功，空指针=失败
char* JADEVIEW_CALL yaml_get_str(const char* file_name, const char* key_path);
// YAML: 读取整个文件，返回 JSON 字符串到 buffer
// 返回值：1=成功，≥2=所需字节数(含NUL)，0=文件不存在，-1=IO错误，-4=格式解析失败
int32_t JADEVIEW_CALL yaml_get_all(const char* file_name, char* buffer, size_t buffer_size);
// YAML: 检查路径是否存在
// 返回值：1=存在，0=不存在，-1=IO错误，-4=格式解析失败
int32_t JADEVIEW_CALL yaml_has(const char* file_name, const char* key_path);
// YAML: 删除指定路径
// 返回值：1=成功，0=路径不存在，-1=IO错误，-2=类型不匹配，-4=格式解析失败
int32_t JADEVIEW_CALL yaml_delete(const char* file_name, const char* key_path);
// YAML: 清空文件为 {}
// 返回值：1=成功，-1=IO错误
int32_t JADEVIEW_CALL yaml_clear(const char* file_name);
// YAML: 删除文件
// 返回值：1=成功，0=文件不存在，-1=IO错误
int32_t JADEVIEW_CALL yaml_delete_file(const char* file_name);
// YAML: 列出路径下的所有 key
// 返回值：1=成功，≥2=所需字节数(含NUL)，0=路径不存在/非映射非序列，-1=IO错误，-4=格式解析失败
int32_t JADEVIEW_CALL yaml_keys(const char* file_name, const char* key_path, char* buffer, size_t buffer_size);
// YAML: 返回数组长度/对象 key 数
// 返回值：≥0=长度，-1=IO错误，-2=类型不匹配(非映射非序列)，-4=格式解析失败
int32_t JADEVIEW_CALL yaml_len(const char* file_name, const char* key_path);

// --- from system.rs ---
// 打开打印对话框（WebView2 内置）
int32_t JADEVIEW_CALL jade_print(uint32_t window_id);
// 使用系统关联程序打印文件（Windows=ShellExecute "print"、Linux=CUPS lp），调用点无 cfg。
int32_t JADEVIEW_CALL jade_print_dialog(const char* file_path);
// 获取系统打印机列表（Windows=EnumPrintersW、Linux=CUPS lpstat -e）。
// 写入 JSON 字符串数组到 buffer，返回打印机数量。
int32_t JADEVIEW_CALL jade_get_printer_list(char* buffer, int32_t buffer_size);
// 智能转码：自动检测输入文本编码，转换为目标编码
// # 参数
// - `input_data`: 输入字节流
// - `input_len`: 输入字节长度
// - `target_encoding`: 目标编码名称，大小写不敏感，支持 WhatWG 标准名称和别名：
// - UTF-8: `"utf-8"`, `"utf8"`, `"unicode-1-1-utf-8"`
// - GBK: `"gbk"`, `"gb2312"`, `"gb18030"`, `"chinese"`, `"csgb2312"`
// - Shift_JIS: `"shift_jis"`, `"sjis"`, `"shift-jis"`, `"ms_kanji"`
// - Big5: `"big5"`, `"big5-hkscs"`, `"cn-big5"`, `"csbig5"`
// - EUC-KR: `"euc-kr"`, `"cseuckr"`, `"korean"`
// - windows-1252: `"windows-1252"`, `"ascii"`, `"latin1"`, `"iso-8859-1"`
// - 完整列表见 https://encoding.spec.whatwg.org/
// - `output_buffer`: 输出缓冲区（NUL 结尾）
// - `buffer_size`: 输出缓冲区大小（字节）
// - `detected_encoding`: [输出] 检测到的来源编码名称（可为 NULL）
// - `detected_encoding_size`: detected_encoding 缓冲区大小
// # 返回值
// - `>0`: 成功，写入输出缓冲区的字节数（不含 NUL）
// - `0`: 失败（参数无效、编码检测失败、目标编码不支持等）
// - `<0`: 缓冲区不足，绝对值为所需缓冲区大小（含 NUL），不写入正文
// - 输出长度（含 NUL）无法用 i32 表示时返回 0
// # 检测逻辑
// 1. BOM 优先（UTF-8 BOM / UTF-16 LE BOM / UTF-16 BE BOM）
// 2. 无 BOM 时使用 chardetng 智能检测
// 3. 源编码与目标编码相同时直接复制
int32_t JADEVIEW_CALL smart_convert_encoding(const uint8_t* input_data, int32_t input_len, const char* target_encoding, char* output_buffer, int32_t buffer_size, char* detected_encoding, int32_t detected_encoding_size);
// 设置自定义协议服务路径（原 create_local_server）
// hot_reload: 是否启用热载模式 (0=禁用, 1=启用)，仅文件系统模式有效
int32_t JADEVIEW_CALL set_protocol_service_path(const char* root_path, char* url_buffer, size_t buffer_size, int32_t hot_reload);
// 获取常用系统 / 应用路径，结果写入 `buffer`（UTF-8，NUL 结尾）。
// `name` 取值：
// - `home`：用户主目录
// - `appData`：应用数据目录（Windows=%LOCALAPPDATA%，Linux=~/.config）
// - `userData`：JadeView 数据目录（`sessionData` / `logs` 的父目录，不存在时自动创建）
// - `sessionData`：WebView 会话目录（`userData/EBWebView`）
// - `logs`：日志目录（`userData/logs`，不存在时自动创建）
// - `temp`：临时目录
// - `exe`：主进程完整路径
// - `app`：主进程所在目录
// - `desktop` / `documents` / `downloads` / `music` / `pictures` / `videos`：标准用户目录
// # 返回值
// - `>0`：成功，写入 `buffer` 的字节数（不含 NUL）
// - `0`：失败（参数无效或该名称无对应路径）
// - `<0`：缓冲区不足，绝对值为所需字节数
int32_t JADEVIEW_CALL getPath(const char* name, char* buffer, size_t buffer_size);
// 显示器信息 JSON 数组写入 `buffer`（UTF-8 NUL 结尾）。每项含 bounds、work_area、scale_factor、dpi_x/y、is_primary。
// `buffer` 过小则返回 0。
int32_t JADEVIEW_CALL get_displays_info(char* buffer, size_t buffer_size);
// 获取系统语言（BCP 47，如 zh-CN）
int32_t JADEVIEW_CALL getLocale(char* buffer, size_t buffer_size);
// 清空数据目录（安全确认）
// confirm_token 必须等于 "I_UNDERSTAND_CLEAR_DATA" 才会执行
int32_t JADEVIEW_CALL clear_data_directory(const char* confirm_token);
// JadeView 版本字符串：`Cargo` 语义化版本 + `.` + 构建号（见 `build.rs`）。
int32_t JADEVIEW_CALL jadeview_version(char* buffer, size_t buffer_size);
int32_t JADEVIEW_CALL register_url_scheme(const char* scheme);
int32_t JADEVIEW_CALL unregister_url_scheme(const char* scheme);
int32_t JADEVIEW_CALL register_file_association(const char* extension, const char* friendly_name);
int32_t JADEVIEW_CALL unregister_file_association(const char* extension);
uint32_t JADEVIEW_CALL register_global_hotkey(uint32_t modifiers, uint32_t vk);
int32_t JADEVIEW_CALL unregister_global_hotkey(uint32_t hotkey_id);
// 启用/取消开机自启。
// 参数：enable（0=取消，1=启用）、args（追加到可执行路径后的启动参数，NULL/空=无）。
// 返回：1=成功，0=失败。
int32_t JADEVIEW_CALL set_login_autostart(int32_t enable, const char* args);
// 查询是否已设置开机自启。返回：1=已设置，0=未设置。
int32_t JADEVIEW_CALL get_login_autostart(void);
// 提取任意路径（exe/lnk/普通文件/文件夹）的图标为 PNG，注册为 jade:// 资源，URL 写入 url_buffer。
// 参数：size 目标边长（16/32/48/64/128/256，<=0 取 48）；window_id 资源归属 WebView ID（0=全局）；
// ttl_seconds 资源过期秒（0=默认，u32::MAX=永不过期）。
// 返回：1=成功，0=失败。Windows 用 SHGetFileInfo 取系统图标；Linux 用 gio 查询 + gtk IconTheme。
int32_t JADEVIEW_CALL get_file_icon(const char* path, int32_t size, uint32_t window_id, uint32_t ttl_seconds, char* url_buffer, size_t buffer_size);
// 获取WebView版本号
int32_t JADEVIEW_CALL get_webview_version(char* buffer, size_t buffer_size);
// 检查当前系统是否为Windows 11
int32_t JADEVIEW_CALL is_windows_11(void);
void JADEVIEW_CALL jade_text_free(char* ptr);
char* JADEVIEW_CALL jade_text_create(const char* text);
// 注册本地文件为安全资源，返回可供页面使用的相对 URL 或完整 HTTPS URL。
// Windows 的永久私有资源会先尝试原生快速模式；条件不满足时使用受控模式。
// 原生模式可能锁定源文件，HEAD、条件请求和路径替换的行为与受控模式不同。
// 需要严格受控语义时请使用 register_resource_ex 和 JADE_RESOURCE_CONTROLLED。
// path: 本地文件绝对路径
// window_id: 所属 WebView ID (0=全局)
// ttl_seconds: 过期时间秒数 (0=默认 6 秒，u32::MAX=永不过期；长媒体应显式指定)
// url_buffer: 输出URL字符串的缓冲区
// buffer_size: 缓冲区大小
// 返回: 1=成功, 0=失败
int32_t JADEVIEW_CALL register_resource(const char* path, uint32_t window_id, uint32_t ttl_seconds, char* url_buffer, size_t buffer_size);
// 注销已注册的安全资源
// token_or_url: token字符串或完整jade://---jade---resource--?token=xxx URL
// 返回: 1=成功, 0=失败或未找到。原生模式 BUSY 时先停止页面消费者再重试；
// 已缓存或已打开的数据无法追回，成功后新的未缓存读取不可再打开该文件。
int32_t JADEVIEW_CALL unregister_resource(const char* token_or_url);
// 清理指定窗口的所有已注册资源
// window_id: WebView ID
// 返回: 实际清理数量（含受控和原生资源）。原生项可能因占用部分失败；
// 检查 jade_last_error_code，BUSY 项仍登记，停止消费者后可再次清理。
int32_t JADEVIEW_CALL clear_window_resources(uint32_t window_id);
// Resource mode capabilities for a Ready WebView. Returns JADE_OK or an error.
// supported_modes bits: JADE_RESOURCE_CAP_CONTROLLED / JADE_RESOURCE_CAP_NATIVE.
// Native mode needs a private, permanent resource and an exclusive named Profile.
// Availability covers engine/policy; the file's filesystem must support hard links.
int32_t JADEVIEW_CALL get_resource_capabilities(uint32_t webview_id, struct ResourceCapabilities* out);
// Query the actual mode of an active resource. Returns JADE_OK or an error.
// Use its full returned URL and its owner WebView ID; global resources use 0.
int32_t JADEVIEW_CALL get_resource_delivery_mode(const char* url, uint32_t webview_id, uint32_t* out_mode);
// Register a local file with explicit ResourceOptions (NULL = controlled).
// Controlled mode always uses JadeView's strict resource handler. Native mode:
// Windows, private Ready WebView, ttl_seconds=UINT32_MAX, exclusive named Profile,
// unrestricted CORS; a hard link pins the registered file object without copying.
// Prepared roots cover local temp/data-directory volumes with hard-link support;
// other volumes fail explicitly. New registrations need no page reload.
// In-place writers are denied while registered; replacing the source path does
// not rebind this resource. Native MIME/Range/cache/conditions are engine-owned:
// HEAD may return a full body and conditional headers may be ignored. Use
// controlled mode for JadeView's strict HTTP semantics, TTL or global resources.
// Native output is HTTPS. Release with unregister_resource; BUSY means detach
// consumers and retry. Cached/already-open browser data cannot be recalled.
// Returns 1 on success, 0 on failure (jade_last_error_code/message).
int32_t JADEVIEW_CALL register_resource_ex(const char* path, uint32_t webview_id, uint32_t ttl_seconds, const struct ResourceOptions* options, char* url_buffer, size_t buffer_size);
int32_t JADEVIEW_CALL clipboard_read_text(char* buffer, int32_t buffer_size);
int32_t JADEVIEW_CALL clipboard_write_text(const char* text);
int32_t JADEVIEW_CALL get_cursor_position(char* buffer, int32_t buffer_size);
// 获取当前网络时间戳（毫秒）
// # 参数
// - `ntp_server`: NTP 服务器地址（如 `"ntp.aliyun.com"`）。传 `NULL` 或空字符串时，
// 使用内置的 NTP 服务器列表逐个尝试。
// # 返回值
// - `>= 0`：Unix 时间戳（毫秒）
// - `-1`：获取失败
// # Safety
// `ntp_server` 可为 `NULL`；若非空必须指向以 NUL 结尾的有效 C 字符串。
int64_t JADEVIEW_CALL jade_ntp_now(const char* ntp_server);

#ifdef __cplusplus
}
#endif

#endif  /* JADEVIEW_H */
