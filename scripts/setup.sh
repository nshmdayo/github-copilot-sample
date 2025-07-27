#!/bin/bash

# Web Todo App - Setup Script
# フロントエンド（Next.js）、バックエンド（Go）、インフラ（Terraform）の初期セットアップ

set -e

echo "🚀 Web Todo App のセットアップを開始します..."

# カラー定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# ヘルパー関数
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

# 必要なツールのチェック
check_requirements() {
    print_info "必要なツールの確認中..."
    
    local missing_tools=()
    
    # Node.js
    if ! command -v node &> /dev/null; then
        missing_tools+=("Node.js (v18以上)")
    else
        NODE_VERSION=$(node --version | cut -d'v' -f2 | cut -d'.' -f1)
        if [ "$NODE_VERSION" -lt 18 ]; then
            missing_tools+=("Node.js (現在: v$NODE_VERSION, 必要: v18以上)")
        fi
    fi
    
    # Go
    if ! command -v go &> /dev/null; then
        missing_tools+=("Go (v1.21以上)")
    else
        GO_VERSION=$(go version | cut -d' ' -f3 | cut -d'o' -f2 | cut -d'.' -f2)
        if [ "$GO_VERSION" -lt 21 ]; then
            missing_tools+=("Go (現在: $(go version), 必要: v1.21以上)")
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
        missing_tools+=("Terraform (v1.5以上)")
    fi
    
    # Git
    if ! command -v git &> /dev/null; then
        missing_tools+=("Git")
    fi
    
    if [ ${#missing_tools[@]} -ne 0 ]; then
        print_error "以下のツールが不足しています:"
        for tool in "${missing_tools[@]}"; do
            echo "  - $tool"
        done
        echo ""
        echo "必要なツールをインストールしてから再実行してください。"
        exit 1
    fi
    
    print_success "すべての必要なツールが確認できました"
}

# フロントエンドセットアップ
setup_frontend() {
    print_info "フロントエンド（Next.js + Tailwind CSS）のセットアップ中..."
    
    if [ ! -d "frontend" ]; then
        mkdir -p frontend
        cd frontend
        
        # Next.jsプロジェクトの初期化
        npx create-next-app@latest . --typescript --tailwind --eslint --app --src-dir --import-alias "@/*"
        
        # 追加パッケージのインストール
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
        print_success "フロントエンドのセットアップが完了しました"
    else
        print_warning "フロントエンドディレクトリが既に存在します。スキップします。"
    fi
}

# バックエンドセットアップ
setup_backend() {
    print_info "バックエンド（Go + Gin + GORM）のセットアップ中..."
    
    if [ ! -d "backend" ]; then
        mkdir -p backend
        cd backend
        
        # Go モジュールの初期化
        go mod init github.com/nshmdayo/github-copilot-sample/backend
        
        # 必要なパッケージのインストール
        go get github.com/gin-gonic/gin
        go get gorm.io/gorm
        go get gorm.io/driver/postgres
        go get github.com/golang-jwt/jwt/v5
        go get github.com/go-playground/validator/v10
        go get github.com/spf13/viper
        go get github.com/sirupsen/logrus
        go get golang.org/x/crypto/bcrypt
        
        # 開発・テスト用パッケージ
        go get github.com/stretchr/testify
        go get github.com/DATA-DOG/go-sqlmock
        go get github.com/golang-migrate/migrate/v4
        
        # Swagger
        go get github.com/swaggo/gin-swagger
        go get github.com/swaggo/files
        
        cd ..
        print_success "バックエンドのセットアップが完了しました"
    else
        print_warning "バックエンドディレクトリが既に存在します。スキップします。"
    fi
}

# インフラセットアップ
setup_infrastructure() {
    print_info "インフラ（Terraform + AWS）のセットアップ中..."
    
    if [ ! -d "infrastructure" ]; then
        mkdir -p infrastructure/{environments/{dev,staging,prod},modules/{networking,security,database,ecs,alb,cloudfront,route53},scripts}
        
        print_success "インフラディレクトリ構造を作成しました"
    else
        print_warning "インフラディレクトリが既に存在します。スキップします。"
    fi
}

# Docker設定ファイル作成
setup_docker() {
    print_info "Docker設定ファイルの作成中..."
    
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
        print_success "docker-compose.yml を作成しました"
    fi
}

# 環境変数ファイル作成
setup_env_files() {
    print_info "環境変数ファイルの作成中..."
    
    # フロントエンド環境変数
    if [ ! -f "frontend/.env.local" ]; then
        mkdir -p frontend
        cat > frontend/.env.local << 'EOF'
# フロントエンド環境変数
NEXT_PUBLIC_API_URL=http://localhost:8080/api
NEXT_PUBLIC_APP_NAME=Todo App
EOF
        print_success "frontend/.env.local を作成しました"
    fi
    
    # バックエンド環境変数
    if [ ! -f "backend/.env" ]; then
        mkdir -p backend
        cat > backend/.env << 'EOF'
# バックエンド環境変数
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
        print_success "backend/.env を作成しました"
    fi
}

# Git設定
setup_git() {
    print_info "Git設定の確認中..."
    
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
        print_success ".gitignore を作成しました"
    fi
}

# 開発用スクリプト作成
setup_dev_scripts() {
    print_info "開発用スクリプトの作成中..."
    
    # 開発環境起動スクリプト
    cat > scripts/dev-start.sh << 'EOF'
#!/bin/bash

echo "🚀 開発環境を起動中..."

# Docker Composeで起動
docker-compose up -d postgres

echo "PostgreSQLの起動を待機中..."
sleep 10

# バックエンド起動（バックグラウンド）
cd backend
go run cmd/server/main.go &
BACKEND_PID=$!
cd ..

echo "バックエンドの起動を待機中..."
sleep 5

# フロントエンド起動
cd frontend
npm run dev &
FRONTEND_PID=$!
cd ..

echo "✅ 開発環境が起動しました!"
echo "📱 フロントエンド: http://localhost:3000"
echo "🔧 バックエンド: http://localhost:8080"
echo "📚 API ドキュメント: http://localhost:8080/swagger/index.html"
echo ""
echo "停止するには Ctrl+C を押してください"

# トラップでプロセスを終了
trap "kill $BACKEND_PID $FRONTEND_PID; docker-compose down" EXIT

wait
EOF
    
    chmod +x scripts/dev-start.sh
    
    # ビルドスクリプト
    cat > scripts/build-all.sh << 'EOF'
#!/bin/bash

echo "🔨 全体ビルドを開始します..."

# バックエンドビルド
echo "🔧 バックエンドビルド中..."
cd backend
go build -o bin/server ./cmd/server
cd ..

# フロントエンドビルド
echo "📱 フロントエンドビルド中..."
cd frontend
npm run build
cd ..

echo "✅ 全体ビルドが完了しました!"
EOF
    
    chmod +x scripts/build-all.sh
    
    print_success "開発用スクリプトを作成しました"
}

# READMEファイル更新
update_readme() {
    print_info "README.md の更新中..."
    
    cat > README.md << 'EOF'
# Web Todo Application

フロントエンド（Next.js + Tailwind CSS）、バックエンド（Go言語）、インフラ（AWS + Terraform）で構成されるWeb Todoアプリケーション。

## 🏗️ アーキテクチャ

- **フロントエンド**: Next.js 14 + Tailwind CSS + TypeScript
- **バックエンド**: Go + Gin + GORM + PostgreSQL
- **インフラ**: AWS ECS + RDS + CloudFront + Terraform
- **CI/CD**: GitHub Actions

## 🚀 クイックスタート

### 1. 初期セットアップ

```bash
# リポジトリをクローン
git clone https://github.com/nshmdayo/github-copilot-sample.git
cd github-copilot-sample

# セットアップスクリプトを実行
./scripts/setup.sh
```

### 2. 開発環境起動

```bash
# 開発環境を起動
./scripts/dev-start.sh
```

### 3. アクセス

- **フロントエンド**: http://localhost:3000
- **バックエンドAPI**: http://localhost:8080
- **API ドキュメント**: http://localhost:8080/swagger/index.html

## 📁 プロジェクト構造

```
├── frontend/              # Next.js アプリケーション
├── backend/               # Go API サーバー
├── infrastructure/        # Terraform 設定
├── .github/               # GitHub Actions & 開発指示書
├── scripts/               # 開発・運用スクリプト
└── docker-compose.yml     # ローカル開発環境
```

## 🔧 開発ガイド

詳細な開発指示については、以下のファイルを参照してください：

- [全体指示書](.github/instructions/project.instructions.md)
- [フロントエンド](.github/instructions/frontend.instructions.md)
- [バックエンド](.github/instructions/backend.instructions.md)
- [インフラ](.github/instructions/infrastructure.instructions.md)

## 📝 ライセンス

MIT License
EOF
    
    print_success "README.md を更新しました"
}

# メイン実行
main() {
    echo "🎯 Web Todo App セットアップスクリプト"
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
    echo "🎉 セットアップが完了しました！"
    echo ""
    echo "次のステップ:"
    echo "1. ./scripts/dev-start.sh で開発環境を起動"
    echo "2. GitHub Copilot を活用してアプリ開発を開始"
    echo "3. .github/instructions/ 内の指示書を参照"
    echo ""
    echo "Happy coding! 🚀"
}

# スクリプト実行
main "$@"
EOF
