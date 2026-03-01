// Package cloudprovider implements the Kubernetes cloud provider interface
// for bare-metal/local clusters (k3s).
package cloudprovider

import (
	"context"
	"fmt"
	"os"

	v1 "k8s.io/api/core/v1"
	cloudprovider "k8s.io/cloud-provider"
)

// InstancesV2 implements the cloud provider instances interface for local/bare-metal clusters.
// This is a simplified implementation that works with k3s without requiring actual cloud API calls.
//
// For a real cloud provider (AWS, GCP, Azure, etc.), you would need to:
// - Call the cloud provider's API (EC2, Compute Engine, Azure VM, etc.) to verify instance existence
// - Fetch instance metadata (IPs, instance type, zone, region) from the cloud API
// - Handle authentication (IAM roles, service principals, etc.)
//
// The InstancesV2 interface is preferred over the legacy Instances interface as it:
// - Reduces API calls to the cloud provider
// - Provides more efficient metadata retrieval
// - Disables the deprecated Zones interface when implemented
var _ cloudprovider.InstancesV2 = (*Instances)(nil)

// Instances provides a simplified cloud provider implementation for local clusters.
// In production, this would connect to actual cloud provider APIs.
type Instances struct{}

// NewInstances creates a new Instances implementation for local/bare-metal clusters.
func NewInstances() *Instances {
	return &Instances{}
}

// InstanceExists checks if the instance for the given node exists in the cloud provider.
//
// For real cloud providers, this would typically:
// - Call the cloud API (e.g., EC2 DescribeInstances) using the node's providerID
// - Return an error if the instance no longer exists
// - Handle rate limiting and retries
//
// Parameters:
//   - ctx: Context for the API call
//   - node: The Kubernetes node object (use node.Name or node.Spec.ProviderID to identify)
//
// Returns:
//   - bool: Whether the instance exists
//   - error: Any error that occurred (return cloudprovider.InstanceNotFound if instance is gone)
func (i *Instances) InstanceExists(ctx context.Context, node *v1.Node) (bool, error) {
	fmt.Printf("[InstancesV2] InstanceExists: node=%s, providerID=%s\n", node.Name, node.Spec.ProviderID)
	// For local/bare-metal, assume all nodes exist
	return true, nil
}

// InstanceShutdown checks if the instance is shutdown in the cloud provider.
//
// For real cloud providers, this would:
// - Call the cloud API to check instance state (e.g., EC2 DescribeInstances with state)
// - Return true if the instance is stopped/terminated
// - This affects Kubernetes node lifecycle events
//
// Parameters:
//   - ctx: Context for the API call
//   - node: The Kubernetes node object
//
// Returns:
//   - bool: Whether the instance is shutdown
//   - error: Any error that occurred
func (i *Instances) InstanceShutdown(ctx context.Context, node *v1.Node) (bool, error) {
	fmt.Printf("[InstancesV2] InstanceShutdown: node=%s\n", node.Name)
	// For local/bare-metal, assume nodes are running
	return false, nil
}

// InstanceMetadata returns the instance's metadata including addresses, zone, region, and instance type.
// This metadata is applied to the Kubernetes Node object during registration.
//
// For real cloud providers, this would:
// - Call the cloud API to fetch instance metadata
// - Return details like: private/public IPs, instance type, zone, region, tags
// - Handle authentication and API rate limiting
//
// The returned metadata sets the following on the Node:
//   - spec.providerID: The unique instance ID
//   - status.addresses: Node IP addresses (InternalIP, ExternalIP)
//   - labels: node.kubernetes.io/instance-type, topology.kubernetes.io/zone, topology.kubernetes.io/region
//
// Parameters:
//   - ctx: Context for the API call
//   - node: The Kubernetes node object (check node.Spec.ProviderID first, then node.Name)
//
// Returns:
//   - *cloudprovider.InstanceMetadata: The instance metadata
//   - error: Any error that occurred
func (i *Instances) InstanceMetadata(ctx context.Context, node *v1.Node) (*cloudprovider.InstanceMetadata, error) {
	fmt.Printf("[InstancesV2] InstanceMetadata: node=%s\n", node.Name)

	providerID := getProviderID(node)
	addresses := collectNodeAddresses(node)

	metadata := &cloudprovider.InstanceMetadata{
		ProviderID:    providerID,
		InstanceType:  "k3s-node",
		NodeAddresses: addresses,
		Zone:          getZone(),
		Region:        getRegion(),
		AdditionalLabels: map[string]string{
			"ccm": getProviderName(),
		},
	}

	fmt.Printf("[InstancesV2] InstanceMetadata: providerID=%s, zone=%s, region=%s, addresses=%v\n",
		providerID, metadata.Zone, metadata.Region, addresses)

	return metadata, nil
}

// getProviderID returns the cloud provider ID for the node.
// Uses the node's existing providerID if set, otherwise generates one from the provider prefix.
func getProviderID(node *v1.Node) string {
	if node.Spec.ProviderID != "" {
		return node.Spec.ProviderID
	}
	return getProviderIDPrefix() + node.Name
}

// collectNodeAddresses collects node addresses from the node's status.
// Falls back to hostname-based addresses if none are present.
func collectNodeAddresses(node *v1.Node) []v1.NodeAddress {
	addresses := node.Status.Addresses

	if len(addresses) == 0 {
		hostname, _ := os.Hostname()
		addresses = []v1.NodeAddress{
			{Type: v1.NodeHostName, Address: hostname},
			{Type: v1.NodeInternalIP, Address: node.Name},
		}
	}

	return addresses
}

// getZone returns the zone from environment variable or default.
func getZone() string {
	if z := os.Getenv("CCM_ZONE"); z != "" {
		return z
	}
	return "local-zone"
}

// getRegion returns the region from environment variable or default.
func getRegion() string {
	if r := os.Getenv("CCM_REGION"); r != "" {
		return r
	}
	return "local-region"
}

// getProviderName returns the cloud provider name from environment variable or default.
func getProviderName() string {
	if name := os.Getenv("CCM_PROVIDER_NAME"); name != "" {
		return name
	}
	return ProviderName
}

// getProviderIDPrefix returns the provider ID prefix from environment variable or default.
func getProviderIDPrefix() string {
	if prefix := os.Getenv("CCM_PROVIDER_ID_PREFIX"); prefix != "" {
		return prefix
	}
	return ProviderName + "://"
}
