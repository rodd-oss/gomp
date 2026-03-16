/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package stdsystems

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime/pprof"
	"runtime/trace"

	"github.com/felixge/fgprof"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const gpprof = false

func NewDebugSystem() DebugSystem {
	return DebugSystem{}
}

type DebugSystem struct {
	pprofEnabled bool
	traceEnabled bool
	traceCounter int
}

func (s *DebugSystem) Init() {
	if gpprof {
		http.DefaultServeMux.Handle("/debug/fgprof", fgprof.Handler())
		go func() {
			log.Println(http.ListenAndServe(":6060", nil))
		}()
	}

}
func (s *DebugSystem) Run() {
	if rl.IsKeyPressed(rl.KeyF9) {
		if s.pprofEnabled {
			pprof.StopCPUProfile()
			log.Println("CPU Profile Stopped")

			// Create a memory profile file
			memProfileFile, err := os.Create("mem.out")
			if err != nil {
				panic(err)
			}
			defer memProfileFile.Close()

			// Write memory profile to file
			if err := pprof.WriteHeapProfile(memProfileFile); err != nil {
				panic(err)
			}
			log.Println("Memory profile written to mem.prof")
		} else {
			f, err := os.Create("cpu.out")
			if err != nil {
				log.Fatal(err)
			}

			err = pprof.StartCPUProfile(f)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("CPU Profile Started")
		}

		s.pprofEnabled = !s.pprofEnabled
	}

	if rl.IsKeyPressed(rl.KeyF10) {
		if !s.traceEnabled {
			f, err := os.Create("trace.out")
			if err != nil {
				log.Fatal(err)
			}

			err = trace.Start(f)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("Trace Profile Started")

			s.traceEnabled = true
		}
	}

	if s.traceEnabled {
		s.traceCounter++
		if s.traceCounter == 50 {
			trace.Stop()
			s.traceEnabled = false
			s.traceCounter = 0
			log.Println("Trace Profile Stopped")
		}
	}
}
func (s *DebugSystem) Destroy() {}
