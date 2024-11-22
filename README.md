# KOL Resource Service

By implementing basic API applications, it is not only for practicing DDD model, but also as a template for any future projects.

## Tech Stack

- Go 1.23+
- PostgreSQL 17.0
- SQLBoiler ORM
- Docker & Docker Compose
- Zerolog for structured logging
- Comprehensive test coverage with race/leak detection
- Trivy/Govulncheck for security scanning
- Makefile for build, test, lint, etc.
- Github Actions for CI

## Prerequisites

- Go 1.23+
- Docker and Docker Compose
- Make

## Getting Started

1. Clone the repository:
    ```sh
    git clone https://github.com/StanleyIsMe/kol-resource.git
    ```

2. Clone the `config/api/base.yaml` file to `config/api/local_docker.yaml` and fill in the details:
    ```sh
    cp config/api/base.yaml config/api/local_docker.yaml
    ```

3. Build the Docker image:
    ```sh
    make build
    ```

4. Run the Docker Compose file:
    ```sh
    make up
    ```
