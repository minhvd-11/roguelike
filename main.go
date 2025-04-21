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

var logMessages []string

var enemies []*entities.Enemy

var potions []*entities.Potion

var inventory []string
var showInventory bool

var floor int = 1

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
	for range 3 {
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

		// Draw potions
		for _, p := range potions {
			screen.SetContent(p.X, p.Y, '!', nil, style)
		}

		// Draw player HP
		hpStr := "HP: " + strconv.Itoa(player.HP)
		for i, r := range hpStr {
			screen.SetContent(i, dungeon.MapHeight, r, nil, style)
		}

		logY := dungeon.MapHeight + 1
		for i, msg := range logMessages {
			for j, ch := range msg {
				screen.SetContent(j, logY+i, ch, nil, style)
			}
		}

		// Inventory UI
		if showInventory {
			invTitle := "Inventory:"
			for i, ch := range invTitle {
				screen.SetContent(i, 0, ch, nil, style)
			}

			for i, item := range inventory {
				for j, ch := range item {
					screen.SetContent(j, i+1, ch, nil, style)
				}
			}
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
			case 'i':
				showInventory = !showInventory
				addLog("Open Inventory")
			}

			switch ev.Rune() {
			case 'i':
				showInventory = !showInventory
				addLog("Toggled inventory")
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
					addLog("-1 HP from enemy attack!")
				}
			}
		}

		dungeon.UpdateVisibility(player.X, player.Y)
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
			e.HP -= 3
			addLog("Enemy -1 HP!")
			if e.HP == 0 {
				addLog("You killed an enemy!")
			}
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

			inventory = append(inventory, "Health Potion")
			addLog("You picked up and drank a health potion.")
			break
		}
	}

	if dungeon.IsWalkable(newX, newY) {
		player.X = newX
		player.Y = newY
	}

	tile := dungeon.GameMap[player.Y][player.X]
	if tile == '>' {
		floor++
		addLog("You descend to floor " + strconv.Itoa(floor) + "...")
		dungeon.GenerateDungeon()
		player.X = 2
		player.Y = 2
		dungeon.UpdateVisibility(player.X, player.Y)

		// Respawn enemies & potions
		enemies = spawnEnemies(5)
		potions = spawnPotions(3)
	}

}

func drawMap(screen tcell.Screen, style tcell.Style) {
	for y, row := range dungeon.GameMap {
		for x, ch := range row {
			if dungeon.Visible[y][x] {
				screen.SetContent(x, y, ch, nil, style)
			} else {
				screen.SetContent(x, y, ' ', nil, style)
			}

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

func addLog(msg string) {
	logMessages = append(logMessages, msg)
	if len(logMessages) > 3 {
		logMessages = logMessages[1:] // keep last 3 messages
	}
}

func spawnEnemies(n int) []*entities.Enemy {
	var list []*entities.Enemy
	for i := 0; i < n; i++ {
		e := entities.NewEnemy(dungeon.RandomFloorTile())
		list = append(list, e)
	}
	return list
}

func spawnPotions(n int) []*entities.Potion {
	var list []*entities.Potion
	for i := 0; i < n; i++ {
		list = append(list, entities.NewPotion(dungeon.RandomFloorTile()))
	}
	return list
}
