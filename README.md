# Notifier - Payment Notification System

<div align="center">
  <img src="https://raw.githubusercontent.com/rober0xf/notifier-front/master/public/bimage.png" alt="Notifier Dashboard" width="900">
</div>

A full-stack payment tracker application with automated email notifications and dashboard view. Built with Golang Gin Framework, SQLite database, and Vuejs frontend for managing payments and users.

## Features

### Backend
- **Automated Scheduling** - Daily, weekly, and monthly payment notifications via cron job.
- **Email Notifications** - Sends personalized payment reminders to users.
- **SQLite** - Lightweight, embedded database with no external dependencies.
- **REST API** - Clean and secure APIs.
- **Logging** - Some logs for monitoring and debugging.

### Frontend
- **User Authentication** - Secure account creation and login through JWT.
- **Payment Management** - Create, update, and delete payments.
- **Dashboard** - Visual view of all payments.

---
## Prerequisites

- Go 1.20 or higher
- Node.js 20+ and npm (for frontend)
- SMTP server credentials for sending emails

---
## Configuration
**Configure the env variables in a file `.env`:**
    
    DB_HOST=localhost # ignore for sqlite
    DB_PORT=5432 # ignore for sqlite
    DB_USER=user # ignore for sqlite
    DB_PASSWORD=password # ignore for sqlite
    DB_NAME=notifier # ignore for sqlite
    SMTP_HOST=smtp.example.com
    SMTP_PORT=587
    SMTP_USER=mail@example.com
    SMTP_PASSWORD=your_app_password
    JWT_KEY=your_jwt_key

---
## Installation

#### 1. Clone the repository
```
# clone with frontend submodule
git clone --recurse-submodules https://github.com/rober0xf/notifier.git
cd notifier
```

#### 2. Backend setup
```
# install Go dependencies
go mod tidy

# Initialize the database (creates tables automatically on first run)
```
#### 3. Frontend setup
```
cd frontend
pnpm install
```
#### 4. Environment configuration
```
# create a .env file in the root directory
cp .env.example .env
```
#### 5. Run the application
```
air
```
#### 6. Test an endpoint
```
curl -X POST http://localhost:3000/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "username": "username",
    "password": "securepassword123"
}'
```
