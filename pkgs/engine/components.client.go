//go:build !renderless

/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package engine

import (
	"github.com/hajimehoshi/ebiten/v2"
	ecs "github.com/yohamta/donburi"
)

type drawableController interface {
	Init()
	Update(dt float64)
	Draw(screen *ebiten.Image)
}

func CreateComponentDrawable[T drawableController](opts ...T) *ecs.ComponentType[T] {
	return ecs.NewComponentType[T](opts)
}
