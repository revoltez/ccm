# CCM - Cloud Controller Manager

A Kubernetes Cloud Controller Manager implementation for local/bare-metal clusters (k3s).

This project demonstrates how to implement a Kubernetes Cloud Controller Manager by implementing the cloud provider interface. It provides a minimal working CCM that enables nodes to become "Ready" in a bare-metal k3s cluster.

## Prerequisites

- Go 1.21+
- Docker
- kubectl
- helm
- [go-task](https://taskfile.dev/)
- k3s (installed via this project's tasks)

## Quick Start

### Step 1: Install k3s with External Cloud Provider

```bash
task k3s-install
```

This installs k3s with the following configuration:

- `--disable-cloud-controller` - Disables the in-tree cloud controller
- `--kubelet-arg=cloud-provider=external` - Tells kubelet to use an external CCM
- `--write-kubeconfig-mode 644` - Allows access to the kubeconfig file

**What happens:** k3s starts as a single-node cluster, but the node will have a taint `node.cloudprovider.kubernetes.io/uninitialized` that prevents workloads from being scheduled. This is because Kubernetes expects the cloud provider (CCM) to provide node information.

### Step 2: Build and Deploy CCM

```bash
task dev
```

This runs a sequence of tasks:

1. `task build` - Compiles the Go binary
2. `task docker-build` - Creates a Docker image
3. `task k3s-load` - Loads the Docker image into k3s
4. `task apply` - Deploys the CCM using Helm

### Step 3: Verify Node is Ready

```bash
kubectl get nodes
```

You should see the node transition from `NotReady` (due to the taint) to `Ready`.

The CCM achieves this by implementing the `InstancesV2` interface. When a node joins the cluster, the CCM watches for node events through the Kubernetes API server and provides node information (providerID, addresses, zone, region) via its `InstanceMetadata()` implementation.

## Configuration

Environment variables can be configured in `helm/ccm/values.yaml`:

```yaml
config:
  providerName: "salih-ccm"
  providerIDPrefix: "salih-ccm://"
  zone: "local-zone"
  region: "local-region"
```

Or via helm upgrade:

```bash
helm upgrade --install ccm ./helm/ccm --set config.zone=us-east-1a,config.region=us-east-1
```

## How It Works

### The InstancesV2 Interface

This CCM implements `InstancesV2`, which is the modern interface for node management. It has three methods:

1. **InstanceExists** - Checks if the instance exists
2. **InstanceShutdown** - Checks if the instance is shutdown
3. **InstanceMetadata** - Returns metadata about the instance

When `InstanceMetadata` is called, it returns an `InstanceMetadata` struct containing:

- **ProviderID** - Unique identifier for the node (set on `node.spec.providerID`)
- **NodeAddresses** - IP addresses (set on `node.status.addresses`)
- **InstanceType** - Type of instance (set as node label)
- **Zone** - Availability zone (set as node label)
- **Region** - Region (set as node label)
- **AdditionalLabels** - Any additional labels

### Why CCM Can Be Scheduled

The CCM DaemonSet includes a toleration for the `node.cloudprovider.kubernetes.io/uninitialized` taint:

```yaml
tolerations:
  - key: "node.cloudprovider.kubernetes.io/uninitialized"
    operator: "Equal"
    value: "true"
    effect: "NoSchedule"
```

This is critical - without this toleration, the CCM pods would never be scheduled on nodes that have this taint. It's a chicken-and-egg problem: the node needs the CCM to become ready, but the CCM needs to run on the node to provide information. The toleration breaks this cycle by allowing the CCM to run despite the taint.

### Implementing for Real Cloud Providers

To extend this CCM for a real cloud provider (AWS, GCP, Azure, etc.), you would:

1. **Update the configuration** in `helm/ccm/values.yaml`:
   - Set `config.providerName` to your cloud provider name (e.g., "aws", "gce", "azure")
   - Set `config.providerIDPrefix` to the appropriate prefix (e.g., "aws://", "gce://", "azure://")
   - Update zone and region as needed

2. **Implement the `InstancesV2` methods** in `instances.go` to call the cloud provider's API:
   - Use AWS SDK for EC2, GCP SDK for Compute Engine, Azure SDK for Azure VM, etc.
   - Fetch real instance metadata from the cloud API
   - Handle authentication (IAM roles, service principals)

3. **Implement other interfaces** as needed:
   - `LoadBalancer` - For Service type=LoadBalancer
   - `Routes` - For cluster networking
   - `Clusters` - For multi-cluster support
