// +build integration

package controllers

import (
	"testing"
)

func TestFailed(t *testing.T) {
	t.Error("integration failed")
}
