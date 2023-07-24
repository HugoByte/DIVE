package common

import (
	"fmt"
	"time"
)

const (
	DiveEnclave                  = "dive"
	DiveRemotePackagePath        = "github.com/hugobyte/dive"
	DiveIconNodeScript           = "services/jvm/icon/src/node-setup/start_icon_node.star"
	DiveIconDecentraliseScript   = "services/jvm/icon/src/node-setup/setup_icon_node.star"
	DiveEthHardhatNodeScript     = "services/evm/eth/src/node-setup/start-eth-node.star"
	DiveBridgeScript             = "main.star"
	DiveDryRun                   = false
	DiveDefaultParallelism       = 4
	DiveEthNodeAlreadyRunning    = "Eth Node Already Running"
	DiveHardhatNodeAlreadyRuning = "Hardhat Node Already Running"
	DiveIconNodeAlreadyRunning   = "Icon Node Already Running"
	DiveLogDirectory             = "/logs/"

	DiveOutFile = "dive.json"
)

const (
	linuxOSName   = "linux"
	macOSName     = "darwin"
	windowsOSName = "windows"

	openFileLinuxCommandName   = "xdg-open"
	openFileMacCommandName     = "open"
	openFileWindowsCommandName = "rundll32"

	openFileWindowsCommandFirstArgumentDefault = "url.dll,FileProtocolHandler"
)

// !!!!!!!!!!! DO NOT UPDATE! WILL BE UPDATED DURING THE RELEASE PROCESS !!!!!!!!!!!!!!!!!!!!!!
var DiveVersion = "v0.0.1-beta"

var DiveLogs bool

var DiveDiwLogFile = fmt.Sprintf("dive-%d.log", time.Now().Unix())
var DiveErorLogFile = fmt.Sprintf("dive-error-%d.log", time.Now().Unix())
