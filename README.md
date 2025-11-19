# packman-api
Pack calculation service API built with Go and deployed to Heroku using Docker containers.

## Overview
The Packman API is a RESTful service that optimizes package fulfillment by calculating the most efficient combination of pack sizes for any given order quantity. The system aims to minimize the number of packs while ensuring all items are shipped.

## Features
- Calculate packs for given quantities
- Manage pack size configurations
- RESTful API with JSON responses
- PostgreSQL database storage
- Comprehensive test coverage for services and handlers
- Docker containerized deployment
- Automated CI/CD with GitHub Actions

### Core Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/v1/calculate` | Calculate optimal pack combination for an order quantity |
| `GET` | `/api/v1/pack-sizes` | Retrieve current pack size configuration |
| `PUT` | `/api/v1/pack-sizes` | Update pack size configuration |
| `GET` | `/health` | Check service and database health status |

### Technology Stack
- Language: Go 1.25
- Database: PostgreSQL
- Deployment: Docker, Heroku


For API docs please refer to [API.md](API.md)


## Prerequisities
- Go 1.25
- Docker (for local testing)
- Postgres(optional for local testing)

## Local Development

### Quick Start

```bash
# Clone the repository
git clone https://github.com/nsaltun/packman.git
cd packman

# Run with Go
make postgres-up
make run

# Or build and run
make postgres-up
make build
make run
```

### Running with Docker

```bash
# to run
make docker-up
# to stop container
make docker-down
# to remove containers and volumes
make docker-clean
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -v -race -coverprofile=coverage.out ./...

# View coverage report
go tool cover -html=coverage.out
```

## API Endpoints
For API docs please refer to [API.md](API.md)

## Deployment

1. **Install Heroku CLI**
   ```bash
   brew tap heroku/brew && brew install heroku
   heroku login
   ```

2. **Create Heroku App**
   ```bash
   heroku create packman-api
   heroku stack:set container -a packman-api
   ```

3. **Setup Database (Free Option)**
   - Create a free PostgreSQL database at [Supabase](https://supabase.com)
   - Get your connection string from Project Settings → Database
   - Set it in Heroku:
   ```bash
   heroku config:set DATABASE_URL="postgresql://postgres:[PASSWORD]@db.[PROJECT-REF].supabase.co:5432/postgres" -a packman-api
   ```

4. **Connect GitHub (via Dashboard)**
   - Go to: https://dashboard.heroku.com/apps/packman-api/deploy/github
   - Connect your GitHub repository
   - Enable "Automatic Deploys" from `main` branch
   - Enable "Wait for CI to pass before deploy"

## Deploy

Just push to GitHub:
```bash
git push origin main
```

## Monitor

```bash
heroku logs --tail -a packman-api       # View logs
heroku open -a packman-api              # Open app
heroku ps -a packman-api                # Check status
heroku releases -a packman-api          # View releases
heroku rollback -a packman-api          # Rollback if needed
```

