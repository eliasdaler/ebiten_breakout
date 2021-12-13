package main

type Paddle struct {
	Entity
}

func makePaddle() (*Paddle, error) {
	e := &Paddle{
		Entity: NewEntity(),
	}

	e.Movement.Speed = Vec2f{90, 0}

	img, _, err := NewImageFromFile("data/images/paddle.png")
	if err != nil {
		return nil, err
	}

	e.Graphics.Sprite = NewSprite(img)

	return e, nil
}
