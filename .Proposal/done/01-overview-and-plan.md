# SurfManager v2.0 - Go + Wails Migration

## Overview

Migrasi SurfManager dari Python/PyQt6 ke Go + Wails untuk performa lebih baik, single binary, dan UI modern.

### Why Migrate?

| Python (Old) | Go + Wails (New) |
|--------------|------------------|
| ~40MB binary + Python runtime | ~15-20MB single binary |
| Startup ~1-2s | Startup <500ms |
| PyQt6 dependency hell | Zero dependencies |
| Complex build process | Simple `wails build` |
| Limited UI customization | Full web UI (HTML/CSS/JS) |

---

## Tech Stack

```
Backend:  Go 1.21+
Frontend: Svelte + TailwindCSS (lightweight & fast)
Build:    Wails v2
Icons:    Lucide Icons (modern, lightweight)
```

---

## Project Structure (Minimal & Powerful)

```
SurfManager/
├── main.go                 # Entry point
├── app.go                  # Wails app struct & methods
├── wails.json              # Wails config
├── go.mod
│
├── internal/
│   ├── config/
│   │   └── config.go       # Config manager + platform paths
│   ├── process/
│   │   └── killer.go       # Process management
│   ├── backup/
│   │   └── backup.go       # Backup/restore logic
│   └── apps/
│       └── loader.go       # App config loader (JSON)
│
├── frontend/
│   ├── src/
│   │   ├── App.svelte      # Main app
│   │   ├── main.js
│   │   ├── app.css         # TailwindCSS
│   │   └── lib/
│   │       ├── Sidebar.svelte
│   │       ├── ResetTab.svelte
│   │       ├── SessionsTab.svelte
│   │       ├── ConfigTab.svelte
│   │       └── NotesTab.svelte
│   ├── index.html
│   ├── package.json
│   └── tailwind.config.js
│
└── .old_python/            # Reference (keep for now)
```

**Total: ~15 files** (vs Python's 20+ files)

---

## Feature Mapping

### Core Features (Must Have)

| Feature | Python | Go Implementation |
|---------|--------|-------------------|
| Reset Data | ✅ | `internal/backup/reset.go` |
| Session Backup | ✅ | `internal/backup/backup.go` |
| Session Restore | ✅ | `internal/backup/restore.go` |
| Auto-Backup | ✅ | Built into reset flow |
| Process Kill | ✅ | `internal/process/killer.go` |
| App Config | ✅ | `internal/apps/loader.go` |
| Platform Paths | ✅ | `internal/config/paths.go` |
| Notepad | ✅ | Simple JSON storage |

### UI Features

| Feature | Python | Wails Implementation |
|---------|--------|----------------------|
| Dark Theme | PyQt6 CSS | TailwindCSS dark mode |
| Tab Navigation | QTabWidget | Svelte components |
| Progress Bar | QProgressBar | HTML progress + events |
| Context Menu | QMenu | Custom dropdown |
| Search/Filter | QLineEdit | Svelte reactive |
| Multi-select | QTableWidget | Checkbox list |

---

## Migration Phases

### Phase 1: Foundation (Week 1)
- [x] Setup Git repo
- [ ] Initialize Wails project
- [ ] Create Go backend structure
- [ ] Implement platform paths
- [ ] Basic frontend shell

### Phase 2: Core Features (Week 2)
- [ ] Process killer
- [ ] App config loader
- [ ] Backup/restore logic
- [ ] Reset functionality

### Phase 3: UI (Week 3)
- [ ] Reset Data tab
- [ ] Sessions tab
- [ ] App Config tab
- [ ] Notes tab

### Phase 4: Polish (Week 4)
- [ ] Progress events
- [ ] Error handling
- [ ] Testing
- [ ] Build & release

---

## Go Backend API

```go
// app.go - Wails bindings

type App struct {
    ctx     context.Context
    config  *config.Manager
    process *process.Killer
    backup  *backup.Manager
    apps    *apps.Loader
}

// Reset Data
func (a *App) ResetApp(appKey string) error
func (a *App) GenerateNewID(appKey string) (int, error)
func (a *App) LaunchApp(appKey string) error
func (a *App) OpenFolder(appKey string) error

// Sessions
func (a *App) GetSessions(appKey string) ([]Session, error)
func (a *App) CreateBackup(appKey, name string) error
func (a *App) RestoreBackup(appKey, name string) error
func (a *App) DeleteSession(appKey, name string) error
func (a *App) SetActiveSession(appKey, name string) error

// App Config
func (a *App) GetApps() ([]AppConfig, error)
func (a *App) AddApp(config AppConfig) error
func (a *App) RemoveApp(appKey string) error
func (a *App) ToggleApp(appKey string) error

// Notes
func (a *App) GetNotes() ([]Note, error)
func (a *App) SaveNote(note Note) error
func (a *App) DeleteNote(id string) error

// Utils
func (a *App) GetPlatformInfo() PlatformInfo
func (a *App) KillProcess(names []string) error
```

---

## Data Compatibility

**100% compatible** dengan data Python version:
- Same folder structure: `Documents/SurfManager/`
- Same JSON format untuk app configs
- Same session storage format
- Same notes format

User bisa langsung upgrade tanpa kehilangan data!

---

## Performance Targets

| Metric | Python | Go Target |
|--------|--------|-----------|
| Startup | ~1.5s | <500ms |
| Memory | ~80MB | <30MB |
| Binary | ~40MB | <20MB |
| Backup 100MB | ~5s | <2s |

---

## Next Steps

1. **Run**: `wails init -n SurfManager -t svelte`
2. **Setup**: TailwindCSS + Lucide icons
3. **Implement**: Backend Go modules
4. **Build**: Frontend components
5. **Test**: Cross-platform
6. **Release**: Single binary

---

*Light as fuck, powerful as hell* 🚀
