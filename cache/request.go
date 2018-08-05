package cache

import (
	"strings"

	"fmt"

	"github.com/go-redis/redis"
)

// CountWebRequest returns the number of failed web request pending for analysis
func CountWebRequest() (int64, error) {
	return rc.LLen("reqs").Result()
}

// GetWebRequest returns the next web request logged from list
func GetWebRequest(first bool) (reqID string, b []byte, err error) {
	var s string
	if first {
		s, err = rc.LPop("reqs").Result()
	} else {
		s, err = rc.RPop("reqs").Result()
	}

	if err != nil {
		if err == redis.Nil {
			err = nil
		}
		return
	}

	buf := strings.Split(s, "\n|\n")
	if len(buf) != 2 {
		return reqID, b, fmt.Errorf("unable to split request result")
	}

	b = []byte(buf[0])
	reqID = buf[1]

	return
}

// LogWebRequest saves a web request for further analysis
func LogWebRequest(reqID string, b []byte) error {
	r := []byte(fmt.Sprintf("\n|\n%s", reqID))
	b = append(b, r...)
	return rc.RPush("reqs", b).Err()
}
