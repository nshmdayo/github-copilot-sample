# Frontend Development Instructions - Next.js + Tailwind CSS

## Project Overview
Build the frontend of a Web Todo application using Next.js 14 (App Router) and Tailwind CSS.

## Technology Stack
- **Framework**: Next.js 14 (App Router)
- **Styling**: Tailwind CSS
- **TypeScript**: Use strict type definitions
- **State Management**: React hooks + Context API
- **HTTP Client**: fetch API / axios
- **UI Components**: Headless UI / Radix UI (as needed)
- **Form Handling**: React Hook Form + Zod
- **Testing**: Jest + React Testing Library
- **Linting**: ESLint + Prettier

## Project Structure
```
frontend/
├── app/                          # App Router directory
│   ├── (auth)/                   # Authentication route groups
│   │   ├── login/
│   │   └── register/
│   ├── todos/                    # Todo-related pages
│   ├── api/                      # API routes (if needed)
│   ├── globals.css               # Global styles
│   ├── layout.tsx                # Root layout
│   └── page.tsx                  # Home page
├── components/                   # Reusable components
│   ├── ui/                       # Basic UI components
│   ├── forms/                    # Form components
│   ├── layouts/                  # Layout components
│   └── features/                 # Feature-specific components
│       └── todos/
├── lib/                          # Utilities & configuration
│   ├── api.ts                    # API communication functions
│   ├── auth.ts                   # Authentication related
│   ├── utils.ts                  # Utility functions
│   └── validations.ts            # Validation schemas
├── hooks/                        # Custom hooks
├── types/                        # TypeScript type definitions
├── middleware.ts                 # Next.js middleware
├── tailwind.config.js
├── next.config.js
└── package.json
```

## Coding Standards

### TypeScript
- Use strict type definitions (`strict: true`)
- Avoid using `any` type
- Properly distinguish between interfaces and type aliases
- Utilize Generics when necessary

### Component Design
- Use function components
- Explicitly define Props types
- Handle default props with initial value settings
- Create pure components whenever possible

```typescript
interface TodoItemProps {
  todo: Todo;
  onToggle: (id: string) => void;
  onDelete: (id: string) => void;
}

export function TodoItem({ todo, onToggle, onDelete }: TodoItemProps) {
  // Component implementation
}
```

### Tailwind CSS
- Prioritize utility classes over custom classes
- Keep responsive design in mind (mobile-first)
- Consider dark mode support
- Organize classes within components

```typescript
const buttonStyles = {
  base: "px-4 py-2 rounded-md font-medium transition-colors",
  variants: {
    primary: "bg-blue-600 text-white hover:bg-blue-700",
    secondary: "bg-gray-200 text-gray-800 hover:bg-gray-300"
  }
};
```

### API Communication
- Manage API endpoints with environment variables
- Implement proper error handling
- Manage loading states
- Consider caching strategies

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

## Data Type Definitions

### Todo Related
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

### User Related
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

## Environment Variables
```env
# .env.local
NEXT_PUBLIC_API_URL=http://localhost:8080/api
NEXT_PUBLIC_APP_NAME=Todo App
```

## Performance Optimization
- Dynamic imports for code splitting
- Image optimization with next/image
- Font optimization
- Bundle analyzer for monitoring

## Accessibility
- Use semantic HTML elements
- Set appropriate ARIA attributes
- Support keyboard navigation
- Consider color contrast

## Error Handling
- Set up global error boundary
- Customize 404/500 pages
- Display appropriate error messages
- Implement logging

## Security
- XSS protection (proper escaping)
- CSRF protection
- Secure authentication token management
- Sensitive information management with environment variables

## Testing Strategy
- Unit tests: Components and utility functions
- Integration tests: API communication and form processing
- E2E tests: Main user flows

## Naming Conventions
- Files: kebab-case
- Components: PascalCase
- Functions & Variables: camelCase
- Constants: UPPER_SNAKE_CASE
- CSS Classes: Tailwind utility-focused

## Git Related
- Develop in feature/* branches
- Proper commit messages
- Use PR templates when necessary
