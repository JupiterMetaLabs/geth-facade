# Scripts

This folder contains deployment scripts, examples, and configuration files for the JMDT Geth Facade.

## Files

### `Dockerfile`
Docker configuration for containerized deployment:

- **Multi-stage Build**: Optimized build process
- **Security**: Non-root user execution
- **Performance**: Minimal image size
- **Configuration**: Environment variable support

### `examples/`
Example implementations and usage patterns:

#### `simple/`
Basic usage example showing how to integrate the facade:

- **Simple Integration**: Minimal setup example
- **Custom Backend**: Example of using a custom backend
- **Configuration**: Basic configuration examples

#### `custom-backend/`
Advanced example showing custom backend implementation:

- **Backend Interface**: How to implement the Backend interface
- **Custom Logic**: Example of custom blockchain logic
- **Integration**: How to integrate with the facade

## Usage

### Docker Deployment
```bash
# Build the image
docker build -t jmdt-geth-facade .

# Run the container
docker run -p 8545:8545 -p 8546:8546 jmdt-geth-facade
```

### Examples
```bash
# Run simple example
cd Scripts/examples/simple
go run main.go

# Run custom backend example
cd Scripts/examples/custom-backend
go run main.go
```

## Configuration

### Environment Variables
- `HTTP_PORT`: HTTP server port (default: 8545)
- `WS_PORT`: WebSocket server port (default: 8546)
- `LOG_LEVEL`: Logging level (debug, info, warn, error)
- `BACKEND_TYPE`: Backend type (memory, custom)

### Docker Compose
Example docker-compose.yml for production deployment:

```yaml
version: '3.8'
services:
  jmdt-geth-facade:
    build: .
    ports:
      - "8545:8545"
      - "8546:8546"
    environment:
      - HTTP_PORT=8545
      - WS_PORT=8546
      - LOG_LEVEL=info
    restart: unless-stopped
```

## Future Enhancements

- **Kubernetes**: Helm charts and K8s manifests
- **Monitoring**: Prometheus metrics and Grafana dashboards
- **Load Balancing**: Nginx configuration examples
- **SSL/TLS**: HTTPS and WSS configuration
- **Rate Limiting**: API rate limiting examples

## Comments

Scripts include standardized comments:

- `//debugging`: Debug configuration and logging
- `//future`: Planned script enhancements
- `//test`: Test-related configurations
- `//conversions`: Configuration format conversions
