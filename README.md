# SecureOps Lite

SecureOps Lite is a security-focused full-stack platform for inventorying assets, importing vulnerability intelligence, and explaining risk across organizations, applications, home networks, and other inventory sources.

The project is designed as a practical cybersecurity workflow that combines:

- asset inventory with product-aware details
- vulnerability intelligence from NVD/NIST
- AI-assisted asset extraction from pasted text or file content
- an asset-scoped chatbot for security posture explanation
- a secure backend trust boundary with centralized authorization

This is not intended to be a full SIEM or enterprise vulnerability platform. The goal is a focused application that shows how secure design, external vulnerability data, and AI-assisted workflows can work together in one system.

## What SecureOps Lite Does

SecureOps Lite helps answer questions like:

- What assets do I have?
- Which vulnerabilities affect them?
- Which assets are riskiest right now?
- Why is a given asset high risk?
- Which vulnerabilities came from official NVD data?
- What changed after an import, refresh, or recalculation?

The platform is intended to support a broader inventory model than traditional organization-only tools. It can represent:

- organizations and offices
- application portfolios
- home network device inventories
- raw asset lists imported from manifests, documents, or scan outputs

## Architecture Overview

The system is intentionally split into clear responsibilities:

- Angular handles the browser UI
- the Go Gin/GORM backend handles authentication, validation, business logic, data orchestration, AI orchestration, and chat orchestration
- PostgreSQL stores application data
- focused Go services handle narrow, isolated tasks

The backend remains the main security and orchestration boundary. That means:

- Angular never calls NVD directly
- Angular never calls AI providers directly
- Angular never calls internal Go services directly
- authorization, validation, and persistence are enforced server-side

Planned architecture:

```text
Browser
  |
  v
Angular frontend
  |
  v
Go Gin/GORM API
  |
  +--> PostgreSQL
  |
  +--> organization/application/home network scoping
  |
  +--> alert-service-go
  |
  +--> cve-sync-service-go
  |
  +--> NVD / NIST APIs
  |
  `--> AI provider API
```

## Key Features

### Asset Inventory

Assets are modeled as both business inventory items and product fingerprints. Example fields include:

- name
- type
- vendor
- product
- version
- operating system
- IP address
- owner
- criticality
- risk score
- risk level

Product-aware details matter because labels like `Firewall-01` are not sufficient for vulnerability matching. To generate useful NVD results, the system needs vendor/product/version context.

### Vulnerability Management

Vulnerabilities can be:

- manually created records
- imported records from NVD/NIST

Records are intended to capture:

- CVE ID
- title
- severity
- description
- status
- CVSS details
- source metadata
- publish and update timestamps

### AI-Assisted Ingestion

SecureOps Lite is intended to accept raw asset descriptions from pasted text or file content and convert them into structured assets.

Supported input examples include:

- `package.json` dependency metadata
- network asset lists and topology documents
- local scan exports and inventory notes

The ingestion agent is intended to run entirely on the backend, where raw input is validated, sanitized, and converted into asset candidates before persistence.

### NVD / NIST Integration

The planned import flow is:

1. create an asset
2. provide product-aware details
3. trigger `Find Vulnerabilities from NVD`
4. map the asset to a likely CPE
5. import matching CVEs
6. store them locally
7. assign them to the asset
8. update risk data later when risk scoring is introduced

NVD/NIST is the vulnerability source of truth. AI may assist with normalization, ranking, and explanation, but it is not intended to invent CVEs or override official data.

### Remediation Workflow

The future roadmap includes a remediation workflow that supports:

- work orders tied to organization, asset, and vulnerability
- status, priority, due date, and remediation metadata
- checklist items for remediation steps
- suppression, exception, false-positive, and risk-acceptance handling
- remediation timelines and write-ups
- internal comments and discussion threads
- chatbot-ready context for remediation history and current state

### Risk Scoring

Risk scoring is planned for a later phase. Initial asset fields are prepared for scoring, but the scoring engine is deferred until after the core inventory and vulnerability flows are stable.

Example risk factors:

- critical vulnerability count
- high vulnerability count
- medium vulnerability count
- low vulnerability count
- asset criticality

### Alerting and CVE Refresh

Future services include:

- `alert-service-go` for security event notification
- `cve-sync-service-go` for refreshing imported CVEs so local data remains aligned with upstream NVD changes

### Asset Chatbot

A planned asset-scoped chatbot will provide read-only explanations for security posture using local application data. Example questions include:

- What vulnerabilities affect this asset?
- Why is this asset critical?
- Which CVEs matter most here?
- What changed after the last import or sync?

## Current Status

The repository is in active development.

Already present in the codebase:

- Go Gin/GORM backend foundation
- JWT-based authentication
- role and permission middleware foundation
- asset CRUD backend
- vulnerability CRUD backend
- asset-to-vulnerability assignment flow
- controller -> service -> repository layering
- GORM AutoMigrate schema provisioning
- Docker Compose support for PostgreSQL and backend
- Angular frontend project structure

Planned next work:

- organization- and application-aware multi-tenant scoping
- frontend authentication and main screens
- NVD/NIST vulnerability import
- AI-assisted product matching and ingestion
- asset-scoped chatbot
- remediation workflow management
- `alert-service-go`
- `cve-sync-service-go`
- full multi-service Docker integration
- AWS service integration in a later phase

For implementation-level details, see `Roadmap.md` and `architecture.md`.

## Security Approach

SecureOps Lite is built with a security-first mindset.

Core security principles:

- hash passwords with BCrypt
- use JWT for backend authentication
- enforce authorization on the backend
- keep admin-only permissions in backend middleware, not frontend-only checks
- validate security-relevant input server-side
- keep AI provider keys and external service credentials on the server
- use DTOs instead of exposing internal entities directly
- keep secrets out of source control
- use environment-based configuration

- keep AI and external API keys server-side only
- treat AI output as advisory text, not trusted system truth
- keep external vulnerability data grounded in official NVD records
- use safe error handling and avoid leaking stack traces or secrets

There is also a lightweight `RequestFilter` middleware in the backend plan to block obviously suspicious request patterns such as simple SQL injection-like strings, XSS-like input, and path traversal attempts.

## Repository Structure

```text
secureops-lite/
|-- frontend-angular/
|-- backend-Go/
|-- docker-compose.yml
|-- .env
|-- README.md
|-- Roadmap.md
`-- architecture.md
```

