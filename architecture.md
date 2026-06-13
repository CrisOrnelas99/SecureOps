# SecureOps Lite Architecture

## Purpose

This document explains how SecureOps Lite is intended to work, how its parts fit together, and how each major feature should be implemented. It is the technical implementation guide for the project, while `Roadmap.md` remains the progress and execution tracker.

> Implementation note: this file is intentionally more detailed than `README.md`. `README.md` explains the project at a high level. This file explains how the system should actually be built.

## System Overview

SecureOps Lite is a full-stack cybersecurity application for:

- organizing assets by organization, office, application, home network, or other inventory source
- tracking assets
- importing relevant CVEs from NVD/NIST
- assigning vulnerabilities to affected assets
- managing work orders and remediation workflows
- tracking asset risk fields for later scoring work
- monitoring risk changes over time
- refreshing vulnerability intelligence
- raising alerts for important security events
- presenting the current security picture in a dashboard
- explaining asset security posture through an asset-aware chatbot

The system is also designed to ingest raw asset descriptions from pasted text or file contents, such as package manifests, network inventory documents, or local scan exports. An AI-assisted ingestion agent will convert that content into normalized asset records and place them into the appropriate organizational, application, or home network scope.

The system is intentionally split into clear responsibilities:

- Angular handles the browser UI
- the Go Gin/GORM backend handles authentication, validation, business logic, data orchestration, NVD integration, AI orchestration, and chat orchestration
- PostgreSQL stores application data
- focused Go services handle narrow, isolated tasks that are fast, focused, and easy to reason about

> Design comment: the backend should stay the main system boundary. Even when AI or external services are added, the main Go backend should remain the place where trust decisions, authorization checks, and persistence rules are enforced.

## High-Level Architecture

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

### Request flow

1. A user interacts with the Angular frontend.
2. Angular sends HTTP requests to the Go backend API.
3. The Go backend authenticates the request, validates input, and runs business logic.
4. The Go backend reads or writes data in PostgreSQL as needed.
5. The Go backend checks organization membership and scoped authorization before touching data.
6. For NVD import flows, the Go backend calls NVD APIs and optionally an AI-assisted matching layer.
7. For security events, the Go backend may notify `alert-service-go`.
8. For scheduled CVE refresh work, the Go backend coordinates with or receives updates from `cve-sync-service-go`.
9. The Go backend returns a safe response to Angular.

> Security comment: Angular should never call NVD directly, never call AI providers directly, and never call Go services directly. This keeps secrets on the server side and keeps authorization decisions centralized.

## Repository Structure

The repository should continue following the current naming style:

