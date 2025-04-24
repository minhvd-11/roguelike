package game

import (
	"image"
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

	for _, potion := range g.Potions {
		if g.Player.X == potion.X && g.Player.Y == potion.Y {
			g.Player.HP = min(g.Player.HP+potion.HealAmount, g.Player.MaxHP)
			potion.X, potion.Y = -100, -100 // Move potion off-screen after use
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{120, 180, 255, 255})

	opts := ebiten.DrawImageOptions{}

	// loop over the layers
	for _, layer := range g.tilemapJSON.Layers {
		// loop over the tiles in the layer data
		for index, id := range layer.Data {

			// get the tile position of the tile
			x := index % layer.Width
			y := index / layer.Width

			// convert the tile position to pixel position
			x *= 16
			y *= 16

			// get the position on the image where the tile id is
			srcX := (id - 1) % 22
			srcY := (id - 1) / 22

			// convert the src tile pos to pixel src position
			srcX *= 16
			srcY *= 16

			// set the drawimageoptions to draw the tile at x, y
			opts.GeoM.Translate(float64(x), float64(y))

			// draw the tile
			screen.DrawImage(
				// cropping out the tile that we want from the spritesheet
				g.tilemapImg.SubImage(image.Rect(srcX, srcY, srcX+16, srcY+16)).(*ebiten.Image),
				&opts,
			)

			// reset the opts for the next tile
			opts.GeoM.Reset()
		}
	}

	g.Player.Draw(screen)

	for _, enemy := range g.Enemies {
		enemy.Draw(screen)
	}

	for _, potion := range g.Potions {
		potion.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 320, 540
}
