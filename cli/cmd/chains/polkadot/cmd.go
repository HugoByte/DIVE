package polkadot

import (
	"fmt"

	"github.com/hugobyte/dive-core/cli/common"
	"github.com/spf13/cobra"
)

var (
	configFilePath string
	paraChain      string
	network        string
	paraNodes      string
	relayNodes     string
	explorer       bool
	metrics        bool
)

const (
	runPolkadotFunctionName       = "run_polkadot_setup"
	runPolkadotRelayLocal         = "start_relay_chains_local"
	runPolkadotRelayTestnetMainet = "start_test_main_net_relay_nodes"
)

var PolkadotCmd = common.NewDiveCommandBuilder().
	SetUse("polkadot").
	SetShort("Build, initialize and start a polkadot node").
	SetLong("The command starts the polkadot relay chain and polkadot parachain if -p flag is given").
	SetRun(polkadot).
	AddStringFlagWithShortHand(&paraChain, "parachain", "p", "", "specify the parachain to spwan parachain node").
	AddStringFlagWithShortHand(&network, "network", "n", "", "specify the which network to run. local/testnet/mainnet. default will be local.").
	AddStringFlag(&paraNodes, "para-nodes", "", "specify the nodes for parachain, default will be '[full, collator]'").
	AddStringFlag(&relayNodes, "relay-nodes", "", "specify the nodes for relaychain, default will be '[full, validator]'").
	AddStringFlagWithShortHand(&configFilePath, "config", "c", "", "path to custom config json file to stat polakdot relaychain and parachain nodes.").
	AddBoolFlag(&explorer, "explorer", false, "specify the bool flag if you wanna start polakdot js explorer service").
	AddBoolFlag(&metrics, "metrics", false, "specify the bool flag if you wanna start prometheus metrics service").
	Build()

func polkadot(cmd *cobra.Command, args []string) {
	cliContext := common.GetCliWithKurtosisContext()

	err := common.ValidateArgs(args)

	if err != nil {
		cliContext.Fatalf("Error %s. %s", err, cmd.UsageString())
	}

	cliContext.StartSpinnerIfNotVerbose("Starting Polkadot Node", common.DiveLogs)

	response, err := RunPolkadot(cliContext)

	if err != nil {
		cliContext.Fatal(err)
	}

	enclaves, err := cliContext.Context().GetEnclaves()
	if err != nil {
		cliContext.Fatal(err)
	}

	var ShortUuid string
	for _, enclave := range enclaves {
		if enclave.Name == common.EnclaveName {
			ShortUuid = enclave.ShortUuid
		}
	}

	serviceFileName := fmt.Sprintf(common.ServiceFilePath, common.EnclaveName, ShortUuid)

	fmt.Print(response.Dive)
	for serviceName := range response.Dive {
		err = common.WriteServiceResponseData(response.Dive[serviceName].ServiceName, *response.Dive[serviceName], cliContext, serviceFileName)

		if err != nil {
			cliContext.Fatal(err)
		}
	}
	stopMessage := fmt.Sprintf("Polkadot Node Started. Please find service details in current working directory(%s)\n", serviceFileName)
	cliContext.StopSpinnerIfNotVerbose(stopMessage, common.DiveLogs)
}
