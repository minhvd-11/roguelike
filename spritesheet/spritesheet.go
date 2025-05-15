package spritesheet

import "image"

type SpriteSheet struct {
	WidthInlines  int
	HeightInlines int
	Tilesize      int
}

func (s *SpriteSheet) Rect(index int) image.Rectangle {
	x := (index % s.WidthInlines) * s.Tilesize
	y := (index / s.WidthInlines) * s.Tilesize

	return image.Rect(x, y, x+s.Tilesize, y+s.Tilesize)
}

func NewSpriteSheet(w, h, t int) *SpriteSheet {
	return &SpriteSheet{
		w, h, t,
	}
}
