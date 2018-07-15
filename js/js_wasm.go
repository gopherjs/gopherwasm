// Copyright 2018 The GopherWasm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build wasm

package js

import (
	"syscall/js"
)

type Type = js.Type

const (
	TypeUndefined Type = js.TypeUndefined
	TypeNull      Type = js.TypeNull
	TypeBoolean   Type = js.TypeBoolean
	TypeNumber    Type = js.TypeNumber
	TypeString    Type = js.TypeString
	TypeSymbol    Type = js.TypeSymbol
	TypeObject    Type = js.TypeObject
	TypeFunction  Type = js.TypeFunction
)

func Global() Value {
	return js.Global()
}

func Null() Value {
	return js.Null()
}

func Undefined() Value {
	return js.Undefined()
}

type Callback = js.Callback

type EventCallbackFlag = js.EventCallbackFlag

const (
	PreventDefault           = js.PreventDefault
	StopPropagation          = js.StopPropagation
	StopImmediatePropagation = js.StopImmediatePropagation
)

func NewCallback(f func([]Value)) Callback {
	return js.NewCallback(f)
}

func NewEventCallback(flags EventCallbackFlag, fn func(event Value)) Callback {
	return js.NewEventCallback(flags, fn)
}

type Error = js.Error

type Value = js.Value

func ValueOf(x interface{}) Value {
	return js.ValueOf(x)
}

type TypedArray = js.TypedArray

func TypedArrayOf(slice interface{}) TypedArray {
	return js.TypedArrayOf(slice)
}

func GetInternalObject(v Value) interface{} {
	return v
}
