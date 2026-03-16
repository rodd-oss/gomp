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
	"fmt"
	"github.com/bytecodealliance/wasmtime-go/v32"
	"log"
	"os"
)

func main() {
	engine := wasmtime.NewEngine()
	defer engine.Close()
	store := wasmtime.NewStore(engine)
	defer store.Close()

	// Create a linker with WASI functions defined within it
	linker := wasmtime.NewLinker(engine)
	err := linker.DefineWasi()
	check(err)

	wasiConfig := wasmtime.NewWasiConfig()
	defer wasiConfig.Close()

	wasiConfig.InheritStderr()
	wasiConfig.InheritStdout()

	store.SetWasi(wasiConfig)

	wasm, err := os.ReadFile("hello.wasm")
	check(err)

	module, err := wasmtime.NewModule(engine, wasm)
	check(err)
	defer module.Close()

	instance, err := linker.Instantiate(store, module)
	check(err)

	//get_game := wasmtime.WrapFunc(store, func(gameRef uint32) {
	//
	//})
	//
	//set_game := wasmtime.WrapFunc(store, func(gameRef uint32) {
	//
	//})
	//instance, err := wasmtime.NewInstance(store, module, []wasmtime.AsExtern{get_game, set_game})

	// Init
	initFn := instance.GetFunc(store, "_start")
	if initFn == nil {
		initFn = instance.GetFunc(store, "_initialize")
		if initFn == nil {
			panic("Init function not found")
		}
	}

	_, err = initFn.Call(store)
	check(err)
	log.Println("Initialized")

	//call
	gcd := instance.GetFunc(store, "add")
	val, err := gcd.Call(store, 6, 27)
	check(err)
	fmt.Printf("add(6, 27) = %d\n", val.(int32))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
