# SecureOps Lite

SecureOps Lite is a security-focused full-stack application for tracking assets, importing vulnerability intelligence, calculating risk, and explaining security posture in one place.

The project is being built around a practical cybersecurity workflow:

- store assets with meaningful product details
- import relevant CVEs from NVD / NIST
- assign vulnerabilities to affected assets
- calculate asset risk
- surface important alert conditions
- refresh vulnerability intelligence over time
- explain asset risk through an asset-aware chatbot

This is not intended to be a full SIEM, scanner, or enterprise vulnerability platform. The goal is a focused application that shows how secure backend design, external vulnerability data, AI-assisted matching, and small supporting services can work together in one system.

## What The Product Does

SecureOps Lite is designed to help a user answer questions like:

- What assets do I have?
- Which vulnerabilities affect them?
- Which assets are riskiest right now?
- Why is a given asset high risk?
- Which vulnerabilities came from official NVD data?
- What changed after an import, refresh, or recalculation?

The core product direction is:

- Angular frontend for the user interface
- Go with Gin and GORM as the main backend and trust boundary
- PostgreSQL for persistence
- focused Go services for narrow supporting tasks
- NVD / NIST as the vulnerability source of truth
- AI as a helper for matching, ranking, and explanation

## Current Repository Status

This repository is in active development.

What is already present in the codebase:

- Go Gin/GORM backend foundation
- JWT-based authentication flow
- asset CRUD backend
- vulnerability CRUD backend
- asset-to-vulnerability assignment backend
- layered backend boundaries with controller -> service -> repository -> database flow
- controller-owned service interfaces and service-owned repository interfaces
- package-local sentinel error files for repository, service, middleware, and security concerns
- Go risk scoring service
- Docker Compose for PostgreSQL, backend, and risk service
- Angular frontend project structure

What is planned next or still in progress:

- Angular authentication and main screens
- NVD / NIST vulnerability import
- AI-assisted product matching and relevance support
- asset-scoped chatbot
- `alert-service-go`
- `cve-sync-service-go`
- full multi-service Docker integration

If you want the implementation-level plan, see `Roadmap.md` and `architecture.md`.

## Planned Architecture

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
  +--> risk-service
  |
  +--> alert-service-go
  |
  +--> cve-sync-service-go
  |
  +--> NVD / NIST APIs
  |
  `--> AI provider API
```

The backend remains the main security and orchestration boundary.

That means:

- Angular should not call NVD directly
- Angular should not call AI providers directly
- Angular should not call Go services directly
- the Go backend handles authorization, validation, persistence, and external-service coordination

## Main Features

### Authentication

The application is designed around:

- user registration
- user login
- JWT-based authentication
- protected backend routes
- protected frontend routes

### Asset Inventory

Assets are intended to hold both internal tracking information and product identity information.

Examples of asset fields:

- asset name
- asset type
- vendor
- product
- version
- IP address
- operating system
- owner
- criticality
- risk score
- risk level

Why the product fields matter:

An internal label like `Firewall-01` is useful for inventory, but not enough for vulnerability matching. To pull relevant CVEs from NVD, the system needs product-aware fields such as `vendor`, `product`, and `version`.

### Vulnerability Tracking

Vulnerabilities can exist in the system in two ways:

- manually created records
- imported records from NVD / NIST

Vulnerability records are intended to include:

- CVE ID
- title
- severity
- description
- status
- CVSS details where available
- source metadata
- publish and update timestamps where useful

### NVD / NIST Vulnerability Import

The planned import flow is:

1. create an asset
2. provide product-aware details
3. trigger `Find Vulnerabilities from NVD`
4. map the asset to a likely CPE
5. import matching CVEs
6. store them locally
7. assign them to the asset
8. recalculate risk

NVD is intended to be the vulnerability source of truth.

AI may help with:

- product name normalization
- candidate CPE ranking
- CVE relevance review
- user-facing explanation

AI is not intended to invent vulnerabilities or silently override official NVD data.

### Risk Scoring

The current dedicated Go service calculates an asset risk score from summarized vulnerability data.

Example factors:

- number of critical vulnerabilities
- number of high vulnerabilities
- number of medium vulnerabilities
- number of low vulnerabilities
- asset criticality

The risk score is intended to help prioritize attention across assets.

### Alerting

