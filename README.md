# Portfolio Website Backend

![Go](https://img.shields.io/badge/go-1.25-blue.svg)
![Gin](https://img.shields.io/badge/Gin-1.9-green.svg)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-blue.svg)
![GORM](https://img.shields.io/badge/GORM-V2-orange.svg)

A modernized multi-tenant backend utilizing the Go Gin web framework and Clean Architecture. This API manages user portfolios including projects, experiences, and categories with JWT authentication and rigorous test coverage.

## 🚀 Features (Implemented So Far)

- **User Management**: Auth logic, JWT issuance, profile updates, and admin CLI logic (via secure handler overrides).
- **Project Management**: CRUD operations for portfolio projects with visibility toggling.
- **Project Categories**: Classify projects into distinct groups utilizing proper foreign key relations natively.
- **Experience/Education**: Track your active work and educational history.
- **Local & GitHub Projects**:
   - `github` types: On creation, the backend automatically contacts the GitHub API fetching data, stars, watchers, and processes markdown README items.
   - 24-hour Expiration logic efficiently refetches stale GitHub data dynamically.
- **Visibility Control**: (`is_visible`) determines whether records appear during unauthenticated public fetches.
- **JWT Authentication**: Middleware securely validates tokens across admin and private endpoints.
- **Swagger Documentation**: Declarative annotations generating a dynamic OpenAPI Swagger Dashboard.
- **Health Check**: Database connectivity monitoring endpoint via `/healthz`.

## 🔧 API Endpoints

### Authentication
- `POST /api/v1/login` - Login and get JWT token
- `GET /api/v1/profile` - Get user profile (requires auth)
- `PUT /api/v1/profile` - Update user profile (requires auth)

### Public Portfolio Data
- `GET /api/v1/public-data/:user_id` - Fetches combined payload representing user portfolio data.

### Projects
- `GET /api/v1/projects` - List all projects including hidden ones (requires auth)
- `POST /api/v1/projects` - Create a new project (requires auth)
- `GET /api/v1/projects/:id` - Get a specific project including if hidden (requires auth)
- `PUT /api/v1/projects/:id` - Update a project (requires auth)
- `PATCH /api/v1/projects/:id/visibility` - Update project visibility (requires auth)
- `DELETE /api/v1/projects/:id` - Delete a project (requires auth)
- `GET /api/v1/projects/public/:user_id` - List strictly visible projects locally mapped (public access)

### Project Categories
- `GET /api/v1/project-categories` - List all categories including hidden ones (requires auth)
- `POST /api/v1/project-categories` - Create a new category (requires auth)
- `GET /api/v1/project-categories/:id` - Get a category by ID (requires auth)
- `PUT /api/v1/project-categories/:id` - Update a category (requires auth)
- `PATCH /api/v1/project-categories/:id/visibility` - Update visibility (requires auth)
- `DELETE /api/v1/project-categories/:id` - Delete a category (requires auth)
- `GET /api/v1/project-categories/public/:user_id` - List visible categories only (public access)

### Experiences
- `GET /api/v1/experiences` - List all experiences including hidden ones (requires auth)
- `POST /api/v1/experiences` - Create a new experience entry (requires auth)
- `GET /api/v1/experiences/:id` - Get a specific experience entry by ID (requires auth)
- `PUT /api/v1/experiences/:id` - Update an experience entry (requires auth)
- `PATCH /api/v1/experiences/:id/visibility` - Update experience visibility (requires auth)
- `DELETE /api/v1/experiences/:id` - Delete an experience entry (requires auth)
- `GET /api/v1/experiences/public/:user_id` - List visible experiences only (public access)

### System
- `GET /healthz` - Health check endpoint verifying DB availability
- `GET /swagger/*any` - Interactive API definition explorer

## 🔐 Authentication
The API utilizes JWT tokens for active validation across internal endpoints.
- Add the token mapped to Authorization headers as: `Bearer <token>`

## ⚙️ Setup and Installation

### Prerequisites
- Go 1.25
- PostgreSQL database
- Git

### Local Development Setup

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd portfolio-website-backend-go
   ```

2. **Set up environment variables:**
   Create a `.env` file manually referencing properties matching local runtime needs (see `env_example.txt`):
   ```env
   DATABASE_URL=postgresql://<db_user>:<db_password>@<db_host>:<db_port>/<db_name>?sslmode=require&channel_binding=require
   API_PREFIX=/api/v1
   DEBUG=True
   MAX_DB_RETRIES=3
   RETRY_BACKOFF=0.5
   GITHUB_TOKEN=<github_token>
   ADMIN_API_KEY=<admin_api_key>
   JWT_SECRET_KEY=<jwt_secret_key>
   ACCESS_TOKEN_EXPIRE_MINUTES=120
   CORS_ORIGINS=http://[IP_ADDRESS]
   PORT=8000
   ```

3. **Install Dependencies and Format:**
   ```bash
   go mod download
   go fmt ./...
   ```

4. **Regenerate Swagger API Docs (If updating Handler Annotations)**
   ```bash
   go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/server/main.go --parseDependency --parseInternal
   ```

5. **Start the application:**
   ```bash
   go run cmd/server/main.go
   ```
   *Note: GORM `AutoMigrate` runs exclusively on server boot establishing any unrecorded database schemas instantaneously.*

## 🧪 Testing

This project leverages native local Go testing infrastructure supplemented by robust Mocking (`testify/mock`) and test containers simulating physical connections.

To trigger validations across repositories, services, and handlers:
```bash
go clean -testcache
go test -p 1 ./... -v -failfast
```
*(Note: Full integration tests for the Service and Repository layers require a live Postgres instance. By default, the testing suite falls back to `postgresql://postgres:postgres@localhost:5432/portfolio_test?sslmode=disable`. You can explicitly override this by setting `TEST_DATABASE_URL` in your `.env` file.)*