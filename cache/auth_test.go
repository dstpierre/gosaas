package cache

import (
	"testing"

	"time"
)

type auth struct {
	AccountID int64
	UserID    int64
	Email     string
	Role      int
}

func TestAuth_Exists_NotPresent(t *testing.T) {
	t.Parallel()

	ca := &Auth{}
	var keys auth
	if err := ca.Exists("key-that-does-not-exists", &keys); err != nil {
		t.Error(err)
	} else if keys.AccountID != 0 {
		t.Errorf("keys.AccountID is not empty %d", keys.AccountID)
	}
}

func TestAuth_Exists_Present(t *testing.T) {
	t.Parallel()

	ca := &Auth{}
	keys := auth{AccountID: 123,
		UserID: 321,
		Email:  "unit@test.com",
		Role:   99,
	}

	key := "testing-presence"
	if err := ca.Set(key, keys, 5*time.Second); err != nil {
		t.Fatal(err)
	}

	var checks auth
	if err := ca.Exists(key, &checks); err != nil {
		t.Fatal(err)
	} else if checks.AccountID != keys.AccountID {
		t.Errorf("received id %d but should have been %d", checks.AccountID, keys.AccountID)
	}
}

func TestAuth_Exists_Expiration(t *testing.T) {
	t.Parallel()

	ca := &Auth{}
	keys := auth{AccountID: 321,
		UserID: 123,
		Email:  "unit@test.com",
		Role:   99,
	}

	key := "testing-expiration"
	if err := ca.Set(key, keys, 1*time.Second); err != nil {
		t.Fatal(err)
	}

	time.Sleep(2 * time.Second)

	var checks auth
	if err := ca.Exists(key, &checks); err != nil {
		t.Fatal(err)
	} else if checks.AccountID > 0 {
		t.Errorf("received empty id 0 length for ID: %d", checks.AccountID)
	}
}
