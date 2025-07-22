# ⚡️ FASTDB

				███████╗ █████╗ ███████╗████████╗██████╗ ██████╗
				██╔════╝██╔══██╗██╔════╝╚══██╔══╝██╔══██╗██╔══██╗
				█████╗  ███████║███████╗   ██║   ██║  ██║██████╔╝
				██╔══╝  ██╔══██║╚════██║   ██║   ██║  ██║██╔══██╗
				██║     ██║  ██║███████║   ██║   ██████╔╝██████╔╝
				╚═╝     ╚═╝  ╚═╝╚══════╝   ╚═╝   ╚═════╝ ╚═════╝
					⚡ FastDB — blazing fast KV store

A high-performance, in-memory sharded key-value store written in Go, with built-in Pub/Sub capabilities. Designed for concurrent access and minimal latency.

---

## 📦 Features

- Ultra-fast `GET` and `SET` operations
- Concurrent client handling
- Lightweight architecture with minimal overhead
- Redis-compatible benchmarking
- Reactive

---

## ⚙️ Test Parameters

Benchmarking was performed using `redis-benchmark`:

- **Commands Tested:** `SET`, `GET`  
- **Requests:** 10,000  
- **Concurrency:** 50  
- **Payload Size:** 3 bytes  
- **Keep-alive:** Enabled  

---

## 📊 Performance Comparison

### 🔄 Throughput (requests/sec)

| Command | FastDB     | Redis      | Winner                    |
| ------- | ---------- | ---------- | ------------------------- |
| SET     | 112,359.11 | 109,279.06 | **FastDB** (≈2.8% faster) |
| GET     | 169,491.53 | 188,679.25 | **Redis** (≈11.3% faster) |

---

### ⏱️ SET Command – Latency (ms)

| Metric | FastDB | Redis | Winner                            |
| ------ | ------ | ----- | --------------------------------- |
| avg    | 0.305  | 0.301 | **Redis** (slightly better)       |
| min    | 0.072  | 0.072 | Tie                               |
| p50    | 0.191  | 0.191 | Tie                               |
| p95    | 1.183  | 0.911 | **Redis**                         |
| p99    | 1.295  | 1.391 | **FastDB**                        |
| max    | 1.359  | 3.055 | **FastDB** (significantly better) |

---

### 🔍 GET Command – Latency (ms)

| Metric | FastDB | Redis | Winner    |
| ------ | ------ | ----- | --------- |
| avg    | 0.158  | 0.148 | **Redis** |
| min    | 0.064  | 0.056 | **Redis** |
| p50    | 0.159  | 0.151 | **Redis** |
| p95    | 0.199  | 0.183 | **Redis** |
| p99    | 0.271  | 0.207 | **Redis** |
| max    | 0.463  | 0.311 | **Redis** |
## ✨ WorkInProgress

