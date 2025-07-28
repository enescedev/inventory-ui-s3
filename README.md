# Inventory UI S3

This project provides a React frontend and a Go backend that store and retrieve table data from an S3 compatible storage (MinIO). Everything is orchestrated using Docker Compose.

## Usage

1. Copy `.env.example` to `.env` and adjust credentials if needed.
2. Ensure `web/frontend/.env` contains `REACT_APP_API_URL=http://localhost:8080/api` so the browser can reach the backend.
3. Run `docker-compose up --build` to start the backend, frontend and MinIO services.
4. Access the frontend at `http://localhost:3000`, the backend at `http://localhost:8080` and the MinIO console at `http://localhost:9001` (default credentials are from `.env`).

The backend exposes several endpoints under `/api`:

- `GET /api/tabs` – list directories in the bucket which become tabs in the UI
- `GET /api/table/{tab}` – fetch the latest JSON or Excel file for the given tab
- `PUT /api/table/{tab}` – upload table data back as `<tab>.json`

If the configured bucket does not exist, it is created on startup. The backend exits with an error log if it cannot connect to the S3 service.
