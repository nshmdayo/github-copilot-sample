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
