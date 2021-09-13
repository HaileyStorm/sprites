package sprites

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	"strconv"
)

// subImager is used to ensure the Image provided for a Sheet is one of the image types which receives the SubImage
// and Set methods.
// For simplicity / user comprehensibility the image parameter in the NewSheet factory is not a subImager but an Image,
// and the factory uses a type assertion and returns (nil, error) if it fails (as opposed to (Sheet, nil)).
type subImager interface {
	draw.Image
	SubImage(r image.Rectangle) image.Image
}

type SheetDimensions struct {
	// EntitiesPerRow is the number of Entities each row of a Sheet contains - the number of columns,
	// where on a given row each column is an Entity
	// (there are EntitiesPerRow * ModesPerEntity Sprite columns in a Sheet).
	// There are EntitiesPerRow * EntitiesPerColumn Entities in a Sheet.
	EntitiesPerRow int
	// EntitiesPerColumn is the number of Entities each column of a Sheet contains - the number of rows,
	// where on a given column each row is an Entity
	// (there are EntitiesPerColumn * FramesPerAnimation Sprite rows in a Sheet).
	// There are EntitiesPerRow * EntitiesPerColumn Entities in a Sheet.
	EntitiesPerColumn int

	// ModesPerEntity is the number Modes (version, unique animation, etc.; e.g. directions of movement) each Entity has.
	// An Entity has one Mode per column. An Entity Mode / column may be blank / unused.
	ModesPerEntity int
	// The number of (Sprite) frames each Entity Mode animation has. An Entity has one frame (Sprite) per row.
	// A frame/row may be blank/unused (the Entity must specify the number of frames for each Mode, or it defaults to
	// FramesPerAnimation; if a frame is blank and the Entity Mode frame count includes it,
	// the blank frame will be shown / included in the animation - there is no logic to check if a frame is blank).
	FramesPerAnimation int
	// FramesRunRows controls the orientation of Modes and their frames within an Entity.
	// False (default) = Each Mode in the entity is a column, and the frames for that Mode run down the column.
	// True = Each Mode in the entity is a row, and the frames for that Mode run along the row.
	FramesRunRows                   bool
	numEntityColumns, numEntityRows int

	// An individual Sprite (frame) is SpriteWidth * SpriteHeight pixels.
	// The sheet image must be EntitiesPerRow * ModesPerEntity * SpriteWidth pixels wide and
	// EntitiesPerColumn * FramesPerAnimation * SpriteHeight pixels high.

	// SpriteWidth is the width of each Sprite (Frame) in pixels.
	SpriteWidth int
	// SpriteHeight is the height of each Sprite (Frame) in pixels.
	SpriteHeight int
}

// EntityAndModeNames contains the name for an Entity and the names for each of its Modes. It is used in the Sheet
// factories to supply names. The length of ModeNames determines how many Modes will be read/created from the Sheet,
// and it must be <= the number of Modes available for each Entity in the sheet layout (SheetDimensions.ModesPerEntity).
type EntityAndModeNames struct {
	// EntityName is the name used to identify the Entity (as a whole, including all its Modes)
	EntityName string
	// ModeNames is a slice of names for each of the Entity's Modes.
	ModeNames []string
}

// init takes the provided SheetDimensions and assigns the non-exported fields which are used during Sheet creation to
// control whether Modes are columns and Frames of that Mode rows, or vice versa, based on the supplied
// SheetDimensions.FramesRunRows field.
func (d *SheetDimensions) init() {
	if d.FramesRunRows {
		d.numEntityColumns = d.FramesPerAnimation
		d.numEntityRows = d.ModesPerEntity
	} else {
		d.numEntityColumns = d.ModesPerEntity
		d.numEntityRows = d.FramesPerAnimation
	}
}

// Sheet holds the Entities of the sheets, along with an Entity name lookup map. An Entity is a unit of Sprites (such
// as a character), and it has Modes which are different states or views (such as direction character is walking), and
// each Mode has a slice of Sprite (image) Frames comprising its animation.
// When using a Sheet (which should be created only once for a given sprite sheet image / file), one should get an
// Instance of an Entity (multiple copies of an entity may exist in an environment), which has a current Mode state
// and an underlying animation which is used by the Instance to control what its current frame is;
// said current frame may be requested directly, or Instance.PlaceSprite may be used to place the current frame on a
// provided image.
type Sheet struct {
	// entities is a map of index->GetEntity (pointer). Index is the position on the Sheet, which starts at upper-left and
	// wraps back to the left at the end of a row of Entities.
	entities map[int]*Entity
	// entityNamesToIndex is a map of Entity.name -> index, where index is a key in entities.
	entityNamesToIndex map[string]int
}

