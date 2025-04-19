package main

import (
	"log"

	"github.com/gdamore/tcell/v2"
)

var gameMap = [][]rune{
	[]rune("####################"),
	[]rune("#..................#"),
	[]rune("#..................#"),
	[]rune("#.........######...#"),
	[]rune("#..................#"),
	[]rune("####################"),
}

func isWalkable(x, y int) bool {
	if y < 0 || y >= len(gameMap) || x < 0 || x >= len(gameMap[y]) {
		return false
	}
	return gameMap[y][x] != '#'
}

func drawMap(screen tcell.Screen, style tcell.Style) {
	for y, row := range gameMap {
		for x, ch := range row {
			screen.SetContent(x, y, ch, nil, style)
		}
	}
}

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("Error creating screen: %v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("Error initializing screen: %v", err)
	}
	defer screen.Fini()

	style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)

	// Start player in an open spot
	x, y := 1, 1

	for {
		screen.Clear()
		drawMap(screen, style)
		screen.SetContent(x, y, '@', nil, style)
		screen.Show()

		ev := screen.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventKey:
			newX, newY := x, y
			switch ev.Key() {
			case tcell.KeyEscape, tcell.KeyCtrlC:
				return
			case tcell.KeyUp:
				newY--
			case tcell.KeyDown:
				newY++
			case tcell.KeyLeft:
				newX--
			case tcell.KeyRight:
				newX++
			default:
				switch ev.Rune() {
				case 'w':
					newY--
				case 's':
					newY++
				case 'a':
					newX--
				case 'd':
					newX++
				}
			}

			if isWalkable(newX, newY) {
				x, y = newX, newY
			}

		case *tcell.EventResize:
			screen.Sync()
		}
	}
}
