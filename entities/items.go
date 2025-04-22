package entities

type Potion struct {
	X, Y   int
	Symbol rune
}

func NewPotion(x, y int) Potion {
	return Potion{X: x, Y: y, Symbol: '!'}
}

type Gold struct {
	X, Y   int
	Symbol rune
}

func NewGold(x, y int) Gold {
	return Gold{X: x, Y: y, Symbol: '$'}
}
