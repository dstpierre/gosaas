package model

import (
	"testing"
)

func Test_Model_NewToken(t *testing.T) {
	token := NewToken(123)
	id, uid := ParseToken(token)
	if id != 123 {
		t.Error("token id not valid expected 123 got", id)
	} else if uid == "" {
		t.Error("token uuid empty, expected a value")
	}
}

func Test_Model_ParseToken(t *testing.T) {
	token := "1|e21ce1fd-0e20-4fbe-b014-378767bb2e97"
	id, tok := ParseToken(token)
	if id != 1 {
		t.Errorf("expected 1 as id got %d", id)
	} else if tok != "e21ce1fd-0e20-4fbe-b014-378767bb2e97" {
		t.Errorf("expected e21ce1fd-0e20-4fbe-b014-378767bb2e97 as token got %s", tok)
	}
}
