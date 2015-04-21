package main

import "testing"

func TestValidName(t *testing.T) {
	validNames := []string{"JoeBob99", "JoeBob_99", "joeBoB99"}
	invalidNames := []string{"99JoeBob", " JoeBob99", "_jobBob99", "Joe Bob99"}

	for _, validName := range validNames {
		if !isValidName(validName) {
			t.Errorf("%q should be valid", validName)
		}
	}

	for _, invalidName := range invalidNames {
		if isValidName(invalidName) {
			t.Errorf("%q should be invalid", invalidName)
		}
	}
}
