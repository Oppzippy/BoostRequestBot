package roll

type WeightedRollResultsIterator[T any] struct {
	results WeightedRollResults[T]
	index   int
}

func (iter *WeightedRollResultsIterator[T]) HasNext() bool {
	return iter.index < len(iter.results.items)
}

func (iter *WeightedRollResultsIterator[T]) Next() (item T, weight float64, isChosenItem bool) {
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
