package sprites

import (
	"fmt"
	"image"
	"image/draw"

	ccsl_graphics "github.com/HaileyStorm/CCSL_go/graphics"
)

type Instance struct {
	//optional
	name string

	*Entity

	*animation
}

func (i *Instance) Name() string {
	return i.name
}

func (i *Instance) SetName(name string) {
	i.name = name
}

//note in docstrings that changing mode does NOT stop or restart the animation
// (if it was running, it still will be, and the currentFrame will be the same and Frame will get that frame from the
// new mode - except that currentFrame is modulo'd with the len(frames) to ensure it's in range)
func (i *Instance) SetModeByIndex(index int) error {
	if mode, ok := i.modes[index]; ok {
		i.Mode = mode
		return nil
	} else {
		return fmt.Errorf("mode with index %d does not exist in instance Entity", index)
	}
}

//note in docstrings that changing mode does NOT stop or restart the animation
// (if it was running, it still will be, and the currentFrame will be the same and Frame will get that frame from the
// new mode - except that currentFrame is modulo'd with the len(frames) to ensure it's in range)
func (i *Instance) SetModeByName(name string) error {
	idx, ok := i.modeNamesToIndex[name]
	if ok {
		mode, ok := i.modes[idx]
		if ok {
			i.Mode = mode
			return nil
		} else {
			panic(fmt.Errorf("internal error: Mode with index %d does not exist in Entity; Entity is corrupted", idx))
		}
	} else {
		return fmt.Errorf("mode with name %s does not exist in Entity", name)
	}

}

// note that placeAt is expected to be within canvas.Bounds() (that is, not necessarily relative to (0,0))
// note that it gets next frame and places that. To not advance the animation, first stop it and then call this (and then start it
func (i *Instance) PlaceSprite(canvas draw.Image, placeAt image.Point) {
	frame := i.Frame()

	// SpriteSize (Rect) + Point = rect translated (placed at) Point. This is placement location on dst. The zero point + frame.Bounds().Min is the rect in source to grab
	// (this is the only area on the source - frame - that has data, but has to be done because Bounds() does not always start at (0,0) - indeed if made from a SubImage it doesn't unless the location on the original started at (0,0))
	// If frame is fully opaque, we can use one of two faster methods to place it on canvas. If not, we must use
	// draw.Draw with draw.Over to respect the transparencies in combining it with canvas.
	if i.Mode.fullyOpaque {
		var img *ccsl_graphics.Image
		var ok bool
		// If canvas is a ccsl_graphics.Image, we can use the specialized/simplified PlaceAtPoint instead of draw.Draw,
		// which is much faster (even with draw.Src and nil mask).
		if img, ok = canvas.(*ccsl_graphics.Image); ok {
			img.PlaceAtPoint(frame.(*image.RGBA), placeAt)
		} else {
			draw.Draw(canvas, i.SpriteSize().Add(placeAt), frame, frame.Bounds().Min, draw.Src)
		}
	} else {
		draw.Draw(canvas, i.SpriteSize().Add(placeAt), frame, frame.Bounds().Min, draw.Over)
	}
}
