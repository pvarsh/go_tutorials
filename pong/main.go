// Tutorial from https://earthly.dev/blog/pongo/
package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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
	width, height := screen.Size()
	paddleSpeed := 5

	game := Game{
		screen: screen,
		ball: Ball{
			X:      6,
			Y:      1,
			Xspeed: 1,
			Yspeed: 1,
		},
		player1: Paddle{
			X:      5,
			Y:      0,
			height: 6,
			Yspeed: paddleSpeed,
		},
		player2: Paddle{
			X:      width - 5,
			Y:      0,
			height: 6,
			Yspeed: paddleSpeed,
		},
		player1Serves: true,
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
			} else if event.Key() == tcell.KeyUp {
				game.player2.MoveUp()
			} else if event.Key() == tcell.KeyDown {
				game.player2.MoveDown(height)
			} else if event.Rune() == 'w' {
				game.player1.MoveUp()
			} else if event.Rune() == 's' {
				game.player1.MoveDown(height)
			}
		}
	}
}

type Game struct {
	screen        tcell.Screen
	ball          Ball
	player1       Paddle
	player2       Paddle
	player1Serves bool
}

func (g *Game) Run() {
	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorGreen)
	paddleStyle := tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorWhite)
	g.screen.SetStyle(defStyle)
	g.NewBall()
	for {
		g.screen.Clear()
		width, height := g.screen.Size()
		g.ball.CheckEdges(width, height)
		if g.ball.X <= 0 {
			g.player2.score += 1
			g.NewBall()
		}
		if g.ball.X >= width {
			g.player1.score += 1
			g.NewBall()
		}
		if g.player1.score >= 2 || g.player2.score >= 2 {
			break
		}
		if g.ball.intersects(g.player1) || g.ball.intersects(g.player2) {
			g.ball.Xspeed *= -1
		}
		g.ball.Update()
		drawSprite(g.screen, 10, 1, 10, 1, defStyle, strconv.Itoa(g.player1.score))
		drawSprite(g.screen, width-10, 1, width-10, 1, defStyle, strconv.Itoa(g.player2.score))
		drawSprite(g.screen, g.ball.X, g.ball.Y, g.ball.X, g.ball.Y, defStyle, g.ball.Display())
		drawSprite(
			g.screen,
			g.player1.X,
			g.player1.Y,
			g.player1.X+g.player1.width,
			g.player1.Y+g.player1.height,
			paddleStyle,
			g.player1.Display(),
		)
		drawSprite(
			g.screen,
			g.player2.X,
			g.player2.Y,
			g.player2.X+g.player2.width,
			g.player2.Y+g.player2.height,
			paddleStyle,
			g.player2.Display(),
		)
		g.screen.Show()
		time.Sleep(30 * time.Millisecond)
	}
	g.screen.Clear()
	gameOver := "Game over!"
	winner := 1
	if g.player2.score > g.player1.score {
		winner = 2
	}
	playerWins := fmt.Sprintf("Player %d wins!", winner)
	width, height := g.screen.Size()
	gameOverLeft := width/2 - len(gameOver)/2
	drawSprite(g.screen, gameOverLeft, height/2, gameOverLeft+len(gameOver), height/2, defStyle, gameOver)
	playerWinsLeft := width/2 - len(playerWins)/2
	drawSprite(g.screen, playerWinsLeft, height/2+1, playerWinsLeft+len(playerWins), height/2+1, defStyle, playerWins)
	g.screen.Show()
	time.Sleep(time.Second)
}

func (g *Game) NewBall() {
	if g.player1Serves {
		g.ball = Ball{
			6,
			1,
			1,
			1,
		}
		g.player1Serves = false
	} else {
		screenWidth, _ := g.screen.Size()
		g.ball = Ball{
			screenWidth - 8,
			1,
			-1,
			1,
		}
		g.player1Serves = true
	}
}

type Paddle struct {
	X      int
	Y      int
	Yspeed int
	width  int
	height int
	score  int
}

func (p *Paddle) Display() string {
	return strings.Repeat(" ", p.height)
}

func (p *Paddle) MoveUp() {
	// TODO: what if speed > 1?
	if p.Y > 0 {
		p.Y -= p.Yspeed
	}
}

func (p *Paddle) MoveDown(windowHeight int) {
	if p.Y < windowHeight-p.height {
		p.Y += p.Yspeed
	}
}

type Ball struct {
	X, Y           int
	Xspeed, Yspeed int
}

func (b *Ball) Display() string {
	return string('\u25CF')
}

func (b *Ball) Update() {
	b.X += b.Xspeed
	b.Y += b.Yspeed
}

func (b *Ball) CheckEdges(maxWidth, maxHeight int) {
	// if b.X <= 0 || b.X >= maxWidth {
	// 	b.Xspeed *= -1
	// }
	if b.Y <= 0 || b.Y >= maxHeight {
		b.Yspeed *= -1
	}
}

func (b *Ball) intersects(p Paddle) bool {
	return b.X == p.X && b.Y >= p.Y && b.Y <= (p.Y+p.height)
	// return b.X >= p.X && b.X <= p.X+p.width && b.Y >= p.Y && b.Y <= p.Y+p.height
}

// x1, y1 are upper left coordinates, x2, y2 are bottom right coordinates for the sprite
func drawSprite(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range text {
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}
