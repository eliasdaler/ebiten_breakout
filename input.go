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

// GetStickState returns a virtual stick state for button controls
// E.g. pressing only W returns (0, -1), pressing W and D returns (1, -1)
// Pressing opposite directions cancels out movement, so pressing W and S returns (0, 0)
func GetStickState() Vec2f {
	stick := Vec2f{}

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		stick.X += -1.
	}

	if ebiten.IsKeyPressed(ebiten.KeyD) {
		stick.X += 1.
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		stick.Y += -1.
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) {
		stick.Y += 1.
	}

	return stick
}
