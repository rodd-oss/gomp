//go:build !renderless

/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package engine

type Render struct {
	Engine *Engine
}

func (r *Render) Update(dt float64) {
	for i := range r.Engine.LoadedScenes {
		scene := r.Engine.LoadedScenes[i]

		if !scene.ShouldRender {
			continue
		}

		// scene.Render(dt)
	}
}
