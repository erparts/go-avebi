package avebi

import "github.com/hajimehoshi/ebiten/v2"

// A utility function to draw a frame into the given viewport, scaling
// as required with [ebiten.FilterLinear] to take as much space as possible
// while preserving the aspect ratio.
//
// If there's extra space in the viewport, the frame will be drawn centered,
// but black bars won't be explicitly drawn, so whatever was on the background
// of the viewport will remain visible.
//
// Common usage:
//
//	frame, err := videoPlayer.CurrentFrame()
//	if err != nil { /* handle error */ }
//	avebi.Draw(screen, frame)
func Draw(viewport, frame *ebiten.Image) {
	geom, filter := CalcProjection(viewport, frame)
	var opts ebiten.DrawImageOptions
	opts.GeoM = geom
	opts.Filter = filter
	viewport.DrawImage(frame, &opts)
}

// CalcProjection returns the GeoM and recommended ebiten.Filter to project
// the frame into the given viewport. If you don't need the specific parameters,
// see [Draw]() instead.
func CalcProjection(viewport, frame *ebiten.Image) (ebiten.GeoM, ebiten.Filter) {
	// get frame and viewport sizes
	frameBounds := frame.Bounds()
	viewBounds := viewport.Bounds()
	vwWidth, vwHeight := viewBounds.Dx(), viewBounds.Dy()
	frWidth, frHeight := frameBounds.Dx(), frameBounds.Dy()

	// prepare variables for translation to viewport origin
	tx, ty := float64(viewBounds.Min.X), float64(viewBounds.Min.Y)

	// whatever annoying calculations are needed
	var geom ebiten.GeoM
	var filter ebiten.Filter = ebiten.FilterLinear
	wf, hf := float64(vwWidth)/float64(frWidth), float64(vwHeight)/float64(frHeight)
	sf := wf
	if hf < wf {
		sf = hf
	}
	if sf == 1.0 {
		offx := (float64(vwWidth) - float64(frWidth)) / 2
		offy := (float64(vwHeight) - float64(frHeight)) / 2
		geom.Translate(tx+offx, ty+offy)
	} else {
		sfrWidth := float64(frWidth) * sf
		sfrHeight := float64(frHeight) * sf
		geom.Scale(sf, sf)
		geom.Translate(tx+(float64(vwWidth)-sfrWidth)/2, ty+(float64(vwHeight)-sfrHeight)/2)
		filter = ebiten.FilterLinear // TODO: use better filters for new ebitengine versions
	}
	return geom, filter
}
