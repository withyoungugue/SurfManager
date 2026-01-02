# Changelog

All notable changes to SurfManager will be documented in this file.

## [2.0.1] - 2026-01-03

### 🔧 Fixes & Improvements

**GitHub Actions Workflow**
- Fixed invalid workflow file error on `build.yml`
- Resolved `matrix` context not accessible at job-level `if` condition
- Added step-level platform filtering for manual workflow dispatch
- Now properly supports selective platform builds (windows/linux/macos)

**Settings UX Improvement**
- Renamed "Skip Close App" to "Keep App Running" for better clarity
- Updated description to be more intuitive
- Toggle ON = app stays running during backup/restore/reset
- Toggle OFF = app will be closed before operations

---

## [2.0.0] - 2025-12-06

### 🚀 Complete Rewrite - Go + Wails + Svelte

**SurfManager v2.0** is a complete rewrite from Python/PyQt6 to Go + Wails + Svelte + TailwindCSS for better performance, smaller binary size, and modern UI.

### ✨ New Features

**🎨 Theme System**
- Dark theme (default)
- Solarized Dark theme
- Solarized Light theme
- Theme persistence across sessions

**⚙️ Settings System**
- Comprehensive settings with categories: Appearance, Behavior, Sessions, Notes, Advanced
- All settings persistent via localStorage
- New "Skip Close App" option - perform operations without closing target app
- Auto-backup toggle
- Confirmation dialogs (customizable)
- Remember last active tab

**🖥️ UI Improvements**
- JetBrains Mono font for better readability
- Realtime clock display in header
- Custom confirmation modals (no more browser alerts)
- Right-click context menus throughout the app
- CTRL+Click multi-selection in Sessions tab
- Split panel design for Reset tab
- Grid layout for Add App dialog
- Text selection disabled (native app feel)

**📝 Notes Tab**
- Create and manage notes
- Auto-save option
- Markdown support

**🔧 Config Tab Enhancements**
- VSCode Preset vs Custom app type
- Customizable backup items (choose what to backup)
- Additional folders support (backup extra directories)
- Smart file dialogs with default directories:
  - Executable: Opens from `AppData/Local/Programs`
  - Data Folder: Opens from `AppData/Roaming`
  - Additional Folders: Opens from user home
- Right-click context menu (Set Active, Edit JSON, Open Folder, Delete)

**📊 Sessions Tab Enhancements**
- Index number column
- App filter dropdown
- CTRL+Click selection (no checkboxes)
- Right-click context menu (Restore, Set Active, Open Folder, Delete)
- Icon + text action buttons

**🔄 Reset Tab Redesign**
- Split panel layout (Pro style)
- Left panel: App list with status indicators
- Right panel: App details, actions, stats
- Session count, Last Reset, Auto-Backup status
- Add App button navigates to Config tab

### 🛠️ Technical Changes

- **Framework**: Go + Wails v2 (from Python + PyQt6)
- **Frontend**: Svelte + TailwindCSS + Vite
- **Binary Size**: ~15MB (from 40+MB)
- **Startup Time**: <0.5s (from ~1s)
- **Memory Usage**: Significantly reduced
- **Cross-platform**: Windows, macOS, Linux support

### 📦 Migration Notes

- Settings are stored in localStorage (browser-based)
- App configs stored in `~/.surfmanager/AppConfigs/`
- Backups stored in `Documents/SurfManager/backup/`
- Notes stored in `Documents/SurfManager/notes/`

---

## [1.0.1] - 2025-12-05

### Fixed
- **Backup Path Structure**: Fixed incorrect backup directory naming
  - Changed from `SurfManagerBackups` to `Documents/SurfManager/backup`
  - Auto-backups now stored in `Documents/SurfManager/auto-backup`
  - Notes stored in `Documents/SurfManager/notes`

### Improved
- **Auto-Backup Toggle Button**: Better UX with badge counter
- **Reset Data Tab Sync**: Fixed app list not syncing with App Configuration
- **Session Scanning**: Skip internal folders when listing backups

---

## [1.0.0] - 2025-12-04

### Highlights
**New Redesign, User Friendly, and Fast - No More Freezing!**

### Added
- **Background Threading**: No more UI freezing during operations
- **Cross-Platform Support**: Windows, Linux, and macOS
- **Optimized Backup Items**: Reduced to 9 essential files/folders
- **Duplicate Protection**: Prevents overwriting existing sessions

### Fixed
- Variable name errors in backup/restore functions
- Filter combo changing after backup
- Session active status not displaying correctly

---

## [0.0.3-beta-windows] - 2025-12-04

### Added
- **Addon Backup Folders**: Optional additional folders to backup
- **Smart App Close**: Graceful termination before operations
- **Auto-Backup Protection**: Automatic backup before reset
- **Dual View Mode**: Separate views for manual sessions vs auto-backups

---

## [0.0.2-beta] - 2025-12-04

### Added
- **Session Backup System**: Full backup/restore functionality
- **Progress bar** in Reset Data tab

---

## [0.0.1-beta] - 2025-12-04

### Initial Release
- GUI foundation (PyQt6)
- Reset Data tab with program list
- Sessions and App Configuration tabs
- Dark theme UI
