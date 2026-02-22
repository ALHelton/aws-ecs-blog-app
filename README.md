# aws-ecs-blog-app

A Go blog API with two models (BlogPost and Comment) in a one-to-many relationship. Containerized with Docker and deployed to AWS ECS via Fargate with Amazon RDS (PostgreSQL).

## Tech Stack

- Go with Chi router
- PostgreSQL (Amazon RDS)
- Docker
- AWS ECS (Fargate)

## Models

- BlogPost: id, title, content, created_at, updated_at
- Comment: id, blog_post_id, content, created_at, updated_at

## Endpoints

### Blog Posts
- GET /blog_posts
- GET /blog_posts/{id}
- POST /blog_posts
- PUT /blog_posts/{id}
- DELETE /blog_posts/{id}

### Comments
- GET /comments?blog_post_id={id}
- GET /comments/{id}
- POST /comments
- PUT /comments/{id}
- DELETE /comments/{id}

## Running Locally

Copy the example env file and fill in your values:

```bash
cp .env.example .env
```

Run with Docker Compose:

```bash
docker-compose up --build
```

API Address: `http://localhost:8080`.

## Running Tests

```bash
go test ./...
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| PORT | Port the server listens on |
| DB_HOST | Database host |
| DB_PORT | Database port |
| DB_USER | Database user |
| DB_PASSWORD | Database password |
| DB_NAME | Database name |

## Deployment

See [aws-ecs-deployment-guide.md](aws-ecs-deployment-guide.md) for full deployment instructions.

```
						┌─────────────────────────────────────────┐
						│           AWS (us-east-1)               │
						│                                         │
						│   ┌─────────────────────────────────┐   │
						│   │        Default VPC              │   │
						│   │                                 │   │
Internet ──── port 8080 ──► │   ┌──────────────────────┐      │   │
						│   │   │   ECS Fargate Task   │      │   │
						│   │   │                      │      │   │
						│   │   │   Go Blog App        │      │   │
						│   │   │   (port 8080)        │      │   │
						│   │   └──────────┬───────────┘      │   │
						│   │              │                  │   │
						│   │         port 5432               │   │
						│   │         (internal)              │   │
						│   │              │                  │   │
						│   │   ┌──────────▼───────────┐      │   │
						│   │   │   Amazon RDS         │      │   │
						│   │   │   PostgreSQL         │      │   │
						│   │   │   (private)          │      │   │
						│   │   └──────────────────────┘      │   │
						│   │                                 │   │
						│   └─────────────────────────────────┘   │
						│                                         │
						└─────────────────────────────────────────┘
