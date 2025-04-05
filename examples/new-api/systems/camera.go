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
	"gomp/stdcomponents"
	"time"
)

// CameraSystem is simple system responsible for managing the camera entities in the game.
type CameraSystem struct {
	EntityManager           *ecs.EntityManager
	MainCameraComponents    *stdcomponents.MainCameraComponentManager
	PipCameraComponents     *stdcomponents.PipCameraComponentManager
	MinimapCameraComponents *stdcomponents.MinimapCameraComponentManager
	Player                  *components.PlayerTagComponentManager
	Position                *stdcomponents.PositionComponentManager

	mainCamera             ecs.Entity
	pipCamera              ecs.Entity
	minimapCamera          ecs.Entity
	mainCameraComponent    *stdcomponents.Camera
	pipCameraComponent     *stdcomponents.PipCamera
	minimapCameraComponent *stdcomponents.MinimapCamera
}

func NewCameraSystem() CameraSystem {
	return CameraSystem{}
}

func (s *CameraSystem) Init() {
	s.mainCamera = s.EntityManager.Create()
	s.mainCameraComponent = s.MainCameraComponents.Create(s.mainCamera, stdcomponents.Camera{})

	s.pipCamera = s.EntityManager.Create()
	s.pipCameraComponent = s.PipCameraComponents.Create(s.pipCamera, stdcomponents.PipCamera{})

	s.minimapCamera = s.EntityManager.Create()
	s.minimapCameraComponent = s.MinimapCameraComponents.Create(s.minimapCamera, stdcomponents.MinimapCamera{})
}

func (s *CameraSystem) Run(dt time.Duration) {
	// Update camera positions or other logic here
	// For example, you might want to move the main camera based on player input
	// or other game events.
	// This is just a placeholder for the actual camera logic.

	// Follow player
	s.Player.EachEntity(func(entity ecs.Entity) bool {
		playerPosition := s.Position.Get(entity)
		s.mainCameraComponent.Target = playerPosition.XY
		s.minimapCameraComponent.Target = playerPosition.XY
		return false
	})

}

func (s *CameraSystem) Destroy() {
	s.MainCameraComponents.Remove(s.mainCamera)
	s.PipCameraComponents.Remove(s.pipCamera)
	s.MinimapCameraComponents.Remove(s.minimapCamera)
}
