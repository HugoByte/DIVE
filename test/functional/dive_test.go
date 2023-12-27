package dive_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	dive "github.com/HugoByte/DIVE/test/functional"
	"github.com/hugobyte/dive-core/cli/cmd/utility"
	"github.com/hugobyte/dive-core/cli/common"
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
	// ginkgo.AfterSuite(func() {
	// 	dive.Clean("all")
	// })
	// ginkgo.AfterAll(func ()  {
	// 	dive.Clean("all")
	// })
}

func (w *testWriter) Write(p []byte) (n int, err error) {
	w.buffer.Write(p)
	os.Stdout.Write(p)
	return len(p), nil
}

var _ = ginkgo.Describe("DIVE CLI App", func() {
	var cmd *exec.Cmd
	var stdout bytes.Buffer

	//    var runKusama bool
	//     var runPolkadot bool

	// BeforeSuite hook
	// ginkgo.BeforeSuite(func() {
	//     flag.BoolVar(&runKusama, "kusama", false, "Run tests for Kusama")
	//     flag.BoolVar(&runPolkadot, "polkadot", false, "Run tests for Polkadot")
	//     flag.Parse()
	// })

	// run clean before each test
	ginkgo.BeforeEach(func() {
		cmd = dive.GetBinaryCommand()
		cmd.Stdout = &testWriter{}
		cmd.Stderr = &testWriter{}
	})

	ginkgo.Describe("Smoke Tests", func() {
		ginkgo.It("should display the correct version", func() {
			cmd.Args = append(cmd.Args, "version")
			cmd.Stdout = &stdout
			err := cmd.Run()
			fmt.Println(stdout.String())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			cli := common.GetCli()
			latestVersion := utility.GetLatestVersion(cli)
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

		ginkgo.It("should start bridge between icon and eth correctly-1", func() {
			// dive.Clean()
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "eth", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between icon and hardhat but with icon bridge set to true-1", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "hardhat", "--bmvbridge", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between icon and icon-1", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "icon", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between archway and archway using ibc-1", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "archway", "--chainB", "archway", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

	})

	ginkgo.Describe("Bridge command Test", func() {
		ginkgo.It("should start bridge between icon and eth but with icon bridge set to true-1", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "eth", "--bmvbridge", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between icon and eth with verbose flag enabled", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "eth", "--verbose", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)

		})

		ginkgo.It("should start bridge between icon and eth but with icon bridge set to true with verbose flag enabled", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "eth", "--bmvbridge", "--verbose", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between icon and hardhat", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "hardhat", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between icon and hardhat with verbose flag enabled", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "hardhat", "--verbose", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between icon and hardhat but with icon bridge set to true with verbose flag enabled", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "hardhat", "--bmvbridge", "--verbose", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between icon and icon with verbose flag enabled", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "icon", "--verbose", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between icon and eth by running each chain individually-2abc ", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunDecentralizedIconNode(enclaveName)
			dive.RunEthNode(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "eth", "--chainAServiceName", "icon-node-0xacbc4e", "--chainBServiceName", "el-1-geth-lighthouse", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between icon and hardhat by running each chain individually -1bcd", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunDecentralizedIconNode(enclaveName)
			dive.RunHardhatNode(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "hardhat", "--chainAServiceName", "icon-node-0xacbc4e", "--chainBServiceName", "hardhat-node", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between icon and eth by running icon node first and then decentralising it", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunIconNode(enclaveName)
			dive.DecentralizeIconNode(enclaveName)
			dive.RunEthNode(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "eth", "--chainAServiceName", "icon-node-0xacbc4e", "--chainBServiceName", "el-1-geth-lighthouse", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between icon and hardhat by running icon node first and then decentralising it-1-abc", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunIconNode(enclaveName)
			dive.DecentralizeIconNode(enclaveName)
			dive.RunHardhatNode(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "hardhat", "--chainAServiceName", "icon-node-0xacbc4e", "--chainBServiceName", "hardhat-node", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between icon and icon by running one custom icon chain-1-abc", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunDecentralizedIconNode(enclaveName)
			dive.RunDecentralizedCustomIconNode1(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "icon", "--chainAServiceName", "icon-node-0xacbc4e", "--chainBServiceName", "icon-node-0x42f1f3", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between icon and running custom icon later decentralising it", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunDecentralizedIconNode(enclaveName)
			dive.RunCustomIconNode(enclaveName)
			dive.DecentralizeCustomIconNode()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "icon", "--chainAServiceName", "icon-node-0xacbc4e", "--chainBServiceName", "icon-node-0x42f1f3", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between icon and icon by running one icon chain and later decentralsing it. Running another custom icon chain and then decentralising it", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunIconNode(enclaveName)
			dive.DecentralizeIconNode(enclaveName)
			dive.RunCustomIconNode(enclaveName)
			dive.DecentralizeCustomIconNode()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "icon", "--chainAServiceName", "icon-node-0xacbc4e", "--chainBServiceName", "icon-node-0x42f1f3", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between 2 custom icon chains", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunDecentralizedCustomIconNode0(enclaveName)
			dive.RunDecentralizedCustomIconNode1(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "icon", "--chainAServiceName", "icon-node-0xacbc4e", "--chainBServiceName", "icon-node-0x42f1f3", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between 2 custom icon chains by running them first and then decentralising it later", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunCustomIconNode_0()
			dive.DecentralizeCustomIconNode_0()
			dive.RunCustomIconNode(enclaveName)
			dive.DecentralizeCustomIconNode()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "icon", "--chainAServiceName", "icon-node-0xacbc4e", "--chainBServiceName", "icon-node-0x42f1f3", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between 2 chains when all nodes are running-1", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunDecentralizedIconNode(enclaveName)
			dive.RunEthNode(enclaveName)
			dive.RunHardhatNode(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "eth", "--chainAServiceName", "icon-node-0xacbc4e", "--chainBServiceName", "el-1-geth-lighthouse", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should handle invalid input for bridge command", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "invalid_input", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should handle invalid input for bridge command", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "invalid_input", "--chainB", "eth", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should handle invalid input ibc bridge command", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "archway", "--chainB", "invalid", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between archway and archway by running one custom archway chain-1", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunArchwayNode(enclaveName)
			dive.RunCustomArchwayNode1(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "archway", "--chainB", "archway", "--chainAServiceName", "node-service-constantine-3", "--chainBServiceName", "node-service-archway-node-1", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between 2 custom archway chains", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunCustomArchwayNode0(enclaveName)
			dive.RunCustomArchwayNode1(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "archway", "--chainB", "archway", "--chainAServiceName", "node-service-archway-node-0", "--chainBServiceName", "node-service-archway-node-1", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between archway to archway with 1 custom chain parameter", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunCustomArchwayNode1(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "archway", "--chainB", "archway", "--chainBServiceName", "node-service-archway-node-1", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between neutron and neutron by running one custom neutron chain-1", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunNeutronNode(enclaveName)
			dive.RunCustomNeutronNode1(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "neutron", "--chainB", "neutron", "--chainAServiceName", "neutron-node-test-chain1", "--chainBServiceName", "neutron-node-test-chain3", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between 2 custom neutron chains", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunCustomNeutronNode0(enclaveName)
			dive.RunCustomNeutronNode1(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "neutron", "--chainB", "neutron", "--chainAServiceName", "neutron-node-test-chain2", "--chainBServiceName", "neutron-node-test-chain3", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start IBC relay between nuetron to neutron with one 1 custom chain.", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunCustomNeutronNode1(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "neutron", "--chainB", "neutron", "--chainBServiceName", "neutron-node-test-chain3", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between archway and neutron chains", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "archway", "--chainB", "neutron", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between already running archway and neutron chains", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunArchwayNode(enclaveName)
			dive.RunNeutronNode(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "archway", "--chainB", "neutron", "--chainAServiceName", "node-service-constantine-3", "--chainBServiceName", "neutron-node-test-chain1", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between already running archway and neutron chains with custom configuration", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunCustomNeutronNode0(enclaveName)
			dive.RunCustomArchwayNode0(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "archway", "--chainB", "neutron", "--chainAServiceName", "node-service-archway-node-0", "--chainBServiceName", "neutron-node-test-chain2", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start IBC relay between icon and archway", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "icon", "--chainB", "archway", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start IBC relay between icon and neutron-1", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "icon", "--chainB", "neutron", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start IBC relay between already running icon and archway chain-1", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunIconNode(enclaveName)
			dive.RunArchwayNode(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "icon", "--chainB", "archway", "--chainAServiceName", "icon-node-0xacbc4e", "--chainBServiceName", "node-service-constantine-3", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start IBC relay between already running icon and neutron chain", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunIconNode(enclaveName)
			dive.RunNeutronNode(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "icon", "--chainB", "neutron", "--chainAServiceName", "icon-node-0xacbc4e", "--chainBServiceName", "neutron-node-test-chain1", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start IBC relay between already running icon and custom archway chain", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunIconNode(enclaveName)
			dive.RunCustomArchwayNode0(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "icon", "--chainB", "archway", "--chainAServiceName", "icon-node-0xacbc4e", "--chainBServiceName", "node-service-archway-node-0", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start IBC relay between already running icon and custom neutron chain", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunIconNode(enclaveName)
			dive.RunCustomNeutronNode0(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "icon", "--chainB", "neutron", "--chainAServiceName", "icon-node-0xacbc4e", "--chainBServiceName", "neutron-node-test-chain2", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start IBC relay between already running custom icon and archway chain", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunCustomIconNode(enclaveName)
			dive.RunArchwayNode(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "icon", "--chainB", "archway", "--chainAServiceName", "icon-node-0xacbc4e", "--chainBServiceName", "node-service-constantine-3", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start IBC relay between already running custom icon and neutron chain", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunCustomIconNode(enclaveName)
			dive.RunNeutronNode(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "icon", "--chainB", "neutron", "--chainAServiceName", "icon-node-0x42f1f3", "--chainBServiceName", "neutron-node-test-chain1", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start IBC relay between already running custom icon and custom archway chain", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunCustomIconNode(enclaveName)
			dive.RunCustomArchwayNode0(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "icon", "--chainB", "archway", "--chainAServiceName", "icon-node-0x42f1f3", "--chainBServiceName", "node-service-archway-node-0", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start IBC relay between already running custom icon and custom neutron chain", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunCustomIconNode(enclaveName)
			dive.RunCustomNeutronNode0(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "icon", "--chainB", "neutron", "--chainAServiceName", "icon-node-0x42f1f3", "--chainBServiceName", "neutron-node-test-chain2", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between icon and hardhat by running icon node first and running bridge command directly", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunDecentralizedCustomIconNode0(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "hardhat", "--chainAServiceName", "icon-node-0xacbc4e", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between icon and hardhat by running hardhat node first and running bridge command directly", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunHardhatNode(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "hardhat", "--chainB", "icon", "--chainAServiceName", "hardhat-node", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between icon and eth by running icon node first and running bridge command directly", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunDecentralizedCustomIconNode0(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "eth", "--chainAServiceName", "icon-node-0xacbc4e", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between icon and eth by running eth node first and running bridge command directly", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunEthNode(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "eth", "--chainB", "icon", "--chainAServiceName", "el-1-geth-lighthouse", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between icon and icon by running icon node first and running bridge command directly", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunDecentralizedCustomIconNode0(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "icon", "--chainAServiceName", "icon-node-0xacbc4e", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})
	})

	ginkgo.Describe("Other commands", func() {
		// ginkgo.It("should handle error when trying to clean if no enclaves are running", func() {
		// 	dive.Clean()
		// 	dive.Clean()
		// })

		ginkgo.It("should handle error when trying to clean if kurtosis engine is not running", func() {
			cmd1 := exec.Command("kurtosis", "engine", "stop")
			cmd1.Run()
			bin := dive.GetBinPath()
			cmd2 := exec.Command(bin, "clean")
			err := cmd2.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
			cmd3 := exec.Command("kurtosis", "engine", "start")
			cmd3.Run()
		})

		ginkgo.It("should handle invalid input for chain command", func() {
			cmd.Args = append(cmd.Args, "chain", "invalid_input")
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})
	})

	ginkgo.Describe("Icon chain commands", func() {
		ginkgo.It("should run single icon node testing", func() {
			time.Sleep(1 * time.Second)
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "icon", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should run single icon node with verbose flag enabled", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "icon", "--verbose", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should run single icon node along with decentralisation", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "icon", "-d", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should run single icon node along with decentralisation with verbose flag enabled", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "icon", "-d", "--verbose", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should run custom Icon node-0", func() {
			filepath := "../../cli/sample-jsons/config0.json"
			updated_path := dive.UpdatePublicPort(filepath)
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "icon", "-c", updated_path, "-g", "./config/genesis-icon-0.zip", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(updated_path)
			dive.Clean(enclaveName)
		})

		ginkgo.It("should run custom Icon node-0  with verbose flag enabled", func() {
			time.Sleep(3 * time.Second)
			filepath := "../../cli/sample-jsons/config0.json"
			updated_path := dive.UpdatePublicPort(filepath)
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "icon", "-c", updated_path, "-g", "./config/genesis-icon-0.zip", "--verbose", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(updated_path)
			dive.Clean(enclaveName)
		})

		ginkgo.It("should run custom Icon node-1", func() {
			filepath := "../../cli/sample-jsons/config1.json"
			updated_path := dive.UpdatePublicPort(filepath)
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "icon", "-c", updated_path, "-g", "./config/genesis-icon-1.zip", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(updated_path)
			dive.Clean(enclaveName)
		})

		ginkgo.It("should run custom Icon node-1  with verbose flag enabled", func() {
			filepath := "../../cli/sample-jsons/config1.json"
			time.Sleep(6 * time.Second)
			updated_path := dive.UpdatePublicPort(filepath)
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "icon", "-c", "../../cli/sample-jsons/config1.json", "-g", "./config/genesis-icon-1.zip", "--verbose", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(updated_path)
			dive.Clean(enclaveName)
		})

		ginkgo.It("should run icon node first and then decentralise it", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunIconNode(enclaveName)
			cmd.Args = append(cmd.Args, "chain", "icon", "decentralize", "-p", "gochain", "-k", "keystores/keystore.json", "-n", "0x3", "-e", "http://172.16.0.3:9080/api/v3/icon_dex", "-s", "icon-node-0xacbc4e", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should run icon node first and then decentralise it with verbose flag enabled", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunIconNode(enclaveName)
			cmd.Args = append(cmd.Args, "chain", "icon", "decentralize", "-p", "gochain", "-k", "keystores/keystore.json", "-n", "0x3", "-e", "http://172.16.0.3:9080/api/v3/icon_dex", "-s", "icon-node-0xacbc4e", "--verbose", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should handle invalid input for chain command", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "icon", "-c", "invalid.json", "-g", "./config/genesis-icon-0.zip", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should handle invalid input for chain command", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "icon", "-c", "../../cli/sample-jsons/config0.json", "-g", "./config/invalid-icon-3.zip", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should handle invalid input for chain command", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "icon", "-c", "../../cli/sample-jsons/invalid_config.json", "-g", "./config/invalid-icon-3.zip", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should handle invalid input for chain command", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunIconNode(enclaveName)
			cmd.Args = append(cmd.Args, "chain", "icon", "decentralize", "-p", "invalidPassword", "-k", "keystores/keystore.json", "-n", "0x3", "-e", "http://172.16.0.3:9080/api/v3/icon_dex", "-s", "icon-node-0xacbc4e", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should handle invalid input for chain command", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunIconNode(enclaveName)
			cmd.Args = append(cmd.Args, "chain", "icon", "decentralize", "-p", "gochain", "-k", "keystores/invalid.json", "-n", "0x3", "-e", "http://172.16.0.3:9080/api/v3/icon_dex", "-s", "icon-node-0xacbc4e", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should handle invalid input for chain command", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunIconNode(enclaveName)
			cmd.Args = append(cmd.Args, "chain", "icon", "decentralize", "-p", "gochain", "-k", "keystores/keystore.json", "-n", "0x9", "-e", "http://172.16.0.3:9080/api/v3/icon_dex", "-s", "icon-node-0xacbc4e", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should handle invalid input for chain command", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunIconNode(enclaveName)
			cmd.Args = append(cmd.Args, "chain", "icon", "decentralize", "-p", "gochain", "-k", "keystores/keystore.json", "-n", "0x3", "-e", "http://172.16.0.3:9081/api/v3/icon_dex", "-s", "icon-node-0xacbc4e", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should handle invalid input for chain command", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunIconNode(enclaveName)
			cmd.Args = append(cmd.Args, "chain", "icon", "decentralize", "-p", "gochain", "-k", "keystores/keystore.json", "-n", "0x3", "-e", "http://172.16.0.3:9080/api/v3/icon_dex", "-s", "icon-node", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should output user that chain is already running when trying to run icon chain that is already running", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunIconNode(enclaveName)
			cmd.Args = append(cmd.Args, "chain", "icon", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

	})

	ginkgo.Describe("Eth chain commands", func() {
		ginkgo.It("should run single eth node", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "eth", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should run single eth node with verbose flag enabled", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "eth", "--verbose", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should output user that chain is already running when trying to run eth chain that is already running", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunEthNode(enclaveName)
			cmd.Args = append(cmd.Args, "chain", "eth", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

	})

	ginkgo.Describe("Hardhat chain commands", func() {
		ginkgo.It("should run single hardhat node-1", func() {
			time.Sleep(3 * time.Second)
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "hardhat", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should run single hardhat node with verbose flag enabled", func() {
			time.Sleep(3 * time.Second)
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "hardhat", "--verbose", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should output user that chain is already running when trying to run hardhat chain that is already running", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunHardhatNode(enclaveName)
			cmd.Args = append(cmd.Args, "chain", "hardhat", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})
	})

	ginkgo.Describe("Archway chain commands", func() {
		ginkgo.It("should run single archway node", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "archway", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)

		})

		ginkgo.It("should run single archway node with verbose flag enabled", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "archway", "--verbose", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should run single custom archway node-1", func() {
			filepath1 := "../../cli/sample-jsons/archway.json"
			updated_path1 := dive.UpdatePublicPorts(filepath1)
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "archway", "-c", updated_path1, "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(updated_path1)
			dive.Clean(enclaveName)

		})

		ginkgo.It("should run single custom archway node with verbose flag enabled", func() {
			filepath1 := "../../cli/sample-jsons/archway.json"
			updated_path1 := dive.UpdatePublicPorts(filepath1)
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "archway", "-c", updated_path1, "--verbose", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(updated_path1)
			dive.Clean(enclaveName)

		})

		ginkgo.It("should run single custom archway node with invalid json path", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "archway", "-c", "../../cli/sample-jsons/invalid_archway.json", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
			dive.Clean(enclaveName)

		})

		ginkgo.It("should output user that chain is already running when trying to run archway chain that is already running", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunArchwayNode(enclaveName)
			cmd.Args = append(cmd.Args, "chain", "archway", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})
	})

	ginkgo.Describe("Neutron chain commands", func() {
		ginkgo.It("should run single neutron node", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "neutron", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)

		})

		ginkgo.It("should run single neutron node with verbose flag enabled", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "neutron", "--verbose", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)

		})

		ginkgo.It("should run single custom neutron node", func() {
			filepath2 := "../../cli/sample-jsons/neutron.json"
			updated_path2 := dive.UpdateNeutronPublicPorts(filepath2)
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "neutron", "-c", updated_path2, "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(updated_path2)
			dive.Clean(enclaveName)
		})

		ginkgo.It("should run single custom neutron node with verbose flag enabled", func() {
			filepath2 := "../../cli/sample-jsons/neutron.json"
			updated_path2 := dive.UpdateNeutronPublicPorts(filepath2)
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "neutron", "-c", updated_path2, "--verbose", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(updated_path2)
			dive.Clean(enclaveName)
		})

		ginkgo.It("should run single custom neutron node with invalid json path", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "neutron", "-c", "../../cli/sample-jsons/neutron5.json", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should output user that chain is already running when trying to run neutron chain that is already running", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunNeutronNode(enclaveName)
			cmd.Args = append(cmd.Args, "chain", "neutron", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})
	})

	ginkgo.Describe("Relaychain commands", func() {
		var selectedChain string

		if envChain := os.Getenv("relayChainName"); envChain != "" {
			selectedChain = envChain
		} else {
			selectedChain = "default" // Provide a default value if not set
		}

		relayChainNames := []string{"kusama", "polkadot"}

		for _, relayChainName := range relayChainNames {
			if selectedChain != "default" && selectedChain != relayChainName {
				continue // Skip tests if the selected chain doesn't match the loop chain
			}

			relayChainName := relayChainName // Capture the loop variable

			ginkgo.It("should run single relaychain with verbose flag enabled for "+relayChainName, func() {
				enclaveName := dive.GenerateRandomName()
				cmd.Args = append(cmd.Args, "chain", relayChainName, "--verbose", "--enclaveName", enclaveName)
				err := cmd.Run()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				dive.Clean(enclaveName)
			})

			ginkgo.It("should run single relaychain in localnet for "+relayChainName, func() {
				enclaveName := dive.GenerateRandomName()
				cmd.Args = append(cmd.Args, "chain", relayChainName, "-n", "local", "--enclaveName", enclaveName)
				err := cmd.Run()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				dive.Clean(enclaveName)
			})

			ginkgo.It("should run single relaychain in mainnet for "+relayChainName, func() {
				enclaveName := dive.GenerateRandomName()
				cmd.Args = append(cmd.Args, "chain", relayChainName, "-n", "mainnet", "--enclaveName", enclaveName)
				err := cmd.Run()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				dive.Clean(enclaveName)
			})

			ginkgo.It("should run single relaychain in testnet for "+relayChainName, func() {
				enclaveName := dive.GenerateRandomName()
				cmd.Args = append(cmd.Args, "chain", relayChainName, "-n", "testnet", "--enclaveName", enclaveName)
				err := cmd.Run()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				dive.Clean(enclaveName)
			})

			ginkgo.It("should run custom relaychain in localnet with verbose enabled for "+relayChainName, func() {
				enclaveName := dive.GenerateRandomName()
				config := dive.NewConfigurationWithChainType("local")

				// Convert the config struct to a JSON-formatted string
				configString, err := json.Marshal(config)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				cmd.Args = append(cmd.Args, "chain", relayChainName, "-c", string(configString), "--verbose", "--enclaveName", enclaveName)
				err = cmd.Run()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				dive.Clean(enclaveName)
			})

			

			ginkgo.It("should run single relaychain with explorer service for "+relayChainName, func() {
				enclaveName := dive.GenerateRandomName()
				cmd.Args = append(cmd.Args, "chain", relayChainName, "--explorer", "--enclaveName", enclaveName)
				err := cmd.Run()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				dive.Clean(enclaveName)
			})

			ginkgo.It("should run single relaychain with metrics service for "+relayChainName, func() {
				enclaveName := dive.GenerateRandomName()
				cmd.Args = append(cmd.Args, "chain", relayChainName, "--metrics", "--enclaveName", enclaveName)
				err := cmd.Run()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				dive.Clean(enclaveName)
			})

			ginkgo.It("should run custom relaychain with explorer service for "+relayChainName, func() {
				enclaveName := dive.GenerateRandomName()
				cmd.Args = append(cmd.Args, "chain", relayChainName, "-c", "./sample-jsons/local.json", "--enclaveName", enclaveName)
				err := cmd.Run()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				dive.Clean(enclaveName)
			})

			ginkgo.It("should run custom relaychain with metrics service for "+relayChainName, func() {
				enclaveName := dive.GenerateRandomName()
				cmd.Args = append(cmd.Args, "chain", relayChainName, "-c", "./sample-jsons/local.json", "--enclaveName", enclaveName)
				err := cmd.Run()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				dive.Clean(enclaveName)
			})
		}
	})

	

	ginkgo.Describe("Parachain commands", func() {
		var selectedRelayChain string
		var selectedParaChain string

		if envChain := os.Getenv("relayChainName"); envChain != "" {
			selectedRelayChain = envChain
		} else {
			selectedRelayChain = "default" // Provide a default value if not set
		}

		if envSelectedParaChain := os.Getenv("paraChainName"); envSelectedParaChain != "" {
			selectedParaChain = envSelectedParaChain
		} else {
			// If not provided, set a default value
			selectedParaChain = "default"
		}

		relayChainNames := []string{"kusama", "polkadot"}
		for _, relayChainName := range relayChainNames {
			if selectedRelayChain != "default" && selectedRelayChain != relayChainName {
				continue // Skip tests if the selected chain doesn't match the loop chain
			}

			relayChainName := relayChainName // Capture the loop variable

			var paraChainNames []string

			if relayChainName == "kusama" {
				paraChainNames = []string{"Karura", "Kintsugi-BTC", "Altair", "Bifrost", "Mangata", "Robonomics", "Turing Network", "Encointer Network", "Bajun Networkc", "Calamari", "Khala Network", "Litmus", "Moonriver", "subzero"}
			} else if relayChainName == "polkadot" {
				paraChainNames = []string{"Polkadex", "Zeitgeist", "Acala", "Bifrost", "Clover", "Integritee Shell", "Integritee Shell", "Litentry", "Moonbeam", "Nodle", "Pendulum", "Ajuna Network", "Centrifuge", "Frequency", "Interlay", "Kylin", "Manta", "Moonsama", "Parallel", "Phala Network", "Subsocial"}
			} else {
				paraChainNames = []string{"Polkadex", "Zeitgeist", "Karura", "Kintsugi-BTC", "Altair", "Bifrost", "Mangata", "Robonomics", "Turing Network", "Encointer Network", "Bajun Networkc", "Calamari", "Khala Network", "Litmus", "Moonriver", "subzero", "Acala", "Bifrost", "Clover", "Integritee Shell", "Integritee Shell", "Litentry", "Moonbeam", "Nodle", "Pendulum", "Ajuna Network", "Centrifuge", "Frequency", "Interlay", "Kylin", "Manta", "Moonsama", "Parallel", "Phala Network", "Subsocial"}
			}

			for _, paraChainName := range paraChainNames {
				if selectedParaChain != "default" && selectedParaChain != paraChainName {
					continue // Skip tests if the selected parachain doesn't match the loop parachain
				}

				ginkgo.It("should run single parachain  in testnet with verbose flag enabled for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := dive.GenerateRandomName()
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-p", paraChainName, "--no-relay", "-n", "testnet", "--vebose", "--enclaveName", enclaveName)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					dive.Clean(enclaveName)
				})
				ginkgo.It("should run single parachain in mainnet for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := dive.GenerateRandomName()
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-p", paraChainName, "--no-relay", "-n", "mainnet", "--enclaveName", enclaveName)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					dive.Clean(enclaveName)
				})
				ginkgo.It("should run single parachain in mainnet with explorer services for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := dive.GenerateRandomName()
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-p", paraChainName, "--no-relay", "-n", "mainnet", "--explorer", "--enclaveName", enclaveName)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					dive.Clean(enclaveName)
				})
				ginkgo.It("should run single parachain in mainnet with metrics services for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := dive.GenerateRandomName()
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-p", paraChainName, "--no-relay", "-n", "mainnet", "--metrics", "--enclaveName", enclaveName)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					dive.Clean(enclaveName)
				})
				ginkgo.It("should run custom parachain in testnet with verbose flag enabled for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := dive.GenerateRandomName()
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-p", paraChainName, "--no-relay", "-c", "./sample-jsons/local.json", "--enclaveName", enclaveName)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					dive.Clean(enclaveName)
				})
				ginkgo.It("should run custom parachain in mainnet for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := dive.GenerateRandomName()
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-p", paraChainName, "--no-relay", "-c", "./sample-jsons/local.json", "--enclaveName", enclaveName)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					dive.Clean(enclaveName)
				})
				ginkgo.It("should run custom parachain in mainnet  with explorer services for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := dive.GenerateRandomName()
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-p", paraChainName, "--no-relay", "-c", "./sample-jsons/local.json", "--enclaveName", enclaveName)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					dive.Clean(enclaveName)
				})
				ginkgo.It("should run custom parachain in mainnet  with metrics services for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := dive.GenerateRandomName()
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-p", paraChainName, "--no-relay", "-c", "./sample-jsons/local.json", "--enclaveName", enclaveName)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					dive.Clean(enclaveName)
				})
			}
		}
	})

	ginkgo.Describe("Relaychain and parachain commands", func() {
		var selectedRelayChain string
		var selectedParaChain string

		if envChain := os.Getenv("relayChainName"); envChain != "" {
			selectedRelayChain = envChain
		} else {
			selectedRelayChain = "default" // Provide a default value if not set
		}

		if envSelectedParaChain := os.Getenv("paraChainName"); envSelectedParaChain != "" {
			selectedParaChain = envSelectedParaChain
		} else {
			// If not provided, set a default value
			selectedParaChain = "default"
		}

		relayChainNames := []string{"kusama", "polkadot"}
		for _, relayChainName := range relayChainNames {
			if selectedRelayChain != "default" && selectedRelayChain != relayChainName {
				continue // Skip tests if the selected chain doesn't match the loop chain
			}

			relayChainName := relayChainName // Capture the loop variable

			var paraChainNames []string

			if relayChainName == "kusama" {
				paraChainNames = []string{"Karura", "Kintsugi-BTC", "Altair", "Bifrost", "Mangata", "Robonomics", "Turing Network", "Encointer Network", "Bajun Networkc", "Calamari", "Khala Network", "Litmus", "Moonriver", "subzero"}
			} else if relayChainName == "polkadot" {
				paraChainNames = []string{"Polkadex", "Zeitgeist", "Acala", "Bifrost", "Clover", "Integritee Shell", "Integritee Shell", "Litentry", "Moonbeam", "Nodle", "Pendulum", "Ajuna Network", "Centrifuge", "Frequency", "Interlay", "Kylin", "Manta", "Moonsama", "Parallel", "Phala Network", "Subsocial"}
			} else {
				paraChainNames = []string{"Polkadex", "Zeitgeist", "Karura", "Kintsugi-BTC", "Altair", "Bifrost", "Mangata", "Robonomics", "Turing Network", "Encointer Network", "Bajun Networkc", "Calamari", "Khala Network", "Litmus", "Moonriver", "subzero", "Acala", "Bifrost", "Clover", "Integritee Shell", "Integritee Shell", "Litentry", "Moonbeam", "Nodle", "Pendulum", "Ajuna Network", "Centrifuge", "Frequency", "Interlay", "Kylin", "Manta", "Moonsama", "Parallel", "Phala Network", "Subsocial"}
			}

			for _, paraChainName := range paraChainNames {
				if selectedParaChain != "default" && selectedParaChain != paraChainName {
					continue // Skip tests if the selected parachain doesn't match the loop parachain
				}

				ginkgo.It("should run single relaychain and parachain  in testnet with verbose flag enabled for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := dive.GenerateRandomName()
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-p", paraChainName, "-n", "testnet", "--vebose", "--enclaveName", enclaveName)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					dive.Clean(enclaveName)
				})
				ginkgo.It("should run single relaychain and parachain in mainnet for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := dive.GenerateRandomName()
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-p", paraChainName, "-n", "mainnet", "--enclaveName", enclaveName)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					dive.Clean(enclaveName)
				})
				ginkgo.It("should run single relaychain and parachain in local for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := dive.GenerateRandomName()
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-p", paraChainName, "-n", "local", "--enclaveName", enclaveName)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					dive.Clean(enclaveName)
				})
				ginkgo.It("should run single relaychain and parachain in mainnet with explorer services for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := dive.GenerateRandomName()
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-p", paraChainName,  "-n", "mainnet", "--explorer", "--enclaveName", enclaveName)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					dive.Clean(enclaveName)
				})
				ginkgo.It("should run single relaychain and  parachain in mainnet with metrics services for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := dive.GenerateRandomName()
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-p", paraChainName, "-n", "mainnet", "--metrics", "--enclaveName", enclaveName)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					dive.Clean(enclaveName)
				})
				ginkgo.It("should run custom relaychain and parachain in testnet with verbose flag enabled for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := dive.GenerateRandomName()
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-p", paraChainName, "-c", "./sample-jsons/local.json", "--enclaveName", enclaveName)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					dive.Clean(enclaveName)
				})
				ginkgo.It("should run custom relaychain and  parachain in mainnet for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := dive.GenerateRandomName()
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-p", paraChainName, "-c", "./sample-jsons/local.json", "--enclaveName", enclaveName)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					dive.Clean(enclaveName)
				})
				ginkgo.It("should run custom relaychain and  parachain in local for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := dive.GenerateRandomName()
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-p", paraChainName, "-c", "./sample-jsons/local.json", "--enclaveName", enclaveName)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					dive.Clean(enclaveName)
				})
				ginkgo.It("should run custom relaychain and  parachain in mainnet  with explorer services for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := dive.GenerateRandomName()
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-p", paraChainName, "-c", "./sample-jsons/local.json", "--enclaveName", enclaveName)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					dive.Clean(enclaveName)
				})
				ginkgo.It("should run custom relaychain and parachain in mainnet  with metrics services for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := dive.GenerateRandomName()
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-p", paraChainName, "-c", "./sample-jsons/local.json", "--enclaveName", enclaveName)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					dive.Clean(enclaveName)
				})
			}
		}
	})
})


