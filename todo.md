# CCM Development TODO

## Phase 1: Initial Setup

- [x] Initialize Go module
- [x] Create basic cloud provider interface implementation
- [x] Create main entry point using cloud-provider scaffolding
- [x] Create test script
- [x] Create Dockerfile (multi-stage build inside container)
- [x] Create Taskfile.yml with go-task
- [x] Create Helm chart for CCM deployment
- [x] Configure k3s to disable in-tree cloud controller

## Phase 2: Cloud Provider Interface Implementation

### Bare-metal

- [x] Implement InstancesV2 interface (node registration, metadata, taint removal)
- [x] Document the methods for future implementations
- [x] Add environment variable support for provider configuration

- [ ] Implement LoadBalancer interface (for Service type LoadBalancer)
- [ ] Implement Routes interface (for cluster networking)
- [ ] Implement Clusters interface
