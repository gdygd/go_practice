package general

type errorString struct {
	s string
}

func UserError(text string) error {
	return &errorString{text}
}

func (e *errorString) Error() string {
	return e.s
}
