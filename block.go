package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Block struct {
	Entity

	IsBroken bool
}

func makeBlock(pos Vec2f, img *ebiten.Image) (*Block, error) {
	e := &Block{
		Entity: NewEntity(),
	}

	e.Transform.Position = pos

	e.Graphics.Sprite = NewSprite(img)
	return e, nil
}
