package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	X     int
	Y     int
	Image *ebiten.Image
}

func NewPlayer(x int, y int) (Player, error) {
	image := mustLoadImage("assets/graphics/player.png")

	player := Player{
		X:     x,
		Y:     y,
		Image: image,
	}
	return player, nil
}

func (player *Player) Draw(screen *ebiten.Image, g *Game) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(player.X*gd.TileSize), float64(player.Y*gd.TileSize+(gd.TileSize/2)))
	screen.DrawImage(player.Image, op)
}

func (player *Player) MoveRight(g *Game) {
	// miramos si se puede mover a la nueva casilla
	tileType := g.CurrentLevel.Tiles[player.Y][player.X+1].TileType
	if tileType != TileFloor && tileType != TileGoal {
		return
	}

	// miramos si hay una caja en la casilla
	for i, box := range g.CurrentLevel.Boxes {
		if player.Y == box.Y && player.X+1 == box.X {
			if box.CanMoveRight(g.CurrentLevel) {
				g.CurrentLevel.AddPushMovement(player.X, player.Y, i, g.CurrentLevel.Boxes[i].X, g.CurrentLevel.Boxes[i].Y)

				g.CurrentLevel.Boxes[i].MoveRight()
				g.CurrentLevel.Pushes++

				player.X++
				g.CurrentLevel.Steps++
			}

			stepAudio.Rewind()
			stepAudio.Play()

			return
		}
	}

	g.CurrentLevel.AddMovement(player.X, player.Y)

	player.X++
	g.CurrentLevel.Steps++

	stepAudio.Rewind()
	stepAudio.Play()
}

func (player *Player) MoveLeft(g *Game) {
	tileType := g.CurrentLevel.Tiles[player.Y][player.X-1].TileType
	if tileType != TileFloor && tileType != TileGoal {
		return
	}

	// miramos si hay una caja en la casilla
	for i, box := range g.CurrentLevel.Boxes {
		if player.Y == box.Y && player.X-1 == box.X {
			if box.CanMoveLeft(g.CurrentLevel) {
				g.CurrentLevel.AddPushMovement(player.X, player.Y, i, g.CurrentLevel.Boxes[i].X, g.CurrentLevel.Boxes[i].Y)

				g.CurrentLevel.Boxes[i].MoveLeft()
				g.CurrentLevel.Pushes++

				player.X--
				g.CurrentLevel.Steps++
			}

			stepAudio.Rewind()
			stepAudio.Play()

			return
		}
	}

	g.CurrentLevel.AddMovement(player.X, player.Y)

	player.X--
	g.CurrentLevel.Steps++

	stepAudio.Rewind()
	stepAudio.Play()
}

func (player *Player) MoveUp(g *Game) {
	tileType := g.CurrentLevel.Tiles[player.Y-1][player.X].TileType
	if tileType != TileFloor && tileType != TileGoal {
		return
	}

	// miramos si hay una caja en la casilla
	for i, box := range g.CurrentLevel.Boxes {
		if player.Y-1 == box.Y && player.X == box.X {
			if box.CanMoveUp(g.CurrentLevel) {
				g.CurrentLevel.AddPushMovement(player.X, player.Y, i, g.CurrentLevel.Boxes[i].X, g.CurrentLevel.Boxes[i].Y)

				g.CurrentLevel.Boxes[i].MoveUp()
				g.CurrentLevel.Pushes++

				player.Y--
				g.CurrentLevel.Steps++
			}

			stepAudio.Rewind()
			stepAudio.Play()

			return
		}
	}

	g.CurrentLevel.AddMovement(player.X, player.Y)

	player.Y--
	g.CurrentLevel.Steps++

	stepAudio.Rewind()
	stepAudio.Play()
}

func (player *Player) MoveDown(g *Game) {
	tileType := g.CurrentLevel.Tiles[player.Y+1][player.X].TileType
	if tileType != TileFloor && tileType != TileGoal {
		return
	}

	// miramos si hay una caja en la casilla
	for i, box := range g.CurrentLevel.Boxes {
		if player.Y+1 == box.Y && player.X == box.X {
			if box.CanMoveDown(g.CurrentLevel) {
				g.CurrentLevel.AddPushMovement(player.X, player.Y, i, g.CurrentLevel.Boxes[i].X, g.CurrentLevel.Boxes[i].Y)

				g.CurrentLevel.Boxes[i].MoveDown()
				g.CurrentLevel.Pushes++

				player.Y++
				g.CurrentLevel.Steps++
			}

			stepAudio.Rewind()
			stepAudio.Play()

			return
		}
	}

	g.CurrentLevel.AddMovement(player.X, player.Y)

	player.Y++
	g.CurrentLevel.Steps++

	stepAudio.Rewind()
	stepAudio.Play()
}
