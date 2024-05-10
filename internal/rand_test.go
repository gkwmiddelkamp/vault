package internal

import (
	"testing"
)

func TestRandomString(t *testing.T) {
	result := RandomString(50)
	if len(result) != 50 {
		t.Fatal("RandomString(50), should result in string of 50 long")
	}
	result2 := RandomString(50)
	if string(result) == string(result2) {
		t.Fatal("RandomString() should always return unique result")
	}
}
