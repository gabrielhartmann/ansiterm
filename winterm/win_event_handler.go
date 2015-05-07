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

// setCursorPosition sets the cursor to the specified position, bounded to the buffer size
func (h *WindowsAnsiEventHandler) setCursorPosition(position COORD, sizeBuffer COORD) error {
	position.X = ensureInRange(position.X, 0, sizeBuffer.X-1)
	position.Y = ensureInRange(position.Y, 0, sizeBuffer.Y-1)
	return SetConsoleCursorPosition(h.fd, position)
}

func (h *WindowsAnsiEventHandler) Print(b byte) error {
	log.Infof("Print: [%v]", []string{string(b)})
	return nil
}

func (h *WindowsAnsiEventHandler) CUU(param int) error {
	log.Infof("CUU: [%v]", []string{strconv.Itoa(param)})
	return nil
}

func (h *WindowsAnsiEventHandler) CUD(param int) error {
	log.Infof("CUD: [%v]", []string{strconv.Itoa(param)})
	return nil
}

func (h *WindowsAnsiEventHandler) CUF(param int) error {
	log.Infof("CUF: [%v]", []string{strconv.Itoa(param)})
	return nil
}

func (h *WindowsAnsiEventHandler) CUB(param int) error {
	log.Infof("CUB: [%v]", []string{strconv.Itoa(param)})
	return nil
}

func (h *WindowsAnsiEventHandler) CNL(param int) error {
	log.Infof("CNL: [%v]", []string{strconv.Itoa(param)})
	return nil
}

func (h *WindowsAnsiEventHandler) CPL(param int) error {
	log.Infof("CPL: [%v]", []string{strconv.Itoa(param)})
	return nil
}

func (h *WindowsAnsiEventHandler) CHA(param int) error {
	log.Infof("CHA: [%v]", []string{strconv.Itoa(param)})
	return nil
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

func (h *WindowsAnsiEventHandler) HVP(x int, y int) error {
	xS, yS := strconv.Itoa(x), strconv.Itoa(y)
	log.Infof("HVP: [%v]", []string{xS, yS})
	return nil
}

func (h *WindowsAnsiEventHandler) DECTCEM(visible bool) error {
	log.Infof("DECTCEM: [%v]", []string{strconv.FormatBool(visible)})
	return nil
}

func (h *WindowsAnsiEventHandler) ED(param int) error {
	log.Infof("ED: [%v]", []string{strconv.Itoa(param)})
	return nil
}

func (h *WindowsAnsiEventHandler) EL(param int) error {
	log.Infof("EL: [%v]", []string{strconv.Itoa(param)})
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
