package entities

type Player struct {
	X, Y                     int
	HP, MaxHP                int
	Atk, Def                 int
	Level, XP, XPToNextLevel int
	Equipped                 map[string]bool
}

func NewPlayer(startX, startY int) *Player {
	return &Player{
		X:             startX,
		Y:             startY,
		HP:            10,
		MaxHP:         10,
		Atk:           2,
		Def:           0,
		Level:         1,
		XP:            0,
		XPToNextLevel: 10,
		Equipped:      make(map[string]bool),
	}
}
