/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package game

import (
	"gomp/internal/tomb-mates-demo-v2/effects"
	"gomp/internal/tomb-mates-demo-v2/protos"
)

type AreaOfEffects struct {
	Area    *protos.Area
	Caster  *protos.Unit
	Effects []*effects.Effect
	Ttd     uint32
}

func (aoe *AreaOfEffects) AddAffectedUnit(unit *protos.Unit) {
	aoe.Area.AffectedUnitIds[unit.Id] = &protos.Empty{}
}

func (aoe *AreaOfEffects) RemoveAffectedUnit(unit *protos.Unit) {
	delete(aoe.Area.AffectedUnitIds, unit.Id)
}
