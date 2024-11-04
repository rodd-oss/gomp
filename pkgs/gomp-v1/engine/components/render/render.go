/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package render

import (
	"gomp_game/pkgs/gomp-v1/engine"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/setanarut/anim"
)

type RendererData struct {
	spriteSheet *ebiten.Image
	AnimPlayer  *anim.AnimationPlayer
}

var RenderComponent = engine.CreateComponent[RendererData]
