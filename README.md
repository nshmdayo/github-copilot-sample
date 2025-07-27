# Web Todo Application

フロントエンド（Next.js + Tailwind CSS）、バックエンド（Go言語）、インフラ（AWS + Terraform）で構成されるWeb Todoアプリケーション。

[![CI/CD](https://github.com/nshmdayo/github-copilot-sample/actions/workflows/frontend-ci.yml/badge.svg)](https://github.com/nshmdayo/github-copilot-sample/actions)
[![Infrastructure](https://github.com/nshmdayo/github-copilot-sample/actions/workflows/infrastructure-cd.yml/badge.svg)](https://github.com/nshmdayo/github-copilot-sample/actions)

## 🏗️ アーキテクチャ

- **フロントエンド**: Next.js 14 + Tailwind CSS + TypeScript
- **バックエンド**: Go + Gin + GORM + PostgreSQL  
- **インフラ**: AWS ECS + RDS + CloudFront + Terraform
- **CI/CD**: GitHub Actions
- **開発環境**: Docker Compose

## 🚀 クイックスタート

### 前提条件

- Node.js 18以上
- Go 1.21以上
- Docker & Docker Compose
- Terraform 1.5以上

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
# 基本的な開発環境
docker-compose up -d

# 開発ツール付きで起動
docker-compose --profile tools up -d

# 監視ツール付きで起動  
docker-compose --profile monitoring up -d
```

### 3. アクセス

| サービス | URL | 説明 |
|---------|-----|------|
| フロントエンド | http://localhost:3000 | Next.js アプリケーション |
| バックエンドAPI | http://localhost:8080 | Go API サーバー |
| API ドキュメント | http://localhost:8080/swagger | Swagger UI |
| データベース管理 | http://localhost:8081 | Adminer |
| pgAdmin | http://localhost:8082 | PostgreSQL 管理ツール |
| Grafana | http://localhost:3001 | メトリクス可視化 |

## 📁 プロジェクト構造

```
github-copilot-sample/
├── .github/
│   ├── instructions/              # 各コンポーネントの開発指示書
│   │   ├── project.instructions.md
│   │   ├── frontend.instructions.md
│   │   ├── backend.instructions.md
│   │   └── infrastructure.instructions.md
│   └── workflows/                 # GitHub Actions CI/CD
│       ├── frontend-ci.yml
│       ├── backend-ci.yml
│       └── infrastructure-cd.yml
├── frontend/                      # Next.js + Tailwind CSS
├── backend/                       # Go API + PostgreSQL
├── infrastructure/                # Terraform + AWS
├── scripts/                       # 開発・運用スクリプト
├── docker-compose.yml             # ローカル開発環境
├── general.instructure.md         # GitHub Copilot 指示書
└── README.md
```

## 🔧 開発ガイド

### GitHub Copilot を活用した開発

このプロジェクトはGitHub Copilotを最大限活用できるよう設計されています。

1. **指示書の確認**: `.github/instructions/` 内の各コンポーネント指示書を参照
2. **コンテキストの活用**: 既存のコードパターンを参考にCopilotが適切な提案を行います
3. **プロンプトの最適化**: コメントで具体的な要件を記述してください

### 開発フロー

```bash
# 1. 機能ブランチを作成
git checkout -b feature/todo-crud-api

# 2. 開発環境で実装・テスト
docker-compose up -d
./scripts/dev-start.sh

# 3. テスト実行
cd backend && go test ./...
cd frontend && npm test

# 4. プルリクエスト作成
git push origin feature/todo-crud-api
```

## 📝 API 仕様

### 認証エンドポイント
- `POST /api/auth/register` - ユーザー登録
- `POST /api/auth/login` - ログイン
- `POST /api/auth/refresh` - トークンリフレッシュ

### Todo エンドポイント
- `GET /api/todos` - Todo一覧取得
- `POST /api/todos` - Todo作成
- `GET /api/todos/:id` - Todo詳細取得  
- `PUT /api/todos/:id` - Todo更新
- `DELETE /api/todos/:id` - Todo削除

詳細は [API ドキュメント](http://localhost:8080/swagger) を参照してください。

## 🏗️ インフラ構成

### AWS アーキテクチャ

```
Internet
    ↓
CloudFront (CDN)
    ↓
Route 53 (DNS)
    ↓
Application Load Balancer
    ↓
ECS Fargate (Frontend & Backend)
    ↓
RDS PostgreSQL (Multi-AZ)
```

### デプロイメント

```bash
# 環境別デプロイ
cd infrastructure/environments/dev
terraform init
terraform plan
terraform apply

# 本番環境はGitHub Actionsで自動デプロイ
```

## 🧪 テスト

### フロントエンド
```bash
cd frontend
npm test              # ユニットテスト
npm run test:e2e      # E2Eテスト
npm run test:coverage # カバレッジ測定
```

### バックエンド
```bash
cd backend
go test ./...                    # ユニットテスト
go test -tags=integration ./...  # 統合テスト
go test -coverprofile=coverage.out ./...  # カバレッジ
```

## 📊 監視・ログ

- **アプリケーションログ**: CloudWatch Logs
- **メトリクス**: CloudWatch + Prometheus
- **可視化**: Grafana ダッシュボード
- **アラート**: CloudWatch Alarms + SNS

## 🔒 セキュリティ

- JWT認証
- HTTPS強制
- CORS設定
- SQL Injection対策
- XSS対策
- AWS WAF

## 🚀 本番環境

### 環境構成
- **Development**: 機能開発・テスト
- **Staging**: 本番環境と同等構成
- **Production**: 本番環境

### CI/CD パイプライン
1. コードプッシュ
2. 自動テスト実行
3. セキュリティスキャン
4. Docker イメージビルド
5. ECRプッシュ
6. ECS デプロイ

## 🤝 コントリビューション

1. Issueを作成
2. フィーチャーブランチで開発
3. テストを追加/更新
4. プルリクエストを作成

## 📚 ドキュメント

- [プロジェクト全体指示](.github/instructions/project.instructions.md)
- [フロントエンド開発指示](.github/instructions/frontend.instructions.md)
- [バックエンド開発指示](.github/instructions/backend.instructions.md)
- [インフラ管理指示](.github/instructions/infrastructure.instructions.md)
- [GitHub Copilot活用指示](general.instructure.md)

## 📄 ライセンス

MIT License

## 🙏 謝辞

このプロジェクトはGitHub Copilotの機能を最大限活用できるよう設計されています。各指示書を参考に、効率的な開発を行ってください。