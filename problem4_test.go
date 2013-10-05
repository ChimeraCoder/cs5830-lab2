package main

import (
	"testing"
)

func Test_Exponentiation(t *testing.T) {
	if Exp(3, 16, 3000) != 2721 {
		t.Error("exp failed")
	} else {
		t.Log("one test passed.") // log some info if you want
	}
	if Exp(2, 10, 5) != 4 {
		t.Error("exp failed")
	}
}

func Test_invert(t *testing.T) {
	if invert(101, 102) != -1 {
		t.Error("Invert failed")
	}
}
