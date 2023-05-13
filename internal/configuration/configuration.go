package configuration

import "github.com/logrusorgru/aurora/v3"

const (
	// NAME is this projet's name
	NAME string = "xurl"
	// VERSION is this projet's version
	VERSION string = "0.0.0"
	// DESCRIPTION is this projet's description
	DESCRIPTION string = "A CLI utility to pull out bits of URLs."
)

var (
	// BANNER is this project's CLI display banner
	BANNER = aurora.Sprintf(
		aurora.BrightBlue(`
                 _ 
__  ___   _ _ __| |
\ \/ / | | | '__| |
 >  <| |_| | |  | |
/_/\_\\__,_|_|  |_| %s

%s`).Bold(),
		aurora.BrightYellow("v"+VERSION).Bold(),
		aurora.BrightGreen(DESCRIPTION).Italic(),
	)
)
