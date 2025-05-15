/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.

===-===-===-===-===-===-===-===-===-===
Donations during this file development:
-===-===-===-===-===-===-===-===-===-===

none :)

Thank you for your support!
*/

package main

import (
	"context"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"os"
	"testing"
)

func TestWASM(t *testing.T) {
	ctx := context.Background()
	r := wazero.NewRuntime(ctx)
	defer r.Close(ctx)

	// Instantiate WASI
	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	config := wazero.NewModuleConfig().
		WithStartFunctions("_start"). // Or "_initialize" based on the module type
		//WithStartFunctions("_initialize"). // Or "_initialize" based on the module type
		WithStdout(os.Stdout).
		WithStderr(os.Stderr)

	// Load WASM module
	wasm, err := os.ReadFile("hello.wasm")
	if err != nil {
		t.Fatal(err)
	}

	compiledWasm, err := r.CompileModule(ctx, wasm)
	if err != nil {
		t.Fatal(err)
	}

	mod, err := r.InstantiateModule(ctx, compiledWasm, config)
	if err != nil {
		t.Fatal(err)
	}

	// Call the function
	add := mod.ExportedFunction("add")
	_, err = add.Call(ctx, api.EncodeU32(1), api.EncodeU32(2))
	if err != nil {
		t.Fatal(err)
	}
}
