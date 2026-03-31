# kbd

A terminal UI for the [beads (bd)](https://github.com/derailed/beads) issue tracker, inspired by [k9s](https://github.com/derailed/k9s).

kbd wraps the `bd` CLI and presents issues in an interactive TUI with filtering, sorting, shortcuts, and tmux integration.

## Prerequisites

- **[bd](https://github.com/derailed/beads)** — beads CLI (required). kbd shells out to `bd` for all data operations.
- **tmux** — optional. Enables `w` and `c` commands to send messages to an adjacent pane.
- **[tmuxinator](https://github.com/tmuxinator/tmuxinator)** — optional. Useful for setting up a two-pane layout (kbd + Claude) with a single command.

Example `~/.tmuxinator/beads.yml`:

```yaml
name: beads
root: ~/Projects/my-project

windows:
  - main:
      layout: even-horizontal
      panes:
        - claude
        - kbd
```

Start with `tmuxinator start beads`.

## Install

```bash
brew tap kostine/tap
brew install kbd
```

Or build from source:

```bash
git clone https://github.com/kostine/kbd.git
cd kbd
make install
```

## Usage

```bash
kbd                          # auto-detect .beads in current dir
kbd --db path/to/.beads/dolt # specify database path
kbd version                  # print version info
```

On first run without a database, kbd shows a folder picker to browse and select a beads directory. Selected paths are saved to `~/.kbd/contexts.json` for quick access.

## Screen Layout

```
┌─────────────────────────────┬─────────────────────────────┐
│  Key Hints (left)           │  Numbered Shortcuts (right)  │
│  Enter  View    x  Close    │  0:kostine    5:epic           │
│  j/k    Nav     o  Reopen   │  1:claude   6:task           │
│  /      Search  r  Refresh  │  2:open     7:context        │
│  :      Type    h  Help     │  3:closed                    │
├─────────────────────────────┴─────────────────────────────┤
│ ID              TYPE PRI STATUS      TITLE         AGE     │
│ bd-123          epic P1  open        My project    2d      │
│ bd-456          task P2  in-progress Fix login     5h      │
├───────────────────────────────────────────────────────────┤
│ task open:820 closed:45 │ epic open:86 closed:12          │
└───────────────────────────────────────────────────────────┘
```

## Key Bindings

### Navigation
| Key | Action |
|-----|--------|
| `j` / `k` | Move down / up |
| `g` / `G` | Jump to top / bottom |
| `Enter` | View issue detail |
| `Esc` | Pop filter or go back |
| `q` | Quit |

### Filtering
| Key | Action |
|-----|--------|
| `/` | Search by title |
| `:` | Filter by type (`epic`, `task`, `bug`, `all`) |
| `0-9` | Quick filter from shortcut panel |

Filters stack. Press `Esc` to remove them one at a time.

### Actions
| Key | Action |
|-----|--------|
| `x` | Close issue |
| `o` | Reopen issue |
| `r` | Refresh data |
| `h` | Show help screen |

### Tmux Integration (auto-detected)
| Key | Action |
|-----|--------|
| `w` | Send "lets work on \<id\>" to adjacent pane |
| `c` | Send a custom message to adjacent pane |

### Commands (via `:`)
| Command | Action |
|---------|--------|
| `:epic` | Show epics |
| `:task` | Show tasks |
| `:all` | Show all types |
| `:context` | Switch database |
| `:q` | Quit |

## Features

- **Auto-refresh** — data refreshes every 5 seconds with delta highlighting (green = new, yellow = changed)
- **Smart sorting** — by status (in-progress first), priority (P0 highest), then type-specific rules
- **Epic view** — shows child issue progress (closed/total) in an ISSUES column
- **Context management** — browse folders, save databases to `~/.kbd/contexts.json`, switch with `:context`
- **Tmux integration** — send issue IDs or messages to an adjacent pane (e.g., Claude)
- **Data-driven shortcuts** — top assignees, types, and statuses auto-populate the 0-9 shortcut panel

## Configuration

| Path | Purpose |
|------|---------|
| `~/.kbd/contexts.json` | Saved database contexts |
| `~/.kbd/kbd.log` | Application log (errors, panics) |

## Built With

- [Go](https://go.dev/)
- [tview](https://github.com/rivo/tview) / [tcell](https://github.com/gdamore/tcell) — terminal UI
- [cobra](https://github.com/spf13/cobra) — CLI framework
- [bd](https://github.com/derailed/beads) — beads issue tracker

## License

MIT
