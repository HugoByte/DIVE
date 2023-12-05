package icon

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hugobyte/dive-core/cli/common"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/enclaves"
)

func genesismanager(enclaveContext *enclaves.EnclaveContext) (*genesisHandler, error) {

	gm := genesisHandler{}

	var genesisFilePath = genesis

	if genesisFilePath != "" {
		genesisFileName := filepath.Base(genesisFilePath)
		if _, err := os.Stat(genesisFilePath); err != nil {
			return nil, err
		}

		_, d, err := enclaveContext.UploadFiles(genesisFilePath, genesisFileName)
		if err != nil {
			return nil, err
		}

		gm.uploadedFiles = fmt.Sprintf(`{"file_path":"%s","file_name":"%s"}`, d, genesisFileName)
	} else {
		gm.genesisFile = filepath.Base(DefaultIconGenesisFile)
		gm.genesisPath = DefaultIconGenesisFile
		gm.uploadedFiles = `{}`

	}

	return &gm, nil
}

func LoadConfig(cliContext *common.Cli, config ConfigLoader, filePath string) error {
	if filePath == "" {
		config.LoadDefaultConfig()
	} else {
		err := config.LoadConfigFromFile(cliContext, filePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetDecentralizeParams(serviceName, nodeEndpoint, keystorePath, keystorepassword, networkID string) string {
	return fmt.Sprintf(`{"service_name":"%s","uri":"%s","keystorepath":"%s","keypassword":"%s","nid":"%s"}`, serviceName, nodeEndpoint, keystorePath, keystorepassword, networkID)
}
