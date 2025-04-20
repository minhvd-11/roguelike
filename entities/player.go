package entities

type Player struct {
	X  int
	Y  int
	HP int
}

func NewPlayer(x, y int) *Player {
	return &Player{X: x, Y: y, HP: 10}
}
