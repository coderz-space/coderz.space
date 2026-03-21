



# 1. AUTH API DESIGN: 

| URL                   | Method | Auth            | Description                                                              | Version |
| --------------------- | ------ | --------------- | ------------------------------------------------------------------------ | ------- |
| /auth/signup          | POST   | ❌               | Register a new user using email & password and create a platform account |         |
| /auth/login           | POST   | ❌               | Authenticate user via email & password and issue access + refresh tokens |         |
| /auth/google          | POST   | ❌               | Authenticate or register user using Google OAuth ID token                |         |
| /auth/refresh         | POST   | ❌ (uses cookie) | Generate a new access token using a valid refresh token                  |         |
| /auth/logout          | POST   | ✅               | Invalidate current session by revoking refresh token                     |         |
| /auth/me              | GET    | ✅               | Get currently authenticated user profile                                 |         |
| /auth/forgot-password | POST   | ❌               | Initiate password reset flow by sending reset link to email              |         |
| /auth/reset-password  | POST   | ❌               | Reset password using secure reset token                                  |         |
| /auth/verify-email    | POST   | ❌               | Verify user email using verification token                               | v2      |

---
## 1. SIGNUP : 

#### URL

`POST /auth/signup`

#### Purpose

Create a new user account using email/password.

#### Header : 

```http
Content-Type: application/json
```
#### Request Body:

```json
{
  "name": "Suraj",
  "email": "suraj@example.com",
  "password": "StrongPassword123"
}
```
#### Validation Rule : 

- name: 2–100 chars
- email: valid format, unique
- password:
	- min 8 chars
	- at least 1 letter + 1 number

#### Response Body :

```json
{
  "success": true,
  "data": {
    "user": {
      "id": "user_123",
      "name": "Suraj",
      "email": "suraj@example.com",
      "emailVerified": false
    }
  }
}
```
#### Status : `201 Created`



## 2. LOGIN :

#### URL : 
`POST /auth/login`

#### Purpose
Authenticate user and issue tokens.

#### Request Body
```json
{
  "email": "suraj@example.com",
  "password": "StrongPassword123"
}
```

#### Success Response : 
```json
{
  "success": true,
  "data": {
    "accessToken": "jwt_token",
	"refreshToken": "refresh_token",
    "user": {
      "id": "user_123",
      "name": "Suraj",
      "email": "suraj@example.com"
    }
  }
}
```

#### Access Token Cookie :
```http
Set-Cookie: access_token=<JWT>;
Path=/;
HttpOnly;
Secure;
SameSite=Strict;
Max-Age=900;
```
#### Token Behaviour : 
- Access token → response body and HttpOnly cookie
- Refresh token → HttpOnly cookie


#### Error : 
| Status | Code                |
| ------ | ------------------- |
| 400    | VALIDATION_ERROR    |
| 401    | INVALID_CREDENTIALS |
| 403    | EMAIL_NOT_VERIFIED  |


## 3. GOOGLE AUTH : V2

#### URL
`POST /auth/google`

#### Purpose
Login or signup using Google OAuth.

#### Request Body 
```json
{
  "idToken": "google_id_token"
}
```

#### Behavior
- Verify token with Google
- If user exists → login
- Else → create user

#### Success Response
```json
{
  "success": true,
  "data": {
    "accessToken": "jwt_token",
    "user": {
      "id": "user_123",
      "name": "Suraj",
      "email": "suraj@gmail.com"
    }
  }
}
```


#### Error : 
| Status | Code                 |
| ------ | -------------------- |
| 401    | INVALID_GOOGLE_TOKEN |


## 4. REFRESH TOKEN :

#### URL
`POST /auth/refresh`

#### Purpose
Generate new access token.

#### Headers
Cookie required:
```HTTP
refresh_token=xyz
```

#### Success Response
```json
{
  "success": true,
  "data": {
    "accessToken": "new_jwt_token"
  }
}
```

#### Behavior
- Validate refresh token (DB lookup)
- Rotate refresh token (recommended)


#### Errors :
| Status | Code                  |
| ------ | --------------------- |
| 401    | INVALID_REFRESH_TOKEN |
| 401    | EXPIRED_REFRESH_TOKEN |


## 5. LOGOUT : 

#### URL
`POST /auth/logout`

#### Purpose
Invalidate session.

#### Headers
`Authorization: Bearer <token>`

#### Behavior
- Delete refresh token from DB
- Clear cookie


#### Success Response : 
```json
{
  "success": true,
  "data": {}
}
```



## 6. GET CURRENT USER 

#### URL
`GET /auth/me`

#### Purpose
Fetch authenticated user profile.

#### Success Response
```JSON
{
  "success": true,
  "data": {
    "id": "user_123",
    "name": "Suraj",
    "email": "suraj@example.com",
    "avatarUrl": null
  }
}
```


## 7. FORGOT PASSWORD

#### URL
`POST /auth/forgot-password`

#### Purpose
Send password reset link.

#### Behavior
- Always return success (avoid email enumeration)

#### Request
```JSON
{
  "email": "suraj@example.com"
}
```

#### Response
```JSON
{
  "success": true,
  "data": {}
}
```


## 8. RESET PASSWORD : 

#### URL
`POST /auth/reset-password`

#### Request 
```json
{
  "token": "reset_token",
  "newPassword": "NewPassword123"
}
```

#### Success Response :
```json
{
  "success": true,
  "data": {}
}
```


## 9. VERIFY EMAIL : V2


---

# 2. Organization Module

