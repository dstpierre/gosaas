package cache

import (
	"fmt"
	"time"
)

// Throttle increments the requests count for a specific key and set expiration if it's a new period.
func Throttle(key string, expire time.Duration) (int64, error) {
	key = fmt.Sprintf("%s_t", key)
	return increaseThrottle(key, expire)
}

// RateLimit increments the requests count for a specific key and set expiration if it's a new period.
func RateLimit(key string, expire time.Duration) (int64, error) {
	key = fmt.Sprintf("%s_rl", key)
	return increaseThrottle(key, expire)
}

func increaseThrottle(key string, expire time.Duration) (int64, error) {
	i, err := rc.Incr(key).Result()
	if err != nil {
		return 0, err
	}

	if i == 1 {
		// the key was created, we set the expire
		ok, err := rc.Expire(key, expire).Result()
		if err != nil {
			// try to remove the key
			if _, e := rc.Del(key).Result(); err != nil {
				return 0, fmt.Errorf("unable to remove key %s: %s and expire failed: %s", key, e.Error(), err.Error())
			}
			return 0, err
		} else if !ok {
			return 0, fmt.Errorf("unable to set expiration on key %s", key)
		}
	}

	return i, nil
}

// GetThrottleExpiration returns the duration before a key expire for throttling
func GetThrottleExpiration(key string) (time.Duration, error) {
	key = fmt.Sprintf("%s_t", key)

	return rc.TTL(key).Result()
}

// GetRateLimitExpiration returns the duration before a key expire for rate limit.
func GetRateLimitExpiration(key string) (time.Duration, error) {
	key = fmt.Sprintf("%s_rl", key)

	return rc.TTL(key).Result()
}
