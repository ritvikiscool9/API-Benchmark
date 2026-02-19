# API Benchmark

**API Benchmark** is a high-performance, distributed load testing tool written in Go. It is designed to simulate massive concurrent traffic against API endpoints, measure latency with microsecond precision, and aggregate real-time performance metrics.

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
![Kubernetes](https://img.shields.io/badge/kubernetes-%23326ce5.svg?style=for-the-badge&logo=kubernetes&logoColor=white)
![AWS](https://img.shields.io/badge/AWS-%23FF9900.svg?style=for-the-badge&logo=amazon-aws&logoColor=white)

## 🚀 Project Overview

The goal of this project is to build a scalable alternative to tools like JMeter or verify system reliability under stress. Unlike simple scripts, **API Benchmark** is architected as a distributed microservice system capable of coordinated attacks from multiple geographic regions using Kubernetes.

### Core Features (In Development)
* **High-Concurrency Load Generation:** Utilizes Go goroutines to spawn thousands of lightweight "virtual users" per node with minimal memory overhead.
* **Distributed Architecture:** Designed to run on AWS EKS (Elastic Kubernetes Service), scaling worker nodes horizontally to generate massive throughput.
* **Real-Time Metrics:** Aggregates latency distributions (P50, P95, P99), error rates, and throughput (RPS) in real-time.
* **Non-Blocking Data Pipeline:** Uses Go channels and buffered queues to process metric streams without locking or bottlenecking the load generator.

## 🏗 Architecture

The system is composed of three main microservices:

1.  **The Controller (Brain):** A REST API that accepts test configurations (Target URL, QPS, Duration) and orchestrates the worker nodes via gRPC.
2.  **The Worker (Muscle):** A highly optimized Go service that generates the actual HTTP traffic. It uses a worker pool pattern to manage concurrency and reports results back asynchronously.
3.  **The Aggregator (Historian):** Ingests the stream of results from all workers, calculates windowed statistics, and persists historical data to a time-series database.

## 🛠️ Tech Stack

* **Language:** Go (Golang) for its superior concurrency model and raw performance.
* **Containerization:** Docker (Multi-stage builds for minimal image size).
* **Orchestration:** Kubernetes (AWS EKS) for auto-scaling worker nodes.
* **Communication:** gRPC for internal service-to-service communication.
* **Database:** PostgreSQL (for test configurations) and TimescaleDB (for metric storage).

## ⚡ Quick Start (Local Dev)

Prerequisites:
* Go 1.22+
* Docker

### Running a Local Benchmark
Currently, the core engine can be run locally to test basic concurrency.

1.  **Clone the repository:**
    ```bash
    git clone [https://github.com/ritvikiscool9/api-benchmark.git](https://github.com/ritvikiscool9/api-benchmark.git)
    cd api-benchmark
    ```

2.  **Run the load generator:**
    ```bash
    go run main.go
    ```
    *By default, this targets `https://httpbin.org/get` with 50 concurrent requests.*

