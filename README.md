# Golang Load Balancer

A high-performance load balancer written in Go that supports multiple backend servers and health checks.

Created by Rudra

## Features

- Dynamic server pool management
- Health check monitoring
- Round-robin load balancing
- Configurable backend servers
- Real-time server status monitoring

## Prerequisites

- Go 1.21 or higher
- Docker (optional)

## Installation

### Using Go

```bash
git clone https://github.com/yourusername/golang-load-balancer.git
cd golang-load-balancer
go mod download
go run main.go
```

### Using Docker

```bash
docker build -t golang-load-balancer .
docker run -p 8080:8080 golang-load-balancer
```

## Configuration

Edit `config.yaml` to configure your backend servers:

```yaml
port: 8080
backends:
  - url: "http://localhost:8081"
  - url: "http://localhost:8082"
  - url: "http://localhost:8083"
```

## Usage

1. Start the load balancer
2. Configure your backend servers
3. Access your application through the load balancer at `http://localhost:8080`

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
