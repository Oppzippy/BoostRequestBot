package roll

import (
	"log"
	"math/rand"
)

type WeightedRoll[T any] struct {
	items       []T
	weights     []float64
	totalWeight float64
}

func NewWeightedRoll[T any](expectedSize int) *WeightedRoll[T] {
	return &WeightedRoll[T]{
		items:   make([]T, 0, expectedSize),
		weights: make([]float64, 0, expectedSize),
	}
}

func (roll *WeightedRoll[T]) AddItem(item T, weight float64) {
	roll.items = append(roll.items, item)
	roll.weights = append(roll.weights, weight)
	roll.totalWeight += weight
}

func (roll *WeightedRoll[T]) Roll() (results *WeightedRollResults[T], ok bool) {
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
		log.Printf("WeightedRoll failed to choose an item! %f/%f", chosenWeight, weightAccumulator)
		return nil, false
	}

	return &WeightedRollResults[T]{
		items:       roll.items,
		weights:     roll.weights,
		chosenIndex: chosenIndex,
		roll:        chosenWeight,
	}, true
}
