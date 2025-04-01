#!/bin/bash

# Script to simplify Docker operations with profiles

set -e

# Default environment is development
PROFILE=${1:-"dev"}

# Ensure we have .env file
if [ ! -f .env ]; then
  echo "No .env file found. Creating from .env.example..."
  cp .env.example .env
  echo "Please update the .env file with your actual configuration values."
fi

case $PROFILE in
  "dev")
    echo "Starting development environment..."
    docker compose up -d
    ;;
  "prod")
    echo "Starting production environment..."
    docker compose --profile prod up -d api-prod postgres
    ;;
  "down")
    echo "Stopping all containers..."
    docker compose down
    ;;
  "clean")
    echo "Stopping and removing all containers, volumes, and built images..."
    docker compose down -v --rmi local
    ;;
  *)
    echo "Unknown profile: $PROFILE"
    echo "Usage: $0 [dev|prod|down|clean]"
    exit 1
    ;;
esac 