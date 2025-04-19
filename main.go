package main

import (
	"log"
	"math/rand"
	"time"

	"roguelike/dungeon"
	"roguelike/entities"

	"github.com/gdamore/tcell/v2"
)

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("Error creating screen: %v", err)
	}
	if err = screen.Init(); err != nil {
		log.Fatalf("Error initializing screen: %v", err)
	}
	defer screen.Fini()

	screen.Clear()
	style := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorWhite)

	rand.Seed(time.Now().UnixNano())
	playerX, playerY := dungeon.GenerateDungeon()
	x, y := playerX, playerY

	// Spawn enemies
	var enemies []*entities.Enemy
	for i := 0; i < 5; i++ {
		for {
			ex := rand.Intn(dungeon.MapWidth)
			ey := rand.Intn(dungeon.MapHeight)
			if dungeon.IsWalkable(ex, ey) && (ex != x || ey != y) {
				enemies = append(enemies, entities.NewEnemy(ex, ey))
				break
			}
		}
	}

	for {
		screen.Clear()

		// Draw dungeon
		drawMap(screen, style)

		// Draw player
		screen.SetContent(x, y, '@', nil, style)

		// Draw enemies
		for _, e := range enemies {
			screen.SetContent(e.X, e.Y, e.Symbol, nil, style)
		}

		screen.Show()

		// Handle input
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape:
				return
			case tcell.KeyUp:
				if dungeon.IsWalkable(x, y-1) {
					y--
				}
			case tcell.KeyDown:
				if dungeon.IsWalkable(x, y+1) {
					y++
				}
			case tcell.KeyLeft:
				if dungeon.IsWalkable(x-1, y) {
					x--
				}
			case tcell.KeyRight:
				if dungeon.IsWalkable(x+1, y) {
					x++
				}
			}
		}

		// Move enemies
		for _, e := range enemies {
			e.MoveRandom()
		}
	}
}

func drawMap(screen tcell.Screen, style tcell.Style) {
	for y, row := range dungeon.GameMap {
		for x, ch := range row {
			screen.SetContent(x, y, ch, nil, style)
		}
	}
}
