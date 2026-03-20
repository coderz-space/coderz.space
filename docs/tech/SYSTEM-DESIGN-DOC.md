# Modules :

Database architecture consists of 
**6 layers**:

1. Auth Module
2. Organization Module
3. Bootcamp Module
4. Problem content Module
5. Assignment system Module
6. Tracking / interaction Module 
7. Analytics Module : `V2`

---

# AUTH Module : 

## 1. System Design :

### Token Strategy : 

- **Access Token (JWT)**
    - Short-lived (15 min recommended)
    - Sent via: `Authorization: Bearer <token>` and `cookies`
- **Refresh Token**
    - Long-lived (7–30 days)
    - Stored in:
        - DB (hashed)
        - HttpOnly cookie

**Access Token Cookie Config :** 
```go
Set-Cookie: access_token=<JWT>;
Path=/;
HttpOnly;
Secure;
SameSite=Strict;
Max-Age=900;
```

**Refresh Token Cookie Config**
```http
Set-Cookie: refresh_token=xyz;
HttpOnly;
Secure;
SameSite=Strict;
Path=/auth/refresh;
```
### Rate Limiting : ( V2 )

| Endpoint              | Limit     |
| --------------------- | --------- |
| /auth/login           | 5 req/min |
| /auth/signup          | 5 req/min |
| /auth/forgot-password | 3 req/min |

---

## 2. DATABASE DESIGN 

Check the User Table in the ORGANIZATION module below.

---


---
## 2. Organization Module

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

#### Why this design is important ?

This allows multiple bootcamps in one organization.
```go
Algo University
   ├── DSA Bootcamp
   ├── Backend Bootcamp
   └── System Design Cohort
```

Users can join multiple bootcamps :
```go
Suraj
   ├── DSA Bootcamp
   └── Backend Bootcamp
```

Users can have Different roles in different bootcamps
```go
Suraj
   Mentor in Bootcamp A
   Mentee in Bootcamp B
```


### 1. Bootcamp

Represents a structured learning program inside an organization. A bootcamp groups mentors and mentees working toward a common learning goal such as a DSA cohort or a system design training.
Examples:

- 8 Week DSA Bootcamp
- Graph Algorithms Sprint
- Backend Engineering Cohort

A bootcamp belongs to an **organization**.

| Field           | Type      | Notes                               |
| --------------- | --------- | ----------------------------------- |
| id              | UUID      | PK                                  |
| organization_id | UUID      | FK → organizations.id, NOT NULL     |
| name            | string    | Bootcamp name, NOT NULL             |
| description     | text      | Optional description of the program |
| start_date      | date      | nullable                            |
| end_date        | date      | nullable                            |
| is_active       | boolean   | default true                        |
| created_at      | timestamp | NOT NULL, default CURRENT_TIMESTAMP |
| updated_at      | timestamp | NOT NULL, default CURRENT_TIMESTAMP |
#### Constraints

- **FOREIGN KEY:**  
    `organization_id → organizations.id`

### 2. bootcamp_enrollment

Represents a user's participation in a specific bootcamp.

Join table between **Bootcamp and OrganizationMember**.
Reason:  
Not every org member is part of every bootcamp.
Tracks:

- mentor / mentee role within bootcamp
- active / inactive

A member of one organization cannot join the bootcamp of another organization unless he joins that organization.
A user might belong to an organization but **not participate in every bootcamp**.

```go
Org Members
Suraj
Rahul
Priya

Bootcamp: DSA Cohort

Participants:
Suraj
Rahul
```


| Field                  | Type      | Notes                                  |
| ---------------------- | --------- | -------------------------------------- |
| id                     | UUID      | PK                                     |
| bootcamp_id            | UUID      | FK → bootcamps.id, NOT NULL            |
| organization_member_id | UUID      | FK → organization_members.id, NOT NULL |
| role                   | enum      | ENUM('mentor','mentee'), NOT NULL      |
| enrolled_at            | timestamp | NOT NULL, default CURRENT_TIMESTAMP    |


#### Constraints

- **FOREIGN KEY:**  
    `bootcamp_id → bootcamps.id`
    
- **FOREIGN KEY:**  
    `organization_member_id → organization_members.id`
    
- **UNIQUE:**  
    `(bootcamp_id, organization_member_id)`  
    Prevents duplicate enrollments.

