package dive

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/google/uuid"
	"github.com/onsi/gomega"
)

type NodeInfo struct {
	ServiceName    string `json:"service_name"`
	EndpointPublic string `json:"endpoint"`
	Nid            string `json:"nid"`
}

type CosmosNodeInfo struct{
	EndpointPublic string `json:"endpoint_public"`
}

type Configuration1 struct {
	ChainType  string `json:"chain_type"`
	RelayChain struct {
		Name  string `json:"name"`
		Nodes []struct {
			Name       string `json:"name"`
			NodeType   string `json:"node_type"`
			Prometheus bool   `json:"prometheus"`
		} `json:"nodes"`
	} `json:"relaychain"`
	Parachains []struct {
		Name  string `json:"name"`
		Nodes []struct {
			Name       string `json:"name"`
			NodeType   string `json:"node_type"`
			Prometheus bool   `json:"prometheus"`
		} `json:"nodes"`
	} `json:"Parachains"`
	Explorer bool `json:"explorer"`
}

var mutex = &sync.Mutex{}
var mutex3 = &sync.Mutex{}

func GetCosmosLatestBlock(nodeURI string) (height int64, err error) {
	http, _ := client.NewClientFromNode(nodeURI)
	cliCtx := client.Context{}.WithClient(http)
	height, err = rpc.GetChainHeight(cliCtx)
	return height, err
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
	cmd.Args = append(cmd.Args, "chain", "icon", "-c", ICON_CONFIG1, "-g", ICON_GENESIS1, "-d", "--enclaveName", enclaveName)
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunCustomIconNode0(enclaveName string) {
	cmd := GetBinaryCommand()
	cmd.Args = append(cmd.Args, "chain", "icon", "-c", ICON_CONFIG0, "-g", ICON_GENESIS0, "--enclaveName", enclaveName)
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunCustomIconNode1(enclaveName string) {
	cmd := GetBinaryCommand()
	cmd.Args = append(cmd.Args, "chain", "icon", "-c", ICON_CONFIG1, "-g", ICON_GENESIS1, "--enclaveName", enclaveName)
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunCustomArchwayNode1(enclaveName string) {
	cmd := GetBinaryCommand()
	cmd.Args = append(cmd.Args, "chain", "archway", "-c", ARCHWAY_CONFIG1, "--enclaveName", enclaveName)
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunCustomArchwayNode0(enclaveName string) {
	cmd := GetBinaryCommand()
	cmd.Args = append(cmd.Args, "chain", "archway", "-c", ARCHWAY_CONFIG0, "--enclaveName", enclaveName)
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunCustomNeutronNode1(enclaveName string) {
	cmd := GetBinaryCommand()
	//updated_path2 := UpdateNeutronPublicPorts(NEUTRON_CONFIG1)
	cmd.Args = append(cmd.Args, "chain", "neutron", "-c", NEUTRON_CONFIG1, "--enclaveName", enclaveName)
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunCustomNeutronNode0(enclaveName string) {
	cmd := GetBinaryCommand()
	//updated_path2 := UpdateNeutronPublicPorts(NEUTRON_CONFIG0)
	cmd.Args = append(cmd.Args, "chain", "neutron", "-c", NEUTRON_CONFIG0, "--enclaveName", enclaveName)
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func DecentralizeCustomIconNode(nid string, endpoint string, serviceName string, enclaveName string) {
	cmd := GetBinaryCommand()
	cmd.Args = append(cmd.Args, "chain", "icon", "decentralize", "-p", "gochain", "-k", "keystores/keystore.json", "-n", nid, "-e", endpoint, "-s", serviceName, "--verbose", "--enclaveName", enclaveName)
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
	cmd.Args = append(cmd.Args, "chain", "icon", "-c", ICON_CONFIG0, "-g", ICON_GENESIS0, "-d", "--enclaveName", enclaveName)
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func GetServiceDetails(servicesJson string, service string) (serviceName string, endpoint string, nid string) {
	var data map[string]NodeInfo
	mutex3.Lock()
	defer mutex3.Unlock()

	fileContent2, err := os.ReadFile(servicesJson)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(fileContent2, &data)
	if err != nil {
		panic(err)
	}

	for key, value := range data {
		if key == service {
			serviceName = value.ServiceName
			endpoint = value.EndpointPublic
			nid = value.Nid
		}
	}
	return serviceName, endpoint, nid

}

func GetServiceDetailsCosmos(servicesJson string, service string) (endpoint string) {
	var data map[string]CosmosNodeInfo
	mutex3.Lock()
	defer mutex3.Unlock()

	fileContent2, err := os.ReadFile(servicesJson)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(fileContent2, &data)
	if err != nil {
		panic(err)
	}

	for key, value := range data {
		if key == service {
			endpoint = value.EndpointPublic
		}
	}
	return endpoint

}

func UpdateRelayChain(filePath, newChainType, newRelayChainName, enclaveName string, newNodeType1, newNodeType2 string, relayChain string) string {
	mutex.Lock()
	defer mutex.Unlock()

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	var local Configuration1
	err = json.Unmarshal(fileContent, &local)
	if err != nil {
		panic(err)
	}

	// Update ChainType and RelayChain Name
	local.ChainType = newChainType
	local.RelayChain.Name = newRelayChainName

	// Update NodeType for RelayChain Nodes
	for i := range local.RelayChain.Nodes {
		if i%2 == 0 {
			local.RelayChain.Nodes[i].NodeType = newNodeType1
		} else {
			local.RelayChain.Nodes[i].NodeType = newNodeType2
		}
	}

	// Update Prometheus for Para Nodes
	for i := range local.Parachains {
		if relayChain == "kusama" {
			local.Parachains[i].Name = "karura"
		} else {
			local.Parachains[i].Name = "acala"
		}
	}

	// Remove content inside RelayChain
	local.Parachains = []struct {
		Name  string `json:"name"`
		Nodes []struct {
			Name       string `json:"name"`
			NodeType   string `json:"node_type"`
			Prometheus bool   `json:"prometheus"`
		} `json:"nodes"`
	}{}

	updatedJSON, err := json.MarshalIndent(local, "", "    ")
	if err != nil {
		panic(err)
	}

	tmpfilePath := fmt.Sprintf("updated-config-%s.json", enclaveName)
	tmpfile, err := os.Create(tmpfilePath)
	if err != nil {
		panic(err)
	}
	defer tmpfile.Close()

	_, err = tmpfile.Write(updatedJSON)
	if err != nil {
		panic(err)
	}

	return tmpfile.Name()
}

func UpdateParaChain(filePath, newChainType, newParaName string) string {
	mutex.Lock()
	defer mutex.Unlock()

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	var local Configuration1
	err = json.Unmarshal(fileContent, &local)
	if err != nil {
		panic(err)
	}
	//update chain type
	local.ChainType = newChainType
	// Remove content inside RelayChain
	local.RelayChain = struct {
		Name  string `json:"name"`
		Nodes []struct {
			Name       string `json:"name"`
			NodeType   string `json:"node_type"`
			Prometheus bool   `json:"prometheus"`
		} `json:"nodes"`
	}{}

	// Update Name  Para Nodes
	for i := range local.Parachains {
		local.Parachains[i].Name = newParaName
	}

	updatedJSON, err := json.MarshalIndent(local, "", "    ")
	if err != nil {
		panic(err)
	}

	tmpfilePath := fmt.Sprintf("updated-local.json")
	tmpfile, err := os.Create(tmpfilePath)
	if err != nil {
		panic(err)
	}

	defer tmpfile.Close()

	_, err = tmpfile.Write(updatedJSON)
	if err != nil {
		panic(err)
	}

	return tmpfile.Name()
}

func UpdateChainInfo(filePath, newChainType, newRelayChainName, newParaName string, newNodeType1, newNodeType2 string) string {
	mutex.Lock()
	defer mutex.Unlock()

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	var local Configuration1
	err = json.Unmarshal(fileContent, &local)
	if err != nil {
		panic(err)
	}

	// Update ChainType and RelayChain Name
	local.ChainType = newChainType
	local.RelayChain.Name = newRelayChainName

	// Update Name and Prometheus for Para Nodes
	for i := range local.Parachains {
		local.Parachains[i].Name = newParaName

	}
	// Update NodeType for RelayChain Nodes
	for i := range local.RelayChain.Nodes {
		if i%2 == 0 {
			local.RelayChain.Nodes[i].NodeType = newNodeType1
		} else {
			local.RelayChain.Nodes[i].NodeType = newNodeType2
		}
	}

	updatedJSON, err := json.MarshalIndent(local, "", "    ")
	if err != nil {
		panic(err)
	}

	tmpfilePath := fmt.Sprintf("updated-local.json")
	tmpfile, err := os.Create(tmpfilePath)
	if err != nil {
		panic(err)
	}

	defer tmpfile.Close()

	_, err = tmpfile.Write(updatedJSON)
	if err != nil {
		panic(err)
	}

	return tmpfile.Name()
}
func CheckInvalidTestnet(selectedParaChain string, invalidParaChainlist []string) bool {
	for _, paraChainName := range invalidParaChainlist {
		if selectedParaChain == paraChainName {
			return true
		}
	}
	return false
}
