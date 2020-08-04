package plants

import (
	"encoding/json"
	"fmt"
	"github.com/EngoEngine/engo"
	"github.com/ulule/deepcopier"
	"gogame/assets"
	"gogame/calendar"
	"gogame/data"
	"gogame/messages"
)

var (
	plants    *Plants
	PlantById map[int]*Plant
)

type Activity uint8

const (
	Growing Activity = iota
	Resting
	Dead
)

func (a Activity) String() string {
	return [...]string{"growing", "resting", "dead"}[a]
}

type Plant struct {
	*data.Tile `deepcopier:"skip"`

	// Species properties, immutable
	ID          int     `json:"id"`
	ObjectID    int     `json:"object_id"`
	Species     string  `json:"species"`
	Name        string  `json:"name"`
	GrownID     int     `json:"grown_id"`
	GrowthRate  float32 `json:"growth_rate"`
	GrowthSpeed float32 `json:"growth_speed"`
	MaxGrowth   float32 `json:"max_growth"`

	// Live properties, mutable
	IsAlive  bool     `json:"is_alive"`
	Activity Activity `json:"activity"`
	Growth   float32  `json:"growth"`
}

type Plants struct {
	Plants []*Plant `json:"plants"`
}

func initPlants() {
	// Load plants
	byteValue := assets.ReadJSON("assets/meta/plants.json")
	json.Unmarshal(byteValue, &plants)

	PlantById = make(map[int]*Plant)
	for _, c := range plants.Plants {
		PlantById[c.ID] = c
	}
}

func GetPlantByID(plantID int) *Plant {
	if PlantById == nil {
		initPlants()
	}

	plant, ok := PlantById[plantID]
	if ok {
		return plant
	}
	panic(fmt.Sprintf("Unknown plant '%d'", plantID))
}

func (self *Plant) GetGrowthSpeed() float32 {
	// TODO affected by the environment, weather etc.
	return self.GrowthSpeed
}

func (self *Plant) GetGrowthRate() float32 {
	// TODO affected by the environment, weather etc.
	return self.GrowthRate
}

func (self *Plant) IsFullyGrown() bool {
	return self.Growth >= self.MaxGrowth
}

func (self *Plant) Mature() {
	// Replace with the mature plant or a new growth stage
	newPlant := GetPlantByID(self.GrownID)
	oldGrowth := self.Growth
	deepcopier.Copy(newPlant).To(self)
	self.Growth = oldGrowth
	self.MaxGrowth += oldGrowth

	// Update plant's visual representation
	engo.Mailbox.Dispatch(messages.TileReplaceMessage{
		Entity:   self.Tile.BasicEntity,
		ObjectID: newPlant.ObjectID,
	})
}

func (self *Plant) CurrentGrowth() string {
	return fmt.Sprintf(
		"Growth: %d/%d",
		int(self.Growth), int(self.MaxGrowth),
	)
}

func (self *Plant) CurrentPosition() string {
	p := self.Tile.SpaceComponent.Position
	return fmt.Sprintf("At (%d, %d)", int(p.X), int(p.Y))
}

func (self *Plant) GetTextStatus() string {
	return fmt.Sprintf("%s, %s\n%s\n%s\n%s", self.Name, self.Species,
		self.CurrentGrowth(),
		self.Activity,
		self.CurrentPosition(),
	)
}

func (self *Plant) Update(currentTime *calendar.Time) {
	if !self.IsAlive {
		self.Activity = Dead
	}
	if self.IsFullyGrown() {
		if self.GrownID != 0 {
			self.Mature()
		} else {
			// TODO should eventually die
		}
	} else if self.Activity == Growing {
		// Handle growth
		self.Growth += self.GetGrowthSpeed()
		self.Tile.AccessibleResource.Amount += self.GetGrowthRate()
	}
	// Handle rest TODO
}
