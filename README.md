## Chirpy

Chirpy is a custom-built social media backend API written in Go. This project serves as a robust exercise in building scalable HTTP servers, handling persistent data, and implementing secure user authentication.

## Features

- **RESTful API**: Clean endpoints for managing "chirps" and user profiles.
- **Authentication**: Secure user login using JWT (JSON Web Tokens) and refresh token logic.
- **Security**: Password hashing using Bcrypt to ensure user data protection.
- **Input Validation**: Server-side filtering (e.g., removing profanity from chirps).

## Tech Stack

- **Language**: Go (Golang)
- **Database**: Local JSON-based storage (or PostgreSQL if you've migrated)
- **Middleware**: Custom logging and authentication middleware

## Getting Started

### Prerequisites

- Go 1.21 or higher installed on your machine.

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/chirpy.git
   cd chirpy

Build the application:
```bash
go build -o out
```

Run the server:

./out

The server will start on port 8080 by default.

API Documentation
GET /api/healthz: Check server health.
POST /api/users: Create a new user.
POST /api/login: Authenticate and receive a JWT.
POST /api/chirps: Create a new chirp (authenticated).

## Configuration

This project relies on environment variables for security and integration. Create a `.env` file in the root of the project and populate it with the following:

- `DB_URL`: The path to your local database file or your database connection string.
- `PLATFORM`: Set to `dev` to enable development-only features (like the reset endpoint).
- `SECRET`: A secure, private string used for signing and verifying JSON Web Tokens.
- `POLKA_KEY`: The API key required to authenticate webhooks from the Polka payment gateway.

### Setup

```bash
cp .env.example .env
# Then edit .env with your actual secrets.
```

## Lessons Learned

- **RESTful API Design**: Built structured endpoints using Go's `net/http` and handled complex JSON request/response cycles.
- **Relational Databases**: Integrated SQL for persistent storage, focusing on schema design and efficient querying.
- **Stateless Authentication**: Implemented secure user flows using JWTs, refresh tokens, and password hashing with Bcrypt.
- **Security Best Practices**: Managed sensitive credentials via environment variables and implemented secure webhook validation.
- **Scalable Architecture**: Developed a modular backend capable of handling concurrent requests and middleware logic.


