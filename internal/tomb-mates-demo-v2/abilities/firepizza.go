/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package abilities

import (
	"gomp/internal/tomb-mates-demo-v2/effects"
	"gomp/internal/tomb-mates-demo-v2/protos"
	"time"
)

type FirePizza = Ability

var damageEffect = &effects.Effect{
	Name:         "Fire Pizza Damage",
	Type:         effects.TypeDamage,
	Damage:       10,
	FriendlyFire: false,
}

func NewFirePizza(caster *protos.Unit) *FirePizza {
	return &FirePizza{
		Name:        "Fire Pizza",
		Description: "Cast fire pizza on the gorund to burn your enemies!",
		Caster:      caster,
		Type:        TypeActive,
		Cooldown:    time.Second * 3,
		Cost: AbilityCost{
			Type:  CostTypeNone,
			Value: 0,
		},
		Target: AbilityTarget{
			Type: TargetTypeArea,
			Area: &TargetArea{
				Range:  100,
				Radius: 20,
			},
		},
		Effects: []*effects.Effect{
			damageEffect,
		},
	}
}
