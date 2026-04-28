# SecureOps Lite

SecureOps Lite is a security-focused full-stack application for tracking assets, managing vulnerabilities, and calculating risk across an environment.

The project is currently being built. This repository represents the direction of the application, the architecture behind it, and the feature set it is intended to support as development continues.

## Project Overview

The goal is to bring a few related security workflows into one system:

- Track assets such as servers, workstations, firewalls, routers, databases, and cloud services
- Record vulnerabilities and their severity
- Assign vulnerabilities to assets
- Calculate a risk score based on vulnerability data and asset criticality
- Show the overall security picture through a dashboard

This is not meant to be a full enterprise SIEM or scanner. The focus is a clean, practical application that shows how asset data, vulnerability data, authentication, and service-to-service communication fit together.

## Core Idea

Each asset in the system carries operational details and a security context. Vulnerabilities can be linked to assets, and those relationships feed into a risk score that helps identify what needs attention first.

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

The backend handles authentication, validation, API logic, and persistence. PostgreSQL stores the application data. The Go service is intended to handle isolated risk calculations, and the Angular frontend provides the user interface.

## Main Features

### Authentication

The application is designed to support:

- User registration
- User login
- JWT-based authentication
- Protected backend endpoints
- Protected frontend routes

Initial access can stay simple, with room to expand later into separate roles such as admin, analyst, or viewer.

### Dashboard

The dashboard is intended to give a quick view of the environment, including:

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

### Asset Inventory

Each asset is intended to include:

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

The system is being built to support:

- Creating assets
- Viewing all assets
- Viewing a single asset
- Updating assets
- Deleting assets
- Searching and filtering assets

### Vulnerability Tracking

Each vulnerability entry is intended to include:

- Vulnerability ID
- CVE ID
- Title
- Severity
- Description
- Status
- Created date
- Updated date

The application is meant to support:

- Creating vulnerabilities
- Viewing all vulnerabilities
- Viewing a single vulnerability
- Updating vulnerabilities
- Deleting vulnerabilities
- Filtering by severity
- Filtering by status

The data can be entered manually or seeded with mock records. The project is focused on managing and connecting the data, not on performing live scanning.

### Asset-to-Vulnerability Assignment

The data model is built around a many-to-many relationship:

- One asset can have many vulnerabilities
- One vulnerability can affect many assets

Example:

```text
Server-01
- CVE-2024-12345 | High | Open
- CVE-2024-77777 | Critical | Open
- CVE-2024-88888 | Medium | Fixed
```

### Asset Details View

An asset details page is intended to show:

- Asset name
- Asset type
- IP address
- Operating system
- Owner
- Criticality
- Current risk score
- Current risk level
- Assigned vulnerabilities

It should also support actions such as assigning vulnerabilities, removing vulnerabilities, and recalculating risk.

## Risk Scoring Service

The project includes a separate Go service for risk scoring.

Its responsibility is straightforward:

- Accept summarized vulnerability data for an asset
- Return a risk score and risk level

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

## Security Direction

The application is being built with a few security fundamentals in mind:

- Password hashing
- JWT authentication
- Protected routes
- Input validation
- Safe error handling
- Environment-based configuration
- Database constraints
- Parameterized data access

There is also room for a lightweight request-filtering layer to catch obviously suspicious input and log blocked requests, but the main priority is getting the core application structure right first.

## Data Model

The main tables are intended to be:

- `users`
- `assets`
- `vulnerabilities`
- `asset_vulnerabilities`

An optional later addition:

- `waf_events`

Relationship goals:

- One user can create many assets
- One asset can have many vulnerabilities
- One vulnerability can belong to many assets

## API Design

The API is being structured around these core routes.

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

## Technology Roles

### Angular

Angular is intended to provide:

- Login page
- Register page
- Dashboard page
- Assets page
- Asset details page
- Vulnerabilities page

### Spring Boot

Spring Boot is the main backend API and is intended to handle:

- Authentication
- Spring Security configuration
- JWT validation
- Asset CRUD
- Vulnerability CRUD
- Asset-vulnerability assignment
- Go service integration
- PostgreSQL persistence
- Input validation
- Error handling

### PostgreSQL

PostgreSQL stores the application data and supports:

- Relational data modeling
- Foreign keys
- Many-to-many relationships
- Querying and filtering

### Go Service

The Go service is intended to expose:

- `POST /calculate-risk`

Its role is narrow by design:

- Parse JSON
- Validate input
- Calculate risk
- Return JSON

### Docker

Docker Compose is used to run the project locally.

Repository structure:

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

## Current Build Status

This repository is actively being built out toward the architecture and feature set above. Some pieces are already in place, and others are still being added.

At the moment, the codebase includes:

- A Spring Boot backend
- Authentication endpoints
- PostgreSQL configuration through Docker Compose
- Repository structure for the frontend and Go service

## Notes

- The backend currently uses `spring.jpa.hibernate.ddl-auto=update`, which is convenient for development but should be replaced with managed migrations before production use.
- Secrets in `.env` should be treated as local development values, not production credentials.
