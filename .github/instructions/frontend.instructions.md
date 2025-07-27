# Frontend Development Instructions - Next.js + Tailwind CSS

## プロジェクト概要
Web TodoアプリのフロントエンドをNext.js 14（App Router）とTailwind CSSで構築します。

## 技術スタック
- **Framework**: Next.js 14 (App Router)
- **Styling**: Tailwind CSS
- **TypeScript**: 厳密な型定義を使用
- **State Management**: React hooks + Context API
- **HTTP Client**: fetch API / axios
- **UI Components**: Headless UI / Radix UI (必要に応じて)
- **Form Handling**: React Hook Form + Zod
- **Testing**: Jest + React Testing Library
- **Linting**: ESLint + Prettier

## プロジェクト構造
```
frontend/
├── app/                          # App Router用ディレクトリ
│   ├── (auth)/                   # 認証関連のルートグループ
│   │   ├── login/
│   │   └── register/
│   ├── todos/                    # Todo関連ページ
│   ├── api/                      # API routes (必要に応じて)
│   ├── globals.css               # グローバルスタイル
│   ├── layout.tsx                # ルートレイアウト
│   └── page.tsx                  # ホームページ
├── components/                   # 再利用可能コンポーネント
│   ├── ui/                       # 基本UIコンポーネント
│   ├── forms/                    # フォームコンポーネント
│   ├── layouts/                  # レイアウトコンポーネント
│   └── features/                 # 機能別コンポーネント
│       └── todos/
├── lib/                          # ユーティリティ・設定
│   ├── api.ts                    # API通信関数
│   ├── auth.ts                   # 認証関連
│   ├── utils.ts                  # ユーティリティ関数
│   └── validations.ts            # バリデーションスキーマ
├── hooks/                        # カスタムフック
├── types/                        # TypeScript型定義
├── middleware.ts                 # Next.js middleware
├── tailwind.config.js
├── next.config.js
└── package.json
```

## コーディング規則

### TypeScript
- 厳密な型定義を使用（`strict: true`）
- `any`型の使用を避ける
- インターフェースとタイプエイリアスを適切に使い分け
- 必要に応じてGenericsを活用

### コンポーネント設計
- 関数コンポーネントを使用
- Props型を明示的に定義
- デフォルトpropsは初期値設定で対応
- 可能な限り純粋コンポーネントを作成

```typescript
interface TodoItemProps {
  todo: Todo;
  onToggle: (id: string) => void;
  onDelete: (id: string) => void;
}

export function TodoItem({ todo, onToggle, onDelete }: TodoItemProps) {
  // コンポーネント実装
}
```

### Tailwind CSS
- カスタムクラスよりもユーティリティクラスを優先
- レスポンシブデザインを意識（mobile-first）
- ダークモード対応を考慮
- コンポーネント内でのクラス組織化

```typescript
const buttonStyles = {
  base: "px-4 py-2 rounded-md font-medium transition-colors",
  variants: {
    primary: "bg-blue-600 text-white hover:bg-blue-700",
    secondary: "bg-gray-200 text-gray-800 hover:bg-gray-300"
  }
};
```

### API通信
- 環境変数でAPIエンドポイントを管理
- エラーハンドリングを適切に実装
- ローディング状態を管理
- キャッシュ戦略を検討

```typescript
// lib/api.ts
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL;

export async function fetchTodos(): Promise<Todo[]> {
  try {
    const response = await fetch(`${API_BASE_URL}/todos`, {
      headers: {
        'Authorization': `Bearer ${getToken()}`,
        'Content-Type': 'application/json',
      },
    });
    
    if (!response.ok) {
      throw new Error('Failed to fetch todos');
    }
    
    return response.json();
  } catch (error) {
    console.error('Error fetching todos:', error);
    throw error;
  }
}
```

## データ型定義

### Todo関連
```typescript
// types/todo.ts
export interface Todo {
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

export interface CreateTodoRequest {
  title: string;
  description?: string;
  priority: Todo['priority'];
  dueDate?: string;
}

export interface UpdateTodoRequest extends Partial<CreateTodoRequest> {
  completed?: boolean;
}
```

### User関連
```typescript
// types/user.ts
export interface User {
  id: string;
  email: string;
  name: string;
  createdAt: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest extends LoginRequest {
  name: string;
}

export interface AuthResponse {
  user: User;
  token: string;
}
```

## 環境変数
```env
# .env.local
NEXT_PUBLIC_API_URL=http://localhost:8080/api
NEXT_PUBLIC_APP_NAME=Todo App
```

## パフォーマンス最適化
- Dynamic imports for code splitting
- Image optimization with next/image
- Font optimization
- Bundle analyzer for監視

## アクセシビリティ
- セマンティックHTML要素を使用
- 適切なARIA属性を設定
- キーボード操作をサポート
- カラーコントラストに配慮

## エラーハンドリング
- グローバルエラーバウンダリを設定
- 404/500ページのカスタマイズ
- 適切なエラーメッセージの表示
- ログ記録の実装

## セキュリティ
- XSS対策（適切なエスケープ）
- CSRF対策
- 認証トークンの安全な管理
- 環境変数での機密情報管理

## テスト戦略
- ユニットテスト: コンポーネントとユーティリティ関数
- 統合テスト: API通信とフォーム処理
- E2Eテスト: 主要なユーザーフロー

## 命名規則
- ファイル: kebab-case
- コンポーネント: PascalCase
- 関数・変数: camelCase
- 定数: UPPER_SNAKE_CASE
- CSS クラス: Tailwindユーティリティ中心

## Git関連
- feature/* ブランチで開発
- 適切なコミットメッセージ
- 必要に応じてPRテンプレートを活用
