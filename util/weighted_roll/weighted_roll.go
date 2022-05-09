package weighted_roll

import (
	"math"
	"math/rand"
	"sort"
)

type WeightedRoll[T any] struct {
	items []*weightedItem[T]
}

type weightedItem[T any] struct {
	item   T
	weight float64
}

type WeightedItemWithRoll[T any] struct {
	Item   T
	Weight float64
	Roll   float64
}

func NewWeightedRoll[T any](expectedSize int) *WeightedRoll[T] {
	return &WeightedRoll[T]{
		items: make([]*weightedItem[T], 0, expectedSize),
	}
}

func (roll *WeightedRoll[T]) AddItem(item T, weight float64) {
	roll.items = append(roll.items, &weightedItem[T]{
		item:   item,
		weight: weight,
	})
}

func (roll *WeightedRoll[T]) Roll() []*WeightedItemWithRoll[T] {
	// Implemented based on the following paper:
	// Weighted Random Sampling (2005; Efraimidis, Spirakis)
	// https://web.archive.org/web/20210506225452/https://utopia.duth.gr/~pefraimi/research/data/2007EncOfAlg.pdf

	results := make([]*WeightedItemWithRoll[T], len(roll.items))
	for i, item := range roll.items {
		results[i] = &WeightedItemWithRoll[T]{
			Item:   item.item,
			Weight: item.weight,
			Roll:   math.Pow(rand.Float64(), 1/item.weight),
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Roll > results[j].Roll
	})

	return results
}
