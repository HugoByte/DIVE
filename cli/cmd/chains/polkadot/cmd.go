package polkadot

import (
	"fmt"
	"strings"

	"github.com/hugobyte/dive/cli/cmd/chains/utils"
	"github.com/hugobyte/dive/cli/common"
	"github.com/spf13/cobra"
)

var (
	configFilePath string
	paraChain      []string
	network        string
	noRelay        bool
	explorer       bool
	metrics        bool
)

const (
	runPolkadotFunctionName             = "run_polkadot_setup"
	runPolkadotRelayLocal               = "start_relay_chains_local"
	runPolkadotRelayTestnetMainet       = "start_test_main_net_relay_nodes"
	runUploadFiles                      = "upload_files"
	runPolkadotParaLocalFunctionName    = "start_nodes"
	runPolkadotParaTestMainFunctionName = "run_testnet_mainnet"
	runPolkadotExplorer                 = "run_pokadot_js_app"
	runPolkadotPrometheus               = "launch_prometheus"
	runPolkadotGrafana                  = "launch_grafana"
)

var (
	polkadotParachains = []string{"acala", "ajuna", "bifrost", "centrifuge", "clover", "frequency", "kilt", "kylin", "litentry", "manta", "moonbeam", "moonsama", "nodle", "parallel", "pendulum", "phala", "polkadex", "subsocial", "zeitgeist", "robonomics"}
)

var PolkadotCmd = common.NewDiveCommandBuilder().
	SetUse("polkadot").
	SetShort("Build, initialize and start a Polkadot node").
	SetLong("The command starts the polkadot relaychain and polkadot parachain if -p flag is given").
	SetRun(polkadot).
	AddStringSliceFlagWithShortHand(&paraChain, "parachain", "p", []string{}, "specify the list of parachains to spawn parachain node").
	AddStringFlagWithShortHand(&network, "network", "n", "", "specify the network to run (local/testnet/mainnet). Default will be local.").
	AddBoolFlag(&noRelay, "no-relay", false, "specify the bool flag to run parachain only (only for testnet and mainnet)").
	AddStringFlagWithShortHand(&configFilePath, "config", "c", "", "path to custom config json file to start polkadot relaychain and parachain nodes.").
	AddBoolFlag(&explorer, "explorer", false, "specify the bool flag if you want to start polkadot js explorer service").
	AddBoolFlag(&metrics, "metrics", false, "specify the bool flag if you want to start prometheus and grafana metrics service").
	Build()

func polkadot(cmd *cobra.Command, args []string) {
	cliContext := common.GetCliWithKurtosisContext(common.EnclaveName)

	err := common.ValidateArgs(args)
	if err != nil {
		cliContext.Fatalf("Error %s. %s", err, cmd.UsageString())
	}
	cliContext.StartSpinnerIfNotVerbose("Starting Polkadot Node", common.DiveLogs)

	response, err := RunPolkadot(cliContext)
	if err != nil {
		if strings.Contains(err.Error(), "already running") {
			cliContext.Error(err)
			cliContext.Context().Exit(0)
		} else {
			cliContext.Fatal(err)
		}
	}

	shortUuid, err := cliContext.Context().GetShortUuid(common.EnclaveName)
	if err != nil {
		cliContext.Fatal(err)
	}

	serviceFileName := fmt.Sprintf(common.ServiceFilePath, common.EnclaveName, shortUuid)

	for serviceName := range response.Dive {
		err = common.WriteServiceResponseData(response.Dive[serviceName].ServiceName, *response.Dive[serviceName], cliContext, serviceFileName)
		if err != nil {
			cliContext.Fatal(err)
		}
	}

	stopMessage, err := utils.GetStopMessage(cliContext, configFilePath, "Polkadot", paraChain)
	if err != nil {
		cliContext.Fatal(err)
	}
	stopMessage = stopMessage + fmt.Sprintf("Please find the service details in the output folder present in current working directory - (output/%s/%s)\n", common.EnclaveName, serviceFileName)
	cliContext.StopSpinnerIfNotVerbose(stopMessage, common.DiveLogs)
}
