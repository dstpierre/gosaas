package cache

import (
	"testing"

	"time"

	"gopkg.in/mgo.v2/bson"
)

type auth struct {
	ID        string
	Email     string
	Token     string
	AccountID bson.ObjectId
	LoginID   bson.ObjectId
	Role      int
}

func TestAuth_Exists_NotPresent(t *testing.T) {
	t.Parallel()

	ca := &Auth{}
	var keys auth
	if err := ca.Exists("key-that-does-not-exists", &keys); err != nil {
		t.Error(err)
	} else if len(keys.ID) != 0 {
		t.Errorf("keys.ID is not empty %s", keys.ID)
	}
}

func TestAuth_Exists_Present(t *testing.T) {
	t.Parallel()

	ca := &Auth{}
	keys := auth{ID: "my-id",
		Email:     "me@domain.com",
		Token:     "super-cool-tok",
		AccountID: bson.NewObjectId(),
		LoginID:   bson.NewObjectId(),
		Role:      1,
	}

	key := "testing-presence"
	if err := ca.Set(key, keys, 5*time.Second); err != nil {
		t.Fatal(err)
	}

	var checks auth
	if err := ca.Exists(key, &checks); err != nil {
		t.Fatal(err)
	} else if checks.ID != keys.ID {
		t.Errorf("received id %s but should have been %s", checks.ID, keys.ID)
	}
}

func TestAuth_Exists_Expiration(t *testing.T) {
	t.Parallel()

	ca := &Auth{}
	keys := auth{ID: "my-id",
		Email:     "me@domain.com",
		Token:     "super-cool-tok",
		AccountID: bson.NewObjectId(),
		LoginID:   bson.NewObjectId(),
		Role:      1,
	}

	key := "testing-expiration"
	if err := ca.Set(key, keys, 1*time.Second); err != nil {
		t.Fatal(err)
	}

	time.Sleep(2 * time.Second)

	var checks auth
	if err := ca.Exists(key, &checks); err != nil {
		t.Fatal(err)
	} else if len(checks.ID) > 0 {
		t.Errorf("received non 0 length for ID: %s", checks.ID)
	}
}
