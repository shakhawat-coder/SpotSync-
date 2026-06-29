# SpotSync - Smart Parking Management System

SpotSync is a robust backend API built with Go to manage parking zones and reservations in real-time. It features secure JWT-based authentication, role-based access control, and atomic concurrency control to prevent overbooking.

## 🚀 Live Demo

**Production API:** `https://spotsync-d7v2.onrender.com`  
**Local Development:** `http://localhost:8080/api/v1`

---

## 🛠️ Tech Stack

- **Language:** Go 1.22+
- **Framework:** Echo v5
- **ORM:** GORM
- **Database:** PostgreSQL
- **Authentication:** JWT & bcrypt

---

## ⚙️ Setup Instructions

### 1. Clone Repository
```bash
git clone https://github.com/shakhawat-coder/SpotSync-.git
cd SpotSync-
```

### 2. Install Dependencies
```bash
go mod download
go mod tidy
```

### 3. Configure `.env`
```env
DATABASE_URL=postgresql://user:password@localhost:5432/spotsync?sslmode=disable
JWT_SECRET=your-32-character-secret-key-here
ENVIRONMENT=development
PORT=8080
```

### 4. Run Application
```bash
go run ./cmd/main.go
```

**Expected Output:**
```
Migration completed successfully
🚀 Server starting on port 8080
```

### 5. Verify Setup
```bash
curl http://localhost:8080/health
```

**Expected Response:**
```json
{"status":"ok","environment":"development"}
```

---

## 🏗️ Architecture

SpotSync follows Clean Architecture with clear layer separation:

- **Handler Layer:** HTTP routing and request handling
- **Service Layer:** Business logic and validation
- **Repository Layer:** Database operations with GORM
- **Models Layer:** Database entities

---

## 📑 API Endpoints

### Authentication
| Method | Endpoint | Description | Auth |
|:------:|----------|-------------|------|
| POST | `/api/v1/auth/register` | Register new user | Public |
| POST | `/api/v1/auth/login` | Login & get JWT | Public |

### Parking Zones
| Method | Endpoint | Description | Auth |
|:------:|----------|-------------|------|
| GET | `/api/v1/zones` | List all zones | Public |
| GET | `/api/v1/zones/:id` | Get zone details | Public |
| POST | `/api/v1/zones` | Create new zone | Admin |

### Reservations
| Method | Endpoint | Description | Auth |
|:------:|----------|-------------|------|
| POST | `/api/v1/reservations` | Create reservation | JWT |
| GET | `/api/v1/reservations/my-reservations` | View own reservations | JWT |
| GET | `/api/v1/reservations` | View all reservations | Admin |
| DELETE | `/api/v1/reservations/:id` | Cancel reservation | JWT |

---

## 🔒 Security Features

- ✅ Password Hashing (bcrypt cost 10-12)
- ✅ JWT Authentication (24-hour expiry)
- ✅ Role-Based Access Control
- ✅ SQL Injection Prevention (GORM)
- ✅ Atomic Reservations (Database locking)
- ✅ Input Validation (Struct-based)
- ✅ HTTPS/SSL (Render auto-enabled)

---

## JWT Token Errors

```bash
# Token expires in 24 hours
# Generate new token: POST /api/v1/auth/login
# Include in all requests: Authorization: Bearer <token>
```

## 📂 Project Structure

```
SpotSync/
├── cmd/main.go
├── config/
├── dto/
├── errors/
├── handler/
├── middleware/
├── migrations/
├── models/
├── repository/
├── routes/
├── service/
├── go.mod
├── go.sum
└── README.md
```

---
