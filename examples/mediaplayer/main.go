package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/erparts/avebi"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func main() {
	// get video path
	if len(os.Args) != 2 {
		fmt.Printf("Usage: go run main.go path/to/video.mp4\n")
		os.Exit(1)
	}

	// check that file exists
	path, err := filepath.Abs(os.Args[1])
	if err != nil {
		panic(err)
	}
	_, err = os.Stat(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			fmt.Printf("'%s' not found.", path)
			os.Exit(1)
		}
		panic(err)
	}

	// create video player
	err = avebi.CreateAudioContextForMedia(path)
	if err != nil {
		panic(err)
	}
	videoPlayer, err := avebi.NewPlayer(path)
	if err != nil {
		panic(err)
	}
	videoPlayer.Play()

	// window config
	ebiten.SetWindowTitle("avebi/mediaplayer")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(1280, 720)

	// create test app
	err = ebiten.RunGame(&MediaPlayer{
		videoPath:   path,
		videoPlayer: videoPlayer,
		duration:    videoPlayer.Duration(),
	})
	if err != nil {
		panic(err)
	}
}

type MediaPlayer struct {
	videoPath   string
	videoPlayer *avebi.Player
	videoFrame  *ebiten.Image

	lastPosition time.Duration
	duration     time.Duration
}

func (m *MediaPlayer) Layout(_, _ int) (int, int) {
	panic("Layout() should not be called when LayoutF() exists")
}

func (m *MediaPlayer) LayoutF(w, h float64) (float64, float64) {
	scaleFactor := ebiten.Monitor().DeviceScaleFactor()
	return w * scaleFactor, h * scaleFactor
}

func (m *MediaPlayer) Draw(canvas *ebiten.Image) {
	avebi.Draw(canvas, m.videoFrame)
	m.drawGUI(canvas)
}

func (m *MediaPlayer) Update() error {
	var err error
	m.videoFrame, err = m.videoPlayer.CurrentFrame()
	if err != nil {
		return err
	}

	m.lastPosition, err = m.videoPlayer.Position()
	if err != nil {
		return err
	}
	// ebiten.SetWindowTitle(fmt.Sprintf("%s - %0.2fs", filepath.Base(m.videoPath), position.Seconds()))

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		err := m.videoPlayer.Close()
		if err != nil {
			return err
		}
		return ebiten.Termination
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyP) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		state, err := m.videoPlayer.State()
		if err != nil {
			return err
		}
		if state == avebi.Playing {
			err := m.videoPlayer.Pause()
			if err != nil {
				return err
			}
		} else {
			err := m.videoPlayer.Play()
			if err != nil {
				return err
			}
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		err := m.videoPlayer.Stop()
		if err != nil {
			return err
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyI) {
		state, err := m.videoPlayer.State()
		if err != nil {
			return err
		}
		fmt.Printf("Video state: %s\n", state.String())
	}

	return nil
}

// --- additional info and instructions ---

// TODO: a clean GUI would use a faded darkened area, then light colors and icons for bars and text
func (m *MediaPlayer) drawGUI(canvas *ebiten.Image) {
	bounds := canvas.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	playWidth := (w * 2) / 3
	playHeight := h / 48
	ox := (w - playWidth) / 2
	oy := h - playHeight*2
	playRect := image.Rect(ox, oy, ox+playWidth, oy+playHeight)
	canvas.SubImage(playRect).(*ebiten.Image).Fill(color.RGBA{255, 255, 255, 255})
	const BorderThickness = 3
	playRect.Min.X += BorderThickness
	playRect.Max.X -= BorderThickness
	playRect.Min.Y += BorderThickness
	playRect.Max.Y -= BorderThickness
	canvas.SubImage(playRect).(*ebiten.Image).Fill(color.RGBA{0, 0, 0, 255})
	const InnerMargin = 2
	playRect.Min.X += InnerMargin
	playRect.Max.X -= InnerMargin
	playRect.Min.Y += InnerMargin
	playRect.Max.Y -= InnerMargin
	t := float64(m.lastPosition) / float64(m.duration)
	playRect.Max.X = playRect.Min.X + int(float64(playRect.Dx())*t)
	canvas.SubImage(playRect).(*ebiten.Image).Fill(color.RGBA{255, 255, 255, 255})

	positionStr := durationToMMSS(m.lastPosition)
	durationStr := durationToMMSS(m.duration)
	ebitenutil.DebugPrintAt(canvas, positionStr+" / "+durationStr+" (SPACE to pause, S to stop)", ox, oy-16)
}

func durationToMMSS(duration time.Duration) string {
	millis := duration.Milliseconds()
	rem := millis
	seconds := rem / 1000
	minutes := seconds / 60
	seconds = seconds % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}
