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
	"github.com/hajimehoshi/go-steamworks"
	"golang.org/x/text/language"
	"gomp"
	"gomp/examples/new-api/scenes"
	"os"
)

const appID = 12 // Rewrite this

func initd() {
	if steamworks.RestartAppIfNecessary(appID) {
		os.Exit(1)
	}
	err := steamworks.Init()
	if err != nil {
		panic("steamworks.Init failed")
	}
}

func SystemLang() language.Tag {
	switch steamworks.SteamApps().GetCurrentGameLanguage() {
	case "russian":
		return language.Russian
	case "english":
		return language.English
	case "japanese":
		return language.Japanese
	}
	return language.Und
}

func main() {
	//log.Println(SystemLang())
	initialScene := scenes.NewAssteroddScene()

	game := NewGame(&initialScene)

	engine := gomp.NewEngine(&game)
	engine.Run(20, 0)
}
