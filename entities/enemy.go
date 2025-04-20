package entities

import (
	"roguelike/dungeon"
)

type Enemy struct {
	X, Y   int
	Symbol rune
	HP     int
}

func NewEnemy(x, y int) *Enemy {
	return &Enemy{X: x, Y: y, Symbol: 'E', HP: 3}
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
