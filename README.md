# Todo List

A full-stack todo list application with PC and mobile (H5) support.

## Tech Stack

| Layer | Technology |
|-------|------------|
| Frontend | Vue 3 (Composition API), TypeScript, Vite |
| Styling | Tailwind CSS v4 (class-based dark mode) |
| State | Pinia |
| Routing | Vue Router 4 |
| HTTP | Axios |
| PWA | vite-plugin-pwa (offline + installable) |
| Backend | Go, Gin, GORM |
| Database | SQLite (pure Go driver, no CGO needed) |
| Testing | Vitest (frontend), Go stdlib (backend) |

## Features

### Core
- CRUD todos with title, description, priority (low/medium/high), tags, due date
- Toggle completion with checkbox
- Inline editing with full field support

### Organization
- **Multi-list**: Create/rename/delete color-coded lists with sidebar navigation
- **Tags**: Comma-separated tags with auto-aggregated filter dropdown
- **Subtasks**: Expandable checklist within each todo, progress indicator
- **Archive**: Archive completed items instead of deleting
- **Drag & drop**: Reorder todos via native HTML5 drag and drop, persisted to backend

### Search & Filter
- Full-text search across title and description
- Filter by priority (high/medium/low)
- Filter by tag (dynamic options from existing tags)
- Sort by newest, priority, or due date
- Bottom navigation tabs: All / Active / Completed

### User Experience
- **Undo delete**: 10-second undo window with auto-restore
- **Confirm dialogs**: Before destructive actions
- **Toast notifications**: Success/error feedback for all operations
- **Skeleton loading**: Animated placeholders during data fetch
- **Keyboard shortcuts**: `Ctrl+N` new, `Ctrl+K` search, `Esc` clear, `?` panel
- **Due date highlighting**: Overdue (red border), due today (amber border)
- **Browser notifications**: Push notifications for due/overdue items (every 5 minutes)

### Appearance
- **Dark mode**: System-aware with manual toggle, persisted in localStorage
- **Responsive**: PC sidebar + full layout, mobile bottom nav + tap-to-detail
- **Animated transitions**: List enter/leave animations

### Data
- **Export**: JSON and CSV download (hidden in "more" dropdown)
- **Import**: n/a (CSV can be re-imported via API)
- **PWA**: Installable on mobile/home screen with service worker caching

### Statistics Dashboard
- Total / Active / Completed counters
- Completion rate progress bar
- Priority distribution chart
- Tag distribution chart

## API Reference

### Todos
| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/todos` | List todos (query: `?list_id=&status=&priority=&tag=&search=`) |
| `POST` | `/api/todos` | Create todo |
| `GET` | `/api/todos/:id` | Get todo |
| `PUT` | `/api/todos/:id` | Update todo |
| `PATCH` | `/api/todos/:id/toggle` | Toggle completed |
| `PATCH` | `/api/todos/:id/archive` | Archive todo |
| `PATCH` | `/api/todos/:id/unarchive` | Unarchive todo |
| `DELETE` | `/api/todos/:id` | Delete todo |
| `PUT` | `/api/todos/reorder` | Reorder todos `{ "ids": [3,1,2] }` |
| `GET` | `/api/todos/stats` | Get statistics |
| `GET` | `/api/todos/export` | Export (query: `?format=json|csv&list_id=`) |

### Subtasks
| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/todos/:id/subtasks` | List subtasks |
| `POST` | `/api/todos/:id/subtasks` | Create subtask `{ "title": "..." }` |
| `PATCH` | `/subtasks/:id/toggle` | Toggle subtask |
| `DELETE` | `/subtasks/:id` | Delete subtask |

### Attachments
| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/todos/:id/attachments` | List attachments |
| `POST` | `/api/todos/:id/attachments` | Upload file (multipart) |
| `GET` | `/attachments/:id` | Download attachment |
| `DELETE` | `/attachments/:id` | Delete attachment |

### Lists
| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/lists` | List all lists |
| `POST` | `/api/lists` | Create list `{ "name": "...", "color": "#6366f1" }` |
| `PUT` | `/api/lists/:id` | Update list |
| `DELETE` | `/api/lists/:id` | Delete list (moves todos to uncategorized) |

## Project Structure

```
todo-list/
├── backend/
│   ├── main.go                    # Entry point
│   ├── config/config.go           # Environment config
│   ├── database/database.go       # SQLite + AutoMigrate
│   ├── models/                    # Data models (Todo, List, Subtask, Attachment, DateOnly)
│   ├── repository/                # Database access layer
│   ├── services/                  # Business logic layer
│   ├── handlers/                  # HTTP handlers (CRUD, export, stats, attachments)
│   ├── router/router.go           # Gin router + CORS
│   └── Dockerfile
├── frontend/
│   ├── src/
│   │   ├── api/todos.ts           # Axios API client
│   │   ├── stores/                # Pinia stores (todos, lists, toast)
│   │   ├── composables/           # useDarkMode, useResponsive, useNotifications
│   │   ├── components/
│   │   │   ├── layout/            # AppHeader, AppNav, Sidebar
│   │   │   ├── ui/                # ConfirmDialog, ToastContainer, Skeleton, ShortcutPanel
│   │   │   ├── TodoForm.vue
│   │   │   ├── TodoItem.vue
│   │   │   └── TodoList.vue
│   │   ├── views/                 # HomeView, TodoDetailView, StatsView
│   │   ├── router/index.ts
│   │   └── types/todo.ts
│   ├── nginx.conf                 # Production nginx config
│   └── Dockerfile
├── docker-compose.yml
└── README.md
```

## Getting Started

### Development

```bash
# Terminal 1: Backend
cd backend
go run main.go
# Listening on :8080

# Terminal 2: Frontend
cd frontend
npm install
npm run dev
# Listening on :5173, proxies /api → :8080
```

### Docker

```bash
docker compose up -d
# Frontend → :80
# Backend  → :8080
```

### Running Tests

```bash
# Backend
cd backend && go test ./...

# Frontend
cd frontend && npx vitest run
```

## Configuration

| Env | Default | Description |
|-----|---------|-------------|
| `PORT` | `8080` | Backend listen port |
| `DB_PATH` | `todo.db` | SQLite database path |
| `CORS_ORIGIN` | `http://localhost:5173` | Allowed CORS origin |
