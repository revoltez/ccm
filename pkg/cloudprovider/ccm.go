package cloudprovider

import (
	"encoding/json"
	"fmt"
	"io"

	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/klog/v2"
)

var _ cloudprovider.Interface = (*CCM)(nil)

func init() {
	cloudprovider.RegisterCloudProvider(ProviderName, NewCloud)
}

type CCM struct {
	cfg       *CCMConfig
	instances *Instances
}

type CCMConfig struct {
	ClientConfig string `json:"clientConfig"`
}

func NewCloud(config io.Reader) (cloudprovider.Interface, error) {
	ccm := &CCM{
		instances: NewInstances(),
	}

	if config == nil {
		klog.InfoS("No config provided, initializing CCM with defaults")
		return ccm, nil
	}

	cfg, err := readConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	ccm.cfg = cfg

	klog.InfoS("Initializing CCM cloud provider", "config", cfg.ClientConfig)

	return ccm, nil
}

func (c *CCM) Initialize(_ cloudprovider.ControllerClientBuilder, _ <-chan struct{}) {
	klog.Info("CCM cloud provider initialized")
}

func (c *CCM) LoadBalancer() (cloudprovider.LoadBalancer, bool) {
	return nil, false
}

func (c *CCM) Instances() (cloudprovider.Instances, bool) {
	return nil, false
}

func (c *CCM) InstancesV2() (cloudprovider.InstancesV2, bool) {
	return c.instances, true
}

func (c *CCM) Zones() (cloudprovider.Zones, bool) {
	return nil, false
}

func (c *CCM) Clusters() (cloudprovider.Clusters, bool) {
	return nil, false
}

func (c *CCM) Routes() (cloudprovider.Routes, bool) {
	return nil, false
}

func (c *CCM) ProviderName() string {
	return ProviderName
}

func (c *CCM) HasClusterID() bool {
	return true
}

func readConfig(configReader io.Reader) (*CCMConfig, error) {
	if configReader == nil {
		return nil, fmt.Errorf("no config provided")
	}

	data, err := io.ReadAll(configReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg CCMConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}
