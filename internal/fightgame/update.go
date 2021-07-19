package fightgame

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) Update() error {
	switch g.mode {
	case ModeTitle:
		if g.jumpFighterOne() {
			g.mode = ModeGame
		}
	case ModeGame:
		g.count++
		//g.xFighterOne = 32
		//g.cameraX += 2
		//g.gravityFighterOne = 100
		if g.jumpFighterOne() {
			g.gravityFighterOne = -96
			jumpPlayer.Rewind()
			jumpPlayer.Play()
		}
		g.yFighterOne += g.gravityFighterOne

		if g.yFighterOne > 3500 {
			g.yFighterOne = 3500
		}

		if g.jumpFighterTwo() {
			g.gravityFighterTwo = -96
			jumpPlayer.Rewind()
			jumpPlayer.Play()
		}

		g.yFighterTwo += g.gravityFighterTwo

		if g.yFighterTwo > 3500 {
			g.yFighterTwo = 3500
		}

		//moving
		keys := inpututil.PressedKeys()
		if len(keys) > 0 && keys[0] == ebiten.KeyArrowLeft {
			g.xFighterOne -= 32
		}
		if len(keys) > 0 && keys[0] == ebiten.KeyArrowRight {
			g.xFighterOne += 32
		}
		if len(keys) > 0 && keys[0] == ebiten.KeyA {
			g.xFighterTwo -= 32
		}
		if len(keys) > 0 && keys[0] == ebiten.KeyD {
			g.xFighterTwo += 32
		}

		// Gravity
		g.gravityFighterOne += 4
		if g.gravityFighterOne > 96 {
			g.gravityFighterOne = 96
		}
		g.gravityFighterTwo += 4
		if g.gravityFighterTwo > 96 {
			g.gravityFighterTwo = 96
		}

		if g.hit() {
			hitPlayer.Rewind()
			hitPlayer.Play()
			g.mode = ModeGameOver
		}
	case ModeGameOver:
	}
	return nil
}
