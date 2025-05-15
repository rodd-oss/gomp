/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.

===-===-===-===-===-===-===-===-===-===
Donations during this file development:
-===-===-===-===-===-===-===-===-===-===

<- HromRU Donated 2 500 RUB
<- Еблан Donated 228 RUB

Thank you for your support!
*/

package main

import (
	"context"
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
	"time"
)

func main() {
	var ctx = context.Background()
	var nodeManager = NodeManager{}
	var game = wasyan.Game{
		Position: vectors.Vec2{X: 1, Y: 2},
		Velocity: vectors.Vec2{X: 3, Y: 4},
	}

	nodeManager.runtime = wazero.NewRuntime(ctx)
	defer nodeManager.runtime.Close(ctx)

	//var sizeOfGame = int(unsafe.Sizeof(game))
	//var getGameModule = Module{
	//	Fn: api.GoModuleFunc(func(ctx context.Context, mod api.Module, stack []uint64) {
	//		gameBytes := unsafe.Slice((*byte)(unsafe.Pointer(&game)), sizeOfGame)
	//		gameRef := api.DecodeU32(stack[0])
	//		mod.Memory().Write(gameRef, gameBytes)
	//	}),
	//	Params:  []api.ValueType{api.ValueTypeI32},
	//	Results: []api.ValueType{},
	//}
	//
	//var setGameModule = Module{
	//	Fn: api.GoModuleFunc(func(ctx context.Context, mod api.Module, stack []uint64) {
	//		gameRef := api.DecodeU32(stack[0])
	//		r, _ := mod.Memory().Read(gameRef, uint32(unsafe.Sizeof(game)))
	//		localGame := (*wasyan.Game)(unsafe.Pointer(unsafe.SliceData(r)))
	//		game = *localGame
	//	}),
	//	Params:  []api.ValueType{api.ValueTypeI32},
	//	Results: []api.ValueType{},
	//}

	_, err := nodeManager.runtime.NewHostModuleBuilder("env").
		//NewFunctionBuilder().
		//WithGoModuleFunction(getGameModule.Fn, getGameModule.Params, getGameModule.Results).
		//Export("get_game").
		//NewFunctionBuilder().
		//WithGoModuleFunction(setGameModule.Fn, setGameModule.Params, setGameModule.Results).
		//Export("set_game").
		Instantiate(ctx)
	if err != nil {
		log.Panicln(err)
	}

	wasi_snapshot_preview1.MustInstantiate(ctx, nodeManager.runtime)

	// Configure the module to initialize the reactor.
	nodeManager.config = wazero.NewModuleConfig().
		WithStdout(os.Stdout).
		WithStderr(os.Stderr)

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
		log.Println("Instance updated")

		return c.SendString("Instance updated")
	})

	app.Get("/internal/add", func(c *fiber.Ctx) error {
		a, err := strconv.Atoi(c.Query("a"))
		if err != nil {
			return c.SendString("Error parsing a:" + err.Error())
		}
		b, err := strconv.Atoi(c.Query("b"))
		if err != nil {
			return c.SendString("Error parsing b:" + err.Error())
		}

		start := time.Now()
		wasyan.UpdateGame(&game, vectors.Vec2{X: float32(a), Y: float32(b)})
		duration := time.Since(start).String()

		return c.SendString("Result is " + fmt.Sprint(game) + "\n" + duration)
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

		instance := &nodeManager.Instances
		module := instance.Module
		if module == nil {
			return c.SendString("Module not found")
		}

		fn := module.ExportedFunction(name)
		if fn == nil {
			return c.SendString("RPC not found")
		}

		start := time.Now()
		results, err := fn.Call(c.Context(), api.EncodeU32(uint32(a)), api.EncodeU32(uint32(b)))
		duration := time.Since(start).String()
		if err != nil {
			log.Println(err)
			return c.SendString("Error calling RPC\n" + err.Error())
		}

		return c.SendString("Result is " + fmt.Sprint(results) + "\n" + "Game is " + fmt.Sprint(game) + "\n" + duration)
	})

	log.Fatal(app.Listen(":3000"))
}
