package gke

import (
	"fmt"

	"github.com/libopenstorage/cloudops"
	"github.com/libopenstorage/cloudops/gce"
	"github.com/portworx/torpedo/drivers/node"
	compute "google.golang.org/api/compute/v1"

	// https://github.com/kubernetes/client-go/issues/242
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

const (
	// DriverName is the name of the gke driver
	DriverName = "gke"
)

type gke struct {
	node.Driver
	ops cloudops.Ops
}

func (g *gke) String() string {
	return DriverName
}

func (g *gke) Init() error {
	ops, err := gce.NewClient()
	if err != nil {
		return err
	}
	g.ops = ops
	return nil
}

func (g *gke) AttachedDisks(node node.Node) ([]string, error) {

	attachedDisks := []string{}
	inst, err := g.ops.Describe()
	if err != nil {
		return nil, err
	}

	computeInst := inst.(*compute.Instance)

	for _, disk := range computeInst.Disks {
		attachedDisks = append(attachedDisks, disk.DeviceName)
	}

	fmt.Printf("List of disks attached to node %s is %v", node.Name, attachedDisks)
	return attachedDisks, nil
}

func (g *gke) Tags(volumeID string) (map[string]string, error) {

	volumeTags, err := g.ops.Tags(volumeID)
	if err != nil {
		return nil, err
	}
	return volumeTags, nil
}

func init() {
	g := &gke{
		Driver: node.NotSupportedDriver,
	}

	node.Register(DriverName, g)
}
