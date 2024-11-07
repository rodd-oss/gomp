package gomp

import (
	"image"

	"github.com/jakecoffman/cp/v2"
)

var BodyComponent = CreateComponent[cp.Body]()

type SpriteData struct {
	Image image.Image
}

var SpriteComponent = CreateComponent[SpriteData]()
