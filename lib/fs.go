package lib

import (
	"os"
)

// Check exist file
func FileExists(name string) bool {

	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}
