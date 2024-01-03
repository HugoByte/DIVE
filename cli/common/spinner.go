package common

import (
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
)

// The diveSpinner type is a struct that contains a pointer to a spinner.Spinner object.
// spinner - The `spinner` is a pointer to an instance of the `spinner.Spinner`
//
// struct.
type diveSpinner struct {
	spinner *spinner.Spinner
}

// The function NewDiveSpinner returns a new instance of the diveSpinner struct.
func NewDiveSpinner() *diveSpinner {

	return &diveSpinner{
		spinner: spinner.New(spinner.CharSets[80], 100*time.Millisecond, spinner.WithWriter(os.Stdin)),
	}

}

// The `SetSuffixMessage` function is setting the suffix message of the spinner. It takes two
// parameters: `message` and `color`.
func (ds *diveSpinner) SetSuffixMessage(message, color string) {
	ds.spinner.Suffix = message
	ds.SetColor(color)
}

// The `SetPrefixMessage` function is setting the prefix message of the spinner. It takes a `message`
// parameter and assigns it to the `Prefix` field of the `spinner` object in the `diveSpinner` struct.
// This prefix message will be displayed before the spinner animation.

func (ds *diveSpinner) SetPrefixMessage(message string) {
	ds.spinner.Prefix = message
}

// The `SetColor` function is setting the color of the spinner. It takes a `color` parameter and
// assigns it to the `Color` field of the `spinner` object in the `diveSpinner` struct. This color will
// be used to display the spinner animation.
func (ds *diveSpinner) SetColor(color string) {
	ds.spinner.Color(color)
}

// The `Start` function is a method of the `diveSpinner` struct. It takes a `color` parameter and
// starts the spinner animation with the specified color.
func (ds *diveSpinner) Start(color string) {
	ds.SetColor(color)
	ds.spinner.Start()
}

// The `StartWithMessage` function is a method of the `diveSpinner` struct. It takes two parameters:
// `message` and `color`.
func (ds *diveSpinner) StartWithMessage(message, color string) {

	ds.SetSuffixMessage(fmt.Sprint(" ", message), color)
	ds.Start(color)
}

// The `Stop()` function is a method of the `diveSpinner` struct. It checks if the spinner is currently
// active and if so, it stops the spinner animation by calling the `Stop()` method of the `spinner`
// object in the `diveSpinner` struct.
func (ds *diveSpinner) Stop() {
	if ds.spinner.Active() {
		ds.spinner.Stop()
	}
}

// The `StopWithMessage` function is a method of the `diveSpinner` struct. It takes a `message`
// parameter and performs the following actions:
func (ds *diveSpinner) StopWithMessage(message string) {
	ds.setFinalMessage(message)
	ds.Stop()
}

// The `setFinalMessage` function is a method of the `diveSpinner` struct. It takes a `message`
// parameter and sets the final message of the spinner. This final message will be displayed after the
// spinner animation stops. The function assigns the `message` parameter to the `FinalMSG` field of the
// `spinner` object in the `diveSpinner` struct.
func (ds *diveSpinner) setFinalMessage(message string) {
	ds.spinner.FinalMSG = fmt.Sprint(" ", message)
}
