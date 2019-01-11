package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

// NewToken returns a token combining an id with a unique identifier
func NewToken(id Key) string {
	return fmt.Sprintf("%s|%s", KeyToString(id), uuid.NewV4())
}

// ParseToken returns the id and uuid for a given token
func ParseToken(token string) (string, string) {
	pairs := strings.Split(token, "|")
	if len(pairs) != 2 {
		return "", ""
	}
	return pairs[0], pairs[1]
}

func NewFriendlyID() string {
	n := time.Now()
	i, _ := strconv.Atoi(
		fmt.Sprintf("%d%d%d%d%d%d",
			n.Year()-2000,
			int(n.Month()),
			n.Day(),
			n.Hour(),
			n.Minute(),
			n.Second()))
	return fmt.Sprintf("%x", i)
}
