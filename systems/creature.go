package systems

import (
	"fmt"
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/ulule/deepcopier"
	"gogame/assets"
	"gogame/data"
	"gogame/messages"
	"gogame/util"
	"log"
)

type CreatureMouseTracker struct {
	ecs.BasicEntity
	common.MouseComponent
}

type CreatureSpawningSystem struct {
	world         *ecs.World
	mouseTracker  CreatureMouseTracker
	dtFullSeconds float32
	entities      []*data.Creature
}

func NewCreature(creatureID int, position *engo.Point) *data.Creature {
	creature := assets.GetCreatureById(creatureID)
	tile := NewTile(creature.ObjectID, position, 4, &common.CollisionComponent{Main: 1, Group: 0})
	entity := &data.Creature{ID: creatureID, Tile: tile}
	// Initialise creature's stats from its initial record
	deepcopier.Copy(creature).To(entity)
	return entity
}

func (self *CreatureSpawningSystem) Add(entity *data.Creature) {
	self.entities = append(self.entities, entity)

	// Add the entity to the various systems
	for _, system := range self.world.Systems() {
		switch sys := system.(type) {
		case *WorldTilesSystem:
			sys.Add(entity.Tile)
		}
	}
}

// New is the initialisation of the System
func (self *CreatureSpawningSystem) New(w *ecs.World) {
	log.Println("CreatureSpawningSystem was added to the Scene")

	self.world = w
	self.mouseTracker.BasicEntity = ecs.NewBasic()
	self.mouseTracker.MouseComponent = common.MouseComponent{Track: true}

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.MouseSystem:
			sys.Add(&self.mouseTracker.BasicEntity, &self.mouseTracker.MouseComponent, nil, nil)
		}
	}

	engo.Mailbox.Listen("CollisionMessage", self.HandleCollisionMessage)
	engo.Mailbox.Listen(messages.ControlMessageType, self.HandleControlMessage)
	engo.Mailbox.Listen(messages.SpacialResponseMessageType, self.HandleSpacialResponseMessage)
	engo.Mailbox.Listen(messages.CreatureHoveredMessageType, self.HandleCreatureHoveredMessage)
}

// Update is ran every frame, with `dt` being the time
// in seconds since the last frame
func (self *CreatureSpawningSystem) Update(dt float32) {
	for _, entity := range self.entities {
		if self.dtFullSeconds > 1 {
			entity.Update(self.dtFullSeconds)
			self.dtFullSeconds = 0
		}
		// TODO might be the good place to implement speed of in-game time
		// TODO and the PauseSystem
		self.dtFullSeconds += dt
	}
}

// Remove is called whenever an Creature is removed from the World, in order to remove it from this sytem as well
func (self *CreatureSpawningSystem) Remove(e ecs.BasicEntity) {
	delete := -1
	for index, entity := range self.entities {
		if entity.BasicEntity.ID() == e.ID() {
			delete = index
		}
	}
	if delete >= 0 {
		self.entities = append(self.entities[:delete], self.entities[delete+1:]...)
	}
}

func (self *CreatureSpawningSystem) Get(entityID uint64) *data.Creature {
	for _, e := range self.entities {
		if e.BasicEntity.ID() == entityID {
			return e
		}
	}
	return nil
}

func (self *CreatureSpawningSystem) HandleSpacialResponseMessage(m engo.Message) {
	log.Printf("SpacialResponseMessage %+v", m)
	msg, ok := m.(messages.SpacialResponseMessage)
	if !ok {
		return
	}
	entity := self.Get(msg.EntityID)
	if entity == nil {
		log.Println(
			fmt.Sprintf("[SpacialResponseMessage] Got a message for an unknown entity %+v", msg))
	}
	if len(msg.Result) > 0 {
		for _, v := range msg.Result {
			if tile, ok := v.(*data.Tile); ok {
				log.Println(
					"Found tile", tile, tile.SpaceComponent.Position,
					"our position", entity.Tile.SpaceComponent.Position)
				entity.MovementTarget = tile
				break
			}
		}
	}
}

func (self *CreatureSpawningSystem) HandleCollisionMessage(message engo.Message) {
	_, isCollision := message.(common.CollisionMessage)

	if isCollision {
		//log.Println("COLLISION")
	}
}

func (self *CreatureSpawningSystem) HandleControlMessage(m engo.Message) {
	log.Printf("%+v", m)
	msg, ok := m.(messages.ControlMessage)
	if !ok {
		return
	}
	if msg.Action == "add_creature" {
		x, y := util.ToGridPosition(self.mouseTracker.MouseX, self.mouseTracker.MouseY)
		c := NewCreature(msg.CreatureID, &engo.Point{x, y})
		self.Add(c)

	}
}

func (self *CreatureSpawningSystem) HandleCreatureHoveredMessage(m engo.Message) {
	msg, ok := m.(messages.CreatureHoveredMessage)
	if !ok {
		return
	}
	entity := self.Get(msg.EntityID)
	lines := []string{
		fmt.Sprintf("%s, %s", entity.Name, entity.Subspecies),
		fmt.Sprintf("%v", entity.Needs),
		fmt.Sprintf("%s", entity.Activity),
		fmt.Sprintf("%v", entity),
	}
	engo.Mailbox.Dispatch(messages.HUDTextUpdateMessage{
		Name:  "HoverInfo",
		Lines: lines,
	})
}

func (self *CreatureSpawningSystem) UpdateSave(saveFile *data.SaveFile) {
	for _, e := range self.entities {
		if e.Tile.Object.Type == "creature" {
			saveFile.Creatures = append(saveFile.Creatures, e)
		}
	}
}

func (self *CreatureSpawningSystem) LoadSave(saveFile *data.SaveFile) {
	log.Printf("[CreatureSpawningSystem] Creatures in the save file: %d\n", len(saveFile.Creatures))
	for _, c := range saveFile.Creatures {
		self.Add(c)
	}
}
