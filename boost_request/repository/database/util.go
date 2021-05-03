package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

func SQLSet(length int) string {
	if length <= 0 {
		if length < 0 {
			log.Printf("Tried to create a sql placeholder set of length %d; should be >=1", length)
		}
		return "(NULL)"
	}
	return fmt.Sprintf("(?%s)", strings.Repeat(", ?", length-1))
}

func SQLSets(setLen, numSets int) string {
	if numSets <= 0 {
		log.Printf("Tried to create %d sql placeholder sets; should be >=1", numSets)
		return ""
	}
	set := SQLSet(setLen)

	sets := set + strings.Repeat(fmt.Sprintf(", %s", set), numSets-1)
	return sets
}

func rollbackIfErr(tx *sql.Tx, err error) error {
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("%v; caused by %v", rbErr, err)
		}
		return err
	}
	return nil
}