The planned `alert-service-go` will focus on security events such as:

- newly imported critical CVEs
- assets crossing a risk threshold
- repeated sync failures
- important state changes that should be surfaced in the UI

### CVE Refresh

The planned `cve-sync-service-go` will refresh locally imported vulnerability data so the application does not drift too far from upstream NVD changes over time.

### Asset Chatbot

The planned chatbot is intended to be:

- asset-scoped
- grounded in local application data
- read-only in the first version

It should answer questions like:

- What vulnerabilities affect this asset?
- Why is this asset critical?
- Which CVEs matter most here?
- What changed after the last import or sync?

## Security Approach

SecureOps Lite is being built with a security-first mindset.

Key security principles in this repository:

- hash passwords with BCrypt
- use JWT for authenticated backend access
- enforce authorization on the backend
- validate security-relevant input server-side
- use DTOs instead of exposing internal entities directly
- keep secrets out of source control
- use environment-based configuration
- keep AI and external API keys server-side only
- treat AI output as advisory text, not trusted system truth
- keep external vulnerability data grounded in official NVD records
- use safe error handling and avoid leaking stack traces or secrets

There is also a lightweight WAF-style filter in the backend plan to block obviously suspicious request patterns such as simple SQL injection-like strings, XSS-like input, and path traversal attempts.

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
|   |   |-- database/
|   |   |-- middleware/
|   |   |-- model/
|   |   |-- repository/
|   |   |-- response/
|   |   |-- security/
|   |   `-- service/
|   |-- main.go
|   |-- Dockerfile
|   |-- go.mod
|   `-- go.sum
`-- risk-service/
    |-- api/
    |-- main.go
    |-- Dockerfile
    `-- go.mod
```

Current backend package rules:

- controllers handle HTTP request parsing, route parameters, and responses
- controllers depend on service interfaces defined in `controller/service_interfaces.go`
- services handle business validation and repository-error translation
- services depend on repository interfaces defined in `service/repository_interfaces.go`
- repositories handle GORM/database access only
- DTO structs currently live in `model/*_dto.go`, separate from controller logic
- `errors.go` files stay simple: error struct with `Message string`, one `Error()` method, and sentinel vars
- config does not currently need `config_errors.go` because config loading does not return errors yet

Planned additions:

- `alert-service-go/`
- `cve-sync-service-go/`

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
- `risk-service`

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
- service URLs
- later: NVD API key
- later: AI provider API key

Important:

- do not commit real secrets
- do not expose API keys to the frontend
- treat `.env` as local development configuration only

## API Direction

Current and planned API areas include:

### Authentication

- `POST /api/auth/register`
- `POST /api/auth/login`

### Assets

- `GET /api/assets`
- `GET /api/assets/{id}`
- `POST /api/assets`
- `PUT /api/assets/{id}`
- `DELETE /api/assets/{id}`
- `POST /api/assets/{id}/import-nvd-vulnerabilities`
- `POST /api/assets/{id}/calculate-risk`
- `GET /api/assets/{id}/alerts`
- `POST /api/assets/{id}/chat`

### Vulnerabilities

- `GET /api/vulnerabilities`
- `GET /api/vulnerabilities/{id}`
- `POST /api/vulnerabilities`
- `PUT /api/vulnerabilities/{id}`
- `DELETE /api/vulnerabilities/{id}`

### Assignment

- `POST /api/assets/{assetId}/vulnerabilities/{vulnerabilityId}`
- `DELETE /api/assets/{assetId}/vulnerabilities/{vulnerabilityId}`

## Data Model Direction

The main data model is built around:

- `users`
- `assets`
- `vulnerabilities`
- `asset_vulnerabilities`

Likely additions:

- `alerts`
- optional `chat_sessions`
- optional `chat_messages`
- optional `waf_events`
- optional sync history tables

## Documentation

Use the docs in this repo like this:

- `README.md`: product overview and current usage guidance
- `architecture.md`: technical implementation guide
- `Roadmap.md`: progress tracker and implementation sequence
- `Agents.md`: repository-specific working instructions for the coding assistant

## Notes

- The Go backend currently uses GORM auto-migration, which is convenient for development but should later be replaced with managed migrations.
- Some features described here are planned and documented, not fully implemented yet. This README reflects both the current repo state and the intended product direction so other users can understand what exists now and what is being built next.
