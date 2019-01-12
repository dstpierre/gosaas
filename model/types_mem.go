// +build mem

package model

import (
	"fmt"
	"strconv"
)

// Connection is the database connection type.
type Connection = bool

// Key is the primary key type.
type Key = int

// Open do nothing since this data provider is in-memory.
func Open(options ...string) (bool, error) {
	return true, nil
}

// KeyToString converts a Key to a string.
func KeyToString(id Key) string {
	return fmt.Sprintf("%d", id)
}

// StringToKey converts a string to a Key.
func StringToKey(id string) int {
	i, err := strconv.Atoi(id)
	if err != nil {
		return -1
	}
	return i
}

// NewID returns a new Key.
func NewID() Key {
	return 1
}
