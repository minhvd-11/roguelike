package entities

type Potion struct {
	*Sprite
	HealAmt uint
}

type Gold struct {
	X, Y   int
	Symbol rune
}

func NewGold(x, y int) Gold {
	return Gold{X: x, Y: y, Symbol: '$'}
}

type Equipment struct {
	X, Y   int
	Name   string
	Symbol rune
}

func NewEquipment(name string, x, y int) Equipment {
	symbol := '?'
	if name == "Sword" {
		symbol = '/'
	} else if name == "Shield" {
		symbol = ']'
	}
	return Equipment{X: x, Y: y, Name: name, Symbol: symbol}
}
