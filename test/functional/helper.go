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
	workingDir, _ := os.Getwd()
	binaryPath := filepath.Join(workingDir, "/../../cli/dive")
	return exec.Command(binaryPath)
}

// function to test and clean encalve created by DIVE
func Clean() {
	var cmd *exec.Cmd
	var stdout bytes.Buffer
	cmd = GetBinaryCommand()
	cmd.Args = append(cmd.Args, "clean")
	cmd.Stdout = &stdout
	err := cmd.Run()
	fmt.Println(stdout.String())
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunIconNode() {
	var cmd *exec.Cmd
	cmd = GetBinaryCommand()
	cmd.Args = append(cmd.Args, "chain", "icon")
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunDecentralizedIconNode() {
	var cmd *exec.Cmd
	cmd = GetBinaryCommand()
	cmd.Args = append(cmd.Args, "chain", "icon", "-d")
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunDecentralizedCustomIconNode() {
	var cmd *exec.Cmd
	cmd = GetBinaryCommand()
	cmd.Args = append(cmd.Args, "chain", "icon", "-c", "../../cli/sample-jsons/config1.json", "-g", "../../services/jvm/icon/static-files/config/genesis-icon-1.zip", "-d")
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunCustomIconNode() {
	var cmd *exec.Cmd
	cmd = GetBinaryCommand()
	cmd.Args = append(cmd.Args, "chain", "icon", "-c", "../../cli/sample-jsons/config1.json", "-g", "../../services/jvm/icon/static-files/config/genesis-icon-1.zip")
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func DecentralizeCustomIconNode() {
	var cmd *exec.Cmd
	cmd = GetBinaryCommand()
	cmd.Args = append(cmd.Args, "chain", "icon", "decentralize", "-p", "gochain", "-k", "keystores/keystore.json", "-n", "0x101", "-e", "http://172.16.0.4:9081/api/v3/icon_dex", "-s", "icon-node-0x42f1f3", "--verbose")
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func DecentralizeIconNode() {
	var cmd *exec.Cmd
	cmd = GetBinaryCommand()
	cmd.Args = append(cmd.Args, "chain", "icon", "decentralize", "-p", "gochain", "-k", "keystores/keystore.json", "-n", "0x3", "-e", "http://172.16.0.3:9080/api/v3/icon_dex", "-s", "icon-node-0xacbc4e", "--verbose")
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunEthNode() {
	var cmd *exec.Cmd
	cmd = GetBinaryCommand()
	cmd.Args = append(cmd.Args, "chain", "eth")
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunHardhatNode() {
	var cmd *exec.Cmd
	cmd = GetBinaryCommand()
	cmd.Args = append(cmd.Args, "chain", "hardhat")
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunDecentralizedCustomIconNode_0() {
	var cmd *exec.Cmd
	cmd = GetBinaryCommand()
	cmd.Args = append(cmd.Args, "chain", "icon", "-c", "../../cli/sample-jsons/config0.json", "-g", "../../services/jvm/icon/static-files/config/genesis-icon-0.zip", "-d")
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func RunCustomIconNode_0() {
	var cmd *exec.Cmd
	cmd = GetBinaryCommand()
	cmd.Args = append(cmd.Args, "chain", "icon", "-c", "../../cli/sample-jsons/config0.json", "-g", "../../services/jvm/icon/static-files/config/genesis-icon-0.zip")
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func DecentralizeCustomIconNode_0() {
	var cmd *exec.Cmd
	cmd = GetBinaryCommand()
	cmd.Args = append(cmd.Args, "chain", "icon", "decentralize", "-p", "gochain", "-k", "keystores/keystore.json", "-n", "0x3", "-e", "hhttp://172.16.0.3:9080/api/v3/icon_dex", "-s", "icon-node-0xacbc4e", "--verbose")
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}