// NewSheet is a basic factory to create a new Sheet from a sprite sheet image and SheetDimensions info about how it is
// organized.
// img is the underlying image.Image which contains all the sub images / pixel data for each Sprite.
// The image must implement SubImage() and be EntitiesPerRow * ModesPerEntity * SpriteWidth pixels wide and
// EntitiesPerColumn * FramesPerAnimation * SpriteHeight pixels high.
func NewSheet(img image.Image, dimensions SheetDimensions) (*Sheet, error) {
	dimensions.init()
	spriteSheet, err := createSpriteSheet(img, dimensions)
	if err != nil {
		return nil, err
	}

	newSheet := new(Sheet)

	modeNames := generateModeNames(dimensions.ModesPerEntity)
	var names []EntityAndModeNames
	for i := 0; i < dimensions.EntitiesPerRow*dimensions.EntitiesPerColumn; i++ {
		names = append(names, EntityAndModeNames{"GetEntity" + strconv.Itoa(i), modeNames})
	}

	newSheet.generateEntities(spriteSheet, dimensions, names)

	return newSheet, nil
}

// note that len(names) defines the number of populated/used entities
//describe entity index order in docstring
func NewSheetWithEntityNames(img image.Image, dimensions SheetDimensions, entityNames []string) (*Sheet, error) {
	modeNames := generateModeNames(dimensions.ModesPerEntity)

	return NewSheetWithEntityAndSharedModeNames(img, dimensions, entityNames, modeNames)
}

// Mode names for each Entity are the same
func NewSheetWithEntityAndSharedModeNames(img image.Image, dimensions SheetDimensions, entityNames []string, modeNames []string) (*Sheet, error) {
	if len(entityNames) > dimensions.EntitiesPerRow*dimensions.EntitiesPerColumn {
		return nil, fmt.Errorf("length of entityNames (%d) is greater than number of Entities in Sheet, i.e. EntitiesPerRow * EntitiesPerColumn (%d)",
			len(entityNames), dimensions.EntitiesPerRow*dimensions.EntitiesPerColumn)
	}

	dimensions.init()
	spriteSheet, err := createSpriteSheet(img, dimensions)
	if err != nil {
		return nil, err
	}

	newSheet := new(Sheet)

	var names []EntityAndModeNames
	for _, entityName := range entityNames {
		names = append(names, EntityAndModeNames{entityName, modeNames})
	}

	newSheet.generateEntities(spriteSheet, dimensions, names)

	return newSheet, nil
}

//note that len(names) defines the number of populated/used Entities, and len of each key defines the number of populate/used modes for the given Entity
//describe entity and mode index order in docstring
func NewSheetWithNames(img image.Image, dimensions SheetDimensions, names []EntityAndModeNames) (*Sheet, error) {
	if len(names) > dimensions.EntitiesPerRow*dimensions.EntitiesPerColumn {
		return nil, fmt.Errorf("length of names (%d) is greater than number of Entities in Sheet, i.e. EntitiesPerRow * EntitiesPerColumn (%d)",
			len(names), dimensions.EntitiesPerRow*dimensions.EntitiesPerColumn)
	}

	dimensions.init()
	spriteSheet, err := createSpriteSheet(img, dimensions)
	if err != nil {
		return nil, err
	}

	newSheet := new(Sheet)

	// In this case the panic from generateEntities is due to a value in names being the wrong length, which is an issue
	// created by the caller, not this package. So we recover from it and pass it along as an error instead.
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("names has more keys (%d) than spriteSheet has Entities (%d)",
				len(names), dimensions.EntitiesPerRow*dimensions.EntitiesPerColumn)
		}
	}()
	newSheet.generateEntities(spriteSheet, dimensions, names)
	if err != nil {
		return nil, err
	}

	return newSheet, nil
}

func createSpriteSheet(img image.Image, dimensions SheetDimensions) (subImager, error) {
	spriteSheet, ok := img.(subImager)
	if ok {
		if dimensions.EntitiesPerRow <= 0 || dimensions.EntitiesPerColumn <= 0 || dimensions.ModesPerEntity <= 0 ||
			dimensions.FramesPerAnimation <= 0 || dimensions.SpriteWidth <= 0 || dimensions.SpriteHeight <= 0 {
			return nil, errors.New("all SheetDimensions fields must be > 0")
		}
		if spriteSheet.Bounds().Dx() != dimensions.EntitiesPerRow*dimensions.numEntityColumns*dimensions.SpriteWidth {
			return nil, fmt.Errorf("image width (%d) is not EntitiesPerRow * #cols/GetEntity * SpriteWidth (%d)",
				spriteSheet.Bounds().Dx(), dimensions.EntitiesPerRow*dimensions.numEntityColumns*dimensions.SpriteWidth)
		}
		if spriteSheet.Bounds().Dy() != dimensions.EntitiesPerColumn*dimensions.numEntityRows*dimensions.SpriteHeight {
			return nil, fmt.Errorf("image height (%d) is not EntitiesPerColumn * #rows/GetEntity * SpriteHeight (%d)",
				spriteSheet.Bounds().Dy(), dimensions.EntitiesPerColumn*dimensions.numEntityRows*dimensions.SpriteHeight)
		}

		return spriteSheet, nil
	} else {
		return nil, errors.New("provided image is a format which implements the SubImage and/or Set methods")
	}
}

