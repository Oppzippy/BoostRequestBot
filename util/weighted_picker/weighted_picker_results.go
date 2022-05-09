package weighted_picker

type WeightedPickerResults[T any] struct {
	items        []T
	weights      []float64
	chosenIndex  int
	chosenNumber float64
}

func (results *WeightedPickerResults[T]) Iterator() *WeightedPickerResultsIterator[T] {
	return &WeightedPickerResultsIterator[T]{
		results: *results,
	}
}

func (results *WeightedPickerResults[T]) HasChosenItem() bool {
	return results.chosenIndex != -1
}

func (results *WeightedPickerResults[T]) ChosenItem() (item T) {
	if !results.HasChosenItem() {
		var zero T
		return zero
	}
	return results.items[results.chosenIndex]
}

func (results *WeightedPickerResults[T]) ChosenItemAndWeight() (item T, weight float64) {
	if !results.HasChosenItem() {
		var zero T
		return zero, 0
	}
	item = results.items[results.chosenIndex]
	weight = results.weights[results.chosenIndex]
	return item, weight
}

// ChosenNumber Returns the random number used to determine the chosen item
func (results *WeightedPickerResults[T]) ChosenNumber() float64 {
	return results.chosenNumber
}
