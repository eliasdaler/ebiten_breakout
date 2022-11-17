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
	"fmt"
	"image"
)

type Vec2i struct {
	X, Y int
}

func (v Vec2i) String() string {
	return fmt.Sprintf("(%d, %d)", v.X, v.Y)
}

type Vec2f struct {
	X, Y float64
}

func (v Vec2f) String() string {
	return fmt.Sprintf("(%.2f, %.2f)", v.X, v.Y)
}

type FloatRect struct {
	X, Y, W, H float64
}

func (r FloatRect) String() string {
	return fmt.Sprintf("(%.2f, %.2f, %.2f)", r.X, r.Y, r.W, r.H)
}

func (r FloatRect) ToImageRect() image.Rectangle {
	return image.Rect(int(r.X), int(r.Y), int(r.X+r.W), int(r.Y+r.H))
}

func (r FloatRect) Corners() (float64, float64, float64, float64) {
	return r.X, r.Y, r.X + r.W, r.Y + r.H
}
