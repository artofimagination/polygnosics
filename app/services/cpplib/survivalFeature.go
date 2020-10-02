package cpplib

// #include "evoGame/GolangInterface.h"
// #cgo CFLAGS: -I/go/src/aiplayground/app/services/cpplib/evoGame
// #cgo LDFLAGS: -L/go/src/aiplayground/app/services/cpplib/evoGame -lstdc++ -lGolangInterface -lUtilsLib -lEvoGame
import "C"
import (
	"math"
	"math/rand"
	"strconv"
)

// EntityData is the golang representation of the C++ iCircleEntity class.
type EntityData struct {
	ID          uint32  `json:"id"`
	Type        uint16  `json:"type"`
	Size        float64 `json:"size"`
	PosX        int32   `json:"posx"`
	PosY        int32   `json:"posy"`
	ThrustR     float64 `json:"thrustr"`
	ThrustTheta float64 `json:"thrusttheta"`
}

type WorldConfig struct {
	FoodCount          uint32
	FoodProductionRate uint16
	CreatureCount      uint32
}

var timerData = make(map[uint32]int32)

func UpdateThrust(id uint32) {
	if _, ok := timerData[id]; ok == false {
		timerData[id] = rand.Int31() >> 24
	}

	if timerData[id] <= 0 {
		r := rand.Int31() >> 24
		theta := float64(rand.Int31()>>22) * math.Pi / 360.0
		timerData[id] = rand.Int31() >> 24
		C.SetThrust(C.int(id), C.double(r), C.double(theta))
	}
	timerData[id] = timerData[id] - 1
}

func Generate(config map[string]interface{}) {
	var configC C.WorldConfig
	dataStr := config["config"].(map[string]interface{})["Food Count"].(string)
	foodCount, _ := strconv.Atoi(dataStr)
	dataStr = config["config"].(map[string]interface{})["Creature Count"].(string)
	creatureCount, _ := strconv.Atoi(dataStr)
	configC.foodCount = C.ulong(uint32(foodCount))
	configC.foodProductionRate = 0
	configC.creatureCount = C.ulong(uint32(creatureCount))
	C.Generate(configC)
}

func GetAttributes(index int) EntityData {
	attributeC := C.GetPublicAttribute(C.int(index))
	var entityData EntityData
	entityData.ID = uint32(attributeC.id)
	entityData.Type = uint16(attributeC.entityType)
	entityData.Size = float64(attributeC.size)
	entityData.PosX = int32(attributeC.posX)
	entityData.PosY = int32(attributeC.posY)
	entityData.ThrustR = float64(attributeC.thrustR)
	entityData.ThrustTheta = float64(attributeC.thrustTheta)
	return entityData
}

func GetEntityCount() int {
	return int(C.GetEntityCount())
}

func Execute() {
	C.Execute()
}
