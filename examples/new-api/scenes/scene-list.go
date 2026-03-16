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

package scenes

func NewSceneList() SceneList {
	return SceneList{
		Main:      NewMainScene(),
		Assterodd: NewAssteroddScene(),
	}
}

type SceneList struct {
	Main      MainScene
	Assterodd AssteroddScene
}
