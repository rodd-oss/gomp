/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.

===-===-===-===-===-===-===-===-===-===
Donations during this file development:
-===-===-===-===-===-===-===-===-===-===

<- Basserti Donated 250 RUB
<- ubiblasosnalim Donated 100 RUB

Thank you for your support!
*/

package ecs

import "github.com/negrel/assert"

type EntityVersion uint

type entityType = uint32
type Entity entityType

func (e *Entity) IsVersion(version EntityVersion) bool {
	return e.GetVersion() == version
}

func (e *Entity) SetVersion(version EntityVersion) {
	assert.True(version <= MaxEntityGenerationId, "version is too high")
	*e = Entity(entityType(*e) - entityType(e.GetVersion()<<(entityPower-generationPower)) | entityType(version)<<(entityPower-generationPower))
}

func (e *Entity) GetVersion() EntityVersion {
	return EntityVersion(*e >> (entityPower - generationPower))
}

const (
	entityPower                         = 32
	generationPower                     = 2
	MaxEntityGenerationId EntityVersion = 1<<generationPower - 1
	NumOfGenerations                    = MaxEntityGenerationId + 1
	MaxEntityId           Entity        = 1<<(entityPower-generationPower) - 1
	ent                   Entity        = 35 | 3<<(entityPower-generationPower)
	ent2                  Entity        = 3221225507 - 3<<(entityPower-generationPower) | 1<<(entityPower-generationPower)
)
