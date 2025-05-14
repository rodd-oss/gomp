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
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"gomp/internal/wasyan"
	"gomp/vectors"
	"log"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

type NodeInstance struct {
	Module api.Module
}

type NodeManager struct {
	Instances             [2]NodeInstance
	LastUsedInstanceIndex atomic.Int32
	config                wazero.ModuleConfig
	runtime               wazero.Runtime
	mx                    sync.Mutex
}

func (nm *NodeManager) LoadWasm(ctx context.Context, file string) error {
	wasm, err := os.ReadFile(file)
	if err != nil {
		return errors.New("Error opening file:" + err.Error())
	}

	compiledWasm, err := nm.runtime.CompileModule(ctx, wasm)
	if err != nil {
		return errors.New("Error compiling module:" + err.Error())
	}

	instance1, err := nm.runtime.InstantiateModule(ctx, compiledWasm, nm.config)
	if err != nil {
		return errors.New("Error creating instance1:" + err.Error())
	}

	instance2, err := nm.runtime.InstantiateModule(ctx, compiledWasm, nm.config)
	if err != nil {
		return errors.New("Error creating instance2:" + err.Error())
	}

	{
		nm.mx.Lock()
		defer nm.mx.Unlock()
		nm.Instances[0].Module = instance1
		nm.Instances[1].Module = instance2
		nm.LastUsedInstanceIndex.Store(0)
	}

	return nil
}

func main() {
	var ctx = context.Background()
	var nodeManager = NodeManager{}
	var game = wasyan.Game{
		Position: vectors.Vec2{X: 0, Y: 0},
		Velocity: vectors.Vec2{X: 50, Y: 10},
	}

	nodeManager.runtime = wazero.NewRuntime(ctx)
	defer nodeManager.runtime.Close(ctx)

	_, err := nodeManager.runtime.NewHostModuleBuilder("env").
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context) uint64 {
			return api.EncodeExternref(uintptr(unsafe.Pointer(&game)))
		}).
		Export("get_game").
		Instantiate(ctx)
	if err != nil {
		log.Panicln(err)
	}

	wasi_snapshot_preview1.MustInstantiate(ctx, nodeManager.runtime)

	// Configure the module to initialize the reactor.
	nodeManager.config = wazero.NewModuleConfig().WithStartFunctions("_initialize")

	var app = fiber.New()

	//app.Use(recover.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/update/:file", func(c *fiber.Ctx) error {
		file := c.Params("file")
		err := nodeManager.LoadWasm(c.Context(), file)
		if err != nil {
			return c.SendString("Error loading module:\n" + err.Error())
		}

		return c.SendString("Instance updated")
	})

	app.Get("/add", func(c *fiber.Ctx) error {
		a, err := strconv.Atoi(c.Query("a"))
		if err != nil {
			return c.SendString("Error parsing a:" + err.Error())
		}
		b, err := strconv.Atoi(c.Query("b"))
		if err != nil {
			return c.SendString("Error parsing b:" + err.Error())
		}

		start := time.Now()
		result := wasyan.Add(int32(a), int32(b))
		duration := time.Since(start).String()

		return c.SendString("Result is " + fmt.Sprint(result) + "\n" + duration)
	})

	app.Get("/call/:name", func(c *fiber.Ctx) error {
		name := c.Params("name")
		a, err := strconv.Atoi(c.Query("a"))
		if err != nil {
			return c.SendString("Error parsing a:" + err.Error())
		}
		b, err := strconv.Atoi(c.Query("b"))
		if err != nil {
			return c.SendString("Error parsing b:" + err.Error())
		}

		index := nodeManager.LastUsedInstanceIndex.Add(1) % int32(len(nodeManager.Instances))
		instance := &nodeManager.Instances[index]
		module := instance.Module
		if module == nil {
			return c.SendString("Module not found")
		}

		fn := module.ExportedFunction(name)
		if fn == nil {
			return c.SendString("RPC not found")
		}

		start := time.Now()
		results, err := fn.Call(c.Context(), api.EncodeI32(int32(a)), api.EncodeI32(int32(b)))
		duration := time.Since(start).String()
		if err != nil {
			return c.SendString("Error calling RPC" + err.Error())
		}
		result := (*float32)(unsafe.Pointer(api.DecodeExternref(results[0])))
		game.Velocity.X += *result

		return c.SendString("Result is " + fmt.Sprint(*result) + "\n" + "Game is " + fmt.Sprint(game) + "\n" + duration)
	})

	log.Fatal(app.Listen(":3000"))
}
