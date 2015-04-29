package ansiterm

type AnsiEventHandler interface {
	// CUrsor Up
	CUU([]string) error

	// CUrsor Down
	CUD([]string) error

	// CUrsor Forward
	CUF([]string) error

	// CUrsor Backward
	CUB([]string) error

	// Cursor to Next Line
	CNL([]string) error

	// Cursor to Previous Line
	CPL([]string) error

	// Cursor Horizontal position Absolute
	CHA([]string) error

	// CUrsor Position
	CUP([]string) error

	// Horizontal and Vertical Position (depends on PUM)
	HVP([]string) error
}
