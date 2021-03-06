package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"gogame/calendar"
	"gogame/messages"
	"log"
)

type TimeSystem struct {
	world         *ecs.World
	dtFullSeconds float32

	speed         float32
	previousSpeed float32

	Time  *calendar.Time
	Month calendar.Month
}

func (self *TimeSystem) New(w *ecs.World) {
	self.world = w
	self.Time = &calendar.Time{}
	self.speed = 1.0
	self.previousSpeed = 1.0
	engo.Mailbox.Dispatch(messages.TimeSecondPassedMessage{
		Time: self.Time,
	})

	engo.Mailbox.Listen(messages.ControlMessageType, self.HandleControlMessage)
}

func (*TimeSystem) Add() {}

func (self *TimeSystem) Update(dt float32) {
	if self.speed == 0 {
		// Paused
		return
	}
	if self.dtFullSeconds*self.speed >= 1 {
		self.dtFullSeconds = 0
		self.Time.AddSecond()
		engo.Mailbox.Dispatch(messages.TimeSecondPassedMessage{
			Time: self.Time,
			Dt:   self.dtFullSeconds,
		})
	}
	self.dtFullSeconds += dt
	// TODO might be the good place to implement speed of in-game time
	// TODO and the PauseSystem
}

func (self *TimeSystem) HandleControlMessage(m engo.Message) {
	msg, ok := m.(messages.ControlMessage)
	if !ok {
		return
	}
	log.Printf("[TimeSystem] %+v", m)
	switch msg.Action {
	}
}

func (*TimeSystem) Remove(ecs.BasicEntity) {}
