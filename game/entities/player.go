package entities

import "roguelike/animation"

type PlayerState uint8

const (
	Down PlayerState = iota
	Up
	Left
	Right
)

type Player struct {
	*Sprite
	HP, MaxHP uint
	Animation map[PlayerState]*animation.Animation
}

func (p *Player) ActiveAnimation(dx, dy int) *animation.Animation {
	if dx > 0 {
		return p.Animation[Right]
	} else if dx < 0 {
		return p.Animation[Left]
	} else if dy > 0 {
		return p.Animation[Down]
	} else if dy < 0 {
		return p.Animation[Up]
	}
	return nil
}
