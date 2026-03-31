# Progress Module (Doubts)

This module manages doubt/question tracking and resolution for the bootcamp management platform.

## Overview

The Progress module enables mentees to raise doubts on assigned problems and allows mentors/admins to resolve them. It implements role-based access control, cursor-based pagination, and rate limiting to prevent spam.

## Structure

```
progress/
├── dto.go          # Request/Response structures with validation tags and Swagger examples
├── handler.go      # HTTP handlers with comprehensive Swagger annotations
├── service.go      # Business logic layer with SQLC integration
├── routes.go       # Route registration with authentication middleware
├── helper.go       # Utility functions for pagination, filtering, and formatting
└── README.md       # This file
```

## Data Model

```sql
doubts
├── id (UUID, PK)
├── assignment_problem_id (UUID, FK → assignment_problems)
├── raised_by (UUID, FK → bootcamp_enrollments)
├── message (TEXT)
├── resolved (BOOLEAN, default: false)
├── resolved_by (UUID, FK → bootcamp_enrollments, nullable)
├── resolved_at (TIMESTAMPTZ, nullable)
├── resolution_note (TEXT, nullable)
├── created_at (TIMESTAMPTZ)
└── updated_at (TIMESTAMPTZ)
```

## API Endpoints

### 1. Create Doubt

- **Endpoint**: `POST /v1/doubts`
- **Auth**: Mentee only
- **Rate Limit**: 10 requests per minute per user
- **Description**: Create a doubt for an assignment problem

### 2. List Doubts

- **Endpoint**: `GET /v1/doubts`
- **Auth**: Mentor/Admin
- **Query Params**:
    - `assignmentProblemId` (UUID, optional)
    - `resolved` (boolean, optional)
    - `cursor` (string, optional)
    - `limit` (int, optional, default: 20, max: 100)
- **Description**: List doubts with filtering and cursor-based pagination

### 3. Get My Doubts

- **Endpoint**: `GET /v1/doubts/me`
- **Auth**: Mentee only
- **Query Params**:
    - `resolved` (boolean, optional)
    - `cursor` (string, optional)
    - `limit` (int, optional)
- **Description**: Get all doubts raised by the authenticated mentee

### 4. Get Doubt Details

- **Endpoint**: `GET /v1/doubts/:doubtId`
- **Auth**: Mentee (own doubts only), Mentor/Admin (all doubts)
- **Description**: Retrieve full details of a specific doubt

### 5. Resolve Doubt

- **Endpoint**: `PATCH /v1/doubts/:doubtId/resolve`
- **Auth**: Mentor/Admin only
- **Description**: Mark a doubt as resolved with optional resolution note
- **Note**: Idempotent operation

### 6. Delete Doubt

- **Endpoint**: `DELETE /v1/doubts/:doubtId`
- **Auth**: Mentor/Admin only
- **Description**: Permanently delete a doubt (mentees cannot delete for audit purposes)

## Authorization Rules

- **Mentees**:
    - Can create doubts on their assigned problems
    - Can only view their own doubts
    - Cannot delete doubts (audit trail)
    - Cannot resolve doubts

- **Mentors/Admins**:
    - Can view all doubts in their organization
    - Can resolve doubts with optional notes
    - Can delete doubts
    - Cannot create doubts (they are not solving problems)

## Features

### Cursor-Based Pagination

- Efficient for large datasets
- Returns `nextCursor` and `hasMore` in metadata
- Default limit: 20, max limit: 100

### Rate Limiting

- Doubt creation: 10 requests per minute per user
- Prevents spam and abuse

### Multi-Tenant Isolation

- All queries filtered by organization context
- Cross-organization access prevented
- Enrollment validation ensures proper access control

### Validation

- Message length: minimum 10 characters, maximum 2000 characters
- Assignment problem ID must be valid UUID
- Resolution note: maximum 1000 characters
- All UUIDs validated before database queries

## Implementation Status

**Current Status**: Structure created, handlers and service methods are stubs

**Next Steps**:

1. Implement SQLC queries for doubt operations
2. Implement service layer business logic
3. Implement handler logic with proper error handling
4. Add rate limiting middleware
5. Write comprehensive unit tests
6. Integrate with main router and container

## Requirements Mapping

This module implements the following requirements:

- **Requirement 10**: Doubt Management (10.1-10.10)
- **Requirement 11**: Doubt Resolution (11.1-11.10)
- **Requirement 16**: Module Structure and Code Organization (16.1-16.5)
- **Requirement 17**: Request Validation (17.1-17.10)
- **Requirement 23**: Pagination and Filtering (23.6)
- **Requirement 26**: Rate Limiting and Security (26.1)
- **Requirement 31**: API Documentation with Swagger/OpenAPI (31.1-31.20)

## Testing

Unit tests should cover:

- Doubt creation with validation
- Ownership verification
- Role-based access control
- Cursor-based pagination
- Resolution idempotency
- Multi-tenant isolation
- Rate limiting

## Dependencies

- Echo v5 (web framework)
- SQLC (type-safe SQL queries)
- pgx/v5 (PostgreSQL driver)
- go-playground/validator (input validation)