```text
secureops-lite/
|-- frontend-angular/
|-- backend-Go/
|-- alert-service-go/
|-- cve-sync-service-go/
|-- docker-compose.yml
|-- .env
|-- README.md
|-- Roadmap.md
`-- architecture.md
```

> Naming comment: Future focused services such as `alert-service-go` and `cve-sync-service-go` should use the same narrow-service style.

## Project Direction

The updated project idea is:

SecureOps Lite is a cybersecurity asset risk platform that:

- separates assets by organization, application portfolio, home network, or other inventory source
- stores technical asset details
- identifies likely product matches for those assets
- imports vulnerabilities from NVD/NIST
- assigns imported CVEs to assets
- manages remediation work orders and workflow notes
- can add risk scoring later after the base app is stable
- uses AI to improve product matching, relevance review, and raw-text/file asset ingestion
- uses a chatbot to explain asset risk and vulnerability posture in plain English

> Architecture comment: the biggest shift is that vulnerabilities are no longer only manually entered data. They become intelligence-backed records that can be imported, refreshed, prioritized, explained, and monitored.

## Frontend Architecture

### Responsibilities

The Angular frontend should be responsible for:

- organization-aware navigation when a user belongs to more than one organization
- login and registration screens
- route navigation
- dashboard rendering
- asset management screens
- vulnerability management screens
- asset detail view
- NVD import actions
- alert and risk views
- chatbot UI for asset-level questions
- work order and remediation screens when those workflows are introduced
- calling backend APIs through Angular services
- attaching JWTs to protected requests

### Recommended structure

The frontend should stay feature-based and simple:

```text
frontend-angular/
|-- src/app/
|   |-- core/
|   |   |-- services/
|   |   |-- guards/
|   |   `-- interceptors/
|   |-- features/
|   |   |-- auth/
|   |   |-- dashboard/
|   |   |-- assets/
|   |   |-- vulnerabilities/
|   |   |-- alerts/
|   |   `-- chat/
|   |-- shared/
|   |   |-- models/
|   |   `-- components/
|   `-- app.routes.ts
```

### Frontend implementation details

#### Authentication

- `register` and `login` forms should post to the backend auth endpoints
- JWT handling should be centralized in an auth service
- an HTTP interceptor should attach `Authorization: Bearer <token>` on protected requests
- route guards should prevent unauthenticated access to private pages
- logout should clear local auth state and redirect to login

> Security comment: frontend route guards improve user flow, but they do not enforce real security. Backend authorization still decides what a user can access.

#### Assets

Assets should now be treated as both business objects and product fingerprints.

The asset UI should support:

- asset name
- asset type
- vendor
- product
- version
- operating system
- owner
- criticality
- exposure-related flags if added later

> Implementation note: `asset name` is usually an internal label like `Firewall-01`. That is not enough for CVE matching. `vendor`, `product`, and `version` are the fields that make NVD integration realistic.

#### Asset details

The asset detail page should become the main workflow page for:

- viewing assigned vulnerabilities
- seeing risk score and risk level
- running `Find Vulnerabilities from NVD`
- reviewing NVD match confidence
- reviewing risk scoring later when that feature is reintroduced
- reading risk explanations
- chatting about the asset

#### Vulnerabilities

- vulnerability list page should support severity and status filtering
- imported NVD vulnerabilities should be visually distinguishable from manually created records
- the UI should show source metadata such as `NVD`, `cvssScore`, `publishedAt`, and last sync time where useful

#### Dashboard

The dashboard should eventually show:

- total assets
- total vulnerabilities
- open vulnerabilities
- critical vulnerabilities
- high-risk assets
- recent alert events
- recent NVD sync activity
- risk trend highlights

## Backend Architecture

### Responsibilities

The Go Gin/GORM backend is the trust boundary and main orchestration layer. It should handle:

- registration and login
- password hashing
- JWT generation and validation
- access control
- route-level permission middleware for elevated actions
- organization membership and tenant scoping
- request validation
- organization CRUD and membership-aware data access
- asset CRUD
- vulnerability CRUD
- asset-vulnerability assignment
- work order and remediation workflow orchestration
- notes, comments, and exception handling
- NVD integration
- AI-assisted product and CVE relevance support
- AI agent ingestion for raw text and file-based asset extraction
- chatbot orchestration
- alert service integration
- sync job coordination
- safe error handling
- CORS policy
- `RequestFilter` request filtering
- audit-relevant event logging

### Recommended package layout

The project should stay separated by layer and concern:

```text
backend-Go/
|-- main/
|   |-- main.go
|   `-- api/
|       |-- config/
|       |-- controller/
|       |   |-- asset_controller.go
|       |   |-- auth_controller.go
|       |   `-- vulnerability_controller.go
|       |-- dto/
|       |   |-- asset_dto.go
|       |   |-- auth_dto.go
|       |   `-- vulnerability_dto.go
|       |-- middleware/
|       |-- model/
|       |   |-- asset.go
|       |   |-- organization.go
|       |   |-- user.go
|       |   `-- vulnerability.go
|       |-- repository/
|       |-- security/
|       |-- service/
|           |-- asset_service.go
|           |-- service_helpers.go
|           |-- repository_interfaces.go
|           `-- service_errors.go
|       `-- utils/
```

Future workflow data files can join the model and DTO layers as needed:

```text
api/model/
|-- work_order.go
|-- work_order_checklist_item.go
|-- vulnerability_exception.go
|-- remediation_entry.go
`-- comment.go
```

Future workflow-focused packages can live alongside the existing layers, such as:

