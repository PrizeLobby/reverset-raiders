package core

import (
	"fmt"
	"math/rand"
)

type TargetType int

const (
	TILE TargetType = iota
	NEARBY
	ALL
)

type Effect struct {
	X                int
	Y                int
	Targets          TargetType
	SpeciesCondition Species
	ColorCondition   CreatureColor
	Value            int
}

func (e Effect) String() string {
	target := ""
	if e.Targets == ALL {
		target = "all"
	} else if e.Targets == NEARBY {
		target = "nearby"
	}

	if e.SpeciesCondition == NO_SPECIES && e.ColorCondition == NO_COLOR {
		return fmt.Sprintf("%s +%d", target, e.Value)
	}
	if e.SpeciesCondition != NO_SPECIES {
		return fmt.Sprintf("%s %s +%d", target, e.SpeciesCondition, e.Value)
	}
	if e.ColorCondition != NO_COLOR {
		return fmt.Sprintf("%s %s +%d", target, e.ColorCondition, e.Value)
	}

	return ""
}

var weightedValues []int = []int{3, 4, 4, 4, 4, 5, 5, 5, 5, 5, 5, 6, 6, 6, 6, 7}
var weightedColors []CreatureColor = []CreatureColor{NO_COLOR, Red, Red, Red, Green, Green, Green, Blue, Blue, Blue}
var weightedSpecies []Species = []Species{NO_SPECIES, Duck, Duck, Duck, Tortoise, Tortoise, Tortoise, Capybara, Capybara, Capybara}
var weightedTargetType []TargetType = []TargetType{ALL, TILE, TILE, TILE, TILE, TILE, TILE, TILE, TILE, TILE}

func RandomEffect(x, y int, random *rand.Rand) Effect {
	conditionType := random.Intn(2)
	sCond := NO_SPECIES
	cCond := NO_COLOR
	if conditionType == 0 {
		sCond = weightedSpecies[random.Intn(len(weightedSpecies))]
	} else if conditionType == 1 {
		cCond = weightedColors[random.Intn(len(weightedColors))]
	}

	targetType := weightedTargetType[random.Intn(len(weightedTargetType))]
	value := weightedValues[random.Intn(len(weightedValues))]

	if sCond == NO_SPECIES && cCond == NO_COLOR {
		value -= 1
	}
	if targetType == NEARBY {
		value -= 1
	} else if targetType == ALL {
		value -= 2
	}
	if value < 0 {
		value = 0
	}

	return Effect{
		X:                x,
		Y:                y,
		Targets:          targetType,
		SpeciesCondition: sCond,
		ColorCondition:   cCond,
		Value:            value,
	}
}
