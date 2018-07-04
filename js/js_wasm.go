// Copyright 2018 The GopherWasm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build wasm

package js

import (
	"reflect"
	"syscall/js"
	"unsafe"
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

var (
	int8Array    = js.Global().Get("Int8Array")
	int16Array   = js.Global().Get("Int16Array")
	int32Array   = js.Global().Get("Int32Array")
	uint16Array  = js.Global().Get("Uint16Array")
	uint32Array  = js.Global().Get("Uint32Array")
	float32Array = js.Global().Get("Float32Array")
	float64Array = js.Global().Get("Float64Array")
)

func ValueOf(x interface{}) Value {
	var xh *reflect.SliceHeader
	var class js.Value
	size := 0
	// TODO: Now slices must be passed to TypedArrayOf. Remove this.
	switch x := x.(type) {
	case []int8:
		size = 1
		xh = (*reflect.SliceHeader)(unsafe.Pointer(&x))
		class = int8Array
	case []int16:
		size = 2
		xh = (*reflect.SliceHeader)(unsafe.Pointer(&x))
		class = int16Array
	case []int32:
		size = 4
		xh = (*reflect.SliceHeader)(unsafe.Pointer(&x))
		class = int32Array
	case []uint16:
		size = 2
		xh = (*reflect.SliceHeader)(unsafe.Pointer(&x))
		class = uint16Array
	case []uint32:
		size = 4
		xh = (*reflect.SliceHeader)(unsafe.Pointer(&x))
		class = uint32Array
	case []float32:
		size = 4
		xh = (*reflect.SliceHeader)(unsafe.Pointer(&x))
		class = float32Array
	case []float64:
		size = 8
		xh = (*reflect.SliceHeader)(unsafe.Pointer(&x))
		class = float64Array
	default:
		return js.ValueOf(x)
	}

	var b []byte
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh.Data = xh.Data
	bh.Len = xh.Len * size
	bh.Cap = xh.Cap * size

	u8 := js.ValueOf(b)
	return class.New(u8.Get("buffer"), u8.Get("byteOffset"), xh.Len)
}

type TypedArray = js.TypedArray

func TypedArrayOf(slice interface{}) TypedArray {
	return js.TypedArrayOf(slice)
}

func GetInternalObject(v Value) interface{} {
	return v
}
