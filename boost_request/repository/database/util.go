package database

import (
	"fmt"
	"log"
	"strings"
)

func SQLSet(length int) string {
	if length <= 0 {
		log.Printf("Tried to create a sql placeholder set of length %d; should be >=1", length)
		return "(NULL)"
	}
	return fmt.Sprintf("(?%s)", strings.Repeat(", ?", length-1))
}
