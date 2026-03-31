# kbd — TUI for beads (bd)

A terminal UI for the `bd` issue tracker, inspired by how k9s wraps kubectl.

## Concept

kbd shells out to `bd` CLI commands (with `--json` flag) and presents the results in an interactive terminal interface. It does **not** access the database directly — `bd` is the single source of truth.

---

## Data Access Strategy

All data flows through `bd` CLI with `--json` output:

```
kbd (TUI)  →  exec bd <command> --json  →  parse JSON  →  render to screen
```

**Core commands used:**

| View | bd command | Purpose |
|------|-----------|---------|
| Issue list | `bd list --json` | Main table view |
| Issue detail | `bd show <id> --json` | Detail pane |
| Status overview | `bd status --json` | Dashboard |
| Dependencies | `bd graph <id>` | Dependency graph |
| Children | `bd children <id> --json` | Child issues |
| Comments | `bd comments <id> --json` | Issue comments |
| Labels | `bd label list --json` | Label management |
| Search | `bd search <query> --json` | Full-text search |
| SQL | `bd sql <query> --json` | Raw query |

**Write operations:**

| Action | bd command |
|--------|-----------|
| Create | `bd create --title "..." --type ...` |
| Update | `bd update <id> --field value` |
| Close | `bd close <id>` |
| Reopen | `bd reopen <id>` |
| Add comment | `bd comments <id> add "..."` |
| Add label | `bd label add <id> <label>` |
| Add dep | `bd dep add <from> <to>` |

---

## Technology

| Choice | Rationale |
|--------|-----------|
| **Go** | Same as k9s; excellent for TUI + CLI exec |
| **tview/tcell** | Proven TUI framework (k9s uses it at scale) |
| **cobra** | CLI flags (--db, --readonly) pass through to bd |
| **os/exec** | Shell out to `bd` for all data operations |

---

## Architecture

```
kbd/
├── main.go                    # Entry point
├── cmd/
│   └── root.go                # Cobra root command, flags, app launch
├── internal/
│   ├── app/
│   │   └── app.go             # Main app: layout, navigation stack, key routing
│   ├── bd/
│   │   ├── client.go          # bd CLI executor: run commands, parse JSON
│   │   └── types.go           # Go structs matching bd --json output
│   ├── view/
│   │   ├── issues.go          # Issue list table with filtering + actions
│   │   └── detail.go          # Issue detail view (scrollable)
│   ├── ui/
│   │   ├── header.go          # Header: hints panel (left) + shortcuts panel (right) + status
│   │   ├── hints_panel.go     # Letter/nav key hints in columnar grid
│   │   ├── shortcuts.go       # Numbered data-driven shortcuts (0-9) in columnar grid
│   │   ├── statusbar.go       # Footer: per-type per-status counts
│   │   ├── table.go           # Table widget with vim keys (j/k/g/G)
│   │   ├── keybinds.go        # Key binding system with hints
│   │   ├── menu.go            # Generic menu bar (text hints)
│   │   ├── prompt.go          # Command (:) and filter (/) input bar
│   │   └── theme.go           # Color theme with status colors
│   └── model/
│       ├── issue.go           # Issue → table row conversion, age formatting
│       └── filter.go          # Multi-field filter (title, type, assignee, status)
├── Makefile                   # build, run, install, clean, fmt, vet
├── go.mod
└── go.sum
```

### Layer Responsibilities

