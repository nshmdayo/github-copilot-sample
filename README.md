# Web Todo Application

A Web Todo application built with Frontend (Next.js + Tailwind CSS), Backend (Go), and Infrastructure (AWS + Terraform).

## 🏗️ Architecture

- **Frontend**: Next.js 14 + Tailwind CSS + TypeScript
- **Backend**: Go + Gin + GORM + PostgreSQL
- **Infrastructure**: AWS ECS + RDS + CloudFront + Terraform
- **CI/CD**: GitHub Actions

## 🚀 Quick Start

### 1. Initial Setup

```bash
# Clone repository
git clone https://github.com/nshmdayo/github-copilot-sample.git
cd github-copilot-sample

# Run setup script
./scripts/setup.sh
```

### 2. Start Development Environment

```bash
# Start development environment
./scripts/dev-start.sh
```

### 3. Access

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **API Documentation**: http://localhost:8080/swagger/index.html

## 📁 Project Structure

```
├── frontend/              # Next.js application
├── backend/               # Go API server
├── infrastructure/        # Terraform configuration
├── .github/               # GitHub Actions & development instructions
├── scripts/               # Development & operation scripts
└── docker-compose.yml     # Local development environment
```

## 🔧 Development Guide

For detailed development instructions, refer to the following files:

- [Project Instructions](.github/instructions/project.instructions.md)
- [Frontend](.github/instructions/frontend.instructions.md)
- [Backend](.github/instructions/backend.instructions.md)
- [Infrastructure](.github/instructions/infrastructure.instructions.md)

## 📝 License

MIT License
