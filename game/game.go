package game

import (
	"roguelike/assets"
	"roguelike/game/entities"
)

type Game struct {
	Player  *entities.Sprite
	Enemies []*entities.Enemy
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

	return &Game{
		Player: &entities.Sprite{
			Img: playerImg,
			X:   100,
			Y:   100,
		},
		Enemies: []*entities.Enemy{
			{Sprite: &entities.Sprite{Img: enemyImg, X: 200, Y: 200}, FollowPlayer: true},
			{Sprite: &entities.Sprite{Img: enemyImg, X: 150, Y: 150}, FollowPlayer: false},
		},
	}, nil
}
