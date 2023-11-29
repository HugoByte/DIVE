package common

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

type ErrorCode int

const (
	ErrorCodeGeneral ErrorCode = iota + 1000
)

const (
	UnknownError ErrorCode = ErrorCodeGeneral + iota
	InvalidEnclaveNameError
	UnsupportedOSError
	FileError
	InvalidCommandError
	InvalidEnclaveContextError
	InvalidEnclaveConfigError
)

var (
	ErrUnknown               = NewBase(UnknownError, "UnknownError")
	ErrInvalidEnclaveName    = NewBase(InvalidEnclaveNameError, "InvalidEnclaveName")
	ErrUnsupportedOS         = NewBase(UnsupportedOSError, "UnsupportedOS")
	ErrFile                  = NewBase(FileError, "FileError")
	ErrInvalidCommand        = NewBase(InvalidCommandError, "InvalidCommand")
	ErrInvalidEnclaveContext = NewBase(InvalidEnclaveContextError, "InvalidEnclaveContext")
	ErrInvalidEnclaveConfig  = NewBase(InvalidEnclaveConfigError, "InvalidEnclaveConfig")
)

func (c ErrorCode) New(msg string) error {
	return Errorc(c, msg)
}

func (c ErrorCode) Errorf(f string, args ...interface{}) error {
	return Errorcf(c, f, args...)
}

func (c ErrorCode) Wrap(e error, msg string) error {
	return WrapCodeToError(e, c, msg)
}

func (c ErrorCode) Wrapf(e error, f string, args ...interface{}) error {
	return WrapCodeToErrorf(e, c, f, args...)
}

func (c ErrorCode) Equals(e error) bool {
	if e == nil {
		return false
	}
	return CodeOf(e) == c
}

// Represent Error With Code and Message

type baseError struct {
	code ErrorCode
	msg  string
}

func (e *baseError) Error() string {
	return e.msg
}

func (e *baseError) ErrorCode() ErrorCode {
	return e.code
}

func (e *baseError) Format(f fmt.State, c rune) {
	switch c {
	case 'v', 's', 'q':
		fmt.Fprintf(f, "E%04d:%s", e.code, e.msg)
	}
}

func (e *baseError) Equals(err error) bool {
	if err == nil {
		return false
	}
	return CodeOf(err) == e.code
}

func NewBase(code ErrorCode, msg string) *baseError {
	return &baseError{code, msg}
}

/*
Associate an error code with a standard Go error.
Create a new error with a code and a message
*/

type codedError struct {
	code ErrorCode
	error
}

func (e *codedError) ErrorCode() ErrorCode {
	return e.code
}

func (e *codedError) Unwrap() error {
	return e.error
}

func Errorc(code ErrorCode, msg string) error {
	return &codedError{
		code:  code,
		error: errors.New(msg),
	}
}

func Errorcf(code ErrorCode, f string, args ...interface{}) error {
	return &codedError{
		code:  code,
		error: errors.Errorf(f, args...),
	}
}

func WithCode(err error, code ErrorCode) error {
	if _, ok := CoderOf(err); ok {
		return WrapCodeToError(err, code, err.Error())
	}
	return &codedError{
		code:  code,
		error: err,
	}
}

/*

Wrapping Existing Error with Additional Context like an error code and message.
Preserves Original error but adds the addtional Context to error messages
*/

type wrappedError struct {
	error
	code   ErrorCode
	origin error
}

func (e *wrappedError) Format(f fmt.State, c rune) {
	switch c {
	case 'v':
		if f.Flag('+') {
			fmt.Fprintf(f, "E%04d:%+v", e.code, e.error)
			fmt.Fprintf(f, "\nWrapping %+v", e.origin)
			return
		}
		fallthrough
	case 'q', 's':
		fmt.Fprintf(f, "E%04d:%s", e.code, e.error)
	}
}

func (e *wrappedError) Unwrap() error {
	return e.origin
}

func (e *wrappedError) ErrorCode() ErrorCode {
	return e.code
}

func WrapCodeToError(e error, c ErrorCode, msg string) error {
	return &wrappedError{
		error:  errors.New(msg),
		code:   c,
		origin: e,
	}
}

func WrapCodeToErrorf(e error, c ErrorCode, f string, args ...interface{}) error {
	return &wrappedError{
		error:  errors.Errorf(f, args...),
		code:   c,
		origin: e,
	}
}

// To Add Messaage to an error without changing the error code.

type messageError struct {
	error
	origin error
}

func (e *messageError) Format(f fmt.State, c rune) {
	switch c {
	case 'v':
		if f.Flag('+') {
			fmt.Fprintf(f, "%+v", e.error)
			fmt.Fprintf(f, "\nWrapping %+v", e.origin)
			return
		}
		fallthrough
	case 's', 'q':
		fmt.Fprintf(f, "%s", e.error)
	}
}

func (e *messageError) Unwrap() error {
	return e.origin
}

func WrapMessageToError(e error, msg string) error {
	return &messageError{
		error:  errors.New(msg),
		origin: e,
	}
}

func WrapMessageToErrorf(e error, f string, args ...interface{}) error {
	return &messageError{
		error:  errors.Errorf(f, args...),
		origin: e,
	}
}

// Extract Code from Custom Error Messages
type ErrorCoder interface {
	error
	ErrorCode() ErrorCode
}

func CoderOf(e error) (ErrorCoder, bool) {
	var coder ErrorCoder
	if AsValue(&coder, e) {
		return coder, true
	}
	return nil, false
}

func CodeOf(e error) ErrorCode {
	if coder, ok := CoderOf(e); ok {
		return coder.ErrorCode()
	}
	return UnknownError
}

func AsValue(ptr interface{}, err error) bool {
	type causer interface {
		Cause() error
	}

	type unwrapper interface {
		Unwrap() error
	}

	value := reflect.ValueOf(ptr)
	if value.Kind() != reflect.Ptr {
		return false
	} else {
		value = value.Elem()
	}
	valueType := value.Type()

	for {
		errValue := reflect.ValueOf(err)
		if errValue.Type().AssignableTo(valueType) {
			value.Set(errValue)
			return true
		}
		if cause, ok := err.(causer); ok {
			err = cause.Cause()
		} else if unwrap, ok := err.(unwrapper); ok {
			err = unwrap.Unwrap()
		} else {
			return false
		}
	}
}

// Is checks whether err is caused by the target.
func Is(err, target error) bool {
	type causer interface {
		Cause() error
	}

	type unwrapper interface {
		Unwrap() error
	}

	for {
		if err == target {
			return true
		}
		if cause, ok := err.(causer); ok {
			err = cause.Cause()
		} else if unwrap, ok := err.(unwrapper); ok {
			err = unwrap.Unwrap()
		} else {
			return false
		}
	}
}
