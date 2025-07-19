# Chirpy — RESTful API Server in Go

## 🚀 Project Overview

**Chirpy** is a production-style RESTful HTTP API server built with Go and PostgreSQL. It powers a simple microblogging platform where users can register, log in, post chirps, and manage their data securely.

This project was built without frameworks to deeply understand HTTP servers, JWT-based authentication, and secure database operations in Go.

---

## 🎯 Project Goals

- Understand HTTP server internals and RESTful APIs
- Build a backend server in Go using the standard library
- Communicate with clients using JSON, HTTP headers, and status codes
- Use SQLC for type-safe Postgres database access
- Implement secure authentication with JWT and refresh tokens
- Handle webhook events and external system integrations

---

## 🗺️ System Architecture

```
Client → HTTP Handlers → Business Logic → Database (PostgreSQL)
```

- **HTTP Handlers:** Handle incoming HTTP requests
- **Business Logic:** Validation, authentication, database interaction
- **Database:** Managed with SQLC and Goose migrations

---

## 📋 API Reference

| Method   | Endpoint                | Description                                                  |
| -------- | ----------------------- | ------------------------------------------------------------ |
| `POST`   | `/api/users`            | Register a new user                                          |
| `POST`   | `/api/login`            | Authenticate user and get JWT + Refresh Token                |
| `POST`   | `/api/chirps`           | Create a new chirp (Authenticated)                           |
| `GET`    | `/api/chirps`           | Retrieve chirps (supports `author_id` & `sort` query params) |
| `DELETE` | `/api/chirps/{chirpID}` | Delete a chirp (Author only)                                 |
| `POST`   | `/api/polka/webhooks`   | Handle user upgrade events (Webhook)                         |
| `POST`   | `/admin/reset`          | Reset users and hit counter (Admin only)                     |

---

## 🗄️ Database Schema

### `users` table

| Column            | Type        | Description                              |
| ----------------- | ----------- | ---------------------------------------- |
| `id`              | `UUID`      | Primary key                              |
| `email`           | `TEXT`      | User email (unique)                      |
| `hashed_password` | `TEXT`      | Bcrypt hashed password                   |
| `is_chirpy_red`   | `BOOLEAN`   | Chirpy Red membership (default: `false`) |
| `created_at`      | `TIMESTAMP` | Creation time                            |
| `updated_at`      | `TIMESTAMP` | Last update time                         |

### `chirps` table

| Column       | Type        | Description            |
| ------------ | ----------- | ---------------------- |
| `id`         | `UUID`      | Primary key            |
| `body`       | `TEXT`      | Chirp content          |
| `user_id`    | `UUID`      | Foreign key to `users` |
| `created_at` | `TIMESTAMP` | Creation time          |
| `updated_at` | `TIMESTAMP` | Last update time       |

### `refresh_tokens` table

| Column       | Type        | Description            |
| ------------ | ----------- | ---------------------- |
| `token`      | `TEXT`      | Refresh token          |
| `user_id`    | `UUID`      | Foreign key to `users` |
| `created_at` | `TIMESTAMP` | Creation time          |
| `updated_at` | `TIMESTAMP` | Last update time       |
| `expires_at` | `TIMESTAMP` | Expiration time        |

---

## 🔐 Authentication

- Access Tokens: JWTs valid for **1 hour**
- Refresh Tokens: Stored in DB, valid for **60 days**
- Passwords hashed with **bcrypt**
- Access control on protected endpoints

---

## 📨 Webhooks

- Accepts `user.upgraded` event
- Upgrades `is_chirpy_red` flag on the corresponding user
- Responds with `204 No Content` if successful

---

## 🛠️ Setup & Run

```bash
# Clone your fork
git clone https://github.com/<your-username>/chirpy.git

# Navigate to project directory
cd chirpy

# Install dependencies
go mod tidy

# Run database migrations
goose postgres "postgres://<user>:<password>@localhost:5432/chirpy" up

# Generate type-safe queries
sqlc generate

# Run the server
go run .
```

---

## 🧪 Testing

- Use `bootdev run <test_id>` for integration tests
- Reset state with `POST /admin/reset` during testing

---

## 📝 Lessons Learned

- Writing idiomatic Go HTTP servers
- Working with JWTs securely
- Type-safe DB access with SQLC
- Secure password handling with bcrypt
- Building production-ready REST APIs

---

## 🚀 Future Improvements

- Middleware for logging, rate limiting, and metrics
- HTTPS support and deployment
- CI/CD pipeline with unit tests
- Improved error handling and API documentation

---

Feel free to contribute and make Chirpy even better!

