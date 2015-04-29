package ansiterm

import (
	"testing"
)

func stateTransitionHelper(t *testing.T, start State, end State, bytes []byte) {
	for _, b := range bytes {
		bytes := []byte{byte(b)}
		parser, _ := createTestParser(start)
		parser.Parse(bytes)
		validateState(t, parser.state, end)
	}
}

func anyToXHelper(t *testing.T, bytes []byte, expectedState State) {
	for i := 0; i < int(stateCount); i++ {
		s := stateMap[StateId(i)]
		stateTransitionHelper(t, s, expectedState, bytes)
	}
}

func funcCallParamHelper(t *testing.T, bytes []byte, expectedState State, expectedCalls []string) {
	parser, evtHandler := createTestParser(CsiEntry)
	parser.Parse(bytes)
	validateState(t, parser.state, expectedState)
	validateFuncCalls(t, evtHandler.FunctionCalls, expectedCalls)
}

func parseParamsHelper(t *testing.T, bytes []byte, expectedParams []string) {
	params, err := parseParams(bytes)

	if err != nil {
		t.Errorf("Parameter parse error: %v", err)
		return
	}

	if len(params) != len(expectedParams) {
		t.Errorf("Parsed   parameters: %v", params)
		t.Errorf("Expected parameters: %v", expectedParams)
		t.Errorf("Parameter length failure: %d != %d", len(params), len(expectedParams))
		return
	}

	for i, v := range expectedParams {
		if v != params[i] {
			t.Errorf("Parsed   parameters: %v", params)
			t.Errorf("Expected parameters: %v", expectedParams)
			t.Errorf("Parameter parse failure: %s != %s at position %d", v, params[i], i)
		}
	}
}
