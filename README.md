# Inventory UI S3

This project provides a simple Go backend that stores and retrieves a `table.json` file from an S3 compatible storage (MinIO). It is orchestrated using Docker Compose.

## Usage

1. Copy `.env.example` to `.env` and adjust credentials if needed.
2. Run `docker-compose up --build` to start the backend and MinIO services.
3. Access the backend at `http://localhost:8080` and MinIO console at `http://localhost:9001` (default credentials are from `.env`).

The backend exposes two endpoints:

- `GET /table` – retrieve `table.json` from the configured bucket
- `PUT /table` – upload a new `table.json` to the bucket

If the configured bucket does not exist, it is created on startup. The backend exits with an error log if it cannot connect to the S3 service.
