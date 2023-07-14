/*
Copyright © 2023 Hugobyte AI Labs<hello@hugobyte.com>
*/
package main

import (
	"github.com/hugobyte/dive/commands"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	commands.Execute()

}
