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
var golds []entities.Gold
var equipments []entities.Equipment

var player *entities.Player

var inventory []string
var showInventory bool

var floor int = 1

var playerWasHitLastTurn bool

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
	golds = spawnGold(5)
	equipments = spawnEquipment(2)

	for {
		screen.Clear()

		// Draw map
		drawMap(screen, style)

		// Draw player
		playerStyle := style
		if playerWasHitLastTurn {
			playerStyle = playerStyle.Foreground(tcell.ColorRed)
			playerWasHitLastTurn = false
		}
		screen.SetContent(player.X, player.Y, '@', nil, playerStyle)

		// Draw enemies
		for _, e := range enemies {
			if e.IsAlive() {
				screen.SetContent(e.X, e.Y, e.Symbol, nil, style)
			}
		}

		// Draw potions
		for _, p := range potions {
			screen.SetContent(p.X, p.Y, p.Symbol, nil, style)
		}

		// Draw golds
		for _, g := range golds {
			screen.SetContent(g.X, g.Y, g.Symbol, nil, style)
		}

		goldCount := 0
		for _, item := range inventory {
			if item == "Gold" {
				goldCount++
			}
		}

		goldStr := "Gold: x" + strconv.Itoa(goldCount)
		for i, r := range goldStr {
			screen.SetContent(i+15, dungeon.MapHeight, r, nil, style)
		}

		//Draw equipments
		for _, eq := range equipments {
			screen.SetContent(eq.X, eq.Y, eq.Symbol, nil, style)
		}

		// Draw player HP
		hpStr := "HP: " + strconv.Itoa(player.HP)
		for i, r := range hpStr {
			screen.SetContent(i, dungeon.MapHeight, r, nil, style)
		}

		levelStr := "Lvl: " + strconv.Itoa(player.Level) + "  XP: " + strconv.Itoa(player.XP) + "/" + strconv.Itoa(player.XPToNextLevel)
		for i, r := range levelStr {
			screen.SetContent(i+30, dungeon.MapHeight, r, nil, style)
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

			if player.Equipped["Sword"] {
				drawStr(screen, "Equipped: Sword", 2)
			}
			if player.Equipped["Shield"] {
				drawStr(screen, "Equipped: Shield", 3)
			}

			hpsCount := 0
			for _, item := range inventory {
				if item == "Health Potion" {
					hpsCount++
				}
			}

			if hpsCount > 0 {
				hpsStr := "Health Potion: x" + strconv.Itoa(hpsCount)
				for j, ch := range hpsStr {
					screen.SetContent(j, 1, ch, nil, style)
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
			case 'e':
				handleEquip("Sword")
				handleEquip("Shield")
			}
		}

		// Enemy movement and attack
		for _, e := range enemies {
			if !e.IsAlive() {
				continue
			}

			e.TickCount++
			if e.TickCount%e.Speed != 0 {
				continue // wait based on speed
			}

			// Skeleton ranged attack
			if e.Type == "skeleton" && e.Cooldown <= 0 && abs(e.X-player.X)+abs(e.Y-player.Y) <= 6 {
				playerWasHitLastTurn = true
				player.HP -= 1
				addLog("An arrow hits you! -1 HP")
				e.Cooldown = 3
				continue
			}
			if e.Type == "skeleton" {
				e.Cooldown--
			}

			// Move toward player
			e.MoveToward(player.X, player.Y, func(x, y int) bool {
				return isOccupied(x, y, player, enemies, potions, golds, equipments)
			})

			// Melee attack
			if abs(e.X-player.X)+abs(e.Y-player.Y) == 1 {
				playerWasHitLastTurn = true
				damage := max(1-player.Def, 1)
				player.HP -= damage
				addLog("-" + strconv.Itoa(damage) + " HP from " + e.Type + "!")
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

func checkForGoldCollect() {
	for i := range golds {
		if golds[i].X == player.X && golds[i].Y == player.Y {
			inventory = append(inventory, "Gold")
			addLog("You collect 1 gold.")
			golds = slices.Delete(golds, i, i+1)
			return
		}
	}
}

func checkForEquipmentPickup() {
	for i := range equipments {
		if equipments[i].X == player.X && equipments[i].Y == player.Y {
			item := equipments[i]
			inventory = append(inventory, item.Name)
			addLog("You picked up a " + item.Name + ".")
			equipments = slices.Delete(equipments, i, i+1)
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

func handleEquip(item string) {
	for i, it := range inventory {
		if it == item && !player.Equipped[item] {
			player.Equipped[item] = true
			inventory = slices.Delete(inventory, i, i+1)
			if item == "Sword" {
				player.Atk += 2
				addLog("You equipped a Sword. (+2 ATK)")
			} else if item == "Shield" {
				player.Def += 1
				addLog("You equipped a Shield. (+1 DEF)")
			}
			return
		}
	}
	addLog("No " + item + " in inventory.")
}

func checkLevelUp() {
	for player.XP >= player.XPToNextLevel {
		player.XP -= player.XPToNextLevel
		player.Level++
		player.XPToNextLevel += 5

		player.MaxHP += 3
		player.HP = player.MaxHP
		player.Atk += 1
		player.Def += 1

		addLog("LEVEL UP! You are now level " + strconv.Itoa(player.Level) + "!")
		addLog("+3 MaxHP, +1 ATK, +1 DEF")
	}
}

func tryMovePlayer(dx, dy int) {
	newX := player.X + dx
	newY := player.Y + dy

	// Check for enemy at target location
	for _, e := range enemies {
		if e.X == newX && e.Y == newY && e.IsAlive() {
			e.HP -= player.Atk
			addLog("Enemy -" + strconv.Itoa(player.Atk) + " HP. (" + strconv.Itoa(max(e.HP, 0)) + "/" + strconv.Itoa(e.MaxHP) + ")")
			if e.HP <= 0 {
				addLog("You killed a " + e.Type + "!")
				player.XP += e.XPGain
				addLog("You gained" + strconv.Itoa(e.XPGain) + " XP. (" + strconv.Itoa(player.XP) + "/" + strconv.Itoa(player.XPToNextLevel) + ")")

				// Slime splits!
				if e.Type == "Slime" && e.CanSplit {
					spawnBaby := func(dx, dy int) {
						nx, ny := e.X+dx, e.Y+dy
						splitStats := e.MaxHP / 2
						if dungeon.IsWalkable(nx, ny) && !isOccupied(nx, ny, player, enemies, potions, golds, equipments) {
							enemies = append(enemies, &entities.Enemy{
								X: nx, Y: ny, Type: "Slime", HP: splitStats, MaxHP: splitStats, Symbol: 's', Speed: 1, CanSplit: splitStats > 0, XPGain: splitStats / 2,
							})
							addLog("A Baby Slime emerges!")
						}
					}
					spawnBaby(1, 0)
					spawnBaby(-1, 0)
				}

				e.HP = 0 // mark as dead
				checkLevelUp()
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
	checkForGoldCollect()
	checkForEquipmentPickup()
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

func drawStr(screen tcell.Screen, s string, y int) {
	for i, ch := range s {
		screen.SetContent(i, y, ch, nil, tcell.StyleDefault)
	}
}

func isOccupied(x, y int, player *entities.Player, enemies []*entities.Enemy, potions []entities.Potion, golds []entities.Gold, equipment []entities.Equipment) bool {
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
	for _, g := range golds {
		if g.X == x && g.Y == y {
			return true
		}
	}
	for _, eq := range equipment {
		if eq.X == x && eq.Y == y {
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

func spawnGold(n int) []entities.Gold {
	var list []entities.Gold
	for range n {
		list = append(list, entities.NewGold(dungeon.RandomFloorTile()))
	}
	return list
}

func spawnEquipment(n int) []entities.Equipment {
	var list []entities.Equipment
	items := []string{"Sword", "Shield"}
	for i := range n {
		name := items[i%len(items)]
		x, y := dungeon.RandomFloorTile()
		list = append(list, entities.NewEquipment(name, x, y))
	}
	return list
}
