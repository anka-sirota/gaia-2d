package controls

import (
	"fmt"
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"gogame/assets"
	"gogame/messages"
	"log"
	"time"
)

type controlEntity struct {
	*ecs.BasicEntity
	*common.MouseComponent
	*common.SpaceComponent
	*common.RenderComponent
}

type MouseTracker struct {
	*ecs.BasicEntity
	*common.MouseComponent
}

type ControlsSystem struct {
	world         *ecs.World
	entities      []*controlEntity
	hoveredEntity *controlEntity
	*MouseTracker
	paused bool
}

func (self *ControlsSystem) New(w *ecs.World) {
	entity := ecs.NewBasic()
	self.MouseTracker = &MouseTracker{&entity, &common.MouseComponent{Track: true}}
	self.world = w
	self.paused = false

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.MouseSystem:
			sys.Add(self.MouseTracker.BasicEntity, self.MouseTracker.MouseComponent, nil, nil)
		}
	}
	engo.Mailbox.Listen(messages.InteractionMessageType, self.HandleInteractMessage)
}

func (self *ControlsSystem) Add(basic *ecs.BasicEntity, mouse *common.MouseComponent, space *common.SpaceComponent, render *common.RenderComponent) {
	entity := &controlEntity{basic, mouse, space, render}
	self.entities = append(self.entities, entity)

	for _, system := range self.world.Systems() {
		switch sys := system.(type) {
		case *common.MouseSystem:
			sys.Add(entity.BasicEntity, entity.MouseComponent, entity.SpaceComponent, entity.RenderComponent)
		}
	}
}

func (self *ControlsSystem) Update(dt float32) {
	if engo.Input.Button("ExitToDesktop").JustPressed() {
		engo.Mailbox.Dispatch(messages.ControlMessage{
			Action: "exit",
		})
	}
	if engo.Input.Button("AddCreature").JustPressed() {
		engo.Mailbox.Dispatch(messages.ControlMessage{
			Action:     "add_creature",
			CreatureID: 1,
		})
	}
	if engo.Input.Button("TogglePause").JustPressed() {
		if !self.paused {
			engo.Time.Pause()
			self.paused = true
			engo.Mailbox.Dispatch(messages.HUDTextUpdateMessage{
				Name: "EventMessage",
				GetText: func() string {
					return "Paused"
				},
			})
		} else {
			engo.Time.Unpause()
			self.paused = false

			engo.Mailbox.Dispatch(messages.HUDTextUpdateMessage{
				Name:      "EventMessage",
				HideAfter: 3 * time.Second,
				GetText: func() string {
					return "Resumed"
				},
			})
		}
		engo.Mailbox.Dispatch(messages.ControlMessage{
			Action: "TogglePause",
		})

	}
	if engo.Input.Button("QuickSave").JustPressed() {
		engo.Mailbox.Dispatch(messages.SaveMessage{
			Filepath: assets.WorkDir + "/quick.save",
		})
	}
	if engo.Input.Button("NewWorld").JustPressed() {
		engo.Mailbox.Dispatch(messages.ControlMessage{
			Action: "ReloadWorld",
			Data:   "",
		})
	}
	if engo.Input.Button("QuickLoad").JustPressed() {
		engo.Mailbox.Dispatch(messages.ControlMessage{
			Action: "ReloadWorld",
			Data:   assets.WorkDir + "/quick.save",
		})
	}
	if engo.Input.Button("AddObject").JustPressed() {
		engo.Mailbox.Dispatch(messages.ControlMessage{
			Action:   "add_object",
			ObjectID: 8,
		})
	}

	var newHoveredEntity *controlEntity
	for _, entity := range self.entities {
		if entity.MouseComponent.Hovered || entity.MouseComponent.Enter {
			newHoveredEntity = entity
		}
		if entity.MouseComponent.Leave && self.hoveredEntity != nil && self.hoveredEntity.ID() == entity.BasicEntity.ID() {
			self.hoveredEntity = nil
		}
	}
	if newHoveredEntity != nil && self.hoveredEntity == nil {
		log.Printf("Hovering over an entity: %+v #%d\n", newHoveredEntity, newHoveredEntity.ID())
		engo.Mailbox.Dispatch(messages.InteractionMessage{
			Action:      "mouse_hover",
			BasicEntity: newHoveredEntity.BasicEntity,
		})
		engo.SetCursor(engo.CursorHand)
	}
	self.hoveredEntity = newHoveredEntity
}

func (self *ControlsSystem) Remove(e ecs.BasicEntity) {
	delete := -1
	for index, entity := range self.entities {
		if entity.BasicEntity.ID() == e.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		self.entities = append(self.entities[:delete], self.entities[delete+1:]...)
	}
}

func (self *ControlsSystem) GetEntityByID(basicEntityID uint64) *controlEntity {
	for _, e := range self.entities {
		if e.BasicEntity.ID() == basicEntityID {
			return e
		}
	}
	return nil
}

func (self *controlEntity) GetTextStatus() string {
	return fmt.Sprintf("#%d\n%v\n<ControlsSystem>", self.BasicEntity.ID(), self)
}

func (self *ControlsSystem) HandleInteractMessage(m engo.Message) {
	log.Printf("ControlsSystem: %+v", m)
	msg, ok := m.(messages.InteractionMessage)
	if !ok {
		return
	}
	if msg.Action == "mouse_hover" && msg.BasicEntity != nil {
		entity := self.GetEntityByID(msg.BasicEntity.ID())
		log.Printf("%+v", entity)
		/*
			if entity != nil {
				engo.Mailbox.Dispatch(messages.HUDTextUpdateMessage{
					Name:    "HoverInfo",
					GetText: entity.GetTextStatus,
				})
			}
		*/
	}
}