- **Role Values**
	Roles inside a bootcamp may differ from organization roles.
	Examples:
	- **mentor** → guides students, assigns problems
	- **mentee** → solves assigned problems

---

## 3. Problem Content Layer

This layer represents the **question bank mentors use to create assignments**.
The important principle here is:
Problems should be **reusable** across assignments, bootcamps, and mentees.

so this layer will contains : 
```go
problems
tags
problem_tags
problem_resources
```



### 1. Problem

Represents a coding or conceptual problem in the organization's master question bank. Mentors create these problems and use them when assigning tasks to mentees.
Problems belong to an **organization** so that it can be reused across multiple bootcamps in an organization.

Example:
- Two Sum
- LRU Cache
- Design Twitter

| Field           | Type      | Notes                                  |
| --------------- | --------- | -------------------------------------- |
| id              | UUID      | PK                                     |
| organization_id | UUID      | FK → organizations.id, NOT NULL        |
| created_by      | UUID      | FK → organization_members.id, NOT NULL |
| title           | string    | Problem title, NOT NULL                |
| description     | text      | Problem statement                      |
| difficulty      | enum      | ENUM('easy','medium','hard'), NOT NULL |
| external_link   | string    | Optional external problem reference    |
| created_at      | timestamp | NOT NULL, default CURRENT_TIMESTAMP    |
| updated_at      | timestamp | NOT NULL, default CURRENT_TIMESTAMP    |
#### Constraints

FOREIGN KEY:  
`organization_id → organizations.id`

FOREIGN KEY:  
`created_by → organization_members.id`


### 2. Tags

**Description:**  
Represents a concept or topic used to categorize problems. Tags allow mentors and mentees to filter problems based on concepts such as algorithms, data structures, or system design topics.

Examples:
- binary-search
- sliding-window
- bit-manipulation
- cap-theorem
- dynamic-programming
- graphs
    
Tags are **global to the organization**.

| Field           | Type      | Notes                               |
| --------------- | --------- | ----------------------------------- |
| id              | UUID      | PK                                  |
| organization_id | UUID      | FK → organizations.id               |
| name            | string    | Tag name, NOT NULL                  |
| created_at      | timestamp | NOT NULL, default CURRENT_TIMESTAMP |
#### Constraints

UNIQUE:
	(organization_id, name)
	Prevents duplicate tags in the same organization.

FOREIGN KEY:
	organization_id → organizations.id



### 3. Problem tags : 

**Description:**  
Join table connecting **problems and tags**, allowing each problem to have multiple tags and each tag to be used by many problems.

| Field      | Type      | Notes                               |
| ---------- | --------- | ----------------------------------- |
| problem_id | UUID      | FK → problems.id, NOT NULL          |
| tag_id     | UUID      | FK → tags.id, NOT NULL              |
| created_at | timestamp | NOT NULL, default CURRENT_TIMESTAMP |
#### Constraints

PRIMARY KEY:
	(problem_id, tag_id)

FOREIGN KEY:
	problem_id → problems.id  
	tag_id → tags.id


### 4. problem_resources

**Description:**  
Optional learning resources attached to a problem to help mentees understand concepts or solutions.
Examples:
- YouTube explanation
- Blog article
- Documentation
- Editorial

Relationship:  
Problem → many resources

| Field      | Type      | Notes                               |
| ---------- | --------- | ----------------------------------- |
| id         | UUID      | PK                                  |
| problem_id | UUID      | FK → problems.id, NOT NULL          |
| title      | string    | Resource title                      |
| url        | string    | Resource link                       |
| created_at | timestamp | NOT NULL, default CURRENT_TIMESTAMP |
|            |           |                                     |
#### Constraints

FOREIGN KEY:
	problem_id → problems.id

---

## 4. Assignment layer

The main principles we follow:
1. **Problem sets should be reusable**
2. **Assignments should be per mentee**
3. **Deadlines should be flexible**

Personalized deadlines
- Mentor can assign the same assignment group to two students with different deadlines.

It includes the following tables : 
```go
assignment_groups
assignment_group_problems
assignments
assignment_problems
```

### 1. Assignment Group

**Description:**  
Represents a set of problems created by a mentor with a defined deadline duration. This acts as a reusable template that can be assigned to multiple mentees.

Examples:
- Sliding Window Practice
- Graph Algorithms Sprint
- Week 1 DSA
- Dynamic Programming Set

