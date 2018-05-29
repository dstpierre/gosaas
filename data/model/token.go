package model

import (
	"fmt"
	"strings"

	uuid "github.com/satori/go.uuid"
)

// NewToken returns a token combining an id with a unique identifier
func NewToken(id Key) string {
	return fmt.Sprintf("%s|%s", keyToString(id), uuid.NewV4().String())
}

// ParseToken returns the id and uuid for a given token
func ParseToken(token string) (string, string) {
	pairs := strings.Split(token, "|")
	if len(pairs) != 2 {
		return "", ""
	}
	return pairs[0], pairs[1]
}
