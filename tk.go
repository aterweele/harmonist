// +build tk

package main

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/draw"
	"image/png"
	"unicode/utf8"

	"github.com/nsf/gothic"
)

type gameui struct {
	g         *game
	ir        *gothic.Interpreter
	cursor    position
	stty      string
	cache     map[UICell]*image.RGBA
	width     int
	height    int
	mousepos  position
	menuHover menu
	itemHover int
	canvas    *image.RGBA
}

func (ui *gameui) Init() error {
	ui.canvas = image.NewRGBA(image.Rect(0, 0, UIWidth*16, UIHeight*24))
	ui.ir = gothic.NewInterpreter(`
set width [expr {16 * 100}]
set height [expr {24 * 26}]
set can [canvas .c -width $width -height $height -background black]
grid $can -row 0 -column 0
focus $can
image create photo gamescreen -width $width -height $height -palette 256/256/256
image create photo bufscreen -width $width -height $height -palette 256/256/256
$can create image 0 0 -anchor nw -image gamescreen
bind $can <Key> {
	GetKey %A %K
}
bind $can <Motion> {
	MouseMotion %x %y
}
bind $can <ButtonPress> {
	MouseDown %x %y %b
}
`)
	ui.ir.RegisterCommand("GetKey", func(c, keysym string) {
		var s string
		if c != "" {
			s = c
		} else {
			s = keysym
		}
		ch <- uiInput{key: s}
	})
	ui.ir.RegisterCommand("MouseDown", func(x, y, b int) {
		ch <- uiInput{mouse: true, mouseX: (x - 1) / ui.width, mouseY: (y - 1) / ui.height, button: b - 1}
	})
	ui.ir.RegisterCommand("MouseMotion", func(x, y int) {
		nx := (x - 1) / ui.width
		ny := (y - 1) / ui.height
		if nx != ui.mousepos.X || ny != ui.mousepos.Y {
			ui.mousepos.X = nx
			ui.mousepos.Y = ny
			ch <- uiInput{mouse: true, mouseX: nx, mouseY: ny, button: -1}
		}
	})
	ui.menuHover = -1
	ui.InitElements()

	SolarizedPalette()
	ui.HideCursor()
	settingsActions = append(settingsActions, toggleTiles)
	gameConfig.Tiles = true
	return nil
}

func (ui *gameui) InitElements() error {
	ui.width = 16
	ui.height = 24
	ui.cache = make(map[UICell]*image.RGBA)
	return nil
}

var ch chan uiInput
var interrupt chan bool

func init() {
	ch = make(chan uiInput, 100)
	interrupt = make(chan bool)
}

func (ui *gameui) Close() {
}

func (ui *gameui) Flush() {
	ui.DrawLogFrame()
	xmin := UIWidth - 1
	xmax := 0
	ymin := UIHeight - 1
	ymax := 0
	for _, cdraw := range ui.g.DrawLog[len(ui.g.DrawLog)-1].Draws {
		cell := cdraw.Cell
		i := cdraw.I
		x, y := ui.GetPos(i)
		ui.Draw(cell, x, y)
		if x < xmin {
			xmin = x
		}
		if x > xmax {
			xmax = x
		}
		if y < ymin {
			ymin = y
		}
		if y > ymax {
			ymax = y
		}
	}
	if xmin > xmax || ymin > ymax {
		return
	}
	pngbuf := &bytes.Buffer{}
	subimg := ui.canvas.SubImage(image.Rect(xmin*16, ymin*24, (xmax+1)*16, (ymax+1)*24))
	png.Encode(pngbuf, subimg)
	png := base64.StdEncoding.EncodeToString(pngbuf.Bytes())
	ui.ir.Eval("gamescreen put %{0%s} -format png -to %{1%d} %{2%d} %{3%d} %{4%d}", png,
		xmin*16, ymin*24, (xmax+1)*16, (ymax+1)*24) // TODO: optimize this more
}

func (ui *gameui) ApplyToggleLayout() {
	gameConfig.Small = !gameConfig.Small
	if gameConfig.Small {
		ui.Clear()
		ui.Flush()
		UIHeight = 24
		UIWidth = 80
	} else {
		UIHeight = 26
		UIWidth = 100
	}
	ui.cache = make(map[UICell]*image.RGBA)
	ui.g.DrawBuffer = make([]UICell, UIWidth*UIHeight)
	ui.Clear()
}

func (ui *gameui) Draw(cell UICell, x, y int) {
	var img *image.RGBA
	if im, ok := ui.cache[cell]; ok {
		img = im
	} else {
		img = getImage(cell)
		ui.cache[cell] = img
	}
	draw.Draw(ui.canvas, image.Rect(x*ui.width, ui.height*y, (x+1)*ui.width, (y+1)*ui.height), img, image.Point{0, 0}, draw.Over)
}

func (ui *gameui) KeyToRuneKeyAction(in uiInput) rune {
	switch in.key {
	case "Enter":
		in.key = "."
	case "Left", "KP_Left":
		in.key = "4"
	case "Right", "KP_Right":
		in.key = "6"
	case "Up", "KP_Up":
		in.key = "8"
	case "Down", "KP_Down":
		in.key = "2"
	case "KP_Home":
		in.key = "7"
	case "KP_End":
		in.key = "1"
	case "KP_Prior":
		in.key = "9"
	case "KP_Next":
		in.key = "3"
	case "KP_Begin", "KP_Delete":
		in.key = "5"
	}
	if utf8.RuneCountInString(in.key) > 1 {
		return 0
	}
	return ui.ReadKey(in.key)
}
