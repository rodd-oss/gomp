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

package ecs

import "github.com/negrel/assert"

type EntityVersion uint

type entityType = uint32
type Entity entityType

func (e *Entity) IsVersion(version EntityVersion) bool {
	return e.GetVersion() == version
}

func (e *Entity) SetVersion(version EntityVersion) {
	assert.True(version <= MaxEntityVersionId, "version is too high")
	*e = Entity(entityType(*e) - entityType(e.GetVersion()<<(entityPower-versionPower)) | entityType(version)<<(entityPower-versionPower))
}

func (e *Entity) GetVersion() EntityVersion {
	return EntityVersion(*e >> (entityPower - versionPower))
}

const (
	entityPower                      = 32
	versionPower                     = 2
	MaxEntityVersionId EntityVersion = 1<<versionPower - 1
	MaxEntityId        Entity        = 1<<(entityPower-versionPower) - 1
	ent                Entity        = 35 | 3<<(entityPower-versionPower)
	ent2               Entity        = 3221225507 - 3<<(entityPower-versionPower) | 1<<(entityPower-versionPower)
)
