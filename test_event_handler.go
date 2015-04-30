package ansiterm

import (
	"fmt"
	"strconv"
)

type TestAnsiEventHandler struct {
	FunctionCalls []string
}

func CreateTestAnsiEventHandler() TestAnsiEventHandler {
	evtHandler := TestAnsiEventHandler{}
	evtHandler.FunctionCalls = make([]string, 0)
	return evtHandler
}

func (h *TestAnsiEventHandler) recordCall(call string, params []string) {
	s := fmt.Sprintf("%s(%v)", call, params)
	h.FunctionCalls = append(h.FunctionCalls, s)
}

func (h *TestAnsiEventHandler) CUU(param int) error {
	h.recordCall("CUU", []string{strconv.Itoa(param)})
	return nil
}

func (h *TestAnsiEventHandler) CUD(param int) error {
	h.recordCall("CUD", []string{strconv.Itoa(param)})
	return nil
}

func (h *TestAnsiEventHandler) CUF(param int) error {
	h.recordCall("CUF", []string{strconv.Itoa(param)})
	return nil
}

func (h *TestAnsiEventHandler) CUB(param int) error {
	h.recordCall("CUB", []string{strconv.Itoa(param)})
	return nil
}

func (h *TestAnsiEventHandler) CNL(param int) error {
	h.recordCall("CNL", []string{strconv.Itoa(param)})
	return nil
}

func (h *TestAnsiEventHandler) CPL(param int) error {
	h.recordCall("CPL", []string{strconv.Itoa(param)})
	return nil
}

func (h *TestAnsiEventHandler) CHA(param int) error {
	h.recordCall("CHA", []string{strconv.Itoa(param)})
	return nil
}

func (h *TestAnsiEventHandler) CUP(x int, y int) error {
	xS, yS := strconv.Itoa(x), strconv.Itoa(y)
	h.recordCall("CUP", []string{xS, yS})
	return nil
}

func (h *TestAnsiEventHandler) HVP(x int, y int) error {
	xS, yS := strconv.Itoa(x), strconv.Itoa(y)
	h.recordCall("HVP", []string{xS, yS})
	return nil
}

func (h *TestAnsiEventHandler) DECTCEM(visible bool) error {
	h.recordCall("DECTCEM", []string{strconv.FormatBool(visible)})
	return nil
}

func (h *TestAnsiEventHandler) ED(param int) error {
	h.recordCall("ED", []string{strconv.Itoa(param)})
	return nil
}

func (h *TestAnsiEventHandler) EL(param int) error {
	h.recordCall("EL", []string{strconv.Itoa(param)})
	return nil
}

func (h *TestAnsiEventHandler) SGR(params []int) error {
	strings := []string{}
	for _, v := range params {
		strings = append(strings, strconv.Itoa(v))
	}

	h.recordCall("SGR", strings)
	return nil
}

func (h *TestAnsiEventHandler) SU(param int) error {
	h.recordCall("SU", []string{strconv.Itoa(param)})
	return nil
}

func (h *TestAnsiEventHandler) SD(param int) error {
	h.recordCall("SD", []string{strconv.Itoa(param)})
	return nil
}
