package common

// !!!!!!!!!!! DO NOT UPDATE! WILL BE UPDATED DURING THE RELEASE PROCESS !!!!!!!!!!!!!!!!!!!!!!
const DiveVersion = "v0.0.14-beta"

const (
	DiveEnclave                  = "dive"
	DiveRemotePackagePath        = "github.com/hugobyte/dive-packages"
	DiveIconNodeScript           = "services/jvm/icon/src/node-setup/start_icon_node.star"
	DiveIconDecentralizeScript   = "services/jvm/icon/src/node-setup/setup_icon_node.star"
	DiveEthHardhatNodeScript     = "services/evm/eth/src/node-setup/start-eth-node.star"
	DiveArchwayNodeScript        = "services/cosmvm/archway/src/node-setup/start_node.star"
	DiveCosmosDefaultNodeScript  = "services/cosmvm/cosmos_chains.star"
	DiveNeutronNodeScript        = "services/cosmvm/neutron/src/node-setup/start_node.star"
	RelayServiceNameIconToCosmos = "ibc-relayer"
	DiveNeutronDefaultNodeScript = "services/cosmvm/neutron/neutron.star"
	DiveBridgeBtpScript          = "/services/bridges/btp/src/bridge.star"
	DiveBridgeIbcScript          = "/services/bridges/ibc/src/bridge.star"
	DiveDryRun                   = false
	DiveDefaultParallelism       = 4
	DiveLogDirectory             = "/logs/"
	DiveDitLogFile               = "dive.log"
	DiveErrorLogFile             = "error.log"
	DiveOutFile                  = "dive_%s.json"
	ServiceFilePath              = "services_%s.json"
	DiveAppDir                   = ".dive"
	removeServiceStarlarkScript  = `
def run(plan,args):
		plan.remove_service(name=args["service_name"])
`
	stopServiceStarlarkScript = `
def run(plan, args):
	plan.stop_service(name=args["service_name"])
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

const (
	ErrorCodeGeneral ErrorCode = iota + 1000
)

const (
	UnknownError ErrorCode = ErrorCodeGeneral + iota
	FileReadError
	FileWriteError
	FileOpenError
	FileNotExistError
	KurtosisInitError
	CLIInitError
	InvalidEnclaveNameError
	UnsupportedOSError
	InvalidCommandError
	InvalidEnclaveError
	EnclaveNotExistError
	InvalidEnclaveContextError
	InvalidEnclaveConfigError
	InvalidCommandArgumentsError
	InvalidKurtosisContextError
	DataMarshallError
	DataUnMarshallError
	StarlarkRunFailedError
	NotFoundError
	StarlarkResponseError
	InvalidPathError
	InvalidFileError
	KurtosisServiceError
	InvalidChain
	PortError
)

var DiveLogs bool
var EnclaveName string
