# Go API Example - Builder + Clean Architecture

This project is an example HTTP API written in Go. It uses a Clean Architecture-inspired structure and a `Builder` that assembles domain entities, routes, health checks, and shutdown hooks.

The main goal is to keep application composition separate from domain behavior:

- `cmd/server`: API entrypoint.
- `env`: environment loading and validation.
- `internal/builder`: domain entity assembly on top of the main router.
- `internal/domain`: shared contract for domain entities.
- `internal/domain/app`: example domain with handlers, services, repository, DTOs, and routes.
- `internal/infra/db`: database connection infrastructure.
- `internal/errors`: structured API error responses.
- `internal/json_schema`: request payload validation with JSON Schema.
- `internal/logger`: simple logging abstraction.

## Architecture

The server starts in `cmd/server/main.go` and follows this flow:

1. Create the context and logger.
2. Load configuration with `env.New`.
3. Create the `chi` router.
4. Open the database connection through `internal/infra/db`.
5. Create the `builder`.
6. Initialize the `app` domain entity.
7. Register the domain entity with `AddDomainEntity`.
8. Run `Build`, which mounts routes under `/api/v1`.
9. Start the HTTP server and handle graceful shutdown with `SIGINT` and `SIGTERM`.

The builder contract is defined in `internal/domain/domain_interfaces.go`:

```go
type DomainIF interface {
	Build(ctx context.Context, r *chi.Mux) error
	Health(ctx context.Context) error
	Close(ctx context.Context) error
}
```

With this contract, each new domain can expose its own routes, health check, and close routine without coupling that logic to `main.go`.

## Builder

The `internal/builder` package centralizes API assembly:

- receives the main router;
- receives domain entities through `AddDomainEntity`;
- creates an internal API router;
- calls `Build` on each domain entity;
- mounts everything under the `/api/v1` prefix;
- registers standardized `MethodNotAllowed` and `NotFound` handlers;
- adds the global `/health` route;
- calls `Close` on domain entities during shutdown.

This pattern keeps `main.go` focused on composition while each domain owns its internal layers.

## App Domain

The example domain lives in `internal/domain/app` and is split into the following layers:

- `app.go`: initializes repository, services, and handlers.
- `app_routes.go`: registers HTTP routes for the domain.
- `app_handlers.go`: receives HTTP requests, reads the body, and validates the schema.
- `app_services.go`: contains application use cases.
- `app_repository.go`: isolates database access.
- `dto/create.go`: defines DTOs and the JSON Schema for the create endpoint.
- `app_interfaces.go`: defines internal domain interfaces.

The current create endpoint validates the request payload and returns a simple handler response. The repository `Create` method is still a placeholder.

## Endpoints

### Health Check

```http
GET /health
```

Runs `Health` on every domain entity registered in the builder. In the `app` domain, this performs a `PingContext` against the database.

Successful response:

```text
OK
```

### Create App

```http
POST /api/v1/app
Content-Type: application/json
```

Valid payload:

```json
{
    "name": "My App",
    "nestedField": {
        "data": "Some data"
    }
}
```

Current successful response:

```text
hello from App handler
```

When the payload does not pass JSON Schema validation, the API returns `400` with a list of errors:

```json
{
    "errors": [
        {
            "field": "name",
            "message": "name is required"
        }
    ]
}
```

## Configuration

The application reads environment variables and also accepts a `.env` file at the project root.

Required variables:

```env
ENVIRONMENT=development
PORT=8080
DB_DRIVER_NAME=sqlite
DB_CONNECTION_STRING=./db.sqlite
```

The project imports the `modernc.org/sqlite` driver, so the configuration above works for local SQLite execution.

## Running

Requirements:

- Go `1.26.2`, as defined in `go.mod`.

Install dependencies:

```sh
go mod download
```

Run the API:

```sh
make dev
```

Or run it directly:

```sh
go run ./cmd/server/main.go
```

With the default configuration, the API is available at:

```text
http://localhost:8080
```

## curl Examples

Health check:

```sh
curl -i http://localhost:8080/health
```

Create with a valid payload:

```sh
curl -i -X POST http://localhost:8080/api/v1/app \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "My App",
    "nestedField": {
      "data": "Some data"
    }
  }'
```

Create with an invalid payload:

```sh
curl -i -X POST http://localhost:8080/api/v1/app \
  -H 'Content-Type: application/json' \
  -d '{}'
```

## Build

Build the server binary:

```sh
make build-server
```

This creates a `server` binary at the project root.

## Tests

Run:

```sh
go test ./...
```

Current observed state:

- tests in `internal/domain/app` pass;
- tests in `internal/errors` fail because they expect the `code` field in the serialized JSON, but the current `AppError.MarshalJSON` implementation omits that field.

## Adding a New Domain

To add another domain using the same model:

1. Create a package under `internal/domain/<name>`.
2. Implement handlers, services, repository, and routes.
3. Make the entity implement `domain.DomainIF`.
4. Initialize the entity in `cmd/server/main.go`.
5. Register it in the builder:

```go
builder.AddDomainEntity(newEntity)
```

The builder will include the new domain routes under `/api/v1` and include the entity in `/health` and shutdown handling.
