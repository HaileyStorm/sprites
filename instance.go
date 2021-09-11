package sprites

import (
	"fmt"
	"image"
	"image/draw"
)

type Instance struct {
	//optional
	name						string

	entity						*Entity

	animation					Animation
}

func (i *Instance) Name() string {
	return i.name
}

func (i *Instance) SetName(name string) {
	i.name = name
}

func (i *Instance) Entity() *Entity {
	return i.entity
}

func (i *Instance) Mode() *Mode {
	return i.animation.entityMode
}

//note in docstrings that changing mode does NOT stop or restart the animation
// (if it was running, it still will be, and the currentFrame will be the same and Frame will get that frame from the
// new mode - except that currentFrame is modulo'd with the len(frames) to ensure it's in range)
func (i *Instance) SetModeByIndex(index int) error {
	if mode, ok := i.entity.modes[index]; ok {
		i.animation.entityMode = mode
		return nil
	} else {
		return fmt.Errorf("mode with index %d does not exist in instance Entity", index)
	}
}

//note in docstrings that changing mode does NOT stop or restart the animation
// (if it was running, it still will be, and the currentFrame will be the same and Frame will get that frame from the
// new mode - except that currentFrame is modulo'd with the len(frames) to ensure it's in range)
func (i *Instance) SetModeByName(name string) error {
	idx, ok := i.entity.modeNamesToIndex[name]; if ok {
		mode, ok := i.entity.modes[idx]; if ok {
			i.animation.entityMode = mode
			return nil
		} else {
			panic(fmt.Errorf("internal error: Mode with index %d does not exist in Entity; Entity is corrupted", idx))
		}
	} else {
		return fmt.Errorf("mode with name %s does not exist in Entity", name)
	}

}

func (i *Instance) ModeCount() int {
	return i.entity.ModeCount()
}

func (i *Instance) StartAnimation() {
	i.animation.startAnimation()
}

func (i *Instance) RestartAnimation() {
	i.animation.restartAnimation()
}

func (i *Instance) ResetAnimation() {
	i.animation.resetAnimation()
}

func (i *Instance) StopAnimation() {
	i.animation.stopAnimation()
}

func (i *Instance) Frame() Sprite {
	return i.animation.frame()
}

func (i *Instance) FrameCount() int {
	return i.animation.frameCount()
}

func (i *Instance) GetFrame(idx int) (Sprite, error) {
	return i.Mode().GetFrame(idx)
}

func (i *Instance) SpriteSize() image.Rectangle {
	return i.animation.spriteSize()
}

// note that it gets next frame and places that. To not advance the animation, first stop it and then call this (and then start it
func (i *Instance) PlaceSprite(canvas draw.Image, placeAt image.Point) {
	frame := i.Frame()

	// SpriteSize (Rect) + Point = rect translated (placed at) Point. This is placement location on dst. The zero point + frame.Bounds().Min is the rect in source to grab
	// (this is the only area on the source - frame - that has data, but has to be done because Bounds() does not always start at (0,0) - indeed if made from a SubImage it doesn't unless the location on the original started at (0,0))
	draw.Draw(canvas, i.SpriteSize().Add(placeAt), frame, image.Pt(0, 0).Add(frame.Bounds().Min), draw.Over)
}