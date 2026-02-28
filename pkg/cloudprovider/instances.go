package cloudprovider

import (
	"context"
	"fmt"
	"os"
	"strings"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	cloudprovider "k8s.io/cloud-provider"
)

const (
	ProviderIDPrefix = "salih-ccm://"
)

var _ cloudprovider.InstancesV2 = (*Instances)(nil)

type Instances struct{}

func NewInstances() *Instances {
	return &Instances{}
}

func (i *Instances) InstanceExists(ctx context.Context, node *v1.Node) (bool, error) {
	fmt.Printf("InstanceExists called for node: %s\n", node.Name)
	return true, nil
}

func (i *Instances) InstanceShutdown(ctx context.Context, node *v1.Node) (bool, error) {
	fmt.Printf("InstanceShutdown called for node: %s\n", node.Name)
	return false, nil
}

func (i *Instances) InstanceMetadata(ctx context.Context, node *v1.Node) (*cloudprovider.InstanceMetadata, error) {
	fmt.Printf("InstanceMetadata called for node: %s\n", node.Name)

	providerID := getProviderID(node)
	addresses := collectNodeAddresses(node)

	metadata := &cloudprovider.InstanceMetadata{
		ProviderID:    providerID,
		InstanceType:  "k3s-node",
		NodeAddresses: addresses,
		Zone:          getZone(),
		Region:        getRegion(),
		AdditionalLabels: map[string]string{
			"ccm": ProviderName,
		},
	}

	fmt.Printf("InstanceMetadata returning for node: %s, providerID: %s, zone: %s, region: %s\n",
		node.Name, providerID, metadata.Zone, metadata.Region)

	return metadata, nil
}

func getProviderID(node *v1.Node) string {
	if node.Spec.ProviderID != "" {
		return node.Spec.ProviderID
	}
	return ProviderIDPrefix + node.Name
}

func collectNodeAddresses(node *v1.Node) []v1.NodeAddress {
	addresses := node.Status.Addresses

	if len(addresses) == 0 {
		hostname, _ := os.Hostname()
		addresses = []v1.NodeAddress{
			{Type: v1.NodeHostName, Address: hostname},
			{Type: v1.NodeInternalIP, Address: node.Name},
		}
	}
	fmt.Println("addresses:", addresses)

	return addresses
}

func getZone() string {
	if z := os.Getenv("CCM_ZONE"); z != "" {
		return z
	}
	return "local-zone"
}

func getRegion() string {
	if r := os.Getenv("CCM_REGION"); r != "" {
		return r
	}
	return "local-region"
}

func ProviderIDToNodeName(providerID string) (types.NodeName, error) {
	if !strings.HasPrefix(providerID, ProviderIDPrefix) {
		return "", fmt.Errorf("invalid provider ID format: %s", providerID)
	}
	nodeName := strings.TrimPrefix(providerID, ProviderIDPrefix)
	return types.NodeName(nodeName), nil
}

func NodeNameToProviderID(nodeName types.NodeName) string {
	return ProviderIDPrefix + string(nodeName)
}
