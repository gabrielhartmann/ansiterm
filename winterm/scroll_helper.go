// +build windows

package winterm

func (h *WindowsAnsiEventHandler) scroll(param int) error {

	info, err := GetConsoleScreenBufferInfo(h.fd)
	if err != nil {
		return err
	}

	rect := info.Window

	// Current scroll region in Windows backing buffer coordinates
	top := rect.Top + SHORT(h.sr.top)
	bottom := rect.Top + SHORT(h.sr.bottom)

	// Area from backing buffer to be copied
	scrollRect := SMALL_RECT{
		Top:    top + SHORT(param),
		Bottom: bottom + SHORT(param),
		Left:   rect.Left,
		Right:  rect.Right,
	}

	// Clipping region should be the original scroll region
	clipRegion := SMALL_RECT{
		Top:    top,
		Bottom: bottom,
		Left:   rect.Left,
		Right:  rect.Right,
	}

	// Origin to which area should be copied
	destOrigin := COORD{
		X: rect.Left,
		Y: top,
	}

	char := CHAR_INFO{
		UnicodeChar: ' ',
		Attributes:  0,
	}

	if err := ScrollConsoleScreenBuffer(h.fd, scrollRect, clipRegion, destOrigin, char); err != nil {
		return err
	}

	return nil
}
