// Tutorial from https://earthly.dev/blog/pongo/
package main

import (
	"log"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)

	}
	if err := screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	game := Game{
		screen: screen,
		ball: Ball{
			X:      1,
			Y:      1,
			Xspeed: 1,
			Yspeed: 1,
		},
	}
	go game.Run()

	for {
		switch event := screen.PollEvent().(type) {
		case *tcell.EventResize:
			screen.Sync()
		case *tcell.EventKey:
			if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyCtrlC {
				screen.Fini()
				os.Exit(0)
			}
		}
	}
}

type Game struct {
	screen tcell.Screen
	ball   Ball
}

func (g *Game) Run() {
	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorGreen)

	g.screen.SetStyle(defStyle)
	for {
		g.screen.Clear()
		width, height := g.screen.Size()
		g.ball.CheckEdges(width, height)
		g.ball.Update()
		g.screen.SetContent(g.ball.X, g.ball.Y, g.ball.Display(), nil, defStyle)
		g.screen.Show()
		time.Sleep(30 * time.Millisecond)
	}
}

type Ball struct {
	X, Y           int
	Xspeed, Yspeed int
}

func (b *Ball) Display() rune {
	return '\u25CF'
}

func (b *Ball) Update() {
	b.X += b.Xspeed
	b.Y += b.Yspeed
}

func (b *Ball) CheckEdges(maxWidth, maxHeight int) {
	if b.X <= 0 || b.X >= maxWidth {
		b.Xspeed *= -1
	}
	if b.Y <= 0 || b.Y >= maxHeight {
		b.Yspeed *= -1
	}
}