func generateModeNames(count int) []string {
	var names []string
	for i := 0; i < count; i++ {
		names = append(names, "Mode"+strconv.Itoa(i))
	}
	return names
}

func (s *Sheet) generateEntities(spriteSheet subImager, dimensions SheetDimensions, names []EntityAndModeNames) {
	if len(names) > dimensions.EntitiesPerRow*dimensions.EntitiesPerColumn {
		panic(fmt.Errorf("internal error: names has more keys (%d) than spriteSheet has Entities (%d)",
			len(names), dimensions.EntitiesPerRow*dimensions.EntitiesPerColumn))
	}
	var x, y, dx, dy int
	var frame image.Image
	spriteSize := image.Rect(0, 0, dimensions.SpriteWidth, dimensions.SpriteHeight)
	s.entities = make(map[int]*Entity)
	s.entityNamesToIndex = make(map[string]int)
	for i, emNames := range names {
		if len(emNames.ModeNames) > dimensions.ModesPerEntity {
			panic(fmt.Errorf("names value, the slice of Mode names, has more entries (%d) than dimensions.ModesPerEntity (%d)",
				len(emNames.ModeNames), dimensions.ModesPerEntity))
		}
		x = (i % dimensions.EntitiesPerRow) * dimensions.numEntityColumns * dimensions.SpriteWidth
		y = (i / dimensions.EntitiesPerRow) * dimensions.numEntityRows * dimensions.SpriteHeight
		s.entities[i] = &Entity{
			name: emNames.EntityName,
		}
		s.entities[i].modes = make(map[int]*Mode)
		s.entities[i].modeNamesToIndex = make(map[string]int)
		for j, modeName := range emNames.ModeNames {
			s.entities[i].modes[j] = &Mode{
				name:       modeName,
				spriteSize: spriteSize,
			}
			for f := 0; f < dimensions.FramesPerAnimation; f++ {
				if dimensions.FramesRunRows {
					dx = f
					dy = j
				} else {
					dx = j
					dy = f
				}
				frame = spriteSheet.SubImage(spriteSize.Add(image.Pt(x+dx*dimensions.SpriteWidth, y+dy*dimensions.SpriteHeight)))
				s.entities[i].modes[j].frames = append(s.entities[i].modes[j].frames, frame)
			}
			s.entities[i].modeNamesToIndex[modeName] = j
		}
		s.entityNamesToIndex[emNames.EntityName] = i
	}
}

//describe index order in docstring
func (s *Sheet) GetEntityByIndex(idx int) (*Entity, error) {
	entity, ok := s.entities[idx]
	if ok {
		return entity, nil
	} else {
		return nil, fmt.Errorf("entity with index %d does not exist in Sheet", idx)
	}
}

func (s *Sheet) GetEntityByName(name string) (*Entity, error) {
	if idx, ok := s.entityNamesToIndex[name]; ok {
		if entity, ok := s.entities[idx]; ok {
			return entity, nil
		} else {
			panic(fmt.Errorf("internal error: GetEntity with index %d does not exist in Sheet; Sheet is corrupted", idx))
		}
	} else {
		return nil, fmt.Errorf("entity with name %s does not exist in Sheet", name)
	}
}

func (s *Sheet) RenameEntity(oldName, newName string) error {
	idx, ok := s.entityNamesToIndex[oldName]
	if ok {
		entity, ok := s.entities[idx]
		if ok {
			entity.name = newName
			s.entityNamesToIndex[newName] = idx
			delete(s.entityNamesToIndex, oldName)
		} else {
			panic(fmt.Errorf("internal error: GetEntity with index %d does not exist in Sheet; Sheet is corrupted", idx))
		}
	} else {
		return fmt.Errorf("entity with name %s does not exist in Sheet", oldName)
	}
	return nil
}

func (s *Sheet) EntityCount() int {
	return len(s.entities)
}

//only decrease
func (s *Sheet) SetEntityCount(count int) error {
	if count > 0 && count <= len(s.entities) {
		var delList []string
		for k, v := range s.entityNamesToIndex {
			if v >= count {
				delList = append(delList, k)
			}
		}
		for _, d := range delList {
			delete(s.entityNamesToIndex, d)
		}
		for i := count; i < len(s.entities); i++ {
			delete(s.entities, i)
		}
		return nil
	} else {
		return fmt.Errorf("new GetEntity count (%d) must be <= the current GetEntity count (%d) and > 0", count, len(s.entities))
	}
}
