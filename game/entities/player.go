package entities

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	Img  *ebiten.Image
	X, Y float64
}

func (s *Sprite) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(s.X, s.Y)
	screen.DrawImage(s.Img.SubImage(image.Rect(0, 0, 16, 16)).(*ebiten.Image), opts)
}
