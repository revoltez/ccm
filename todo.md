# CCM Development TODO

## Phase 1: Initial Setup

- [x] Initialize Go module
- [x] Create basic cloud provider interface implementation
- [x] Create main entry point using cloud-provider scaffolding
- [x] Create test scripts
- [x] Create Dockerfile
- [x] Create Taskfile.yml & bootstraping k3s with external cloud provider configs script
- [x] Create Helm chart for CCM deployment
- [x] Configure k3s to disable in-tree cloud controller

## Phase 2: Cloud Provider Interface Implementation

- [x] Implement InstancesV2 interface (node registration, metadata, taint removal)
- [x] Document the methods for future implementations
- [x] Add environment variable support for provider configuration

- [ ] Implement LoadBalancer interface (for Service type LoadBalancer)
