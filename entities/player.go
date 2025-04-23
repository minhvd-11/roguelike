package entities

type Player struct {
	X, Y      int
	HP, MaxHP int
	Atk       int
	Def       int
	Equipped  map[string]bool
}

func NewPlayer(startX, startY int) *Player {
	return &Player{
		X:        startX,
		Y:        startY,
		HP:       10,
		MaxHP:    10,
		Atk:      2,
		Def:      0,
		Equipped: make(map[string]bool),
	}
}
