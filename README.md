# Distributed Load Tester & API Benchmarker

A high-performance, containerized load-testing suite written in Go and orchestrated with Kubernetes. This system benchmarks API reliability by distributing concurrent load generation across worker nodes and aggregating real-time telemetry via gRPC into a stateful persistence layer.

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
![Kubernetes](https://img.shields.io/badge/kubernetes-%23326ce5.svg?style=for-the-badge&logo=kubernetes&logoColor=white)
![gRPC](https://img.shields.io/badge/gRPC-%23244c5a.svg?style=for-the-badge&logo=grpc&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)

## Core Capabilities

*   **Massive Concurrency:** Leverages Go’s **Goroutines and channels** to execute thousands of simultaneous HTTP requests while maintaining strict memory safety.
*   **High-Fidelity Observability:** Calculates exact **P95 and P99 latency percentiles** using Mutex-protected data structures to provide deep insights into system tail latency.
*   **Low-Latency Telemetry:** Utilizes a custom **gRPC networking contract** for type-safe, high-speed data transmission between workers and the central aggregator.
*   **Stateful Storage:** Persists every benchmark run and aggregated metric into a **PostgreSQL** database for historical analysis and trend reporting.

## Architecture

The system follows a distributed producer-consumer model orchestrated by Kubernetes:

1.  **Worker Nodes (Go):** Deployed as Kubernetes `Jobs`. These containers orchestrate load generation, calculate local performance metrics, and stream payloads to the aggregator.
2.  **Central Aggregator (Go):** A Kubernetes `Deployment` acting as the gRPC server. It consumes telemetry streams, performs final statistical aggregation, and manages the database transaction layer.
3.  **Persistence Layer (PostgreSQL):** A stateful storage service ensuring all benchmark data survives container restarts and cluster updates.

## Tech Stack

*   **Language:** Go (Golang)
*   **Protocols:** gRPC, Protocol Buffers (Protobuf)
*   **Containerization:** Docker (Multi-stage Alpine builds)
*   **Orchestration:** Kubernetes (Deployments, StatefulSets, Jobs)
*   **Database:** PostgreSQL

## Quick Start (Local Development)

Execute the following commands to initialize the environment, build the engine, and launch a benchmark:
```bash
# 1. Initialize the Environment (Database & Aggregator)
kubectl apply -f postgres.yaml
kubectl apply -f aggregator.yaml

# 2. Build the Worker Engine
docker build -t worker-service -f services/worker/Dockerfile.worker .

# 3. Launch a Benchmark Run (Clean previous jobs first)
kubectl delete job worker-job --ignore-not-found=true
kubectl apply -f worker.yaml

# 4. Analyze Results (View P95/P99 metrics and DB insertions)
kubectl logs deployment/aggregator-deployment