# Distributed Load Tester & API Benchmarker

A high-performance, containerized load-testing suite written in Go and orchestrated with Kubernetes. This system benchmarks API reliability by distributing concurrent load generation across worker nodes and aggregating real-time telemetry via gRPC into a stateful persistence layer.

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
![Kubernetes](https://img.shields.io/badge/kubernetes-%23326ce5.svg?style=for-the-badge&logo=kubernetes&logoColor=white)
![gRPC](https://img.shields.io/badge/gRPC-%23244c5a.svg?style=for-the-badge&logo=grpc&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)

## Core Capabilities

* **Massive Concurrency:** Leverages Go’s **Goroutines and channels** to execute thousands of simultaneous HTTP requests with a non-blocking fan-in pattern.
* **High-Fidelity Observability:** Calculates exact **P95 and P99 latency percentiles** to provide deep insights into system tail latency.
* **Low-Latency Telemetry:** Utilizes a custom **gRPC networking contract** for type-safe, high-speed batch data transmission between workers and the aggregator.
* **Stateful Storage:** Persists every benchmark run and aggregated metric into **PostgreSQL** for historical reporting.

## Architecture

The system follows a distributed producer-consumer model orchestrated by Kubernetes:

1.  **Worker Nodes (Go):** Deployed as Kubernetes `Jobs`. They read dynamic environment variables, execute concurrent HTTP load, and stream batched results.
2.  **Central Aggregator (Go):** A Kubernetes `Deployment` (gRPC server) that consumes telemetry, calculates percentiles, and manages database transactions.
3.  **Persistence Layer (PostgreSQL):** A stateful service ensuring all benchmark data survives container lifecycle events.

## Tech Stack

* **Language:** Go (Golang)
* **Protocols:** gRPC, Protocol Buffers (Protobuf)
* **Containerization:** Docker (Multi-stage Alpine builds)
* **Orchestration:** Kubernetes (Deployments, Services, Jobs)
* **Automation:** PowerShell (Configuration Templating)

## Quick Start (Local Development)

Execute the following block to initialize infrastructure, build the worker engine, and run a dynamic benchmark:

```bash
# 1. Initialize Infrastructure (Database & Aggregator)
kubectl apply -f postgres.yaml
kubectl apply -f aggregator.yaml

# 2. Build the Worker Engine Image
docker build -t worker-service -f services/worker/Dockerfile.worker .

# 3. Execute Dynamic Benchmark Run via Automation Script
# The script templates worker.yaml with your URL and Request Count
# Usage: .\stress.ps1 -url "<TARGET_URL>" -count <NUMBER_OF_REQUESTS>
powershell.exe -File .\stress.ps1 -url "http://www.google.com" -count 100

# 4. Analysis (View P95/P99 metrics and DB insertions)
kubectl logs deployment/aggregator-deployment