| NO  | URL                             | Method | Auth            | Description                   |     |
| --- | ------------------------------- | ------ | --------------- | ----------------------------- | --- |
| 1   | `/orgs`                         | POST   | USER            | Create org (pending approval) |     |
| 2   | `/orgs`                         | GET    | USER            | List orgs user belongs to     |     |
| 3   | `/orgs/{slug}`                  | GET    | USER            | Get org details               |     |
| 4   | `/orgs/{slug}`                  | PATCH  | ADMIN           | Update org details            |     |
| 5   | `/orgs/{slug}/members`          | GET    | ADMIN / MENTORS | List members                  |     |
| 6   | `/orgs/{slug}/members`          | POST   | ADMIN / MENTORS | Add member                    |     |
| 7   | `/orgs/{slug}/members/{userId}` | PATCH  | ADMIN           | Update member role            |     |
| 8   | `/orgs/{slug}/members/{userId}` | DELETE | ADMIN/MENTOR    | Remove member                 |     |
| 9   | `/orgs/{slug}/join`             | POST   | USER            | Request/join org              |     |
| 10  | `/orgs/{slug}/leave`            | POST   | USER            | Leave org                     |     |
| 11  | `/orgs/pending`                 | GET    | SUPER_ADMIN     | View pending orgs             |     |
| 12  | `/orgs/{id}/approve`            | POST   | SUPER_ADMIN     | Approve org                   |     |
| 13  | `/orgs/{id}/suspend`            | POST   | SUPER_ADMIN     | Suspend org                   |     |

## 1. CREATE ORGANIZATION : `POST /orgs`

#### **Purpose:** 
- Use to Create organization.
- pending approval from SUPER_ADMIN

#### **Access Control**
- Any authenticated USER

#### Header : 
```http
Authorization: Bearer <token>
```

#### Request Body : 
```json
{
  "name": "Algo University",
  "slug": "algo-university",
  "description": "DSA focused bootcamp"
}
```

#### Validation  :
```json
{
  "name": "Algo University",
  "slug": "algo-university",
  "description": "DSA focused bootcamp"
}
```


#### Success Response : 
```json
{
  "id": "org_uuid",
  "status": "PENDING_APPROVAL"
}
```



#### Error 
|Status|Code|Condition|
|---|---|---|
|409|ORG_SLUG_EXISTS|slug taken|
|400|INVALID_INPUT|bad slug|


### 2. GET ORGANIZATION :  GET `/orgs`

#### Purpose : 
- List all the orgs that user belongs to.

#### Access Control :
- USER

## 3. GET ORGANIZATION :  `GET /orgs/{slug}`

#### Purpose
- Get details of the organization.

#### Access Control
- Only members of the organization : MENTEE, MENTOR, ADMIN, SUPER_ADMIN




# 3. BOOTCAMP : 

| #   | URL                                                          | Method | Auth                                   | Description                                                               |
| --- | ------------------------------------------------------------ | ------ | -------------------------------------- | ------------------------------------------------------------------------- |
| 1   | `/orgs/{org_id}/b`                                           | POST   | super_admin, org_admin                 | Create a bootcamp inside an organization                                  |
| 2   | `/orgs/{org_id}/b`                                           | GET    | super_admin, org_admin, mentor, mentee | List bootcamps for an organization                                        |
| 3   | `/orgs/{org_id}/b/{bootcamp_id}`                             | GET    | super_admin, org_admin, mentor, mentee | Get bootcamp details                                                      |
| 4   | `/orgs/{org_id}/b/{bootcamp_id}`                             | PATCH  | org_admin                              | Update bootcamp metadata or active state                                  |
| 5   | `/orgs/{org_id}/b/{bootcamp_id}/enrollments`                 | GET    | super_admin, org_admin, mentor,mentee  | List bootcamp participants                                                |
| 6   | `/orgs/{org_id}/b/{bootcamp_id}/enrollments`                 | POST   | org_admin                              | Enroll an org member into the bootcamp                                    |
| 7   | `/orgs/{org_id}/b/{bootcamp_id}/enrollments/{enrollment_id}` | PATCH  | org_admin                              | Change participant role or status.<br>mentee -> mentor<br>mentor -> admin |
| 8   | `/orgs/{org_id}/b/{bootcamp_id}/enrollments/{enrollment_id}` | DELETE | org_admin                              | Remove a participant from the bootcamp                                    |
## 1. GLOBAL BEHAVIOR
```mermaid
erDiagram
    ORGANIZATION ||--o{ BOOTCAMP : owns
    ORGANIZATION ||--o{ ORGANIZATION_MEMBER : has
    BOOTCAMP ||--o{ BOOTCAMP_ENROLLMENT : contains
    ORGANIZATION_MEMBER ||--o{ BOOTCAMP_ENROLLMENT : participates_in

    ORGANIZATION {
        uuid id PK
        string name
    }

    BOOTCAMP {
        uuid id PK
        uuid organization_id FK
        string name
        text description
        date start_date
        date end_date
        boolean is_active
    }

    ORGANIZATION_MEMBER {
        uuid id PK
        uuid organization_id FK
        uuid user_id FK
    }

    BOOTCAMP_ENROLLMENT {
        uuid id PK
        uuid bootcamp_id FK
        uuid organization_member_id FK
        enum role
        timestamp enrolled_at
    }
```
### Role hierarchy

- **super_admin**: can act across all organizations.
- **org_admin**: manages bootcamps and enrollments only inside their organization.
- **mentor / mentee**: can read bootcamp data only when they are part of that organization or bootcamp.

### Important invariants

