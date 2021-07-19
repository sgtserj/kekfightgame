package fightgame

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) jumpFighterOne() bool {
	keys := inpututil.PressedKeys()

	if len(keys) > 0 && keys[0] == ebiten.KeyArrowUp {
		return true
	}
	return false
}

func (g *Game) jumpFighterTwo() bool {
	keys := inpututil.PressedKeys()

	if len(keys) > 0 && keys[0] == ebiten.KeyW {
		return true
	}
	return false
}

func (g *Game) punchFighterOne() bool {
	keys := inpututil.PressedKeys()
	if len(keys) > 0 && keys[0] == ebiten.KeyL {
		return true
	}
	return false
}
func (g *Game) punchFighterTwo() bool {
	keys := inpututil.PressedKeys()
	if len(keys) > 0 && keys[0] == ebiten.KeyE {
		return true
	}
	return false
}
