/*
Copyright © 2023 Hugobyte AI Labs <hello@hugobyte.com>
*/
package bridge

import (
	"fmt"

	"github.com/spf13/cobra"
)

// bridgeCmd represents the bridge command
var BridgeCmd = &cobra.Command{
	Use:   "bridge",
	Short: "Command for cross chain communication between two different chains",
	Long: `To connect two different chains using any of the supported cross chain communication protocols. 
This will create an relay to connect two different chains and pass any messages between them.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("bridge called")
	},
}