Inside `backend-Go/`:

```text
backend-Go/
|-- main/
|   |-- api/
|   |   |-- config/
|   |   |-- controller/
|   |   |-- dto/
|   |   |-- middleware/
|   |   |-- model/
|   |   |-- repository/
|   |   |-- security/
|   |   |-- service/
|   |   `-- utils/
|   |-- main.go
|   |-- Dockerfile
|   |-- go.mod
|   `-- go.sum
```

Current backend package rules:

- controllers handle HTTP request parsing, route parameters, and responses
- controllers receive service interfaces directly from constructors
- controllers call service interfaces exposed by the `service` package
- services handle business validation, ownership checks, and repository-error translation
- services depend on repository interfaces defined in `service/repository_interfaces.go`
- repositories handle GORM/database access only
- DTO structs live in `dto/`, separate from controller logic and database models
- backend uses GORM AutoMigrate at startup instead of maintaining separate SQL migration scripts
- `errors.go` files stay simple: error struct with `Message string`, one `Error()` method, and sentinel vars
- `permissions.go` contains route-level permission middleware such as `RequireAdmin`
- normal registered users default to the `user` role; admin access must not be granted by client-controlled registration data
- config does not currently need `config_errors.go` because config loading does not return errors yet

Planned additions:

- `alert-service-go/`
- `cve-sync-service-go/`
- service-to-service authentication for internal Go service calls, using a server-side shared token or stronger handshake pattern later

## Running The Current Project

The current runnable backend stack in this repository is based on Docker Compose.

### Prerequisites

- Docker Desktop
- Node.js for local Angular work
- Go for local backend and service work if not using Docker only

### Current Compose Services

The current `docker-compose.yml` defines:

- `postgres`
- `backend`

Start them with:

```bash
docker compose up --build
```

Current port notes:

- backend: `http://localhost:8080`
- PostgreSQL: mapped from `${POSTGRES_PORT}` to container `5432`

### Frontend

The Angular frontend exists in `frontend-angular/`, but it is not yet wired into Docker Compose in the current state of the repo.

At this stage, the frontend should be treated as in-progress application work rather than a finished production-ready UI.

## Environment Configuration

This repository uses a local `.env` file for development configuration.

