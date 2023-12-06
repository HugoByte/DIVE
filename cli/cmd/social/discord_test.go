package social

import (
	"bytes"
	"os"
	"os/exec"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestDiscordWithInvalidArgs(t *testing.T) {

	if os.Getenv("FLAG") == "1" {
		cmd := &cobra.Command{}
		args := []string{"arg1", "arg2"}
		discord(cmd, args)
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestDiscord")
	cmd.Env = append(os.Environ(), "FLAG=1")
	err := cmd.Run()
	e, ok := err.(*exec.ExitError)
	expectedErrorString := "exit status 1"
	assert.Equal(t, true, ok)
	assert.Equal(t, expectedErrorString, e.Error())
}

func TestDiscord(t *testing.T) {
	expectedLog := "Redirecting to DIVE discord channel..."
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	defer func() {
		w.Close()
		os.Stdout = old
	}()

	cmd := &cobra.Command{}
	args := []string{}

	discord(cmd, args)

	w.Close()
	var capturedLog bytes.Buffer
	_, _ = capturedLog.ReadFrom(r)

	assert.Contains(t, capturedLog.String(), expectedLog)

}
