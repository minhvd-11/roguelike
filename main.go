package main

import (
	"log"
	"math/rand"
	"strconv"
	"time"

	"roguelike/dungeon"
	"roguelike/entities"

	"slices"

	"github.com/gdamore/tcell/v2"
)

var enemies []*entities.Enemy

var potions []*entities.Potion

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
	player := entities.NewPlayer(playerX, playerY)

	// Spawn enemies
	for range 5 {
		for {
			ex := rand.Intn(dungeon.MapWidth)
			ey := rand.Intn(dungeon.MapHeight)
			if dungeon.IsWalkable(ex, ey) && (ex != player.X || ey != player.Y) {
				enemies = append(enemies, entities.NewEnemy(ex, ey))
				break
			}
		}
	}

	//spawn potions
	for range 3 {
		for {
			x := rand.Intn(dungeon.MapWidth)
			y := rand.Intn(dungeon.MapHeight)
			if dungeon.IsWalkable(x, y) && (x != player.X || y != player.Y) {
				potions = append(potions, &entities.Potion{X: x, Y: y})
				break
			}
		}
	}

	for {
		screen.Clear()

		// Draw map
		drawMap(screen, style)

		// Draw player
		screen.SetContent(player.X, player.Y, '@', nil, style)

		// Draw enemies
		for _, e := range enemies {
			if e.IsAlive() {
				screen.SetContent(e.X, e.Y, e.Symbol, nil, style)
			}
		}

		for _, p := range potions {
			screen.SetContent(p.X, p.Y, '!', nil, style)
		}

		// Draw player HP
		hpStr := "HP: " + strconv.Itoa(player.HP)
		for i, r := range hpStr {
			screen.SetContent(i, dungeon.MapHeight, r, nil, style)
		}

		screen.Show()

		// Exit if player is dead
		if player.HP <= 0 {
			break
		}

		// Handle input
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape:
				return
			case tcell.KeyUp:
				tryMovePlayer(player, 0, -1)
			case tcell.KeyDown:
				tryMovePlayer(player, 0, 1)
			case tcell.KeyLeft:
				tryMovePlayer(player, -1, 0)
			case tcell.KeyRight:
				tryMovePlayer(player, 1, 0)
			}
		}

		// Enemy movement and attack
		for _, e := range enemies {
			if e.IsAlive() {
				e.MoveToward(player.X, player.Y, func(x, y int) bool {
					return isOccupied(x, y, player, enemies, potions)
				})
				if abs(e.X-player.X)+abs(e.Y-player.Y) == 1 {
					player.HP -= 1
				}

			}
		}
	}

	// Show death message
	screen.Clear()
	msg := "You Died! Press ESC to exit..."
	for i, r := range msg {
		screen.SetContent(i+10, dungeon.MapHeight/2, r, nil, style)
	}
	screen.Show()

	for {
		ev := screen.PollEvent()
		if key, ok := ev.(*tcell.EventKey); ok && key.Key() == tcell.KeyEscape {
			return
		}
	}
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func tryMovePlayer(player *entities.Player, dx, dy int) {
	newX := player.X + dx
	newY := player.Y + dy

	// Check for enemy at target location
	for _, e := range enemies {
		if e.X == newX && e.Y == newY && e.IsAlive() {
			e.HP -= 1
			return // attack instead of moving
		}
	}

	for i, p := range potions {
		if p.X == newX && p.Y == newY {
			if player.HP < 10 {
				player.HP += 3
				if player.HP > 10 {
					player.HP = 10
				}
			}

			potions = slices.Delete(potions, i, i+1)
			break
		}
	}

	if dungeon.IsWalkable(newX, newY) {
		player.X = newX
		player.Y = newY
	}
}

func drawMap(screen tcell.Screen, style tcell.Style) {
	for y, row := range dungeon.GameMap {
		for x, ch := range row {
			screen.SetContent(x, y, ch, nil, style)
		}
	}
}

func isOccupied(x, y int, player *entities.Player, enemies []*entities.Enemy, potions []*entities.Potion) bool {
	if x == player.X && y == player.Y {
		return true
	}
	for _, e := range enemies {
		if e.IsAlive() && e.X == x && e.Y == y {
			return true
		}
	}
	for _, p := range potions {
		if p.X == x && p.Y == y {
			return true
		}
	}
	return false
}
