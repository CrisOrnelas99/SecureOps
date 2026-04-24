# SecureOps Lite

Cyber Asset Risk Tracker for practicing secure full-stack engineering.

## Project Summary

SecureOps Lite is a cybersecurity-themed full-stack application for tracking assets, assigning vulnerabilities, calculating risk, and viewing security status from a simple dashboard.

This project is intentionally internship-sized:

- Big enough to demonstrate real backend, frontend, database, container, and security skills
- Small enough to finish, explain, demo, and defend in an interview

## Why This Project Matters

This project gives you hands-on practice with:

- Angular
- Spring Boot
- PostgreSQL
- Docker
- Go microservices
- Secure API design
- Basic WAF and backend request filtering concepts

It is not meant to be an enterprise SIEM or vulnerability scanner. It is a realistic learning project that shows you understand secure backend systems and how services connect.

## Core Idea

Users log in and manage cybersecurity assets such as:

- Servers
- Workstations
- Firewalls
- Routers
- Cameras
- Databases
- Cloud services

Each asset can have vulnerabilities assigned to it. The system calculates a risk score based on:

- Vulnerability severity
- Asset criticality

Example asset:

```text
Asset Name: Firewall-01
Asset Type: Firewall
IP Address: 192.168.1.1
Operating System: pfSense
Owner: IT Department
Criticality: High
Assigned Vulnerabilities: 4
Risk Score: 82
Risk Level: Critical
```

## Planned Architecture

```text
Angular frontend
    |
    v
Spring Boot API
    |
    +--> PostgreSQL
    |
    +--> Go risk-scoring service
```

Optional later addition:

- Lightweight WAF-style request filter or gateway

## Main Features

### 1. Authentication

Users should be able to:

- Register
- Log in
- Log out
- Access protected dashboard pages

Security expectations:

- JWT-based authentication
- Password hashing
- Protected backend endpoints
- Protected Angular routes

Initial scope:

- One authenticated user type is enough

Later optional roles:

- Admin
- Analyst
- Viewer

### 2. Dashboard

The dashboard gives a quick security overview.

Example dashboard widgets:

- Total assets
- Total vulnerabilities
- High-risk assets
- Critical vulnerabilities
- Average risk score
- Recently updated assets

Example:

```text
Assets: 12
Open Vulnerabilities: 28
Critical Assets: 3
Average Risk Score: 64
```

### 3. Asset Inventory

Each asset should include:

- Asset ID
- Asset name
- Asset type
- IP address
- Operating system
- Owner
- Criticality
- Risk score
- Risk level
- Created date
- Updated date

Allowed criticality values:

- Low
- Medium
- High

Allowed risk levels:

- Low
- Medium
- High
- Critical

Users should be able to:

- Create assets
- View all assets
- View one asset
- Update assets
- Delete assets
- Search and filter assets

### 4. Vulnerability Tracker

Each vulnerability should include:

- Vulnerability ID
- CVE ID
- Title
- Severity
- Description
- Status
- Created date
- Updated date

Allowed severity values:

- Low
- Medium
- High
- Critical

Allowed status values:

- Open
- Fixed

Users should be able to:

- Create vulnerabilities
- View all vulnerabilities
- View one vulnerability
- Update vulnerabilities
- Delete vulnerabilities
- Filter by severity
- Filter by status

The entries can be manual or mock data. Real scanning is not required for this project.

### 5. Asset-to-Vulnerability Assignment

The system should support many-to-many relationships:

- One asset can have many vulnerabilities
- One vulnerability can affect many assets

Example:

```text
Server-01
- CVE-2024-12345 | High | Open
- CVE-2024-77777 | Critical | Open
- CVE-2024-88888 | Medium | Fixed
```

This is a good practice area for relational database design in PostgreSQL.

### 6. Asset Details Page

Each asset details page should show:

- Asset name
- Asset type
- IP address
- Operating system
- Owner
- Criticality
- Current risk score
- Current risk level
- Assigned vulnerabilities

Actions:

- Assign vulnerability
- Remove vulnerability
- Calculate risk

This page connects the frontend, backend, database, and Go service.

## Risk Scoring Service

The project includes a separate Go service for risk scoring.

Responsibility:

- Accept summarized vulnerability data
- Return a risk score and risk level

Spring Boot flow:

1. Angular user clicks `Calculate Risk`
2. Spring Boot receives the request
3. Spring Boot counts vulnerabilities for the asset
4. Spring Boot sends the data to the Go service
5. Go calculates the result
6. Spring Boot stores the score in PostgreSQL
7. Angular displays the updated result

Example request:

```json
{
  "assetId": 5,
  "criticality": "High",
  "criticalVulnerabilities": 1,
  "highVulnerabilities": 2,
  "mediumVulnerabilities": 3,
  "lowVulnerabilities": 1
}
```

Example response:

```json
{
  "assetId": 5,
  "riskScore": 78,
  "riskLevel": "High"
}
```

Risk formula:

- Critical vulnerabilities x 25
- High vulnerabilities x 15
- Medium vulnerabilities x 8
- Low vulnerabilities x 3

Criticality bonus:

- Low = +0
- Medium = +10
- High = +20

Maximum score:

- 100

Risk levels:

