package customerrors

import "fmt"

type NotFoundError struct {
	ID string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("item with ID %s not found", e.ID)
}

type NotInitError struct {
	ID string
}

func (e *NotInitError) Error() string {
	return fmt.Sprintf("camera with ID %s not initialized", e.ID)
}

type AlreadyInitError struct {
	ID string
}

func (e *AlreadyInitError) Error() string {
	return fmt.Sprintf("camera with ID %s is already initialized", e.ID)
}
