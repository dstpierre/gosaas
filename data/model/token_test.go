// +build mem

package model

import "testing"

func Test_Model_NewToken(t *testing.T) {
	token := NewToken(123)
	id, uid := ParseToken(token)
	if id != "123" {
		t.Error("token id not valid expected 123 got", id)
	} else if uid == "" {
		t.Error("token uuid empty, expected a value")
	}
}
