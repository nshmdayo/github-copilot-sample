# Project Instructions - Web Todo Application

## Project Overview
Development and operations instructions for a Web Todo application consisting of Frontend (Next.js + Tailwind CSS), Backend (Go language), and Infrastructure (AWS + Terraform).

## Architecture Overview
```
[Frontend: Next.js + Tailwind] 
           ↓ (HTTPS/API)
[Backend: Go + Gin + PostgreSQL]
           ↓ (Deploy)
[Infrastructure: AWS ECS + RDS + CloudFront]
           ↓ (IaC)
[Terraform + GitHub Actions]
```

## Project Structure
```
github-copilot-sample/
├── .github/
│   ├── instructions/             # Development instructions for each component
│   │   ├── frontend.instructions.md
│   │   ├── backend.instructions.md
│   │   ├── infrastructure.instructions.md
│   │   └── project.instructions.md
│   └── workflows/               # GitHub Actions CI/CD
│       ├── frontend-ci.yml
│       ├── backend-ci.yml
│       └── infrastructure-cd.yml
├── frontend/                    # Next.js + Tailwind CSS
│   ├── app/
│   ├── components/
│   ├── lib/
│   ├── types/
│   ├── hooks/
│   ├── middleware.ts
│   ├── next.config.js
│   ├── tailwind.config.js
│   ├── package.json
│   └── Dockerfile
├── backend/                     # Go + Gin + GORM
│   ├── cmd/
│   ├── internal/
│   ├── pkg/
│   ├── migrations/
│   ├── docker/
│   ├── go.mod
│   └── go.sum
├── infrastructure/              # Terraform + AWS
│   ├── environments/
│   ├── modules/
│   └── scripts/
├── docs/                        # Project documentation
│   ├── api/                     # API specifications
│   ├── deployment/              # Deployment procedures
│   └── architecture/            # Architecture diagrams
├── scripts/                     # Various scripts
│   ├── setup.sh                 # Initial setup
│   ├── dev-start.sh            # Development environment startup
│   └── build-all.sh            # Full build
├── docker-compose.yml           # Local development environment
├── README.md
└── general.instructure.md
```

## Development Flow

### 1. Environment Setup
```bash
# Clone repository
git clone https://github.com/nshmdayo/github-copilot-sample.git
cd github-copilot-sample

# Run initial setup
./scripts/setup.sh
```

### 2. Start Local Development Environment
```bash
# Start local environment with Docker Compose
docker-compose up -d

# Or individual startup
./scripts/dev-start.sh
```

### 3. Development Workflow
1. **Create Feature Branch**
   ```bash
   git checkout -b feature/todo-crud-api
   ```

2. **Development Implementation**
   - Backend API implementation
   - Frontend UI implementation
   - Infrastructure configuration update

3. **Run Tests**
   ```bash
   # Backend tests
   cd backend && go test ./...
   
   # Frontend tests
   cd frontend && npm test
   ```

4. **Create Pull Request**
   - CI/CD pipeline runs automatically
   - Conduct code review

5. **Merge & Deploy**
   - Merge to main branch
   - Execute automatic deployment

## API Design

### Endpoint List
```
# Authentication
POST   /api/auth/register      # User registration
POST   /api/auth/login         # Login
POST   /api/auth/refresh       # Token refresh

# Todo
GET    /api/todos              # Get Todo list
POST   /api/todos              # Create Todo
GET    /api/todos/:id          # Get Todo details
PUT    /api/todos/:id          # Update Todo
DELETE /api/todos/:id          # Delete Todo

# User
GET    /api/users/me           # Get user information
PUT    /api/users/me           # Update user information

# Health Check
GET    /health                 # Check application status
```

### Data Models
```typescript
interface User {
  id: string;
  email: string;
  name: string;
  createdAt: string;
}

interface Todo {
  id: string;
  title: string;
  description?: string;
  completed: boolean;
  priority: 'low' | 'medium' | 'high';
  dueDate?: string;
  createdAt: string;
  updatedAt: string;
  userId: string;
}
```

## Environment Configuration

### Development Environment
- **Frontend**: http://localhost:3000
- **Backend**: http://localhost:8080
- **Database**: PostgreSQL (localhost:5432)
- **Documentation**: http://localhost:8080/swagger/index.html

### Environment Variables
```env
# Frontend (.env.local)
NEXT_PUBLIC_API_URL=http://localhost:8080/api
NEXT_PUBLIC_APP_NAME=Todo App

# Backend (.env)
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=todoapp
JWT_SECRET=your-secret-key
GO_ENV=development

# Infrastructure
AWS_REGION=ap-northeast-1
TF_STATE_BUCKET=todoapp-terraform-state
```

## CI/CD Pipeline

