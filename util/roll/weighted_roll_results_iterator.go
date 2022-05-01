package roll

type WeightedRollResultsIterator struct {
	results WeightedRollResults
	index   int
}

func (iter *WeightedRollResultsIterator) HasNext() bool {
	return iter.index < len(iter.results.items)
}

func (iter *WeightedRollResultsIterator) Next() (item string, weight float64, isChosenItem bool) {
	if iter.index == len(iter.results.items) {
		return "", 0, false
	}
	item = iter.results.items[iter.index]
	weight = iter.results.weights[iter.index]
	isChosenItem = iter.results.chosenIndex == iter.index
	iter.index++
	return item, weight, isChosenItem
}
