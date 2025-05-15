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
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"os"
	"sync"
	"sync/atomic"
)

type NodeInstance struct {
	Module api.Module
}

type NodeManager struct {
	Instances             NodeInstance
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

	if nm.Instances.Module != nil {
		err := nm.Instances.Module.Close(ctx)
		if err != nil {
			return errors.New("Error closing module:" + err.Error())
		}
	}

	mod, err := nm.runtime.InstantiateModule(ctx, compiledWasm, nm.config)
	if err != nil {
		return errors.New("Error creating instance:" + err.Error())
	}

	{
		nm.mx.Lock()
		defer nm.mx.Unlock()
		nm.Instances.Module = mod
		nm.LastUsedInstanceIndex.Store(0)
	}

	return nil
}
