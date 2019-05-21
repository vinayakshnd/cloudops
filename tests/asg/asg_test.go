package tests

import (
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
	"github.com/portworx/torpedo/drivers/node"
	. "github.com/portworx/torpedo/tests"
)

const (
	ASG_VOLUME_TAG_KEY = "pwx_cluster_id"
)

func TestASG(t *testing.T) {
	RegisterFailHandler(Fail)

	var specReporters []Reporter
	junitReporter := reporters.NewJUnitReporter("/testresults/junit_Reboot.xml")
	specReporters = append(specReporters, junitReporter)
	RunSpecsWithDefaultAndCustomReporters(t, "Torpedo : ASG", specReporters)
}

var _ = BeforeSuite(func() {
	InitInstance()
})

var _ = Describe("{ValidateTags}", func() {

	It("has to validate tags for each cloud drive", func() {
		nodeMap := node.GetNodesByName()

		for _, node := range nodeMap {
			disks, err := Inst().N.AttachedDisks(node)
			Expect(err).NotTo(HaveOccurred())
			Expect(disks).NotTo(BeEmpty())

			for _, disk := range disks {
				Step(fmt.Sprintf("Validate [%s] disk have tag [%s]", disk, ASG_VOLUME_TAG_KEY), func() {
					tags, err := Inst().N.Tags(disk)
					Expect(err).NotTo(HaveOccurred())
					Expect(tags).To(HaveKey(ASG_VOLUME_TAG_KEY))
				})
			}
		}
	})
})

var _ = AfterSuite(func() {
	PerformSystemCheck()
	CollectSupport()
	ValidateCleanup()
})

func init() {
	ParseFlags()
}
