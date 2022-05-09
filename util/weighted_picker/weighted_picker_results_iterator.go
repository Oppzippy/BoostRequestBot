package weighted_picker

type WeightedPickerResultsIterator[T any] struct {
	results WeightedPickerResults[T]
	index   int
}

func (iter *WeightedPickerResultsIterator[T]) HasNext() bool {
	return iter.index < len(iter.results.items)
}

func (iter *WeightedPickerResultsIterator[T]) Next() (item T, weight float64, isChosenItem bool) {
	if iter.index == len(iter.results.items) {
		var zero T
		return zero, 0, false
	}
	item = iter.results.items[iter.index]
	weight = iter.results.weights[iter.index]
	isChosenItem = iter.results.chosenIndex == iter.index
	iter.index++
	return item, weight, isChosenItem
}
