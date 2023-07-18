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
	Long:  `The command opens the official HugoByte Twitter homepage. It launches a web browser and directs users to the designated Twitter profile of HugoByte, providing
access to the latest updates, announcements, news, and insights shared by the official HugoByte Twitter account. Users can stay informed about HugoByte's activities, engage with the community, and follow our social media presence directly from the Twitter homepage.`,
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("Redirecting to twitter...")
		if err := common.OpenFile(twitterURL); err != nil {
			logrus.Errorf("Failed to open HugoByte twitter with error %v", err)
		}
	},
}
