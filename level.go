package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image/color"
)

const (
	TileFloor  string = "floor"
	TileWall   string = "wall"
	TileBox    string = "box"
	TileGoal   string = "goal"
	TilePlayer string = "player"
	TileEmpty  string = "empty"
)

type Tile struct {
	X        int
	Y        int
	TileType string
	Image    *ebiten.Image
}

type Movement struct {
	PlayerLastX int
	PlayerLastY int
	HasPush     bool
	BoxIndex    int
	BoxLastX    int
	BoxLastY    int
}

type Level struct {
	Tiles         [][]Tile
	Boxes         []Box
	Player        Player
	Steps         int
	Pushes        int
	IsCompleted   bool
	LastMovements []Movement
}

func NewLevel(numLevel int) Level {
	l := Level{}
	l.createTiles(numLevel)

	return l
}

func (level *Level) Draw(screen *ebiten.Image, g *Game) {
	for x := 0; x < gd.TilesX; x++ {
		for y := 0; y < gd.TilesY; y++ {
			tile := level.Tiles[y][x]
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(tile.X*gd.TileSize), float64(tile.Y*gd.TileSize+(gd.TileSize/2)))
			screen.DrawImage(tile.Image, op)
		}
	}

	op := &text.DrawOptions{}
	op.GeoM.Translate(20, 10)
	op.ColorScale.ScaleWithColor(color.White)
	score := fmt.Sprintf("Level: %d/%d", g.CurrentLevelNum+1, len(g.Levels))
	text.Draw(screen, score, &text.GoTextFace{Source: mplusFaceSource, Size: 16}, op)

	op = &text.DrawOptions{}
	op.GeoM.Translate(850, 10)
	steps := fmt.Sprintf("Steps: %d", level.Steps)
	text.Draw(screen, steps, &text.GoTextFace{Source: mplusFaceSource, Size: 16}, op)

	op = &text.DrawOptions{}
	op.GeoM.Translate(1050, 10)
	pushes := fmt.Sprintf("Pushes: %d", level.Pushes)
	text.Draw(screen, pushes, &text.GoTextFace{Source: mplusFaceSource, Size: 16}, op)
}

func (level *Level) IsLevelCompleted() bool {
	for _, box := range level.Boxes {
		if level.Tiles[box.Y][box.X].TileType != TileGoal {
			return false
		}
	}

	return true
}

func newTile(x int, y int, tileType string) (Tile, error) {
	image := mustLoadImage("assets/graphics/" + tileType + ".png")

	tile := Tile{
		X:        x,
		Y:        y,
		TileType: tileType,
		Image:    image,
	}
	return tile, nil
}

func (level *Level) createTiles(numLevel int) {

	tiles := make([][]Tile, gd.TilesY)
	for i := range tiles {
		tiles[i] = make([]Tile, gd.TilesX)
	}
	boxes := make([]Box, 0)
	player, _ := NewPlayer(0, 0)

	for y := 0; y < gd.TilesY; y++ {
		for x := 0; x < gd.TilesX; x++ {
			switch levelsDefinition[numLevel][y][x] {
			case '#':
				wall, err := newTile(x, y, TileWall)
				if err != nil {
					panic(err)
				}
				tiles[y][x] = wall
			case '-':
				floor, err := newTile(x, y, TileFloor)
				if err != nil {
					panic(err)
				}
				tiles[y][x] = floor
			case '$':
				floor, err := newTile(x, y, TileFloor)
				if err != nil {
					panic(err)
				}
				tiles[y][x] = floor

				box, err := NewBox(x, y)
				if err != nil {
					panic(err)
				}
				boxes = append(boxes, box)
			case '*':
				goal, err := newTile(x, y, TileGoal)
				if err != nil {
					panic(err)
				}
				tiles[y][x] = goal

				box, err := NewBox(x, y)
				if err != nil {
					panic(err)
				}
				boxes = append(boxes, box)
			case '.':
				goal, err := newTile(x, y, TileGoal)
				if err != nil {
					panic(err)
				}
				tiles[y][x] = goal
			case '@':
				floor, err := newTile(x, y, TileFloor)
				if err != nil {
					panic(err)
				}
				tiles[y][x] = floor

				player, err = NewPlayer(x, y)
				if err != nil {
					panic(err)
				}
			case '+':
				goal, err := newTile(x, y, TileGoal)
				if err != nil {
					panic(err)
				}
				tiles[y][x] = goal

				player, err = NewPlayer(x, y)
				if err != nil {
					panic(err)
				}
			default:
				empty, err := newTile(x, y, TileEmpty)
				if err != nil {
					panic(err)
				}
				tiles[y][x] = empty

			}
		}
	}

	level.Tiles = tiles
	level.Boxes = boxes
	level.Player = player
}

func (level *Level) AddMovement(x int, y int) {
	m := Movement{PlayerLastX: x, PlayerLastY: y}

	level.LastMovements = append(level.LastMovements, m)

	// only 3 last movements can be removed
	if len(level.LastMovements) > 3 {
		level.LastMovements = level.LastMovements[1:4]
	}

}

func (level *Level) AddPushMovement(playerX int, playerY int, boxIndex int, boxX int, boxY int) {
	m := Movement{PlayerLastX: playerX, PlayerLastY: playerY, HasPush: true, BoxIndex: boxIndex, BoxLastX: boxX, BoxLastY: boxY}

	level.LastMovements = append(level.LastMovements, m)

	// only 3 last movements can be removed
	if len(level.LastMovements) > 3 {
		level.LastMovements = level.LastMovements[1:4]
	}
}

func (level *Level) RemoveMovement() {
	if len(level.LastMovements) == 0 {
		return
	}

	m := level.LastMovements[len(level.LastMovements)-1]
	level.LastMovements = level.LastMovements[0 : len(level.LastMovements)-1]

	level.Player.X = m.PlayerLastX
	level.Player.Y = m.PlayerLastY

	level.Steps--

	if m.HasPush {
		level.Boxes[m.BoxIndex].X = m.BoxLastX
		level.Boxes[m.BoxIndex].Y = m.BoxLastY
		level.Pushes--
	}
}
