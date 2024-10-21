package status

type Error struct {
	Status int
	Error  string
}

func (se *Error) StatusCode() int {
	return se.Status
}

func (se *Error) Message() string {
	return se.Error
}
