# Secure Website Builder Backend

## Overview

This project provides the backend for the Secure Website Builder application, including a Go API and a PostgreSQL database. The setup supports easy configuration via environment variables and is ready for local development with Docker.

---

## Prerequisites

Before running the project, ensure your machine has:

- **Go 1.22** or higher installed:

  ```bash
  go version
  ```

- **Docker** installed and running:

  ```bash
  docker --version
  ```

- **Docker Compose** installed and running:

  ```bash
  docker compose version
  ```

- **Git** (to clone the repository)

---

## Clone Repository and Create `.env` File

1. Clone the repository:

```bash
git clone https://github.com/Secure-Website-Builder/Backend.git
cd Backend
```

2. Create a `.env` file in the **project root** with your configuration. Replace placeholders with your own credentials:

```env
# Application
APP_ENV=<development|production>
APP_PORT=<host-port-for-backend>

# Database
DB_USER=<your-db-user>
DB_PASSWORD=<your-db-password>
DB_NAME=<your-db-name>
DB_HOST=db

# MinIO / S3 Storage
MINIO_ENDPOINT=minio:<your-minio-port>
MINIO_PORT=<your-minio-port>
MINIO_USER=<your-minio-user>
MINIO_PASS=<your-minio-password>
MINIO_BUCKET=<your-bucket-name>

# Auth
JWT_SECRET=<your-jwt-secret>
```

> - This file stores secrets and host-specific configuration. **Do not commit it to version control.**
> - MinIO is used for image storage; local development uses the minio container.

---

## Start Backend and Database

1. Launch services with Docker Compose:

```bash
docker compose up
```

2. Access the backend on your host machine:

```
http://localhost:<APP_PORT>
```

> **Port mapping explanation:**
>
> - The container listens internally on the port exposed in the Dockerfile (`EXPOSE 8080`).
> - The host port comes from `.env` (`APP_PORT`). Docker Compose maps host port → container port.
> - Example: `.env` has `APP_PORT=9090`, container exposes `8080`. Access via `http://localhost:9090`.

3. Stop the services:

```bash
docker compose down
```

---

## MinIO / Image Storage (Local Development)

MinIO is a lightweight S3-compatible storage server used for storing product images.

### Access MinIO Web UI:

```bash
http://localhost:<MINIO_PORT>
```

---

## Optional: Seeding an Initial Admin (Local Development Only)

For local development, you may want to seed an initial admin account.

1. Create a file at:

```pgsql
internal/database/seed_admin.sql
```

2. Add the following (replace placeholders):

```sql
-- Insert initial admin (local development only)
INSERT INTO admin (email, password_hash)
VALUES (
'your-admin-email@example.com',
'<bcrypt-hashed-password>'
);
```

> ⚠️ Notes:
>
> - This file is ignored by git and must not be committed.
> - The password must be bcrypt-hashed, not plain text.
> - This script runs only on first database initialization.

---

## Notes

- The backend container mounts your local code for **live code updates**, so you don’t need to rebuild the image after code changes.
- Ensure the `.env` file is in the project root before running `docker compose up`.
