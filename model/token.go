package model

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

// NewToken returns a token combining an id with a unique identifier.
func NewToken(id int64) string {
	return fmt.Sprintf("%d|%s", id, uuid.NewV4().String())
}

// ParseToken returns the id and uuid for a given token.
func ParseToken(token string) (int64, string) {
	pairs := strings.Split(token, "|")
	if len(pairs) != 2 {
		return -1, ""
	}

	id, err := strconv.ParseInt(pairs[0], 10, 64)
	if err != nil {
		return -1, ""
	}
	return id, pairs[1]
}

// NewFriendlyID returns a ~somewhat unique friendly id.
func NewFriendlyID(id int64, key string) string {
	n := time.Now()
	i, _ := strconv.Atoi(
		fmt.Sprintf("%d%d%d%d%d%d%d%d%d",
			id,
			len(key),
			n.Year()-2000,
			int(n.Month()),
			n.Day(),
			n.Hour(),
			n.Minute(),
			n.Second(),
			n.Nanosecond()))
	return fmt.Sprintf("%x", i)
}

func StringToKey(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Printf("error converting %s to int64\n", s)
		return -1
	}
	return i
}
