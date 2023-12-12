package neutron

import (
	"fmt"

	"github.com/hugobyte/dive-core/cli/common"
	"github.com/spf13/cobra"
)

var (
	configFilePath string
)

const (
	runNeutronNodeWithDefaultConfigFunctionName = "start_node_service"
)

var NeutronCmd = common.NewDiveCommandBuilder().
	SetUse("neutron").
	SetShort("Build, initialize and start a neutron node").
	SetLong("The command starts the neutron network and allows node in executing contracts").
	SetRun(neutron).
	AddStringFlagWithShortHand(&configFilePath, "config", "c", "", "path to custom config json file to start archway node ").
	Build()

func neutron(cmd *cobra.Command, args []string) {

	cliContext := common.GetCliWithKurtosisContext()

	err := common.ValidateArgs(args)
	if err != nil {
		cliContext.Fatalf("Error %s. %s", err, cmd.UsageString())
	}

	cliContext.Spinner().StartWithMessage("Starting Neutron Node", "green")

	response, err := RunNeutron(cliContext)

	if err != nil {
		cliContext.Fatal(err)
	}
	serviceFileName := fmt.Sprintf(common.ServiceFilePath, common.EnclaveName)

	err = common.WriteServiceResponseData(response.ServiceName, *response, cliContext, serviceFileName)
	if err != nil {
		cliContext.Fatal(err)
	}

	cliContext.Spinner().StopWithMessage("Neutron Node Started. Please find service details in current working directory(services.json)")
}
