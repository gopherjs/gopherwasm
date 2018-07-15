// Copyright 2018 The GopherWasm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package js_test

import (
	"math"
	"testing"

	"github.com/gopherjs/gopherwasm/js"
)

func TestNull(t *testing.T) {
	want := "null"
	if got := js.Null().String(); got != want {
		t.Errorf("got %#v, want %#v", got, want)
	}
}

func TestCallback(t *testing.T) {
	ch := make(chan int)
	c := js.NewCallback(func(args []js.Value) {
		ch <- args[0].Int() + args[1].Int()
	})
	defer c.Release()

	js.ValueOf(c).Invoke(1, 2)
	got := <-ch
	want := 3
	if got != want {
		t.Errorf("got %#v, want %#v", got, want)
	}
}

func TestCallbackObject(t *testing.T) {
	ch := make(chan string)
	c := js.NewCallback(func(args []js.Value) {
		ch <- args[0].Get("foo").String()
	})
	defer c.Release()

	js.ValueOf(c).Invoke(js.Global().Call("eval", `({"foo": "bar"})`))
	got := <-ch
	want := "bar"
	if got != want {
		t.Errorf("got %#v, want %#v", got, want)
	}
}

func TestString(t *testing.T) {
	obj := js.Global().Call("eval", "'Hello'")
	got := obj.String()
	if want := "Hello"; got != want {
		t.Errorf("got %#v, want %#v", got, want)
	}

	obj = js.Global().Call("eval", "123")
	got = obj.String()
	if want := "123"; got != want {
		t.Errorf("got %#v, want %#v", got, want)
	}
}

func TestInt64(t *testing.T) {
	var i int64 = math.MaxInt64
	got := js.ValueOf(i).String()
	// js.Value keeps the value only in 53-bit precision.
	if want := "9223372036854776000"; got != want {
		t.Errorf("got %#v, want %#v", got, want)
	}
}

func TestInstanceOf(t *testing.T) {
	arr := js.Global().Call("eval", "[]")
	got := arr.InstanceOf(js.Global().Call("eval", "Array"))
	want := true
	if got != want {
		t.Errorf("got %#v, want %#v", got, want)
	}

	got = arr.InstanceOf(js.Global().Call("eval", "Object"))
	want = true
	if got != want {
		t.Errorf("got %#v, want %#v", got, want)
	}

	got = arr.InstanceOf(js.Global().Call("eval", "String"))
	want = false
	if got != want {
		t.Errorf("got %#v, want %#v", got, want)
	}

	str := js.Global().Call("eval", "String").New()
	got = str.InstanceOf(js.Global().Call("eval", "Array"))
	want = false
	if got != want {
		t.Errorf("got %#v, want %#v", got, want)
	}

	got = str.InstanceOf(js.Global().Call("eval", "Object"))
	want = true
	if got != want {
		t.Errorf("got %#v, want %#v", got, want)
	}

	got = str.InstanceOf(js.Global().Call("eval", "String"))
	want = true
	if got != want {
		t.Errorf("got %#v, want %#v", got, want)
	}
}

func TestTypedArrayOf(t *testing.T) {
	testTypedArrayOf(t, "[]int8", []int8{0, -42, 0}, -42)
	testTypedArrayOf(t, "[]int16", []int16{0, -42, 0}, -42)
	testTypedArrayOf(t, "[]int32", []int32{0, -42, 0}, -42)
	testTypedArrayOf(t, "[]uint8", []uint8{0, 42, 0}, 42)
	testTypedArrayOf(t, "[]uint16", []uint16{0, 42, 0}, 42)
	testTypedArrayOf(t, "[]uint32", []uint32{0, 42, 0}, 42)
	testTypedArrayOf(t, "[]float32", []float32{0, -42.5, 0}, -42.5)
	testTypedArrayOf(t, "[]float64", []float64{0, -42.5, 0}, -42.5)
}

func testTypedArrayOf(t *testing.T, name string, slice interface{}, want float64) {
	t.Run(name, func(t *testing.T) {
		a := js.TypedArrayOf(slice)
		got := a.Index(1).Float()
		a.Release()
		if got != want {
			t.Errorf("got %#v, want %#v", got, want)
		}

		a2 := js.TypedArrayOf(slice)
		v := js.ValueOf(a2)
		got = v.Index(1).Float()
		a2.Release()
		if got != want {
			t.Errorf("got %#v, want %#v", got, want)
		}
	})
}

func TestType(t *testing.T) {
	if got, want := js.Undefined().Type(), js.TypeUndefined; got != want {
		t.Errorf("got %s, want %s", got, want)
	}
	if got, want := js.Null().Type(), js.TypeNull; got != want {
		t.Errorf("got %s, want %s", got, want)
	}
	if got, want := js.ValueOf(true).Type(), js.TypeBoolean; got != want {
		t.Errorf("got %s, want %s", got, want)
	}
	if got, want := js.ValueOf(42).Type(), js.TypeNumber; got != want {
		t.Errorf("got %s, want %s", got, want)
	}
	if got, want := js.ValueOf("test").Type(), js.TypeString; got != want {
		t.Errorf("got %s, want %s", got, want)
	}
	if got, want := js.Global().Get("Symbol").Invoke("test").Type(), js.TypeSymbol; got != want {
		t.Errorf("got %s, want %s", got, want)
	}
	if got, want := js.Global().Get("Array").New().Type(), js.TypeObject; got != want {
		t.Errorf("got %s, want %s", got, want)
	}
	if got, want := js.Global().Get("Array").Type(), js.TypeFunction; got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
