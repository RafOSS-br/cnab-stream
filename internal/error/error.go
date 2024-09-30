package internal_error

import (
	"errors"
	"fmt"
	"reflect"
)

type InternalError interface {
	error
	FirstClass() error
	Unwrap() []error
}
type IError struct {
	Err [2]error
}

func (e *IError) Error() string {
	if e.Err[0] != nil && e.Err[1] != nil {
		return fmt.Sprintf("%v: %v", e.Err[0], e.Err[1])
	}
	if e.Err[0] != nil {
		return e.Err[0].Error()
	}
	return e.Err[1].Error()
}

func (e *IError) Unwrap() []error {
	return e.Err[:]
}

func (e *IError) FirstClass() error {
	return e.Err[0]
}

func NewError(errMessage string, childOrInfo error) InternalError {
	return &IError{Err: [2]error{errors.New(errMessage), childOrInfo}}
}

func IsError(errExpected, err error) bool {
	if errExpected == nil {
		return err == nil
	}
	if errors.Is(err, errExpected) {
		return true
	}

	errType := reflect.TypeOf(errExpected)
	target := reflect.New(errType).Interface()

	return errors.As(err, &target)
}

func MatchError(errExpected error) func(error) bool {
	return func(err error) bool {
		return IsError(errExpected, err)
	}
}

type Encapsulator[T any] interface {
	CreateError(error) error
}

// NewCreator returns a function that creates an error using the provided Encapsulator
func NewCreator[T any](e Encapsulator[T]) func(error) error {
	return func(inner error) error {
		if e == nil {
			return nil
		}
		return e.CreateError(inner)
	}
}
