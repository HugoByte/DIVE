package common

const (
	DiveEnclave                = "dive"
	DiveRemotePackagePath      = "github.com/hugobyte/dive"
	DiveIconNodeScript         = "services/jvm/icon/src/node-setup/start_icon_node.star"
	DiveIconDecentraliseScript = "services/jvm/icon/src/node-setup/setup_icon_node.star"
	DiveEthHardhatNodeScript   = "services/evm/eth/src/node-setup/start-eth-node.star"
	DiveBridgeScript           = "main.star"
	DiveDryRun                 = false
	DiveDefaultParallelism     = 4
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

const (
	// !!!!!!!!!!! DO NOT UPDATE! WILL BE MANUALLY UPDATED DURING THE RELEASE PROCESS !!!!!!!!!!!!!!!!!!!!!!
	DiveVersion = "v 0.0.1(alpha)"
	// !!!!!!!!!!! DO NOT UPDATE! WILL BE MANUALLY UPDATED DURING THE RELEASE PROCESS !!!!!!!!!!!!!!!!!!!!!!
)
