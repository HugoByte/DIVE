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

func (s *diveSpinner) Start(message string)      {}
func (s *diveSpinner) Stop()                     {}
func (s *diveSpinner) SetMessage(message string) {}
func (s *diveSpinner) SetColor(color string)     {}
