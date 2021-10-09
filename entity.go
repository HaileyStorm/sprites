package sprites

import (
	"errors"
	"fmt"
	"image"
)

type Entity struct {
	name string

	modes            map[int]*Mode
	modeNamesToIndex map[string]int
}

func (e *Entity) Name() string {
	return e.name
}

//describe index order in docstring
func (e *Entity) GetModeByIndex(idx int) (*Mode, error) {
	mode, ok := e.modes[idx]
	if ok {
		return mode, nil
	} else {
		return nil, fmt.Errorf("mode with index %d does not exist in Entity", idx)
	}
}

func (e *Entity) GetModeByName(name string) (*Mode, error) {
	idx, ok := e.modeNamesToIndex[name]
	if ok {
		mode, ok := e.modes[idx]
		if ok {
			return mode, nil
		} else {
			panic(fmt.Errorf("internal error: Mode with index %d does not exist in Entity; Entity is corrupted", idx))
		}

	} else {
		return nil, fmt.Errorf("mode with name %s does not exist in Entity", name)
	}
}

func (e *Entity) RenameMode(oldName, newName string) error {
	idx, ok := e.modeNamesToIndex[oldName]
	if ok {
		mode, ok := e.modes[idx]
		if ok {
			mode.name = newName
			e.modeNamesToIndex[newName] = idx
			delete(e.modeNamesToIndex, oldName)
		} else {
			panic(fmt.Errorf("internal error: Mode with index %d does not exist in Entity; Entity is corrupted", idx))
		}
	} else {
		return fmt.Errorf("mode with name %s does not exist in Entity", oldName)
	}
	return nil
}

func (e *Entity) ModeCount() int {
	return len(e.modes)
}

//only decrease
func (e *Entity) SetModeCount(count int) error {
	if count > 0 && count <= len(e.modes) {
		var delList []string
		for k, v := range e.modeNamesToIndex {
			if v >= count {
				delList = append(delList, k)
			}
		}
		for _, d := range delList {
			delete(e.modeNamesToIndex, d)
		}
		for i := count; i < len(e.modes); i++ {
			delete(e.modes, i)
		}
		return nil
	} else {
		return fmt.Errorf("new mode count (%d) must be <= the current mode count (%d) and > 0", count, len(e.modes))
	}
}

func (e *Entity) SpriteSize() image.Rectangle {
	return e.modes[0].SpriteSize()
}

func (e *Entity) NewInstance(initialMode int, advanceEvery int) (*Instance, error) {
	if mode, ok := e.modes[initialMode]; ok {
		if advanceEvery <= 0 {
			return nil, errors.New("advanceEvery must be > 0")
		}
		return &Instance{
			Entity: e,
			animation: &animation{
				Mode:         mode,
				running:      false,
				advanceEvery: advanceEvery,
				advanceCt:    0,
				currentFrame: 0,
			},
		}, nil
	} else {
		return nil, fmt.Errorf("mode with index %d does not exist in Entity", initialMode)
	}
}

func (e *Entity) NewInstanceWithModeName(initialMode string, advanceEvery int) (*Instance, error) {
	if idx, ok := e.modeNamesToIndex[initialMode]; ok {
		if instance, err := e.NewInstance(idx, advanceEvery); err == nil {
			return instance, nil
		} else {
			panic(fmt.Errorf("internal error:%v", err))
		}
	} else {
		return nil, fmt.Errorf("mode with name %s does not exist in Entity", initialMode)
	}
}
