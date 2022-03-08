package util

import "github.com/satori/go.uuid"

func UUId() string {
	return uuid.NewV4().String()
}
