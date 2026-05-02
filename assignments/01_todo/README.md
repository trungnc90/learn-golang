# Todo App

A simple todo application built in Go with multiple storage backends.

## Prerequisites

- [Go 1.22+](https://go.dev/dl/)
- [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/)

## Starting PostgreSQL with Docker

1. Start the Postgres container:

```bash
docker compose up -d
```

2. Verify it's running and healthy:

```bash
docker compose ps
```

3. Run the migration to create the tasks table:

```bash
docker compose exec -T postgres psql -U todouser -d tododb < migrations/001_create_tasks_table.sql
```

Or connect manually and run it:

```bash
docker compose exec postgres psql -U todouser -d tododb
```

Then paste the contents of `migrations/001_create_tasks_table.sql`.

4. The connection string for the Go app is:

```
postgres://todouser:todopass@localhost:5432/tododb?sslmode=disable
```

## Stopping PostgreSQL

```bash
docker compose down
```

To also remove the persisted data:

```bash
docker compose down -v
```

## Storage Backends

- **FileStore** — persists tasks to a local JSON file
- **MemoryStore** — in-memory storage (lost on restart)
- **PostgresStore** — persists tasks to a PostgreSQL database
