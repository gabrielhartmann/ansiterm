// +build windows

package winterm

func (h *WindowsAnsiEventHandler) setCursorOnScreen(x int, y int) error {
	info, err := GetConsoleScreenBufferInfo(h.fd)
	if err != nil {
		return err
	}

	xs := SHORT(x) + info.Window.Left
	ys := SHORT(y) + info.Window.Top

	position := COORD{xs, ys}

	if err := h.setCursorPosition(position, info.Size); err != nil {
		return err
	}

	return nil
}
