# SecureOps Lite Go Backend

This is the main Go backend for SecureOps Lite.

It now uses:

- Gin for routing and HTTP middleware
- GORM for PostgreSQL persistence

It follows the same current backend logic:

- `GET /api/health`
- `POST /api/auth/register`
- `POST /api/auth/login`
- JWT-protected asset and vulnerability routes
- basic WAF-style request blocking
- asset CRUD
- vulnerability CRUD
- asset-to-vulnerability assignment
- `POST /api/assets/{id}/calculate-risk`

This backend is the main API and trust boundary for the application.

## Structure

The layout keeps the backend packages at the project root for a simple application structure:

```text
backend-Go/
|-- main/
|   |-- main.go
|   `-- api/
|       |-- config/
|       |-- controller/
|       |-- database/
|       |-- middleware/
|       |-- model/
|       |-- repository/
|       |-- response/
|       |-- security/
|       `-- service/
`-- risk-service/
    |-- main.go
    `-- api/
        |-- config/
        |-- controller/
        |-- model/
        |-- response/
        `-- service/
```

Package roles:

- `controller/`: Gin handlers
- `service/`: business logic
- `repository/`: GORM persistence
- `model/`: database models and request/response structs
- `security/`: JWT generation and authentication middleware
- `middleware/`: request pipeline middleware such as the WAF filter
- `database/`: PostgreSQL connection and schema setup
- `response/`: shared API error/response helpers

## Environment

The service reads these environment values:

- `DB_HOST`
- `POSTGRES_PORT`
- `POSTGRES_DB`
- `POSTGRES_USER`
- `POSTGRES_PASSWORD`
- `JWT_SECRET`
- `JWT_EXPIRATION_MS`
- `RISK_SERVICE_URL`

Default backend port: `8080`.
