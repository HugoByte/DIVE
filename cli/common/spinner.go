package common

import (
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
)

type diveSpinner struct {
	spinner *spinner.Spinner
}

func NewDiveSpinner() *diveSpinner {

	return &diveSpinner{
		spinner: spinner.New(spinner.CharSets[80], 100*time.Millisecond, spinner.WithWriter(os.Stdin)),
	}

}

func (ds *diveSpinner) SetSuffixMessage(message, color string) {
	ds.spinner.Suffix = message
	ds.SetColor(color)
}
func (ds *diveSpinner) SetPrefixMessage(message string) {
	ds.spinner.Prefix = message
}

func (ds *diveSpinner) SetColor(color string) {
	ds.spinner.Color(color)
}

func (ds *diveSpinner) Start(color string) {
	ds.SetColor(color)
	ds.spinner.Start()
}

func (ds *diveSpinner) StartWithMessage(message, color string) {

	ds.SetSuffixMessage(fmt.Sprint(" ", message), color)
	ds.Start(color)
}

func (ds *diveSpinner) Stop() {
	ds.spinner.Stop()
}

func (ds *diveSpinner) StopWithMessage(message string) {
	ds.setFinalMessage(message)
	ds.Stop()
}

func (ds *diveSpinner) setFinalMessage(message string) {
	ds.spinner.FinalMSG = fmt.Sprint(" ", message)
}