- A bootcamp always belongs to exactly one organization.
- A bootcamp enrollment is valid only if the user is already an **organization member** of the same organization.
- `(bootcamp_id, organization_member_id)` must be unique.
- `bootcamp_enrollment.role` is **bootcamp-scoped** and does **not** have to match org-level role.
- `start_date <= end_date` when both are present.
- `is_active=false` means the bootcamp is closed for new participation, but historical data stays intact. This is the safer production choice.


## 1. CREATE BOOTCAMP : `POST /orgs/{org_id}/b`


**Purpose:** Create a new bootcamp under the given organization.
**Access Control:** `super_admin`, `org_admin`
- `org_admin` can only create bootcamps in their own organization.
- `super_admin` can create anywhere.

**Headers**
- `Authorization: Bearer <token>`
- `Idempotency-Key: <optional>` for retry-safe creation

**Path Params**
- `org_id` — organization UUID

#### Request body :
```json
{
  "name": "8 Week DSA Bootcamp",
  "description": "Core DSA cohort for new learners",
  "start_date": "2026-04-01",
  "end_date": "2026-05-27",
  "is_active": true
}
```
**Validation Rules**

- `name` required, trimmed, 3–120 chars
- `description` optional, max e.g. 2000 chars
- `start_date` and `end_date` optional
- if both dates exist, `start_date <= end_date`
- `is_active` defaults to `true` if omitted

**Success Response**

- `201 Created`

```json
{
  "id": "bootcamp_uuid",
  "organization_id": "org_uuid",
  "name": "8 Week DSA Bootcamp",
  "description": "Core DSA cohort for new learners",
  "start_date": "2026-04-01",
  "end_date": "2026-05-27",
  "is_active": true,
  "created_at": "2026-03-21T10:00:00Z",
  "updated_at": "2026-03-21T10:00:00Z"
}
```

**Error :**

|Status|Code|Condition|
|---|---|---|
|400|invalid_body|Malformed JSON or validation failure|
|401|unauthorized|Missing/invalid token|
|403|forbidden|Caller cannot manage this org|
|404|org_not_found|Organization does not exist|
|409|duplicate_bootcamp_name|Same org already has conflicting bootcamp name, if you enforce that rule|


## 2. LIST BOOTCAMPS : `GET /orgs/{org_id}/b`


**Purpose:** 
- List bootcamps within an organization.

**Access Control:**
- `super_admin`, `org_admin`, `mentor`, `mentee`

**Headers**
- `Authorization: Bearer <token>`

**Path Params**
- `org_id` — organization UUID

**Query Params**
- `q` — search by name
- `is_active` — `true|false`
- `page`, `limit`

**Success Response**
- `200 OK`

```json
{
  "items": [
    {
      "id": "bootcamp_uuid",
      "name": "8 Week DSA Bootcamp",
      "is_active": true,
      "start_date": "2026-04-01",
      "end_date": "2026-05-27"
    }
  ],
  "page": 1,
  "limit": 20,
  "total": 1
}
```


#### Error :
| Status | Code          | Condition                             |
| ------ | ------------- | ------------------------------------- |
| 401    | unauthorized  | Missing/invalid token                 |
| 403    | forbidden     | Caller is not allowed to see this org |
| 404    | org_not_found | Organization does not exist           |

## 3. GET BOOTCAMP DETAILS :  `GET /orgs/{org_id}/b/{bootcamp_id}`


**Purpose:** Fetch one bootcamp with full metadata.

**Access Control:** `super_admin`, `org_admin`, `mentor`, `mentee`
- The caller must belong to the org, or be super admin.
- If you want stricter access, require the caller to be an org member.

**Headers**
- `Authorization: Bearer <token>`

**Path Params**
- `org_id` — organization UUID
- `bootcamp_id` — bootcamp UUID

**Success Response**
- `200 OK`

```json
{  
  "id": "bootcamp_uuid",  
  "organization_id": "org_uuid",  
  "name": "8 Week DSA Bootcamp",  
  "description": "Core DSA cohort for new learners",  
  "start_date": "2026-04-01",  
  "end_date": "2026-05-27",  
  "is_active": true,  
  "created_at": "2026-03-21T10:00:00Z",  
  "updated_at": "2026-03-21T10:00:00Z"  
}
```

**Errors**

|Status|Code|Condition|
|---|---|---|
|401|unauthorized|Missing/invalid token|
|403|forbidden|Bootcamp does not belong to this org or caller has no access|
|404|bootcamp_not_found|Bootcamp does not exist|


## 4. UPDATE BOOTCAMP : `PATCH /orgs/{org_id}/b/b{bootcamp_id}`

**Purpose:** Update bootcamp metadata or active state.
**Access Control:** `org_admin`

**Headers**
- `Authorization: Bearer <token>`

**Path Params**
- `org_id` — organization UUID
- `bootcamp_id` — bootcamp UUID

**Request Body**
```json
{
  "name": "Backend Engineering Cohort",
  "description": "Cohort for Go and system design",
  "start_date": "2026-04-10",
  "end_date": "2026-06-10",
  "is_active": false
}
```

**Validation Rules**

- All fields optional, but at least one must be present
- `name` if present: 3–120 chars
- `description` if present: max length enforced
- dates must be valid and ordered
- `is_active` toggles availability for new enrollment

**Success Response**

- `200 OK`
```json
{
  "id": "bootcamp_uuid",
  "organization_id": "org_uuid",
  "name": "Backend Engineering Cohort",
  "description": "Cohort for Go and system design",
  "start_date": "2026-04-10",
  "end_date": "2026-06-10",
  "is_active": false,
  "updated_at": "2026-03-21T10:20:00Z"
}
```

