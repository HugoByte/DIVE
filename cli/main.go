/*
Copyright © 2023 Hugobyte AI Labs<hello@hugobyte.com>
*/
package main

import (
	"github.com/hugobyte/dive/commands"
	"github.com/hugobyte/dive/styles"
)

func main() {
	styles.RenderBanner()
	commands.Execute()

}