```text
api/
|-- organizations/
|-- workflow/
|-- work_orders/
|-- checklist_items/
|-- vulnerability_exceptions/
|-- remediation_notes/
|-- comments/
`-- chatbot_context/
```

> Design comment: `nvd`, `ai`, and `chat` belong in the main Go backend because they depend on authorization, database rules, DTO mapping, and application-level trust decisions. They should not be pushed into separate services just because they talk to external services.

### Backend layer rules

The main Go backend should keep this dependency direction:

```text
controller -> service -> repository -> database
```

Layer responsibilities:

- `controller`: HTTP-only concerns such as Gin context handling, JSON binding, route parameter parsing, and response calls
- `service`: business validation, ownership checks, repository-error translation, and use-case coordination
- `repository`: GORM/database reads and writes only
- `utils`: database connection helpers and database error helpers. The current backend provisions the schema with GORM AutoMigrate at startup rather than a separate SQL migration folder.
- `middleware`: request filtering and Gin middleware behavior
- `security`: JWT generation, parsing, and authentication filtering
- `config`: environment-backed construction of config and dependencies

Interfaces should be owned by the consuming layer:

- service interfaces are exposed by the `service` package for controller use
- controllers receive service interfaces directly through constructors in `main.go`
- `service/repository_interfaces.go` defines the repository interfaces services need
- repository structs satisfy those interfaces implicitly
- service structs satisfy their package interfaces implicitly

This keeps controllers unaware of repository implementations and keeps repositories unaware of HTTP.

### DTO and model placement

The current code separates database/domain structs from request/response DTO structs:

- database/domain structs live in `api/model`: `asset.go`, `user.go`, `vulnerability.go`
- DTO structs live in `api/dto`: `asset_dto.go`, `auth_dto.go`, `vulnerability_dto.go`

DTO files should not live in `controller`. Controllers use DTOs, but DTOs are not controller behavior.

### Error handling layout

Package-level `errors.go` files should stay simple:

```go
type ServiceError struct {
	Message string
}

func (e ServiceError) Error() string {
	return e.Message
}

