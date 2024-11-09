package gomp

import (
	"image"

	"github.com/jakecoffman/cp/v2"
)

var BodyComponent = CreateComponent[cp.Body]()

type SpriteData struct {
	Image  image.Image
	Config image.Config
}

var SpriteComponent = CreateComponent[SpriteData]()
