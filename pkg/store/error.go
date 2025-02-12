package store

type NotFoundError struct {
	label string
}

func newNotFoundError(label string) *NotFoundError {
	return &NotFoundError{label}
}

func (e *NotFoundError) Error() string {
	return e.label + " not found"
}
