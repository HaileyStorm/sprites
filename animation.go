package sprites

type animation struct {
	*Mode

	running      bool
	currentFrame int
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
	if a.running {
		a.currentFrame %= a.FrameCount()
		frame, err := a.GetFrame(a.currentFrame)
		if err != nil {
			panic(err)
		}
		a.currentFrame++
		// We do this after as well so that any changes to the Mode frame count result in the appropriate next frame
		a.currentFrame %= a.FrameCount()
		return frame
	}
	a.currentFrame %= a.FrameCount()
	frame, err := a.GetFrame(a.currentFrame)
	if err != nil {
		panic(err)
	}
	return frame
}
