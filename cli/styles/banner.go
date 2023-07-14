package styles

import "fmt"

var banner = `
  ___ _____   _____ 
 |   \_ _\ \ / / __|
 | |) | | \ V /| _| 
 |___/___| \_/ |___|
	
	%s
`

func RenderBanner() {

	banner := fmt.Sprintf(BANNER_COLOR(banner), TAG_COLOR("Developed by Hugobyte AI Labs and Powered by Kurtosis"))

	fmt.Println(banner)

}
