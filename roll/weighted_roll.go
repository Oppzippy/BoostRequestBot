package roll

import (
	"fmt"
	"math/rand"
)

type WeightedRoll struct {
	items       []string
	weights     []float64
	totalWeight float64
}

func NewWeightedRoll(expectedSize int) *WeightedRoll {
	return &WeightedRoll{
		items:   make([]string, 0, expectedSize),
		weights: make([]float64, 0, expectedSize),
	}
}

func (roll *WeightedRoll) AddItem(item string, weight float64) {
	roll.items = append(roll.items, item)
	roll.weights = append(roll.weights, weight)
	roll.totalWeight += weight
}

func (roll *WeightedRoll) Roll() (results *WeightedRollResults, ok bool) {
	if len(roll.items) == 0 {
		return nil, false
	}
	var weightAccumulator float64
	chosenWeight := rand.Float64() * roll.totalWeight
	chosenIndex := -1
	for i := 0; i < len(roll.items); i++ {
		weightAccumulator += roll.weights[i]
		if chosenWeight < weightAccumulator {
			chosenIndex = i
			break
		}
	}

	if chosenIndex == -1 {
		fmt.Printf("WeightedRoll failed to choose an item! %f/%f", chosenWeight, weightAccumulator)
		return nil, false
	}

	return &WeightedRollResults{
		items:       roll.items,
		weights:     roll.weights,
		chosenIndex: chosenIndex,
		roll:        chosenWeight,
	}, true
}
