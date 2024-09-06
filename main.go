package main

import (
	"bytes"
	"embed"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image"
	"image/color"
	"os"
	"strconv"
)

type Scene int64
type SelectedMode int64

const (
	CoverScene Scene = iota
	PlayingScene
	ExcelScene
	EndScene
	QuitScene
)

const (
	EasyMode SelectedMode = iota
	OriginalMode
	QuitMode
)

type Game struct {
	Levels          []Level
	CurrentLevel    *Level
	CurrentLevelNum int
	CurrentScene    Scene
	Mode            string
	ShowHelp        bool
}

type GameData struct {
	TilesX   int
	TilesY   int
	TileSize int
}

var gd = GameData{
	TilesX:   20,
	TilesY:   17,
	TileSize: 64,
}

var mplusFaceSource *text.GoTextFaceSource
var coverImage *ebiten.Image
var excelImage *ebiten.Image
var endImage *ebiten.Image
var audioContext *audio.Context
var stepAudio *audio.Player
var loopAudio *audio.Player
var levelsDefinition [][]string
var coverSelectedMode SelectedMode

//go:embed all:assets
var assets embed.FS

func NewGame() *Game {
	g := &Game{}

	g.CurrentScene = CoverScene

	return g
}

func (g *Game) Start() {
	switch coverSelectedMode {
	case EasyMode:
		levelsDefinition = easyLevelsDefinition
	case OriginalMode:
		levelsDefinition = originalLevelsDefinition
	}

	g.Levels = g.Levels[:0]
	for i := range levelsDefinition {
		g.Levels = append(g.Levels, NewLevel(i))
	}

	switch coverSelectedMode {
	case EasyMode:
		// if file doesn't exist, levelNum will be empty
		// and CurrentLevelNum will be 0
		levelNum, _ := os.ReadFile("current_level_easy.dat")
		g.CurrentLevelNum, _ = strconv.Atoi(string(levelNum))
	case OriginalMode:
		// if file doesn't exist, levelNum will be empty
		// and CurrentLevelNum will be 0
		levelNum, _ := os.ReadFile("current_level_original.dat")
		g.CurrentLevelNum, _ = strconv.Atoi(string(levelNum))
	}

	loopAudio.Close()

	g.CurrentLevel = &g.Levels[g.CurrentLevelNum]
	g.CurrentScene = PlayingScene
}

func (g *Game) RestartMode() {
	g.CurrentScene = CoverScene
}

func (g *Game) RestartLevel() {
	levels := g.Levels[:g.CurrentLevelNum]
	levels = append(levels, NewLevel(g.CurrentLevelNum))
	levels = append(levels, g.Levels[g.CurrentLevelNum+1:]...)
	g.Levels = levels
	g.ShowHelp = false
}

func (g *Game) PreviousLevel() {
	if g.CurrentLevelNum > 0 {
		g.CurrentLevelNum--
		g.CurrentLevel = &g.Levels[g.CurrentLevelNum]
		g.RestartLevel()
	}
}

func (g *Game) NextLevel() {
	if g.CurrentLevelNum < len(g.Levels)-1 {
		g.ShowHelp = false
		g.CurrentLevelNum++
		g.CurrentLevel = &g.Levels[g.CurrentLevelNum]
		g.RestartLevel()

		// save current level so we can load it later
		data := []byte(strconv.Itoa(g.CurrentLevelNum))

		switch coverSelectedMode {
		case EasyMode:
			err := os.WriteFile("current_level_easy.dat", data, 0777)
			if err != nil {
				panic(err)
			}
		case OriginalMode:
			err := os.WriteFile("current_level_original.dat", data, 0777)
			if err != nil {
				panic(err)
			}
		}
	} else {
		g.CurrentLevelNum = 0

		// save current level so we can load it later
		data := []byte(strconv.Itoa(g.CurrentLevelNum))
		switch coverSelectedMode {
		case EasyMode:
			err := os.WriteFile("current_level_easy.dat", data, 0777)
			if err != nil {
				panic(err)
			}
		case OriginalMode:
			err := os.WriteFile("current_level_original.dat", data, 0777)
			if err != nil {
				panic(err)
			}
		}

		g.CurrentScene = EndScene
	}
}

