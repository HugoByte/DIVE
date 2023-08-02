package common

var DiveLogs bool

// !!!!!!!!!!! DO NOT UPDATE! WILL BE UPDATED DURING THE RELEASE PROCESS !!!!!!!!!!!!!!!!!!!!!!
var DiveVersion = "v0.0.2-beta"

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
	DiveDitLogFile               = "divelog.log"
	DiveErorLogFile              = "error.log"
	DiveOutFile                  = "dive.json"
	ServiceFilePath              = "services.json"
	starlarkScript               = `
def run(plan, args):
	plan.stop_service(name=args["service_name"])
	plan.print(args["uuid"]) # we add this print of a random UUID to make sure the single stop_service above won't get cached
`
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
