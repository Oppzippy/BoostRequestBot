package database_test

import (
	"testing"

	"github.com/oppzippy/BoostRequestBot/boost_request/repository/database"
)

func TestSQLSet(t *testing.T) {
	t.Run("length -1", func(t *testing.T) {
		set := database.SQLSet(0)
		if set != "(NULL)" {
			t.Errorf("Set should be (NULL), got %s", set)
		}
	})
	t.Run("length 1", func(t *testing.T) {
		set := database.SQLSet(1)
		if set != "(?)" {
			t.Errorf("Set should be (?), got %s", set)
		}
	})
	t.Run("length 3", func(t *testing.T) {
		set := database.SQLSet(3)
		if set != "(?, ?, ?)" {
			t.Errorf("Set should be (?, ?, ?), got %s", set)
		}
	})
}
