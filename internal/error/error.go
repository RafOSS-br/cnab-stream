package internal_error

import (
	"fmt"
	"reflect"
)

type IError struct {
	Msg string
	Err error
}

func (e *IError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Msg, e.Err)
	}
	return e.Msg
}

func (e *IError) Unwrap() error {
	return e.Err
}

func NewError(msg string, err error) *IError {
	return &IError{Msg: msg, Err: err}
}

func IsError(errExpected, err error) bool {
	if errExpected == nil {
		return err == nil
	}
	expectedType := reflect.TypeOf(errExpected)
	return reflect.TypeOf(err) == expectedType
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
			return fmt.Errorf("Encapsulator is nil")
		}
		return e.CreateError(inner)
	}
}
