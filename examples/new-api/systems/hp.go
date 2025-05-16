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

package systems

import (
	"gomp/examples/new-api/components"
	"gomp/pkg/ecs"
	"time"
)

func NewHpSystem() HpSystem {
	return HpSystem{}
}

type HpSystem struct {
	EntityManager        *ecs.EntityManager
	Hps                  *components.HpComponentManager
	Asteroids            *components.AsteroidComponentManager
	Players              *components.PlayerTagComponentManager
	Hp                   *components.HpComponentManager
	AsteroidSceneManager *components.AsteroidSceneManagerComponentManager
}

func (s *HpSystem) Init() {}
func (s *HpSystem) Run(dt time.Duration) {
	s.Hps.EachEntity(func(e ecs.Entity) bool {
		hp := s.Hps.GetUnsafe(e)

		if hp.Hp <= 0 {
			asteroid := s.Asteroids.GetUnsafe(e)
			player := s.Players.GetUnsafe(e)
			s.AsteroidSceneManager.EachComponent(func(a *components.AsteroidSceneManager) bool {
				if asteroid != nil {
					a.PlayerScore += hp.MaxHp
				}
				if player != nil {
					playerHp := s.Hp.GetUnsafe(e)
					a.PlayerHp = playerHp.Hp
				}
				return false
			})

			s.EntityManager.Delete(e)
		}
		return true
	})
}
func (s *HpSystem) Destroy() {}