Typical values include:

- PostgreSQL database name
- PostgreSQL username
- PostgreSQL password
- PostgreSQL port
- JWT secret
- JWT expiration
- later: service URLs
- later: NVD API key
- later: AI provider API key

Important:

- do not commit real secrets
- do not expose API keys to the frontend
- treat `.env` as local development configuration only

## API Direction

Implemented API areas:

### Authentication

- `POST /api/auth/register`
- `POST /api/auth/login`

### Assets

- `GET /api/assets`
- `GET /api/assets/{id}`
- `POST /api/assets`
- `PUT /api/assets/{id}`
- `DELETE /api/assets/{id}`

Asset endpoints are currently scoped to the authenticated user. The future multi-tenant model will scope assets, vulnerabilities, work orders, and comments to the organization the user belongs to.

### Vulnerabilities

- `GET /api/vulnerabilities`
- `GET /api/vulnerabilities/{id}`
- `POST /api/vulnerabilities`
- `PUT /api/vulnerabilities/{id}`
- `DELETE /api/vulnerabilities/{id}`

### Assignment

- `POST /api/assets/{assetId}/vulnerabilities/{vulnerabilityId}`
- `DELETE /api/assets/{assetId}/vulnerabilities/{vulnerabilityId}`

Planned API areas:

- `POST /api/assets/{id}/import-nvd-vulnerabilities`
- `GET /api/assets/{id}/alerts`
- `POST /api/assets/{id}/chat`
- `POST /api/organizations/{orgId}/work-orders`
- `GET /api/organizations/{orgId}/work-orders`
- `GET /api/organizations/{orgId}/work-orders/{id}`
- `PATCH /api/organizations/{orgId}/work-orders/{id}`
- `DELETE /api/organizations/{orgId}/work-orders/{id}`
- `POST /api/organizations/{orgId}/work-orders/{id}/checklist-items`
- `PATCH /api/organizations/{orgId}/work-orders/{id}/checklist-items/{itemId}`
- `DELETE /api/organizations/{orgId}/work-orders/{id}/checklist-items/{itemId}`
- `POST /api/organizations/{orgId}/work-orders/{id}/exceptions`
- `PATCH /api/organizations/{orgId}/work-orders/{id}/exceptions/{exceptionId}`
- `DELETE /api/organizations/{orgId}/work-orders/{id}/exceptions/{exceptionId}`
- `POST /api/organizations/{orgId}/work-orders/{id}/entries`
- `GET /api/organizations/{orgId}/work-orders/{id}/entries`
- `PATCH /api/organizations/{orgId}/work-orders/{id}/entries/{entryId}`
- `DELETE /api/organizations/{orgId}/work-orders/{id}/entries/{entryId}`
- `POST /api/organizations/{orgId}/assets/{assetId}/comments`
- `POST /api/organizations/{orgId}/vulnerabilities/{vulnerabilityId}/comments`
- `POST /api/organizations/{orgId}/work-orders/{id}/comments`
- `GET /api/organizations/{orgId}/assets/{assetId}/comments`
- `GET /api/organizations/{orgId}/vulnerabilities/{vulnerabilityId}/comments`
- `GET /api/organizations/{orgId}/work-orders/{id}/comments`
- `POST /api/sync/nvd`
- `GET /api/alerts`
- `PATCH /api/alerts/{id}/acknowledge`

## Data Model Direction

The main data model is built around:

- `organizations`
- `users`
- `assets`
- `vulnerabilities`
- `asset_vulnerabilities`

Likely additions:

- `alerts`
- `work_orders`
- `work_order_checklist_items`
- `vulnerability_exceptions`
- `remediation_entries`
- `comments`
- optional `chat_sessions`
- optional `chat_messages`
- optional sync history tables

## Documentation

Use the docs in this repo like this:

- `README.md`: product overview and current usage guidance
- `architecture.md`: technical implementation guide
- `Roadmap.md`: progress tracker and implementation sequence
- `Agents.md`: repository-specific working instructions for the coding assistant

## Notes

- The Go backend currently uses GORM AutoMigrate at startup instead of separate embedded SQL migration files. A more advanced migration tool can be added later if the schema becomes more complex.
- Some features described here are planned and documented, not fully implemented yet. This README reflects both the current repo state and the intended product direction so other users can understand what exists now and what is being built next.
