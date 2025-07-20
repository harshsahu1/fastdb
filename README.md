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
## Benchmarks

| Benchmark              | Ops/sec (approx) | Latency (ns/op) |
| ---------------------- | ---------------- | --------------- |
| `BenchmarkSet`         | ~2.86 million    | **389.3 ns**    |
| `BenchmarkGet`         | ~5.83 million    | **236.8 ns**    |
| `BenchmarkSetParallel` | ~8.53 million    | **147.9 ns**    |
| `BenchmarkPubSubSet`   | ~5.90 million    | **208.3 ns**    |

## ⚠️Note: These benchmarks reflect in-process performance only. FastDB was not tested in server mode for this benchmark run, so results do not include TCP/network or inter-process communication overhead.
## ✨ WorkInProgress

