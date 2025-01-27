/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package systems

import (
	"fmt"
	"gomp/pkgs/ecs"
	"log"
	"os"
	"runtime/pprof"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type debugController struct {
	pprofEnabled bool
}

func (s *debugController) Init(world *ecs.World) {}
func (s *debugController) Update(world *ecs.World) {
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
func (s *debugController) FixedUpdate(world *ecs.World) {}
func (s *debugController) Destroy(world *ecs.World)     {}
