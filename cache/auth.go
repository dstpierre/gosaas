package cache

import (
	"bytes"
	"encoding/gob"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

// Auth is used to get/set authentication related keys
type Auth struct{}

// Exists returns the authentication in cache
func (x *Auth) Exists(key string, v interface{}) error {
	s, err := rc.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return err
	}

	dec := gob.NewDecoder(strings.NewReader(s))
	return dec.Decode(v)
}

// Set cache this key for 30 minutes
func (x *Auth) Set(key string, v interface{}, expiration time.Duration) error {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(v); err != nil {
		return err
	}

	return rc.Set(key, buf.String(), expiration).Err()
}