An assignment group belongs to a **bootcamp**.

| Field         | Type      | Notes                                             |
| ------------- | --------- | ------------------------------------------------- |
| id            | UUID      | PK                                                |
| bootcamp_id   | UUID      | FK → bootcamps.id, NOT NULL                       |
| created_by    | UUID      | FK → organization_members.id, NOT NULL            |
| title         | string    | Assignment set name, NOT NULL                     |
| description   | text      | Optional explanation                              |
| deadline_days | integer   | Number of days allowed to complete the assignment |
| created_at    | timestamp | NOT NULL, default CURRENT_TIMESTAMP               |
| updated_at    | timestamp | NOT NULL, default CURRENT_TIMESTAMP               |

#### Constraints

FOREIGN KEY  
`bootcamp_id → bootcamps.id`

FOREIGN KEY  
`created_by → organization_members.id`


### 2. AssignmentGroupProblem

**Description:**  
Join table connecting problems to an assignment group.  
Defines which problems belong to a specific assignment group.

This allows **reusing problems across many assignments**.

| Field               | Type      | Notes                               |
| ------------------- | --------- | ----------------------------------- |
| assignment_group_id | UUID      | FK → assignment_groups.id, NOT NULL |
| problem_id          | UUID      | FK → problems.id, NOT NULL          |
| position            | integer   | Optional ordering of problems       |
| created_at          | timestamp | NOT NULL, default CURRENT_TIMESTAMP |
#### Constraints

PRIMARY KEY
	(assignment_group_id, problem_id)

FOREIGN KEY
	assignment_group_id → assignment_groups.id  
	problem_id → problems.id

### 3. assignment

**Description:**  
Represents an **assignment instance for a specific mentee**.  
Created when a mentor assigns an assignment group to a mentee.

Each mentee gets their **own assignment instance**, allowing personalized deadlines and progress tracking.

```go
AssignmentGroup: Graph Practice
Assigned to: Suraj
Deadline: 12 March
```

| Field                  | Type      | Notes                                   |
| ---------------------- | --------- | --------------------------------------- |
| id                     | UUID      | PK                                      |
| assignment_group_id    | UUID      | FK → assignment_groups.id               |
| bootcamp_enrollment_id | UUID      | FK → bootcamp_enrollments.id            |
| assigned_by            | UUID      | FK → organization_members.id            |
| assigned_at            | timestamp | default CURRENT_TIMESTAMP               |
| deadline_at            | timestamp | Calculated based on assignment duration |
| status                 | enum      | ENUM('active','completed','expired')    |
| created_at             | timestamp | default CURRENT_TIMESTAMP               |
| updated_at             | timestamp | default CURRENT_TIMESTAMP               |

#### Constraints

FOREIGN KEY
	assignment_group_id → assignment_groups.id

FOREIGN KEY
	bootcamp_enrollment_id → bootcamp_enrollments.id

FOREIGN KEY
	assigned_by → organization_members.id

### 4. assignment problem

**Description:**  
Represents the **problems assigned within a specific assignment instance**.

This table tracks **mentee progress per problem**.

#### Status Values

- pending
- attempted
- completed

| Field         | Type      | Notes                                   |
| ------------- | --------- | --------------------------------------- |
| id            | UUID      | PK                                      |
| assignment_id | UUID      | FK → assignments.id                     |
| problem_id    | UUID      | FK → problems.id                        |
| status        | enum      | ENUM('pending','attempted','completed') |
| solution_link | string    | Optional link to solution               |
| notes         | text      | Optional notes from mentee              |
| completed_at  | timestamp | nullable                                |
| created_at    | timestamp | default CURRENT_TIMESTAMP               |
| updated_at    | timestamp | default CURRENT_TIMESTAMP               |

#### Constraints

FOREIGN KEY
	assignment_id → assignments.id  
	problem_id → problems.id

UNIQUE
	(assignment_id, problem_id)

Prevents duplicate problems in the same assignment.


---

### Note :

#### 1. Why do we have 2 fields to track deadline ? One in `assignment_group` and other in `assignment`.

One is a **template rule**, the other is a **real deadline for a specific mentee**.
`assignment_groups.deadline_days` :
- This defines the **default duration of the task set**.
	It belongs to the **template created by the mentor** and not tied to any specific mentee.
