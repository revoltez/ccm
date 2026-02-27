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

- [ ] Implement LoadBalancer interface (for Service type LoadBalancer)
- [ ] Implement Instances interface (for node management)
- [ ] Implement InstancesV2 interface (optimized instance API)
- [ ] Implement Zones interface (for zone/region info)
- [ ] Implement Routes interface (for cluster networking)
- [ ] Implement Clusters interface

## Phase 3: Controller Implementation

- [ ] Implement Node Controller logic
- [ ] Implement Service Controller logic (LoadBalancer)
- [ ] Implement Route Controller logic

## Phase 4: Production Readiness

- [ ] Add proper logging
- [ ] Add metrics
- [ ] Implement graceful shutdown
- [ ] Add health checks
