# G-Tunnel-Lite

A lightweight, high-performance TCP tunneling tool written in Go. Designed for Red Team operations and network penetration testing.

## 🚀 Features

- **High Concurrency**: Built on Go's goroutines, capable of handling thousands of concurrent connections with minimal memory footprint.
- **Zero Copy**: Utilizes `io.Copy` (splice/sendfile) for kernel-level data transfer, maximizing throughput.
- **Graceful Shutdown**: Properly handles TCP half-closed states to ensure data integrity.
- **Compact**: Single binary with no external dependencies.

## 🛠 Usage

### Build
```bash
go build -o gtun main.go