`assignments.deadline_at` :
	- This is the actual deadline for a specific mentee.




---

## 5. Progress Tracking Layer

### 1. Doubt

Description:  
Represents a question or doubt raised by a mentee while solving an assigned problem.
Doubts allow mentors to track **problematic questions across mentees** and identify commonly difficult concepts.
A doubt is linked to a specific **assignment problem**.

| Field                 | Type      | Notes                                  |
| --------------------- | --------- | -------------------------------------- |
| id                    | UUID      | PK                                     |
| assignment_problem_id | UUID      | FK → assignment_problems.id            |
| raised_by             | UUID      | FK → organization_members.id           |
| message               | text      | Description of the doubt               |
| resolved              | boolean   | default false                          |
| resolved_by           | UUID      | FK → organization_members.id, nullable |
| resolved_at           | timestamp | nullable                               |
| created_at            | timestamp | default CURRENT_TIMESTAMP              |
## Constraints

FOREIGN KEY:
	assignment_problem_id → assignment_problems.id

FOREIGN KEY:
	raised_by → organization_members.id

FOREIGN KEY:
	resolved_by → organization_members.id

---

## 6. Analytics Layer

### 1. LeaderboardEntry

Description:  
Represents aggregated performance metrics of a mentee within a bootcamp. This table stores periodic snapshots used for leaderboards and mentor analytics dashboards.
This prevents recalculating heavy analytics every time.

Example metrics:
- total problems solved
- completion rate
- streak
- rank

| Field                  | Type      | Notes                        |
| ---------------------- | --------- | ---------------------------- |
| id                     | UUID      | PK                           |
| bootcamp_id            | UUID      | FK → bootcamps.id            |
| bootcamp_enrollment_id | UUID      | FK → bootcamp_enrollments.id |
| problems_completed     | integer   | total completed problems     |
| problems_attempted     | integer   | attempted problems           |
| completion_rate        | float     | percentage completion        |
| streak_days            | integer   | current solving streak       |
| score                  | integer   | leaderboard score            |
| rank                   | integer   | rank within bootcamp         |
| calculated_at          | timestamp | when metrics were generated  |
#### Constraints

FOREIGN KEY
	bootcamp_id → bootcamps.id

FOREIGN KEY
	bootcamp_enrollment_id → bootcamp_enrollments.id

UNIQUE
	(bootcamp_id, bootcamp_enrollment_id)

### 2. Poll

Description:  
Mentors can create polls to determine which problems were difficult for the cohort.
Example:
	Was LRU Cache difficult?

| Field       | Type      | Notes                        |
| ----------- | --------- | ---------------------------- |
| id          | UUID      | PK                           |
| bootcamp_id | UUID      | FK → bootcamps.id            |
| problem_id  | UUID      | FK → problems.id             |
| question    | string    | poll question                |
| created_by  | UUID      | FK → organization_members.id |
| created_at  | timestamp | default CURRENT_TIMESTAMP    |

Foreign Key: 
```go
bootcamp_id → bootcamps.id
problem_id → problems.id
created_by → organization_members.id
```

### 3. PollVote

Description:  
Stores votes from mentees on a poll.
Example options:
- easy
- medium
- hard

| Field      | Type      | Notes                        |
| ---------- | --------- | ---------------------------- |
| id         | UUID      | PK                           |
| poll_id    | UUID      | FK → polls.id                |
| voter_id   | UUID      | FK → bootcamp_enrollments.id |
| vote       | enum      | ENUM('easy','medium','hard') |
| created_at | timestamp | default CURRENT_TIMESTAMP    |
|            |           |                              |

#### Constraints

FOREIGN KEY
	poll_id → polls.id  
	voter_id → bootcamp_enrollments.id

UNIQUE
	(poll_id, voter_id)

Each mentee votes once per poll.


---
## NOTE :

1. Never hard delete important entities like:
	- problems
	- assignments
	- bootcamps
	
	Because deleting them can break historical data.
	Instead add : `archived_at TIMESTAMP NULL`
	meaning : 

```go
archived_at = NULL → active
archived_at != NULL → archived
```


2. Add created_by for important entities : 

`created_by UUID → organization_members.id`

Add this to :
```go
bootcamps
assignment_groups
problems
tags
```

You can easily build features like:
Mentor dashboard:
```go
Show problems created by Suraj
Show assignments created by Rahul
```


