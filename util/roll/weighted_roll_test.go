package roll_test

import (
	"math"
	"testing"

	"github.com/oppzippy/BoostRequestBot/util/roll"
)

func TestWeightedRollWithNoItems(t *testing.T) {
	t.Parallel()
	roll := roll.NewWeightedRoll[string](0)
	results, ok := roll.Roll()
	checkNotOK(t, results, ok)
}

func TestWeightedRollWithOneItem(t *testing.T) {
	t.Parallel()
	roll := roll.NewWeightedRoll[string](1)
	roll.AddItem("test", 1)
	results, ok := roll.Roll()
	if !checkOK(t, results, ok) {
		return
	}

	chosenItem, weight := results.ChosenItemAndWeight()
	if chosenItem != "test" || weight != 1 {
		t.Errorf("the one item was not chosen: %s with weight %f", chosenItem, weight)
	}
}

func TestWeightedRollZeroWeight(t *testing.T) {
	t.Parallel()
	roll := roll.NewWeightedRoll[string](1)
	roll.AddItem("test", 0)
	results, ok := roll.Roll()
	checkNotOK(t, results, ok)
}

func TestWeightedRollResultsIterator(t *testing.T) {
	t.Parallel()
	roll := roll.NewWeightedRoll[string](2)
	roll.AddItem("one", 1)
	roll.AddItem("two", 2)
	results, ok := roll.Roll()
	if !checkOK(t, results, ok) {
		return
	}

	var numChosenItems int
	seenItems := make(map[string]struct{})
	for iter := results.Iterator(); iter.HasNext(); {
		item, _, isChosenItem := iter.Next()
		if _, ok := seenItems[item]; ok {
			t.Errorf("%s was seen twice", item)
			return
		}
		seenItems[item] = struct{}{}

		if isChosenItem {
			numChosenItems++
			if numChosenItems > 1 {
				t.Errorf("isChosenItem was true for %d items, expected 1", numChosenItems)
				return
			}
		}
	}
}

func TestRandomness(t *testing.T) {
	t.Parallel()
	iterations := 100000

	items := []string{"zero", "one", "two", "three"}
	roll := roll.NewWeightedRoll[string](len(items))
	for i, item := range items {
		roll.AddItem(item, float64(i))
	}

	wins := make(map[string]int)
	for i := 0; i < iterations; i++ {
		results, ok := roll.Roll()
		if !ok {
			t.Errorf("iteration %d failed, not ok", i)
			return
		}
		wins[results.ChosenItem()]++
	}

	for i, item := range items {
		expectedWinRate := float64(i) / float64(triangleNumber(len(items)-1))
		winRate := float64(wins[item]) / float64(iterations)

		if math.Abs(expectedWinRate-winRate) > 0.01 {
			t.Errorf(
				"expected %s to win %f%%, got %f%%",
				item,
				expectedWinRate*100,
				winRate*100,
			)
		}
	}
}

// triangleNumber returns the nth triangle number.
// It's basically factorial but with addition instead of multiplication.
func triangleNumber(n int) int {
	return (n*n + n) / 2
}

func checkOK[T any](t *testing.T, results *roll.WeightedRollResults[T], ok bool) bool {
	if !ok {
		t.Errorf("roll wasn't ok")
		return false
	}
	if results == nil {
		t.Errorf("results is nil")
		return false
	}
	return true
}

func checkNotOK[T any](t *testing.T, results *roll.WeightedRollResults[T], ok bool) bool {
	if ok {
		t.Errorf("roll was ok, expected not ok")
		return false
	}
	if results != nil {
		t.Errorf("results wasn't nil, expected nil")
		return false
	}
	return true
}