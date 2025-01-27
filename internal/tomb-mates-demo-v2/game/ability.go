/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package game

import (
	"gomp/internal/tomb-mates-demo-v2/abilities"
	"time"
)

func (g *Game) CastAbility(a *abilities.Ability) (err error) {
	// Check cooldown
	if a.PreviousCastAt != nil {
		if time.Since(*a.PreviousCastAt) <= a.Cooldown {
			return abilities.ErrorAbilityOnCooldown()
		}
	}

	// Check cost
	switch a.Cost.Type {
	case abilities.CostTypeNone:
		// No cost
	default:
		err = abilities.ErrorUnknownCostType()
	}
	if err != nil {
		return err
	}

	// Get type of ability
	switch a.Type {
	case abilities.TypeActive:
		err = g.castActive(a)
	default:
		err = abilities.ErrorUnknownAbilityType()
	}
	if err != nil {
		return err
	}

	// Success
	now := time.Now()
	a.PreviousCastAt = &now

	return nil
}

func (g *Game) castActive(a *abilities.Ability) (err error) {
	switch a.Target.Type {
	case abilities.TargetTypeArea:
		err = g.castArea(a)
	default:
		err = abilities.ErrorUnknownTargetType()
	}
	if err != nil {
		return err
	}

	return nil
}

func (g *Game) castArea(a *abilities.Ability) (err error) {
	area := a.Target.Area
	if area == nil {
		panic("No cast area object while casting in area")
	}

	return nil

}