- 0-25 = Low
- 26-50 = Medium
- 51-75 = High
- 76-100 = Critical

## Basic WAF Concept

The project can include a lightweight WAF-style security layer.

The goal is not a full enterprise WAF. The goal is to demonstrate secure backend request filtering and logging.

Simple implementation options:

- Spring Boot security filter
- Separate lightweight Go gateway service

Recommended starting point:

- Spring Boot request filter

Responsibilities:

- Inspect incoming API requests
- Block obviously suspicious patterns
- Log blocked attempts

Examples of patterns to detect:

- SQL injection-like payloads
- XSS-like payloads
- Command injection characters
- Path traversal attempts
- Unusually long request values
- Invalid or missing content types
- Excessive repeated requests from one IP

Examples:

```text
' OR '1'='1
DROP TABLE
<script>
../
; rm -rf
UNION SELECT
```

Example block response:

```json
{
  "error": "Request blocked by WAF",
  "reason": "Suspicious input detected"
}
```

Blocked request logging should capture:

- Timestamp
- IP address
- HTTP method
- Endpoint
- Reason blocked
- Request path

Initial implementation can log to the backend console. Later, logs can be stored in a `waf_events` table.

## Data Model

Main tables:

- `users`
- `assets`
- `vulnerabilities`
- `asset_vulnerabilities`

Optional table:

- `waf_events`

Relationship goals:

- One user can create many assets
- One asset can have many vulnerabilities
- One vulnerability can belong to many assets

## API Design

### Authentication

- `POST /api/auth/register`
- `POST /api/auth/login`

### Assets

- `GET /api/assets`
- `GET /api/assets/{id}`
- `POST /api/assets`
- `PUT /api/assets/{id}`
- `DELETE /api/assets/{id}`
- `POST /api/assets/{id}/calculate-risk`
- `POST /api/assets/{assetId}/vulnerabilities/{vulnerabilityId}`
- `DELETE /api/assets/{assetId}/vulnerabilities/{vulnerabilityId}`

### Vulnerabilities

- `GET /api/vulnerabilities`
- `GET /api/vulnerabilities/{id}`
- `POST /api/vulnerabilities`
- `PUT /api/vulnerabilities/{id}`
- `DELETE /api/vulnerabilities/{id}`

### Optional WAF Events

- `GET /api/waf-events`
- `GET /api/waf-events/{id}`

## Technology Roles

### Angular

Angular is the frontend and should include:

- Login page
- Register page
- Dashboard page
- Assets page
- Asset details page
- Vulnerabilities page

Main skills practiced:

- Components
- Routing
- Forms
- Validation
- HTTP requests
- Services
- Route guards
- Token handling
- Tables and filtering

### Spring Boot

Spring Boot is the main backend API and should handle:

- Authentication
- JWT validation
- Asset CRUD
- Vulnerability CRUD
- Asset-vulnerability assignment
- Go service integration
- PostgreSQL persistence
- Basic WAF request filtering
- Input validation
- Error handling

Recommended backend structure:

- Controllers
- Services
- Repositories
- Entities
- DTOs
- Security classes
- Configuration classes
- Exception handlers

### PostgreSQL

PostgreSQL stores application data and gives practice with:

- Tables
- Primary keys
- Foreign keys
- Many-to-many relationships
- Indexes
- SQL queries
- Migrations

### Go Service

The Go service should expose:

- `POST /calculate-risk`

Its role is intentionally narrow:

- Parse JSON
- Validate input
- Calculate risk
- Return JSON

### Docker

Use Docker Compose to run the project locally.

Containers:

- Angular frontend
- Spring Boot backend
- PostgreSQL database
- Go risk service

Optional:

- pgAdmin

Planned structure:

```text
secureops-lite/
|-- frontend-angular/
|-- backend-springboot/
|-- risk-service-go/
|-- docker-compose.yml
|-- .env
`-- README.md
```

Run target:

```bash
docker compose up --build
```

## Security Goals

This project should demonstrate these security fundamentals:

- Password hashing
- JWT authentication
- Protected backend routes
- Protected frontend routes
- Input validation
- Basic WAF filtering
- Blocked request logging
- Safe error messages
- CORS configuration
- Environment variables for secrets
- Database constraints
- Parameterized data access through JPA
- Basic rate-limiting concepts

Secure defaults:

- Never store plain-text passwords
- Never hardcode JWT secrets
- Never trust frontend validation alone
- Validate on the backend
- Use DTOs instead of exposing entities directly
- Use environment variables for credentials
- Return safe error messages instead of stack traces
- Prefer Docker internal networking over exposing every service publicly

## Resume Value

If built well, this project can demonstrate:

- Secure full-stack development
- API and database design
- Service-to-service communication
- Dockerized local environments
- Authentication and authorization fundamentals
- Defensive backend engineering
- Cybersecurity-focused system design

## Current State

Right now this repository is a planning and architecture foundation.

The next build stages are likely:

1. Scaffold the project structure
2. Stand up PostgreSQL and Spring Boot
3. Add authentication and secured API endpoints
4. Add Angular pages and API integration
5. Add Go risk service
6. Add WAF filtering and security logging

## Notes

Keep the project focused. A secure, clean, explainable implementation is more valuable than adding too many features.
