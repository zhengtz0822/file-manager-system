# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Common Commands

```bash
# Development
pnpm run dev              # Start development server

# Build
pnpm run build            # Type check and build for production

# Code Quality
pnpm run lint             # Run ESLint
pnpm run format:check     # Check Prettier formatting
pnpm run format           # Fix Prettier formatting
pnpm run knip             # Find unused dependencies

# Preview
pnpm run preview          # Preview production build locally
```

## Architecture Overview

### Routing Structure

This project uses **TanStack Router** with file-based routing and automatic code-splitting (configured in [vite.config.ts](vite.config.ts)).

- **`__root.tsx`**: Root route with global providers (QueryClient, Toaster, Devtools)
- **`_authenticated/`**: Protected route layout - wraps all authenticated pages with sidebar, search provider, and layout providers. Access to these routes requires authentication (check [authenticated-layout.tsx](src/components/layout/authenticated-layout.tsx))
- **`(auth)/`**: Authentication route group (sign-in, OTP verification, etc.)
- **`(errors)/`**: Error pages (401, 403, 404, 500, 503)

When adding new routes:
- Protected pages go under `_authenticated/`
- Auth pages go under `(auth)/`
- Route files use `route.tsx` for layouts and `index.tsx` for page content

### State Management

Three layers of state management:

1. **Server State**: TanStack Query for data fetching, caching, and synchronization
2. **Client State**: Zustand store for authentication state ([auth-store.ts](src/stores/auth-store.ts))
3. **Local State**: React Context for theme, direction (RTL/LTR), and font preferences

The auth store uses cookie-based token storage with a hardcoded cookie key (`thisisjustarandomstring`).

### Component Organization

- **`src/components/ui/`**: Shadcn UI components - many are customized for RTL support
- **`src/components/layout/`**: Layout components (Header, Sidebar, AppSidebar, etc.)
- **`src/features/`**: Feature-based modules (dashboard, settings, tasks, auth, errors)
- **`src/context/`**: React context providers (theme, direction, font, layout, search)

### RTL Support Warning

**Important**: Many Shadcn UI components in this project are customized for RTL (Right-to-Left) language support. When running `npx shadcn@latest add <component>`:

- **Safe to update via CLI**: Most non-listed components
- **Manual merge required**:
  - Modified Components: `scroll-area`, `sonner`, `separator`
  - RTL Updated Components: `alert-dialog`, `calendar`, `command`, `dialog`, `dropdown-menu`, `select`, `table`, `sheet`, `sidebar`, `switch`

See [README.md](README.md#L24-L59) for full details on customized components.

### Path Aliases

The `@/` alias maps to `./src/` (configured in [vite.config.ts](vite.config.ts:18-20) and [tsconfig.app.json](tsconfig.app.json)). Use this for all imports from src.

### Build Configuration

- **Bundler**: Vite with SWC for fast React compilation
- **TypeScript**: Strict mode enabled, builds must pass type checks
- **Styling**: TailwindCSS v4 with custom utilities and CSS variables for theming
- **Linting**: ESLint flat config with React hooks, TanStack Query, and custom import rules
- **Formatting**: Prettier with Trivago import sorting and Tailwind CSS class sorting

### Error Handling

- Global error boundaries at root route level
- Custom error components for different status codes in `(errors)/` routes
- Error handling utilities in [lib/handle-server-error.ts](src/lib/handle-server-error.ts)
- Toast notifications via Sonner for user feedback
