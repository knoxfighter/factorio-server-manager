package bootstrap

import (
	"fmt"
	"strings"
)

// errors are basically strings, that implement the error interface
type ConstError string

// Implement error interface, so this is a valid error.
func (err ConstError) Error() string {
	return string(err)
}

// Implement newer Is method to check if wrapped error is the desired error.
func (err ConstError) Is(target error) bool {
	targetError := target.Error()
	errorString := string(err)
	return targetError == errorString || strings.HasPrefix(targetError, errorString+": ")
}

// Wrap suberror with this error. `Is` can be checked if wrapped errors is of type
func (err ConstError) Wrap(inner error) error {
	return wrapError{msg: string(err), err: inner}
}

// If an error is wrapped, we change the type to this
type wrapError struct {
	err error
	msg string
}

// Also implement Error interface, to use wrapErrors as error
func (err wrapError) Error() string {
	if err.err != nil {
		return fmt.Sprintf("%s: %v", err.msg, err.err)
	}
	return err.msg
}

// Make it possible to Unwrap a wrapped error again.
func (err wrapError) Unwrap() error {
	return err.err
}

// Implement newer Is method to check unwrapped error
func (err wrapError) Is(target error) bool {
	return ConstError(err.msg).Is(target)
}

type ErrorWithData struct {
	Msg  string
	Data interface{}
}

func (err ErrorWithData) Error() string {
	return fmt.Sprintf("%s: %v", err.Msg, err.Data)
}

// Implement newer Is method to check if wrapped error is the desired error.
func (err ErrorWithData) Is(target error) bool {
	targetError := target.Error()
	return targetError == err.Msg || strings.HasPrefix(targetError, err.Msg+": ")
}

// Wrap suberror with this error. `Is` can be checked if wrapped errors is of type
func (err ErrorWithData) Wrap(inner error) error {
	return wrapErrorWithData{msg: err.Msg, err: inner, data: err.Data}
}

// If an error is wrapped, we change the type to this
type wrapErrorWithData struct {
	msg  string
	err  error
	data interface{}
}

// Also implement Error interface, to use wrapErrors as error
func (err wrapErrorWithData) Error() string {
	if err.err != nil {
		return fmt.Sprintf("%s: %v: %v", err.msg, err.err, err.data)
	}
	return err.msg
}

// Make it possible to Unwrap a wrapped error again.
func (err wrapErrorWithData) Unwrap() error {
	return err.err
}

// Implement newer Is method to check unwrapped error
func (err wrapErrorWithData) Is(target error) bool {
	return ConstError(err.msg).Is(target)
}
