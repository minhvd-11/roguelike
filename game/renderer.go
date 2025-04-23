package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.Player.X += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.Player.X -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.Player.Y -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.Player.Y += 2
	}

	for _, enemy := range g.Enemies {
		enemy.Update(g.Player)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{120, 180, 255, 255})

	g.Player.Draw(screen)

	for _, enemy := range g.Enemies {
		enemy.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 320, 540
}
