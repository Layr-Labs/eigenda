package integration_test

import (
	"os"
	"os/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Operator Clock Drift Reporter Integration", func() {
	It("should successfully report clock drift to the DataAPI endpoint", func() {
		// Set up environment variables for the script
		os.Setenv("OPERATOR_PRIVATE_KEY", "0x0000000000000000000000000000000000000000000000000000000000000000")
		os.Setenv("EIGENDA_API_ENDPOINT", "http://localhost:8080/api/v1/clock_drift")
		os.Setenv("NTP_SERVER", "pool.ntp.org")

		cmd := exec.Command("python3", "tools/operator_scripts/clock_drift_reporter.py")
		cmd.Env = os.Environ()
		output, err := cmd.CombinedOutput()
		GinkgoWriter.Println(string(output))
		Expect(err).To(BeNil(), "Script should run without error")
		Expect(string(output)).To(ContainSubstring("Offset reported"))
		Expect(string(output)).To(ContainSubstring("HTTP 200"))
	})
})
