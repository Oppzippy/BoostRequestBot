package weighted_picker

import (
	"log"
	"math/rand"
)

type WeightedPicker[T any] struct {
	items       []T
	weights     []float64
	totalWeight float64
}

func NewWeightedPicker[T any](expectedSize int) *WeightedPicker[T] {
	return &WeightedPicker[T]{
		items:   make([]T, 0, expectedSize),
		weights: make([]float64, 0, expectedSize),
	}
}

func (picker *WeightedPicker[T]) AddItem(item T, weight float64) {
	picker.items = append(picker.items, item)
	picker.weights = append(picker.weights, weight)
	picker.totalWeight += weight
}

func (picker *WeightedPicker[T]) Pick() (results *WeightedPickerResults[T], ok bool) {
	if len(picker.items) == 0 {
		return nil, false
	}
	var weightAccumulator float64
	chosenWeight := rand.Float64() * picker.totalWeight
	chosenIndex := -1
	for i := 0; i < len(picker.items); i++ {
		weightAccumulator += picker.weights[i]
		if chosenWeight < weightAccumulator {
			chosenIndex = i
			break
		}
	}

	if chosenIndex == -1 {
		log.Printf("WeightedPicker failed to choose an item! %f/%f", chosenWeight, weightAccumulator)
		return nil, false
	}

	return &WeightedPickerResults[T]{
		items:        picker.items,
		weights:      picker.weights,
		chosenIndex:  chosenIndex,
		chosenNumber: chosenWeight,
	}, true
}
