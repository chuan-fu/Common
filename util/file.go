package util

import "os"

func IsExist(path string) bool {
	if _, err := os.Stat(path); err == nil || os.IsExist(err) {
		return true
	}
	return false
}
