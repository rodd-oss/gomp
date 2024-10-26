/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package client

import (
	e "github.com/hajimehoshi/ebiten/v2"
)

type Config struct {
	Width        int
	Height       int
	Title        string
	ResizingMode e.WindowResizingModeType
	EnableDebug  bool
}

func NewConfig(width int, height int, title string, resizingMode e.WindowResizingModeType, enableDebug bool) *Config {
	return &Config{
		Width:        width,
		Height:       height,
		Title:        title,
		ResizingMode: resizingMode,
		EnableDebug:  enableDebug,
	}
}
