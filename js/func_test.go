// Copyright 2019 The GopherWasm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.12 !wasm

package js_test

import (
	"testing"

	"github.com/gopherjs/gopherwasm/js"
)

func TestFuncObject(t *testing.T) {
	got := ""
	f := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		got = args[0].Get("foo").String()
		return nil
	})
	defer f.Release()

	js.ValueOf(f).Invoke(js.Global().Call("eval", `({"foo": "bar"})`))
	want := "bar"
	if got != want {
		t.Errorf("got %#v, want %#v", got, want)
	}
}

func TestValueOfFunc(t *testing.T) {
	f := js.FuncOf(func(this js.Value, args []js.Value) interface{} { return nil })
	got := js.ValueOf(f).Type()
	want := js.TypeFunction
	if got != want {
		t.Fail()
	}
}
