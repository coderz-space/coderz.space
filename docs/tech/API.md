



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

