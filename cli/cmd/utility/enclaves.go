package utility

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/hugobyte/dive/cli/common"
	"github.com/spf13/cobra"
)

var EnclavesCmd = common.NewDiveCommandBuilder().
	SetUse("enclaves").
	SetShort("Prints The Info About Enclaves").
	SetLong("Info About Enclaves Name,UUID,Short UUID,Status,Created Time.").
	SetRun(enclave).Build()

func enclave(cmd *cobra.Command, args []string) {
	cliContext := common.GetCliWithKurtosisContext(common.EnclaveName)

	err := common.ValidateArgs(args)
	if err != nil {
		cliContext.Fatal(err)
	}

	enclaves, err := cliContext.Context().GetEnclaves()
	if err != nil {
		cliContext.Fatal(err)
	}

	if len(enclaves) == 0 {
		cliContext.Logger().SetOutputToStdout()
		cliContext.Logger().Info("No Enclaves Running")
		cliContext.Context().Exit(0)

	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	// Print headers
	fmt.Fprintln(w, "Name\tUUID\tShort UUID\tCreated Time\tStatus")

	// Print each row
	for _, enclave := range enclaves {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", enclave.Name, enclave.Uuid, enclave.ShortUuid, enclave.CreatedTime, enclave.Status)
	}

	// Flush the buffer
	w.Flush()

	cliContext.StopSpinnerIfNotVerbose("Clean Completed", common.DiveLogs)
}
