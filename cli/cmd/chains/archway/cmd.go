package archway

import (
	"fmt"

	"github.com/hugobyte/dive-core/cli/common"
	"github.com/spf13/cobra"
)

var (
	configFilePath string
)

const (
	constructServiceConfigFunctionName          = "get_service_config"
	runArchwayNodeWithCustomServiceFunctionName = "start_cosmos_node"
	runArchwayNodeWithDefaultConfigFunctionName = "start_node_service"
)

var ArchwayCmd = common.NewDiveCommandBuilder().
	SetUse("archway").
	SetShort("Build, initialize and start a archway node").
	SetLong("The command starts the archway network and allows node in executing contracts").
	SetRun(archway).
	AddStringFlagWithShortHand(&configFilePath, "config", "c", "", "path to custom config json file to start archway node ").
	Build()

func archway(cmd *cobra.Command, args []string) {

	cliContext := common.GetCliWithKurtosisContext(common.EnclaveName)

	err := common.ValidateArgs(args)
	if err != nil {
		cliContext.Fatalf("Error %s. %s", err, cmd.UsageString())
	}

	cliContext.StartSpinnerIfNotVerbose("Starting Archway Node", common.DiveLogs)

	response, err := RunArchway(cliContext)

	if err != nil {
		cliContext.Fatal(err)
	}

	shortUuid, err := cliContext.Context().GetShortUuid(common.EnclaveName)
	if err != nil {
		cliContext.Fatal(err)
	}

	serviceFileName := fmt.Sprintf(common.ServiceFilePath, common.EnclaveName, shortUuid)

	err = common.WriteServiceResponseData(response.ServiceName, *response, cliContext, serviceFileName)
	if err != nil {
		cliContext.Fatal(err)

	}

	stopMessage := fmt.Sprintf("Archway Node Started. Please find service details in current working directory(%s)\n", serviceFileName)
	cliContext.StopSpinnerIfNotVerbose(stopMessage, common.DiveLogs)

}
