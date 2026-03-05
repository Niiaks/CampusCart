# CampusCart -- Client Application

The CampusCart client is a Next.js 16 web application that serves as the user-facing frontend for the CampusCart student marketplace. It provides interfaces for buyers to browse and discover listings, and for sellers to manage their brands, create listings, and communicate with buyers. The application uses React 19 with the App Router, Tailwind CSS 4 for styling, shadcn/ui for a polished component library, and TanStack React Query for server state management.

---

## Table of Contents

- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
- [Environment Variables](#environment-variables)
- [Available Scripts](#available-scripts)
- [Architecture](#architecture)
- [Routing and Layouts](#routing-and-layouts)
- [State Management](#state-management)
- [API Communication](#api-communication)
- [Authentication](#authentication)
- [Form Handling and Validation](#form-handling-and-validation)
- [UI and Styling](#ui-and-styling)
- [Component Library](#component-library)
- [Pages and Features](#pages-and-features)
- [Configuration Files](#configuration-files)
- [Code Conventions](#code-conventions)
- [Deployment](#deployment)

---

## Tech Stack

| Technology               | Version   | Purpose                                                            |
| ------------------------ | --------- | ------------------------------------------------------------------ |
| Next.js                  | 16.1.6    | React framework with App Router and server components              |
| React                    | 19.2.3    | UI component library                                               |
| TypeScript               | 5.x       | Static type checking across the entire codebase                    |
| Tailwind CSS             | 4.x       | Utility-first CSS framework                                        |
| shadcn/ui                | 3.x (CLI) | Pre-built accessible components on top of Radix UI primitives      |
| TanStack React Query     | 5.x       | Asynchronous server state management, caching, and synchronization |
| Axios                    | latest    | HTTP client for API communication                                  |
| Zod                      | 4.x       | Runtime schema validation for forms and API payloads               |
| React Hook Form          | 7.x       | Performant form state management with minimal re-renders           |
| Sonner                   | latest    | Lightweight toast notification system                              |
| Lucide React             | latest    | Consistent SVG icon set                                            |
| class-variance-authority | latest    | Variant-based component styling utility                            |
| tailwind-merge           | latest    | Intelligent Tailwind class merging to prevent conflicts            |
| tw-animate-css           | latest    | Tailwind CSS animation utilities                                   |

## Project Structure

```
apps/client/
  |-- app/                          # Next.js App Router directory
  |    |-- globals.css              # Global styles, Tailwind imports, CSS custom properties
  |    |-- layout.tsx               # Root layout (Inter font, metadata, Providers wrapper)
  |    |-- providers.tsx            # Client-side providers (QueryClient, Toaster, DevTools)
  |    |
  |    |-- (auth)/                  # Authentication route group
  |    |    |-- layout.tsx          # Centered layout with muted background
  |    |    |-- login/
  |    |    |    |-- page.tsx       # Login page
  |    |    |-- register/
  |    |    |    |-- page.tsx       # Registration page
  |    |    |-- verify-email/
  |    |         |-- page.tsx       # Email verification page
  |    |
  |    |-- (dashboard)/             # Authenticated dashboard route group
  |    |    |-- layout.tsx          # Navbar + floating chat button
  |    |    |-- brand/
  |    |    |    |-- page.tsx       # Brand management page
  |    |    |    |-- sections/      # Brand page sections
  |    |    |-- chat/
  |    |    |    |-- page.tsx       # Messaging/chat page
  |    |    |    |-- page.module.css
  |    |    |-- profile/
  |    |         |-- page.tsx       # User profile page
  |    |         |-- sections/      # Profile page sections
  |    |
  |    |-- (root)/                  # Public-facing route group
  |    |    |-- layout.tsx          # Navbar + Footer + floating chat button
  |    |    |-- page.tsx            # Home page (categories, trending, recommended)
  |    |    |-- [id]/
  |    |    |    |-- page.tsx       # Individual listing detail page
  |    |    |-- [id].tsx            # Alternative listing route
  |    |
  |    |-- components/              # Route-specific components
  |         |-- layout/             # Layout components (Navbar, Footer)
  |         |-- sections/           # Page sections organized by feature area
  |         |    |-- Authentication/  # Auth form components
  |         |    |-- Common/          # Shared section components
  |         |    |-- Dashboard/       # Dashboard-specific sections
  |         |    |-- Home/            # Home page sections
  |         |    |-- Listing/         # Listing display components
  |         |-- ui/                 # App-level UI components (ListingCard, etc.)
  |
  |-- common/                       # Shared utilities and types
  |    |-- types.ts                 # TypeScript interfaces (User, Brand, Listing, etc.)
  |    |-- schemas.ts               # Zod validation schemas
  |    |-- utils.ts                 # Error handling utilities
  |
  |-- components/                   # Base UI components (shadcn/ui)
  |    |-- ui/
  |         |-- button.tsx          # Button component with variants
  |         |-- dropdown-menu.tsx   # Dropdown menu component
  |         |-- input.tsx           # Input component
  |
  |-- constants/
  |    |-- constants.ts             # Application constants (API_URL)
  |
  |-- hooks/
  |    |-- useAuth.ts               # Authentication hooks (useAuth, useLogin, etc.)
  |
  |-- lib/
  |    |-- utils.ts                 # Utility functions (cn class merger)
  |
  |-- public/                       # Static assets served at root
  |
  |-- services/                     # API service layer
  |    |-- auth/
  |    |    |-- authService.ts      # Authentication API calls
  |    |-- listing/                 # Listing API calls (placeholder)
  |    |-- user/
  |         |-- userService.ts      # User API calls (placeholder)
  |
  |-- components.json               # shadcn/ui configuration
  |-- declarations.d.ts             # TypeScript module declarations
  |-- eslint.config.mjs             # ESLint 9 flat config
  |-- next-env.d.ts                 # Next.js TypeScript environment
  |-- next.config.ts                # Next.js configuration
  |-- package.json                  # Dependencies and scripts
  |-- postcss.config.mjs            # PostCSS configuration (Tailwind plugin)
  |-- tsconfig.json                 # TypeScript compiler configuration
```

## Prerequisites

- **Node.js** v20 or later
- **npm** (bundled with Node.js)
- A running instance of the CampusCart backend server (see `apps/server/README.md`)

## Getting Started

### 1. Install dependencies

```bash
cd apps/client
npm install
```

### 2. Configure environment variables

Create a `.env.local` file in the `apps/client` directory:

```env
NEXT_PUBLIC_API_URL=http://localhost:<server-port>
```

Replace `<server-port>` with the port your backend server is running on.

### 3. Start the development server

```bash
npm run dev
```

The application starts at `http://localhost:3000` with Turbopack enabled for fast refresh.

## Environment Variables

| Variable              | Required | Description                                                                                                                                                      |
| --------------------- | -------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `NEXT_PUBLIC_API_URL` | Yes      | The base URL of the CampusCart backend API server. This is a public variable exposed to the browser. If not set, the application will throw an error at startup. |

The variable is consumed in `constants/constants.ts` and used throughout the `services/` layer for all API calls.

## Available Scripts

| Command         | Description                                                                                               |
| --------------- | --------------------------------------------------------------------------------------------------------- |
| `npm run dev`   | Start the Next.js development server with Turbopack. Hot module replacement and fast refresh are enabled. |
| `npm run build` | Create an optimized production build. Runs type checking and tree shaking.                                |
| `npm run start` | Start the production server from the build output. Requires `npm run build` to have been run first.       |
| `npm run lint`  | Run ESLint with the Next.js core web vitals and TypeScript rule sets.                                     |

## Architecture

The client application follows a layered architecture pattern:

```
Pages (app/)
    |
    v
Components (app/components/, components/)
    |
    v
Hooks (hooks/)
    |
    v
Services (services/)
    |
    v
API Server (via Axios + withCredentials)
```

### Layer Responsibilities

1. **Pages** (`app/`): Define routes, compose section components, and handle page-level concerns like metadata and data fetching.
2. **Components**: Split into two categories:
   - `app/components/`: Route-specific components organized by feature area (Authentication, Dashboard, Home, Listing, etc.) and layout components (Navbar, Footer).
   - `components/ui/`: Base-level shadcn/ui primitives (Button, Input, DropdownMenu) shared across the entire application.
3. **Hooks** (`hooks/`): Custom React hooks encapsulating business logic sequences. Currently focused on authentication (`useAuth.ts`) providing `useAuth`, `useLogin`, `useRegister`, `useVerifyEmail`, and `useLogout` hooks built on React Query mutations and queries.
4. **Services** (`services/`): Thin API communication layer. Each service module exports functions that make Axios HTTP requests to the backend with `withCredentials: true` for cookie-based session authentication.
5. **Common** (`common/`): Shared TypeScript interfaces, Zod schemas, and utility functions used across multiple layers.

## Routing and Layouts

The application uses the Next.js App Router with route groups to organize pages by access level and layout requirements.

### Route Groups

**`(auth)/`** -- Authentication pages

- Applies a centered layout with a muted background.
- Contains: `/login`, `/register`, `/verify-email`.
- These pages are accessible without authentication.

**`(root)/`** -- Public-facing pages

- Applies the main site layout with Navbar at the top, Footer at the bottom, and a floating chat button.
- Contains: Home page (`/`), individual listing pages (`/[id]`).
- Browse functionality is available without authentication.

**`(dashboard)/`** -- Authenticated user pages

- Applies a dashboard layout with Navbar and floating chat button (no Footer).
- Contains: `/brand` (brand management), `/chat` (messaging), `/profile` (user profile).
- Requires an active session to access meaningful content.

### Root Layout

The root `layout.tsx` wraps the entire application with:

- The Inter font loaded via `next/font/google`.
- Metadata: title set to "CampusCart".
- The `<Providers>` component, which provides the React Query client and toast notifications.

## State Management

### Server State (React Query)

All server-side data is managed through TanStack React Query 5 via the `<Providers>` component. The QueryClient is configured with:

- **Stale time**: 5 minutes -- data is considered fresh for 5 minutes before background refetching is triggered.
- **Retry**: 1 attempt on failure before surfacing the error.

Query keys are used to scope and invalidate cached data. For example, the `["auth"]` query key manages the current user session state. On login success, the auth query data is set directly. On logout, the auth query is invalidated to clear cached user data.

React Query Devtools are included in development builds for inspecting cache state, query lifecycles, and mutation status.

### Client State

Local UI state is managed with standard React hooks (`useState`, `useReducer`). There is no global client state store; the application relies on React Query for any state that originates from the server.

## API Communication

All API communication is handled through the `services/` layer using Axios.

### Service Layer Pattern

Each service module exports plain async functions that:

1. Accept typed request parameters.
2. Make an HTTP request to the backend using Axios with `withCredentials: true`.
3. Return the typed response data.

Example from `authService.ts`:

```typescript
export const login = async (data: LoginRequest) => {
  const response = await axios.post<LoginResponse>(
    `${API_URL}/auth/login`,
    data,
    {
      withCredentials: true,
    },
  );
  return response.data;
};
```

The `withCredentials: true` flag ensures the browser sends and receives the `cc_refresh_token` HttpOnly cookie with every request, enabling cookie-based session authentication.

### Error Handling

API errors are handled through the `handleError` utility in `common/utils.ts`. It catches Axios errors and extracts the structured `APIError` response from the backend (which includes a code, message, and optional field-level errors). Non-Axios errors are re-thrown as generic errors.

The `APIError` interface:

```typescript
interface APIError {
  code: string;
  message: string;
  status: number;
  errors?: Record<string, string>;
}
```

## Authentication

Authentication is managed entirely through the `useAuth` hook family in `hooks/useAuth.ts`, which wraps React Query queries and mutations.

### Available Hooks

| Hook               | Type     | Description                                                                                                    |
| ------------------ | -------- | -------------------------------------------------------------------------------------------------------------- |
| `useAuth()`        | Query    | Fetches the currently authenticated user via `GET /auth/me`. Returns the user object or null.                  |
| `useLogin()`       | Mutation | Submits login credentials via `POST /auth/login`. On success, sets the auth query data directly and redirects. |
| `useRegister()`    | Mutation | Submits registration data via `POST /auth/register`. On success, redirects to email verification.              |
| `useVerifyEmail()` | Mutation | Submits the 6-digit verification code via `POST /auth/verify-email`. On success, sets the auth query data.     |
| `useLogout()`      | Mutation | Calls `POST /auth/logout` and invalidates the auth query to clear cached user data.                            |

### Session Management

Sessions are managed server-side with sliding expiration. The client does not store tokens in localStorage or memory. The HttpOnly `cc_refresh_token` cookie is automatically managed by the browser, providing protection against XSS-based token theft.

## Form Handling and Validation

Forms are built using React Hook Form 7 with Zod 4 schema validation.

### Defined Schemas (`common/schemas.ts`)

| Schema                | Fields                                                      | Description                   |
| --------------------- | ----------------------------------------------------------- | ----------------------------- |
| `loginSchema`         | email, password                                             | Login form validation         |
| `registerSchema`      | firstName, lastName, username, email, password              | Registration form validation  |
| `verifyEmailSchema`   | email, code                                                 | Email verification validation |
| `createBrandSchema`   | name, description                                           | Brand creation validation     |
| `createListingSchema` | title, description, price, condition, categoryId, imageUrls | Listing creation validation   |

Each schema exports an inferred TypeScript type (e.g., `LoginSchema`, `RegisterSchema`) for use in form components and service calls.

## UI and Styling

### Tailwind CSS 4

The application uses Tailwind CSS 4 with the `@tailwindcss/postcss` plugin. The configuration is handled through CSS custom properties defined in `app/globals.css` rather than a `tailwind.config.ts` file.

### Theme System

The application defines a complete light and dark theme using CSS custom properties on the `:root` and `.dark` selectors. Key theme tokens include:

- **Brand color**: `#4B3EC4` (indigo-purple), used as the primary accent throughout the application.
- **Background**: Light (`0 0% 100%`) and dark (`240 10% 3.9%`) variants.
- **Foreground, card, popover, muted, accent, destructive**: Full set of semantic color tokens.
- **Additional tokens**: success, warning, info, brand, and chart colors.
- **Sidebar-specific tokens**: For dashboard sidebar navigation.

### Base Typography

Global styles define consistent typography for heading elements (h1 through h6), paragraphs, links, and small text. These are applied as base layer styles in `globals.css`.

### Utility Classes

- `.scrollbar-none`: Hides scrollbars across browsers (standard `scrollbar-width: none` and webkit pseudo-element).

## Component Library

The application uses shadcn/ui with the "new-york" style variant, configured via `components.json`:

- **Style**: new-york
- **RSC (React Server Components)**: Enabled
- **TSX**: Enabled
- **Tailwind CSS**: Using CSS variables for theming
- **Base color**: neutral
- **Icons**: lucide

### Installed Base Components

These live in `components/ui/` and are the foundational building blocks:

- **Button** (`button.tsx`): Multi-variant button with default, destructive, outline, secondary, ghost, and link variants. Supports sm, default, lg, and icon sizes. Built with class-variance-authority.
- **Input** (`input.tsx`): Styled input field with consistent focus states and disabled styling.
- **Dropdown Menu** (`dropdown-menu.tsx`): Full dropdown menu built on Radix UI primitives with sub-menus, checkboxes, radio items, labels, separators, and keyboard shortcuts.

### Application Components

Route-specific and feature-specific components live in `app/components/`:

- **Layout**: `Navbar.tsx`, `Footer.tsx` -- persistent navigation and footer components.
- **Sections**: Organized by feature area:
  - `Authentication/` -- Login, register, and verification form components.
  - `Common/` -- Shared section components used across multiple pages.
  - `Dashboard/` -- Dashboard-specific content sections.
  - `Home/` -- Home page sections (Categories, ListingGrid, ListingCarousel).
  - `Listing/` -- Listing display and detail components.
- **UI**: `ListingCard.tsx` -- Card component for displaying listing previews in grids and carousels.

## Pages and Features

### Home Page (`/`)

The landing page displays several content sections:

- **Trending Categories**: Horizontal category browser.
- **Recommended For You**: Personalized listing suggestions.
- **Top Trending**: Popular listings by view count or engagement.
- **Top Seller**: Featured seller brands.
- **New In**: Recently added listings.

Uses `Categories`, `ListingGrid`, and `ListingCarousel` components for layout.

### Listing Detail Page (`/[id]`)

Displays full details for a single listing. Viewing a listing triggers an automatic view count increment on the backend.

### Authentication Pages

- **Login** (`/login`): Email and password form, redirects to home on success.
- **Register** (`/register`): First name, last name, username, email (must be `@st.ug.edu.gh`), and password. Redirects to verification on success.
- **Verify Email** (`/verify-email`): 6-digit code input. On success, creates a session and redirects to the dashboard.

### Dashboard Pages

- **Brand Management** (`/brand`): View and edit the seller's brand profile, including name, description, profile image, and banner image.
- **Chat** (`/chat`): Messaging interface for buyer-seller communication.
- **Profile** (`/profile`): View and manage user profile information.

## Configuration Files

| File                 | Purpose                                                                                                                                          |
| -------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------ |
| `next.config.ts`     | Next.js configuration. Sets Turbopack root to the parent directory for monorepo support.                                                         |
| `tsconfig.json`      | TypeScript configuration. Target ES2017, ESNext modules, bundler module resolution, strict mode. Defines the `@/*` path alias for clean imports. |
| `eslint.config.mjs`  | ESLint 9 flat config extending `next/core-web-vitals` and `next/typescript` rule sets.                                                           |
| `postcss.config.mjs` | PostCSS configuration using the `@tailwindcss/postcss` plugin for Tailwind CSS 4 processing.                                                     |
| `components.json`    | shadcn/ui configuration defining the style variant (new-york), color scheme (neutral), icon library (lucide), and import aliases.                |
| `declarations.d.ts`  | TypeScript module declarations for non-standard imports.                                                                                         |

## Code Conventions

### TypeScript

- Strict mode is enabled. All code must pass strict type checking.
- Interfaces are defined in `common/types.ts` for shared types and in `common/schemas.ts` for Zod-inferred types.
- Use the `@/*` path alias for absolute imports from the project root (e.g., `@/common/types`, `@/hooks/useAuth`).

### Components

- Use functional components with TypeScript.
- Client components that use hooks or browser APIs must include the `"use client"` directive.
- Prefer composition over prop drilling. Use React Query hooks for data instead of passing data down through many layers.

### API Services

- Each service module corresponds to a backend resource (auth, listing, user, brand).
- All requests use `withCredentials: true` for cookie-based session management.
- Service functions are typed with request and response generics.

### Error Handling

- Use the `handleError` utility from `common/utils.ts` in mutation `onError` callbacks.
- Display user-facing errors via Sonner toast notifications.
- Log unexpected errors to the console in development.

## Deployment

### Production Build

```bash
npm run build
npm run start
```

The build process runs Next.js optimization including:

- Automatic code splitting
- Tree shaking of unused code
- Static generation where possible
- Image optimization

### Vercel Deployment

The application is compatible with Vercel deployment. Set the `NEXT_PUBLIC_API_URL` environment variable in the Vercel dashboard to point to the production backend URL.

### Environment Considerations

- Ensure the backend server's CORS configuration includes the client's production domain in `CAMPUS_CART_SERVER_CORS_ALLOWED_ORIGINS`.
- The session cookie's domain must be configured via `CAMPUS_CART_AUTH_COOKIE_DOMAIN` on the backend to match the client's production domain.
- In production, the `cc_refresh_token` cookie is set with the `Secure` flag, requiring HTTPS.