#### Error : 
| Status | Code                     | Condition                                      |
| ------ | ------------------------ | ---------------------------------------------- |
| 400    | invalid_body             | No valid fields or bad data                    |
| 401    | unauthorized             | Missing/invalid token                          |
| 403    | forbidden                | Caller cannot manage this org                  |
| 404    | bootcamp_not_found       | Bootcamp does not exist                        |
| 409    | invalid_state_transition | Example: trying to activate with invalid dates |

## 5. LIST ENROLLMENTS : `GET /orgs/{org_id}/b/{bootcamp_id}/enrollments`


**Purpose:** List all participants in a bootcamp.

**Access Control:** `super_admin`, `org_admin`, `mentor` `mentee`
- `mentor` can read participants if they are part of the bootcamp or org.

**Headers**
- `Authorization: Bearer <token>`

**Path Params**
- `org_id` — organization UUID
- `bootcamp_id` — bootcamp UUID

**Query Params**
- `role` — `mentor|mentee`
- `page`, `limit`

**Success Response**
- `200 OK`

```JSON
{  
  "items": [  
    {  
      "id": "enrollment_uuid",  
      "bootcamp_id": "bootcamp_uuid",  
      "organization_member_id": "member_uuid",  
      "role": "mentor",  
      "enrolled_at": "2026-03-21T10:00:00Z"  
    }  
  ],  
  "page": 1,  
  "limit": 20,  
  "total": 1  
}
```

**Errors**

|Status|Code|Condition|
|---|---|---|
|401|unauthorized|Missing/invalid token|
|403|forbidden|Caller cannot see this bootcamp|
|404|bootcamp_not_found|Bootcamp does not exist|


## 6. ENROLL MEMBER IN BOOTCAMP : `POST /orgs/{org_id}/b/{bootcamp_id}/enrollments`

**Purpose:** 
- Enroll an organization member into a bootcamp with a bootcamp-specific role.
**Access Control:**
-  `org_admin`

**Headers**
- `Authorization: Bearer <token>`
- `Idempotency-Key: <optional>`

**Path Params**
- `org_id` — organization UUID
- `bootcamp_id` — bootcamp UUID

**Request Body**
```json
{  
  "organization_member_id": "member_uuid",  
  "role": "mentee"  
}
```

**Validation Rules**
- `organization_member_id` required
- `role` required, enum: `mentor | mentee`
- member must belong to the same org as the bootcamp
- bootcamp must be active unless you explicitly allow late joins
- prevent duplicates via unique `(bootcamp_id, organization_member_id)`

**Success Response**
- `201 Created`
```json
{  
  "id": "enrollment_uuid",  
  "bootcamp_id": "bootcamp_uuid",  
  "organization_member_id": "member_uuid",  
  "role": "mentee",  
  "enrolled_at": "2026-03-21T10:30:00Z"  
}
```

**Errors**

|Status|Code|Condition|
|---|---|---|
|400|invalid_body|Bad request or invalid role|
|401|unauthorized|Missing/invalid token|
|403|forbidden|Caller cannot manage this org|
|404|bootcamp_not_found|Bootcamp does not exist|
|404|org_member_not_found|Member does not exist in this org|
|409|duplicate_enrollment|Same member already enrolled|
|409|cross_org_violation|Member belongs to a different organization|
|409|bootcamp_inactive|Enrollment blocked because bootcamp is closed|

## 7. Update Enrollment

**URL:** `PATCH /orgs/{org_id}/b/{bootcamp_id}/enrollments/{enrollment_id}`

**Purpose:** 
- Change a participant’s bootcamp role or status-related fields.
- Org Admin can use it to promote a member (mentee) to mentor 

**Access Control:** `org_admin`

**Headers**
- `Authorization: Bearer <token>`

**Path Params**
- `org_id` — organization UUID
- `bootcamp_id` — bootcamp UUID
- `enrollment_id` — enrollment UUID

**Request Body**
```json
{  
  "role": "mentor"  
}
```

**Validation Rules**
- `role` optional but if present must be `mentor|mentee`
- enrollment must belong to the given bootcamp
- org must match
- if you later add `is_active` to enrollment, keep the same pattern here

**Success Response**
- `200 OK`

```json
{  
  "id": "enrollment_uuid",  
  "bootcamp_id": "bootcamp_uuid",  
  "organization_member_id": "member_uuid",  
  "role": "mentor",  
  "enrolled_at": "2026-03-21T10:00:00Z"  
}
```

**Errors**

|Status|Code|Condition|
|---|---|---|
|400|invalid_body|Invalid role or empty patch|
|401|unauthorized|Missing/invalid token|
|403|forbidden|Caller cannot manage this org|
|404|enrollment_not_found|Enrollment does not exist|
|409|duplicate_enrollment|If changing role collides with an existing enforced rule|


## 8. Remove Enrollment

**URL:** `DELETE /orgs/{org_id}/b/{bootcamp_id}/enrollments/{enrollment_id}`

**Purpose:** 
- Remove a participant from the bootcamp.

**Access Control:** `org_admin`, `mentor`

**Headers**
- `Authorization: Bearer <token>`

**Path Params**
- `org_id` — organization UUID
- `bootcamp_id` — bootcamp UUID
- `enrollment_id` — enrollment UUID

**Success Response**
- `204 No Content`

**Errors**

| Status | Code                      | Condition                                          |
| ------ | ------------------------- | -------------------------------------------------- |
| 401    | unauthorized              | Missing/invalid token                              |
| 403    | forbidden                 | Caller cannot manage this org                      |
| 404    | enrollment_not_found      | Enrollment does not exist                          |
| 409    | cannot_remove_last_mentor | If your business rule requires at least one mentor |



