# Flaking API · [![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://golang.org) [![Gin](https://img.shields.io/badge/Gin-Framework-0096D0?style=flat)](https://gin-gonic.com)

> **A Production-Grade RESTful API Boilerplate in Go. Built for Performance, Security, and Scale.**

## 📖 Overview
**Flaking API** is a robust, modular backend **boilerplate** built with **Go** and the **Gin framework**. It is designed to be cloned and used as a high-performance foundation for building any web or mobile application backend, implementing core production-ready features out of the box.
## ⚡ Why Flaking API? The Engineering Choice

| For Your Application ✅ | For You (The Developer) 🔧 |
| :--- | :--- |
| **Blazing Fast Performance** – Native Go runtime efficiency. | **Structured for Scale** – Clean, modular architecture (config, controllers, models, middleware). |
| **Secure by Default** – JWT authentication, bcrypt hashing, input validation. | **Production-Ready Features** – Rate limiting, CORS, centralized error handling. |
| **Multi Database Support** – Supports MySQL & PostgreSQL via GORM. | **Go & Gin Mastery** – Demonstrates deep understanding of idiomatic Go and a modern framework. |
| **RESTful Standards** – Predictable, well-documented endpoints for easy integration. | **Minimal & Efficient** – No bloated dependencies; a lean, powerful codebase. |

## 🛠️ Tech Stack & Architecture
- **Language:** Go 1.25+
- **Framework:** Gin Gonic
- **ORM:** GORM (MySQL & PostgreSQL support)
- **Authentication:** JWT (JSON Web Tokens)
- **Security:** Bcrypt for password hashing
- **Middleware:** Custom JWT Auth, Rate Limiting

## 📈 Core Features
- **🔐 JWT Authentication:** Secure user registration, login, and protected route management.
- **⏱️ Rate Limiting:** Configurable middleware to protect against abuse and DDoS.
- **🗃️ Database Operations:** Full CRUD for users and posts with GORM, including pagination and data relations.
- **✅ Input Validation:** Robust validation for all incoming requests.
- **🧩 Modular Design:** Separated concerns with clear structure (models, controllers, routes, middleware).
- **♻️ RESTful Design:** Intuitive API endpoints following industry standards.

## 🚀 Getting Started

### Prerequisites
- Go 1.25+
- MySQL or PostgreSQL database

### Installation & Local Development
1.  **Clone the repository:**
    ```bash
    git clone https://github.com/Okeke-Divine/flaking-api-golang.git
    cd flaking-api-golang
    ```

2.  **Install dependencies:**
    ```bash
    go mod download
    ```

3.  **Configure environment variables:**
    ```bash
    cp .env.example .env
    ```
    Edit the `.env` file with your settings:
    ```bash
    # Database
    DB_DRIVER=mysql
    DB_HOST=localhost
    DB_PORT=3306
    DB_USER=your_username
    DB_PASSWORD=your_password
    DB_NAME=your_database_name

    # JWT
    JWT_SECRET=your_very_secure_secret_key_here

    # Server
    PORT=4000
    ```

4.  **Run database migrations (if using GORM AutoMigrate):**
    *   The application will typically create the necessary tables on first run if `AUTO_MIGRATE` is set.

5.  **Start the development server:**
    ```bash
    go run main.go
    ```
    The API server will start on `http://localhost:4000`.

## 📚 API Endpoints & Documentation

### Authentication Endpoints
| Method | Endpoint | Description | Auth Required |
| :--- | :--- | :--- | :---: |
| `POST` | `/api/v1/auth/register` | Register a new user | No |
| `POST` | `/api/v1/auth/login` | Login user | No |

### User Endpoints
| Method | Endpoint | Description | Auth Required |
| :--- | :--- | :--- | :---: |
| `GET` | `/api/v1/users` | Get all users (public) | No |
| `GET` | `/api/v1/users/:id` | Get user by ID (public) | No |
| `GET` | `/api/v1/user/profile` | Get current user profile | **Yes** |
| `PUT` | `/api/v1/user/profile` | Update current user | **Yes** |
| `DELETE` | `/api/v1/user/profile` | Delete current user | **Yes** |

### Post Endpoints
| Method | Endpoint | Description | Auth Required |
| :--- | :--- | :--- | :---: |
| `GET` | `/api/v1/posts` | Get all posts (public) | No |
| `GET` | `/api/v1/posts/:id` | Get post by ID (public) | No |
| `POST` | `/api/v1/posts` | Create a new post | **Yes** |
| `PUT` | `/api/v1/posts/:id` | Update a post (owner only) | **Yes** |
| `DELETE` | `/api/v1/posts/:id` | Delete a post (owner only) | **Yes** |

**🔐 Authentication:** Protected routes require a JWT token in the Authorization header:
```http
Authorization: Bearer <your_jwt_token>
```

## 🧪 Testing the API

1. Register a new user:
```bash
curl -X POST http://localhost:4000/api/v1/auth/register   -H "Content-Type: application/json"   -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

2. Login to get a JWT token:
```bash
curl -X POST http://localhost:4000/api/v1/auth/login   -H "Content-Type: application/json"   -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

3. Create a post (using the JWT token from login):
```bash
curl -X POST http://localhost:4000/api/v1/posts   -H "Content-Type: application/json"   -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE"   -d '{
    "title": "My First Post",
    "content": "This is the content of my first post."
  }'
```

## 🧠 Knowledge Demonstrated

- Go Conventions: Proper project structure, error handling, and modular packages.
- API Design: RESTful principles, clean endpoint design, and HTTP status codes.
- Database Modeling: ORM usage (GORM), migrations, and relationships.
- Security: JWT implementation, password hashing (bcrypt), and middleware for auth/rate limiting.
- Production Readiness: Environment configuration, and scalable architecture.

---
**Built by [Divine-Vessel](https://github.com/Okeke-Divine)**
