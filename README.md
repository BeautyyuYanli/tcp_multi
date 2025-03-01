# TCP Multi Proxy

A TCP proxy that can accept connections on one port and forward them to multiple backend addresses in a round-robin manner.

## Features

- Accept TCP connections on a single port
- Forward connections to multiple backend servers
- Round-robin load balancing
- Configurable connection settings
- Graceful shutdown

## Installation

```bash
# Clone the repository
git clone https://github.com/BeautyyuYanli/tcp_multi.git
cd tcp_multi

# Build the binary
go build -o tcp_proxy cmd/one2many/main.go
```

## Configuration

The application is configured using a YAML file. By default, it looks for `config.yaml` in the current directory, but you can specify a different path using the `-config` flag.

Example configuration:

```yaml
# Server settings
server:
  host: "0.0.0.0"
  port: 8080

# Connection settings
connection:
  buffer_size: 4096
  keep_alive: true
  keep_alive_period: 30s
  
# Logging settings
logging:
  level: "info"  # debug, info, warn, error
  format: "json"  # text, json
  output: "stdout"  # stdout, file
  # file_path: "./logs/tcp_multi.log"  # Optional: only used when output is "file"

# Proxy settings
proxy:
  backends:
    - "127.0.0.1:8081"
    - "127.0.0.1:8082"
    - "127.0.0.1:8083"
```

## Usage

```bash
# Run with default configuration
./tcp_proxy

# Run with a custom configuration file
./tcp_proxy -config /path/to/config.yaml
```

## How It Works

1. The proxy listens for incoming TCP connections on the configured host and port.
2. When a client connects, the proxy selects a backend server using round-robin load balancing.
3. The proxy establishes a connection to the selected backend server.
4. Data is copied bidirectionally between the client and the backend server.
5. When either the client or the backend server closes the connection, the proxy closes both connections.

## Development

```bash
# Run tests
go test ./...

# Build for different platforms
GOOS=linux GOARCH=amd64 go build -o tcp_proxy-linux-amd64 cmd/one2many/main.go
GOOS=darwin GOARCH=amd64 go build -o tcp_proxy-darwin-amd64 cmd/one2many/main.go
GOOS=windows GOARCH=amd64 go build -o tcp_proxy-windows-amd64.exe cmd/one2many/main.go
```
