#!/bin/bash

# Script to simplify Docker operations

set -e

# Default environment is development
ENV=${1:-"dev"}

# Ensure we have .env file
if [ ! -f .env ]; then
  echo "No .env file found. Creating from .env.example..."
  cp .env.example .env
  echo "Please update the .env file with your actual configuration values."
fi

case $ENV in
  "dev")
    echo "Starting development environment..."
    docker compose -f docker-compose.dev.yml up -d
    ;;
  "prod")
    echo "Starting production environment..."
    docker compose -f docker-compose.prod.yml up -d
    ;;
  "down")
    echo "Stopping all containers..."
    docker compose -f docker-compose.dev.yml down
    ;;
  "clean")
    echo "Stopping and removing all containers, volumes, and built images..."
    docker compose -f docker-compose.dev.yml down -v --rmi local
    ;;
  *)
    echo "Unknown environment: $ENV"
    echo "Usage: $0 [dev|prod|down|clean]"
    exit 1
    ;;
esac 