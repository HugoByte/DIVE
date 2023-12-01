package relays

import (
	"github.com/hugobyte/dive-core/cli/common"
	"github.com/spf13/cobra"
)

var IbcRelayCmd = common.NewDiveCommandBuilder().
	SetUse("ibc").
	SetShort("Start connection between Cosmos based chainA and ChainB and initiate communication between them").
	SetLong(`This Command deploy , initialize the contracts and make it ready for ibc.
Along with that setup and starts the ibc relayer to establish communication between chains specified`).
	SetRun(ibcRelay).
	Build()

func ibcRelay(cmd *cobra.Command, args []string) {}
