# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

File Management System (文件管理系统) - A monorepo application with Go backend (Gin) and React frontend (Ant Design Pro), supporting large file chunked uploads, downloads, and previews.

### Technology Stack
- **Backend**: Go 1.21+, Gin, GORM, MySQL 8.0+, JWT (golang-jwt/jwt), bcrypt
- **Frontend**: React 18, Vite 5, Ant Design 5, Ant Design Pro Components, React Router v6, Axios, TypeScript
- **Package Manager**: pnpm workspaces

## Development Commands

### Starting Development Servers
```bash
# Start both API (port 8080) and Web (port 3000) together
pnpm dev

# Or use Makefile
make dev

# Start individually
pnpm dev:api    # Backend only
pnpm dev:web    # Frontend only
```

### Building
```bash
pnpm build          # Build both
pnpm build:api      # Build Go binary to bin/server
pnpm build:web      # Build React app to apps/web/dist
make build          # Alternative via Makefile
```

### Frontend-specific (apps/web)
```bash
cd apps/web
pnpm dev            # Start Vite dev server
pnpm build          # Production build
pnpm lint           # ESLint with fix
pnpm format         # Prettier format
```

### Backend-specific (apps/api)
```bash
cd apps/api
go run cmd/server/main.go        # Run server
go build -o ../../bin/server cmd/server/main.go    # Build binary
go mod download                   # Install dependencies
```

### Database
```bash
# Initialize database
mysql -u root -p < apps/api/sql/init.sql
make migrate                      # Alternative via Makefile
```

### Other
```bash
make install        # Install all dependencies (pnpm + go mod)
make clean          # Clean build artifacts
```

## Architecture

### Monorepo Structure
```
file-manager-system/
├── apps/
│   ├── api/          # Go backend (module: file-manager-service)
│   └── web/          # React frontend (package: @file-manager/web)
├── packages/
│   └── types/        # Shared TypeScript types (currently empty, reserved for future)
├── pnpm-workspace.yaml
└── Makefile
```

### Backend Architecture (apps/api)

**Layered Architecture**: handler → service → repository → model

```
internal/
├── handler/          # HTTP request handlers
├── service/          # Business logic layer
├── repository/       # Data access layer (GORM)
├── model/            # Database models (User, Document, Chunk)
├── middleware/       # Auth middleware (JWT verification)
├── pkg/              # Utilities (jwt, uuid, storage)
├── config/           # Config loading from YAML
└── router/           # Route definitions (Gin router)

cmd/server/main.go    # Application entry point
configs/config.yaml   # Configuration file
```

**Key Design Patterns**:
- Repository pattern for database access (user.go, document.go, chunk.go)
- Service layer separates business logic from HTTP handlers
- JWT middleware protects authenticated routes
- Token blacklist supports logout/revocation

**Storage**:
- Files stored in `uploads/` directory
- `uploads/chunks/` - Temporary chunk storage during upload
- `uploads/documents/` - Final assembled documents
- Uses UUID for document identification

### Frontend Architecture (apps/web)

```
src/
├── features/         # Feature modules (recommended for new features)
│   └── document/     # Document system features
│       └── applications/  # Application management feature
├── routes/           # TanStack Router file-based routes
│   └── _authenticated/  # Protected routes
│       └── document/     # Document system routes
│           └── applications/  # Application management routes
├── api/              # API service layer
│   └── document/     # Document system API calls
│       └── application.ts  # Application API
├── components/       # Reusable UI components
├── lib/              # Utilities (api-client, etc.)
└── stores/           # State management (Zustand)
```

**API Proxy**: Vite proxies `/api` requests to `http://localhost:8080` in development ([vite.config.ts](apps/web/vite.config.ts))

**Path Aliases**: `@/` → `src/`, `@file-manager/types` → `packages/types/src`

### File Upload Flow

