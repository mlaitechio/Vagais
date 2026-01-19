# Migration script for AGAIS.AI Backend
# Usage: .\scripts\migrate.ps1 [migrate|seed|reset]

param(
    [Parameter(Mandatory=$true)]
    [ValidateSet("migrate", "seed", "reset")]
    [string]$Command
)

# Colors for output
$Red = "Red"
$Green = "Green"
$Yellow = "Yellow"
$Blue = "Blue"
$White = "White"

# Function to print colored output
function Write-Status {
    param([string]$Message)
    Write-Host "[INFO] $Message" -ForegroundColor $Blue
}

function Write-Success {
    param([string]$Message)
    Write-Host "[SUCCESS] $Message" -ForegroundColor $Green
}

function Write-Warning {
    param([string]$Message)
    Write-Host "[WARNING] $Message" -ForegroundColor $Yellow
}

function Write-Error {
    param([string]$Message)
    Write-Host "[ERROR] $Message" -ForegroundColor $Red
}

# Check if we're in the backend directory
if (-not (Test-Path "go.mod")) {
    Write-Error "Please run this script from the backend directory"
    exit 1
}

Write-Status "Starting migration process..."

switch ($Command) {
    "migrate" {
        Write-Status "Running database migrations..."
        go run cmd/migrate/main.go migrate
        Write-Success "Migrations completed successfully"
    }
    "seed" {
        Write-Status "Seeding database with initial data..."
        go run cmd/migrate/main.go seed
        Write-Success "Seeding completed successfully"
    }
    "reset" {
        Write-Warning "This will reset the database and all data will be lost!"
        $confirmation = Read-Host "Are you sure you want to continue? (y/N)"
        if ($confirmation -eq "y" -or $confirmation -eq "Y") {
            Write-Status "Resetting database..."
            go run cmd/migrate/main.go reset
            Write-Success "Database reset completed successfully"
        } else {
            Write-Status "Reset cancelled"
            exit 0
        }
    }
}

Write-Success "Migration process completed!"