# 4. PROBLEM CONTENT Module 

| #   | URL                                                            | Method | Auth                               | Description                                   |
| --- | -------------------------------------------------------------- | ------ | ---------------------------------- | --------------------------------------------- |
| 1   | `/orgs/{org_id}/problems`                                      | POST   | admin, mentor                      | Create a new problem in the org question bank |
| 2   | `/orgs/{org_id}/problems`                                      | GET    | super_admin, admin, mentor, mentee | List org problems with filters                |
| 3   | `/orgs/{org_id}/problems/{problem_id}`                         | GET    | super_admin, admin, mentor, mentee | Get problem details                           |
| 4   | `/orgs/{org_id}/problems/{problem_id}`                         | PATCH  | admin, mentor                      | Update a problem                              |
| 5   | `/orgs/{org_id}/problems/{problem_id}`                         | DELETE | admin, mentor                      | Delete a problem                              |
| 6   | `/orgs/{org_id}/tags`                                          | POST   | admin, mentor                      | Create a tag in the org                       |
| 7   | `/orgs/{org_id}/tags`                                          | GET    | super_admin, admin, mentor, mentee | List org tags                                 |
| 8   | `/orgs/{org_id}/tags/{tag_id}`                                 | PATCH  | admin, mentor                      | Rename a tag                                  |
| 9   | `/orgs/{org_id}/tags/{tag_id}`                                 | DELETE | admin, mentor                      | Delete a tag                                  |
| 10  | `/orgs/{org_id}/problems/{problem_id}/tags`                    | POST   | admin, mentor                      | Attach tags to a problem                      |
| 11  | `/orgs/{org_id}/problems/{problem_id}/tags/{tag_id}`           | DELETE | admin, mentor                      | Detach a tag from a problem                   |
| 12  | `/orgs/{org_id}/problems/{problem_id}/resources`               | POST   | admin, mentor                      | Add a learning resource to a problem          |
| 13  | `/orgs/{org_id}/problems/{problem_id}/resources`               | GET    | super_admin, admin, mentor, mentee | List problem resources                        |
| 14  | `/orgs/{org_id}/problems/{problem_id}/resources/{resource_id}` | PATCH  | admin, mentor                      | Update a resource                             |
| 15  | `/orgs/{org_id}/problems/{problem_id}/resources/{resource_id}` | DELETE | admin, mentor                      | Remove a resource                             |


## Entity Relationship : 


```mermaid
erDiagram
    ORGANIZATIONS ||--o{ PROBLEMS : owns
    ORGANIZATIONS ||--o{ TAGS : owns
    ORGANIZATION_MEMBERS ||--o{ PROBLEMS : creates
    PROBLEMS ||--o{ PROBLEM_RESOURCES : has
    PROBLEMS ||--o{ PROBLEM_TAGS : mapped_by
    TAGS ||--o{ PROBLEM_TAGS : mapped_by

    ORGANIZATIONS {
        uuid id PK
        string name
    }

    ORGANIZATION_MEMBERS {
        uuid id PK
        uuid organization_id FK
        uuid user_id FK
        string role
    }

    PROBLEMS {
        uuid id PK
        uuid organization_id FK
        uuid created_by FK
        string title
        text description
        enum difficulty
        string external_link
    }

    TAGS {
        uuid id PK
        uuid organization_id FK
        string name
    }

    PROBLEM_TAGS {
        uuid problem_id FK
        uuid tag_id FK
    }

    PROBLEM_RESOURCES {
        uuid id PK
        uuid problem_id FK
        string title
        string url
    }
```


## 1. CREATE PROBLEM :

- **URL:** `POST /orgs/{org_id}/problems`
- **Purpose:** Create a reusable problem in the org master question bank.
- **Access Control:** `admin`, `mentor`

#### Headers

- `Authorization: Bearer <token>`
- `Content-Type: application/json`
- Cookie auth may also be accepted if your gateway supports it.

#### Path Params

- `org_id` — organization UUID

#### Request Body

```JSON
{
  "title": "Two Sum",
  "description": "Given an array of integers, return indices of the two numbers such that they add up to a target.",
  "difficulty": "easy",
  "external_link": "https://leetcode.com/problems/two-sum"
}
```

### Validation Rules

- `title` required, trimmed, 3–200 chars
- `description` optional
- `difficulty` required, enum: `easy | medium | hard`
- `external_link` optional, must be a valid URL if present
- caller must be a member of the same organization
- `title` uniqueness inside an org is recommended, but not mandatory unless you want strict dedupe
- `created_by` is derived from the authenticated org member, never accepted from client

#### Success Response

- `201 Created`
```json
{
  "id": "prob_uuid",
  "organization_id": "org_uuid",
  "created_by": "org_member_uuid",
  "title": "Two Sum",
  "description": "Given an array of integers, return indices of the two numbers such that they add up to a target.",
  "difficulty": "easy",
  "external_link": "https://leetcode.com/problems/two-sum",
  "created_at": "2026-03-21T10:00:00Z",
  "updated_at": "2026-03-21T10:00:00Z"
}
```

#### Error :
| Status | Code              | Condition                                            |
| ------ | ----------------- | ---------------------------------------------------- |
| 400    | invalid_body      | Bad JSON or validation failed                        |
| 401    | unauthorized      | Missing or invalid token                             |
| 403    | forbidden         | Caller is not allowed to create problems in this org |
| 404    | org_not_found     | Organization does not exist                          |
| 409    | duplicate_problem | Optional if you enforce title uniqueness             |

