package roll

type WeightedRollResults[T any] struct {
	items       []T
	weights     []float64
	chosenIndex int
	roll        float64
}

func (results *WeightedRollResults[T]) Iterator() *WeightedRollResultsIterator[T] {
	return &WeightedRollResultsIterator[T]{
		results: *results,
	}
}

func (results *WeightedRollResults[T]) HasChosenItem() bool {
	return results.chosenIndex != -1
}

func (results *WeightedRollResults[T]) ChosenItem() (item T) {
	if !results.HasChosenItem() {
		var zero T
		return zero
	}
	return results.items[results.chosenIndex]
}

func (results *WeightedRollResults[T]) ChosenItemAndWeight() (item T, weight float64) {
	if !results.HasChosenItem() {
		var zero T
		return zero, 0
	}
	item = results.items[results.chosenIndex]
	weight = results.weights[results.chosenIndex]
	return item, weight
}

// Roll Returns the random number used to determine the chosen item
func (results *WeightedRollResults[T]) Roll() float64 {
	return results.roll
}
