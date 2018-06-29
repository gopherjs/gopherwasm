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
	defer c.Close()

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
	defer c.Close()

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
