# AWS ECS Deployment Guide — aws-ecs-blog-app

A guide for deploying a Go app to AWS ECS with RDS.

---

## Step 1: IAM Setup

```
IAM → Users → Create User
- Attach: AdministratorAccess
- Create access key for CLI
- Set up console password and MFA
```

```bash
aws configure
# Enter: Access Key ID, Secret Access Key, region (us-east-1), output format (leave blank)
```

---

## Step 2: Create Security Groups FIRST

**RDS Security Group:**
```
EC2 → Security Groups → Create security group
- Name: rds-sg
- Description: RDS security group for aws-ecs-blog-app
- VPC: default
- Inbound rule: PostgreSQL, Port 5432, Source: 0.0.0.0/0 (tighten to ECS sg after setup)
```

**ECS Security Group (default):**
```
EC2 → Security Groups → default → Edit inbound rules
- Add: Custom TCP, Port 8080, Source: 0.0.0.0/0
```

---

## Step 3: RDS — Create Database

```
RDS → Create database
- Engine: PostgreSQL (latest version)
- Template: Free tier
- Identifier: aws-ecs-blog-app-db
- Master username: postgres
- Password: your choice (save it immediately)
- Instance: (smallest default from easy setup option)
- VPC: default
- Public access: Yes (temporarily, to create the blog database manually)
- Security group: rds-sg
- Initial database name: blog
```

> Note: Setting the initial database name here avoids having to create it manually after the fact.

After RDS is created, connect and verify the database exists:

```bash
psql "host=<rds-endpoint> user=postgres dbname=postgres sslmode=require"
```

```sql
\l          -- list databases, verify blog exists
\q          -- exit
```

Then set **Public access: No** once you're done.

---

## Step 4: ECR — Create Repository and Push Image

```bash
# Create the repository
aws ecr create-repository --repository-name aws-ecs-blog-app --region us-east-1

# Authenticate Docker with ECR
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin <account-id>.dkr.ecr.us-east-1.amazonaws.com

# Build, tag, and push
docker build -t aws-ecs-blog-app .
docker tag aws-ecs-blog-app <account-id>.dkr.ecr.us-east-1.amazonaws.com/aws-ecs-blog-app:latest
docker push <account-id>.dkr.ecr.us-east-1.amazonaws.com/aws-ecs-blog-app:latest
```

> Note: `--username AWS` is literal — it's not your IAM username. That's how ECR authentication works.

---

## Step 5: ECS Cluster

```
ECS → Clusters → Create Cluster
- Name: aws-ecs-blog-app-cluster
- Infrastructure: Fargate only
```

---

## Step 6: ECS Task Definition

```
ECS → Task Definitions → Create new task definition
- Name: aws-ecs-blog-app-task
- Launch type: Fargate
- OS: Linux/ARM64
- CPU: .5 vCPU, Memory: 1 GB
```

Container configuration:
```
- Name: app
- Image URI: <account-id>.dkr.ecr.us-east-1.amazonaws.com/aws-ecs-blog-app:latest
- Port: 8080
- Environment variables:
    PORT=8080
    DB_HOST=<rds-endpoint>
    DB_PORT=5432
    DB_USER=postgres
    DB_PASSWORD=<your-password>
    DB_NAME=blog
```

> Note: Use `sslmode=require` in your Go database connection string for RDS.

---

## Step 7: ECS Service

```
ECS → Clusters → aws-ecs-blog-app-cluster → Services → Create
- Launch type: Fargate
- Task definition: aws-ecs-blog-app-task
- Service name: aws-ecs-blog-app-service
- Desired tasks: 1
- VPC: default
- Subnets: default
- Auto-assign public IP: Enabled
- Security group: default
```

---

## Step 8: Tighten RDS Security Group

After ECS service is created, get the ECS security group ID:
```
ECS → Clusters → aws-ecs-blog-app-cluster → Services → aws-ecs-blog-app-service → Configuration and networking tab
```

Then update `rds-sg` to only allow traffic from ECS:
```
EC2 → Security Groups → rds-sg → Edit inbound rules
- Change source from 0.0.0.0/0 to the default security group ID (sg-xxxxxxxx)
```

---

## Updating the Deployment

When you make code changes:

```bash
# Rebuild and push new image
docker build -t aws-ecs-blog-app .
docker tag aws-ecs-blog-app <account-id>.dkr.ecr.us-east-1.amazonaws.com/aws-ecs-blog-app:latest
docker push <account-id>.dkr.ecr.us-east-1.amazonaws.com/aws-ecs-blog-app:latest

# Create new task definition revision
ECS → Task Definitions → aws-ecs-blog-app-task → Create new revision

# Update service
ECS → Clusters → aws-ecs-blog-app-cluster → Services → aws-ecs-blog-app-service → Update service
- Select latest task definition revision
- Force new deployment: checked
```

---

## Tearing Down (to avoid charges)

Tear down in reverse order:

**1. Delete ECS Service**
```
ECS → Clusters → aws-ecs-blog-app-cluster → Services → aws-ecs-blog-app-service → Delete
Use "Force delete"
```

**2. Delete ECS Cluster**
```
ECS → Clusters → aws-ecs-blog-app-cluster → Delete cluster
```

**3. Delete Task Definition**
```
ECS → Task Definitions → aws-ecs-blog-app-task → select all revisions → Deregister → Delete
```

**4. Delete ECR Repository**
```bash
aws ecr delete-repository --repository-name aws-ecs-blog-app --region us-east-1 --force
```

**5. Delete RDS Instance**
```
RDS → Databases → aws-ecs-blog-app-db → Actions → Delete
- Uncheck "Create final snapshot"
- Confirm deletion
```

**6. Delete Security Groups**
```
EC2 → Security Groups → delete rds-sg
EC2 → Security Groups → default → remove inbound rules you added
```

**7. Delete CloudFormation stacks**
```
CloudFormation → Stacks → delete any aws-ecs-blog-app related stacks
```

**8. Verify in Billing**
```
AWS Console → Billing → Bills → confirm nothing is still running
```

> RDS charges by the hour even when idle — make sure to delete it when done.
