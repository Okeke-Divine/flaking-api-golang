# Flaking API

[![Tech Stack](https://skillicons.dev/icons?i=go,mysql,postgres,docker)](https://skillicons.dev)

## 📖 Overview
A production-ready RESTful API built with Go and the Gin framework to demonstrate the core backend development concepts, including authentication, database operations, and middleware patterns.

## ✨ Features
- **🔐 Authentication:** Secure JWT-based login & registration
- **📁 RESTful Endpoints:** Clean CRUD operations for users and posts with pagination
- **🗄️ Database Integration:** GORM ORM with support for MySQL & PostgreSQL
- **⚙️ Middleware:** Custom JWT auth, rate limiting, and request validation
- **🛡️ Security:** Bcrypt password hashing and input sanitization
- **📦 Extensible Design:** Modular structure for easy maintenance and scalability

## 🛠️ Tech Stack
- **Language:** Go (Golang)
- **Framework:** Gin
- **ORM:** GORM
- **Database:** MySQL/PostgreSQL
- **Authentication:** JWT
- **Password Hashing:** Bcrypt


## 🚀 Getting Started

### Prerequisites
- Go 1.18+
- MySQL or PostgreSQL
- Git

### Installation & Setup

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/Okeke-Divine/flaking-api-golang
    cd flaking-api-golang
    ```

2.  **Install dependencies:**
    ```bash
    go mod download
    ```

3.  **Set up environment variables:**
    ```bash
    cp .env.example .env
    ```
    Edit the `.env` file with your database credentials and a secure JWT secret.

4.  **Set up your database:**
    Create a new database in your MySQL/PostgreSQL instance matching the name in your `.env` file.

5.  **Run the application:**
    ```bash
    go run main.go
    ```
    The API server will start on `http://localhost:4000`.

## 📚 API Endpoints

| Method | Endpoint | Description | Auth Required |
| :--- | :--- | :--- | :---: |
| `POST` | `/api/v1/auth/register` | Register a new user | No |
| `POST` | `/api/v1/auth/login` | Login user | No |
| `GET` | `/api/v1/users` | Get all users (public) | No |
| `GET` | `/api/v1/users/:id` | Get user by ID (public) | No |
| `GET` | `/api/v1/user/profile` | Get current user profile | Yes |
| `PUT` | `/api/v1/user/profile` | Update current user | Yes |
| `DELETE` | `/api/v1/user/profile` | Delete current user | Yes |
| `GET` | `/api/v1/posts` | Get all posts (public) | No |
| `GET` | `/api/v1/posts/:id` | Get post by ID (public) | No |
| `POST` | `/api/v1/posts` | Create a new post | Yes |
| `PUT` | `/api/v1/posts/:id` | Update a post (owner only) | Yes |
| `DELETE` | `/api/v1/posts/:id` | Delete a post (owner only) | Yes |

## 🔐 Authentication
Protected routes require a JWT token in the Authorization header:
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

## 📋 What I learnt

- Go's HTTP handling and the Gin framework
- Structuring large Go applications effectively
- JWT implementation for stateless authentication
- Middleware chaining and custom middleware creation
- Database modeling with GORM and relationships
- Rate limiting strategies and implementation
- RESTful API design principles

## 📄 License

This project is licensed under the MIT License.

