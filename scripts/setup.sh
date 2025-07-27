#!/bin/bash

# Web Todo App - Setup Script
# Initial setup for Frontend (Next.js), Backend (Go), Infrastructure (Terraform)

set -e

echo "🚀 Starting Web Todo App setup..."

# Color definitions
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Check required tools
check_requirements() {
    print_info "Checking required tools..."
    
    local missing_tools=()
    
    # Node.js
    if ! command -v node &> /dev/null; then
        missing_tools+=("Node.js (v18 or higher)")
    else
        NODE_VERSION=$(node --version | cut -d'v' -f2 | cut -d'.' -f1)
        if [ "$NODE_VERSION" -lt 18 ]; then
            missing_tools+=("Node.js (current: v$NODE_VERSION, required: v18 or higher)")
        fi
    fi
    
    # Go
    if ! command -v go &> /dev/null; then
        missing_tools+=("Go (v1.21 or higher)")
    else
        GO_VERSION=$(go version | cut -d' ' -f3 | cut -d'o' -f2 | cut -d'.' -f2)
        if [ "$GO_VERSION" -lt 21 ]; then
            missing_tools+=("Go (current: $(go version), required: v1.21 or higher)")
        fi
    fi
    
    # Docker
    if ! command -v docker &> /dev/null; then
        missing_tools+=("Docker")
    fi
    
    # Docker Compose
    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
        missing_tools+=("Docker Compose")
    fi
    
    # Terraform
    if ! command -v terraform &> /dev/null; then
        missing_tools+=("Terraform (v1.5 or higher)")
    fi
    
    # Git
    if ! command -v git &> /dev/null; then
        missing_tools+=("Git")
    fi
    
    if [ ${#missing_tools[@]} -ne 0 ]; then
        print_error "The following tools are missing:"
        for tool in "${missing_tools[@]}"; do
            echo "  - $tool"
        done
        echo ""
        echo "Please install the required tools and run again."
        exit 1
    fi
    
    print_success "All required tools have been verified"
}

# Frontend setup
setup_frontend() {
    print_info "Setting up frontend (Next.js + Tailwind CSS)..."
    
    if [ ! -d "frontend" ]; then
        mkdir -p frontend
        cd frontend
        
        # Initialize Next.js project
        npx create-next-app@latest . --typescript --tailwind --eslint --app --src-dir --import-alias "@/*"
        
        # Install additional packages
        npm install --save \
            @hookform/resolvers \
            react-hook-form \
            zod \
            axios \
            @headlessui/react \
            @heroicons/react \
            clsx \
            tailwind-merge
        
        npm install --save-dev \
            @testing-library/react \
            @testing-library/jest-dom \
            @testing-library/user-event \
            jest \
            jest-environment-jsdom \
            @playwright/test \
            prettier \
            prettier-plugin-tailwindcss
        
        cd ..
        print_success "Frontend setup completed"
    else
        print_warning "Frontend directory already exists. Skipping."
    fi
}

# Backend setup
setup_backend() {
    print_info "Setting up backend (Go + Gin + GORM)..."
    
    if [ ! -d "backend" ]; then
        mkdir -p backend
        cd backend
        
        # Initialize Go module
        go mod init github.com/nshmdayo/github-copilot-sample/backend
        
        # Install required packages
        go get github.com/gin-gonic/gin
        go get gorm.io/gorm
        go get gorm.io/driver/postgres
        go get github.com/golang-jwt/jwt/v5
        go get github.com/go-playground/validator/v10
        go get github.com/spf13/viper
        go get github.com/sirupsen/logrus
        go get golang.org/x/crypto/bcrypt
        
        # Development and testing packages
        go get github.com/stretchr/testify
        go get github.com/DATA-DOG/go-sqlmock
        go get github.com/golang-migrate/migrate/v4
        
        # Swagger
        go get github.com/swaggo/gin-swagger
        go get github.com/swaggo/files
        
        cd ..
        print_success "Backend setup completed"
    else
        print_warning "Backend directory already exists. Skipping."
    fi
}

# Infrastructure setup
setup_infrastructure() {
    print_info "Setting up infrastructure (Terraform + AWS)..."
    
    if [ ! -d "infrastructure" ]; then
        mkdir -p infrastructure/{environments/{dev,staging,prod},modules/{networking,security,database,ecs,alb,cloudfront,route53},scripts}
        
        print_success "Infrastructure directory structure created"
    else
        print_warning "Infrastructure directory already exists. Skipping."
    fi
}

# Create Docker configuration files
setup_docker() {
    print_info "Creating Docker configuration files..."
    
    # docker-compose.yml
    if [ ! -f "docker-compose.yml" ]; then
        cat > docker-compose.yml << 'EOF'
version: '3.8'

services:
  postgres:
    image: postgres:15
    container_name: todoapp-postgres
    environment:
      POSTGRES_DB: todoapp
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./backend/migrations:/docker-entrypoint-initdb.d
    networks:
      - todoapp-network

  backend:
    build:
      context: ./backend
      dockerfile: docker/Dockerfile
    container_name: todoapp-backend
    environment:
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: postgres
      DB_PASSWORD: password
      DB_NAME: todoapp
      JWT_SECRET: your-secret-key
      GO_ENV: development
    ports:
      - "8080:8080"
    depends_on:
      - postgres
    volumes:
      - ./backend:/app
    networks:
      - todoapp-network

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    container_name: todoapp-frontend
    environment:
      NEXT_PUBLIC_API_URL: http://localhost:8080/api
    ports:
      - "3000:3000"
    depends_on:
      - backend
    volumes:
      - ./frontend:/app
      - /app/node_modules
    networks:
      - todoapp-network

volumes:
  postgres_data:

networks:
  todoapp-network:
    driver: bridge
EOF
        print_success "docker-compose.yml created"
    fi
}

# Create environment variable files
setup_env_files() {
    print_info "Creating environment variable files..."
    
    # Frontend environment variables
    if [ ! -f "frontend/.env.local" ]; then
        mkdir -p frontend
        cat > frontend/.env.local << 'EOF'
# Frontend environment variables
NEXT_PUBLIC_API_URL=http://localhost:8080/api
NEXT_PUBLIC_APP_NAME=Todo App
EOF
        print_success "frontend/.env.local created"
    fi
    
    # Backend environment variables
    if [ ! -f "backend/.env" ]; then
        mkdir -p backend
        cat > backend/.env << 'EOF'
# Backend environment variables
PORT=8080
HOST=localhost

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=todoapp
DB_SSL_MODE=disable

# JWT
JWT_SECRET=your-secret-key-please-change-in-production
JWT_EXPIRATION=24h

# Environment
GO_ENV=development
EOF
        print_success "backend/.env created"
    fi
}

# Git configuration
setup_git() {
    print_info "Checking Git configuration..."
    
    if [ ! -f ".gitignore" ]; then
        cat > .gitignore << 'EOF'
# Dependencies
node_modules/
vendor/

# Build outputs
dist/
build/
.next/
bin/

# Environment variables
.env
.env.local
.env.*.local

# Logs
*.log
logs/

# Runtime data
pids/
*.pid
*.seed
*.pid.lock

# Coverage directory used by tools like istanbul
coverage/

# Terraform
*.tfstate
*.tfstate.*
.terraform/
*.tfplan
.terraform.lock.hcl

# IDE
.vscode/
.idea/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# Database
*.db
*.sqlite

# Temporary files
tmp/
temp/
EOF
        print_success ".gitignore created"
    fi
}

# Create development scripts
setup_dev_scripts() {
    print_info "Creating development scripts..."
    
    # Development environment startup script
    cat > scripts/dev-start.sh << 'EOF'
#!/bin/bash

echo "🚀 Starting development environment..."

# Start with Docker Compose
docker-compose up -d postgres

echo "Waiting for PostgreSQL to start..."
sleep 10

# Start backend (background)
cd backend
go run cmd/server/main.go &
BACKEND_PID=$!
cd ..

echo "Waiting for backend to start..."
sleep 5

# Start frontend
cd frontend
npm run dev &
FRONTEND_PID=$!
cd ..

echo "✅ Development environment started!"
echo "📱 Frontend: http://localhost:3000"
echo "🔧 Backend: http://localhost:8080"
echo "📚 API Documentation: http://localhost:8080/swagger/index.html"
echo ""
echo "Press Ctrl+C to stop"

# Trap to terminate processes
trap "kill $BACKEND_PID $FRONTEND_PID; docker-compose down" EXIT

wait
EOF
    
    chmod +x scripts/dev-start.sh
    
    # Build script
    cat > scripts/build-all.sh << 'EOF'
#!/bin/bash

echo "🔨 Starting full build..."

# Backend build
echo "🔧 Building backend..."
cd backend
go build -o bin/server ./cmd/server
cd ..

# Frontend build
echo "📱 Building frontend..."
cd frontend
npm run build
cd ..

echo "✅ Full build completed!"
EOF
    
    chmod +x scripts/build-all.sh
    
    print_success "Development scripts created"
}

# Update README file
update_readme() {
    print_info "Updating README.md..."
    
    cat > README.md << 'EOF'
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
EOF
    
    print_success "README.md updated"
}

# Main execution
main() {
    echo "🎯 Web Todo App Setup Script"
    echo "======================================"
    
    check_requirements
    setup_frontend
    setup_backend
    setup_infrastructure
    setup_docker
    setup_env_files
    setup_git
    setup_dev_scripts
    update_readme
    
    echo ""
    echo "🎉 Setup completed!"
    echo ""
    echo "Next steps:"
    echo "1. Run ./scripts/dev-start.sh to start development environment"
    echo "2. Start app development with GitHub Copilot"
    echo "3. Refer to instruction files in .github/instructions/"
    echo ""
    echo "Happy coding! 🚀"
}

# Execute script
main "$@"
EOF
