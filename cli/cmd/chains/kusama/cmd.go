package kusama

import (
	"fmt"
	"strings"

	"github.com/hugobyte/dive-core/cli/common"
	"github.com/spf13/cobra"
)

var (
	configFilePath string
	paraChain      string
	network        string
	noRelay        bool
	explorer       bool
	metrics        bool
)

const (
	runKusamaFunctionName             = "run_polkadot_setup"
	runKusamaRelayLocal               = "start_relay_chains_local"
	runKusamaRelayTestnetMainet       = "start_test_main_net_relay_nodes"
	runUploadFiles                    = "upload_files"
	runKusamaParaLocalFunctionName    = "start_nodes"
	runKusamaParaTestMainFunctionName = "run_testnet_mainnet"
	runKusamaExplorer                 = "run_pokadot_js_app"
	runKusamaPrometheus               = "launch_prometheus"
	runKusamaGrafana                  = "launch_grafana"
)

var KusamaCmd = common.NewDiveCommandBuilder().
	SetUse("kusama").
	SetShort("Build, initialize and start a Kusama node").
	SetLong("The command starts the kusama relay chain and kusama parachain if -p flag is given").
	SetRun(kusama).
	AddStringFlagWithShortHand(&paraChain, "parachain", "p", "", "specify the parachain to spawn parachain node").
	AddStringFlagWithShortHand(&network, "network", "n", "", "specify the network to run (local/testnet/mainnet). Default will be local.").
	AddBoolFlag(&noRelay, "no-relay", false, "specify the bool flag to run parachain only (only for testnet and mainnet)").
	AddStringFlagWithShortHand(&configFilePath, "config", "c", "", "path to custom config json file to start kusama relaychain and parachain nodes.").
	AddBoolFlag(&explorer, "explorer", false, "specify the bool flag if you want to start polkadot js explorer service").
	AddBoolFlag(&metrics, "metrics", false, "specify the bool flag if you want to start prometheus metrics service").
	Build()

func kusama(cmd *cobra.Command, args []string) {
	cliContext := common.GetCliWithKurtosisContext(common.EnclaveName)

	err := common.ValidateArgs(args)
	if err != nil {
		cliContext.Fatalf("Error %s. %s", err, cmd.UsageString())
	}
	cliContext.StartSpinnerIfNotVerbose("Starting Kusama Node", common.DiveLogs)

	response, err := RunKusama(cliContext)
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
	stopMessage := fmt.Sprintf("Kusama Node Started. Please find service details in current working directory(%s)\n", serviceFileName)
	cliContext.StopSpinnerIfNotVerbose(stopMessage, common.DiveLogs)

}