var (
	ErrInvalidRequestData = &ServiceError{Message: "invalid request data"}
)
```

Rules:

- `errors.go` files contain only the error struct, its `Error()` method, and sentinel vars
- repository errors describe repository/database outcomes only
- service errors are more general business outcomes such as invalid request, conflict, not found, invalid credentials, and forbidden
- middleware and security errors follow the same simple sentinel style
- helper functions and mapping logic belong in normal implementation files, not in `errors.go`
- config does not need `config_errors.go` until config loading returns `(Config, error)`

### Core backend flows

#### Register flow

1. Client sends registration data to `POST /api/auth/register`.
2. Backend validates required fields and uniqueness rules.
3. Password is hashed with BCrypt.
4. User is saved in PostgreSQL with the default `user` role.
5. Backend returns a safe success response without sensitive fields.

> Security comment: registration must not let a client choose `admin`. Admin users should be created through a controlled local, seeded, or admin-only process.

#### Login flow

1. Client sends credentials to `POST /api/auth/login`.
2. Backend loads the user by username or email.
3. BCrypt verifies the password.
4. Backend generates a JWT if credentials are valid.
5. Backend returns a login response DTO with token metadata.

#### Protected request flow

1. Angular sends a request with `Authorization: Bearer <token>`.
2. JWT filter extracts the token.
3. Backend validates the token.
4. Gin authentication middleware establishes the authenticated request context.
5. The authenticated database user ID is attached to `GinContext`.
6. Controller logic runs only if access is allowed.

#### Asset creation flow

1. User creates an asset.
2. Backend validates asset fields.
3. Backend attaches the authenticated user's ID to the asset.
4. Backend stores the asset in PostgreSQL.
5. Asset is available for later NVD import and risk work.

> Implementation note: asset creation and NVD import should be separate at first. That keeps the create flow fast and avoids coupling user data entry to external API latency.

#### NVD import flow

This should follow Option B:

1. User creates the asset first.
2. User clicks `Find Vulnerabilities from NVD`.
3. The Go backend loads the asset.
4. The Go backend builds a product fingerprint from `vendor`, `product`, `version`, and related metadata.
5. The Go backend searches NVD CPE data.
6. AI may help rank candidate matches when there is ambiguity.
7. The Go backend selects or confirms the best CPE.
8. The Go backend queries NVD CVEs using that selected CPE.
9. The Go backend maps the returned CVE data to local vulnerability records.
10. The Go backend assigns those vulnerabilities to the asset.
11. Risk data can be updated later when risk scoring is reintroduced.

> Security comment: the NVD API is the vulnerability source of truth. AI may help rank or explain relevance, but AI should not invent vulnerabilities or silently override official data.

#### Asset chat flow

1. User opens an asset detail page.
2. User asks a chat question.
3. Angular sends the question to `POST /api/assets/{id}/chat`.
4. The Go backend verifies the user is allowed to view that asset.
5. The Go backend loads only the relevant asset context.
6. The Go backend builds a constrained prompt from local data plus stored NVD metadata.
7. The Go backend calls the AI provider.
8. The Go backend returns a read-only natural-language response.

> Security comment: the chatbot should explain and summarize. It should not directly mutate assets, vulnerabilities, assignments, or risk scores.

## NVD / NIST Integration Design

### Why this belongs in the Go backend

NVD integration is mostly:

- HTTP integration
- DTO mapping
- validation
- persistence
- deduplication
- asset-to-product matching
- authorization-aware business logic

Those are already main-backend responsibilities.

### Asset fingerprinting

For reliable matching, an asset should include:

- `name`
- `type`
- `vendor`
- `product`
- `version`
- `operatingSystem`
- optional deployment notes or environment tags

> Concept comment: `Firewall-01` is a useful asset label, but it is not a vulnerability lookup key. NVD matching works much better when the app knows what product and version the asset actually runs.

### Recommended stored fields

On `Asset`, plan for fields like:

- `vendor`
- `product`
- `version`
- `cpeName`
- `cpeMatchConfidence`
- `cpeMatchMethod`
- `lastNvdSyncAt`

### Local vulnerability storage

NVD results should be stored locally in PostgreSQL.

Benefits:

- faster UI reads
- stable dashboard data
- historical comparison
- less rate-limit pressure
- easier relevance review
- easier risk recalculation

> Security comment: do not live-query NVD every time a page loads. That would be slow, noisy, harder to secure, and more likely to create inconsistent user-facing results.

## AI-Assisted Matching Design

### Role of AI

AI should be a supporting layer, not the truth layer.

AI can help:

- normalize vendor and product names
- rank likely CPE matches
- explain why a given CPE seems relevant
- rank CVEs by likely relevance to the asset
- explain risk in plain English

AI should not:

- invent CVEs
- replace NVD as source of truth
- bypass authorization
- silently auto-import weak matches with no confidence threshold

### AI Agent Architecture for Asset Ingestion

The AI ingestion agent should be implemented as a backend-only orchestration layer that accepts sanitized raw text or file content and converts it into structured asset records.

Key principles:

- The frontend sends only user-provided text or file payload metadata to the backend. AI provider keys, prompts, and routing remain server-side only.
- The backend validates and sanitizes incoming content before using AI, including file type restrictions, size limits, and secret-stripping from pasted data.
- Supported inputs should include package manifests, network asset lists, inventory documents, local scan exports, and other free-form asset descriptions.
- The agent should extract asset metadata such as name, type, vendor, product, version, operating system, IP address, ownership, environment, and confidence metadata.
- The backend should map extracted results into the correct scope: organization, application portfolio, home network, or other inventory context.
- Audit logs should record ingestion input sources, AI-extracted outputs, and the user who initiated the ingestion.
- AI-generated suggestions should be treated as guidance. The backend should apply deterministic rules and review thresholds before committing extracted assets into the system.

### Recommended pattern

1. deterministically query candidate CPEs from NVD
2. give the AI only the asset details and candidate options
3. ask the AI to rank the candidates and explain confidence
4. if confidence is high, continue
5. if confidence is low, require user confirmation

> Security comment: this is the safe way to use AI in a security tool. The model helps interpret ambiguity, but deterministic systems still control data integrity.

## Chatbot Design

### Purpose

The chatbot should answer asset-specific questions such as:

- what vulnerabilities affect this asset
- why this asset is high or critical risk
- which CVEs should be prioritized
- what changed after a sync or recalculation
- why a vulnerability was imported
- what the risk score means in plain English
- what work order is active for this asset
- which checklist items are done or still open
- what remediation notes or comments the team has already written

### Scope

The first version should be organization-scoped and asset-scoped, not global.

Recommended endpoint:

- `POST /api/assets/{id}/chat`
- later: `POST /api/organizations/{orgId}/chat/context`

### Grounding strategy

The chatbot should answer using:

- organization membership and tenant scope
- asset data from PostgreSQL
- assigned vulnerability data
- CVE details and source metadata
- work order status and priority
- checklist state
- remediation entries
- comments and team notes
- stored NVD metadata
- risk score and risk level
- alert summaries if useful
- sync history if useful

It should not depend on live internet lookups for every answer.

> Security comment: grounded answers are safer than open-ended chat because the model is less likely to hallucinate when it is anchored to local records.

### Chat security rules

- require authentication
- require asset-level authorization
- rate limit chat requests
- limit prompt size
- sanitize user-provided text before prompt construction where needed
- never include secrets, tokens, or hidden internal config in prompts
- treat model output as untrusted text
- render chat output safely in Angular

## Go Services Architecture

Go should remain focused on narrow services that are easy to reason about and easy to secure.

### 1. `alert-service-go`

Purpose:

- detect and process important security events

Examples:

- new critical CVE imported for an asset
- risk score crosses a threshold
- repeated NVD sync failures
- asset becomes newly critical

Possible outputs:

- create alert records
- return event severity
- later send notifications or webhook events

> Design comment: this service should stay focused on alert evaluation and alert generation. It should not become a general-purpose workflow engine.

### 2. `cve-sync-service-go`

Purpose:

- refresh previously imported CVE data on a schedule or on demand

Examples:

- sync modified CVE metadata
- refresh CVSS values
- detect changed NVD records
- update asset vulnerability freshness state

> Implementation note: this service does not replace the main Go backend as the owner of NVD import business rules. It supports refresh and synchronization work. The main Go backend still owns the core application workflow and persistence decisions.

### Why not move NVD or chatbot into Go

Those responsibilities depend heavily on:

- authorization
- DTO control
- asset ownership rules
- database writes
- trust-boundary decisions
- prompt safety

That makes the main Go backend the better home for them.

### Service-to-service authentication

Future internal Go services should not accept privileged work just because a request reaches their port. When the main backend calls `alert-service-go`, `cve-sync-service-go`, or another focused service, the call should include a server-side service credential.

Recommended first version:

- keep the service credential in environment-backed backend and service configuration
- send it only from the Go backend to the internal service
- reject missing or invalid service credentials with a generic `401` or `403`
- expose `/health` for basic liveness checks without secrets
- expose `/ready` for dependency readiness checks if the service has dependencies

This is the practical handshake for this project: the backend proves it is an authorized internal caller before the service performs privileged work. More complex options such as mutual TLS can be considered later during hardening.

## Data Design

### Main tables

- `organizations`
- `users`
- `assets`
- `vulnerabilities`
- `asset_vulnerabilities`
- `work_orders`
- `work_order_checklist_items`
- `vulnerability_exceptions`
- `remediation_entries`
- `comments`
- `alerts`
- optional `chat_sessions`
- optional `chat_messages`
- optional sync history tables

### Relationships

- one organization can own many users, assets, vulnerabilities, work orders, comments, and remediation records
- a user can belong to one or more organizations depending on the tenancy model
- asset CRUD is scoped by organization membership and asset ownership rules
- one asset can have many vulnerabilities
- one vulnerability can be assigned to many assets
- one work order belongs to one organization and usually one asset and one vulnerability context
- one work order can have many checklist items, comments, timeline entries, and exception records
- one asset can produce many alerts
- one asset can have many chat messages if chat persistence is enabled

### Asset shape

An asset should include fields such as:

- `id`
- `organizationId`
- `userId` stored internally as `assets.user_id`
- `name`
- `type`
- `vendor`
- `product`
- `version`
- `ipAddress`
- `operatingSystem`
- `owner`
- `criticality`
- `riskScore`
- `riskLevel`
- `cpeName`
- `cpeMatchConfidence`
- `cpeMatchMethod`
- `lastNvdSyncAt`
- `createdAt`
- `updatedAt`

### Vulnerability shape

A vulnerability should include fields such as:

- `id`
- `organizationId`
- `cveId`
- `title`
- `severity`
- `description`
- `status`
- `source`
- `sourceReference`
- `cvssScore`
- `cvssVector`
- `publishedAt`
- `lastModifiedAt`
- `nvdStatus`
- `createdAt`
- `updatedAt`

### Organization shape

An organization should include fields such as:

- `id`
- `name`
- `slug`
- `status`
- `createdAt`
- `updatedAt`

### Work order shape

A work order should include fields such as:

- `id`
- `organizationId`
- `assetId`
- `vulnerabilityId`
- `assignedUserId`
- `status`
- `priority`
- `dueDate`
- `remediationPlan`
- `verificationNotes`
- `resolutionSummary`
- `createdAt`
- `updatedAt`

Suggested statuses:

- `open`
- `verifying`
- `remediation_planning`
- `patching`
- `testing`
- `resolved`
- `closed`
- `suppressed`
- `exception_granted`

### Checklist item shape

A checklist item should include fields such as:

- `id`
- `workOrderId`
- `title`
- `description`
- `completed`
- `completedBy`
- `completedAt`
- `sortOrder`
- `notes`
- `createdAt`
- `updatedAt`

### Exception shape

A vulnerability exception should include fields such as:

- `id`
- `organizationId`
- `workOrderId`
- `assetId`
- `vulnerabilityId`
- `status`
- `reason`
- `approvedBy`
- `expirationDate`
- `notes`
- `createdAt`
- `updatedAt`

Suggested exception statuses:

- `active`
- `suppressed`
- `exception_granted`
- `false_positive`
- `risk_accepted`
- `remediated`

### Remediation entry shape

A remediation entry should include fields such as:

- `id`
- `workOrderId`
- `authorId`
- `title`
- `bodyMarkdown`
- `entryType`
- `createdAt`
- `updatedAt`

Suggested entry types:

- `investigation`
- `remediation_plan`
- `patch_notes`
- `test_results`
- `resolution`
- `general_note`

### Comment shape

A comment should include fields such as:

- `id`
- `organizationId`
- `assetId`
- `vulnerabilityId`
- `workOrderId`
- `authorId`
- `body`
- `createdAt`
- `updatedAt`

### Asset-vulnerability assignment shape

The join record can eventually hold useful metadata such as:

- `relevanceConfidence`
- `relevanceReason`
- `matchMethod`
- `importedFromNvdAt`
- `lastVerifiedAt`
- `isManuallyReviewed`

> Design comment: this metadata matters because it separates "this CVE exists in NVD" from "this CVE is relevant to this specific asset."

### Alert shape

An alert can include:

- `id`
- `assetId`
- `type`
- `severity`
- `message`
- `status`
- `source`
- `createdAt`
- `acknowledgedAt`

## API Design

### Authentication endpoints

- `POST /api/auth/register`
- `POST /api/auth/login`

### Implemented asset endpoints

- `GET /api/assets`
- `GET /api/assets/{id}`
- `POST /api/assets`
- `PUT /api/assets/{id}`
- `DELETE /api/assets/{id}`
- `POST /api/assets/{assetId}/vulnerabilities/{vulnerabilityId}`
- `DELETE /api/assets/{assetId}/vulnerabilities/{vulnerabilityId}`

These asset endpoints are currently scoped to the authenticated user. The multi-tenant direction is to scope these by organization membership and enforce that users only reach data in organizations they belong to.

### Implemented vulnerability endpoints

- `GET /api/vulnerabilities`
- `GET /api/vulnerabilities/{id}`
- `POST /api/vulnerabilities`
- `PUT /api/vulnerabilities/{id}`
- `DELETE /api/vulnerabilities/{id}`

### Planned endpoints

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

> Implementation note: the sync endpoint should eventually be admin-protected. Background scheduled sync is safer than exposing too much operational control to normal users.

### Response shape rules

- never return password hashes
- use DTOs for request and response payloads
- keep auth errors generic
- keep validation errors readable but not overly revealing
- return safe `404`, `400`, `401`, and `403` responses where appropriate
- do not leak raw upstream stack traces or provider secrets
- map service errors to HTTP responses in the response layer, not in repositories

## Security Architecture

### Authentication and authorization

- passwords should be hashed with BCrypt
- backend routes should be protected with Gin JWT middleware
- authorization checks must be enforced on the backend for all organization, asset, chat, sync, work order, comment, and alert operations
- organization membership must be checked before any org-scoped data is returned or mutated
- admin-only or elevated routes should be explicitly separated from normal user routes
- admin-only routes should use permission middleware such as `RequireAdmin`
- user roles must come from server-side user records, not from frontend state or client-submitted role fields

### Input validation

- validate request DTOs in the backend service layer
- allowlist structured values like severity, status, criticality, and alert types
- allowlist workflow and exception statuses instead of accepting free-form strings
- validate IDs and ownership assumptions before changing records
- validate all imported upstream data before trusting it

### NVD and external API security

- keep the NVD API key in environment variables
- do not expose the NVD key to Angular
- cache and store imported results locally
- respect rate limits
- handle upstream failures safely and generically

### AI security

- keep model API keys server-side only
- never allow the model to access unrestricted internal data
- constrain prompts to asset-specific context
- defend against prompt injection and data exfiltration attempts
- treat model output as advisory text, not executable truth

### RequestFilter

The lightweight `RequestFilter` middleware is meant to block obviously suspicious input patterns early, such as:

- SQL injection-like strings
- XSS-like strings
- path traversal patterns like `../`

This should be treated as an extra defensive layer, not a replacement for normal validation and safe coding.

### Secrets and config

- keep `.env` out of source control
- treat `.env` as local development only
- avoid hardcoded JWT secrets, DB passwords, and service URLs
- prefer environment-based configuration for backend services
- avoid baking secrets into container images

### Logging rules

- log blocked suspicious requests carefully
- log authentication and security-relevant events when useful
- never log plaintext passwords
- never log full JWTs or secret values
- avoid leaking stack traces to API consumers
- keep AI prompt and response logging minimal and sanitized if stored

## Docker and Environment Design

### Compose responsibilities

`docker-compose.yml` should eventually orchestrate:

- PostgreSQL
- Go Gin/GORM backend
- Angular frontend
- `alert-service-go`
- `cve-sync-service-go`

### Networking rules

- Angular should call the Go backend
- the Go backend should call PostgreSQL and Go services by service name inside Docker
- internal Go service calls should use service-to-service authentication for privileged actions
- the Go backend should call external APIs over outbound network access
- Angular should never call NVD directly
- Angular should never call AI providers directly
- Angular should never call Go services directly

### Environment variables

Typical values should include:

- database host
- database port
- database name
- database username
- database password
- JWT secret
- JWT expiration
- NVD API key
- AI provider API key
- alert service base URL
- cve sync service base URL
- frontend API base URL

## Implementation Order

The intended build order is:

1. finish the current frontend auth and basic styling work
2. add organization and membership-aware backend foundations
3. scope asset, vulnerability, and assignment access by organization
4. add work order, checklist, exception, remediation entry, and comment models
5. expose org-scoped workflow endpoints
6. expand the asset model for product-aware matching
7. add Go backend NVD integration
8. add manual `Find Vulnerabilities from NVD`
9. store imported CVEs locally and assign them to assets
10. add risk scoring after the base app is stable
11. add AI-assisted CPE ranking and CVE relevance review
12. add `alert-service-go`
13. add `cve-sync-service-go`
14. add the chatbot context layer
15. add the asset-scoped and org-scoped chatbot experiences
16. add dashboard support for alerts, sync, workflow, and risk trends
17. do Docker integration across all services
18. do hardening, access review, and end-to-end testing

`Roadmap.md` should remain the canonical checklist for what is done and what is next.

## Current Status Summary

Based on the current project plan, the backend foundation is focused on the main Go API, PostgreSQL persistence, authentication, assets, vulnerabilities, and asset-vulnerability assignment. The next major implementation area is still the Angular frontend, with risk scoring deferred until the base app is stable.

## Future Improvements

Possible later improvements include:

- role-based access control
- a dedicated migration tool if GORM AutoMigrate becomes too limited or schema versioning needs more control
- stronger audit logging
- pagination and server-side filtering
- remediation prioritization
- exposure scoring
- better dashboard analytics
- secure secret management beyond local `.env`
- automated tests across frontend, backend, and Go services
- container hardening and production deployment guidance
