package common

import "github.com/spf13/cobra"

type diveCommandBuilder struct {
	cmd *cobra.Command
}

func NewDiveCommandBuilder() *diveCommandBuilder {
	return &diveCommandBuilder{
		cmd: &cobra.Command{},
	}
}

// AddCommand adds a subcommand to the command.
func (d *diveCommandBuilder) AddCommand(cmd *cobra.Command) CommandBuilder {
	panic("not implemented") // TODO: Implement
}

// Add Persistent Bool Flag
func (d *diveCommandBuilder) AddBoolPersistentFlag(p *bool, name string, value bool, usage string) CommandBuilder {
	panic("not implemented") // TODO: Implement
}

// Add Persistent Bool Flag with Short hand
func (d *diveCommandBuilder) AddBoolPersistentFlagWithShortHand(p *bool, name string, value bool, usage string, shorthand string) CommandBuilder {
	panic("not implemented") // TODO: Implement
}

// Add Persistent String Flag
func (d *diveCommandBuilder) AddStringPersistentFlag(p *string, name string, value string, usage string) CommandBuilder {
	panic("not implemented") // TODO: Implement
}

// Add Persistent String Flag with Short hand
func (d *diveCommandBuilder) AddStringPersistentFlagWithShortHand(p *string, name string, shorthand string, value string, usage string) CommandBuilder {
	panic("not implemented") // TODO: Implement
}

// Add StringFlag adds a string flag to the command that persists
func (d *diveCommandBuilder) AddStringFlag(name string, value string, usage string) CommandBuilder {
	panic("not implemented") // TODO: Implement
}

// Add StringFlag adds a string flag to the command that persists with short hand
func (d *diveCommandBuilder) AddStringFlagWithShortHand(p *string, name string, shorthand string, value string, usage string) CommandBuilder {
	panic("not implemented") // TODO: Implement
}

// Add BooFlag adds a boolean flag to the command that persists
func (d *diveCommandBuilder) AddBoolFlag(name string, value bool, usage string) CommandBuilder {
	panic("not implemented") // TODO: Implement
}

func (d *diveCommandBuilder) AddBoolFlagWithShortHand(name string, shorthand string, value bool, usage string) CommandBuilder {
	panic("not implemented") // TODO: Implement
}

// Build constructs and returns the Cobra command.
func (d *diveCommandBuilder) Build() *cobra.Command {
	panic("not implemented") // TODO: Implement
}

// SetUse sets the Use field of the command.
func (d *diveCommandBuilder) SetUse(use string) CommandBuilder {
	panic("not implemented") // TODO: Implement
}

// SetShort sets the Short field of the command.
func (d *diveCommandBuilder) SetShort(short string) CommandBuilder {
	panic("not implemented") // TODO: Implement
}

// SetLong sets the Long field of the command.
func (d *diveCommandBuilder) SetLong(long string) CommandBuilder {
	panic("not implemented") // TODO: Implement
}

// SetRun sets the Run field of the command.
func (d *diveCommandBuilder) SetRun(run func(cmd *cobra.Command, args []string)) CommandBuilder {
	panic("not implemented") // TODO: Implement
}
