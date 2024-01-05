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
func (dc *diveCommandBuilder) AddCommand(cmd *cobra.Command) CommandBuilder {
	dc.cmd.AddCommand(cmd)
	return dc
}

// Add Persistent Bool Flag
func (dc *diveCommandBuilder) AddBoolPersistentFlag(boolV *bool, name string, value bool, usage string) CommandBuilder {
	dc.cmd.PersistentFlags().BoolVar(boolV, name, value, usage)
	return dc
}

// Add Persistent Bool Flag with Short hand
func (dc *diveCommandBuilder) AddBoolPersistentFlagWithShortHand(boolV *bool, name string, value bool, usage string, shorthand string) CommandBuilder {

	dc.cmd.PersistentFlags().BoolVarP(boolV, name, shorthand, value, usage)
	return dc
}

// Add Persistent String Flag
func (dc *diveCommandBuilder) AddStringPersistentFlag(stringV *string, name string, value string, usage string) CommandBuilder {
	dc.cmd.PersistentFlags().StringVar(stringV, name, value, usage)
	return dc
}

// Add Persistent String Flag with Short hand
func (dc *diveCommandBuilder) AddStringPersistentFlagWithShortHand(stringV *string, name string, shorthand string, value string, usage string) CommandBuilder {
	dc.cmd.PersistentFlags().StringVarP(stringV, name, shorthand, value, usage)
	return dc
}

// Add StringFlag adds a string flag to the command that persists
func (dc *diveCommandBuilder) AddStringFlag(stringV *string, name string, value string, usage string) CommandBuilder {
	dc.cmd.Flags().StringVar(stringV, name, value, usage)
	return dc
}

// Add StringFlag adds a string flag to the command that persists with short hand
func (dc *diveCommandBuilder) AddStringFlagWithShortHand(stringV *string, name string, shorthand string, value string, usage string) CommandBuilder {
	dc.cmd.Flags().StringVarP(stringV, name, shorthand, value, usage)
	return dc
}

// Add StringSliceFlag adds a slice of string flag to the command that persists with short hand
func (dc *diveCommandBuilder) AddStringSliceFlagWithShortHand(stringV *[]string, name string, shorthand string, value []string, usage string) CommandBuilder {
	dc.cmd.Flags().StringSliceVarP(stringV, name, shorthand, value, usage)
	return dc
}

// Add BooFlag adds a boolean flag to the command that persists
func (dc *diveCommandBuilder) AddBoolFlag(boolV *bool, name string, value bool, usage string) CommandBuilder {
	dc.cmd.Flags().BoolVar(boolV, name, value, usage)
	return dc
}

func (dc *diveCommandBuilder) AddBoolFlagWithShortHand(boolV *bool, name string, shorthand string, value bool, usage string) CommandBuilder {
	dc.cmd.Flags().BoolVarP(boolV, name, shorthand, value, usage)
	return dc
}

// Build constructs and returns the Cobra command.
func (dc *diveCommandBuilder) Build() *cobra.Command {
	dc.cmd.CompletionOptions.DisableDefaultCmd = true
	dc.cmd.CompletionOptions.DisableNoDescFlag = true
	return dc.cmd
}

// SetUse sets the Use field of the command.
func (dc *diveCommandBuilder) SetUse(use string) CommandBuilder {
	dc.cmd.Use = use
	return dc
}

// SetShort sets the Short field of the command.
func (dc *diveCommandBuilder) SetShort(short string) CommandBuilder {
	dc.cmd.Short = short
	return dc
}

// SetLong sets the Long field of the command.
func (dc *diveCommandBuilder) SetLong(long string) CommandBuilder {
	dc.cmd.Long = long
	return dc
}

// SetRun sets the Run field of the command.
func (dc *diveCommandBuilder) SetRun(run func(cmd *cobra.Command, args []string)) CommandBuilder {

	dc.cmd.Run = run

	return dc
}
func (dc *diveCommandBuilder) ToggleHelpCommand(enable bool) CommandBuilder {

	dc.cmd.SetHelpCommand(&cobra.Command{Hidden: enable})
	return dc
}

func (dc *diveCommandBuilder) SetRunE(run func(cmd *cobra.Command, args []string) error) CommandBuilder {

	dc.cmd.RunE = run

	return dc
}

func (dc *diveCommandBuilder) MarkFlagsAsRequired(flags []string) CommandBuilder {

	dc.cmd.MarkFlagsRequiredTogether(flags...)
	return dc
}

func (dc *diveCommandBuilder) MarkFlagRequired(flag string) CommandBuilder {
	dc.cmd.MarkFlagRequired(flag)
	return dc
}

func (dc *diveCommandBuilder) AddBoolFlagP(name string, shorthand string, value bool, usage string) CommandBuilder {

	dc.cmd.Flags().BoolP(name, shorthand, value, usage)
	return dc
}
