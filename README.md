# go-minitrackr — Minimal Self-Hosted Issue Tracker

## 1. Overview

go-minitrackr is a **minimal, self-hosted issue tracker inspired by Linear**, designed for:
- Solo developers or very small teams
- Extremely low memory usage
- Fast page loads and keyboard-first UX
- Simple deployment on a VPS

This project prioritizes **clarity, speed, and resource efficiency** over features.

Non-goals:
- Realtime collaboration
- Multi-workspace / multi-org support
- Plugin ecosystems
- Complex workflow automation

---

## 2. Core Features

### Issue Management
- Create issues with title, status, priority
- Edit issue titles inline
- Move issues between statuses
- Backlog (list) view
- Kanban board view

### UX (Linear-Inspired)
- Clean, dense UI
- System font stack
- Keyboard-first interactions:
  - `c` → create issue
  - `j / k` → navigate issues
  - `enter` → open issue
  - `esc` → close / blur
- Dark theme by default

### Interaction Model
- Full server-rendered HTML on first load
- Partial updates via HTMX
- No client-side routing
- No global JS state

---

## 3. Tech Stack

### Backend
- **Language:** Go (1.22+)
- **HTTP:** net/http (stdlib)
- **Templates:** html/template
- **Database:** SQLite
- **SQLite Driver:** modernc.org/sqlite (pure Go, no CGO)
- **Auth (optional):**
  - HTTP Basic Auth
  - Signed cookies (HMAC)

### Frontend
- **Rendering:** Server-side HTML
- **CSS:** Handwritten CSS (no framework runtime)
- **JavaScript:**
  - HTMX (~14 KB) for partial updates
  - Minimal custom JS (~1–2 KB) for keyboard shortcuts
- **Assets:** Embedded in Go binary

### Build & Runtime
- Single static Go binary
- No Node.js
- No frontend build step
- No background workers

---

## 4. Linear-Like UI (Without SPA)

### What Is Preserved
- Instant navigation feel
- Minimal visual noise
- Predictable keyboard interactions
- Kanban + list workflow
- Fast perceived performance

### What Is Intentionally Avoided
- React / Vue / Svelte
- Virtual DOM
- Client-side state management
- WebSockets
- Heavy animations

### UI Layout

Sidebar:
- Inbox
- Backlog
- Board
- Settings

Main Area:
- Kanban columns (Todo / Doing / Done)
- Issues rendered as compact rows
- Inline editing
- Drag-and-drop via HTMX + native HTML events

---

## 5. Performance Constraints

### Target Limits

| Resource | Constraint |
|-------|------------|
| RAM | **≤ 30 MB** |
| CPU | ~0–2% idle |
| Startup Time | < 50 ms |
| Binary Size | ~5–10 MB |

### Memory Budget Breakdown

| Component | Estimated RAM |
|--------|---------------|
| Go runtime | 5–8 MB |
| HTTP handlers & templates | 2–4 MB |
| SQLite cache | 2–3 MB |
| Embedded assets | ~1 MB |
| **Total** | **12–20 MB** |

### Hard Rules
- No ORMs
- No in-memory caches
- SQLite WAL mode with small page cache
- GOMEMLIMIT set to 25 MiB

---

## 6. JavaScript Policy

### app.js (Custom)
- Keyboard shortcuts
- Focus / blur helpers
- No state
- No reactivity
- < 5 KB unminified

### htmx.min.js
- Used for:
  - Issue creation
  - Inline edits
  - Status changes
- Zero server-side memory impact
- No framework overhead

---

## 7. Database Design

- SQLite (single file)
- WAL mode enabled
- Integer primary keys
- Minimal indexes
- Optimized for reads

Future-proofing:
- Schema can be migrated to PostgreSQL later if needed

---

## 8. Deployment (Dokploy)

### Why It Fits Well
- One container
- One exposed port
- One volume
- No sidecars (Redis, Mongo, etc.)

### Dokploy Configuration
- Service type: Docker
- Port: 3000
- Volume mount: /data (SQLite DB)
- Health check: GET /health
- Memory limit: 64 MB (safe cap)

### Docker Characteristics
- Multi-stage build
- Scratch final image
- CGO disabled
- Static binary

---

## 9. Expected Runtime Metrics

| Scenario | RAM |
|-------|-----|
| Idle | 10–15 MB |
| Light usage | 18–22 MB |
| Peak (few users) | < 25 MB |

CPU usage remains near zero when idle.

---

## 10. Philosophy

This project follows a simple rule:

> **Every feature must justify its memory cost.**

go-minitrackr is not a clone of Linear's tech stack —  
it is a recreation of its *feel* using fundamentally simpler tools.

---

## 11. Possible Future Enhancements (Still Within Budget)

- Issue detail modal (`<dialog>`)
- Theme toggle (CSS variables)
- Basic auth middleware
- Keyboard command palette
- PostgreSQL support (optional)

---

## 12. Summary

- Linear-like UX
- < 30 MB RAM
- Single Go binary
- SQLite persistence
- HTMX-powered interactivity
- Dokploy-friendly deployment

**Small, fast, and boring — by design.**