**bd/** — CLI executor. Single point of contact with `bd`. Runs commands via `os/exec`, captures JSON output, unmarshals into Go structs. Handles `--db` passthrough and error parsing.

**model/** — Converts bd JSON into display-ready rows. Handles multi-field filtering (title, type, assignee, status), column definitions, and age formatting.

**view/** — One file per screen. Each view calls `bd.Client` methods to get data, renders via ui widgets, and defines its own key bindings.

**ui/** — Reusable widgets. Header (two-panel), table, prompt, status bar. Thin wrappers over tview.

**app/** — Orchestrates views. Stack-based navigation. Global key routing (number keys → shortcuts, `/` → search, `:` → type filter, Esc → pop filter stack).

---

## Screen Layout

```
┌──────────────────────────────┬──────────────────────────────┐
│  LEFT: Key Hints             │  RIGHT: Numbered Shortcuts   │
│  (letter/nav shortcuts       │  (data-driven, 0-9)          │
│   in columnar grid)          │   in columnar grid)          │
│                              │                              │
│  Enter  View      x  Close  │  0:kostine    5:epic           │
│  j/k    Navigate  o  Reopen │  1:claude   6:task            │
│  /      Search    r  Refresh│  2:open     7:application     │
│  :      Type      Esc Pop   │  3:closed   8:blocked         │
│  q      Quit                │  4:in_prog                    │
├──────────────────────────────┴──────────────────────────────┤
│  (status line — shown only during loading/errors)           │
├─────────────────────────────────────────────────────────────┤
│  (prompt line — shown when : or / is pressed)               │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Issue Table                                                │
│  ┌─────────────────────────────────────────────────────────┐│
│  │ ID    TYPE   STATUS   TITLE          ASSIGNEE  DEPS AGE ││
│  │ ...                                                     ││
│  └─────────────────────────────────────────────────────────┘│
│                                                             │
├─────────────────────────────────────────────────────────────┤
│  task  open:820 closed:45 │ epic  open:86 closed:12        │
│  application  open:18 researching:5 applied:8 closed:3     │
└─────────────────────────────────────────────────────────────┘
```

---

## Header

Two-panel layout, both rendering in columnar grids:

**Left panel (HintsPanel):** Static letter/navigation shortcuts. Updated when views change (e.g., detail view has different hints than list view). Rendered in 2 columns, filling rows top-to-bottom then left-to-right.

**Right panel (ShortcutPanel):** Dynamic numbered shortcuts (0-9). Computed from top-count values in the current data, grouped in order: assignees first, then types, then statuses. Each entry color-coded by category (green=assignee, aqua=type, yellow=status). Rendered in a 5×2 grid.

**Status line:** One row below the grid. Empty most of the time. Shows loading/error messages temporarily.

---

## Filtering

Filters stack and are peeled back with Esc, most recent first.

**Entry methods:**
- `/` — title search (fuzzy match on issue title)
- `:` — type filter (e.g., `:epic`, `:task`, `:all` to clear)
- `0-9` — quick filter from shortcut panel (assignee, type, or status)

**Filter stack example:**
1. Press `:` → type `epic` → Enter → filters to epics
2. Press `3` (shortcut for assignee:kostine) → adds assignee filter
3. Press `/` → type `coin` → Enter → adds title filter
4. Esc → removes `/coin` title filter
5. Esc → removes assignee filter
6. Esc → removes `:epic` type filter

**Title bar** reflects all active filters: `epic @kostine /coin (12)`

**Filter model (model/Filters):**
```go
type Filters struct {
    Title    string  // from /
    Type     string  // from : or shortcut
    Assignee string  // from shortcut
    Status   string  // from shortcut
}
```

---

## Footer (Status Bar)

One line per issue type that has data. Each line shows counts per status, color-coded:

```
task        open:820 in_progress:10 blocked:3 closed:45
epic        open:86 closed:12
application open:18 researching:5 applied:8 closed:3
```

Status colors: green=open, yellow=in_progress/hooked, red=blocked, gray=closed, blue=deferred, white=custom.

Counts always reflect the full dataset (not filtered). Bar resizes dynamically.

---

## Navigation

Stack-based. Views push/pop onto a page stack.

```
Enter          → Push detail view for selected issue
Esc            → Pop filter (if active) or pop view
q / Ctrl+C     → Quit (from root view)
```

**Issue list keys:**
```
j/k            → Navigate up/down
g/G            → Jump to top/bottom
Enter          → View issue detail
/              → Title search
:              → Type filter
0-9            → Quick filter (shortcut)
x              → Close issue
o              → Reopen issue
r              → Refresh data
```

**Issue detail keys:**
```
j/k            → Scroll up/down
r              → Refresh
Esc/q          → Back to list
```

---

## Views

### 1. Issue List (default)
Table with columns: ID, Type, Status, Title, Assignee, Deps, Age.
Powered by `bd list --json --limit 0`. Status-colored rows.

### 2. Issue Detail
Full issue view: all fields rendered as key-value pairs, description block below.
Powered by `bd show <id> --json --long`. Returns array; first element used.

---

## MVP Scope (v0.1) — Implemented

1. **bd client** — exec wrapper with JSON parsing
2. **Issue list view** — table with multi-field filtering
3. **Issue detail view** — drill-down on Enter
4. **Stack navigation** — push/pop with layered filter Esc
5. **Key bindings** — vim keys, `:` type filter, `/` search, `0-9` shortcuts
6. **Header** — two-panel layout with hints + data-driven shortcuts
7. **Footer** — per-type per-status summary counts
8. **Theme** — status-colored rows and UI elements

### Deferred to later
- Status dashboard
- Dependency graph visualization
- Write operations (create/edit)
- SQL view
- Comments view
- Custom themes/skins
- Plugin system
- Config file
