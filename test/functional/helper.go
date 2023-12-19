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
	PrivatePort      int    `json:"private_port"`
	PublicPort       int    `json:"public_port"`
	P2PListenAddress string `json:"p2p_listen_address"`
	P2PAddress       string `json:"p2p_address"`
	CID              string `json:"cid"`
}

type Archway struct {
	ChainID     string `json:"chain_id"`
	Key         string `json:"key"`
	PrivateGRPC int    `json:"private_grpc"`
	PrivateHTTP int    `json:"private_http"`
	PrivateTCP  int    `json:"private_tcp"`
	PrivateRPC  int    `json:"private_rpc"`
	PublicGRPC  int    `json:"public_grpc"`
	PublicHTTP  int    `json:"public_http"`
	PublicTCP   int    `json:"public_tcp"`
	PublicRPC   int    `json:"public_rpc"`
	Password    string `json:"password"`
}

type Neutron struct {
	ChainID    string `json:"chain_id"`
	Key        string `json:"key"`
	Password   string `json:"password"`
	PublicGRPC int    `json:"public_grpc"`
	PublicTCP  int    `json:"public_tcp"`
	PublicHTTP int    `json:"public_http"`
	PublicRPC  int    `json:"public_rpc"`
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
	cmd.Args = append(cmd.Args, "chain", "icon", "--enclaveName", enclaveName)
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunArchwayNode(enclaveName string) {
	cmd := GetBinaryCommand()
	cmd.Args = append(cmd.Args, "chain", "archway", "--enclaveName", enclaveName)
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunNeutronNode(enclaveName string) {
	cmd := GetBinaryCommand()
	cmd.Args = append(cmd.Args, "chain", "neutron", "--enclaveName", enclaveName)
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunDecentralizedIconNode(enclaveName string) {
	cmd := GetBinaryCommand()
	cmd.Args = append(cmd.Args, "chain", "icon", "-d", "--enclaveName", enclaveName)
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunDecentralizedCustomIconNode1(enclaveName string) {
	cmd := GetBinaryCommand()
	filepath := "../../cli/sample-jsons/config1.json"
	updated_path := UpdatePublicPort(filepath)
	cmd.Args = append(cmd.Args, "chain", "icon", "-c", updated_path, "-g", "./config/genesis-icon-1.zip", "-d", "--enclaveName", enclaveName)
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunCustomIconNode(enclaveName string) {
	cmd := GetBinaryCommand()
	filepath := "../../cli/sample-jsons/config0.json"
	updated_path := UpdatePublicPort(filepath)
	cmd.Args = append(cmd.Args, "chain", "icon", "-c", updated_path, "-g", "./config/genesis-icon-0.zip", "--enclaveName", enclaveName)
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunCustomArchwayNode1(enclaveName string) {
	cmd := GetBinaryCommand()
	filepath1 := "../../cli/sample-jsons/archway1.json"
	updated_path1 := UpdatePublicPorts(filepath1)
	cmd.Args = append(cmd.Args, "chain", "archway", "-c", updated_path1,"--enclaveName", enclaveName )
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunCustomArchwayNode0(enclaveName string) {
	cmd := GetBinaryCommand()
	filepath1 := "../../cli/sample-jsons/archway.json"
	updated_path1 := UpdatePublicPorts(filepath1)
	cmd.Args = append(cmd.Args, "chain", "archway", "-c",updated_path1,"--enclaveName", enclaveName )
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunCustomNeutronNode1(enclaveName string) {
	cmd := GetBinaryCommand()
	filepath2 := "../../cli/sample-jsons/neutron1.json"
	updated_path2 := UpdateNeutronPublicPorts(filepath2)
	cmd.Args = append(cmd.Args, "chain", "neutron", "-c", updated_path2,"--enclaveName", enclaveName)
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunCustomNeutronNode0(enclaveName string) {
	cmd := GetBinaryCommand()
	filepath2 := "../../cli/sample-jsons/neutron.json"
	updated_path2 := UpdateNeutronPublicPorts(filepath2)
	cmd.Args = append(cmd.Args, "chain", "neutron", "-c", updated_path2,"--enclaveName", enclaveName)
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func DecentralizeCustomIconNode() {
	cmd := GetBinaryCommand()
	cmd.Args = append(cmd.Args, "chain", "icon", "decentralize", "-p", "gochain", "-k", "keystores/keystore.json", "-n", "0x101", "-e", "http://172.16.0.4:9081/api/v3/icon_dex", "-s", "icon-node-0x42f1f3", "--verbose")
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func DecentralizeIconNode(enclaveName string) {
	cmd := GetBinaryCommand()
	cmd.Args = append(cmd.Args, "chain", "icon", "decentralize", "-p", "gochain", "-k", "keystores/keystore.json", "-n", "0x3", "-e", "http://172.16.0.3:9080/api/v3/icon_dex", "-s", "icon-node-0xacbc4e", "--verbose", "--enclaveName", enclaveName)
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunEthNode(enclaveName string) {
	cmd := GetBinaryCommand()
	cmd.Args = append(cmd.Args, "chain", "eth", "--enclaveName", enclaveName)
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunHardhatNode(enclaveName string) {
	cmd := GetBinaryCommand()
	cmd.Args = append(cmd.Args, "chain", "hardhat", "--enclaveName", enclaveName)
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunDecentralizedCustomIconNode0(enclaveName string) {
	cmd := GetBinaryCommand()
	filepath := "../../cli/sample-jsons/config0.json"
	updated_path := UpdatePublicPort(filepath)
	cmd.Args = append(cmd.Args, "chain", "icon", "-c", updated_path, "-g", "./config/genesis-icon-0.zip", "-d", "--enclaveName", enclaveName)
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunCustomIconNode_0() {
	cmd := GetBinaryCommand()
	cmd.Args = append(cmd.Args, "chain", "icon", "-c", "../../cli/sample-jsons/config0.json", "-g", "./config/genesis-icon-0.zip")
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

func UpdatePublicPort(filePath string) string {
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

//var mutex = &sync.Mutex{}

// Assuming Archway struct is defined as mentioned earlier

func UpdatePublicPorts(filePath1 string) string {
	mutex.Lock()
	defer mutex.Unlock()

	// Read the content of the existing JSON file
	fileContent1, err := os.ReadFile(filePath1)
	if err != nil {
		panic(err)
	}

	// Unmarshal JSON into a struct
	var archway Archway
	err = json.Unmarshal(fileContent1, &archway)
	if err != nil {
		panic(err)
	}

	// Get available ports for PublicGRPC, PublicHTTP, PublicTCP, and PublicRPC
	availableGRPC, err := common.GetAvailablePort()
	if err != nil {
		panic(err)
	}
	availableHTTP, err := common.GetAvailablePort()
	if err != nil {
		panic(err)
	}
	availableTCP, err := common.GetAvailablePort()
	if err != nil {
		panic(err)
	}
	availableRPC, err := common.GetAvailablePort()
	if err != nil {
		panic(err)
	}

	// Update the Public ports fields in the Archway struct
	archway.PublicGRPC = availableGRPC
	archway.PublicHTTP = availableHTTP
	archway.PublicTCP = availableTCP
	archway.PublicRPC = availableRPC

	// Marshal the updated struct into JSON with indentation
	updatedJSON1, err := json.MarshalIndent(archway, "", "    ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(updatedJSON1))

	// Generate a random name for the new file
	name := GenerateRandomName()
	tmpfile, err := os.Create(fmt.Sprintf("updated-archway-%s.json", name))
	if err != nil {
		panic(err)
	}
	defer tmpfile.Close()

	// Write the updated JSON to the temporary file
	_, err = tmpfile.Write(updatedJSON1)
	if err != nil {
		panic(err)
	}

	return tmpfile.Name()
}

func UpdateNeutronPublicPorts(filePath2 string) string {
	mutex.Lock()
	defer mutex.Unlock()

	// Read the content of the existing JSON file
	fileContent2, err := os.ReadFile(filePath2)
	if err != nil {
		panic(err)
	}

	// Unmarshal JSON into a Neutron struct
	var neutron Neutron
	err = json.Unmarshal(fileContent2, &neutron)
	if err != nil {
		panic(err)
	}

	// Get available ports for PublicGRPC, PublicHTTP, PublicTCP, and PublicRPC
	availableGRPC, err := common.GetAvailablePort()
	if err != nil {
		panic(err)
	}
	availableHTTP, err := common.GetAvailablePort()
	if err != nil {
		panic(err)
	}
	availableTCP, err := common.GetAvailablePort()
	if err != nil {
		panic(err)
	}
	availableRPC, err := common.GetAvailablePort()
	if err != nil {
		panic(err)
	}

	// Update the Public ports fields in the Neutron struct
	neutron.PublicGRPC = availableGRPC
	neutron.PublicHTTP = availableHTTP
	neutron.PublicTCP = availableTCP
	neutron.PublicRPC = availableRPC

	// Marshal the updated struct into JSON with indentation
	updatedJSON, err := json.MarshalIndent(neutron, "", "    ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(updatedJSON))

	// Generate a random name for the new file
	name := GenerateRandomName()
	tmpfile, err := os.Create(fmt.Sprintf("updated-neutron-%s.json", name))
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