## 2. LIST PROBLEMS: 

- **URL:** `GET /orgs/{org_id}/problems`
- **Purpose:** List reusable problems in the org question bank.
- **Access Control:** `super_admin`, `admin`, `mentor`, `mentee`

#### Headers
- `Authorization: Bearer <token>`

#### Path Params
- `org_id` — organization UUID

#### Query Params
- `q` — search by title
- `difficulty` — `easy | medium | hard`
- `tag_id` — filter by tag
- `page`
- `limit`
- `sort_by` — `created_at | title | difficulty`
- `order` — `asc | desc`

#### Success Response
- `200 OK`

```json
{  
  "items": [  
    {  
      "id": "prob_uuid",  
      "title": "Two Sum",  
      "difficulty": "easy",  
      "external_link": "https://leetcode.com/problems/two-sum",  
      "tag_count": 3,  
      "resource_count": 2,  
      "created_at": "2026-03-21T10:00:00Z"  
    }  
  ],  
  "page": 1,  
  "limit": 20,  
  "total": 1  
}
```

#### Errors

|Status|Code|Condition|
|---|---|---|
|401|unauthorized|Missing or invalid token|
|403|forbidden|Caller cannot access this organization|
|404|org_not_found|Organization does not exist|

## 3. GET PROBLEM DETAILS :

- **URL:** `GET /orgs/{org_id}/problems/{problem_id}`
- **Purpose:** Fetch one problem with its tags and resources.
- **Access Control:** `super_admin`, `admin`, `mentor`, `mentee`

#### Headers
- `Authorization: Bearer <token>`

#### Path Params
- `org_id` — organization UUID
- `problem_id` — problem UUID

#### Success Response
- `200 OK`

```json
{  
  "id": "prob_uuid",  
  "organization_id": "org_uuid",  
  "created_by": "org_member_uuid",  
  "title": "Two Sum",  
  "description": "Given an array of integers, return indices of the two numbers such that they add up to a target.",  
  "difficulty": "easy",  
  "external_link": "https://leetcode.com/problems/two-sum",  
  "tags": [  
    {  
      "id": "tag_uuid_1",  
      "name": "arrays"  
    },  
    {  
      "id": "tag_uuid_2",  
      "name": "hash-map"  
    }  
  ],  
  "resources": [  
    {  
      "id": "res_uuid_1",  
      "title": "Official Editorial",  
      "url": "https://example.com/editorial"  
    }  
  ],  
  "created_at": "2026-03-21T10:00:00Z",  
  "updated_at": "2026-03-21T10:00:00Z"  
}
```

#### Errors

|Status|Code|Condition|
|---|---|---|
|401|unauthorized|Missing or invalid token|
|403|forbidden|Org access denied|
|404|problem_not_found|Problem does not exist in this org|


## 4) Update Problem

- **URL:** `PATCH /orgs/{org_id}/problems/{problem_id}`
- **Purpose:** Update problem metadata.
- **Access Control:** `admin`, `mentor`

#### Headers
- `Authorization: Bearer <token>`
- `Content-Type: application/json`

#### Path Params
- `org_id` — organization UUID
- `problem_id` — problem UUID

#### Request Body

```json
{  
  "title": "Two Sum - Optimized",  
  "description": "Updated statement",  
  "difficulty": "easy",  
  "external_link": "https://leetcode.com/problems/two-sum"  
}
```

#### Validation Rules

- all fields optional, but at least one field must be present
- `title` length 3–200 if provided
- `difficulty` enum if provided
- URL must be valid if provided
- problem must belong to the org in path
- do not allow changing `organization_id` or `created_by`

#### Success Response

- `200 OK`
```
{  
  "id": "prob_uuid",  
  "organization_id": "org_uuid",  
  "created_by": "org_member_uuid",  
  "title": "Two Sum - Optimized",  
  "description": "Updated statement",  
  "difficulty": "easy",  
  "external_link": "https://leetcode.com/problems/two-sum",  
  "updated_at": "2026-03-21T10:15:00Z"  
}
```
#### Errors

| Status | Code              | Condition                                 |
| ------ | ----------------- | ----------------------------------------- |
| 400    | invalid_body      | No valid fields or bad values             |
| 401    | unauthorized      | Missing or invalid token                  |
| 403    | forbidden         | Caller cannot update problems in this org |
| 404    | problem_not_found | Problem does not exist                    |
| 409    | duplicate_problem | If title uniqueness is enforced           |


## 5) Delete Problem

- **URL:** `DELETE /orgs/{org_id}/problems/{problem_id}`
- **Purpose:** Remove a problem from the org question bank.
- **Access Control:** `admin`, `mentor`

#### Headers
- `Authorization: Bearer <token>`

#### Path Params
- `org_id` — organization UUID
- `problem_id` — problem UUID

#### Success Response
- `204 No Content`

#### Important rule
Do not hard-delete if the problem is already used in assignments/submissions unless your system has a proper archival strategy. In production, soft delete is usually safer.

#### Errors

|Status|Code|Condition|
|---|---|---|
|401|unauthorized|Missing or invalid token|
|403|forbidden|Caller cannot delete problems in this org|
|404|problem_not_found|Problem does not exist|
|409|problem_in_use|Problem is referenced by assignments or submissions|

---

## 6) Create Tag

- **URL:** `POST /orgs/{org_id}/tags`
- **Purpose:** Create a reusable org tag.
- **Access Control:** `admin`, `mentor`

#### Headers
- `Authorization: Bearer <token>`
- `Content-Type: application/json`

