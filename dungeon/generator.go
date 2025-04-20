package dungeon

import (
	"math/rand"
	"slices"
	"time"
)

const (
	MapWidth    = 80
	MapHeight   = 24
	MaxRooms    = 8
	RoomMinSize = 5
	RoomMaxSize = 10
)

var GameMap [][]rune

type Room struct {
	X1, Y1, X2, Y2 int
}

func (r Room) Center() (int, int) {
	return (r.X1 + r.X2) / 2, (r.Y1 + r.Y2) / 2
}

func (r Room) Intersects(other Room) bool {
	return r.X1 <= other.X2 && r.X2 >= other.X1 &&
		r.Y1 <= other.Y2 && r.Y2 >= other.Y1
}

func GenerateDungeon() (int, int) {
	rand.Seed(time.Now().UnixNano())
	rooms := []Room{}
	var playerX, playerY int

	GameMap = make([][]rune, MapHeight)
	for y := range GameMap {
		GameMap[y] = make([]rune, MapWidth)
		for x := range GameMap[y] {
			GameMap[y][x] = '#'
		}
	}

	for range MaxRooms {
		w := RoomMinSize + rand.Intn(RoomMaxSize-RoomMinSize+1)
		h := RoomMinSize + rand.Intn(RoomMaxSize-RoomMinSize+1)
		x := rand.Intn(MapWidth - w - 1)
		y := rand.Intn(MapHeight - h - 1)

		newRoom := Room{x, y, x + w, y + h}
		intersects := slices.ContainsFunc(rooms, newRoom.Intersects)

		if !intersects {
			createRoom(newRoom)
			if len(rooms) > 0 {
				prev := rooms[len(rooms)-1]
				prevX, prevY := prev.Center()
				newX, newY := newRoom.Center()

				if rand.Intn(2) == 1 {
					createHTunnel(prevX, newX, prevY)
					createVTunnel(prevY, newY, newX)
				} else {
					createVTunnel(prevY, newY, prevX)
					createHTunnel(prevX, newX, newY)
				}
			} else {
				playerX, playerY = newRoom.Center()
			}

			rooms = append(rooms, newRoom)
		}
	}

	return playerX, playerY
}

func createRoom(r Room) {
	for y := r.Y1; y <= r.Y2; y++ {
		for x := r.X1; x <= r.X2; x++ {
			GameMap[y][x] = '.'
		}
	}
}

func createHTunnel(x1, x2, y int) {
	for x := min(x1, x2); x <= max(x1, x2); x++ {
		GameMap[y][x] = '.'
	}
}

func createVTunnel(y1, y2, x int) {
	for y := min(y1, y2); y <= max(y1, y2); y++ {
		GameMap[y][x] = '.'
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func IsWalkable(x, y int) bool {
	if y < 0 || y >= len(GameMap) || x < 0 || x >= len(GameMap[y]) {
		return false
	}
	return GameMap[y][x] != '#'
}
