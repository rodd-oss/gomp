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

package config

import "gomp/stdcomponents"

const (
	TickRate  = 20
	FrameRate = 0
)

const (
	DefaultCollisionLayer stdcomponents.CollisionLayer = iota
	PlayerCollisionLayer
	BulletCollisionLayer
	EnemyCollisionLayer
	WallCollisionLayer
)

const (
	MainCameraLayer stdcomponents.CameraLayer = 1 << iota
	DebugLayer
	MinimapCameraLayer
)
