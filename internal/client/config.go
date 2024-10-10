package client

import (
	e "github.com/hajimehoshi/ebiten/v2"
)

type Config struct {
	Width        int
	Height       int
	Title        string
	ResizingMode e.WindowResizingModeType
}

func NewConfig(width int, height int, title string, resizingMode e.WindowResizingModeType) *Config {
	return &Config{
		Width:        width,
		Height:       height,
		Title:        title,
		ResizingMode: resizingMode,
	}
}
