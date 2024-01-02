package icon

import (
	"fmt"

	"github.com/hugobyte/dive-core/cli/common"
	"github.com/spf13/cobra"
)

const DefaultIconGenesisFile = "../../static-files/config/genesis-icon-0.zip"

var (
	serviceName    = ""
	ksPath         = ""
	ksPassword     = ""
	networkID      = ""
	nodeEndpoint   = ""
	genesis        = ""
	configFilePath = ""
)

var IconCmd = common.NewDiveCommandBuilder().
	SetUse("icon").
	SetShort("Build, initialize and start a icon node.").
	SetLong(`The command starts an Icon node, initiating the process of setting up and launching a local Icon network.
It establishes a connection to the Icon network and allows the node in executing smart contracts and maintaining the decentralized ledger.`).
	AddCommand(IconDecentralizeCmd).
	AddStringFlagWithShortHand(&genesis, "genesis", "g", "", "path to custom genesis file").
	AddStringFlagWithShortHand(&configFilePath, "config", "c", "", "path to custom config json file").
	AddBoolFlagP("decentralization", "d", false, "decentralize Icon Node").
	SetRun(icon).
	Build()

var IconDecentralizeCmd = common.NewDiveCommandBuilder().
	SetUse("decentralize").
	SetShort("Decentralize already running Icon Node").
	SetLong(`Decentralize Icon Node is necessary if you want to connect your local icon node to BTP network`).
	AddStringFlagWithShortHand(&serviceName, "serviceName", "s", "", "service name").
	AddStringFlagWithShortHand(&nodeEndpoint, "nodeEndpoint", "e", "", "endpoint address").
	AddStringFlagWithShortHand(&ksPath, "keystorePath", "k", "", "keystore path").
	AddStringFlagWithShortHand(&ksPassword, "keyPassword", "p", "", "keypassword").
	AddStringFlagWithShortHand(&networkID, "nid", "n", "", "NetworkId of Icon Node").
	MarkFlagsAsRequired([]string{"serviceName", "nodeEndpoint", "keystorePath", "keyPassword", "nid"}).
	SetRun(iconDecentralization).
	Build()

func icon(cmd *cobra.Command, args []string) {

	cliContext := common.GetCliWithKurtosisContext()

	err := common.ValidateArgs(args)
	if err != nil {
		cliContext.Fatalf("Error %s. %s", err, cmd.UsageString())
	}

	decentralization, err := cmd.Flags().GetBool("decentralization")
	if err != nil {
		cliContext.Fatal(common.WrapMessageToError(common.ErrInvalidCommand, err.Error()))
	}

	var response = &common.DiveServiceResponse{}

	cliContext.Spinner().StartWithMessage("Starting Icon Node", "green")
	if decentralization {
		response, err = RunIconNode(cliContext)

		if err != nil {
			cliContext.Fatal(err)
		}
		params := GetDecentralizeParams(response.ServiceName, response.PrivateEndpoint, response.KeystorePath, response.KeyPassword, response.NetworkId)

		err = RunDecentralization(cliContext, params)

		if err != nil {
			cliContext.Fatal(err)
		}

	} else {
		response, err = RunIconNode(cliContext)

		if err != nil {
			cliContext.Fatal(err)
		}

	}

	shortUuid, err := cliContext.Context().GetShortUuid(common.EnclaveName)
	if err != nil {
		cliContext.Fatal(err)
	}

	serviceFileName := fmt.Sprintf(common.ServiceFilePath, common.EnclaveName, shortUuid)

	err = common.WriteServiceResponseData(response.ServiceName, *response, cliContext, serviceFileName)
	if err != nil {
		cliContext.Error(err)
		cliContext.Context().Exit(1)

	}

	stopMessage := fmt.Sprintf("Icon Node Started. Please find service details in current working directory(%s)\n", serviceFileName)
	cliContext.Spinner().StopWithMessage(stopMessage)

}

func iconDecentralization(cmd *cobra.Command, args []string) {

	cliContext := common.GetCliWithKurtosisContext()

	err := common.ValidateArgs(args)

	if err != nil {
		cliContext.Fatalf("Error %s. %s", err, cmd.UsageString())
	}

	cliContext.Spinner().StartWithMessage("Starting Icon Node Decentralization", "green")

	params := GetDecentralizeParams(serviceName, nodeEndpoint, ksPath, ksPassword, networkID)

	err = RunDecentralization(cliContext, params)

	if err != nil {
		cliContext.Fatal(err)

	}

	cliContext.Spinner().StopWithMessage("Decentralization Completed")
}
