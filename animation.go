package sprites

import (
	"image"
)

type animation struct {
	entityMode *Mode

	running      bool
	currentFrame int
}

func (a *animation) startAnimation() {
	a.running = true
}

func (a *animation) restartAnimation() {
	a.currentFrame = 0
	a.running = true
}

func (a *animation) resetAnimation() {
	a.currentFrame = 0
	a.running = false
}

func (a *animation) stopAnimation() {
	a.running = false
}

func (a *animation) frame() Sprite {
	if a.running {
		a.currentFrame %= a.entityMode.FrameCount()
		frame, err := a.entityMode.GetFrame(a.currentFrame)
		if err != nil {
			panic(err)
		}
		a.currentFrame++
		// We do this after as well so that any changes to the Mode frame count result in the appropriate next frame
		a.currentFrame %= a.entityMode.FrameCount()
		return frame
	}
	a.currentFrame %= a.entityMode.FrameCount()
	frame, err := a.entityMode.GetFrame(a.currentFrame)
	if err != nil {
		panic(err)
	}
	return frame
}

func (a *animation) frameCount() int {
	return a.entityMode.FrameCount()
}

func (a *animation) spriteSize() image.Rectangle {
	return a.entityMode.SpriteSize()
}
