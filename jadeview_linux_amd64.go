//go:build linux && amd64 && jadeview_cgo && cgo

package jadeview

// #cgo LDFLAGS: -L${SRCDIR}/lib/linux_amd64 -lJadeView
import "C"
