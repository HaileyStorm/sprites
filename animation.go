package sprites

import (
	"errors"
	"image"

	ccsl_graphics "github.com/HaileyStorm/CCSL_go/graphics"
)

type animation struct {
	*Mode

	running      bool
	advanceEvery int
	advanceCt    int
	currentFrame int
}

func (a *animation) Running() bool {
	return a.running
}

func (a *animation) StartAnimation() {
	a.advanceCt = 0
	a.running = true
}

func (a *animation) ResumeAnimation() {
	a.running = true
}

func (a *animation) RestartAnimation() {
	a.advanceCt = 0
	a.currentFrame = 0
	a.running = true
}

func (a *animation) ResetAnimation() {
	a.advanceCt = 0
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
	if a.running {
		if a.advanceCt == 0 {
			a.currentFrame++
			// We do this after as well so that any changes to the Mode frame count result in the appropriate next frame
			a.currentFrame %= a.FrameCount()
		}
		a.advanceCt++
		a.advanceCt %= a.advanceEvery
	}

	return frame
}

func (a *animation) FrameResized(w, h uint) Sprite {
	frame := a.Frame()
	return ccsl_graphics.ResizeMaintain(frame.(*image.RGBA), w, h)
}

func (a *animation) Advance() {
	if !a.running {
		return
	}
	a.currentFrame %= a.FrameCount()
	if a.advanceCt == 0 {
		a.currentFrame++
		// We do this after as well so that any changes to the Mode frame count result in the appropriate next frame
		a.currentFrame %= a.FrameCount()
	}
	a.advanceCt++
	a.advanceCt %= a.advanceEvery
}

// NextFrameDiff returns true if the next call to Frame will return a different frame than the previous count (or, the
// next call will be the first since Starting/Restarting/Resetting the animation).
func (a *animation) NextFrameDiff() bool {
	return a.running && a.advanceCt == 0 && len(a.frames) > 0
}

func (a *animation) AdvanceEvery() int {
	return a.advanceEvery
}

func (a *animation) SetAdvanceEvery(ct int) error {
	if ct <= 0 {
		return errors.New("new advanceEvery count must be > 0")
	}
	a.advanceEvery = ct
	return nil
}
