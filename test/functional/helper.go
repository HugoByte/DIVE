package dive

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
	"github.com/onsi/gomega"

	"github.com/hugobyte/dive-core/cli/common"
)


type Configuration struct {
	PrivatePort        int    `json:"private_port"`
	PublicPort         int    `json:"public_port"`
	P2PListenAddress   string `json:"p2p_listen_address"`
	P2PAddress         string `json:"p2p_address"`
	CID                string `json:"cid"`
}

func GetBinaryCommand() *exec.Cmd {
	binaryPath := GetBinPath()
	return exec.Command(binaryPath)
}

func GetBinPath() string {
	workingDir, _ := os.Getwd()
	binaryPath := filepath.Join(workingDir, "/../../cli/dive")
	return binaryPath
}

// function to generate random enclave name
func GenerateRandomName() string {
	id := uuid.New()
	return id.String()
}

// function to test and clean encalve created by DIVE
func Clean(enclaveName string) {
	var stdout bytes.Buffer
	cmd := GetBinaryCommand()
	if enclaveName == "all" {
		cmd.Args = append(cmd.Args, "clean", "-a")
	} else {
		cmd.Args = append(cmd.Args, "clean", "--enclaveName", enclaveName)
	}
	cmd.Stdout = &stdout
	err := cmd.Run()
	fmt.Println(stdout.String())
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunIconNode(enclaveName string) {
	cmd := GetBinaryCommand()
	// enclaveName := generateRandomName()
	cmd.Args = append(cmd.Args, "chain", "icon", "--enclaveName", enclaveName)
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunArchwayNode() {
	cmd := GetBinaryCommand()
	cmd.Args = append(cmd.Args, "chain", "archway")
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunNeutronNode() {
	cmd := GetBinaryCommand()
	cmd.Args = append(cmd.Args, "chain", "neutron")
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunDecentralizedIconNode() {
	cmd := GetBinaryCommand()
	cmd.Args = append(cmd.Args, "chain", "icon", "-d")
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunDecentralizedCustomIconNode1() {
	cmd := GetBinaryCommand()
	cmd.Args = append(cmd.Args, "chain", "icon", "-c", "../../cli/sample-jsons/config1.json", "-g", "../../services/jvm/icon/static-files/config/genesis-icon-1.zip", "-d")
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunCustomIconNode() {
	cmd := GetBinaryCommand()
	cmd.Args = append(cmd.Args, "chain", "icon", "-c", "../../cli/sample-jsons/config1.json", "-g", "../../services/jvm/icon/static-files/config/genesis-icon-1.zip")
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunCustomArchwayNode1() {
	cmd := GetBinaryCommand()
	cmd.Args = append(cmd.Args, "chain", "archway", "-c", "../../cli/sample-jsons/archway1.json")
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunCustomArchwayNode0() {
	cmd := GetBinaryCommand()
	cmd.Args = append(cmd.Args, "chain", "archway", "-c", "../../cli/sample-jsons/archway.json")
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunCustomNeutronNode1() {
	cmd := GetBinaryCommand()
	cmd.Args = append(cmd.Args, "chain", "neutron", "-c", "../../cli/sample-jsons/neutron1.json")
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunCustomNeutronNode0() {
	cmd := GetBinaryCommand()
	cmd.Args = append(cmd.Args, "chain", "neutron", "-c", "../../cli/sample-jsons/neutron.json")
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func DecentralizeCustomIconNode() {
	cmd := GetBinaryCommand()
	cmd.Args = append(cmd.Args, "chain", "icon", "decentralize", "-p", "gochain", "-k", "keystores/keystore.json", "-n", "0x101", "-e", "http://172.16.0.4:9081/api/v3/icon_dex", "-s", "icon-node-0x42f1f3", "--verbose")
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func DecentralizeIconNode() {
	cmd := GetBinaryCommand()
	cmd.Args = append(cmd.Args, "chain", "icon", "decentralize", "-p", "gochain", "-k", "keystores/keystore.json", "-n", "0x3", "-e", "http://172.16.0.3:9080/api/v3/icon_dex", "-s", "icon-node-0xacbc4e", "--verbose")
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunEthNode() {
	cmd := GetBinaryCommand()
	cmd.Args = append(cmd.Args, "chain", "eth")
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunHardhatNode() {
	cmd := GetBinaryCommand()
	cmd.Args = append(cmd.Args, "chain", "hardhat")
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunDecentralizedCustomIconNode0() {
	cmd := GetBinaryCommand()
	cmd.Args = append(cmd.Args, "chain", "icon", "-c", "../../cli/sample-jsons/config0.json", "-g", "../../services/jvm/icon/static-files/config/genesis-icon-0.zip", "-d")
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunCustomIconNode_0() {
	cmd := GetBinaryCommand()
	cmd.Args = append(cmd.Args, "chain", "icon", "-c", "../../cli/sample-jsons/config0.json", "-g", "../../services/jvm/icon/static-files/config/genesis-icon-0.zip")
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func DecentralizeCustomIconNode_0() {
	cmd := GetBinaryCommand()
	cmd.Args = append(cmd.Args, "chain", "icon", "decentralize", "-p", "gochain", "-k", "keystores/keystore.json", "-n", "0x3", "-e", "http://172.16.0.3:9080/api/v3/icon_dex", "-s", "icon-node-0xacbc4e", "--verbose")
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}


var mutex = &sync.Mutex{}
func UpdatePublicPort(filePath string) string{
	mutex.Lock()
	defer mutex.Unlock()
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	// Unmarshal JSON into a struct
	var config Configuration
	err = json.Unmarshal(fileContent, &config)
	if err != nil {
		panic(err)
	}

	availablePort, err := common.GetAvailablePort()
	if err != nil {
		panic(err)
	}
	config.PublicPort = availablePort
	// config.CID =  "0xacbc4e"
	// config.P2PAddress = "8080"
	// config.P2PListenAddress = "7080"
	// config.PrivatePort = 9080	

	updatedJSON, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(updatedJSON))

	name := GenerateRandomName()
	tmpfile, err := os.Create(fmt.Sprintf("updated-config-%s.json", name))
	if err != nil {
		panic(err)
	}
	defer tmpfile.Close()

	// Write the updated JSON to the temporary file
	_, err = tmpfile.Write(updatedJSON)
	if err != nil {
		panic(err)
	}

	return tmpfile.Name()
}
	

	// // Write the updated JSON back to the same file
	// err = os.WriteFile(filePath, updatedJSON, 0644)
	// if err != nil {
	// 	panic(err)
	// }


