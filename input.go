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
