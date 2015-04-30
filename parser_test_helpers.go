package ansiterm

import (
	"fmt"
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

func cursorSingleParamHelper(t *testing.T, command byte, funcName string) {
	funcCallParamHelper(t, []byte{command}, Ground, []string{fmt.Sprintf("%s([1])", funcName)})
	funcCallParamHelper(t, []byte{'2', command}, Ground, []string{fmt.Sprintf("%s([2])", funcName)})
	funcCallParamHelper(t, []byte{'2', '3', command}, Ground, []string{fmt.Sprintf("%s([23])", funcName)})
	funcCallParamHelper(t, []byte{'2', ';', '3', command}, Ground, []string{fmt.Sprintf("%s([2])", funcName)})
	funcCallParamHelper(t, []byte{'2', ';', '3', ';', '4', command}, Ground, []string{fmt.Sprintf("%s([2])", funcName)})
}

func cursorTwoParamHelper(t *testing.T, command byte, funcName string) {
	funcCallParamHelper(t, []byte{command}, Ground, []string{fmt.Sprintf("%s([1 1])", funcName)})
	funcCallParamHelper(t, []byte{'2', command}, Ground, []string{fmt.Sprintf("%s([2 1])", funcName)})
	funcCallParamHelper(t, []byte{'2', '3', command}, Ground, []string{fmt.Sprintf("%s([23 1])", funcName)})
	funcCallParamHelper(t, []byte{'2', ';', '3', command}, Ground, []string{fmt.Sprintf("%s([2 3])", funcName)})
	funcCallParamHelper(t, []byte{'2', ';', '3', ';', '4', command}, Ground, []string{fmt.Sprintf("%s([2 3])", funcName)})
}

func eraseHelper(t *testing.T, command byte, funcName string) {
	funcCallParamHelper(t, []byte{command}, Ground, []string{fmt.Sprintf("%s([0])", funcName)})
	funcCallParamHelper(t, []byte{'0', command}, Ground, []string{fmt.Sprintf("%s([0])", funcName)})
	funcCallParamHelper(t, []byte{'1', command}, Ground, []string{fmt.Sprintf("%s([1])", funcName)})
	funcCallParamHelper(t, []byte{'2', command}, Ground, []string{fmt.Sprintf("%s([2])", funcName)})
	funcCallParamHelper(t, []byte{'3', command}, Ground, []string{fmt.Sprintf("%s([3])", funcName)})
	funcCallParamHelper(t, []byte{'4', command}, Ground, []string{fmt.Sprintf("%s([0])", funcName)})
	funcCallParamHelper(t, []byte{'1', ';', '2', command}, Ground, []string{fmt.Sprintf("%s([1])", funcName)})
}

func panHelper(t *testing.T, command byte, funcName string) {
	funcCallParamHelper(t, []byte{command}, Ground, []string{fmt.Sprintf("%s([0])", funcName)})
	funcCallParamHelper(t, []byte{'0', command}, Ground, []string{fmt.Sprintf("%s([0])", funcName)})
	funcCallParamHelper(t, []byte{'1', command}, Ground, []string{fmt.Sprintf("%s([1])", funcName)})
	funcCallParamHelper(t, []byte{'5', command}, Ground, []string{fmt.Sprintf("%s([5])", funcName)})
	funcCallParamHelper(t, []byte{'4', ';', '6', command}, Ground, []string{fmt.Sprintf("%s([4])", funcName)})
}

func clearOnStateChangeHelper(t *testing.T, start State, end State, bytes []byte) {
	p, _ := createTestParser(start)
	fillContext(p.context)
	p.Parse(bytes)
	validateState(t, p.state, end)
	validateEmptyContext(t, p.context)
}
