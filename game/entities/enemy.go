package entities

type Enemy struct {
	*Sprite
	FollowPlayer bool
}

func (e *Enemy) Update(player *Sprite) {
	if !e.FollowPlayer {
		return
	}
	if e.X < player.X {
		e.X += 1
	} else {
		e.X -= 1
	}
	if e.Y < player.Y {
		e.Y += 1
	} else {
		e.Y -= 1
	}

	// clamp to screen bounds
	if e.X < 0 {
		e.X = 0
	}
	if e.X > 640 {
		e.X = 640
	}
	if e.Y < 0 {
		e.Y = 0
	}
	if e.Y > 480 {
		e.Y = 480
	}
}
