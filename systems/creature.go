package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"

	//"log"
	"fmt"
)

type CreatureMouseTracker struct {
	ecs.BasicEntity
	common.MouseComponent
}

type Creature struct {
	ecs.BasicEntity
	common.AnimationComponent
	common.RenderComponent
	common.SpaceComponent
}

// System is an interface which implements an ECS-System. A System
// should iterate over its Entities on `Update`, in any way
// suitable for the current implementation.
type System interface {
	// Update is ran every frame, with `dt` being the time
	// in seconds since the last frame
	Update(dt float32)

	// Remove removes an Creature from the System
	Remove(ecs.BasicEntity)
}

type CreatureSpawningSystem struct {
	world        *ecs.World
	mouseTracker CreatureMouseTracker
	entityActions []*common.Animation
}

func (self *CreatureSpawningSystem) CreateCreature(point engo.Point, spriteSheet *common.Spritesheet) *Creature {
	entity := &Creature{BasicEntity: ecs.NewBasic()}

	entity.SpaceComponent = common.SpaceComponent{
		Position: point,
		Width:    32,
		Height:   32,
	}
	entity.RenderComponent = common.RenderComponent{
		Drawable: spriteSheet.Cell(0),
		Scale:    engo.Point{1, 1},
	}
	entity.RenderComponent.SetZIndex(10.0)
	entity.AnimationComponent = common.NewAnimationComponent(spriteSheet.Drawables(), 0.25)

	entity.AnimationComponent.AddAnimations(self.entityActions)
	entity.AnimationComponent.AddDefaultAnimation(self.entityActions[0])

	return entity
}

// New is the initialisation of the System
func (self *CreatureSpawningSystem) New(w *ecs.World) {
	fmt.Println("CreatureSpawningSystem was added to the Scene")

	self.world = w

	self.mouseTracker.BasicEntity = ecs.NewBasic()
	self.mouseTracker.MouseComponent = common.MouseComponent{Track: true}

	self.entityActions = []*common.Animation{
		&common.Animation{Name: "feed", Frames: []int{4, 4, 5, 6, 7, 8, 9, 9}},
		&common.Animation{Name: "walk_right", Frames: []int{4, 4, 0, 1, 2, 3, 4, 4, 4}, Loop: true},
		&common.Animation{Name: "walk_left", Frames: []int{9, 9, 10, 11, 12, 13, 9, 9}, Loop: true},
		&common.Animation{Name: "walk_down", Frames: []int{14, 15, 16, 17, 18, 19, 20}, Loop: true},
		&common.Animation{Name: "walk_up", Frames: []int{21, 22, 23, 24, 25, 26, 27}, Loop: true},
	}

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.MouseSystem:
			sys.Add(&self.mouseTracker.BasicEntity, &self.mouseTracker.MouseComponent, nil, nil)
		}
	}
}

// Update is ran every frame, with `dt` being the time
// in seconds since the last frame
func (self *CreatureSpawningSystem) Update(dt float32) {
	if engo.Input.Button("AddCreature").JustPressed() {
		fmt.Println("The gamer pressed F1")
		spriteSheet := common.NewSpritesheetFromFile("textures/chick_32x32.png", 32, 32)
		animal := self.CreateCreature(engo.Point{self.mouseTracker.MouseX, self.mouseTracker.MouseY}, spriteSheet)

		for _, system := range self.world.Systems() {
			switch sys := system.(type) {
			case *common.RenderSystem:
				sys.Add(&animal.BasicEntity, &animal.RenderComponent, &animal.SpaceComponent)
			case *common.AnimationSystem:
				sys.Add(&animal.BasicEntity, &animal.AnimationComponent, &animal.RenderComponent)
			}
		}
	}
}

// Remove is called whenever an Creature is removed from the World, in order to remove it from this sytem as well
func (*CreatureSpawningSystem) Remove(ecs.BasicEntity) {}
