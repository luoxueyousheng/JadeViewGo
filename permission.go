//go:build windows || (linux && amd64 && (!jadeview_cgo || cgo))

package jadeview

import "sync"

var (
	permissionRegisterMu sync.Mutex
	permissionMu         sync.RWMutex
	permissionHandler    PermissionDecisionHandler
)

func permissionDecision(view uint32, data string) PermissionDecision {
	permissionMu.RLock()
	handler := permissionHandler
	permissionMu.RUnlock()
	if handler == nil {
		return PermissionDefault
	}
	return normalizePermission(handler(view, data))
}

// SetWebviewPermissionDecisionHandler 注册同步三态处理器；传 nil 恢复浏览器默认策略。
func SetWebviewPermissionDecisionHandler(handler PermissionDecisionHandler) bool {
	permissionRegisterMu.Lock()
	defer permissionRegisterMu.Unlock()
	if !setPermissionCallback(handler != nil) {
		return false
	}
	permissionMu.Lock()
	permissionHandler = handler
	permissionMu.Unlock()
	return true
}

// SetWebviewPermissionHandler 注册双态允许/拒绝处理器。
func SetWebviewPermissionHandler(handler PermissionHandler) bool {
	if handler == nil {
		return ClearWebviewPermissionHandler()
	}
	return SetWebviewPermissionDecisionHandler(func(view uint32, data string) PermissionDecision {
		if handler(view, data) {
			return PermissionAllow
		}
		return PermissionDeny
	})
}

// ClearWebviewPermissionHandler 清除处理器，恢复浏览器默认策略。
func ClearWebviewPermissionHandler() bool { return SetWebviewPermissionDecisionHandler(nil) }
