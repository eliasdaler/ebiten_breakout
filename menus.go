// Copyright (c) 2022 Elias Daler
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of version 3 of the GNU General Public
// License as published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
// General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"image/color"
	"io"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

// dialogue box position
var dbPos = Vec2i{64, 42}
var fontSize = 8

var (
	menuBGColor     = color.RGBA{77, 170, 182, 200}
	textNormalColor = color.RGBA{34, 32, 32, 255}
	textRedColor    = color.RGBA{172, 50, 50, 255}
	textGreenColor  = color.RGBA{51, 125, 60, 255}
)

type MenuState int

const (
	MenuMain MenuState = iota
	MenuGameOver
	MenuWin
)

type textItem struct {
	text  string
	posY  int
	color color.Color
}

var mainMenuTexts = []textItem{
	{"Ebiten Breakout", dbPos.Y + 18, textRedColor},
	{"Press Z to play", dbPos.Y + 34, textNormalColor},
	{"W - launch ball", dbPos.Y + 50, textNormalColor},
	{"A/D - move", dbPos.Y + 66, textNormalColor},
	{"2021, Elias Daler", 236, textNormalColor},
}

var gameOverTexts = []textItem{
	{"panic:game over", dbPos.Y + 18, textRedColor},
	{"Press Z to", dbPos.Y + 50, textNormalColor},
	{"recover", dbPos.Y + 66, textGreenColor},
}

var wonTexts = []textItem{
	{"YOU WON!", dbPos.Y + 18, textGreenColor},
	{"Press Z to", dbPos.Y + 50, textNormalColor},
	{"play again", dbPos.Y + 66, textNormalColor},
}

func (g *Game) loadMenuResources() error {
	// load font
	f, err := OpenFile("data/fonts/pressstart2p.ttf")
	if err != nil {
		log.Fatal(err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	tt, err := opentype.Parse(b)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72
	font, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    float64(fontSize),
		DPI:     dpi,
		Hinting: font.HintingNone,
	})
	if err != nil {
		return err
	}
	g.font = font

	// load bg
	img, _, err := NewImageFromFile("data/images/bg.png")
	if err != nil {
		return err
	}
	g.bg = NewSprite(img)

	// load dialogue box
	img, _, err = NewImageFromFile("data/images/dialogue_box.png")
	if err != nil {
		return err
	}
	g.dialogueBox = NewSprite(img)

	return nil
}

func (g *Game) drawMenu(screen *ebiten.Image) {
	ebitenutil.DrawRect(screen, 0, 0, gameScreenWidth, gameScreenHeight, menuBGColor)

	m := ebiten.GeoM{}
	m.Translate(float64(dbPos.X), float64(dbPos.Y))
	g.dialogueBox.Draw(screen, m)

	var texts []textItem
	switch g.menuState {
	case MenuMain:
		texts = mainMenuTexts
	case MenuGameOver:
		texts = gameOverTexts
	case MenuWin:
		texts = wonTexts
	default:
		panic("unexpected state")
	}

	for _, ti := range texts {
		text.Draw(screen, ti.text, g.font, gameScreenWidth/2.-len(ti.text)/2.*fontSize, ti.posY, ti.color)
	}
}
