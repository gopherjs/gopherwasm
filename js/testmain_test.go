// Copyright 2019 The GopherWasm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !wasm

package js_test

import (
	"testing"
)

func TestMain(m *testing.M) {
	// Suppress the 'deadlock' error on GopherJS by goroutine
	// (https://github.com/gopherjs/gopherjs/issues/826).
	go func() {
		m.Run()
	}()
}
