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
