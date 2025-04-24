package game

import (
	"roguelike/assets"
	"roguelike/game/dungeon"
	"roguelike/game/entities"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Player      *entities.Player
	Enemies     []*entities.Enemy
	Potions     []*entities.Potion
	tilemapJSON *dungeon.TilemapJSON
	tilemapImg  *ebiten.Image
}

func NewGame() (*Game, error) {
	playerImg, err := assets.LoadImage("assets/images/Player/SpriteSheet.png")
	if err != nil {
		return nil, err
	}
	enemyImg, err := assets.LoadImage("assets/images/Cyclope/SpriteSheet.png")
	if err != nil {
		return nil, err
	}
	potionImg, err := assets.LoadImage("assets/images/Items/LifePot.png")
	if err != nil {
		return nil, err
	}
	tilemapImg, err := assets.LoadImage("assets/images/TilesetInteriorFloor.png")
	if err != nil {
		return nil, err
	}
	tilemapJSON, err := dungeon.NewTilemapJSON("assets/maps/spawn.json")
	if err != nil {
		return nil, err
	}
	return &Game{
		Player: &entities.Player{
			Sprite: &entities.Sprite{
				Img: playerImg,
				X:   100,
				Y:   100,
			},
			HP:    10,
			MaxHP: 10,
		},
		Enemies: []*entities.Enemy{
			{Sprite: &entities.Sprite{Img: enemyImg, X: 200, Y: 200}, FollowPlayer: true},
			{Sprite: &entities.Sprite{Img: enemyImg, X: 150, Y: 150}, FollowPlayer: false},
		},
		Potions: []*entities.Potion{
			{Sprite: &entities.Sprite{Img: potionImg, X: 200, Y: 100}, HealAmount: 5},
		},
		tilemapJSON: tilemapJSON,
		tilemapImg:  tilemapImg,
	}, nil
}
