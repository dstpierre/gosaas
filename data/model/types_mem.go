// +build mem

package model

import (
	"fmt"
	"strconv"
)

type Connection = bool
type Key = int

func Open(options ...string) (bool, error) {
	return true, nil
}

func KeyToString(id Key) string {
	return fmt.Sprintf("%d", id)
}

func StringToKey(id string) int {
	i, err := strconv.Atoi(id)
	if err != nil {
		return -1
	}
	return i
}

func NewID() Key {
	return 1
}
