# Database Migration and Seeding

This directory contains the migration and seeding system for the vagais.ai backend.

## Overview

The migration system provides three main commands:
- **migrate**: Creates database tables and indexes
- **seed**: Populates the database with initial data
- **reset**: Drops all tables, runs migrations, and seeds data

## Usage

### Using Go directly

```bash
# Run migrations only
go run cmd/migrate/main.go migrate

# Seed database with initial data
go run cmd/migrate/main.go seed

# Reset database (drop tables, migrate, seed)
go run cmd/migrate/main.go reset
```

### Using the provided scripts

#### Linux/macOS
```bash
# Make the script executable (first time only)
chmod +x scripts/migrate.sh

# Run migrations
./scripts/migrate.sh migrate

# Seed database
./scripts/migrate.sh seed

# Reset database
./scripts/migrate.sh reset
```

#### Windows (PowerShell)
```powershell
# Run migrations
.\scripts\migrate.ps1 migrate

# Seed database
.\scripts\migrate.ps1 seed

# Reset database
.\scripts\migrate.ps1 reset
```

## Seeded Data

The seeding process creates the following initial data:

### Organizations
- **AGAI Studio** (Enterprise) - Official AGAI development studio
- **AI Research Lab** (Pro) - Research organization focused on AI development
- **Startup Inc** (Basic) - Innovative startup building AI solutions

### Users
- **admin@agai.studio** (Admin) - Admin user for AGAI Studio
- **developer@agai.studio** (Developer) - Developer user for AGAI Studio
- **researcher@airesearchlab.com** (Researcher) - Researcher for AI Research Lab
- **founder@startupinc.ai** (Founder) - Founder for Startup Inc

**Default password for all users**: `password`

### Agents
- **Customer Support Bot** - Free AI-powered customer support agent
- **Data Analysis Assistant** - Paid advanced data analysis agent
- **Content Writer Pro** - Subscription-based content writing assistant

### Licenses
- Enterprise license for AGAI Studio
- Pro license for AI Research Lab
- Basic license for Startup Inc

## Database Schema

The migration creates the following tables:

- `organizations` - Company/team information
- `users` - User accounts and profiles
- `agents` - AI agents and their configurations
- `reviews` - User reviews for agents
- `executions` - Agent execution logs
- `licenses` - Software licenses
- `payments` - Payment transactions
- `subscriptions` - User subscriptions
- `analytics` - Usage analytics
- `webhooks` - Webhook configurations
- `notifications` - User notifications

## Environment Variables

The migration system uses the same environment variables as the main application:

- `DATABASE_TYPE` - Database type (postgres, sqlite)
- `DB_HOST` - Database host
- `DB_PORT` - Database port
- `DB_USER` - Database user
- `DB_PASSWORD` - Database password
- `DB_NAME` - Database name
- `DB_SSLMODE` - Database SSL mode

## Notes

- The seeding process is idempotent - running it multiple times won't create duplicate data
- The reset command will **permanently delete all data** - use with caution
- All seeded users have the password `password` - change this in production
- The system supports both PostgreSQL and SQLite databases

