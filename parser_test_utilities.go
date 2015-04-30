package ansiterm

import (
	"testing"
)

func createTestParser(s State) (AnsiParser, *TestAnsiEventHandler) {
	evtHandler := CreateTestAnsiEventHandler()
	parser := CreateParser(s, &evtHandler)

	return parser, &evtHandler
}

func validateState(t *testing.T, actualState State, expectedState State) {
	actualName := "Nil"
	expectedName := "Nil"

	if actualState != nil {
		actualName = actualState.Name()
	}

	if expectedState != nil {
		expectedName = expectedState.Name()
	}

	if actualState != expectedState {
		t.Errorf("Invalid State: '%s' != '%s'", actualName, expectedName)
	}
}

func validateFuncCalls(t *testing.T, actualCalls []string, expectedCalls []string) {
	actualCount := len(actualCalls)
	expectedCount := len(expectedCalls)

	if actualCount != expectedCount {
		t.Errorf("Actual   calls: %v", actualCalls)
		t.Errorf("Expected calls: %v", expectedCalls)
		t.Errorf("Call count error: %d != %d", actualCount, expectedCount)
	}

	for i, v := range actualCalls {
		if v != expectedCalls[i] {
			t.Errorf("Actual   calls: %v", actualCalls)
			t.Errorf("Expected calls: %v", expectedCalls)
			t.Errorf("Mismatched calls: %s != %s with lengths %d and %d", v, expectedCalls[i], len(v), len(expectedCalls[i]))
		}
	}
}
