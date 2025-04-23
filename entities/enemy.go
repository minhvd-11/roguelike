package entities

import (
	"math/rand"
	"roguelike/dungeon"
)

type Enemy struct {
	X, Y                       int
	Symbol                     rune
	Type                       string
	HP, MaxHP, Speed, Cooldown int
	TickCount                  int
	IsRanged, CanSplit         bool
	XPGain                     int
}

func NewEnemy(x, y int) *Enemy {
	types := []string{"Slime", "Bat", "Skeleton"}
	t := types[rand.Intn(len(types))]

	switch t {
	case "Slime":
		return &Enemy{
			X: x, Y: y, Type: "Slime", HP: 4, MaxHP: 4, Symbol: 'S', Speed: 2, CanSplit: true, XPGain: 2,
		}
	case "Bat":
		return &Enemy{
			X: x, Y: y, Type: "Bat", HP: 2, MaxHP: 2, Symbol: 'B', Speed: 1, XPGain: 1,
		}
	case "Skeleton":
		return &Enemy{
			X: x, Y: y, Type: "Skeleton", HP: 3, MaxHP: 3, Symbol: 'K', Speed: 2, Cooldown: 3, IsRanged: true, XPGain: 1,
		}
	}
	return &Enemy{X: x, Y: y, HP: 3, Symbol: 'e'}
}

// Check if enemy is alive
func (e *Enemy) IsAlive() bool {
	return e.HP > 0
}

// Move enemy toward player
func (e *Enemy) MoveToward(targetX, targetY int, occupiedFn func(x, y int) bool) {
	dx, dy := 0, 0

	if targetX < e.X {
		dx = -1
	} else if targetX > e.X {
		dx = 1
	}

	if targetY < e.Y {
		dy = -1
	} else if targetY > e.Y {
		dy = 1
	}

	// Try horizontal move first
	newX := e.X + dx
	newY := e.Y

	if dungeon.IsWalkable(newX, newY) && !occupiedFn(newX, newY) {
		e.X = newX
		return
	}

	// Then try vertical move
	newX = e.X
	newY = e.Y + dy

	if dungeon.IsWalkable(newX, newY) && !occupiedFn(newX, newY) {
		e.Y = newY
	}
}
