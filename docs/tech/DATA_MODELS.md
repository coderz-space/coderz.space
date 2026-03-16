# Data Model

Database architecture consists of **4 layers**:

1. Organization structure
2. Learning content
3. Assignment system
4. Tracking / interaction

---

## 1. Organization Layer

### 1. organization

**Description:**  
Represents a company, community, or group that runs bootcamps on the platform.  
Acts as the **tenant boundary** for all resources such as bootcamps, problems, and assignments.

- Example : `Team Shiksha` , `Monkeys` , `Real Dev Squad`

| Field       | Type      | Notes                                    |
| ----------- | --------- | ---------------------------------------- |
| id          | UUID      | PK                                       |
| name        | string    | Organization name, NOT NULL              |
| slug        | string    | unique identifier used in URLs, NOT NULL |
| description | text      | optional                                 |
| created_at  | timestamp | NOT NULL, default CURRENT_TIMESTAMP      |
| updated_at  | timestamp | NOT NULL, default CURRENT_TIMESTAMP      |

#### Constraints

- **UNIQUE:** slug
    - Example : coderz-space, algo-university, dsa-club
    - URL : coderzspace.com/org/coderz-space

### 2. user

**Description:**  
Represents a global platform account. A user exists independently of any organization and can belong to multiple organizations with different roles.

**Authentication:**  
Supports multiple login methods such as **email/password** and **Google OAuth**.

| Field          | Type      | Notes                                                 |
| -------------- | --------- | ----------------------------------------------------- |
| id             | UUID      | PK                                                    |
| name           | string    | User's full name, NOT NULL                            |
| email          | string    | unique, nullable                                      |
| email_verified | boolean   | default false, NOT NULL                               |
| password_hash  | string    | nullable, required when google_id is NULL             |
| google_id      | string    | unique, nullable, required when password_hash is NULL |
| avatar_url     | string    | nullable                                              |
| created_at     | timestamp | NOT NULL, default CURRENT_TIMESTAMP                   |
| updated_at     | timestamp | NOT NULL, default CURRENT_TIMESTAMP, auto-update      |

#### Constraints

- **CHECK:** At least one authentication method must exist  
   `(password_hash IS NOT NULL OR google_id IS NOT NULL)`
- **UNIQUE:** email
- **UNIQUE:** google_id

### 3. organization_member

- Join table between user and organization.
- A user can be mentor in one organization and mentee in another.
  **Description:**  
  Represents the relationship between a **user and an organization**.  
  Stores the user's **role within that organization** and enables role-based access control.

#### Role Values

- **owner** → Creator of the organization with full permissions
- **admin** → Organization administrators managing members and bootcamps
- **mentor** → Mentors who assign problems and review progress
- **mentee** → Students participating in bootcamps

| Field           | Type      | Notes                                             |
| --------------- | --------- | ------------------------------------------------- |
| id              | UUID      | PK                                                |
| organization_id | UUID      | FK → organizations.id, NOT NULL                   |
| user_id         | UUID      | FK → users.id, NOT NULL                           |
| role            | enum      | ENUM('owner','admin','mentor','mentee'), NOT NULL |
| joined_at       | timestamp | NOT NULL, default CURRENT_TIMESTAMP               |

#### Constraints

- **UNIQUE:** `(organization_id, user_id)`  
   Prevents a user from joining the same organization multiple times.
- **FOREIGN KEY:**  
   `organization_id → organizations.id`
- **FOREIGN KEY:**  
   `user_id → users.id`

---

## 2. Bootcamp Layer

### 1. Bootcamp

A learning program under an organization.
Example:

- "6 Week DSA Bootcamp"
- "System Design Cohort"

Bootcamp belongs to:

- Organization

### 2. bootcamp_enrollment

Join table between **Bootcamp and OrganizationMember**.
Reason:  
Not every org member is part of every bootcamp.
Tracks:

- mentor / mentee role within bootcamp
- active / inactive

A member of one organization cannot join the bootcamp of another organization uless he joins that organization.

---

## 3. Learning Content Layer

### 1. Problem

Master question bank.
Mentors create these.
Example:

- Two Sum
- LRU Cache
- Design Twitter

Belongs to:

- Organization (important for multi tenant)

### 2. problem_resources

Optional resources attached to a problem.

Examples:

- YouTube link
- Blog
- Docs
- Article

Relationship:  
Problem → many resources

---

## 4. Assignment layer

Personalized deadlines

- Mentor can assign the same assignment group to two students with different deadlines.

Reusable problem set :

- Mentor creates `Graph fundamentals` and reuses it across multiple batches.

### 1. Assignment Group

Represents a **set of problems**.

### 2. AssignmentGroupProblem

Join table : AssignmentGroup ↔ Problem

### 3. assignment

Represents the **assignment for a specific mentee**.
Because mentors may assign the same group to different mentees.

```go
AssignmentGroup: Graph Practice
Assigned to: Suraj
Deadline: 12 March
```

### 4. assignment problem

Tracks **each problem inside the assignment**.

---

## 5. Progress Tractking Layer

### 1. Submission

Tracks mentee progress on an assigned problem.
Example:

- completed
- attempted
- pending

Contains optional:

- solution link
- notes

### 2. Doubt

If a mentee highlights a problem.
Example:  
"I don't understand sliding window here."
Linked to:

- AssignmentProblem
- User

---

## 6. Analytics Layer

### 1. LeaderboardEntry

Tracks score per user per bootcamp.
Example metrics:

- problems solved
- streak
- weekly score

### 2. Poll

Mentor polls problems to see difficulty.
Example:  
"Was LRU Cache hard?"

### 3. PollVote

User vote on poll.

### 2. Poll

Mentor polls problems to see difficulty.
Example:  
"Was LRU Cache hard?"

### 3. PollVote

User vote on poll.
