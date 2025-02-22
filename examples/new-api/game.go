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

package main

import (
	"gomp"
	"gomp/examples/new-api/scenes"
)

func main() {
	sceneList := scenes.NewSceneList()

	game := gomp.NewGame(
		&sceneList.Main,
	)
	game.CurrentSceneId = scenes.MainSceneId

	engine := gomp.NewEngine(&game)
	engine.Run(50, 165)
}
