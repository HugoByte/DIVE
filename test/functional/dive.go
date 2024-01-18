package dive

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"testing"

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
		cmd = GetBinaryCommand()
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
			enclaveName := GenerateRandomName()
			cli := common.GetCli(enclaveName)
			latestVersion := utility.GetLatestVersion(cli)
			gomega.Expect(stdout.String()).To(gomega.ContainSubstring(latestVersion))
		})

		ginkgo.It("should start bridge between icon and eth correctly", func() {
			enclaveName := GenerateRandomName()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "eth", "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between icon and hardhat but with icon bridge set to true", func() {
			enclaveName := GenerateRandomName()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "hardhat", "--bmvbridge", "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between icon and icon", func() {
			enclaveName := GenerateRandomName()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "icon", "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between archway and archway using ibc", func() {
			enclaveName := GenerateRandomName()
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "archway", "--chainB", "archway", "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

	})

	ginkgo.Describe("Bridge command Test", func() {
		ginkgo.It("should start bridge between icon and eth but with icon bridge set to true", func() {
			enclaveName := GenerateRandomName()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "eth", "--bmvbridge", "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between icon and hardhat", func() {
			enclaveName := GenerateRandomName()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "hardhat", "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between icon and eth by running each chain individually", func() {
			enclaveName := GenerateRandomName()
			RunDecentralizedIconNode(enclaveName)
			RunEthNode(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "eth", "--chainAServiceName", ICON_CONFIG0_SERVICENAME, "--chainBServiceName", ETH_SERVICENAME, "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between icon and hardhat by running each chain individually", func() {
			enclaveName := GenerateRandomName()
			RunDecentralizedIconNode(enclaveName)
			RunHardhatNode(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "hardhat", "--chainAServiceName", ICON_CONFIG0_SERVICENAME, "--chainBServiceName", HARDHAT_SERVICENAME, "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between icon and eth by running icon node first and then decentralising it", func() {
			enclaveName := GenerateRandomName()
			RunIconNode(enclaveName)

			serviceName0, endpoint0, nid0 := GetServiceDetails(fmt.Sprintf("services_%s.json", enclaveName), ICON_CONFIG0_SERVICENAME)
			DecentralizeCustomIconNode(nid0, endpoint0, serviceName0, enclaveName)

			RunEthNode(enclaveName)

			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "eth", "--chainAServiceName", ICON_CONFIG0_SERVICENAME, "--chainBServiceName", ETH_SERVICENAME, "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(fmt.Sprintf("updated-config-%s.json", enclaveName))
		})

		ginkgo.It("should start bridge between icon and hardhat by running icon node first and then decentralising it", func() {
			enclaveName := GenerateRandomName()

			RunIconNode(enclaveName)
			serviceName0, endpoint0, nid0 := GetServiceDetails(fmt.Sprintf("services_%s.json", enclaveName), ICON_CONFIG0_SERVICENAME)
			DecentralizeCustomIconNode(nid0, endpoint0, serviceName0, enclaveName)

			RunHardhatNode(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "hardhat", "--chainAServiceName", ICON_CONFIG0_SERVICENAME, "--chainBServiceName", HARDHAT_SERVICENAME, "--enclaveName", enclaveName)
			Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(fmt.Sprintf("updated-config-%s.json", enclaveName))
		})

		ginkgo.It("should start bridge between icon and icon by running one custom icon chain", func() {
			enclaveName := GenerateRandomName()
			RunDecentralizedIconNode(enclaveName)
			RunDecentralizedCustomIconNode1(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "icon", "--chainAServiceName", ICON_CONFIG0_SERVICENAME, "--chainBServiceName", ICON_CONFIG1_SERVICENAME, "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between icon and running custom icon later decentralising it", func() {
			enclaveName := GenerateRandomName()
			RunDecentralizedIconNode(enclaveName)
			RunCustomIconNode1(enclaveName)

			serviceName, endpoint, nid := GetServiceDetails(fmt.Sprintf("services_%s.json", enclaveName), ICON_CONFIG1_SERVICENAME)
			DecentralizeCustomIconNode(nid, endpoint, serviceName, enclaveName)

			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "icon", "--chainAServiceName", ICON_CONFIG0_SERVICENAME, "--chainBServiceName", ICON_CONFIG1_SERVICENAME, "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(fmt.Sprintf("updated-config-%s.json", enclaveName))
		})

		ginkgo.It("should start bridge between icon and icon by running one icon chain and later decentralsing it. Running another custom icon chain and then decentralising it", func() {
			enclaveName := GenerateRandomName()
			RunIconNode(enclaveName)
			RunCustomIconNode1(enclaveName)

			serviceName0, endpoint0, nid0 := GetServiceDetails(fmt.Sprintf("services_%s.json", enclaveName), ICON_CONFIG0_SERVICENAME)
			DecentralizeCustomIconNode(nid0, endpoint0, serviceName0, enclaveName)

			serviceName1, endpoint1, nid1 := GetServiceDetails(fmt.Sprintf("services_%s.json", enclaveName), ICON_CONFIG1_SERVICENAME)
			DecentralizeCustomIconNode(nid1, endpoint1, serviceName1, enclaveName)

			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "icon", "--chainAServiceName", ICON_CONFIG0_SERVICENAME, "--chainBServiceName", ICON_CONFIG1_SERVICENAME, "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(fmt.Sprintf("updated-config-%s.json", enclaveName))
		})

		ginkgo.It("should start bridge between 2 custom icon chains", func() {
			enclaveName := GenerateRandomName()
			RunDecentralizedCustomIconNode0(enclaveName)
			RunDecentralizedCustomIconNode1(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "icon", "--chainAServiceName", ICON_CONFIG0_SERVICENAME, "--chainBServiceName", ICON_CONFIG1_SERVICENAME, "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between 2 custom icon chains by running them first and then decentralising it later", func() {
			enclaveName := GenerateRandomName()
			RunCustomIconNode0(enclaveName)
			RunCustomIconNode1(enclaveName)

			serviceName0, endpoint0, nid0 := GetServiceDetails(fmt.Sprintf("services_%s.json", enclaveName), ICON_CONFIG0_SERVICENAME)
			DecentralizeCustomIconNode(nid0, endpoint0, serviceName0, enclaveName)

			serviceName1, endpoint1, nid1 := GetServiceDetails(fmt.Sprintf("services_%s.json", enclaveName), ICON_CONFIG1_SERVICENAME)
			DecentralizeCustomIconNode(nid1, endpoint1, serviceName1, enclaveName)

			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "icon", "--chainAServiceName", ICON_CONFIG0_SERVICENAME, "--chainBServiceName", ICON_CONFIG1_SERVICENAME, "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(fmt.Sprintf("updated-config-%s.json", enclaveName))
		})

		ginkgo.It("should start bridge between 2 chains when all nodes are running", func() {
			enclaveName := GenerateRandomName()
			RunDecentralizedIconNode(enclaveName)
			RunEthNode(enclaveName)
			RunHardhatNode(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "eth", "--chainAServiceName", ICON_CONFIG0_SERVICENAME, "--chainBServiceName", ETH_SERVICENAME, "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should handle invalid input for bridge command", func() {
			enclaveName := GenerateRandomName()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "invalid_input", "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("should handle invalid input for bridge command", func() {
			enclaveName := GenerateRandomName()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "invalid_input", "--chainB", "eth", "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("should handle invalid input ibc bridge command", func() {
			enclaveName := GenerateRandomName()
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "archway", "--chainB", "invalid", "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between archway and archway by running one custom archway chain", func() {
			enclaveName := GenerateRandomName()
			RunArchwayNode(enclaveName)
			RunCustomArchwayNode1(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "archway", "--chainB", "archway", "--chainAServiceName", DEFAULT_ARCHWAY_SERVICENAME, "--chainBServiceName", ARCHWAY_CONFIG1_SERVICENAME, "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(fmt.Sprintf("updated-config-%s.json", enclaveName))
		})

		ginkgo.It("should start bridge between 2 custom archway chains", func() {
			enclaveName := GenerateRandomName()
			RunCustomArchwayNode0(enclaveName)
			RunCustomArchwayNode1(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "archway", "--chainB", "archway", "--chainAServiceName", ARCHWAY_CONFIG0_SERVICENAME, "--chainBServiceName", ARCHWAY_CONFIG1_SERVICENAME, "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(fmt.Sprintf("updated-config-%s.json", enclaveName))
		})

		ginkgo.It("should start bridge between archway to archway with 1 custom chain parameter", func() {
			enclaveName := GenerateRandomName()
			RunCustomArchwayNode1(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "archway", "--chainB", "archway", "--chainBServiceName", ARCHWAY_CONFIG1_SERVICENAME, "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(fmt.Sprintf("updated-config-%s.json", enclaveName))
		})

		ginkgo.It("should start bridge between neutron and neutron by running one custom neutron chain", func() {
			enclaveName := GenerateRandomName()
			RunNeutronNode(enclaveName)
			RunCustomNeutronNode1(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "neutron", "--chainB", "neutron", "--chainAServiceName", DEFAULT_NEUTRON_SERVICENAME, "--chainBServiceName", NEUTRON_CONFIG1_SERVICENAME, "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(fmt.Sprintf("updated-config-%s.json", enclaveName))
		})

		ginkgo.It("should start bridge between 2 custom neutron chains", func() {
			enclaveName := GenerateRandomName()
			RunCustomNeutronNode0(enclaveName)
			RunCustomNeutronNode1(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "neutron", "--chainB", "neutron", "--chainAServiceName", NEUTRON_CONFIG0_SERVICENAME, "--chainBServiceName", NEUTRON_CONFIG1_SERVICENAME, "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(fmt.Sprintf("updated-config-%s.json", enclaveName))
		})

		ginkgo.It("should start IBC relay between nuetron to neutron with one 1 custom chain.", func() {
			enclaveName := GenerateRandomName()
			RunCustomNeutronNode1(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "neutron", "--chainB", "neutron", "--chainBServiceName", NEUTRON_CONFIG1_SERVICENAME, "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(fmt.Sprintf("updated-config-%s.json", enclaveName))
		})

		ginkgo.It("should start bridge between archway and neutron chains", func() {
			enclaveName := GenerateRandomName()
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "archway", "--chainB", "neutron", "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between already running archway and neutron chains", func() {
			enclaveName := GenerateRandomName()
			RunArchwayNode(enclaveName)
			RunNeutronNode(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "archway", "--chainB", "neutron", "--chainAServiceName", DEFAULT_ARCHWAY_SERVICENAME, "--chainBServiceName", DEFAULT_NEUTRON_SERVICENAME, "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between already running archway and neutron chains with custom configuration", func() {
			enclaveName := GenerateRandomName()
			RunCustomNeutronNode0(enclaveName)
			RunCustomArchwayNode0(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "archway", "--chainB", "neutron", "--chainAServiceName", ARCHWAY_CONFIG0_SERVICENAME, "--chainBServiceName", NEUTRON_CONFIG0_SERVICENAME, "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(fmt.Sprintf("updated-config-%s.json", enclaveName))
		})

		ginkgo.It("should start IBC relay between icon and archway", func() {
			enclaveName := GenerateRandomName()
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "icon", "--chainB", "archway", "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start IBC relay between icon and neutron", func() {
			enclaveName := GenerateRandomName()
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "icon", "--chainB", "neutron", "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start IBC relay between already running icon and archway chain", func() {
			enclaveName := GenerateRandomName()
			RunIconNode(enclaveName)
			RunArchwayNode(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "icon", "--chainB", "archway", "--chainAServiceName", ICON_CONFIG0_SERVICENAME, "--chainBServiceName", DEFAULT_ARCHWAY_SERVICENAME, "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start IBC relay between already running icon and neutron chain", func() {
			enclaveName := GenerateRandomName()
			RunIconNode(enclaveName)
			RunNeutronNode(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "icon", "--chainB", "neutron", "--chainAServiceName", ICON_CONFIG0_SERVICENAME, "--chainBServiceName", DEFAULT_NEUTRON_SERVICENAME, "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start IBC relay between already running icon and custom archway chain", func() {
			enclaveName := GenerateRandomName()
			RunIconNode(enclaveName)
			RunCustomArchwayNode0(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "icon", "--chainB", "archway", "--chainAServiceName", ICON_CONFIG0_SERVICENAME, "--chainBServiceName", ARCHWAY_CONFIG0_SERVICENAME, "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(fmt.Sprintf("updated-config-%s.json", enclaveName))
		})

		ginkgo.It("should start IBC relay between already running icon and custom neutron chain", func() {
			enclaveName := GenerateRandomName()
			RunIconNode(enclaveName)
			RunCustomNeutronNode0(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "icon", "--chainB", "neutron", "--chainAServiceName", ICON_CONFIG0_SERVICENAME, "--chainBServiceName", NEUTRON_CONFIG0_SERVICENAME, "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(fmt.Sprintf("updated-config-%s.json", enclaveName))
		})

		ginkgo.It("should start IBC relay between already running custom icon and archway chain", func() {
			enclaveName := GenerateRandomName()
			RunCustomIconNode0(enclaveName)
			RunArchwayNode(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "icon", "--chainB", "archway", "--chainAServiceName", ICON_CONFIG0_SERVICENAME, "--chainBServiceName", DEFAULT_ARCHWAY_SERVICENAME, "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(fmt.Sprintf("updated-config-%s.json", enclaveName))
		})

		ginkgo.It("should start IBC relay between already running custom icon and neutron chain", func() {
			enclaveName := GenerateRandomName()
			RunCustomIconNode1(enclaveName)
			RunNeutronNode(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "icon", "--chainB", "neutron", "--chainAServiceName", ICON_CONFIG1_SERVICENAME, "--chainBServiceName", DEFAULT_NEUTRON_SERVICENAME, "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(fmt.Sprintf("updated-config-%s.json", enclaveName))
		})

		ginkgo.It("should start IBC relay between already running custom icon and custom archway chain", func() {
			enclaveName := GenerateRandomName()
			RunCustomIconNode1(enclaveName)
			RunCustomArchwayNode0(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "icon", "--chainB", "archway", "--chainAServiceName", ICON_CONFIG1_SERVICENAME, "--chainBServiceName", ARCHWAY_CONFIG0_SERVICENAME, "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(fmt.Sprintf("updated-config-%s.json", enclaveName))
		})

		ginkgo.It("should start IBC relay between already running custom icon and custom neutron chain", func() {
			enclaveName := GenerateRandomName()
			RunCustomIconNode1(enclaveName)
			RunCustomNeutronNode0(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "icon", "--chainB", "neutron", "--chainAServiceName", ICON_CONFIG1_SERVICENAME, "--chainBServiceName", NEUTRON_CONFIG0_SERVICENAME, "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(fmt.Sprintf("updated-config-%s.json", enclaveName))
		})

		ginkgo.It("should start bridge between icon and hardhat by running icon node first and running bridge command directly", func() {
			enclaveName := GenerateRandomName()
			RunDecentralizedCustomIconNode0(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "hardhat", "--chainAServiceName", ICON_CONFIG0_SERVICENAME, "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(fmt.Sprintf("updated-config-%s.json", enclaveName))
		})

		ginkgo.It("should start bridge between icon and hardhat by running hardhat node first and running bridge command directly", func() {
			enclaveName := GenerateRandomName()
			RunHardhatNode(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "hardhat", "--chainB", "icon", "--chainAServiceName", HARDHAT_SERVICENAME, "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between icon and eth by running icon node first and running bridge command directly", func() {
			enclaveName := GenerateRandomName()
			RunDecentralizedCustomIconNode0(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "eth", "--chainAServiceName", ICON_CONFIG0_SERVICENAME, "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(fmt.Sprintf("updated-config-%s.json", enclaveName))
		})

		ginkgo.It("should start bridge between icon and eth by running eth node first and running bridge command directly", func() {
			enclaveName := GenerateRandomName()
			RunEthNode(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "eth", "--chainB", "icon", "--chainAServiceName", ETH_SERVICENAME, "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between icon and icon by running icon node first and running bridge command directly", func() {
			enclaveName := GenerateRandomName()
			RunDecentralizedCustomIconNode1(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "icon", "--chainAServiceName", ICON_CONFIG1_SERVICENAME, "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(fmt.Sprintf("updated-config-%s.json", enclaveName))
		})
	})

	ginkgo.Describe("Other commands", func() {

		ginkgo.It("should handle error when trying to clean if kurtosis engine is not running", func() {
			cmd1 := exec.Command("kurtosis", "engine", "stop")
			cmd1.Run()
			bin := GetBinPath()
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
			enclaveName := GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "icon", "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should run single icon node along with decentralisation", func() {
			enclaveName := GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "icon", "-d", "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should run custom Icon node-0", func() {
			enclaveName := GenerateRandomName()
			updated_path := UpdatePublicPort(enclaveName, ICON_CONFIG0)
			cmd.Args = append(cmd.Args, "chain", "icon", "-c", updated_path, "-g", ICON_GENESIS0, "--enclaveName", enclaveName)
			defer os.Remove(updated_path)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should run custom Icon node-1", func() {
			enclaveName := GenerateRandomName()
			filepath := ICON_CONFIG1
			updated_path := UpdatePublicPort(enclaveName, filepath)
			cmd.Args = append(cmd.Args, "chain", "icon", "-c", updated_path, "-g", ICON_GENESIS1, "--enclaveName", enclaveName)
			defer os.Remove(updated_path)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should run icon node first and then decentralise it", func() {
			enclaveName := GenerateRandomName()
			RunIconNode(enclaveName)
			serviceName0, endpoint0, nid0 := GetServiceDetails(fmt.Sprintf("services_%s.json", enclaveName), ICON_CONFIG0_SERVICENAME)
			cmd.Args = append(cmd.Args, "chain", "icon", "decentralize", "-p", "gochain", "-k", "keystores/keystore.json", "-n", nid0, "-e", endpoint0, "-s", serviceName0, "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should handle invalid input for chain command", func() {
			enclaveName := GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "icon", "-c", "invalid.json", "-g", ICON_GENESIS0, "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("should handle invalid input for chain command", func() {
			enclaveName := GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "icon", "-c", ICON_CONFIG0, "-g", "./config/invalid-icon-3.zip", "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("should handle invalid input for chain command", func() {
			enclaveName := GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "icon", "-c", "../../cli/sample-jsons/invalid_config.json", "-g", "./config/invalid-icon-3.zip", "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("should handle invalid input for chain command", func() {
			enclaveName := GenerateRandomName()
			RunIconNode(enclaveName)
			serviceName0, endpoint0, nid0 := GetServiceDetails(fmt.Sprintf("services_%s.json", enclaveName), ICON_CONFIG0_SERVICENAME)
			cmd.Args = append(cmd.Args, "chain", "icon", "decentralize", "-p", "invalidPassword", "-k", "keystores/keystore.json", "-n", nid0, "-e", endpoint0, "-s", serviceName0, "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("should handle invalid input for chain command", func() {
			enclaveName := GenerateRandomName()
			RunIconNode(enclaveName)
			serviceName0, endpoint0, nid0 := GetServiceDetails(fmt.Sprintf("services_%s.json", enclaveName), ICON_CONFIG0_SERVICENAME)
			cmd.Args = append(cmd.Args, "chain", "icon", "decentralize", "-p", "gochain", "-k", "keystores/invalid.json", "-n", nid0, "-e", endpoint0, "-s", serviceName0, "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("should handle invalid input for chain command", func() {
			enclaveName := GenerateRandomName()
			RunIconNode(enclaveName)
			serviceName0, endpoint0, _ := GetServiceDetails(fmt.Sprintf("services_%s.json", enclaveName), ICON_CONFIG0_SERVICENAME)
			cmd.Args = append(cmd.Args, "chain", "icon", "decentralize", "-p", "gochain", "-k", "keystores/keystore.json", "-n", "0x9", "-e", endpoint0, "-s", serviceName0, "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("should handle invalid input for chain command", func() {
			enclaveName := GenerateRandomName()
			RunIconNode(enclaveName)
			serviceName0, _, nid0 := GetServiceDetails(fmt.Sprintf("services_%s.json", enclaveName), ICON_CONFIG0_SERVICENAME)
			cmd.Args = append(cmd.Args, "chain", "icon", "decentralize", "-p", "gochain", "-k", "keystores/keystore.json", "-n", nid0, "-e", "http://172.16.0.3:9081/api/v3/icon_dex", "-s", serviceName0, "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("should handle invalid input for chain command", func() {
			enclaveName := GenerateRandomName()
			RunIconNode(enclaveName)
			_, endpoint0, nid0 := GetServiceDetails(fmt.Sprintf("services_%s.json", enclaveName), ICON_CONFIG0_SERVICENAME)
			cmd.Args = append(cmd.Args, "chain", "icon", "decentralize", "-p", "gochain", "-k", "keystores/keystore.json", "-n", nid0, "-e", endpoint0, "-s", "icon-node", "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("should output user that chain is already running when trying to run icon chain that is already running", func() {
			enclaveName := GenerateRandomName()
			RunIconNode(enclaveName)
			cmd.Args = append(cmd.Args, "chain", "icon", "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

	})

	ginkgo.Describe("Eth chain commands", func() {
		ginkgo.It("should run single eth node", func() {
			enclaveName := GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "eth", "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should output user that chain is already running when trying to run eth chain that is already running", func() {
			enclaveName := GenerateRandomName()
			RunEthNode(enclaveName)
			cmd.Args = append(cmd.Args, "chain", "eth", "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

	})

	ginkgo.Describe("Hardhat chain commands", func() {
		ginkgo.It("should run single hardhat node-1", func() {
			enclaveName := GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "hardhat", "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should output user that chain is already running when trying to run hardhat chain that is already running", func() {
			enclaveName := GenerateRandomName()
			RunHardhatNode(enclaveName)
			cmd.Args = append(cmd.Args, "chain", "hardhat", "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})
	})

	ginkgo.Describe("Archway chain commands", func() {
		ginkgo.It("should run single archway node", func() {
			enclaveName := GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "archway", "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should run single custom archway node-1", func() {
			filepath1 := "../../cli/sample-jsons/archway.json"
			updated_path1 := UpdatePublicPorts(filepath1)
			enclaveName := GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "archway", "-c", updated_path1, "--enclaveName", enclaveName)
			defer os.Remove(updated_path1)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should run single custom archway node with invalid json path", func() {
			enclaveName := GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "archway", "-c", "../../cli/sample-jsons/invalid_archway.json", "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
			Clean(enclaveName)
		})

		ginkgo.It("should output user that chain is already running when trying to run archway chain that is already running", func() {
			enclaveName := GenerateRandomName()
			RunArchwayNode(enclaveName)
			cmd.Args = append(cmd.Args, "chain", "archway", "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})
	})

	ginkgo.Describe("Neutron chain commands", func() {
		ginkgo.It("should run single neutron node", func() {
			enclaveName := GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "neutron", "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should run single custom neutron node-1", func() {
			updated_path2 := UpdateNeutronPublicPorts(NEUTRON_CONFIG0)
			enclaveName := GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "neutron", "-c", updated_path2, "--enclaveName", enclaveName)
			defer os.Remove(updated_path2)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should run single custom neutron node with invalid json path", func() {
			enclaveName := GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "neutron", "-c", "../../cli/sample-jsons/neutron5.json", "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("should output user that chain is already running when trying to run neutron chain that is already running", func() {
			enclaveName := GenerateRandomName()
			RunNeutronNode(enclaveName)
			cmd.Args = append(cmd.Args, "chain", "neutron", "--enclaveName", enclaveName)
			defer Clean(enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})
	})

	ginkgo.Describe("Relaychain commands", func() {
		var selectedChain string
		if envChain := os.Getenv("relayChain"); envChain != "" {
			selectedChain = envChain
		} else {
			selectedChain = "default" // Provide a default value if not set
		}
		relayChainNames := []string{"kusama", "polkadot"}

		// Add a flag to check if the selected chain is valid
		validChainSelected := false

		for _, relayChainName := range relayChainNames {

			relayChainName := relayChainName // Capture the loop variable

			if selectedChain != "default" && selectedChain != relayChainName {
				// Skip tests for other chains
				continue
			}

			// Set the flag to indicate that a valid chain is selected
			validChainSelected = true

			ginkgo.It("should run single relaychain "+relayChainName, func() {
				enclaveName := GenerateRandomName()
				cmd.Args = append(cmd.Args, "chain", relayChainName, "--enclaveName", enclaveName)
				defer Clean(enclaveName)
				err := cmd.Run()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})

			ginkgo.It("should run single relaychain in mainnet for "+relayChainName, func() {
				enclaveName := GenerateRandomName()
				cmd.Args = append(cmd.Args, "chain", relayChainName, "-n", "mainnet", "--enclaveName", enclaveName)
				defer Clean(enclaveName)
				err := cmd.Run()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})

			ginkgo.It("should run single relaychain in testnet for "+relayChainName, func() {
				enclaveName := GenerateRandomName()
				cmd.Args = append(cmd.Args, "chain", relayChainName, "-n", "testnet", "--enclaveName", enclaveName)
				defer Clean(enclaveName)
				err := cmd.Run()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})

			ginkgo.It("should run custom relaychain in localnet for "+relayChainName, func() {
				enclaveName := GenerateRandomName()
				config := UpdateRelayChain(LOCAL_CONFIG0, "local", "rococo-local", enclaveName, false, false)
				cmd.Args = append(cmd.Args, "chain", relayChainName, "-c", config, "--enclaveName", enclaveName)
				defer os.Remove(config)
				defer Clean(enclaveName)
				err := cmd.Run()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})

			ginkgo.It("should run custom relaychain in testnet "+relayChainName, func() {
				enclaveName := GenerateRandomName()
				config := UpdateRelayChain(LOCAL_CONFIG0, "testnet", "rococo", enclaveName, false, false)
				cmd.Args = append(cmd.Args, "chain", relayChainName, "-c", config, "--enclaveName", enclaveName)
				defer os.Remove(config)
				defer Clean(enclaveName)
				err := cmd.Run()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})

			ginkgo.It("should run custom relaychain in mainnet"+relayChainName, func() {
				enclaveName := GenerateRandomName()
				config := UpdateRelayChain(LOCAL_CONFIG0, "mainnet", "kusama", enclaveName, false, false)
				cmd.Args = append(cmd.Args, "chain", relayChainName, "-c", config, "--enclaveName", enclaveName)
				defer os.Remove(config)
				defer Clean(enclaveName)
				err := cmd.Run()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})

			ginkgo.It("should run single relaychain with explorer and metrix service for "+relayChainName, func() {
				enclaveName := GenerateRandomName()
				cmd.Args = append(cmd.Args, "chain", relayChainName, "--explorer", "--metrics", "--enclaveName", enclaveName)
				defer Clean(enclaveName)
				err := cmd.Run()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})

			ginkgo.It("should run custom relaychain with explorer services and metrics service in testnet for "+relayChainName, func() {
				enclaveName := GenerateRandomName()
				config := UpdateRelayChain(LOCAL_CONFIG0, "testnet", "rococo", enclaveName, true, true)
				cmd.Args = append(cmd.Args, "chain", relayChainName, "-c", config, "--enclaveName", enclaveName)
				defer os.Remove(config)
				defer Clean(enclaveName)
				err := cmd.Run()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})

		}
		if !validChainSelected && selectedChain != "default" {
			// Print an error message if an invalid chain is selected
			fmt.Printf("Error: Invalid relayChain selected: %s. Expected 'kusama' or 'polkadot'\n", selectedChain)
			log.Fatal("Tests cannot be run because an invalid relayChain is selected.")
		}
	})

	ginkgo.Describe("Parachain commands", func() {
		var selectedRelayChain string
		var selectedParaChain string

		if envChain := os.Getenv("relayChain"); envChain != "" {
			selectedRelayChain = envChain
		} else {
			selectedRelayChain = "default" // Provide a default value if not set
		}

		if envSelectedParaChain := os.Getenv("paraChain"); envSelectedParaChain != "" {
			selectedParaChain = envSelectedParaChain
		} else {
			// If not provided, set a default value
			selectedParaChain = "default"
		}

		relayChainNames := []string{"kusama", "polkadot"}

		if selectedRelayChain != "default" {
			validRelayChain := false
			for _, relayChainName := range relayChainNames {
				if selectedRelayChain == relayChainName {
					validRelayChain = true
					break
				}
			}
			if !validRelayChain {
				fmt.Printf("Error: Invalid relayChain selected: %s. Expected one of %v\n", selectedRelayChain, relayChainNames)
				log.Fatal("Tests cannot be run because an invalid relayChain is selected.")
			}
		}

		for _, relayChainName := range relayChainNames {
			if selectedRelayChain != "default" && selectedRelayChain != relayChainName {
				continue
			}

			relayChainName := relayChainName // Capture the loop variable

			var paraChainNames []string

			if relayChainName == "kusama" {
				paraChainNames = []string{"karura", "kintsugi-btc", "altair", "bifrost", "mangata", "robonomics", "turing-network", "encointer-network", "bajun-networkc", "calamari", "Khala-network", "litmus", "moonriver", "subzero"}
			} else if relayChainName == "polkadot" {
				paraChainNames = []string{"polkadex", "zeitgeist", "acala", "bifrost", "clover", "integritee-shell", "integritee-shell", "litentry", "moonbeam", "nodle", "pendulum", "ajuna-network", "centrifuge", "frequency", "interlay", "kylin", "manta", "moonsama", "parallel", "phala-network", "subsocial"}
			} else {
				paraChainNames = []string{"polkadex", "zeitgeist", "karura", "kintsugi-btc", "altair", "bifrost", "mangata", "robonomics", "turing-network", "encointer-network", "bajun-networkc", "calamari", "khala-network", "litmus", "moonriver", "subzero", "acala", "bifrost", "clover", "integritee-shell", "integritee-shell", "litentry", "moonbeam", "nodle", "pendulum", "ajuna-network", "centrifuge", "frequency", "interlay", "kylin", "manta", "moonsama", "parallel", "phala-network", "subsocial"}
			}

			// Validate paraChain before running tests
			if selectedParaChain != "default" {
				validParaChain := false
				for _, paraChainName := range paraChainNames {
					if selectedParaChain == paraChainName {
						validParaChain = true
						break
					}
				}
				if !validParaChain {
					fmt.Printf("Error: Invalid paraChain selected: %s. Expected a valid parachain name\n", selectedParaChain)
					log.Fatal("Tests cannot be run because an invalid paraChain is selected.")
				}
			}

			for _, paraChainName := range paraChainNames {
				if selectedParaChain != "default" && selectedParaChain != paraChainName {
					continue
				}

				paraChainName := paraChainName

				ginkgo.It("should run single parachain  in testnet for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := GenerateRandomName()
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-p", paraChainName, "--no-relay", "-n", "testnet", "--enclaveName", enclaveName)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					Clean(enclaveName)
				})
				ginkgo.It("should run single parachain in mainnet for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := GenerateRandomName()
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-p", paraChainName, "--no-relay", "-n", "mainnet", "--enclaveName", enclaveName)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					Clean(enclaveName)
				})
				ginkgo.It("should run single parachain in mainnet with explorer services for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := GenerateRandomName()
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-p", paraChainName, "--no-relay", "-n", "mainnet", "--explorer", "--enclaveName", enclaveName)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					Clean(enclaveName)
				})
				ginkgo.It("should run single parachain in mainnet with metrics services for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := GenerateRandomName()
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-p", paraChainName, "--no-relay", "-n", "mainnet", "--metrics", "--enclaveName", enclaveName)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					Clean(enclaveName)
				})
				ginkgo.It("should run custom parachain in testnet with for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := GenerateRandomName()
					config := UpdateParaChain(LOCAL_CONFIG0, "karura", false, false)
					cmd.Args = append(cmd.Args, "chain", relayChainName, "--no-relay", "-c", config, "--enclaveName", enclaveName)
					defer os.Remove(config)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					Clean(enclaveName)
				})
				ginkgo.It("should run custom parachain in mainnet for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := GenerateRandomName()
					config := UpdateParaChain(LOCAL_CONFIG0, "karura", false, false)
					cmd.Args = append(cmd.Args, "chain", relayChainName, "--no-relay", "-c", config, "--enclaveName", enclaveName)
					defer os.Remove(config)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					Clean(enclaveName)
				})
				ginkgo.It("should run custom parachain in mainnet with explorer services for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := GenerateRandomName()
					config := UpdateParaChain(LOCAL_CONFIG0, "karura", true, false)
					cmd.Args = append(cmd.Args, "chain", relayChainName, "--no-relay", "-c", config, "--enclaveName", enclaveName)
					defer os.Remove(config)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					Clean(enclaveName)
				})
				ginkgo.It("should run custom parachain in mainnet with metrics services for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := GenerateRandomName()
					config := UpdateParaChain(LOCAL_CONFIG0, "karura", false, true)
					cmd.Args = append(cmd.Args, "chain", relayChainName, "--no-relay", "-c", config, "--enclaveName", enclaveName)
					defer os.Remove(config)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					Clean(enclaveName)
				})
				ginkgo.It("should run custom parachain in mainnet with explorer and metrics services for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := GenerateRandomName()
					config := UpdateParaChain(LOCAL_CONFIG0, "karura", true, true)
					cmd.Args = append(cmd.Args, "chain", relayChainName, "--no-relay", "-c", config, "--enclaveName", enclaveName)
					defer os.Remove(config)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					Clean(enclaveName)
				})
			}
		}
	})

	ginkgo.Describe("Relaychain and parachain commands", func() {

		var selectedRelayChain string
		var selectedParaChain string

		if envChain := os.Getenv("relayChain"); envChain != "" {
			selectedRelayChain = envChain
		} else {
			selectedRelayChain = "default" // Provide a default value if not set
		}

		if envSelectedParaChain := os.Getenv("paraChain"); envSelectedParaChain != "" {
			selectedParaChain = envSelectedParaChain
		} else {
			// If not provided, set a default value
			selectedParaChain = "default"
		}

		relayChainNames := []string{"kusama", "polkadot"}

		// Validate relayChain before running tests
		if selectedRelayChain != "default" {
			validRelayChain := false
			for _, relayChainName := range relayChainNames {
				if selectedRelayChain == relayChainName {
					validRelayChain = true
					break
				}
			}
			if !validRelayChain {
				fmt.Printf("Error: Invalid relayChain selected: %s. Expected one of %v\n", selectedRelayChain, relayChainNames)
				log.Fatal("Tests cannot be run because an invalid relayChain is selected.")
			}
		}

		// if selectedParaChain == "default" && selectedRelayChain == "default" {
		// 	fmt.Println("Error: Atleast relay chain should be given. ")
		// 	log.Fatal("Tests cannot be run because relayChain is missing.")
		// 	return // Added return to stop further execution
		// }

		for _, relayChainName := range relayChainNames {
			if selectedRelayChain != "default" && selectedRelayChain != relayChainName {
				continue

			}

			relayChainName := relayChainName // Capture the loop variable

			var paraChainNames []string

			if relayChainName == "kusama" {
				paraChainNames = []string{"karura", "kintsugi-btc", "altair", "bifrost", "mangata", "robonomics", "turing-network", "encointer-Network", "bajun-networkc", "calamari", "khala-network", "litmus", "moonriver", "subzero"}
			} else if relayChainName == "polkadot" {
				paraChainNames = []string{"Polkadex", "zeitgeist", "acala", "bifrost", "clover", "integritee-shell", "integritee-shell", "litentry", "moonbeam", "nodle", "pendulum", "ajuna-network", "centrifuge", "frequency", "interlay", "kylin", "manta", "moonsama", "parallel", "phala-network", "subsocial"}
			} else {
				paraChainNames = []string{"polkadex", "zeitgeist", "karura", "kintsugi-btc", "altair", "bifrost", "mangata", "robonomics", "turing-network", "encointer-network", "bajun-networkc", "calamari", "khala-network", "litmus", "moonriver", "subzero", "acala", "bifrost", "clover", "integritee-shell", "integritee-shell", "litentry", "moonbeam", "nodle", "pendulum", "ajuna-network", "centrifuge", "frequency", "interlay", "kylin", "manta", "moonsama", "parallel", "phala-network", "subsocial"}
			}

			// Validate paraChain before running tests
			if selectedParaChain != "default" {
				validParaChain := false
				for _, paraChainName := range paraChainNames {
					if selectedParaChain == paraChainName {
						validParaChain = true
						break
					}
				}
				if !validParaChain {
					fmt.Printf("Error: Invalid paraChain selected: %s. Expected a valid parachain name\n", selectedParaChain)
					log.Fatal("Tests cannot be run because an invalid paraChain is selected.")
				}
			}

			for _, paraChainName := range paraChainNames {
				if selectedParaChain != "default" && selectedParaChain != paraChainName {
					continue
				}

				paraChainName := paraChainName
				ginkgo.It("should run single relaychain and parachain  in testnet for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := GenerateRandomName()
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-p", paraChainName, "-n", "testnet", "--enclaveName", enclaveName)
					defer Clean(enclaveName)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
				})
				ginkgo.It("should run single relaychain and parachain in mainnet for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := GenerateRandomName()
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-p", paraChainName, "-n", "mainnet", "--enclaveName", enclaveName)
					defer Clean(enclaveName)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
				})
				ginkgo.It("should run single relaychain and parachain in local for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := GenerateRandomName()
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-p", paraChainName, "-n", "local", "--enclaveName", enclaveName)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					Clean(enclaveName)
				})
				ginkgo.It("should run single relaychain and parachain in mainnet with explorer services for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := GenerateRandomName()
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-p", paraChainName, "-n", "mainnet", "--explorer", "--enclaveName", enclaveName)
					defer Clean(enclaveName)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
				})
				ginkgo.It("should run single relaychain and  parachain in mainnet with metrics services for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := GenerateRandomName()
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-p", paraChainName, "-n", "mainnet", "--metrics", "--enclaveName", enclaveName)
					defer Clean(enclaveName)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
				})
				ginkgo.It("should run custom relaychain and parachain in testnet for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := GenerateRandomName()
					config := UpdateChainInfo(LOCAL_CONFIG0, "testnet", "rococo", "karura", false, false)
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-c", config, "--enclaveName", enclaveName)
					defer os.Remove(config)
					defer Clean(enclaveName)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
				})
				ginkgo.It("should run custom relaychain and  parachain in mainnet for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := GenerateRandomName()
					config := UpdateChainInfo(LOCAL_CONFIG0, "mainnet", "polkadot", "karura", false, false)
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-c", config, "--enclaveName", enclaveName)
					defer os.Remove(config)
					defer Clean(enclaveName)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
				})
				ginkgo.It("should run custom relaychain and  parachain in local for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := GenerateRandomName()
					config := UpdateChainInfo(LOCAL_CONFIG0, "local", "rococo-local", "karura", false, false)
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-c", config, "--enclaveName", enclaveName)
					defer os.Remove(config)
					defer Clean(enclaveName)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
				})
				ginkgo.It("should run custom relaychain and  parachain in mainnet  with explorer services for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := GenerateRandomName()
					config := UpdateChainInfo(LOCAL_CONFIG0, "mainnet", "polkadot", "frequency", true, false)
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-c", config, "--enclaveName", enclaveName)
					defer os.Remove(config)
					defer Clean(enclaveName)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
				})
				ginkgo.It("should run custom relaychain and parachain in mainnet  with metrics services for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := GenerateRandomName()
					config := UpdateChainInfo(LOCAL_CONFIG0, "mainnet", "polkadot", "frequency", true, false)
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-c", config, "--enclaveName", enclaveName)
					defer os.Remove(config)
					defer Clean(enclaveName)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
				})
				ginkgo.It("should run custom relaychain and parachain in mainnet with explorer and metrics services for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := GenerateRandomName()
					config := UpdateChainInfo(LOCAL_CONFIG0, "mainnet", "polkadot", "frequency", true, true)
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-c", config, "--enclaveName", enclaveName)
					defer os.Remove(config)
					defer Clean(enclaveName)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
				})
			}
		}
	})
})
