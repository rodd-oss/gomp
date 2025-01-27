/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package effects

type EffectType int

const (
	TypeDamage EffectType = iota
	TypeHeal
)

type Effect struct {
	Name         string
	Type         EffectType
	Damage       int
	FriendlyFire bool // if true, the effect will be applied to friendly units
}
