# Portfolio Website Backend

![Go](https://img.shields.io/badge/go-1.25-blue.svg)
![Gin](https://img.shields.io/badge/Gin-1.9-green.svg)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-blue.svg)
![GORM](https://img.shields.io/badge/GORM-V2-orange.svg)

A modernized multi-tenant backend utilizing the Go Gin web framework and Clean Architecture. This API manages user portfolios including projects, experiences, skills, and reviews with JWT authentication and rigorous test coverage across all domains.

## 🚀 Features

- **User Management**: Auth logic, JWT issuance, profile updates, and a secure admin CLI tool (`scripts/manage_user.go`) for user provisioning.
- **Project Management**: Full CRUD for portfolio projects with optional GitHub API sync for stats and README content.
- **Project Categories**: Classify projects into distinct groups with proper foreign key relations.
- **Experience / Education**: Track active work history and educational records with date-range support.
- **Skills & Skill Groups**: Hierarchical skill management — skills are organized under groups, both with independent visibility control.
- **Reviews / Testimonials**: Store and manage client/colleague reviews with visibility toggling.
- **Local & GitHub Projects**:
  - `github` types: On creation, the backend automatically contacts the GitHub API fetching stars, watchers, and processes markdown README content.
  - 24-hour expiration logic efficiently refetches stale GitHub data.
- **Visibility Control**: `is_visible` determines whether records appear during unauthenticated public fetches.
- **Multi-tenancy**: All domain records are scoped to an authenticated user via `user_id`, enabling future multi-user support.
- **JWT Authentication**: Middleware securely validates tokens across all admin and private endpoints.
- **Admin Auth**: Separate `X-Admin-Api-Key` header auth for privileged provisioning endpoints.
- **Swagger Documentation**: Declarative annotations generating a dynamic OpenAPI Swagger Dashboard.
- **AI Chatbot & Dynamic RAG**: Built-in websocket chat endpoint powered by Google's Gemini LLM. Uses `pgvector` for semantic search across portfolio records, providing fully dynamic context integration without rigid keyword restrictions.
- **Multitenant Vector Store**: All RAG embeddings are strictly segregated by `user_id`, ensuring the chatbot securely synthesizes data from the active portfolio owner.
- **Health Check**: Database connectivity monitoring via `/healthz`.
- **CI/CD Pipeline**: GitHub Actions workflow runs the full test suite against a real `pgvector/pgvector:pg15` service container on every pull request.

## 🔧 API Endpoints

### Authentication
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `POST` | `/api/v1/login` | — | Login and receive a JWT token |
| `GET` | `/api/v1/profile` | 🔒 JWT | Get authenticated user's profile |
| `PUT` | `/api/v1/profile` | 🔒 JWT | Update authenticated user's profile |

### Public Portfolio Data
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `GET` | `/api/v1/public-data/:user_id` | — | Combined payload of all visible portfolio data for a user |

### Projects
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `GET` | `/api/v1/projects` | 🔒 JWT | List all projects (incl. hidden) |
| `POST` | `/api/v1/projects` | 🔒 JWT | Create a new project |
| `GET` | `/api/v1/projects/:id` | 🔒 JWT | Get a specific project by ID |
| `PUT` | `/api/v1/projects/:id` | 🔒 JWT | Update a project |
| `PATCH` | `/api/v1/projects/:id/visibility` | 🔒 JWT | Toggle project visibility |
| `DELETE` | `/api/v1/projects/:id` | 🔒 JWT | Delete a project |
| `GET` | `/api/v1/projects/public/:user_id` | — | List visible projects (public) |

### Project Categories
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `GET` | `/api/v1/project-categories` | 🔒 JWT | List all categories |
| `POST` | `/api/v1/project-categories` | 🔒 JWT | Create a new category |
| `GET` | `/api/v1/project-categories/:id` | 🔒 JWT | Get a category by ID |
| `PUT` | `/api/v1/project-categories/:id` | 🔒 JWT | Update a category |
| `PATCH` | `/api/v1/project-categories/:id/visibility` | 🔒 JWT | Toggle category visibility |
| `DELETE` | `/api/v1/project-categories/:id` | 🔒 JWT | Delete a category |
| `GET` | `/api/v1/project-categories/public/:user_id` | — | List visible categories (public) |

### Experiences
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `GET` | `/api/v1/experiences` | 🔒 JWT | List all experiences (incl. hidden) |
| `POST` | `/api/v1/experiences` | 🔒 JWT | Create a new experience entry |
| `GET` | `/api/v1/experiences/:id` | 🔒 JWT | Get a specific experience by ID |
| `PUT` | `/api/v1/experiences/:id` | 🔒 JWT | Update an experience entry |
| `PATCH` | `/api/v1/experiences/:id/visibility` | 🔒 JWT | Toggle experience visibility |
| `DELETE` | `/api/v1/experiences/:id` | 🔒 JWT | Delete an experience entry |
| `GET` | `/api/v1/experiences/public/:user_id` | — | List visible experiences (public) |

### Reviews
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `GET` | `/api/v1/reviews` | 🔒 JWT | List all reviews (incl. hidden) |
| `POST` | `/api/v1/reviews` | 🔒 JWT | Create a new review |
| `GET` | `/api/v1/reviews/:review_id` | 🔒 JWT | Get a specific review by ID |
| `PUT` | `/api/v1/reviews/:review_id` | 🔒 JWT | Update a review |
| `PATCH` | `/api/v1/reviews/:review_id/visibility` | 🔒 JWT | Toggle review visibility |
| `DELETE` | `/api/v1/reviews/:review_id` | 🔒 JWT | Delete a review |
| `GET` | `/api/v1/reviews/public/:user_id` | — | List visible reviews (public) |

### Skills & Skill Groups
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `GET` | `/api/v1/skills/groups` | 🔒 JWT | List all skill groups (incl. hidden) |
| `POST` | `/api/v1/skills/groups` | 🔒 JWT | Create a new skill group |
| `GET` | `/api/v1/skills/groups/:group_id` | 🔒 JWT | Get a skill group by ID |
| `PUT` | `/api/v1/skills/groups/:group_id` | 🔒 JWT | Update a skill group |
| `PATCH` | `/api/v1/skills/groups/:group_id/visibility` | 🔒 JWT | Toggle skill group visibility |
| `DELETE` | `/api/v1/skills/groups/:group_id` | 🔒 JWT | Delete a skill group |
| `GET` | `/api/v1/skills` | 🔒 JWT | List all individual skills (incl. hidden) |
| `POST` | `/api/v1/skills` | 🔒 JWT | Create an individual skill |
| `GET` | `/api/v1/skills/:skill_id` | 🔒 JWT | Get a skill by ID |
| `PUT` | `/api/v1/skills/:skill_id` | 🔒 JWT | Update a skill |
| `PATCH` | `/api/v1/skills/:skill_id/visibility` | 🔒 JWT | Toggle skill visibility |
| `DELETE` | `/api/v1/skills/:skill_id` | 🔒 JWT | Delete a skill |
| `GET` | `/api/v1/skills/public/:user_id` | — | List visible skill groups with their skills (public) |

### Chatbot & RAG
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `POST` | `/api/v1/chatbot/sync` | 🔒 JWT | Triggers a sync of user portfolio data into vector embeddings |
| `GET` | `/api/v1/chatbot/sessions` | 🔒 JWT | List all historical chat sessions (pagination supported) |
| `GET` | `/api/v1/chatbot/sessions/:session_id/messages` | 🔒 JWT | Get messages for a specific session |
| `WS`   | `/api/v1/chatbot/ws/chat` | — | Realtime WebSocket chat stream (requires `session_id` and `user_id` query params) |

### Admin (API-Key Protected)
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `POST` | `/api/v1/admin/users` | 🔑 API Key | Create a new user account |
| `GET` | `/api/v1/admin/users` | 🔑 API Key | List all users |
| `GET` | `/api/v1/admin/users/:id` | 🔑 API Key | Get a specific user by ID |

### System
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `GET` | `/healthz` | — | Health check — verifies DB availability |
| `GET` | `/swagger/*any` | — | Interactive Swagger API explorer |

## 🔐 Authentication

The API uses two auth mechanisms:

- **JWT (Bearer Token)** — For all standard user-facing endpoints. Include in the `Authorization` header:
  ```
  Authorization: Bearer <jwt_token>
  ```
- **Admin API Key** — For privileged provisioning endpoints (creating/listing users). Include in the `X-Admin-Api-Key` header:
  ```
  X-Admin-Api-Key: <admin_api_key>
  ```

## ⚙️ Setup and Installation

### Prerequisites
- Go 1.25+
- PostgreSQL 15 with the `pgvector` extension
- Git

### Local Development Setup

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd portfolio-website-backend-go
   ```

2. **Set up environment variables:**
   Create a `.env` file referencing `env_example.txt`:
   ```env
   DATABASE_URL=postgresql://<db_user>:<db_password>@<db_host>:<db_port>/<db_name>?sslmode=require
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

3. **Install dependencies:**
   ```bash
   go mod download
   go fmt ./...
   ```

4. **Regenerate Swagger API docs** (only needed after updating handler annotations):
   ```bash
   go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/server/main.go --parseDependency --parseInternal
   ```

5. **Start the server:**
   ```bash
   go run ./cmd/server/main.go
   ```
   > GORM `AutoMigrate` runs on server boot and applies any missing schema changes automatically.

### User Management CLI

Use `scripts/manage_user.go` to provision users while the server is running. Requires `ADMIN_API_KEY` to be set in your environment.

```bash
# Create a user
go run scripts/manage_user.go create --email user@example.com --password secret --username myuser

# List all users
go run scripts/manage_user.go list

# Get a user by ID
go run scripts/manage_user.go get --id <uuid>
```

## 🧪 Testing

The project has a comprehensive test suite across all 7 domains (Users, Experiences, Projects, Project Categories, Reviews, Skills, Chatbot), each with:

- **Unit tests** — Service layer tested with `testify/mock`.
- **Integration tests** — Repository layer tested against a live Postgres instance.
- **Handler tests** — HTTP handlers tested end-to-end via `httptest`.

### Running Tests Locally

```bash
go clean -testcache
go test -p 1 ./tests/... -v -failfast
```

> Integration tests require a live Postgres instance. The suite falls back to `postgresql://postgres:postgres@localhost:5432/portfolio_test?sslmode=disable` by default. Override with `TEST_DATABASE_URL` in your `.env` file.

### CI/CD

A GitHub Actions workflow (`.github/workflows/test.yml`) runs the full test suite automatically on every pull request to `main`. It spins up a `pgvector/pgvector:pg15` service container for real database integration tests.

To simulate CI locally using [`act`](https://github.com/nektos/act):
```bash
act -j test
```

## 🐳 Docker

Build and run the application using Docker:

```bash
# Build the image
docker build -t portfolio-backend .

# Run with docker-compose (includes Postgres)
docker-compose up
```