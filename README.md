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

## Ready to Use Link (Heroku)
- Heroku: https://packman-api-55ce3a95f0f8.herokuapp.com/health

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


## Design Notes:

### Functional Requirements:
1. Calculate optimal pack combination for given quantity.
2. Allow dynamic configuration of available pack sizes with a persistent layer.
3. Provide RESTful API endpoints for interaction.
4. Ensure data consistency and integrity when updating pack sizes. Avoid race conditions.
5. Keep it simple for pack configurations, so just 1 row in the database to hold pack sizes as an array. For history tracking also have an audit table to log changes with timestamps and user info just in case.

### Non-Functional Requirements(System Qualities):
1. **Availability**: The current deployment uses a single Heroku dyno with a single database instance, providing basic availability (~99.5-99.9% uptime). For higher availability (99.99%+), the system would need: multi-region deployment, database replication, load balancing, and automated failover mechanisms.
2. **Consistency and Partition Tolerance**: The system uses PostgreSQL as a single source of truth, ensuring consistency over partition tolerance (CP in CAP theorem). In case of network partitions between app and database, the service fails gracefully rather than serving stale/inconsistent data. For increasing availability we could consider caching strategies with eventual consistency, but for now consistency is prioritized because there is just a single row to hold pack sizes and not considering having much read-heavy operations.
3. **Performance**: The system is designed to handle calculations efficiently, with response times under 200ms for typical requests. Performance can be further improved by optimizing database queries and using caching strategies if needed.
4. **Scalability**: The application is stateless, allowing horizontal scaling by adding more dynos/instances behind a load balancer. The database can be scaled vertically or horizontally (read replicas) as needed.
5. **Maintainability**: The codebase is structured with clear separation of concerns (handlers, services, repositories). Dependency injection is used to facilitate testing and future enhancements.
6. **Testability**: Comprehensive unit tests cover core functionalities, ensuring reliability and facilitating future changes. Mocking is used for external dependencies to isolate tests.
7. **Security**: For production maybe authorization and authentication should be added but if this is planning to use as a public service then maybe not needed.
8. **Usability**: The API is designed to be intuitive and easy to use, with clear endpoints and JSON responses. Documentation is provided for developers to understand how to interact with the service.
9. **Observability**: Logging is implemented for key actions and errors, aiding in monitoring and debugging. Health check endpoints provide insights into service status.


### Implementation Details:
1. Used Gin framework for building RESTful APIs.
1. Added persistent layer by using PostgresDB for conveniently maintain pack configurations. It will also compatible to run with multiple instances of service.
2. Separated concerns by using handler, service and repository layers, also avoid loosely coupling by doing so.
   2.1 Handlers are responsible for request/response handling and validation.
   2.2 Services contain business logic.
   2.3 Repositories handle data persistence and retrieval.
3. Used dependency injection for better testability and maintainability.
4. Comprehensive unit tests for services and handlers to ensure functionality.
5. Graceful shutdown to handle termination signals and close resources properly. Added app interface to manage start and stop of services.
6. Health check endpoint to monitor service and database status.
7. Used environment variables for configuration management.
8. Implemented logging for better observability and debugging.
9. Middlewares:
   9.1 Logging middleware to log incoming requests and responses.
   9.2 Recovery middleware to handle panics and return appropriate error responses. Used gin's built-in recovery middleware for this.
   9.3 CORS middleware to handle cross-origin requests.
   9.4 Request ID middleware to assign unique IDs to requests for tracing and debugging.
   9.5 error handling middleware to standardize error responses.
10. slog used for structured logging since it's practically standard library and has good performance. It's a good choice for such a lightweight service
11. Dockerized the application for local development and consistent deployment environments.
12. Used Makefile for automating common tasks like building, testing, and running the application.
13. Migrations handled with simple SQL files and executed on application start for simplicity.
14. CI/CD with GitHub Actions for automated testing and deployment to Heroku.

### Improvement Ideas as project matures:
1. Implement Rate limiting to prevent abuse and ensure fair usage.
2. Observability enhancements: integrate with monitoring tools like Prometheus/Grafana for metrics, and use distributed tracing for better request tracking.
3. Caching strategies to improve performance for frequently requested data.
4. Readiness and liveness probes for better orchestration in containerized environments.

## Prerequisities
- Go 1.25
- Docker (for local testing)
- Postgres(for local testing - optional/already in docker)

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

