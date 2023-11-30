package common

import (
	"errors"
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

func TestWrapMessage(t *testing.T) {

	error1 := Errorc(InvalidCommandError, "Usage of Invalid Command")

	error2 := WrapMessageToError(error1, "Invalid Usage")

	if c := CodeOf(error2); c != InvalidCommandError {
		t.Error("Code of Wrap() isn't codeInvalidCommandError")
	}

	err3 := WithCode(error1, UnknownError)

	if c := CodeOf(err3); c != UnknownError {
		t.Errorf("with code doesn't change the Error code")
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

	if Is(error1, error3) {
		t.Errorf("error1 is originated from error3")
	}
}

func TestAsValue(t *testing.T) {
	var coder ErrorCoder
	if AsValue(&coder, Errorc(InvalidCommandError, "Test")) {
		if coder == nil {
			t.Error("Returned object is nil")
		}
		if c := coder.ErrorCode(); c != InvalidCommandError {
			t.Error("Fail to find ErrorCoder")
		}
	} else {
		t.Error("Fail to get ErrorCoder from result of Errorc()")
	}
}

func TestWithCode(t *testing.T) {

	errorCodes := []ErrorCode{
		InvalidCommandError, InvalidEnclaveNameError,
	}

	tests := map[string]error{

		"Errorc":       Errorc(UnsupportedOSError, "OS Error"),
		"NewBaseError": NewBase(FileError, "Invalid File"),
		"WrapMessage":  WrapMessageToError(Errorc(UnknownError, "Error Parsing Arguments"), "Invalid Usage"),
		"WrapCode":     WrapCodeToError(Errorc(UnsupportedOSError, "Unknown Platform"), UnsupportedOSError, "Error"),
		"WithCode":     WithCode(errors.New("test"), FileError),
	}

	for errName, err := range tests {
		t.Run(errName, func(t *testing.T) {
			for _, code := range errorCodes {
				error1 := WithCode(err, code)
				if c := CodeOf(error1); code != c {
					t.Errorf("Returned code=%d exp=%d", code, c)
				}
			}
		})
	}
}

func TestCodeOf(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		expectedCode ErrorCode
	}{
		{
			"test_new_base",
			NewBase(UnsupportedOSError, "Invalid OS"),
			UnsupportedOSError,
		},
		{
			"test_with_code",
			WithCode(errors.New("test"), FileError),
			FileError,
		},
		{
			"test_with_wrapwessage",
			WrapMessageToError(Errorc(UnknownError, "Error Parsing Arguments"), "Invalid Usage"),
			UnknownError,
		},
		{
			"test_with_wrap_code",
			WrapCodeToError(Errorc(UnsupportedOSError, "Unknown Platform"), UnsupportedOSError, "Error"),
			UnsupportedOSError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			code := CodeOf(test.err)

			if code != test.expectedCode {
				t.Errorf("Expected %d got %d", test.expectedCode, code)
			}
		})
	}
}