#### Path Params

- `org_id` — organization UUID

#### Request Body
```json
{  
  "name": "sliding-window"  
}
```

#### Validation Rules
- `name` required, trimmed, 2–80 chars
- store normalized form consistently, e.g. lowercase with hyphens
- must be unique per organization
- do not accept `organization_id` from client

#### Success Response
- `201 Created`

```json
{  
  "id": "tag_uuid",  
  "organization_id": "org_uuid",  
  "name": "sliding-window",  
  "created_at": "2026-03-21T10:00:00Z"  
}
```

### Errors

|Status|Code|Condition|
|---|---|---|
|400|invalid_body|Bad request or invalid tag name|
|401|unauthorized|Missing or invalid token|
|403|forbidden|Caller cannot create tags|
|404|org_not_found|Organization does not exist|
|409|duplicate_tag|Tag already exists in this org|

---

## 7) List Tags

- **URL:** `GET /orgs/{org_id}/tags`
- **Purpose:** List all tags in the org.
- **Access Control:** `super_admin`, `admin`, `mentor`, `mentee`

#### Headers
- `Authorization: Bearer <token>`

#### Path Params

- `org_id` — organization UUID

#### Query Params

- `q` — search tag name
- `page`
- `limit`

#### Success Response

- `200 OK`

```json
{  
  "items": [  
    {  
      "id": "tag_uuid",  
      "name": "sliding-window"  
    }  
  ],  
  "page": 1,  
  "limit": 20,  
  "total": 1  
}
```

#### Errors

| Status | Code          | Condition                              |
| ------ | ------------- | -------------------------------------- |
| 401    | unauthorized  | Missing or invalid token               |
| 403    | forbidden     | Caller cannot access this organization |
| 404    | org_not_found | Organization does not exist            |


## 8) Update Tag

- **URL:** `PATCH /orgs/{org_id}/tags/{tag_id}`
- **Purpose:** Rename a tag.
- **Access Control:** `admin`, `mentor`

#### Headers
- `Authorization: Bearer <token>`
- `Content-Type: application/json`

#### Path Params
- `org_id` — organization UUID
- `tag_id` — tag UUID

#### Request Body
```json
{  
  "name": "two-pointers"  
}
```

#### Validation Rules
- `name` required in patch body
- normalized uniqueness enforced within org
- tag must belong to the path org

#### Success Response

- `200 OK`

```json
{  
  "id": "tag_uuid",  
  "organization_id": "org_uuid",  
  "name": "two-pointers",  
  "created_at": "2026-03-21T10:00:00Z"  
}
```
#### Errors

|Status|Code|Condition|
|---|---|---|
|400|invalid_body|Missing or invalid name|
|401|unauthorized|Missing or invalid token|
|403|forbidden|Caller cannot update tags|
|404|tag_not_found|Tag does not exist|
|409|duplicate_tag|Another tag with same name already exists|

---

## 9) Delete Tag

- **URL:** `DELETE /orgs/{org_id}/tags/{tag_id}`
- **Purpose:** Remove a tag from the org.
- **Access Control:** `admin`, `mentor`

#### Headers

- `Authorization: Bearer <token>`

#### Path Params

- `org_id` — organization UUID
- `tag_id` — tag UUID

#### Success Response

- `204 No Content`

#### Important rule

If the tag is attached to problems, you need a business decision:

- either cascade remove from `problem_tags`, or
- reject deletion with `409 tag_in_use`.

For production clarity, I recommend **rejecting delete when in use** unless you explicitly support cleanup.

#### Errors

|Status|Code|Condition|
|---|---|---|
|401|unauthorized|Missing or invalid token|
|403|forbidden|Caller cannot delete tags|
|404|tag_not_found|Tag does not exist|
|409|tag_in_use|Tag is attached to one or more problems|

---

## 10) Attach Tags to Problem

- **URL:** `POST /orgs/{org_id}/problems/{problem_id}/tags`
- **Purpose:** Attach one or more tags to a problem.
- **Access Control:** `admin`, `mentor`

#### Headers

- `Authorization: Bearer <token>`
- `Content-Type: application/json`

#### Path Params

- `org_id` — organization UUID
- `problem_id` — problem UUID

#### Request Body
```
{  
  "tag_ids": [  
    "tag_uuid_1",  
    "tag_uuid_2"  
  ]  
}
```

#### Validation Rules

- `tag_ids` required and must not be empty
- every tag must belong to the same org
- problem must belong to the same org
- duplicates in request should be deduplicated or rejected cleanly
- existing relations should not error unless your API is strict; idempotent attach is better

#### Success Response

- `200 OK`
```
{  
  "problem_id": "prob_uuid",  
  "attached_tag_ids": [  
    "tag_uuid_1",  
    "tag_uuid_2"  
  ]  
}
```
#### Errors

|Status|Code|Condition|
|---|---|---|
|400|invalid_body|Empty list or invalid IDs|
|401|unauthorized|Missing or invalid token|
|403|forbidden|Caller cannot modify this org|
|404|problem_not_found|Problem does not exist|
|404|tag_not_found|One or more tags do not exist in this org|
|409|cross_org_violation|A tag belongs to a different organization|

---

## 11) Detach Tag from Problem

- **URL:** `DELETE /orgs/{org_id}/problems/{problem_id}/tags/{tag_id}`
- **Purpose:** Remove one tag from one problem.
- **Access Control:** `admin`, `mentor`

#### Headers

- `Authorization: Bearer <token>`

#### Path Params

