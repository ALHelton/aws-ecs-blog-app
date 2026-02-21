# Go App Deployment to AWS ECS — Full Walkthrough

## Overview

This guide documents the full process of scaffolding a Go app, containerizing it with Docker, and deploying it to AWS ECS using Fargate.

**Overall flow:** your code → Docker image → ECR (storage) → ECS Task Definition (recipe) → ECS Service (execution) → Security Group (network access)

---

## Step 1: Scaffold the App with go-blueprint

**What:** Generates a basic Go HTTP server with the Chi router, a Hello World endpoint, and all the boilerplate (go.mod, Makefile, etc.)

**Why:** Instead of writing a Go app from scratch, go-blueprint generates a working, structured app quickly — the whole point of this exercise.

```bash
go install github.com/melkeydev/go-blueprint@latest

export PATH=$PATH:$(go env GOPATH)/bin

GOTOOLCHAIN=local go-blueprint create

# Choose: Chi framework, no database

cd <your-project-name>

go run cmd/api/main.go

# Verify at http://localhost:8080/
```

---

## Step 2: Build a Docker Image

**What:** A two-stage Dockerfile — the first stage compiles the Go binary, the second stage puts just the binary into a lightweight Alpine image.

**Why:** ECS can't take your source code and run it directly. It needs a self-contained image that has everything required to run your app. Docker packages it into that portable unit.

**Dockerfile:**
# Step 1: Build the binary
```dockerfile
# Step 1: Build the binary
FROM golang:1.26.0-alpine AS builder

# Install git or other necessary dependencies
RUN apk add --no-cache git

# Set the working directory
WORKDIR /aws-ecs-blog-app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary for arm64
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o main cmd/api/main.go

# Step 2: Final lightweight image
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from the builder
COPY --from=builder /aws-ecs-blog-app/main .

# Expose port
EXPOSE 8080

CMD ["./main"]
```

# Step 2: Final lightweight image
```dockerfile
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /aws-ecs-blog-app/main .

EXPOSE 8080

CMD ["./main"]
```

**Commands:**
```bash
docker build -t deploy-test .

docker run -p 8080:8080 -e PORT=8080 deploy-test

# Verify at http://localhost:8080/
```

---

## Step 3: Push the Image to ECR

**What:** Creates a repository in AWS's container registry and pushes your Docker image there.

**Why:** ECS needs to pull your image from somewhere. ECR (Elastic Container Registry) is AWS's own Docker registry — ECS has direct access to it.

```bash
# Create the repository
aws ecr create-repository --repository-name deploy-test --region us-east-1

# Authenticate Docker with ECR
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin <account-id>.dkr.ecr.us-east-1.amazonaws.com

# Tag and push the image
docker tag deploy-test <account-id>.dkr.ecr.us-east-1.amazonaws.com/deploy-test
docker push <account-id>.dkr.ecr.us-east-1.amazonaws.com/deploy-test
```

> Note: `--username AWS` is literal — it's not your IAM username, it's how ECR authentication works.

---

## Step 4: Create an ECS Cluster

**What:** A logical grouping in AWS that holds your services and tasks.

**Why:** Think of it as the "environment" your app lives in. It doesn't do anything by itself — it's the container that organizes everything.

```
AWS Console: ECS → Clusters → Create Cluster

- Name: deploy-test-cluster
- Infrastructure: Fargate (AWS manages the servers)
```

---

## Step 5: Create a Task Definition

**What:** A configuration that tells ECS how to run your container — which image to use, how much CPU/memory, what port, and what environment variables.

**Why:** This is the "recipe" for your app. ECS references it every time it needs to spin up an instance of your container.

```
AWS Console: ECS → Task Definitions → Create new task definition

- Name: deploy-test-task
- Launch type: Fargate
- OS: Linux/ARM64
- CPU: .25 vCPU, Memory: .5 GB
- Container name: deploy-test
- Image URI: <account-id>.dkr.ecr.us-east-1.amazonaws.com/deploy-test:latest
- Container port: 8080
- Environment variable: PORT=8080
```

---

## Step 6: Create a Service

**What:** Tells ECS to actually run your task definition and keep it running.

**Why:** The task definition alone is just a blueprint — the service is what launches it and maintains it. If your container crashes, the service restarts it. It's the thing that keeps your app alive.

```
AWS Console: ECS → Clusters → deploy-test-cluster → Services → Create

- Launch type: Fargate
- Task definition: deploy-test-task
- Service name: deploy-test-service
- Desired tasks: 1
```

---

## Step 7: Open Port 8080 in the Security Group

**What:** Adds an inbound rule to the task's security group allowing TCP traffic on port 8080.

**Why:** AWS blocks all inbound traffic by default. Without this rule, requests to your app are silently blocked at the network level before they ever reach your container. The security group acts as a firewall — you have to explicitly open the port your app listens on.

```
AWS Console: ECS → Clusters → deploy-test-cluster → Tasks → <your task> → Network → Security Group → Edit Inbound Rules

- Type: Custom TCP
- Port: 8080
- Source: 0.0.0.0/0
```

> Note: In production, restrict source to known IPs rather than 0.0.0.0/0.

Your app is now accessible at `http://<public-ip>:8080/`

---

## Tearing Down (to avoid charges)

Tear down in reverse order:

**1. Delete the ECS Service**
```
ECS → Clusters → deploy-test-cluster → Services → deploy-test-service → Delete

Use "Force delete" to have ECS scale down to zero automatically
```

**2. Delete the ECS Cluster**
```
ECS → Clusters → deploy-test-cluster → Delete cluster
```

**3. Delete the Task Definition**
```
ECS → Task Definitions → deploy-test-task → select all revisions → Deregister → Delete
```

**4. Delete the ECR Repository**
```bash
aws ecr delete-repository --repository-name deploy-test --region us-east-1 --force
```

**5. Delete the Security Group**
```
EC2 → Security Groups → find the one created for your task → Delete
```

**6. Delete any CloudFormation stacks**
```
CloudFormation → Stacks → delete any deploy-test related stacks
```

**7. Verify in Billing**
```
AWS Console → Billing → Bills → confirm nothing is still running
```

> Fargate charges by the second — once the service is deleted you stop accruing costs immediately.
