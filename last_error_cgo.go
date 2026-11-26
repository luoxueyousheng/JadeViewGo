//go:build linux && amd64 && jadeview_cgo && cgo

package jadeview

// #include "JadeView.h"
import "C"

func LastErrorCode() int32     { return int32(C.jade_last_error_code()) }
func LastErrorMessage() string { return goStringFree(C.jade_last_error_message()) }
