/*
Copyright © 2023 Hugobyte AI Labs<hello@hugobyte.com>
*/
package twitter

import (
	"github.com/hugobyte/dive/common"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const twitterURL = "https://twitter.com/hugobyte"

// twitterCmd redirects users to twitter home page
var TwitterCmd = &cobra.Command{
	Use:   "twitter",
	Short: "Opens official HugoByte twitter home page",
	Long:  `Redirects users to HugoByte twitter home page`,
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("Redirecting to twitter...")
		if err := common.OpenFile(twitterURL); err != nil {
			logrus.Errorf("Failed to open HugoByte twitter with error %v", err)
		}
	},
}
