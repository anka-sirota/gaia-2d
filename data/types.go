package data

import (
	"fmt"
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"gogame/config"
)

type Spritesheet struct {
	ID         int                 `json:"id"`
	FilePath   string              `json:"filepath"`
	Animations []*common.Animation `json:"animations"`
}

type Spritesheets struct {
	Spritesheets []*Spritesheet              `json:"spritesheets"`
	Loaded       map[int]*common.Spritesheet `json:"-"`
}

// Description of a resource
type Resource struct {
	ID   int    `json:"id"`
	Type string `json:"type"`
}

type Resources struct {
	Resources []*Resource `json:"resources"`
}

// Mutable resource, e.g. plant matter available at a specific tile
type AccessibleResource struct {
	ResourceID int     `json:"resource_id"`
	Amount     float32 `json:"amount"`
}

type Object struct {
	ID            int     `json:"id"`
	SpriteID      int     `json:"sprite_id"`
	SpritesheetID int     `json:"spritesheet_id"`
	Name          string  `json:"name"`
	ResourceID    int     `json:"resource_id"`
	Amount        float32 `json:"amount"`

	// Runtime only fields
	Spritesheet *common.Spritesheet `json:"-"`
	Animations  []*common.Animation `json:"-"`
}

type Objects struct {
	Objects []*Object `json:"objects"`
}

type Tile struct {
	*ecs.BasicEntity           `json:"-"` // FIXME? marshalled into an empty object
	*common.RenderComponent    `json:"-"` // FIXME cannot unmarshal .. color.Color
	*common.AnimationComponent `json:"-"`

	SpaceComponent     *common.SpaceComponent
	CollisionComponent *common.CollisionComponent
	MouseComponent     *common.MouseComponent

	Layer              float32
	ObjectID           int
	AccessibleResource *AccessibleResource
	Object             *Object   `json:"-"`
	Resource           *Resource `json:"-"`
}

func (self *Tile) AABB() engo.AABB {
	return self.SpaceComponent.AABB()
}

func (self *Tile) SurroundingAreaAABB(radius float32) engo.AABB {
	var padding float32 = 1 // extra pixel in order to get corner positions into the rectangle
	return engo.AABB{
		Min: engo.Point{
			X: self.SpaceComponent.Position.X - float32(config.SpriteWidth) * radius - padding,
			Y: self.SpaceComponent.Position.Y - float32(config.SpriteHeight) * radius - padding,
		},
		Max: engo.Point{
			X: self.SpaceComponent.Position.X + float32(config.SpriteWidth) * (radius + 1) + padding,
			Y: self.SpaceComponent.Position.Y + float32(config.SpriteHeight) * (radius + 1) + padding,
		},
	}
}

func (self *Tile) GetTextStatus() []string {
	return []string{
		fmt.Sprintf("#%d", self.BasicEntity.ID()),
		fmt.Sprintf("%v", self.Object),
		fmt.Sprintf("%v", self.AccessibleResource),
		fmt.Sprintf("%v", self.Resource),
	}
}
