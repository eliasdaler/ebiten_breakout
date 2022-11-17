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

// Checks if two rects collide and returns a depth of intersection
func getIntersectionDepth(rectA, rectB FloatRect) Vec2f {
	if rectA.W == 0. || rectA.H == 0. ||
		rectB.W == 0. || rectB.H == 0. { // rectA or rectB are zero-sized
		return Vec2f{}
	}

	rectARight := rectA.X + rectA.W
	rectADown := rectA.Y + rectA.H

	rectBRight := rectB.X + rectB.W
	rectBDown := rectB.Y + rectB.H

	// check if rects intersect
	if rectA.X >= rectBRight || rectARight <= rectB.X ||
		rectA.Y >= rectBDown || rectADown <= rectB.Y {
		return Vec2f{}
	}

	var depth Vec2f

	// x
	if rectA.X < rectB.X {
		depth.X = rectARight - rectB.X
	} else {
		depth.X = rectA.X - rectBRight
	}

	// y
	if rectA.Y < rectB.Y {
		depth.Y = rectADown - rectB.Y
	} else {
		depth.Y = rectA.Y - rectBDown
	}

	return depth
}
