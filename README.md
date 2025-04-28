# Makerble Medical System API

A REST API for managing a medical clinic system with doctors, receptionists, and patients.

## Features

- User authentication (login/logout)
- Role-based access control (Doctors and Receptionists)
- Patient management
  - Create patients (Receptionists only)
  - List all patients
  - Get patient details
  - Update patient basic information (Receptionists only)
  - Update patients fully (Doctors only)
  - Delete patients (Receptionists only)

## Tech Stack

- Go 1.24.1
- PostgreSQL
- Chi Router
- JWT Authentication
- Swagger/OpenAPI Documentation

## Prerequisites

- Go 1.24.1 or higher
- PostgreSQL
- Make

## Getting Started

1. Clone the repository:

```bash
git clone git@github.com:Daniel-Brai/makerble.git
cd makerble
```

2. Set up environment variables (Optional):

```bash
cp .env.example .env
# Edit .env with your database credentials
```

3. Run database migrations:

```bash
make migrate-up
```

4. Build and run the application:

```bash
go run cmd/api/main.go
```

The API will be available at `http://localhost:5000/api/v1`

## Demo


https://github.com/user-attachments/assets/eaf3ff38-c2f1-4879-98f5-67a8a4db0d46



## API Documentation

The API documentation is available via Swagger UI at:

```
http://localhost:5000/swagger/index.html
```

To regenerate Swagger documentation after making changes:

```bash
make swagger
```

## Development

### Local Development

For hot-reloading during development:

```bash
air
```

### Docker Development

To run the application in development mode with Docker:

```bash
# Build the development image
docker build -f Dockerfile.dev -t makerble-dev .

# Run the container
docker run -p 5000:5000 -v $(pwd):/app makerble-dev
```

Or using docker-compose:

```bash
docker-compose -f docker-compose.dev.yml up
```

## Available Make Commands

- `make migrate-create name=<migration_name>` - Create a new migration
- `make migrate-up` - Run all pending migrations
- `make migrate-down` - Roll back the last migration
- `make migrate-force version=<version>` - Force migration version
- `make migrate-version` - Show current migration version
- `make swagger` - Generate Swagger documentation

## License

This project is licensed under the MIT License - see the [LICENSE-MIT](LICENSE-MIT) file for details.
