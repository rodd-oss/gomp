/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package engine

import (
	"github.com/yohamta/donburi"
)

// TODO: implement on go 1.24+
// type componentSyncable[T networkController[S], S any] Œ= *ecs.ComponentType[T]

// type component[T controller] *ecs.ComponentType[T]

type IComponent = donburi.IComponentType

func CreateComponent[T any](initData T) *donburi.ComponentType[T] {
	return donburi.NewComponentType[T](initData)
}
