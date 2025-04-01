# Docker Setup Guide

This repository includes a comprehensive Docker setup for both development and production environments.

## Quick Start

### Development Environment

```bash
# Using Docker Compose directly
make dev

# Or using the helper script
make docker-dev
```

### Production Environment

```bash
# Using Docker Compose directly
make prod

# Or using the helper script
make docker-prod
```

### Stopping the Environment

```bash
# Using Docker Compose directly
make down

# Or using the helper script
make docker-down
```

### Cleaning Up

```bash
# Using the helper script
make docker-clean
```

## Docker Configuration Overview

The Docker setup consists of the following components:

### Docker Files

- `docker/Dockerfile` - Production Dockerfile using multi-stage builds
- `docker/dev.Dockerfile` - Development Dockerfile with hot reloading
- `.dockerignore` - Lists files to exclude from Docker context for faster builds

### Docker Compose Files

- `docker-compose.yml` - Base configuration with dev and prod profiles
- `docker-compose.override.yml` - Development-specific overrides (loaded automatically)

### Environment Profiles

The Docker Compose files use profiles to manage different environments:

- **dev** (default): Development environment with hot reloading and debugging tools
- **prod**: Production environment with optimized configuration

## Services

### Development Environment Services

- **api**: Go API with hot reloading via Air
- **postgres**: PostgreSQL database
- **pgadmin**: Web interface for PostgreSQL management

### Production Environment Services

- **api-prod**: Production-ready Go API with optimized settings
- **postgres**: PostgreSQL database

## Security Features

- Non-root users in containers
- Read-only filesystem for production
- Limited container resources
- Various security options like `no-new-privileges`
- Distroless base image for production

## Resource Limits

All containers have resource limits configured to control CPU and memory usage:

- **API (Development & Production)**:

  - CPU Limit: 0.5 cores
  - Memory Limit: 256MB
  - CPU Reservation: 0.25 cores
  - Memory Reservation: 128MB

- **PostgreSQL**:
  - CPU Limit: 0.5 cores
  - Memory Limit: 512MB
  - CPU Reservation: 0.25 cores
  - Memory Reservation: 256MB

## Health Checks

All services have health checks configured to ensure they're running correctly:

- **API**: Checks `/health` endpoint
- **PostgreSQL**: Uses `pg_isready` to verify database is accepting connections

## Volumes

- **postgres_data**: Persists PostgreSQL data
- **go-modules**: Caches Go modules in development

## Networks

Services communicate over a dedicated bridge network with a fixed subnet.

## Environment Variables

Configuration is managed through environment variables. Create a `.env` file based on `.env.example` to customize settings.

## Helper Scripts

The `scripts/docker-start.sh` script provides an easy way to:

- Start the development environment
- Start the production environment
- Stop all containers
- Clean up containers, volumes, and images

## Debugging

For development debugging:

- Port 2345 is exposed for debugger connections
- Debug logs are enabled
- Source code is mounted for immediate reflection of changes