func (g *Game) Update() error {
	switch g.CurrentScene {
	case CoverScene:
		if !loopAudio.IsPlaying() {
			loopAudio.Play()
		}

		HandleInputCover(g)

	case PlayingScene:
		if !g.CurrentLevel.IsCompleted {
			HandleInputPlaying(g)

			if !g.CurrentLevel.IsCompleted {
				g.CurrentLevel.IsCompleted = g.CurrentLevel.IsLevelCompleted()
			}
		} else {
			HandleInputCompleted(g)
		}

	case ExcelScene:
		HandleInputExcel(g)

	case EndScene:
		HandleInputEnd(g)

	case QuitScene:
		return ebiten.Termination
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.CurrentScene {
	case CoverScene:
		screen.DrawImage(coverImage, nil)
		op := &text.DrawOptions{}
		op.GeoM.Translate(850, 600)
		if coverSelectedMode == EasyMode {
			op.ColorScale.ScaleWithColor(color.RGBA{0xff, 0xff, 0x00, 0xff})
		}
		text.Draw(screen, "Easy", &text.GoTextFace{Source: mplusFaceSource, Size: 42}, op)
		op = &text.DrawOptions{}
		op.GeoM.Translate(850, 700)
		if coverSelectedMode == OriginalMode {
			op.ColorScale.ScaleWithColor(color.RGBA{0xff, 0xff, 0x00, 0xff})
		}
		text.Draw(screen, "Original", &text.GoTextFace{Source: mplusFaceSource, Size: 42}, op)
		op = &text.DrawOptions{}
		op.GeoM.Translate(850, 800)
		if coverSelectedMode == QuitMode {
			op.ColorScale.ScaleWithColor(color.RGBA{0xff, 0xff, 0x00, 0xff})
		}
		text.Draw(screen, "Quit", &text.GoTextFace{Source: mplusFaceSource, Size: 42}, op)

	case ExcelScene:
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(gd.TileSize*gd.TilesX)/900, float64(gd.TileSize*gd.TilesY)/679)
		screen.DrawImage(excelImage, op)

	case PlayingScene:
		g.CurrentLevel.Draw(screen, g)

		for _, box := range g.CurrentLevel.Boxes {
			box.Draw(screen, g)
		}

		g.CurrentLevel.Player.Draw(screen, g)

		if g.CurrentLevel.IsCompleted {
			op := &text.DrawOptions{}
			op.GeoM.Translate(450, 550)
			text.Draw(screen, "Level complete!", &text.GoTextFace{Source: mplusFaceSource, Size: 36}, op)
			op = &text.DrawOptions{}
			op.GeoM.Translate(420, 600)
			text.Draw(screen, "Press space to continue...", &text.GoTextFace{Source: mplusFaceSource, Size: 24}, op)
		} else if g.ShowHelp {
			op := &text.DrawOptions{}
			op.GeoM.Translate(700, 250)
			text.Draw(screen, "Help!", &text.GoTextFace{Source: mplusFaceSource, Size: 36}, op)
			op = &text.DrawOptions{}
			op.GeoM.Translate(700, 350)
			op.LayoutOptions.LineSpacing = 40
			text.Draw(screen, "Arrows: move player\nJ: undo movement\nR: restart level\nZ: previous level\nX: toggle Excel\nF: toggle fullscreen\nH: toggle help", &text.GoTextFace{Source: mplusFaceSource, Size: 24}, op)
		}

	case EndScene:
		screen.DrawImage(endImage, nil)
		op := &text.DrawOptions{}
		op.GeoM.Translate(450, 650)
		text.Draw(screen, "Congrats!", &text.GoTextFace{Source: mplusFaceSource, Size: 48}, op)
		op = &text.DrawOptions{}
		op.GeoM.Translate(150, 750)
		text.Draw(screen, "All levels completed", &text.GoTextFace{Source: mplusFaceSource, Size: 48}, op)
		op = &text.DrawOptions{}
		op.GeoM.Translate(380, 900)
		text.Draw(screen, "Press space to continue...", &text.GoTextFace{Source: mplusFaceSource, Size: 24}, op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return gd.TileSize * gd.TilesX, gd.TileSize*gd.TilesY + 20
}

func mustLoadImage(name string) *ebiten.Image {
	f, err := assets.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(img)
}

func mustLoadSingleAudio(name string) *audio.Player {
	input, err := assets.Open(name)
	if err != nil {
		panic(err)
	}

	stream, err := mp3.DecodeWithoutResampling(input)
	if err != nil {
		panic(err)
	}
	singlePlayer, err := audioContext.NewPlayer(stream)
	if err != nil {
		panic(err)
	}

	return singlePlayer
}

func mustLoadLoopAudio(name string) *audio.Player {
	input, err := assets.Open(name)
	if err != nil {
		panic(err)
	}

	stream, err := mp3.DecodeWithoutResampling(input)
	if err != nil {
		panic(err)
	}
	loop := audio.NewInfiniteLoop(stream, stream.Length())
	loopPlayer, err := audioContext.NewPlayer(loop)
	if err != nil {
		panic(err)
	}

	return loopPlayer
}

func main() {
	ebiten.SetWindowSize(800, 690)
	ebiten.SetWindowTitle("SokoMAD")

	ff, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.PressStart2P_ttf))
	if err != nil {
		panic(err)
	}
	mplusFaceSource = ff

	coverImage = mustLoadImage("assets/graphics/cover.png")
	if err != nil {
		panic(err)
	}

	excelImage = mustLoadImage("assets/graphics/excel.png")
	if err != nil {
		panic(err)
	}

	endImage = mustLoadImage("assets/graphics/end.png")
	if err != nil {
		panic(err)
	}

	audioContext = audio.NewContext(24_000)
	stepAudio = mustLoadSingleAudio("assets/sounds/step.mp3")
	loopAudio = mustLoadLoopAudio("assets/sounds/cover.mp3")

	g := NewGame()
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
