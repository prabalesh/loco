# Loco - Competitive Programming Platform

A modern, full-stack competitive programming platform built with **Go** and **React**. Practice coding problems, submit solutions, track your progress, and compete on the leaderboard.

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8.svg)
![React Version](https://img.shields.io/badge/react-19.0+-61DAFB.svg)

---

## 📋 Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Running the Application](#running-the-application)
- [API Documentation](#api-documentation)
- [Security Features](#security-features)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)

---

## ✨ Features

- 🔐 **Secure Authentication** - JWT-based auth with HttpOnly cookies
- 👤 **User Profiles** - Public and private profile views with gamification stats
- 🎮 **Gamification** - Earn XP, level up, unlock achievements, and maintain submission streaks
- 💻 **Problem Library** - Comprehensive collection of coding challenges with multi-language support
- ✅ **Code Submission** - Secure code execution via Piston with real-time feedback
- ✅ **Reference Solution Validation** - Mandatory problem validation against reference solutions before publishing
- 🏆 **Leaderboard** - Global rankings based on XP and problem-solving prowess
- 🛠️ **Admin Portal** - Robust management interface for problems, users, and system monitoring
- 📊 **Progress Tracking** - Detailed statistics and progress visualization
- 🎨 **Modern UI** - Premium design system using Tailwind CSS v4 and Framer Motion
- 🌙 **Dark Mode Ready** - Fully responsive design with native dark theme support

---

## 🛠️ Tech Stack

### Backend
- **Language:** Go 1.21+
- **Framework:** net/http (Standard Library)
- **Database:** PostgreSQL 15+
- **Queue & Cache:** Redis
- **Code Execution:** Piston API
- **Authentication:** JWT with HttpOnly cookies
- **Logging:** Zap
- **Database Driver:** pgx

### Frontend (Main & Admin)
- **Framework:** React 19 + TypeScript
- **Build Tool:** Vite 7
- **Styling:** Tailwind CSS v4
- **State Management:** Zustand
- **Data Fetching:** TanStack Query v5
- **UI Components:** Framer Motion, Lucide React, MUI (Admin)
- **Editor:** Monaco Editor, Tiptap (Admin Rich Text)
- **Form Handling:** React Hook Form + Zod
- **Routing:** React Router v7

---

## 📁 Project Structure

```bash
loco/
├── backend/                      # Go backend (clean architecture)
│   ├── cmd/server/               # App entry point
│   ├── internal/                 # Private application and library code
│   │   ├── delivery/             # HTTP handlers & middleware
│   │   ├── domain/               # Domain entities (Models, Usecases, Repositories interfaces)
│   │   ├── usecase/              # Business logic implementation
│   │   ├── repository/           # Database access (GORM/SQL)
│   │   └── infrastructure/       # External services (Redis, Piston, JWT)
│   ├── migrations/               # Database migrations
│   └── pkg/                      # Shared utility packages
│
├── frontend/                     # Main user application (React 19)
│   ├── src/features/             # Feature-based organization (Auth, Problems, Profile)
│   ├── src/shared/               # Reusable components and hooks
│   └── tailwind.config.js
│
├── admin-frontend/               # Admin dashboard (React 19 + MUI)
│   ├── src/pages/                # Admin views (Dashboard, Problem Management)
│   └── src/features/             # Admin specific logic
│
├── README.md
└── LICENSE

```

---

## 📦 Prerequisites

Before you begin, ensure you have the following installed:

- **Go** 1.21+
- **Node.js** 20+
- **PostgreSQL** 15+
- **Redis** 7+
- **Docker** (Optional, for running execution engine)

---

## 🚀 Installation

### 1. Clone the Repository

```bash
git clone https://github.com/prabalesh/loco.git
cd loco
```

### 2. Backend Setup

```bash
cd backend

# Install Go dependencies
go mod download

# Create PostgreSQL database
createdb coding_platform

# Or using psql
psql -U postgres -c "CREATE DATABASE coding_platform;"

# Run migrations
psql -d coding_platform -f migrations/001_create_users_table.sql
```

### 3. Frontend Setup (Main App)

```bash
cd frontend
npm install
```

### 4. Admin Frontend Setup

```bash
cd admin-frontend
npm install
```

---

## ⚙️ Configuration

### Backend Configuration

Create `backend/.env` file:

```env
# Server Configuration
PORT=8080
ENV=development

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_secure_password_here
DB_NAME=coding_platform
DB_SSL_MODE=disable

# JWT Secrets (CHANGE THESE IN PRODUCTION!)
ACCESS_TOKEN_SECRET=your-super-secret-access-token-key-minimum-32-characters-long
REFRESH_TOKEN_SECRET=your-super-secret-refresh-token-key-minimum-32-characters-long
ACCESS_TOKEN_EXPIRATION=15m
REFRESH_TOKEN_EXPIRATION=168h

# Cookie Settings
COOKIE_SAMESITE=lax
COOKIE_DOMAIN=

# CORS Configuration (comma-separated for multiple origins)
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173

# Logging
LOG_LEVEL=info
```

### Frontend Configuration

Create `frontend/.env` file:

```env
VITE_API_BASE_URL=http://localhost:8080
```

---

## ▶️ Running the Application

### Development Mode

#### Terminal 1: Start Backend

```bash
cd backend
go run cmd/server/main.go
```

✅ Backend runs on **http://localhost:8080**

#### Terminal 2: Start Frontend

```bash
cd frontend
npm run dev
```

✅ Frontend runs on **http://localhost:3000**

### Production Build

#### Backend

```bash
cd backend
go build -o bin/server cmd/server/main.go
./bin/server
```

#### Frontend

```bash
cd frontend
npm run build
npm run preview
```

Or deploy the `dist/` folder to any static hosting service.

---

## 🔌 API Documentation

### Authentication Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/auth/register` | Register new user | ❌ |
| POST | `/auth/login` | Login user | ❌ |
| POST | `/auth/logout` | Logout user | ✅ |
| POST | `/auth/refresh` | Refresh access token | ✅ Cookie |
| GET | `/auth/me` | Get current user info | ✅ |

### User Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/users/me` | Get own profile | ✅ |
| GET | `/users/{id}` | Get user by ID | ✅ |
| GET | `/users/{username}` | Get user by username | ✅ |

### Other Endpoints (Coming Soon)

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/problems` | List coding problems | ✅ |
| GET | `/problems/{slug}` | Get problem details | ✅ |
| POST | `/problems/submit` | Submit solution | ✅ |
| GET | `/submissions` | List user submissions | ✅ |
| GET | `/leaderboard` | Get leaderboard | ✅ |

---

### Problem Publishing Workflow (Admin)

To ensure problem quality and test case correctness, Loco enforces a strict validation workflow:

1. **Create Draft**: Use the V2 Problem Creation Form to define problem metadata, parameters, and test cases.
2. **Implement Reference Solution**: Navigate to the Problem Management page for your draft.
3. **Validate**: Submit a reference solution in any supported language. The system executes this solution against all test cases.
4. **Publish**: Only after a successful validation (100% test cases passed) will the "Publish" button be enabled to make the problem public.

---

### Request/Response Examples

#### Register User

**Request:**
```bash
POST /auth/register
Content-Type: application/json

{
  "email": "john@example.com",
  "username": "johndoe",
  "password": "SecurePass123!"
}
```

**Response:**
```json
{
  "message": "registration successful",
  "user": {
    "id": 1,
    "email": "john@example.com",
    "username": "johndoe",
    "role": "user",
    "email_verified": false,
    "created_at": "2024-11-16T10:00:00Z"
  }
}
```

#### Login

**Request:**
```bash
POST /auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "SecurePass123!"
}
```

**Response:**
- Sets `accessToken` and `refreshToken` HttpOnly cookies
- Returns user object

```json
{
  "message": "login successful",
  "user": {
    "id": 1,
    "email": "john@example.com",
    "username": "johndoe",
    "role": "user",
    "email_verified": false,
    "created_at": "2024-11-16T10:00:00Z"
  }
}
```

#### Get User Profile (Public)

**Request:**
```bash
GET /users/johndoe
Authorization: Cookie (automatically sent)
```

**Response:**
```json
{
  "id": 1,
  "username": "johndoe",
  "role": "user",
  "created_at": "2024-11-16T10:00:00Z"
}
```

Note: Email and verification status are hidden for privacy.

---

## 🔒 Security Features

- ✅ **HttpOnly Cookies** - Prevents XSS attacks from stealing tokens
- ✅ **JWT Token Rotation** - Access + Refresh token pattern
- ✅ **CORS Protection** - Configured allowed origins
- ✅ **Password Hashing** - bcrypt with cost factor 10
- ✅ **Request Validation** - Input sanitization with Zod
- ✅ **SQL Injection Prevention** - Parameterized queries with pgx
- ✅ **Privacy Controls** - Email and sensitive data hidden from public profiles
- ✅ **Secure Cookies** - SameSite and Secure flags in production
- ✅ **Rate Limiting** - (Recommended for production)
- ✅ **HTTPS Only** - (Required for production)

---

## 🧪 Testing

### Backend Tests

```bash
cd backend
go test ./...

# With coverage
go test -cover ./...

# Verbose output
go test -v ./...
```

### Frontend Tests

```bash
cd frontend
npm run test

# With coverage
npm run test:coverage

# Watch mode
npm run test:watch
```

---

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1. **Fork the repository**
2. **Create a feature branch**
   ```bash
   git checkout -b feature/amazing-feature
   ```
3. **Commit your changes**
   ```bash
   git commit -m 'feat: add amazing feature'
   ```
4. **Push to the branch**
   ```bash
   git push origin feature/amazing-feature
   ```
5. **Open a Pull Request**

### Commit Convention

Follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `style:` - Code style changes (formatting, no logic change)
- `refactor:` - Code refactoring
- `test:` - Adding or updating tests
- `chore:` - Maintenance tasks

### Code Style

- **Backend:** Follow [Effective Go](https://go.dev/doc/effective_go) guidelines
- **Frontend:** Use ESLint and Prettier configurations provided

---

## 📝 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

---

## 👥 Authors

- **Prabalesh** - [GitHub](https://github.com/prabalesh)

---

## 🙏 Acknowledgments

- [Go](https://go.dev/) - Backend programming language
- [React](https://react.dev/) - Frontend framework
- [PostgreSQL](https://www.postgresql.org/) - Relational database
- [Tailwind CSS](https://tailwindcss.com/) - CSS framework
- [Zap](https://github.com/uber-go/zap) - Structured logging
- [TanStack Query](https://tanstack.com/query) - Data fetching library
- [Zustand](https://github.com/pmndrs/zustand) - State management
- [Vite](https://vitejs.dev/) - Build tool

---

## 📧 Contact & Support

- **Issues:** [GitHub Issues](https://github.com/prabalesh/loco/issues)
- **Discussions:** [GitHub Discussions](https://github.com/prabalesh/loco/discussions)
- **Email:** your-email@example.com

---

## 🗺️ Roadmap

- [x] User authentication with JWT
- [x] User profiles with gamification
- [x] Problem library and detail views
- [x] Monaco code editor integration
- [x] Multi-language execution via Redis & Piston
- [x] Real-time submission status polling
- [x] Global Leaderboard
- [x] User submission history & statistics
- [x] Achievements & Badges system
- [x] Admin Dashboard for problem management
- [x] Reference Solution Validation System (V2)
- [ ] Contest system and live rankings
- [ ] Contest system and live rankings
- [ ] Discussion forums per problem
- [ ] OAuth integration (Google, GitHub)
- [ ] Advanced code analytics (runtime/memory complexity)

---

**Happy Coding! 🚀**
