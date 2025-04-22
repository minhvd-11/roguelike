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

var potions []entities.Potion

var player *entities.Player

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
	player = entities.NewPlayer(dungeon.GenerateDungeon())

	enemies = spawnEnemies(5)
	potions = spawnPotions(3)

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
				tryMovePlayer(0, -1)
			case tcell.KeyDown:
				tryMovePlayer(0, 1)
			case tcell.KeyLeft:
				tryMovePlayer(-1, 0)
			case tcell.KeyRight:
				tryMovePlayer(1, 0)
			}
			switch ev.Rune() {
			case 'i':
				showInventory = !showInventory
				addLog("Toggled inventory")
			case 'h':
				usePotion()
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

func checkForPotionPickup() {
	for i := range potions {
		if potions[i].X == player.X && potions[i].Y == player.Y {
			inventory = append(inventory, "Health Potion")
			addLog("You picked up a health potion.")
			potions = slices.Delete(potions, i, i+1)
			return
		}
	}
}

func checkForStairs() {
	tile := dungeon.GameMap[player.Y][player.X]
	if tile == '>' {
		floor++
		addLog("You descend to floor " + strconv.Itoa(floor))
		player.X, player.Y = dungeon.GenerateDungeon()
		enemies = spawnEnemies(5)
		potions = spawnPotions(3)
		dungeon.UpdateVisibility(player.X, player.Y)
	}
}

func usePotion() {
	for i, item := range inventory {
		if item == "Health Potion" {
			if player.HP < player.MaxHP {
				player.HP += 5
				if player.HP > player.MaxHP {
					player.HP = player.MaxHP
				}
				inventory = slices.Delete(inventory, i, i+1)
				addLog("You drink a health potion.")
			} else {
				addLog("You're already at full health.")
			}
			return
		}
	}
	addLog("You have no potions.")
}

func tryMovePlayer(dx, dy int) {
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

	if dungeon.IsWalkable(newX, newY) {
		player.X = newX
		player.Y = newY
	}

	checkForPotionPickup()
	checkForStairs()

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

func isOccupied(x, y int, player *entities.Player, enemies []*entities.Enemy, potions []entities.Potion) bool {
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
	for range n {
		e := entities.NewEnemy(dungeon.RandomFloorTile())
		list = append(list, e)
	}
	return list
}

func spawnPotions(n int) []entities.Potion {
	var list []entities.Potion
	for range n {
		list = append(list, entities.NewPotion(dungeon.RandomFloorTile()))
	}
	return list
}
