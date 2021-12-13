package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	Texture          *ebiten.Image // full texture
	Image            *ebiten.Image // displayed portion of image
	MirrorX, MirrorY bool
}

func NewSprite(img *ebiten.Image) *Sprite {
	return &Sprite{
		Texture: img,
		Image:   img,
	}
}

func (s *Sprite) SetTextureRect(r image.Rectangle) {
	s.Image = s.Texture.SubImage(r).(*ebiten.Image)
}

func (s *Sprite) GetTextureRect() image.Rectangle {
	return s.Image.Bounds()
}

func (s *Sprite) GetSize() Vec2f {
	b := s.Image.Bounds()
	return Vec2f{float64(b.Dx()), float64(b.Dy())}
}

func (s *Sprite) Draw(screen *ebiten.Image, parentM ebiten.GeoM) {
	scale := Vec2f{1., 1.}
	offset := Vec2f{}
	if s.MirrorX {
		scale.X = -1
		offset.X = float64(s.GetTextureRect().Dx())
	}
	if s.MirrorY {
		scale.Y = -1
		offset.Y = float64(s.GetTextureRect().Dy())
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale.X, scale.Y)
	op.GeoM.Translate(offset.X, offset.Y)
	op.GeoM.Concat(parentM)

	screen.DrawImage(s.Image, op)
}
