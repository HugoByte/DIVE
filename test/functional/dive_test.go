package dive_test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/hugobyte/dive/common"
	"github.com/onsi/ginkgo"
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
			Clean()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "eth")
			err := cmd.Run()
			Clean()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge icon and hardhat but with icon bridge set to true", func() {
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "hardhat", "--bmvbridge")
			err := cmd.Run()
			Clean()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge icon and icon", func() {
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "icon")
			err := cmd.Run()
			Clean()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

	})
})
