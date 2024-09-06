package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func repeatingKeyPressed(key ebiten.Key) bool {
	const (
		delay    = 30
		interval = 10
	)
	d := inpututil.KeyPressDuration(key)
	if d == 1 {
		return true
	}
	if d >= delay && (d-delay)%interval == 0 {
		return true
	}
	return false
}

func HandleInputCover(g *Game) {
	if repeatingKeyPressed(ebiten.KeyDown) {
		switch coverSelectedMode {
		case EasyMode:
			coverSelectedMode = OriginalMode
		case OriginalMode:
			coverSelectedMode = QuitMode
		case QuitMode:
			coverSelectedMode = EasyMode
		}
	}

	if repeatingKeyPressed(ebiten.KeyUp) {
		switch coverSelectedMode {
		case EasyMode:
			coverSelectedMode = QuitMode
		case OriginalMode:
			coverSelectedMode = EasyMode
		case QuitMode:
			coverSelectedMode = OriginalMode
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		switch coverSelectedMode {
		case EasyMode, OriginalMode:
			g.Start()
		case QuitMode:
			g.CurrentScene = QuitScene
		}
	}
}

func HandleInputExcel(g *Game) {
	if inpututil.IsKeyJustPressed(ebiten.KeyX) {
		g.CurrentScene = PlayingScene
	}
}

func HandleInputPlaying(g *Game) {
	if repeatingKeyPressed(ebiten.KeyDown) {
		g.CurrentLevel.Player.MoveDown(g)
	}

	if repeatingKeyPressed(ebiten.KeyUp) {
		g.CurrentLevel.Player.MoveUp(g)
	}

	if repeatingKeyPressed(ebiten.KeyLeft) {
		g.CurrentLevel.Player.MoveLeft(g)
	}

	if repeatingKeyPressed(ebiten.KeyRight) {
		g.CurrentLevel.Player.MoveRight(g)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyJ) {
		g.CurrentLevel.RemoveMovement()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.RestartLevel()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		g.PreviousLevel()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyX) {
		g.CurrentScene = ExcelScene
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyH) {
		g.ShowHelp = !g.ShowHelp
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		g.RestartMode()
	}
}

func HandleInputCompleted(g *Game) {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.NextLevel()
	}
}

func HandleInputEnd(g *Game) {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.RestartMode()
	}
}
