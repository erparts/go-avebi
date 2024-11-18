# avebi

A video playing library for Ebitengine based on the [erparts/reisen] fork of [zergon321/reisen](https://github.com/zergon321/reisen), with an API design inspired by [tinne26/mpegg](https://github.com/tinne26/mpegg).

Warnings and limitations:
- The library is still quite barebones, lacking testing and only trying to cover primary needs for erparts.
- Reisen uses cgo, so this library inherits the problem (consider purego).
- The dependency on ffmpeg6.1 is very undesirable and problematic for casual projects.
- The `erparts/reisen` fork is only adapted for Linux, so multi-platform support is non-existent.
- In order to play video with audio, the audio channels need to be stereo and the same length as the video channels.

## Dependencies

Reisen depends on ffmpeg6.1, which is currently an outdated ffmpeg version.

On a linux system, the choices are:
- **Keeping only the old ffmpeg version**: which is not viable on up-to-date personal use computers where you have other programs that depend on the newest ffmpeg version.
- **Keeping ffmpeg6.1 alongside newer versions**: and using `PKG_CONFIG_PATH=/usr/lib/ffmpeg6.1/pkgconfig` or similar to point to it.

## Usage

```Golang
func main() {
    // create video player
    videoPlayer, err := avebi.NewPlayer("../test_video.mp4")
    if err != nil {
        panic(err)
    }

    // start playing
    videoPlayer.Play()

    // ... (ebiten.RunGame)
}

// ... (game definition)

func (g *Game) Draw(canvas *ebiten.Image) {
	avebi.Draw(canvas, g.videoFrame)
}

func (g *Game) Update() error {
    var err error
    g.videoFrame, err = g.videoPlayer.CurrentFrame()
	return err
}
```

## TODO

Potential improvements on avebi:
- Implement audio looping.
- Use circular buffers for audio data.
- Add support for mono audio.
- Add support for videos with audio and video channels of different lengths (annoying).

Potential improvements on reisen:
- Add support hardware acceleration, ¿..primarily h264_v4l2m2m for the raspberry pi?
- Use pools for both video and audio frames data.