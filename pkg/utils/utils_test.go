package utils

import "testing"

func TestContainsString(t *testing.T) {
	if !ContainsString([]string{"a", "b"}, "a") {
		t.Fail()
	}
}
