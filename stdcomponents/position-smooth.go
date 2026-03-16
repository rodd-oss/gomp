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

package stdcomponents

import (
	"gomp/pkg/ecs"
)

const (
	TexturePositionSmoothOff TexturePositionSmooth = iota
	TexturePositionSmoothLerp
	TexturePositionSmoothExpDecay
)

// TexturePositionSmooth is the component tag for stdsystems.TexturePositionSmoothSystem
// TODO: refactor or make stable realization
type TexturePositionSmooth uint8

type TexturePositionSmoothComponentManager = ecs.ComponentManager[TexturePositionSmooth]

func NewTexturePositionSmoothComponentManager() TexturePositionSmoothComponentManager {
	return ecs.NewComponentManager[TexturePositionSmooth](TexturePositionSmoothComponentId)
}
