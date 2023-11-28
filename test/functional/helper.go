package dive

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/onsi/gomega"
)

func GetBinaryCommand() *exec.Cmd {
	binaryPath := GetBinPath()
	return exec.Command(binaryPath)
}

func GetBinPath() string {
	workingDir, _ := os.Getwd()
	binaryPath := filepath.Join(workingDir, "/../../cli/dive")
	return binaryPath
}

// function to test and clean encalve created by DIVE
func Clean() {
	var stdout bytes.Buffer
	cmd := GetBinaryCommand()
	cmd.Args = append(cmd.Args, "clean")
	cmd.Stdout = &stdout
	err := cmd.Run()
	fmt.Println(stdout.String())
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunIconNode() {
	cmd := GetBinaryCommand()
	cmd.Args = append(cmd.Args, "chain", "icon")
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
