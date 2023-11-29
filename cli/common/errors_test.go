package common

import (
	"testing"
)

func TestNewBase(t *testing.T) {

	err := NewBase(UnknownError, "UnknownError")

	if c := CodeOf(err); c != UnknownError {
		t.Error("Code of NewBase() isn't codeUnknownError")
	}
}

func TestCodedError(t *testing.T) {
	msg := "Unknown Error Message"
	err := Errorc(UnknownError, msg)
	if c := CodeOf(err); c != UnknownError {
		t.Errorf("Expected code of error to be %d, got %d", UnknownError, c)
	}
}

func TestWrap(t *testing.T) {

	err := Errorc(InvalidCommandError, "Usage of Invalid Command")

	err2 := WrapMessageToError(err, "Invalid Usage")

	if c := CodeOf(err2); c != InvalidCommandError {
		t.Error("Code of Wrap() isn't codeInvalidCommandError")
	}

	err3 := WrapCodeToError(WithCode(err, UnknownError), InvalidCommandError, "InvalidUsage")

	if c := CodeOf(err3); c != InvalidCommandError {
		t.Errorf("Code of WithCode() isn't %d, got %d", InvalidCommandError, c)
	}

}

func TestIs(t *testing.T) {

	error1 := Errorc(InvalidCommandError, "Usage of Invalid Command")

	error2 := WrapMessageToError(error1, "Invalid Usage")

	if Is(error1, error2) {
		t.Error("error1 is not originated from error2")
	}

	if !Is(error2, error1) {
		t.Errorf("error2 is originated from error1")
	}

	error3 := WrapCodeToError(WithCode(error2, UnknownError), InvalidCommandError, "Invalid Usage")

	if !Is(error3, error1) {
		t.Errorf("error3 is originated from error1")
	}
}
