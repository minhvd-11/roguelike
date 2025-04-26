package entities

type Enemy struct {
	*Sprite
	FollowPlayer bool
}

func (e *Enemy) Update(p *Player) {
	if !e.FollowPlayer {
		return
	}
	if e.X < p.X {
		e.Dx += 1
	} else if e.X > p.X {
		e.Dx -= 1
	}
	if e.Y < p.Y {
		e.Dy += 1
	} else if e.Y > p.Y {
		e.Dy -= 1
	}
}
