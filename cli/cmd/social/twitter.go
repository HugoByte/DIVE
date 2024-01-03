package social

import (
	"fmt"

	"github.com/hugobyte/dive-core/cli/common"
	"github.com/spf13/cobra"
)

const twitterURL = "https://twitter.com/hugobyte"

var TwitterCmd = common.NewDiveCommandBuilder().
	SetUse("twitter").
	SetShort("Opens official HugoByte twitter home page").
	SetLong(`The command opens the official HugoByte Twitter homepage. It launches a web browser and directs users to the designated Twitter profile of HugoByte, providing access to the latest updates, announcements, news, and insights
shared by the official HugoByte Twitter account. Users can stay informed about HugoByte's activities, engage with the 
community, and follow our social media presence directly from the Twitter homepage.`,
	).
	SetRun(twitter).Build()

func twitter(cmd *cobra.Command, args []string) {

	cli := common.GetCli()

	cli.Logger().SetOutputToStdout()
	err := common.ValidateArgs(args)

	if err != nil {
		cli.Logger().Error(common.CodeOf(err), common.Errorc(common.CodeOf(err), fmt.Sprintf("error %s \n %s ", err, cmd.UsageString())).Error())
	}

	cli.Logger().Info("Redirecting to twitter...")

	if err := common.OpenFile(twitterURL); err != nil {
		cli.Logger().Fatalf(common.CodeOf(err), "Failed to open HugoByte twitter with error %v", err)
	}
}
