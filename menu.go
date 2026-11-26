package jadeview

// MenuKind 是 v3 菜单项类型。
type MenuKind uint32

const (
	MenuCommand MenuKind = iota
	MenuSubmenu
	MenuSeparator
	MenuCheckbox
	MenuRadio
)

// MenuIcon 是可选的菜单项图标；Codepoint 非零时优先使用字体图标。
type MenuIcon struct {
	Path      string
	Data      []byte
	Codepoint uint32
}

// MenuItem 使用字符串 ID，子项通过 ParentID 引用父菜单项。
type MenuItem struct {
	ID, ParentID, Label, Accelerator string
	Kind                             MenuKind
	Disabled, Checked                bool
	Icon                             MenuIcon
}

func (m MenuItem) flags() uint32 {
	var flags uint32
	if m.Disabled {
		flags |= 1
	}
	if m.Checked {
		flags |= 2
	}
	return flags
}

// ContextMenuItem 为网页右键菜单增加原生动作角色，例如 copy 或 paste。
type ContextMenuItem struct {
	MenuItem
	Role string
}

// MenuPosition 为窗口客户区逻辑像素；无归属窗口时为屏幕坐标。
type MenuPosition struct{ X, Y int32 }
