package assets

import (
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"

	//"fmt"
	"encoding/json"
	//"errors"
	"io/ioutil"
	"math/rand"
	"os"
)

// Description of a resource
type Resource struct {
	ID int              `json:"id"`
	Type string         `json:"type"`
}

// Mutable resource, e.g. plant matter available at a specific tile
type AccessibleResource struct {
	ResourceID int      `json:"resource_id"`
	Amount float32      `json:"amount"`
}

type Object struct {
	ID int              `json:"id"`
	SpriteID int        `json:"sprite_id"`
	Name string         `json:"name"`
	ResourceID int      `json:"resource_id"`
	Amount float32      `json:"amount"`
}

type Objects struct {
    Objects []*Object   `json:"objects"`
}

type Resources struct {
    Resources []*Resource `json:"resource"`
}

// Structures for the save file
type SavedTile struct {
	ObjectID int                            `json:"object_id"`
	AccessibleResource *AccessibleResource  `json:"accessible_resource"`
	Position           *engo.Point          `json:"position"`
	Layer              float32              `json:"layer"`
}

type SavedTiles struct {
	Tiles []*SavedTile `json:"tiles"`
}

var (
	LineHeight = 20
	FontSize = 16
	SpriteWidth = 32
	SpriteHeight = 32
	PreloadList = []string{
		"textures/chick_32x32.png",
		"tilemap/terrain-v7.png",
	}
	FullSpriteSheet *common.Spritesheet
	objects *Objects
	resources *Resources

	ResourceById map[int]*Resource
	ResourceByType map[string]*Resource

	ObjectById map[int]*Object
)

func InitAssets() {
	// Load the spritesheet
	FullSpriteSheet = common.NewSpritesheetFromFile("tilemap/terrain-v7.png", 32, 32)

	// Load objects
	objectsJsonFile, _ := os.Open("assets/meta/objects.json")
	defer objectsJsonFile.Close()
	byteValue, _ := ioutil.ReadAll(objectsJsonFile)
	json.Unmarshal(byteValue, &objects)

	// Load other related metadata
	metaJsonFile, _ := os.Open("assets/meta/resource.json")
	defer metaJsonFile.Close()
	byteValue, _ = ioutil.ReadAll(metaJsonFile)
	json.Unmarshal(byteValue, &resources)

	// Prepare hashes for ease of access to loaded assets
	ResourceById = make(map[int]*Resource)
	for _, r := range resources.Resources {
		ResourceById[r.ID] = r
	}
	ResourceByType = make(map[string]*Resource)
	for _, r := range resources.Resources {
		ResourceByType[r.Type] = r
	}
	ObjectById = make(map[int]*Object)
	for _, o := range objects.Objects {
		ObjectById[o.ID] = o
	}
}

func GetResourceByID(resourceID int) *Resource {
	if resource, ok := ResourceById[resourceID]; ok {
		return resource
	}
	return nil
}

func GetResourceByType(resourceType string) *Resource {
	if resource, ok := ResourceByType[resourceType]; ok {
		return resource
	}
	return nil
}

func GetObjectById(objectID int) *Object {
	if object, ok := ObjectById[objectID]; ok {
		return object
	}
	return nil
}

func GetObjectsByType(resourceType string) []*Object {
	var result []*Object
	for _, v := range objects.Objects {
		if GetResourceByID(v.ResourceID).Type == resourceType {
			result = append(result, v)
		}
	}
	return result
}

func GetRandomObjectOfType(resourceType string) *Object {
	objects := GetObjectsByType(resourceType)
	return objects[rand.Intn(len(objects))]
}
