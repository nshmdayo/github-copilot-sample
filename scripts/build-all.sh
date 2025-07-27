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
