# Project Instructions - Web Todo Application

## プロジェクト概要
フロントエンド（Next.js + Tailwind CSS）、バックエンド（Go言語）、インフラ（AWS + Terraform）で構成されるWeb Todoアプリケーションの開発・運用指示書です。

## アーキテクチャ概要
```
[Frontend: Next.js + Tailwind] 
           ↓ (HTTPS/API)
[Backend: Go + Gin + PostgreSQL]
           ↓ (Deploy)
[Infrastructure: AWS ECS + RDS + CloudFront]
           ↓ (IaC)
[Terraform + GitHub Actions]
```

## プロジェクト構造
```
github-copilot-sample/
├── .github/
│   ├── instructions/             # 各コンポーネントの開発指示書
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
├── docs/                        # プロジェクトドキュメント
│   ├── api/                     # API仕様書
│   ├── deployment/              # デプロイ手順
│   └── architecture/            # アーキテクチャ図
├── scripts/                     # 各種スクリプト
│   ├── setup.sh                 # 初期セットアップ
│   ├── dev-start.sh            # 開発環境起動
│   └── build-all.sh            # 全体ビルド
├── docker-compose.yml           # ローカル開発環境
├── README.md
└── general.instructure.md
```

## 開発フロー

### 1. 環境セットアップ
```bash
# リポジトリクローン
git clone https://github.com/nshmdayo/github-copilot-sample.git
cd github-copilot-sample

# 初期セットアップ実行
./scripts/setup.sh
```

### 2. ローカル開発環境起動
```bash
# Docker Composeでローカル環境起動
docker-compose up -d

# または個別起動
./scripts/dev-start.sh
```

### 3. 開発ワークフロー
1. **機能ブランチ作成**
   ```bash
   git checkout -b feature/todo-crud-api
   ```

2. **開発実施**
   - バックエンドAPI実装
   - フロントエンド画面実装
   - インフラ設定更新

3. **テスト実行**
   ```bash
   # バックエンドテスト
   cd backend && go test ./...
   
   # フロントエンドテスト
   cd frontend && npm test
   ```

4. **プルリクエスト作成**
   - CI/CDパイプラインが自動実行
   - コードレビュー実施

5. **マージ & デプロイ**
   - main ブランチへマージ
   - 自動デプロイ実行

## API設計

### エンドポイント一覧
```
# 認証
POST   /api/auth/register      # ユーザー登録
POST   /api/auth/login         # ログイン
POST   /api/auth/refresh       # トークンリフレッシュ

# Todo
GET    /api/todos              # Todo一覧取得
POST   /api/todos              # Todo作成
GET    /api/todos/:id          # Todo詳細取得
PUT    /api/todos/:id          # Todo更新
DELETE /api/todos/:id          # Todo削除

# ユーザー
GET    /api/users/me           # ユーザー情報取得
PUT    /api/users/me           # ユーザー情報更新

# ヘルスチェック
GET    /health                 # アプリケーション状態確認
```

### データモデル
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

## 環境設定

### 開発環境
- **Frontend**: http://localhost:3000
- **Backend**: http://localhost:8080
- **Database**: PostgreSQL (localhost:5432)
- **Documentation**: http://localhost:8080/swagger/index.html

### 環境変数
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

## CI/CD パイプライン

### フロントエンド CI
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

### バックエンド CI
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

## セキュリティ要件

### 認証・認可
- JWT トークンによる認証
- トークンの適切な有効期限設定
- リフレッシュトークンの実装

### データ保護
- HTTPS通信の強制
- パスワードのハッシュ化（bcrypt）
- 機密情報の環境変数管理

### API セキュリティ
- CORS設定
- Rate Limiting
- Input Validation
- SQL Injection 対策

### インフラセキュリティ
- AWS WAF による Web アプリケーション保護
- Security Groups による適切なネットワーク制御
- IAM ロールによる最小権限アクセス

## パフォーマンス要件

### フロントエンド
- First Contentful Paint < 1.5s
- Largest Contentful Paint < 2.5s
- Cumulative Layout Shift < 0.1
- Time to Interactive < 3.0s

### バックエンド
- API レスポンス時間 < 200ms (95 percentile)
- データベースクエリ最適化
- 適切なインデックス設定

### インフラ
- Auto Scaling による負荷対応
- CloudFront による静的コンテンツ配信
- RDS Multi-AZ による高可用性

## モニタリング & ログ

### アプリケーション監視
- CloudWatch メトリクス
- ECS Container Insights
- RDS Performance Insights

### ログ管理
- 構造化ログ（JSON形式）
- CloudWatch Logs による集約
- 適切なログレベル設定

### アラート設定
- CPU/Memory使用率
- エラー率
- レスポンス時間
- データベース接続数

## テスト戦略

### フロントエンド
- **ユニットテスト**: React Testing Library + Jest
- **統合テスト**: API通信のモック
- **E2Eテスト**: Playwright

### バックエンド
- **ユニットテスト**: Go標準テスト + testify
- **統合テスト**: データベース操作
- **API テスト**: HTTPリクエスト/レスポンス

### インフラ
- **Infrastructure Tests**: Terratest
- **Security Tests**: tfsec, checkov

## ドキュメント管理

### API ドキュメント
- Swagger/OpenAPI仕様書
- 自動生成 + 手動メンテナンス
- 各環境でのアクセス可能

### アーキテクチャドキュメント
- システム構成図
- データフロー図
- セキュリティアーキテクチャ

## ブランチ戦略

### Git Flow
```
main                 # 本番環境
├── develop          # 開発統合ブランチ
├── feature/*        # 機能開発
├── release/*        # リリース準備
└── hotfix/*         # 緊急修正
```

### コミットメッセージ規則
```
type(scope): description

feat(backend): add todo CRUD API endpoints
fix(frontend): resolve authentication token refresh issue
docs(readme): update setup instructions
refactor(infrastructure): optimize ECS task definition
```

## 品質管理

### コードレビュー
- 最低1人のレビュアー必須
- 自動テスト通過が必須
- セキュリティチェック実施

### 品質ゲート
- テストカバレッジ 80% 以上
- Linter エラー 0 件
- 脆弱性スキャン通過

## デプロイメント戦略

### 環境構成
- **Development**: 機能開発・テスト
- **Staging**: 本番環境と同等構成でのテスト
- **Production**: 本番環境

### デプロイ方式
- Blue-Green Deployment
- Rolling Update (開発環境)
- Canary Release (本番環境)

## 運用・保守

### 定期メンテナンス
- 依存関係の更新
- セキュリティパッチ適用
- パフォーマンスチューニング

### バックアップ戦略
- RDS 自動バックアップ
- データの定期的なエクスポート
- 災害復旧計画

## トラブルシューティング

### 一般的な問題と解決方法
1. **アプリケーション起動失敗**
   - ログ確認: CloudWatch Logs
   - 環境変数設定確認
   - データベース接続確認

2. **API レスポンス遅延**
   - CloudWatch メトリクス確認
   - データベースクエリ最適化
   - ECS タスク数増加

3. **認証エラー**
   - JWT トークン有効性確認
   - 環境変数 JWT_SECRET 確認
   - CORS設定確認