Large files use **chunked upload** (5MB chunks, max 5GB):
1. `POST /api/v1/documents/chunks/init` - Initialize upload session, get upload_id
2. `POST /api/v1/documents/chunks/upload` - Upload chunks sequentially
3. `POST /api/v1/documents/chunks/complete` - Finalize upload, assemble chunks, get document_id

## Configuration

### Backend Configuration ([apps/api/configs/config.yaml](apps/api/configs/config.yaml))
- **Server**: Port 8080, mode (debug/release)
- **Database**: MySQL connection settings
- **Storage**: File paths, max file size (5GB), chunk size (5MB), allowed extensions
- **JWT**: Secret key, expiration (24 hours)

**Important**: Update `jwt.secret` and database credentials for production.

### Frontend Configuration
- Vite config: [apps/web/vite.config.ts](apps/web/vite.config.ts)
- TypeScript config: [apps/web/tsconfig.json](apps/web/tsconfig.json)

## Key Documentation

- [Backend API Documentation](apps/api/README.md)
- [Service-to-Service Integration Guide](apps/api/docs/SERVICE-TO-SERVICE.md)
- [Integration Guide](apps/api/docs/INTEGRATION-GUIDE.md)

## Important Notes

### Authentication
- JWT-based authentication with 24-hour token expiration
- Token blacklist enables logout functionality
- All authenticated endpoints require `Authorization: Bearer {token}` header
- Frontend stores token in localStorage (see [apps/web/src/services/auth.ts](apps/web/src/services/auth.ts))

### Database Models
- **User**: id, username, password (bcrypt hashed)
- **Document**: id (UUID), user_id, file_name, file_size, file_path, created_at
- **Chunk**: id, upload_id (UUID), chunk_number, chunk_path, file_size
- **Application**: id, app_name, app_account (auto-generated), app_secret (auto-generated), status

### Supported File Types
Images: jpg, jpeg, png, gif | Documents: pdf, doc, docx, xls, xlsx, ppt, pptx | Text: txt, md

### Go Module
The backend Go module is named `file-manager-service` (see [apps/api/go.mod](apps/api/go.mod)). When adding new dependencies, run:
```bash
cd apps/api && go get <package> && go mod tidy
```

### Environment Requirements
- Go 1.21+
- Node.js 18+
- pnpm 8+
- MySQL 8.0+

## Development Standards

### Code Organization for New Features

When adding new features, follow this directory structure:

#### Backend (apps/api)
```
internal/
├── model/            # Add new model here (e.g., application.go)
├── repository/       # Add data access layer (e.g., application.go)
├── service/          # Add business logic (e.g., application.go)
├── handler/          # Add HTTP handlers (e.g., application.go)
└── router/           # Register routes in router.go
```

#### Frontend (apps/web)
```
src/
├── api/              # API service layer
│   └── {module}/     # Module directory (e.g., document/)
│       └── {feature}/   # Feature API directory (e.g., application/)
│           ├── index.ts    # API methods (main export)
│           └── types.ts    # TypeScript types
├── features/         # Feature modules
│   └── {module}/     # Module directory (e.g., document/)
│       └── {feature}/   # Feature directory (e.g., applications/)
│           ├── index.tsx       # Main feature component
│           └── components/     # Feature sub-components
├── routes/           # TanStack Router file-based routes
│   └── _authenticated/
│       └── {module}/     # Module routes (e.g., document/)
│           └── {feature}/   # Feature routes (e.g., applications/)
│               └── index.tsx
```

**Example**: For "Application Management" in the document system:
- API: `src/api/document/application/` (directory with index.ts and types.ts)
- Feature: `src/features/document/applications/index.tsx`
- Route: `src/routes/_authenticated/document/applications/index.tsx`

### API Response Format
All API responses should follow this structure:
```json
{
  "code": 200,
  "message": "success",
  "data": { ... }
}
```

### Database Changes
When modifying database schema:
1. Update `apps/api/sql/init.sql` with table changes
2. Update `apps/api/internal/model/*.go` with corresponding models
3. Update `apps/api/internal/model/db.go` AutoMigrate() to include new models
