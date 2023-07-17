package styles

import "fmt"

func color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

var ERROR_COLOR = color("\033[0;31m]")
var BANNER_COLOR = color("\033[1;34m%s\033[0m")
var TAG_COLOR = color("\033[3;32m%s\033[0m")
