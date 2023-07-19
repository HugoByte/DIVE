/*
Copyright © 2023 Hugobyte AI Labs<hello@hugobyte.com>
*/
package main

import (
	"github.com/hugobyte/dive/commands"
	"github.com/hugobyte/dive/common"
	"github.com/sirupsen/logrus"
)

func main() {

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	latestVersion := common.GetLatestVersion()
	if common.DiveVersion != latestVersion {
		logrus.Warnf("Update available '%s'. Get the latest version of our DIVE CLI for bug fixes, performance improvements, and new features.", latestVersion)
	}
	commands.Execute()

}
