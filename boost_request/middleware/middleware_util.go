package middleware

import "github.com/lus/dgc"

func commandHasFlag(command *dgc.Command, flag string) bool {
	for _, f := range command.Flags {
		if f == flag {
			return true
		}
	}
	return false
}
