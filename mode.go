package sprites

import (
	"errors"
	"fmt"
	"image"

	"github.com/corona10/goimagehash"
)

type Sprite image.Image

type Mode struct {
	name string

	spriteSize  image.Rectangle
	fullyOpaque bool

	frames []Sprite
}

func (m *Mode) Name() string {
	return m.name
}

func (m *Mode) SpriteSize() image.Rectangle {
	//return (m.frames[0]).Bounds().Sub((m.frames[0]).Bounds().Min)
	return m.spriteSize
}

func (m *Mode) FullyOpaque() bool {
	return m.fullyOpaque
}

//note that unlike Instance.Frame() this does not advance the current frame (there is no current frame in Mode - this is an Instance concept)
func (m *Mode) GetFrame(index int) (Sprite, error) {
	if index < len(m.frames) {
		return m.frames[index], nil
	} else {
		return nil, errors.New("index out of bounds")
	}
}

func (m *Mode) FrameCount() int {
	return len(m.frames)
}

//only decrease
func (m *Mode) SetFrameCount(count int) error {
	if count > 0 && count <= len(m.frames) {
		m.frames = m.frames[0:count]
		return nil
	} else {
		return fmt.Errorf("new frame count (%d) must be <= the current frame count (%d) and > 0", count, len(m.frames))
	}
}

// SpriteHash gets a string hash representation of sprite, using the average hash algorithm.
//
// License(s) - see internal\licenses:
// goimagehash
func SpriteHash(sprite Sprite) string {
	var hashstr string
	defer func() {
		if r := recover(); r != nil {
			hashstr = "hash index out of bounds error"
		}
	}()
	hash, e := goimagehash.AverageHash(sprite)
	if e != nil {
		hashstr = e.Error()
	} else {
		hashstr = hash.ToString()
	}
	return hashstr
}
