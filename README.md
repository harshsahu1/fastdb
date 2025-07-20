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

## ✨ Features

- 🔑 **Sharded key-value storage** using CRC32 for even key distribution
- 🔒 **Thread-safe design** using per-shard `sync.Mutex` for maximum concurrency
- 📢 **Built-in Pub/Sub manager** to broadcast updates to subscribers
- ⚡️ **High performance** with parallel-safe operations
- 🧪 **Benchmarked** to validate performance under load

