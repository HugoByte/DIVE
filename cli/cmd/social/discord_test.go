package social

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestDiscordCommand(t *testing.T) {
	tests := []struct {
		name             string
		args             []string
		expectedError    *string // Use a pointer to string
		expectedExitCode int
	}{
		{
			name:             "ValidArgs",
			args:             []string{},
			expectedError:    nil,
			expectedExitCode: 0,
		},
		{
			name:             "InvalidArgs",
			args:             []string{"arg1", "arg2"},
			expectedError:    getStringPointer("Invalid Usage Of Command Arguments"),
			expectedExitCode: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			cmd := &cobra.Command{}

			err := discord(cmd, tt.args)

			if tt.expectedError != nil {
				assert.NotNil(t, err)
				assert.EqualError(t, err, *tt.expectedError)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func getStringPointer(s string) *string {
	return &s
}
