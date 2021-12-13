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
