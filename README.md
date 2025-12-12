# ⚡ GoCollab: Real-Time Collaborative Editor Backend

![Go Version](https://img.shields.io/github/go-mod/go-version/Kunal-3004/go-collab-editor)
![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Architecture](https://img.shields.io/badge/architecture-Clean-green)
![Status](https://img.shields.io/badge/status-Active-success)

> A high-performance, concurrency-safe backend service that powers real-time collaborative editing (similar to Google Docs), built with **Go** and **WebSockets**.

---

## 🏗️ System Architecture

This project follows **Clean Architecture** principles to separate business logic from the transport layer.

```mermaid
graph TD
    ClientA[User A (Frontend)] <-->|WebSocket| Handler
    ClientB[User B (Frontend)] <-->|WebSocket| Handler
    
    subgraph "Delivery Layer"
        Handler[WebSocket Handler]
        Hub[Connection Hub]
    end
    
    subgraph "Domain & UseCase"
        Service[Editor Service]
        Logic[CRDT / Operational Logic]
    end
    
    subgraph "Repository"
        Repo[(In-Memory / Redis Storage)]
    end

    Handler --> Hub
    Hub --> Service
    Service --> Logic
    Service --> Repo
    Hub -.->|Broadcast Updates| ClientA
    Hub -.->|Broadcast Updates| ClientB