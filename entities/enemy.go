package entities

import (
	"math/rand"
	"roguelike/dungeon"
)

type Enemy struct {
	X      int
	Y      int
	Symbol rune
}

func NewEnemy(x, y int) *Enemy {
	return &Enemy{X: x, Y: y, Symbol: 'E'}
}

// Move enemy in a random direction (no pathfinding... yet)
func (e *Enemy) MoveRandom() {
	dx := []int{-1, 1, 0, 0}
	dy := []int{0, 0, -1, 1}
	i := rand.Intn(4)
	newX := e.X + dx[i]
	newY := e.Y + dy[i]

	if dungeon.IsWalkable(newX, newY) {
		e.X = newX
		e.Y = newY
	}
}