- `org_id` — organization UUID
- `problem_id` — problem UUID
- `tag_id` — tag UUID

#### Success Response

- `204 No Content`

#### Errors

| Status | Code               | Condition                           |
| ------ | ------------------ | ----------------------------------- |
| 401    | unauthorized       | Missing or invalid token            |
| 403    | forbidden          | Caller cannot modify this org       |
| 404    | relation_not_found | Problem-tag relation does not exist |

---

## 12) Add Problem Resource

- **URL:** `POST /orgs/{org_id}/problems/{problem_id}/resources`
- **Purpose:** Add a learning resource to a problem.
- **Access Control:** `admin`, `mentor`

#### Headers

- `Authorization: Bearer <token>`
- `Content-Type: application/json`

#### Path Params

- `org_id` — organization UUID
- `problem_id` — problem UUID

#### Request Body
```
{  
  "title": "Official Editorial",  
  "url": "https://example.com/editorial"  
}
```

#### Validation Rules

- `title` required, 2–150 chars
- `url` required and must be a valid URL
- problem must belong to same org
- no org_id field in body
- resource count can be unlimited unless product policy says otherwise

#### Success Response

- `201 Created`
```
{  
  "id": "res_uuid",  
  "problem_id": "prob_uuid",  
  "title": "Official Editorial",  
  "url": "https://example.com/editorial",  
  "created_at": "2026-03-21T10:30:00Z"  
}
```
#### Errors

| Status | Code              | Condition                     |
| ------ | ----------------- | ----------------------------- |
| 400    | invalid_body      | Bad title or URL              |
| 401    | unauthorized      | Missing or invalid token      |
| 403    | forbidden         | Caller cannot modify this org |
| 404    | problem_not_found | Problem does not exist        |

---

## 13) List Problem Resources

- **URL:** `GET /orgs/{org_id}/problems/{problem_id}/resources`
- **Purpose:** List all resources attached to a problem.
- **Access Control:** `super_admin`, `admin`, `mentor`, `mentee`

#### Headers

- `Authorization: Bearer <token>`

#### Path Params

- `org_id` — organization UUID
- `problem_id` — problem UUID

#### Success Response

- `200 OK`
```
{  
  "items": [  
    {  
      "id": "res_uuid",  
      "title": "Official Editorial",  
      "url": "https://example.com/editorial",  
      "created_at": "2026-03-21T10:30:00Z"  
    }  
  ],  
  "total": 1  
}
```
#### Errors

| Status | Code              | Condition                              |
| ------ | ----------------- | -------------------------------------- |
| 401    | unauthorized      | Missing or invalid token               |
| 403    | forbidden         | Caller cannot access this organization |
| 404    | problem_not_found | Problem does not exist                 |

---

## 14) Update Problem Resource

- **URL:** `PATCH /orgs/{org_id}/problems/{problem_id}/resources/{resource_id}`
- **Purpose:** Update a problem resource.
- **Access Control:** `admin`, `mentor`

#### Headers

- `Authorization: Bearer <token>`
- `Content-Type: application/json`

#### Path Params

- `org_id` — organization UUID
- `problem_id` — problem UUID
- `resource_id` — resource UUID

#### Request Body
```
{  
  "title": "Updated Editorial",  
  "url": "https://example.com/new-editorial"  
}
```
#### Validation Rules

- all fields optional, but at least one must be present
- title/url must be valid if present
- resource must belong to the problem in path
- problem must belong to the org in path

#### Success Response

- `200 OK`
```
{  
  "id": "res_uuid",  
  "problem_id": "prob_uuid",  
  "title": "Updated Editorial",  
  "url": "https://example.com/new-editorial",  
  "created_at": "2026-03-21T10:30:00Z"  
}
```
#### Errors

| Status | Code               | Condition                      |
| ------ | ------------------ | ------------------------------ |
| 400    | invalid_body       | No valid fields or bad values  |
| 401    | unauthorized       | Missing or invalid token       |
| 403    | forbidden          | Caller cannot update resources |
| 404    | resource_not_found | Resource does not exist        |

---

## 15) Delete Problem Resource

- **URL:** `DELETE /orgs/{org_id}/problems/{problem_id}/resources/{resource_id}`
- **Purpose:** Remove a resource from a problem.
- **Access Control:** `admin`, `mentor`

#### Headers

- `Authorization: Bearer <token>`

#### Path Params

- `org_id` — organization UUID
- `problem_id` — problem UUID
- `resource_id` — resource UUID

#### Success Response

- `204 No Content`

#### Errors

|Status|Code|Condition|
|---|---|---|
|401|unauthorized|Missing or invalid token|
|403|forbidden|Caller cannot delete resources|
|404|resource_not_found|Resource does not exist|

### Security considerations

- Every request must verify:
    1. the caller is authenticated,
    2. the caller belongs to the org unless they are super_admin read-only,
    3. the target problem/tag/resource belongs to the same org.
- Never trust `created_by`, `organization_id`, or relationship IDs from the client when the server can infer them.
- Enforce uniqueness at the database level:
    - `tags(organization_id, name)`
    - `problem_tags(problem_id, tag_id)`
- Do not allow `super_admin` to mutate org content. Read-only only, exactly as required.
- If cookies and headers both transport JWT, make sure your middleware resolves precedence consistently and safely.

### Multi-tenant isolation rules

- Organization is the tenant boundary.
- A problem can never reference a tag from another organization.
- A resource cannot escape its parent problem, and that problem cannot escape its org.
- List endpoints must always filter by org first, not by search criteria first.
- Super admin moderation should be read-only for this layer, even if they can inspect all orgs.


# 