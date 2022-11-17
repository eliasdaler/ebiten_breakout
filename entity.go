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

import "github.com/hajimehoshi/ebiten/v2"

type Entity struct {
	Transform *TransformComponent
	Movement  *MovementComponent
	Graphics  *GraphicsComponent
}

func NewEntity() Entity {
	return Entity{
		Transform: &TransformComponent{},
		Movement:  &MovementComponent{},
		Graphics:  &GraphicsComponent{},
	}
}

func (e *Entity) getAABB() FloatRect {
	w := e.Graphics.Sprite.GetSize().X
	h := e.Graphics.Sprite.GetSize().Y
	return FloatRect{e.Transform.Position.X, e.Transform.Position.Y, w, h}
}

func (e *Entity) Draw(screen *ebiten.Image) {
	var m ebiten.GeoM
	m.Translate(
		e.Transform.Position.X,
		e.Transform.Position.Y,
	)
	e.Graphics.Sprite.Draw(screen, m)
}
