package game

import (
	"image"
	"roguelike/animation"
	"roguelike/assets"
	"roguelike/game/dungeon"
	"roguelike/game/entities"
	"roguelike/spritesheet"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	player            *entities.Player
	playerSpriteSheet *spritesheet.SpriteSheet
	enemies           []*entities.Enemy
	potions           []*entities.Potion
	tilemapJSON       *dungeon.TilemapJSON
	tilemapImg        *ebiten.Image
	cam               *Camera
	colliders         []image.Rectangle
}

func CheckCollisionHorizontal(sprite *entities.Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(image.Rect(
			int(sprite.X),
			int(sprite.Y),
			int(sprite.X)+16,
			int(sprite.Y)+16,
		),
		) {
			if sprite.Dx > 0 {
				sprite.X = float64(collider.Min.X) - 16
			} else if sprite.Dx < 0 {
				sprite.X = float64(collider.Max.X)
			}
		}
	}
}

func CheckCollisionVertical(sprite *entities.Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(image.Rect(
			int(sprite.X),
			int(sprite.Y),
			int(sprite.X)+16,
			int(sprite.Y)+16,
		),
		) {
			if sprite.Dy > 0 {
				sprite.Y = float64(collider.Min.Y) - 16
			} else if sprite.Dy < 0 {
				sprite.Y = float64(collider.Max.Y)
			}
		}
	}
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

	playerSpriteSHeet := spritesheet.NewSpriteSheet(4, 7, 16)

	return &Game{
		player: &entities.Player{
			Sprite: &entities.Sprite{
				Img: playerImg,
				X:   100,
				Y:   100,
			},
			HP:    10,
			MaxHP: 10,
			Animation: map[entities.PlayerState]*animation.Animation{
				entities.Up:    animation.NewAnimation(5, 13, 4, 20),
				entities.Down:  animation.NewAnimation(4, 12, 4, 20),
				entities.Left:  animation.NewAnimation(6, 14, 4, 20),
				entities.Right: animation.NewAnimation(7, 15, 4, 20)},
		},
		playerSpriteSheet: playerSpriteSHeet,
		enemies: []*entities.Enemy{
			{Sprite: &entities.Sprite{Img: enemyImg, X: 200, Y: 200}, FollowPlayer: true},
			{Sprite: &entities.Sprite{Img: enemyImg, X: 150, Y: 150}, FollowPlayer: false},
		},
		potions: []*entities.Potion{
			{Sprite: &entities.Sprite{Img: potionImg, X: 200, Y: 100}, HealAmt: 5},
		},
		tilemapJSON: tilemapJSON,
		tilemapImg:  tilemapImg,
		cam:         NewCamera(0, 0),
		colliders: []image.Rectangle{
			image.Rect(144, 144, 160, 160),
		},
	}, nil
}
