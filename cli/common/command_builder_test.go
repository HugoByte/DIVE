package common

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestDiveCommandBuilder(t *testing.T) {
	// Initialize a new DiveCommandBuilder
	builder := NewDiveCommandBuilder()

	// Set command attributes
	builder.SetUse("testCommand").
		SetShort("Short description").
		SetLong("Long description").
		SetRun(func(cmd *cobra.Command, args []string) error {
			cmd.Println("test function")

			return nil
		})

	// Add flags
	boolFlag := false
	builder.AddBoolFlag(&boolFlag, "boolFlag", false, "Bool flag usage")

	stringFlag := ""
	builder.AddStringFlag(&stringFlag, "stringFlag", "default", "String flag usage")

	// Build the command
	cmd := builder.Build()

	// Run tests
	assert.Equal(t, "testCommand", cmd.Use)
	assert.Equal(t, "Short description", cmd.Short)
	assert.Equal(t, "Long description", cmd.Long)

	// Test flag addition
	assert.NotNil(t, cmd.Flags().Lookup("boolFlag"))
	assert.NotNil(t, cmd.Flags().Lookup("stringFlag"))

	// Test flag default values
	assert.False(t, boolFlag)
	assert.Equal(t, "default", stringFlag)

	// Test the Run function

	var outputBuffer bytes.Buffer
	cmd.SetOut(&outputBuffer)

	// Run the command
	cmd.Execute()

	// Check the output
	expectedOutput := "test function\n"
	actualOutput := outputBuffer.String()

	assert.Equal(t, expectedOutput, actualOutput)
}
