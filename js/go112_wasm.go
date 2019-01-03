// Copyright 2019 The GopherWasm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build wasm
// +build go1.12

package js

import (
	"syscall/js"
)

type Func = js.Func

func FuncOf(fn func(this Value, args []Value) interface{}) Func {
	return js.FuncOf(fn)
}

// Callback is for backward compatibility. Use Func instead.
type Callback = js.Func

type EventCallbackFlag int

const (
	PreventDefault EventCallbackFlag = 1 << iota
	StopPropagation
	StopImmediatePropagation
)

// NewCallback is for backward compatibility. Use FuncOf instead.
func NewCallback(fn func([]Value)) Callback {
	return FuncOf(func(this Value, args []Value) interface{} {
		go func() {
			fn(args)
		}()
		return nil
	})
}

// NewEventCallback is for backward compatibility. Use FuncOf instead.
func NewEventCallback(flags EventCallbackFlag, fn func(event Value)) Callback {
	return FuncOf(func(this Value, args []Value) interface{} {
		e := js.Undefined()
		if len(args) > 0 {
			e = args[0]
			if flags&PreventDefault != 0 {
				e.Call("preventDefault")
			}
			if flags&StopPropagation != 0 {
				e.Call("stopPropagation")
			}
			if flags&StopImmediatePropagation != 0 {
				e.Call("stopImmediatePropagation")
			}
		}
		go func() {
			fn(e)
		}()
		return nil
	})
}
