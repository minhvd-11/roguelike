package entities

type Player struct {
	X, Y  int
	HP    int
	MaxHP int
}

func NewPlayer(x, y int) *Player {
	return &Player{X: x, Y: y, HP: 10, MaxHP: 10}
}
