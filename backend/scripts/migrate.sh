#!/bin/bash

# Migration script for merv.one Backend
# Usage: ./scripts/migrate.sh [migrate|seed|reset]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if command is provided
if [ $# -eq 0 ]; then
    echo "Usage: $0 [migrate|seed|reset]"
    echo ""
    echo "Commands:"
    echo "  migrate  - Run database migrations"
    echo "  seed     - Seed database with initial data"
    echo "  reset    - Reset database (drop tables, migrate, seed)"
    exit 1
fi

COMMAND=$1

# Check if we're in the backend directory
if [ ! -f "go.mod" ]; then
    print_error "Please run this script from the backend directory"
    exit 1
fi

print_status "Starting migration process..."

case $COMMAND in
    "migrate")
        print_status "Running database migrations..."
        go run cmd/migrate/main.go migrate
        print_success "Migrations completed successfully"
        ;;
    "seed")
        print_status "Seeding database with initial data..."
        go run cmd/migrate/main.go seed
        print_success "Seeding completed successfully"
        ;;
    "reset")
        print_warning "This will reset the database and all data will be lost!"
        read -p "Are you sure you want to continue? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            print_status "Resetting database..."
            go run cmd/migrate/main.go reset
            print_success "Database reset completed successfully"
        else
            print_status "Reset cancelled"
            exit 0
        fi
        ;;
    *)
        print_error "Unknown command: $COMMAND"
        echo "Usage: $0 [migrate|seed|reset]"
        exit 1
        ;;
esac

print_success "Migration process completed!"



