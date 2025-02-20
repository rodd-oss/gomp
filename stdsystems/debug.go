/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package stdsystems

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func NewDebugSystem() DebugSystem {
	return DebugSystem{}
}

type DebugSystem struct {
	pprofEnabled bool
}

func (s *DebugSystem) Init() {}
func (s *DebugSystem) Run(dt time.Duration) {
	if rl.IsKeyPressed(rl.KeyF9) {
		if s.pprofEnabled {
			pprof.StopCPUProfile()
			fmt.Println("CPU Profile Stopped")
		} else {
			f, err := os.Create("cpu.out")
			if err != nil {
				log.Fatal(err)
			}
			pprof.StartCPUProfile(f)
			fmt.Println("CPU Profile Started")
		}

		s.pprofEnabled = !s.pprofEnabled
	}
}
func (s *DebugSystem) Destroy() {}
