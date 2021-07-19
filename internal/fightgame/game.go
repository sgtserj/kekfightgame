package fightgame

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math/rand"
	"os"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	raudio "github.com/hajimehoshi/ebiten/v2/examples/resources/audio"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

func Run() {
	ebiten.SetWindowTitle("MK na minimalkax")
	ebiten.SetWindowSize(screenWidth, screenHeight)
	if err := ebiten.RunGame(NewGame()); err != nil {
		panic(err)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func floorDiv(x, y int) int {
	d := x / y
	if d*y == x || x >= 0 {
		return d
	}
	return d - 1
}

func floorMod(x, y int) int {
	return x - floorDiv(x, y)*y
}

const (
	frameOX     = 0
	frameOY     = 0
	frameWidth  = 400
	frameHeight = 1100
	frameNum    = 3
)

const (
	screenWidth   = 640
	screenHeight  = 480
	titleFontSize = fontSize * 1.5
	fontSize      = 24
	smallFontSize = fontSize / 2
)

var (
	fighterOneImage *ebiten.Image
	fighterTwoImage *ebiten.Image
	titleArcadeFont font.Face
	arcadeFont      font.Face
	smallArcadeFont font.Face
)

func init() {
	fighters, err := os.Open("sources/player/fighters.png")
	if err != nil {
		log.Fatal(err)
	}

	img, _, err := image.Decode(fighters)
	if err != nil {
		log.Fatal(err)
	}
	fighterOneImage = ebiten.NewImageFromImage(img)
	fighterTwoImage = ebiten.NewImageFromImage(img)
}

func init() {
	tt, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72
	titleArcadeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    titleFontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	arcadeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	smallArcadeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    smallFontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

var (
	audioContext = audio.NewContext(48000)
	jumpPlayer   *audio.Player
	hitPlayer    *audio.Player
)

func init() {
	jumpD, err := vorbis.Decode(audioContext, bytes.NewReader(raudio.Jump_ogg))
	if err != nil {
		log.Fatal(err)
	}
	jumpPlayer, err = audio.NewPlayer(audioContext, jumpD)
	if err != nil {
		log.Fatal(err)
	}

	jabD, err := wav.Decode(audioContext, bytes.NewReader(raudio.Jab_wav))
	if err != nil {
		log.Fatal(err)
	}
	hitPlayer, err = audio.NewPlayer(audioContext, jabD)
	if err != nil {
		log.Fatal(err)
	}
}

type Mode int

const (
	ModeTitle Mode = iota
	ModeGame
	ModeGameOver
)

type Game struct {
	mode Mode

	// The fighter's position
	xFighterOne       int
	yFighterOne       int
	gravityFighterOne int

	xFighterTwo       int
	yFighterTwo       int
	gravityFighterTwo int

	// Camera
	cameraX int
	cameraY int

	keys       []ebiten.Key
	touchIDs   []ebiten.TouchID
	gamepadIDs []ebiten.GamepadID

	count int
}

func NewGame() *Game {
	g := &Game{}
	g.init()
	return g
}

func (g *Game) init() {
	g.xFighterOne = 100
	g.yFighterOne = 100 * 16
	g.xFighterTwo = 0
	g.yFighterTwo = 100 * 16

	g.cameraX = -240
	g.cameraY = 0
}

func (g *Game) isAnyKeyJustPressed() bool {
	return false
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0xc0, 0xff, 0x80, 0xa0})
	if g.mode != ModeTitle {
		g.drawFighterOne(screen)
		g.drawFighterTwo(screen)
	}
	var titleTexts []string
	var texts []string
	switch g.mode {
	case ModeTitle:
		titleTexts = []string{"MK NA MINIMUME"}
		texts = []string{"", "", "", "", "", "", "", "PRESS ANY KEY OR BUTTON"}
	case ModeGameOver:
		texts = []string{"", "GAME OVER!"}
	}
	for i, l := range titleTexts {
		x := (screenWidth - len(l)*titleFontSize) / 2
		text.Draw(screen, l, titleArcadeFont, x, (i+4)*titleFontSize, color.White)
	}
	for i, l := range texts {
		x := (screenWidth - len(l)*fontSize) / 2
		text.Draw(screen, l, arcadeFont, x, (i+4)*fontSize, color.White)
	}

	if g.mode == ModeTitle {
		msg := []string{
			"PO4ti PoluchiloS",
		}
		for i, l := range msg {
			x := (screenWidth - len(l)*smallFontSize) / 2
			text.Draw(screen, l, smallArcadeFont, x, screenHeight-4+(i-1)*smallFontSize, color.White)
		}
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))
}

func (g *Game) move() bool {
	if g.mode != ModeGame {
		return false
	}

	keys := inpututil.PressedKeys()
	if len(keys) > 0 && keys[0] == ebiten.KeyArrowLeft {
		fmt.Println(keys)
	}
	return false
}

func (g *Game) hit() bool {
	if g.mode != ModeGame {
		return false
	}
	return false
}

func (g *Game) drawFighterOne(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)

	if g.punchFighterOne() {
		op.GeoM.Translate(screenWidth, screenHeight)
		sx, sy := 1200, frameOY
		//w, h := fighterOneImage.Size()
		//op.GeoM.Translate(float64(w)/1.0, float64(h)/1.0)
		op.GeoM.Scale(0.3, 0.3)

		op.GeoM.Translate(float64(g.xFighterOne/16.0)-float64(g.cameraX), float64(g.yFighterOne/16.0)-float64(g.cameraY))

		screen.DrawImage(fighterOneImage.SubImage(image.Rect(sx, sy, sx+frameWidth+150, sy+frameHeight)).(*ebiten.Image), op)

		return
	}

	op.GeoM.Translate(screenWidth, screenHeight)
	i := (g.count / 14) % frameNum
	sx, sy := frameOX+i*frameWidth, frameOY
	//w, h := fighterOneImage.Size()
	//op.GeoM.Translate(float64(w)/1.0, float64(h)/1.0)
	op.GeoM.Scale(0.3, 0.3)

	op.GeoM.Translate(float64(g.xFighterOne/16.0)-float64(g.cameraX), float64(g.yFighterOne/16.0)-float64(g.cameraY))

	screen.DrawImage(fighterOneImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)
}

func (g *Game) drawFighterTwo(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)

	if g.punchFighterTwo() {
		op.GeoM.Translate(screenWidth, screenHeight)
		sx, sy := 1200, 1200
		op.GeoM.Scale(0.3, 0.3)

		op.GeoM.Translate(float64(g.xFighterTwo/16.0)-float64(g.cameraX), float64(g.yFighterTwo/16.0)-float64(g.cameraY))

		screen.DrawImage(fighterTwoImage.SubImage(image.Rect(sx+30, sy, sx+frameWidth+350, sy+frameHeight)).(*ebiten.Image), op)

		return
	}

	op.GeoM.Translate(screenWidth, screenHeight)
	i := (g.count / 14) % frameNum
	sx, sy := frameOX+i*frameWidth, 1200
	op.GeoM.Scale(0.3, 0.3)

	op.GeoM.Translate(float64(g.xFighterTwo/16.0)-float64(g.cameraX), float64(g.yFighterTwo/16.0)-float64(g.cameraY))

	screen.DrawImage(fighterTwoImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)
}
