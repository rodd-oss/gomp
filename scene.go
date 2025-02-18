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

package gomp

import (
	"time"
)

type SceneId uint16

type AnyScene interface {
	Init()
	Update(dt time.Duration) SceneId
	FixedUpdate(dt time.Duration)
	Render(dt time.Duration)
	Destroy()
	OnEnter()
	OnExit()
	Id() SceneId
}
