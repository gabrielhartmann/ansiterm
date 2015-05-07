// +build windows

package winterm

import (
	log "github.com/Sirupsen/logrus"
	"strconv"
)

type WindowsAnsiEventHandler struct {
	fd uintptr
}

func CreateWinEventHandler(fileDesc uintptr) *WindowsAnsiEventHandler {
	return &WindowsAnsiEventHandler{fd: fileDesc}
}

func (h *WindowsAnsiEventHandler) Print(b byte) error {
	log.Infof("Print: [%v]", []string{string(b)})
	return nil
}

func (h *WindowsAnsiEventHandler) CUU(param int) error {
	//log.Infof("CUU: [%v]", []string{strconv.Itoa(param)})
	return h.moveCursorVertical(-param)
}

func (h *WindowsAnsiEventHandler) CUD(param int) error {
	//log.Infof("CUD: [%v]", []string{strconv.Itoa(param)})
	return h.moveCursorVertical(param)
}

func (h *WindowsAnsiEventHandler) CUF(param int) error {
	//log.Infof("CUF: [%v]", []string{strconv.Itoa(param)})
	return h.moveCursorHorizontal(param)
	return nil
}

func (h *WindowsAnsiEventHandler) CUB(param int) error {
	//log.Infof("CUB: [%v]", []string{strconv.Itoa(param)})
	return h.moveCursorHorizontal(-param)
}

func (h *WindowsAnsiEventHandler) CNL(param int) error {
	//log.Infof("CNL: [%v]", []string{strconv.Itoa(param)})
	return h.moveCursorLine(param)
}

func (h *WindowsAnsiEventHandler) CPL(param int) error {
	//log.Infof("CPL: [%v]", []string{strconv.Itoa(param)})
	return h.moveCursorLine(-param)
}

func (h *WindowsAnsiEventHandler) CHA(param int) error {
	//log.Infof("CHA: [%v]", []string{strconv.Itoa(param)})
	return h.moveCursorColumn(param)
}

func (h *WindowsAnsiEventHandler) CUP(row int, col int) error {
	//rowS, colS := strconv.Itoa(row), strconv.Itoa(col)
	//log.Infof("CUP: [%v]", []string{rowS, colS})
	info, err := GetConsoleScreenBufferInfo(h.fd)
	if err != nil {
		return err
	}

	rowS := AddInRange(SHORT(row), -1, info.Window.Top, info.Window.Bottom)
	colS := AddInRange(SHORT(col), -1, info.Window.Left, info.Window.Right)
	position := COORD{colS, rowS}

	return h.setCursorPosition(position, info.Size)
}

func (h *WindowsAnsiEventHandler) HVP(row int, col int) error {
	// rowS, colS := strconv.Itoa(row), strconv.Itoa(row)
	// log.Infof("HVP: [%v]", []string{rowS, colS})
	return h.CUP(row, col)
}

func (h *WindowsAnsiEventHandler) DECTCEM(visible bool) error {
	//log.Infof("DECTCEM: [%v]", []string{strconv.FormatBool(visible)})

	return nil
}

func (h *WindowsAnsiEventHandler) ED(param int) error {
	//log.Infof("ED: [%v]", []string{strconv.Itoa(param)})

	// [J  -- Erases from the cursor to the end of the screen, including the cursor position.
	// [1J -- Erases from the beginning of the screen to the cursor, including the cursor position.
	// [2J -- Erases the complete display. The cursor does not move.
	// [3J -- Erases the complete display and backing buffer, cursor moves to (0,0)
	// Notes:
	// -- ANSI.SYS always moved the cursor to (0,0) for both [2J and [3J
	// -- Clearing the entire buffer, versus just the Window, works best for Windows Consoles

	info, err := GetConsoleScreenBufferInfo(h.fd)
	if err != nil {
		return err
	}

	var start COORD
	var end COORD

	switch param {
	case 0:
		start = info.CursorPosition
		end = COORD{info.Size.X - 1, info.Size.Y - 1}

	case 1:
		start = COORD{0, 0}
		end = info.CursorPosition

	case 2:
		start = COORD{0, 0}
		end = COORD{info.Size.X - 1, info.Size.Y - 1}

	case 3:
		start = COORD{0, 0}
		end = COORD{info.Size.X - 1, info.Size.Y - 1}
	}

	err = h.clearRange(info.Attributes, start, end)
	if err != nil {
		return err
	}

	if param == 2 || param == 3 {
		err = h.setCursorPosition(COORD{0, 0}, info.Size)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *WindowsAnsiEventHandler) EL(param int) error {
	//log.Infof("EL: [%v]", []string{strconv.Itoa(param)})

	// [K  -- Erases from the cursor to the end of the line, including the cursor position.
	// [1K -- Erases from the beginning of the line to the cursor, including the cursor position.
	// [2K -- Erases the complete line.

	info, err := GetConsoleScreenBufferInfo(h.fd)
	if err != nil {
		return err
	}

	var start COORD
	var end COORD

	switch param {
	case 0:
		start = info.CursorPosition
		end = COORD{info.Window.Right, info.CursorPosition.Y}

	case 1:
		start = COORD{0, info.CursorPosition.Y}
		end = info.CursorPosition

	case 2:
		start = COORD{0, info.CursorPosition.Y}
		end = COORD{info.Window.Right, info.CursorPosition.Y}
	}

	err = h.clearRange(info.Attributes, start, end)
	if err != nil {
		return err
	}

	return nil
}

func (h *WindowsAnsiEventHandler) SGR(params []int) error {
	//log.Infof("SGR: [%v]", strings)
	// strings := []string{}
	// for _, v := range params {
	// 	log.Infof("SGR: [%v]", strings)
	// 	strings = append(strings, strconv.Itoa(v))
	// }

	return nil
}

func (h *WindowsAnsiEventHandler) SU(param int) error {
	log.Infof("SU: [%v]", []string{strconv.Itoa(param)})
	return nil
}

func (h *WindowsAnsiEventHandler) SD(param int) error {
	log.Infof("SD: [%v]", []string{strconv.Itoa(param)})
	return nil
}
