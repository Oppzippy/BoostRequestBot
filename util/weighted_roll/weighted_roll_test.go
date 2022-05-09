package weighted_roll_test

import (
	"math"
	"testing"

	"github.com/oppzippy/BoostRequestBot/util/weighted_roll"
)

func TestWeightedRollWithNoItems(t *testing.T) {
	t.Parallel()
	roll := weighted_roll.NewWeightedRoll[string](0)
	results := roll.Roll()
	if results == nil || len(results) != 0 {
		t.Errorf("expected empty slice, got length %v", len(results))
	}
}

func TestWeightedRollWithOneItem(t *testing.T) {
	t.Parallel()
	roll := weighted_roll.NewWeightedRoll[string](1)
	roll.AddItem("test", 1)
	results := roll.Roll()

	if results[0].Item != "test" || results[0].Weight != 1 {
		t.Errorf("the one Item was not chosen: %s with Weight %f", results[0].Item, results[0].Weight)
	}
}

func TestWeightedRollZeroWeight(t *testing.T) {
	t.Parallel()
	roll := weighted_roll.NewWeightedRoll[string](1)
	roll.AddItem("test", 0)
	results := roll.Roll()

	if results[0].Item != "test" || results[0].Roll != 0 {
		t.Errorf("first result is %v with roll %v", results[0].Item, results[0].Roll)
	}
}

func TestRandomness(t *testing.T) {
	t.Parallel()
	iterations := 100000

	items := []string{"zero", "one", "two", "three"}
	roll := weighted_roll.NewWeightedRoll[string](len(items))
	for i, item := range items {
		roll.AddItem(item, float64(i))
	}

	wins := make(map[string]int)
	for i := 0; i < iterations; i++ {
		results := roll.Roll()
		wins[results[0].Item]++
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
