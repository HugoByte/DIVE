/*
Copyright © 2023 Hugobyte AI Labs<hello@hugobyte.com>
*/
package chain

import (
	"fmt"

	"github.com/spf13/cobra"
)

// chainCmd represents the chain command
var ChainCmd = &cobra.Command{
	Use:   "chain",
	Short: "Build, initialize and start a given blockchain node.",
	Long: `This command will spin up given chains. It will only run a node and developers can upload and test their smart contracts here. Incase if you want to use cross chain communication, please refer to 'bridge' command.`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

// chainCmd represents the chain command
var iconCmd = &cobra.Command{
	Use:   "icon",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hello")
	},
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// chainCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// chainCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	ChainCmd.AddCommand(iconCmd)
}
