package main

type TransformComponent struct {
	Position  Vec2f
	Direction Vec2f
}

type MovementComponent struct {
	Speed    Vec2f
	Velocity Vec2f
}

type GraphicsComponent struct {
	Sprite *Sprite
}
