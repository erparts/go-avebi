package main

import (
	"fmt"
	"image/color"
	"os"

	"github.com/erparts/go-avebi"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go rtsp://<username>:<password>@<ip>:<port>")
		os.Exit(1)
	}

	path := os.Args[1]

	player, err := avebi.NewStreamPlayer(path)
	if err != nil {
		panic(err)
	}

	defer player.Close()

	if err := player.Play(); err != nil {
		panic(err)
	}

	ebiten.SetWindowTitle("Basic Stream Player")
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	g := &game{p: player}
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}

type game struct {
	p     *avebi.Player
	frame *ebiten.Image
}

func (g *game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	f, err := g.p.CurrentFrame()
	if err != nil {
		fmt.Printf("error getting current frame: %v\n", err)
		return nil
	}
	g.frame = f
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	avebi.Draw(screen, g.frame)
}

func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
