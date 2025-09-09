# Blog Application - Domain-Driven Design Practice

A blog application built with Go to practice Domain-Driven Design (DDD) patterns and clean architecture principles.

## Purpose

This project was created for learning purposes to get a better grasp on DDD. The `pkg/ddd` package was written predominantly with Claude Code, and is mostly being used for niceties with Aggregate IDs, event handling, and validation. Everything else was built predominantly by hand.

Auth is also very lightly implemented, as it was not the focus of this learning exercise.

**⚠️ This is not a production application.** It is purely for learning purposes to try and build a small-medium sized application using DDD principles.

## Features

- **Domain-Driven Design Architecture**
  - Aggregates with business logic
  - Domain events and event handling
  - Repository pattern
  - Value objects and domain services

- **Core Functionality**
  - User registration and authentication
  - Create, read, update, and archive blog posts
  - Comment system with threaded discussions
  - Post rating system (upvote/downvote)
  - Role-based permissions (Admin, Author, Commenter)

- **Technical Stack**
  - Go 1.21+
  - SQLite database
  - Session-based authentication with SCS
  - Database migrations with golang-migrate
  - RESTful API design

## Project Structure

```
├── cmd/                    # Application entry points
├── internal/
│   ├── domain/            # Pure business logic (aggregates, events, interfaces)
│   ├── application/       # Use case orchestration (services)
│   ├── infrastructure/    # External system integration
│   │   └── persistence/   # Database implementations and models
│   └── interfaces/        # Inbound adapters (HTTP handlers, middleware)
├── pkg/ddd/              # Reusable DDD infrastructure
└── migrations/           # Database schema migrations
```

## Setup Instructions

### Prerequisites

- Go 1.21 or higher
- SQLite3
- golang-migrate CLI tool

### Install golang-migrate

```bash
# Install with SQLite support
go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Make sure $GOPATH/bin is in your PATH
export PATH=$PATH:$(go env GOPATH)/bin
```

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd blog
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Run database migrations**
   ```bash
   make migrate-up
   ```

4. **Start the application**
   ```bash
   go run cmd/server/main.go
   ```

The server will start on `http://localhost:8080`

## Available Makefile Commands

```bash
# Database migrations
make migrate-create name=migration_name  # Create new migration
make migrate-up                          # Apply all pending migrations
make migrate-down                        # Rollback last migration
make migrate-version                     # Show current migration version

# Development
make help                               # Show all available commands
```

## API Endpoints

### Authentication
- `POST /api/v1/register` - Register new user
- `POST /api/v1/login` - User login
- `POST /api/v1/logout` - User logout

### Users
- `GET /api/v1/users` - Get all users (authenticated)
- `GET /api/v1/users/{id}` - Get user by ID (authenticated)
- `POST /api/v1/users/{id}/description` - Update user description (authenticated)
- `POST /api/v1/users/{id}/password` - Update user password (authenticated)

### Posts
- `GET /api/v1/posts` - Get all posts (authenticated)
- `GET /api/v1/posts/{id}` - Get post by ID (authenticated)
- `POST /api/v1/posts` - Create new post (authenticated)
- `PATCH /api/v1/posts/{id}/title` - Update post title (authenticated)
- `PATCH /api/v1/posts/{id}/content` - Update post content (authenticated)
- `DELETE /api/v1/posts/{id}` - Archive post (authenticated)

### Health Check
- `GET /health` - Service health status

## Learning Focus Areas

This project demonstrates:

1. **Domain-Driven Design**
   - Aggregates with encapsulated business logic
   - Domain events for side effects
   - Repository interfaces in domain layer
   - Value objects for type safety

2. **Clean Architecture**
   - Dependency inversion
   - Separation of concerns
   - Infrastructure independence

3. **Go Best Practices**
   - Package organization
   - Error handling
   - Interface design
   - Testing strategies

## Database Schema

The application uses SQLite with the following main entities:
- **Users** - User accounts with roles and authentication
- **Posts** - Blog posts with authorship and timestamps
- **Comments** - Threaded comments on posts
- **Ratings** - User ratings (upvote/downvote) on posts

## Development Notes

- Authentication is session-based using SCS (Simple Cookie Sessions)
- Users can only modify their own content (basic authorization)
- Posts and comments use soft deletion (archived_at timestamp)
- All timestamps are handled at the database level
- Domain events are dispatched after successful repository operations

## Contributing

This is a learning project, but feel free to explore the code and structure. The DDD patterns and clean architecture principles demonstrated here can be adapted for other Go projects.

## License

This project is for educational purposes only.
