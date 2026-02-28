# Notifier - Payment Notification System

A full-stack payment tracker application with automated email notifications and dashboard view.
Built with Go (Gin), PostgreSQL, and React.

## Features
### Backend
- Automated daily, weekly, monthly, and yearly payment notifications using scheduled cron jobs
- Email reminders sent via SMTP
- PostgreSQL with sqlc for compile-time type-safe query generation
- REST API with JWT authentication
- Request logging for monitoring and debugging
- Clean architecture focused on scalability and maintainability

### Frontend
- Account creation and login via JWT
- Create, update, and delete payments
- Dashboard with a visual overview of all payments

---
## Prerequisites

- Go 1.20+
- Bun 1.0+
- PostgreSQL instance
- SMTP credentials for email sending

---
## Configuration
Create a `.env` file in the root directory:
    
    POSTGRES_HOST=localhost
    POSTGRES_PORT="5432"
    POSTGRES_USER=user
    POSTGRES_PASSWORD=password
    POSTGRES_NAME=notifier
    SMTP_HOST=smtp.example.com
    SMTP_PORT=587
    SMTP_USERNAME=mail@example.com
    SMTP_PASSWORD=your_app_password
    JWT_KEY=your_jwt_key

---
## Installation

#### 1. Clone the repository
```
git clone --recurse-submodules https://github.com/rober0xf/notifier.git
cd notifier
```

#### 2. Backend
```
go mod tidy
```

#### 3. Frontend setup
```
cd frontend
bun install
```

#### 4. Run
```
air
```

#### 5. Test an endpoint
```
curl -X POST http://localhost:3000/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "username": "username",
    "password": "securepassword123"
}'
```

---

## Testing
The project has two layers of tests.

**Unit tests** cover the business logic using mocks for external dependencies like the database and email sender.

**Integration tests** cover the full HTTP request lifecycle against a real Postgres database, including handlers, use cases, and repository layer.

```
# unit tests
go test ./internal/...

# integration tests
go test ./test/integration/...

# all tests
go test ./...
```
---

## Contributing
Fork the repository and submit a pull request.

1. Fork the project
2. Create your feature branch (`git checkout -b feature/SomeFeature`)
3. Commit your changes (`git commit -m 'Added some Feature'`)
4. Push to the branch (`git push origin feature/SomeFeature`)
5. Open a Pull Request
