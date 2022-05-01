package roll

type WeightedRollResults struct {
	items       []string
	weights     []float64
	chosenIndex int
	roll        float64
}

func (results *WeightedRollResults) Iterator() *WeightedRollResultsIterator {
	return &WeightedRollResultsIterator{
		results: *results,
	}
}

func (results *WeightedRollResults) HasChosenItem() bool {
	return results.chosenIndex != -1
}

func (results *WeightedRollResults) ChosenItem() (item string) {
	if !results.HasChosenItem() {
		return ""
	}
	return results.items[results.chosenIndex]
}

func (results *WeightedRollResults) ChosenItemAndWeight() (item string, weight float64) {
	if !results.HasChosenItem() {
		return "", 0
	}
	item = results.items[results.chosenIndex]
	weight = results.weights[results.chosenIndex]
	return item, weight
}

// Roll Returns the random number used to determine the chosen item
func (results *WeightedRollResults) Roll() float64 {
	return results.roll
}
