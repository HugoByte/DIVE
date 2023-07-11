package styles

import "fmt"

var banner = `
  ___ _____   _____ 
 |   \_ _\ \ / / __|
 | |) | | \ V /| _| 
 |___/___| \_/ |___|                         
`

func RenderBanner() {

	fmt.Println(BANNER_COLOR(banner))

}
