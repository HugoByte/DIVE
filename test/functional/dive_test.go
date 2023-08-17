package dive_test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"testing"

	dive "github.com/HugoByte/DIVE/test/functional"
	"github.com/hugobyte/dive/common"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

// To Print cli output to console
type testWriter struct {
	buffer bytes.Buffer
}

func TestCLIApp(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "DIVE CLI App Suite")
}

func (w *testWriter) Write(p []byte) (n int, err error) {
	w.buffer.Write(p)
	os.Stdout.Write(p)
	return len(p), nil
}

var _ = ginkgo.Describe("DIVE CLI App", func() {
	var cmd *exec.Cmd
	var stdout bytes.Buffer

	ginkgo.BeforeEach(func() {
		cmd = dive.GetBinaryCommand()
		cmd.Stdout = &testWriter{}
		cmd.Stderr = &testWriter{}
	})

	ginkgo.AfterEach(func() {
		dive.Clean()
	})

	ginkgo.Describe("Smoke Tests", func() {
		// Clean before running tests
		dive.Clean()

		ginkgo.It("should display the correct version", func() {
			cmd.Args = append(cmd.Args, "version")
			cmd.Stdout = &stdout
			err := cmd.Run()
			fmt.Println(stdout.String())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			latestVersion := common.GetLatestVersion()
			gomega.Expect(stdout.String()).To(gomega.ContainSubstring(latestVersion))
		})

		ginkgo.It("should open twitter page on browser", func() {
			cmd.Args = append(cmd.Args, "twitter")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should open Youtube Tutorial Channel on browser", func() {
			cmd.Args = append(cmd.Args, "tutorial")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should open Discord channel on browser", func() {
			cmd.Args = append(cmd.Args, "discord")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between icon and eth", func() {
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "eth")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between icon and hardhat but with icon bridge set to true", func() {
			dive.Clean()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "hardhat", "--bmvbridge")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between icon and icon", func() {
			dive.Clean()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "icon")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

	})

	ginkgo.Describe("Functional Tests", func() {
		ginkgo.It("should start bridge between icon and eth but with icon bridge set to true", func() {
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "eth", "--bmvbridge")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

		})

		ginkgo.It("should start bridge between icon and eth with verbose flag enabled", func() {
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "eth", "--verbose")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

		})

		ginkgo.It("should start bridge between icon and eth but with icon bridge set to true with verbose flag enabled", func() {
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "eth", "--bmvbridge", "--verbose")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between icon and hardhat", func() {
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "hardhat")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between icon and hardhat with verbose flag enabled", func() {
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "hardhat", "--verbose")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between icon and hardhat but with icon bridge set to true with verbose flag enabled", func() {
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "hardhat", "--bmvbridge", "--verbose")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between icon and icon with verbose flag enabled", func() {
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "icon", "--verbose")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should run single icon node", func() {
			cmd.Args = append(cmd.Args, "chain", "icon")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should run single icon node with verbose flag enabled", func() {
			cmd.Args = append(cmd.Args, "chain", "icon", "--verbose")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should run single icon node along with decentralisation", func() {
			cmd.Args = append(cmd.Args, "chain", "icon", "-d")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should run single icon node along with decentralisation with verbose flag enabled", func() {
			cmd.Args = append(cmd.Args, "chain", "icon", "-d", "--verbose")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should run single eth node", func() {
			cmd.Args = append(cmd.Args, "chain", "eth")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should run single eth node with verbose flag enabled", func() {
			cmd.Args = append(cmd.Args, "chain", "eth", "--verbose")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should run single hardhat node", func() {
			cmd.Args = append(cmd.Args, "chain", "hardhat")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should run single hardhat node with verbose flag enabled", func() {
			cmd.Args = append(cmd.Args, "chain", "hardhat", "--verbose")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should run custom Icon node", func() {
			cmd.Args = append(cmd.Args, "chain", "icon", "-c", "../../cli/sample-jsons/config0.json", "-g", "../../services/jvm/icon/static-files/config/genesis-icon-0.zip")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should run custom Icon node  with verbose flag enabled", func() {
			cmd.Args = append(cmd.Args, "chain", "icon", "-c", "../../cli/sample-jsons/config0.json", "-g", "../../services/jvm/icon/static-files/config/genesis-icon-0.zip", "--verbose")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should run custom Icon node-1", func() {
			cmd.Args = append(cmd.Args, "chain", "icon", "-c", "../../cli/sample-jsons/config1.json", "-g", "../../services/jvm/icon/static-files/config/genesis-icon-1.zip")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should run custom Icon node-1  with verbose flag enabled", func() {
			cmd.Args = append(cmd.Args, "chain", "icon", "-c", "../../cli/sample-jsons/config1.json", "-g", "../../services/jvm/icon/static-files/config/genesis-icon-1.zip", "--verbose")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should run icon node first and then decentralise it", func() {
			dive.RunIconNode()
			cmd.Args = append(cmd.Args, "chain", "icon", "decentralize", "-p", "gochain", "-k", "keystores/keystore.json", "-n", "0x3", "-e", "http://172.16.0.3:9080/api/v3/icon_dex", "-s", "icon-node-0xacbc4e")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should run icon node first and then decentralise it with verbose flag enabled", func() {
			dive.RunIconNode()
			cmd.Args = append(cmd.Args, "chain", "icon", "decentralize", "-p", "gochain", "-k", "keystores/keystore.json", "-n", "0x3", "-e", "http://172.16.0.3:9080/api/v3/icon_dex", "-s", "icon-node-0xacbc4e", "--verbose")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between icon and eth by running each chain individually ", func() {
			dive.RunDecentralizedIconNode()
			dive.RunEthNode()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "eth", "--chainAServiceName", "icon-node-0xacbc4e", "--chainBServiceName", "el-1-geth-lighthouse")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between icon and hardhat by running each chain individually ", func() {
			dive.RunDecentralizedIconNode()
			dive.RunHardhatNode()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "hardhat", "--chainAServiceName", "icon-node-0xacbc4e", "--chainBServiceName", "hardhat-node")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between icon and eth by running icon node first and then decentralising it", func() {
			dive.RunIconNode()
			dive.DecentralizeIconNode()
			dive.RunEthNode()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "eth", "--chainAServiceName", "icon-node-0xacbc4e", "--chainBServiceName", "el-1-geth-lighthouse")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between icon and hardhat by running icon node first and then decentralising it", func() {
			dive.RunIconNode()
			dive.DecentralizeIconNode()
			dive.RunHardhatNode()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "hardhat", "--chainAServiceName", "icon-node-0xacbc4e", "--chainBServiceName", "el-1-geth-lighthouse")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between icon and icon by running one custom icon chain", func() {
			dive.RunDecentralizedIconNode()
			dive.RunDecentralizedCustomIconNode()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "icon", "--chainAServiceName", "icon-node-0xacbc4e", "--chainBServiceName", "icon-node-0x42f1f3")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between icon and running custom icon later decentralising it", func() {
			dive.RunDecentralizedIconNode()
			dive.RunCustomIconNode()
			dive.DecentralizeCustomIconNode()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "icon", "--chainAServiceName", "icon-node-0xacbc4e", "--chainBServiceName", "icon-node-0x42f1f3")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between icon and icon by running one icon chain and later decentralsing it. Running another custom icon chain and then decentralising it", func() {
			dive.RunIconNode()
			dive.DecentralizeIconNode()
			dive.RunCustomIconNode()
			dive.DecentralizeCustomIconNode()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "icon", "--chainAServiceName", "icon-node-0xacbc4e", "--chainBServiceName", "icon-node-0x42f1f3")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between 2 custom icon chains", func() {
			dive.RunDecentralizedCustomIconNode_0()
			dive.RunDecentralizedCustomIconNode()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "icon", "--chainAServiceName", "icon-node-0xacbc4e", "--chainBServiceName", "icon-node-0x42f1f3")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between 2 custom icon chains by running them first and then decentralising it later", func() {
			dive.RunCustomIconNode_0()
			dive.DecentralizeCustomIconNode_0()
			dive.RunCustomIconNode()
			dive.DecentralizeCustomIconNode()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "icon", "--chainAServiceName", "icon-node-0xacbc4e", "--chainBServiceName", "icon-node-0x42f1f3")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between 2 chains when all nodes are running", func() {
			dive.RunDecentralizedIconNode()
			dive.RunEthNode()
			dive.RunHardhatNode()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "eth", "--chainAServiceName", "icon-node-0xacbc4e", "--chainBServiceName", "el-1-geth-lighthouse")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should handle error when trying to clean if no enclaves are running", func() {
			dive.Clean()
			dive.Clean()
		})

		ginkgo.It("should handle error when trying to clean if kurtosis engine is not running ", func() {
			cmd1 := exec.Command("kurtosis", "engine", "stop")
			cmd1.Run()
			// to add bin path here
			cmd2 := exec.Command("dive", "clean")
			err := cmd2.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
			cmd3 := exec.Command("kurtosis", "engine", "start")
			cmd3.Run()
		})

		ginkgo.It("should handle invalid input for bridge command", func() {
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "invalid_input")
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("should handle invalid input for bridge command", func() {
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "invalid_input", "--chainB", "eth")
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("should handle invalid input for chain command", func() {
			cmd.Args = append(cmd.Args, "chain", "invalid_input")
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("should handle invalid input for chain command", func() {
			cmd.Args = append(cmd.Args, "chain", "icon", "-c", "invalid.json", "-g", "../../services/jvm/icon/static-files/config/genesis-icon-0.zip")
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("should handle invalid input for chain command", func() {
			cmd.Args = append(cmd.Args, "chain", "icon", "-c", "../../cli/sample-jsons/config0.json", "-g", "../../services/jvm/icon/static-files/config/invalid-icon-3.zip")
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("should handle invalid input for chain command", func() {
			cmd.Args = append(cmd.Args, "chain", "icon", "-c", "../../cli/sample-jsons/invalid_config.json", "-g", "../../services/jvm/icon/static-files/config/invalid-icon-3.zip")
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("should handle invalid input for chain command", func() {
			dive.RunIconNode()
			cmd.Args = append(cmd.Args, "chain", "icon", "decentralize", "-p", "invalidPassword", "-k", "keystores/keystore.json", "-n", "0x3", "-e", "http://172.16.0.3:9080/api/v3/icon_dex", "-s", "icon-node-0xacbc4e")
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("should handle invalid input for chain command", func() {
			dive.RunIconNode()
			cmd.Args = append(cmd.Args, "chain", "icon", "decentralize", "-p", "gochain", "-k", "keystores/invalid.json", "-n", "0x3", "-e", "http://172.16.0.3:9080/api/v3/icon_dex", "-s", "icon-node-0xacbc4e")
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("should handle invalid input for chain command", func() {
			dive.RunIconNode()
			cmd.Args = append(cmd.Args, "chain", "icon", "decentralize", "-p", "gochain", "-k", "keystores/keystore.json", "-n", "0x9", "-e", "http://172.16.0.3:9080/api/v3/icon_dex", "-s", "icon-node-0xacbc4e")
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("should handle invalid input for chain command", func() {
			dive.RunIconNode()
			cmd.Args = append(cmd.Args, "chain", "icon", "decentralize", "-p", "gochain", "-k", "keystores/keystore.json", "-n", "0x3", "-e", "http://172.16.0.3:9081/api/v3/icon_dex", "-s", "icon-node-0xacbc4e")
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("should handle invalid input for chain command", func() {
			dive.RunIconNode()
			cmd.Args = append(cmd.Args, "chain", "icon", "decentralize", "-p", "gochain", "-k", "keystores/keystore.json", "-n", "0x3", "-e", "http://172.16.0.3:9080/api/v3/icon_dex", "-s", "icon-node")
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})
	})
})
