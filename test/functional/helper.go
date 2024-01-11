package dive

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/onsi/gomega"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

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
	updated_path := UpdatePublicPort(enclaveName, ICON_CONFIG1)
	cmd.Args = append(cmd.Args, "chain", "icon", "-c", updated_path, "-g", ICON_GENESIS1, "-d", "--enclaveName", enclaveName)
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunCustomIconNode0(enclaveName string) {
	cmd := GetBinaryCommand()
	updated_path := UpdatePublicPort(enclaveName, ICON_CONFIG0)
	cmd.Args = append(cmd.Args, "chain", "icon", "-c", updated_path, "-g", ICON_GENESIS0, "--enclaveName", enclaveName)
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunCustomIconNode1(enclaveName string) {
	cmd := GetBinaryCommand()
	updated_path := UpdatePublicPort(enclaveName, ICON_CONFIG1)
	cmd.Args = append(cmd.Args, "chain", "icon", "-c", updated_path, "-g", ICON_GENESIS1, "--enclaveName", enclaveName)
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunCustomArchwayNode1(enclaveName string) {
	cmd := GetBinaryCommand()
	updated_path1 := UpdatePublicPorts(ARCHWAY_CONFIG1)
	cmd.Args = append(cmd.Args, "chain", "archway", "-c", updated_path1, "--enclaveName", enclaveName)
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunCustomArchwayNode0(enclaveName string) {
	cmd := GetBinaryCommand()
	updated_path1 := UpdatePublicPorts(ARCHWAY_CONFIG0)
	cmd.Args = append(cmd.Args, "chain", "archway", "-c", updated_path1, "--enclaveName", enclaveName)
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunCustomNeutronNode1(enclaveName string) {
	cmd := GetBinaryCommand()
	updated_path2 := UpdateNeutronPublicPorts(NEUTRON_CONFIG1)
	cmd.Args = append(cmd.Args, "chain", "neutron", "-c", updated_path2, "--enclaveName", enclaveName)
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunCustomNeutronNode0(enclaveName string) {
	cmd := GetBinaryCommand()
	updated_path2 := UpdateNeutronPublicPorts(NEUTRON_CONFIG0)
	cmd.Args = append(cmd.Args, "chain", "neutron", "-c", updated_path2, "--enclaveName", enclaveName)
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
	updated_path := UpdatePublicPort(enclaveName, ICON_CONFIG0)
	cmd.Args = append(cmd.Args, "chain", "icon", "-c", updated_path, "-g", ICON_GENESIS0, "-d", "--enclaveName", enclaveName)
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

var mutex = &sync.Mutex{}
var mutex1 = &sync.Mutex{}
var mutex2 = &sync.Mutex{}
var mutex3 = &sync.Mutex{}

func UpdatePublicPort(enclaveName string, filePath string) string {
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
	updatedJSON, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		panic(err)
	}

	tmpfilePath := fmt.Sprintf("updated-config-%s.json", enclaveName)
    tmpfile, err := os.Create(tmpfilePath)
    if err != nil {
        panic(err)
    }
	defer tmpfile.Close()

	// Write the updated JSON to the temporary file
	_, err = tmpfile.Write(updatedJSON)
	if err != nil {
		panic(err)
	}

	// Delete the updated temporary file
    defer os.Remove(tmpfilePath)
	
	return tmpfile.Name()
}

// Assuming Archway struct is defined as mentioned earlier

func UpdatePublicPorts(filePath1 string) string {
	mutex1.Lock()
	defer mutex1.Unlock()

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
	tmpfilePath := fmt.Sprintf("updated-archway-%s.json", name)
    tmpfile, err := os.Create(tmpfilePath)
    if err != nil {
        panic(err)
    }

	defer tmpfile.Close()

	// Write the updated JSON to the temporary file
	_, err = tmpfile.Write(updatedJSON1)
	if err != nil {
		panic(err)
	}

	 // Delete the updated temporary file
	 defer os.Remove(tmpfilePath)

	return tmpfile.Name()
}

func UpdateNeutronPublicPorts(filePath2 string) string {
	mutex2.Lock()
	defer mutex2.Unlock()

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
	tmpfilePath := fmt.Sprintf("updated-neutron-%s.json", name)
    tmpfile, err := os.Create(tmpfilePath)
    if err != nil {
        panic(err)
    }

	defer tmpfile.Close()

	// Write the updated JSON to the temporary file
	_, err = tmpfile.Write(updatedJSON)
	if err != nil {
		panic(err)
	}

     // Delete the updated temporary file
	 defer os.Remove(tmpfilePath) 

	return tmpfile.Name()
}

type NodeInfo struct {
	ServiceName    string `json:"service_name"`
	EndpointPublic string `json:"endpoint"`
	Nid            string `json:"nid"`
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

type Configuration1 struct {
	ChainType  string `json:"chain-type"`
	RelayChain struct {
		Name  string `json:"name"`
		Nodes []struct {
			Name       string `json:"name"`
			NodeType   string `json:"node-type"`
			Prometheus bool   `json:"prometheus"`
		} `json:"nodes"`
	} `json:"relaychain"`
	Para []struct {
		Name  string `json:"name"`
		Nodes []struct {
			Name       string `json:"name"`
			NodeType   string `json:"node-type"`
			Prometheus bool   `json:"prometheus"`
		} `json:"nodes"`
	} `json:"para"`
	Explorer bool `json:"explorer"`
}


func UpdateRelayChain(filePath, newChainType, newRelayChainName, enclaveName string, newExplorer, newPrometheus bool) string {
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
    // Update Explorer
    local.Explorer = newExplorer

    for i := range local.RelayChain.Nodes {
        local.RelayChain.Nodes[i].Prometheus = newPrometheus
    }

    // Update Prometheus for Para Nodes
    for i := range local.Para {
        for j := range local.Para[i].Nodes {
            local.Para[i].Nodes[j].Prometheus = newPrometheus
        }
    }

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

    // Delete the updated temporary file
    defer os.Remove(tmpfilePath)

    return tmpfile.Name()
}


func UpdateParaChain(filePath, newParaName string, newExplorer, newPrometheus bool) string {
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

	// Remove content inside RelayChain
	local.RelayChain = struct {
		Name  string `json:"name"`
		Nodes []struct {
			Name       string `json:"name"`
			NodeType   string `json:"node-type"`
			Prometheus bool   `json:"prometheus"`
		} `json:"nodes"`
	}{}

	// Update Name and Prometheus for Para Nodes
	for i := range local.Para {
		local.Para[i].Name = newParaName
		for j := range local.Para[i].Nodes {
			local.Para[i].Nodes[j].Prometheus = newPrometheus
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

	  // Delete the updated temporary file
	  defer os.Remove(tmpfilePath)

	return tmpfile.Name()
}

func UpdateChainInfo(filePath, newChainType, newRelayChainName, newParaName string, newExplorer, newPrometheus bool) string {
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
	// Update Explorer
	local.Explorer = newExplorer

	for i := range local.RelayChain.Nodes {
		local.RelayChain.Nodes[i].Prometheus = newPrometheus
	}

	// Update Name and Prometheus for Para Nodes
	for i := range local.Para {
		local.Para[i].Name = newParaName
		for j := range local.Para[i].Nodes {
			local.Para[i].Nodes[j].Prometheus = newPrometheus
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

	 // Delete the updated temporary file
	 defer os.Remove(tmpfilePath)

	return tmpfile.Name()
}
