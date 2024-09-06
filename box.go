package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Box struct {
	X     int
	Y     int
	Image *ebiten.Image
}

func NewBox(x int, y int) (Box, error) {
	image := mustLoadImage("assets/graphics/box.png")

	box := Box{
		X:     x,
		Y:     y,
		Image: image,
	}
	return box, nil
}

func (box *Box) Draw(screen *ebiten.Image, g *Game) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(box.X*gd.TileSize), float64(box.Y*gd.TileSize+(gd.TileSize/2)))
	screen.DrawImage(box.Image, op)
}

func (box *Box) CanMoveRight(level *Level) bool {
	tileType := level.Tiles[box.Y][box.X+1].TileType
	if tileType != TileFloor && tileType != TileGoal {
		return false
	}

	// miramos si hay una caja en la casilla
	for _, otherBox := range level.Boxes {
		if box.Y == otherBox.Y && box.X+1 == otherBox.X {
			return false
		}
	}

	return true
}

func (box *Box) MoveRight() {
	box.X++
}

func (box *Box) CanMoveLeft(level *Level) bool {
	tileType := level.Tiles[box.Y][box.X-1].TileType
	if tileType != TileFloor && tileType != TileGoal {
		return false
	}

	// miramos si hay una caja en la casilla
	for _, otherBox := range level.Boxes {
		if box.Y == otherBox.Y && box.X-1 == otherBox.X {
			return false
		}
	}

	return true
}

func (box *Box) MoveLeft() {
	box.X--
}

func (box *Box) CanMoveUp(level *Level) bool {
	tileType := level.Tiles[box.Y-1][box.X].TileType
	if tileType != TileFloor && tileType != TileGoal {
		return false
	}

	// miramos si hay una caja en la casilla
	for _, otherBox := range level.Boxes {
		if box.Y-1 == otherBox.Y && box.X == otherBox.X {
			return false
		}
	}

	return true
}

func (box *Box) MoveUp() {
	box.Y--
}

func (box *Box) CanMoveDown(level *Level) bool {
	tileType := level.Tiles[box.Y+1][box.X].TileType
	if tileType != TileFloor && tileType != TileGoal {
		return false
	}

	// miramos si hay una caja en la casilla
	for _, otherBox := range level.Boxes {
		if box.Y+1 == otherBox.Y && box.X == otherBox.X {
			return false
		}
	}

	return true
}

func (box *Box) MoveDown() {
	box.Y++
}
