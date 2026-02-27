package main

import (
	"os"

	ccmcloud "github.com/revoltez/ccm/pkg/cloudprovider"
	"github.com/spf13/cobra"
	cloudprovider "k8s.io/cloud-provider"
	app "k8s.io/cloud-provider/app"
	cloudcontrollerconfig "k8s.io/cloud-provider/app/config"
	"k8s.io/cloud-provider/options"
	"k8s.io/klog/v2"
)

func main() {
	opts, err := options.NewCloudControllerManagerOptions()
	if err != nil {
		klog.Fatalf("Failed to create options: %v", err)
	}

	cmd := buildCommand(opts)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func buildCommand(opts *options.CloudControllerManagerOptions) *cobra.Command {
	stopCh := make(chan struct{})

	cmdBuilder := app.NewBuilder()
	cmdBuilder.SetOptions(opts)
	cmdBuilder.RegisterDefaultControllers()
	cmdBuilder.SetCloudInitializer(initCloudProvider)
	cmdBuilder.SetStopChannel(stopCh)

	return cmdBuilder.BuildCommand()
}

func initCloudProvider(cfg *cloudcontrollerconfig.CompletedConfig) cloudprovider.Interface {
	cloud, err := cloudprovider.GetCloudProvider(ccmcloud.ProviderName, nil)
	if err != nil {
		klog.Fatalf("Failed to initialize cloud provider: %v", err)
	}
	return cloud
}
