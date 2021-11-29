package utils

import (
	"testing"
)

func TestGenerateNumber(t *testing.T) {
	res := GenerateNumber()

	if res != "KZ" {
		t.Errorf("want: %v, got: %v", "KZ", res)
	}
}
