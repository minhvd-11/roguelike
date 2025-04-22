package entities

type Gold struct {
	X, Y   int
	Symbol rune
}

func NewGold(x, y int) Gold {
	return Gold{X: x, Y: y, Symbol: '$'}
}
