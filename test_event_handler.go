package ansiterm

import (
	"fmt"
)

type TestAnsiEventHandler struct {
	FunctionCalls []string
}

func CreateTestAnsiEventHandler() TestAnsiEventHandler {
	evtHandler := TestAnsiEventHandler{}
	evtHandler.FunctionCalls = make([]string, 0)
	return evtHandler
}

func (h *TestAnsiEventHandler) recordCall(call string, params []int) {
	s := fmt.Sprintf("%s(%v)", call, params)
	h.FunctionCalls = append(h.FunctionCalls, s)
}

func (h *TestAnsiEventHandler) CUU(params []int) error {
	h.recordCall("CUU", params)
	return nil
}

func (h *TestAnsiEventHandler) CUD(params []int) error {
	h.recordCall("CUD", params)
	return nil
}

func (h *TestAnsiEventHandler) CUF(params []int) error {
	h.recordCall("CUF", params)
	return nil
}

func (h *TestAnsiEventHandler) CUB(params []int) error {
	h.recordCall("CUB", params)
	return nil
}

func (h *TestAnsiEventHandler) CNL(params []int) error {
	h.recordCall("CNL", params)
	return nil
}

func (h *TestAnsiEventHandler) CPL(params []int) error {
	h.recordCall("CPL", params)
	return nil
}

func (h *TestAnsiEventHandler) CHA(params []int) error {
	h.recordCall("CHA", params)
	return nil
}

func (h *TestAnsiEventHandler) CUP(params []int) error {
	h.recordCall("CUP", params)
	return nil
}

func (h *TestAnsiEventHandler) HVP(params []int) error {
	h.recordCall("HVP", params)
	return nil
}
