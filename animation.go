package sprites

import (
	"image"
)

type Animation struct {
	entityMode					*Mode

	running						bool
	currentFrame				int
}

func (a *Animation) startAnimation() {
	a.running = true
}

func (a *Animation) restartAnimation() {
	a.currentFrame = 0
	a.running = true
}

func (a *Animation) resetAnimation() {
	a.currentFrame = 0
	a.running = false
}

func (a *Animation) stopAnimation() {
	a.running = false
}

func (a *Animation) frame() Sprite {
	if a.running {
		a.currentFrame %= a.entityMode.FrameCount()
		frame, err := a.entityMode.GetFrame(a.currentFrame)
		if err != nil { panic(err) }
		a.currentFrame++
		// We do this after as well so that any changes to the Mode frame count result in the appropriate next frame
		a.currentFrame %= a.entityMode.FrameCount()
		return frame
	}
	a.currentFrame %= a.entityMode.FrameCount()
	frame, err := a.entityMode.GetFrame(a.currentFrame)
	if err != nil { panic(err) }
	return frame
}

func (a *Animation) frameCount() int {
	return a.entityMode.FrameCount()
}

func (a *Animation) spriteSize() image.Rectangle {
	return a.entityMode.SpriteSize()
}