# Distributed Load Tester & API Benchmarker

A containerized load testing tool written in Go, orchestrated with Kubernetes. It is designed to measure API latency and aggregate results across a distributed network using gRPC.

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
![Kubernetes](https://img.shields.io/badge/kubernetes-%23326ce5.svg?style=for-the-badge&logo=kubernetes&logoColor=white)
![gRPC](https://img.shields.io/badge/gRPC-%23244c5a.svg?style=for-the-badge&logo=grpc&logoColor=white)

## Current Capabilities

Currently, the application acts as the foundational infrastructure for a distributed testing system. 

* **Containerized Load Generation:** Executes batches of HTTP requests against a target API from an isolated Alpine Linux environment.
* **gRPC Telemetry:** Transmits test results (total requests, success rates, latency) across the Kubernetes network using a custom Protobuf contract.
* **Centralized Aggregation:** A dedicated server listens for incoming payloads from worker nodes and calculates the average latency of the benchmark run.

## Architecture

The system is currently composed of two main microservices:

1. **The Worker Node Service (The Muscle):** A lightweight Go container deployed as a Kubernetes `Job`. It wakes up, executes a batch of HTTP requests against a target URL, and sends the performance payload to the Aggregator over gRPC before shutting down.
2. **The Aggregator Service:** A Go server deployed as a Kubernetes `Deployment`. It runs continuously, listening on port 50051 for incoming gRPC connections from Worker pods. It receives the batched data streams, calculates the current average latency, and logs the metrics to standard output.

## Tech Stack

* **Language:** Go (Golang)
* **Communication:** gRPC / Protocol Buffers
* **Containerization:** Docker (Multi-stage builds)
* **Orchestration:** Kubernetes (Deployments, Services, Jobs)

## Quick Start (Local Kubernetes Cluster)

Prerequisites:
* Docker Desktop (with Kubernetes enabled)
* `kubectl` CLI configured

### 1. Start the Aggregator (The Brain)
Deploy the central server that will listen for incoming metrics:
```bash
kubectl apply -f aggregator.yaml
```
Verify the Aggregator is running and its network service is active:
```bash
kubectl logs deployment/aggregator-deployment
```

### 2. Build the Worker (The Muscle)
Compile the Go load-testing engine into an Alpine Linux Docker container:
```bash
docker build -t worker-service -f services/worker/Dockerfile.worker .
```

### 3. Launch the Load Test
Deploy the Kubernetes Job to spin up the Worker pod, execute the requests, and transmit the payload:
```bash
# Ensure any previous completed jobs are cleared
kubectl delete job worker-job --ignore-not-found=true

# Deploy the new benchmark job
kubectl apply -f worker.yaml
```

### 4. View the Results
Check the Worker logs to verify the HTTP requests were successfully executed:
```bash
kubectl logs job/worker-job
```
Check the Aggregator's logs to see the incoming gRPC payloads and latency averages:
```bash
kubectl logs deployment/aggregator-deployment
```