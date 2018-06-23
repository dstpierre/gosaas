package cache

import (
	"fmt"
	"testing"
	"time"
)

func TestThrottle_CreateNewKeyIfNotExists(t *testing.T) {
	count, err := Throttle(fmt.Sprintf("%v", time.Now().Unix()), 1*time.Second)
	if err != nil {
		t.Error(err)
	} else if count != 1 {
		t.Error("expected count to be 1 was", count)
	}
}

func TestThrottle_IncreaseForAKey(t *testing.T) {
	key := "unittest"

	var c int64
	var err error

	for i := 0; i < 3; i++ {
		c, err = Throttle(key, 3*time.Second)
		if err != nil {
			t.Fatal(err)
		}
	}
	if c != 3 {
		t.Error("increase count should be 3 got", c)
	}
}

func TestThrottle_WhenExpireReachedResetToZero(t *testing.T) {
	key := "another_unittest_throttling"
	if _, err := Throttle(key, 1*time.Second); err != nil {
		t.Fatal(err)
	}

	time.Sleep(1100 * time.Millisecond)

	c, err := Throttle(key, 1*time.Second)
	if err != nil {
		t.Fatal(err)
	} else if c != 1 {
		t.Error("thottle expiration should have reset this key to 1 got", c)
	}
}

func TestThrottle_GetExpirationDuration(t *testing.T) {
	key := fmt.Sprintf("ttl_unittest_throttling_%d", time.Now().Unix())
	if _, err := Throttle(key, 5*time.Minute); err != nil {
		t.Fatal(err)
	}

	d, err := GetThrottleExpiration(key)
	if err != nil {
		t.Fatal(err)
	} else if d.Seconds() < 280 {
		t.Log(d)
		t.Error("duration minutes should have been > 4 got", d.Minutes())
	}
}
