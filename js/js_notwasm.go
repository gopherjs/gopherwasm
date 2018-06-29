// Copyright 2018 The GopherWasm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !wasm

package js

import (
	"reflect"
	"unsafe"

	"github.com/gopherjs/gopherjs/js"
)

func Global() Value {
	return Value{v: js.Global}
}

func Null() Value {
	return Value{v: nil}
}

func Undefined() Value {
	return Value{v: js.Undefined}
}

type Callback struct {
	f     func([]Value)
	flags EventCallbackFlag
}

type EventCallbackFlag int

const (
	PreventDefault EventCallbackFlag = 1 << iota
	StopPropagation
	StopImmediatePropagation
)

func NewCallback(f func([]Value)) Callback {
	return Callback{f: f}
}

func NewEventCallback(flags EventCallbackFlag, fn func(event Value)) Callback {
	f := func(args []Value) {
		e := args[0]
		fn(e)
	}
	return Callback{
		f:     f,
		flags: flags,
	}
}

func (c Callback) Release() {
}

type Error struct {
	e *js.Error
}

func (e Error) Error() string {
	return e.e.Error()
}

type Value struct {
	v *js.Object
}

var (
	id         *js.Object
	instanceOf *js.Object
)

func init() {
	if js.Global != nil {
		id = js.Global.Call("eval", "(function(x) { return x; })")
		instanceOf = js.Global.Call("eval", "(function(x, y) { return x instanceof y; })")
	}
}

func ValueOf(x interface{}) Value {
	switch x := x.(type) {
	case Value:
		return x
	case Callback:
		return Value{
			v: id.Invoke(func(args ...*js.Object) {
				if len(args) > 0 {
					e := args[0]
					if x.flags&PreventDefault != 0 {
						e.Call("preventDefault")
					}
					if x.flags&StopPropagation != 0 {
						e.Call("stopPropagation")
					}
					if x.flags&StopImmediatePropagation != 0 {
						e.Call("stopImmediatePropagation")
					}
				}

				// Call the function asyncly to emulate Wasm's Callback more
				// precisely.
				go func() {
					newArgs := []Value{}
					for _, arg := range args {
						newArgs = append(newArgs, Value{v: arg})
					}
					x.f(newArgs)
				}()
			}),
		}
	case nil:
		return Null()
	case bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, unsafe.Pointer, string, []byte:
		return Value{v: id.Invoke(x)}
	case []int8, []int16, []int32, []int64, []uint16, []uint32, []uint64, []float32, []float64:
		return Value{v: id.Invoke(x)}
	default:
		panic(`invalid arg: ` + reflect.TypeOf(x).String())
	}
}

func (v Value) Bool() bool {
	return v.v.Bool()
}

func convertArgs(args []interface{}) []interface{} {
	newArgs := []interface{}{}
	for _, arg := range args {
		v := ValueOf(arg)
		newArgs = append(newArgs, v.v)
	}
	return newArgs
}

func (v Value) Call(m string, args ...interface{}) Value {
	return Value{v: v.v.Call(m, convertArgs(args)...)}
}

func (v Value) Float() float64 {
	return v.v.Float()
}

func (v Value) Get(p string) Value {
	return Value{v: v.v.Get(p)}
}

func (v Value) Index(i int) Value {
	return Value{v: v.v.Index(i)}
}

func (v Value) Int() int {
	return v.v.Int()
}

func (v Value) Invoke(args ...interface{}) Value {
	return Value{v: v.v.Invoke(convertArgs(args)...)}
}

func (v Value) Length() int {
	return v.v.Length()
}

func (v Value) New(args ...interface{}) Value {
	return Value{v: v.v.New(convertArgs(args)...)}
}

func (v Value) Set(p string, x interface{}) {
	v.v.Set(p, x)
}

func (v Value) SetIndex(i int, x interface{}) {
	v.v.SetIndex(i, x)
}

func (v Value) String() string {
	return v.v.String()
}

func (v Value) InstanceOf(t Value) bool {
	return instanceOf.Invoke(v, t).Bool()
}

func GetInternalObject(v Value) interface{} {
	return v.v
}