### Frontend CI
```yaml
name: Frontend CI
on:
  push:
    paths: ['frontend/**']
  pull_request:
    paths: ['frontend/**']

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: '18'
          cache: 'npm'
          cache-dependency-path: frontend/package-lock.json
      
      - name: Install dependencies
        run: cd frontend && npm ci
      
      - name: Run linter
        run: cd frontend && npm run lint
      
      - name: Run tests
        run: cd frontend && npm test
      
      - name: Build application
        run: cd frontend && npm run build
      
      - name: Build Docker image
        run: |
          cd frontend
          docker build -t todoapp-frontend:${{ github.sha }} .
```

### Backend CI
```yaml
name: Backend CI
on:
  push:
    paths: ['backend/**']
  pull_request:
    paths: ['backend/**']

jobs:
  test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: password
          POSTGRES_DB: todoapp_test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5

    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Run tests
        run: |
          cd backend
          go test -v -race -coverprofile=coverage.out ./...
      
      - name: Run linter
        uses: golangci/golangci-lint-action@v3
        with:
          working-directory: backend
          version: latest
      
      - name: Build application
        run: |
          cd backend
          go build -o bin/server ./cmd/server
      
      - name: Build Docker image
        run: |
          cd backend
          docker build -t todoapp-backend:${{ github.sha }} .
```

## Security Requirements

### Authentication & Authorization
- JWT token-based authentication
- Appropriate token expiration settings
- Refresh token implementation

### Data Protection
- Enforce HTTPS communication
- Password hashing (bcrypt)
- Environment variable management for sensitive information

### API Security
- CORS configuration
- Rate Limiting
- Input Validation
- SQL Injection prevention

### Infrastructure Security
- Web application protection with AWS WAF
- Proper network control with Security Groups
- Least privilege access with IAM roles

## Performance Requirements

### Frontend
- First Contentful Paint < 1.5s
- Largest Contentful Paint < 2.5s
- Cumulative Layout Shift < 0.1
- Time to Interactive < 3.0s

### Backend
- API response time < 200ms (95 percentile)
- Database query optimization
- Proper index configuration

### Infrastructure
- Load handling with Auto Scaling
- Static content delivery with CloudFront
- High availability with RDS Multi-AZ

## Monitoring & Logging

### Application Monitoring
- CloudWatch metrics
- ECS Container Insights
- RDS Performance Insights

### Log Management
- Structured logging (JSON format)
- Aggregation with CloudWatch Logs
- Proper log level configuration

### Alert Configuration
- CPU/Memory usage
- Error rate
- Response time
- Database connection count

## Testing Strategy

### Frontend
- **Unit Tests**: React Testing Library + Jest
- **Integration Tests**: Mock API communication
- **E2E Tests**: Playwright

### Backend
- **Unit Tests**: Go standard testing + testify
- **Integration Tests**: Database operations
- **API Tests**: HTTP request/response

### Infrastructure
- **Infrastructure Tests**: Terratest
- **Security Tests**: tfsec, checkov

## Documentation Management

### API Documentation
- Swagger/OpenAPI specification
- Auto-generated + manual maintenance
- Accessible in each environment

### Architecture Documentation
- System architecture diagrams
- Data flow diagrams
- Security architecture

## Branching Strategy

### Git Flow
```
main                 # Production environment
├── develop          # Development integration branch
├── feature/*        # Feature development
├── release/*        # Release preparation
└── hotfix/*         # Emergency fixes
```

### Commit Message Conventions
```
type(scope): description

feat(backend): add todo CRUD API endpoints
fix(frontend): resolve authentication token refresh issue
docs(readme): update setup instructions
refactor(infrastructure): optimize ECS task definition
```

## Quality Management

### Code Review
- Minimum 1 reviewer required
- Automated tests must pass
- Security checks performed

### Quality Gates
- Test coverage 80% or higher
- 0 linter errors
- Vulnerability scan passed

## Deployment Strategy

### Environment Configuration
- **Development**: Feature development and testing
- **Staging**: Production-equivalent configuration testing
- **Production**: Production environment

### Deployment Methods
- Blue-Green Deployment
- Rolling Update (development environment)
- Canary Release (production environment)

## Operations & Maintenance

### Regular Maintenance
- Dependency updates
- Security patch application
- Performance tuning

### Backup Strategy
- RDS automatic backup
- Regular data export
- Disaster recovery plan

## Troubleshooting

### Common Issues and Solutions
1. **Application startup failure**
   - Check logs: CloudWatch Logs
   - Verify environment variable configuration
   - Confirm database connection

2. **API response delays**
   - Check CloudWatch metrics
   - Optimize database queries
   - Increase ECS task count

3. **Authentication errors**
   - Verify JWT token validity
   - Check JWT_SECRET environment variable
   - Verify CORS configuration
