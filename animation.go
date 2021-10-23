package sprites

import (
	"image"

	ccsl_graphics "github.com/HaileyStorm/CCSL_go/graphics"
)

type animation struct {
	*Mode

	running      bool
	currentFrame int
}

func (a *animation) Running() bool {
	return a.running
}

func (a *animation) StartAnimation() {
	a.running = true
}

func (a *animation) RestartAnimation() {
	a.currentFrame = 0
	a.running = true
}

func (a *animation) ResetAnimation() {
	a.currentFrame = 0
	a.running = false
}

func (a *animation) StopAnimation() {
	a.running = false
}

func (a *animation) Frame() Sprite {
	var frame Sprite
	var err error

	a.currentFrame %= a.FrameCount()
	frame, err = a.GetFrame(a.currentFrame)
	if err != nil {
		panic(err)
	}
	a.Advance()

	return frame
}

func (a *animation) FrameResized(w, h uint) Sprite {
	frame := a.Frame()
	return ccsl_graphics.ResizeMaintain(frame.(*image.RGBA), w, h)
}

func (a *animation) Advance() {
	if a.running {
		a.currentFrame++
		// We do this after as well so that any changes to the Mode frame count before the next call to Frame will
		// result in the appropriate next frame
		a.currentFrame %= a.FrameCount()
	}
}
