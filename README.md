# WITS HRMS

HRMS backend project written in Go, using Fiber for APIs and PostgreSQL for persistence.

## Overview

This codebase is currently centered around the **Assets Management** module.  
The project follows a layered structure:

- `routes` for endpoint registration
- `handler` for HTTP request/response logic
- `service` for business rules
- `repository` for SQL/database work
- `model` for request/response/data structs
- `migration` for database schema

## Tech Stack

- Go `1.26.1`
- Fiber `github.com/gofiber/fiber/v2`
- PostgreSQL
- `pgx` / `pgxpool` for database access
- `godotenv` for local env loading

## Current Focus

The main active implementation in this repo is the **Assets module** under:

- [`internal/assets`](/c:/My%20Stuff/WITS/internal/assets)

## Assets Module Features

Current codebase includes flows for:

- create asset
- list assets with filters and pagination
- get asset by id
- update asset
- delete asset
- assign asset
- get my active assets
- update asset status
- return asset
- get assets by employee id
- assignment acknowledge
- assignment history
- maintenance listing
- maintenance creation
- maintenance update
- admin reports
- CSV asset import

## Current Routes

Registered in [`internal/assets/routes/asset_routes.go`](/c:/My%20Stuff/WITS/internal/assets/routes/asset_routes.go):

- `GET /api/v1/assets`
- `POST /api/v1/assets`
- `GET /api/v1/assets/:id`
- `PUT /api/v1/assets/:id`
- `DELETE /api/v1/assets/:id`
- `POST /api/v1/assets/:id/assign`
- `GET /api/v1/assets/my`
- `PATCH /api/v1/assets/:id/status`
- `POST /api/v1/assets/:id/return`
- `GET /api/v1/assets/employee/:employeeId`
- `GET /api/v1/assets/:id/assignments`
- `GET /api/v1/assets/:id/maintenance`
- `POST /api/v1/assets/:id/maintenance`
- `PATCH /api/v1/assets/:id/maintenance/:maintenanceId`
- `PATCH /api/v1/assets/assignments/:id/acknowledge`
- `GET /api/v1/assets/admin/reports`
- `POST /api/v1/assets/import`

## Database Migrations

Current migration files:

- [`005_asset_inventory.up.sql`](/c:/My%20Stuff/WITS/migration/005_asset_inventory.up.sql)
- [`005_asset_assignments.up.sql`](/c:/My%20Stuff/WITS/migration/005_asset_assignments.up.sql)
- [`005_asset_maintenance.up.sql`](/c:/My%20Stuff/WITS/migration/005_asset_maintenance.up.sql)

## Project Structure

Important folders:

- [`cmd/server`](/c:/My%20Stuff/WITS/cmd/server)
- [`cmd/migrate`](/c:/My%20Stuff/WITS/cmd/migrate)
- [`config`](/c:/My%20Stuff/WITS/config)
- [`internal/assets/handler`](/c:/My%20Stuff/WITS/internal/assets/handler)
- [`internal/assets/service`](/c:/My%20Stuff/WITS/internal/assets/service)
- [`internal/assets/repository`](/c:/My%20Stuff/WITS/internal/assets/repository)
- [`internal/assets/model`](/c:/My%20Stuff/WITS/internal/assets/model)
- [`migration`](/c:/My%20Stuff/WITS/migration)
- [`package/database`](/c:/My%20Stuff/WITS/package/database)
- [`package/middleware`](/c:/My%20Stuff/WITS/package/middleware)

## Run Locally

### 1. Create `.env`

Example:

```env
DATABASE_URL=postgres://postgres:password@localhost:5432/hrms
APP_PORT=3000
```

### 2. Apply migrations

```bash
go run ./cmd/migrate
```

### 3. Start server

```bash
go run ./cmd/server
```

Default app URL:

```txt
http://127.0.0.1:3000
```

## Temporary Ownership / Auth

Until shared JWT auth is integrated, the current codebase uses request-header-based ownership helpers in [`package/middleware/middleware.go`](/c:/My%20Stuff/WITS/package/middleware/middleware.go):

- `X-Employee-ID`
- `X-Role`

Examples:

- `X-Role: HR`
- `X-Employee-ID: <employee-uuid>`

## Intern Work Summary

### E1

Mainly aligned with:

- inventory and assignment base schema
- asset listing/filtering
- base inventory structures
- assignment groundwork

### E2

Mainly aligned with:

- maintenance schema and maintenance flows
- create/update/delete asset related APIs
- asset status and return flows
- assignment acknowledge/history related APIs
- employee asset lookup
- report and import endpoints
- temporary ownership handling

## Notes

- The assets module is functional, but still under active refinement.
- Some logic is spread across `asset_inventory`, `asset_assignments`, and `asset_maintenance` files and may be cleaned up later.
- Current ownership checks are temporary and expected to be replaced by shared JWT-based auth in a later integration step.
