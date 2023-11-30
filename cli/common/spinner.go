package common

import (
	"os"
	"time"

	"github.com/briandowns/spinner"
)

type diveSpinner struct {
	spinner *spinner.Spinner
}

func NewDiveSpinner() *diveSpinner {

	spinner := spinner.New(spinner.CharSets[80], 100*time.Millisecond, spinner.WithWriter(os.Stdin))

	return &diveSpinner{spinner: spinner}

}

func (ds *diveSpinner) SetMessage(message string, color string) {
	panic("not implemented") // TODO: Implement
}

func (ds *diveSpinner) SetColor(color string) {
	panic("not implemented") // TODO: Implement
}

func (ds *diveSpinner) Start(message string) {
	panic("not implemented") // TODO: Implement
}

func (ds *diveSpinner) Stop(message string) {
	panic("not implemented") // TODO: Implement
}
