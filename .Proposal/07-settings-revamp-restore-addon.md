# Settings Revamp & Restore Addon Only Feature

**Status**: Proposal  
**Date**: 2026-01-04  
**Author**: Kiro AI Assistant

---

## 📋 Overview

Revamp settings structure dan tambah fitur "Restore Addon Only" untuk restore hanya addon folders dari backup session.

---

## 🎯 Goals

1. **Simplify Settings** - Merge kategori yang isinya sedikit
2. **Add Import/Export Settings** - Backup/restore konfigurasi SurfManager
3. **Add Restore Addon Only** - Restore hanya addon folders (experimental)

---

## 📐 Settings Structure (New)

### **Current Structure (6 Categories)**
```
├── Appearance     → Theme only
├── Behavior       → Few options
├── Sessions       → Few options
├── Notes          → ???
├── Advanced       → ???
└── Experimental   → ???
```

### **Proposed Structure (4 Categories)**
```
├── 🎨 General
├── ⚙️ Behavior
├── 📦 Sessions
└── 🧪 Experimental
```

---

## 🎨 General Settings

| Setting | Type | Default | Description |
|---------|------|---------|-------------|
| **Theme** | Toggle (Dark/Auto/Light) | Dark | Color scheme |
| **Remember Last Tab** | Toggle | ON | Remember last opened tab on startup |
| **Import Settings** | Button | - | Import settings from .json file |
| **Export Settings** | Button | - | Export settings to .json file |
| **Reset Settings** | Button | - | Reset all settings to default |

### Import/Export Settings Feature

**Export Format** (`.surfmanager-settings.json`):
```json
{
  "version": "2.0.1",
  "exported_at": "2026-01-04T16:50:00Z",
  "settings": {
    "theme": "dark",
    "rememberLastTab": true,
    "confirmBeforeReset": true,
    "confirmBeforeDelete": true,
    "autoBackup": true,
    "skipCloseApp": false,
    "showAutoBackups": true,
    "maxAutoBackups": 5,
    "experimental": {
      "showRestoreAddonOnly": false
    }
  }
}
```

**Import Behavior**:
- Validate JSON structure
- Show preview of changes
- Confirm before applying
- Backup current settings before import

---

## ⚙️ Behavior Settings

| Setting | Type | Default | Description |
|---------|------|---------|-------------|
| **Confirm before Reset** | Toggle | ON | Show confirmation dialog before reset |
| **Confirm before Delete** | Toggle | ON | Show confirmation dialog before delete |
| **Auto-Backup on Reset** | Toggle | ON | Create auto-backup before reset |
| **Skip Close App** | Toggle | OFF | Don't close app during backup/reset |

---

## 📦 Sessions Settings

| Setting | Type | Default | Description |
|---------|------|---------|-------------|
| **Show Auto-backups** | Toggle | ON | Display auto-backups in session list |
| **Max Auto-backups** | Number (1-10) | 5 | Maximum auto-backups per app |
| **Open Backup Folder** | Button | - | Open backup folder in explorer |
| **Clean Old Auto-backups** | Button | - | Delete all auto-backups older than X days |

---

## 🧪 Experimental Settings

| Setting | Type | Default | Description |
|---------|------|---------|-------------|
| **Show Restore Addon Only** | Toggle | OFF | Enable "Restore Addon Only" in context menu |
| **Debug Mode** | Toggle | OFF | Show debug logs in console |
| **Enable Beta Features** | Toggle | OFF | Enable experimental features |

---

## 🆕 Feature: Restore Addon Only

### Overview
Restore ONLY addon folders from a backup session, without touching the main data folder.

### Use Case
User wants to restore `.aws` credentials from backup without restoring entire VSCode settings.

### Requirements
1. **Setting enabled**: `Experimental > Show Restore Addon Only = ON`
2. **App has addon paths**: `addon_backup_paths` not empty
3. **Session has addons**: `_addons` folder exists in backup

### UI Flow

**Right-click Context Menu (Sessions Tab)**:
```
┌─────────────────────────────┐
│ 📦 Restore Session          │
│ 👤 Restore Account Only     │
│ 📁 Restore Addon Only       │ ← NEW (conditional)
│ ─────────────────────────── │
│ ✓ Set as Active             │
│ 📂 Open Folder              │
│ ─────────────────────────── │
│ 🗑️ Delete Session           │
└─────────────────────────────┘
```

**Visibility Logic**:
```javascript
showRestoreAddonOnly = 
  settings.experimental.showRestoreAddonOnly === true &&
  app.addon_backup_paths.length > 0 &&
  sessionHasAddonsFolder(session)
```

### Backend Implementation

**New Function**: `RestoreAddonOnly(appKey, sessionName, skipClose)`

```go
func (a *App) RestoreAddonOnly(appKey, sessionName string, skipClose bool) error {
    cfg := a.GetApp(appKey)
    if cfg == nil {
        return fmt.Errorf("app not found: %s", appKey)
    }

    if len(cfg.AddonPaths) == 0 {
        return fmt.Errorf("no addon paths configured")
    }

    // Close app if needed
    if !skipClose {
        // ... close logic
    }

    // Restore ONLY addons from backup
    backupFolder := filepath.Join(m.backupPath, appKey, sessionName)
    return m.restoreAddons(backupFolder, cfg.AddonPaths, progressCb)
}
```

### Confirmation Dialog

```
┌─────────────────────────────────────────┐
│  Restore Addon Folders Only             │
├─────────────────────────────────────────┤
│                                         │
│  Restore addon folders from:            │
│  "session-name"                         │
│                                         │
│  This will restore:                     │
│  • C:\Users\User\.aws                   │
│  • C:\Users\User\.ssh                   │
│                                         │
│  Main data folder will NOT be touched.  │
│                                         │
│  ┌─────────┐  ┌─────────┐              │
│  │ Cancel  │  │ Restore │              │
│  └─────────┘  └─────────┘              │
└─────────────────────────────────────────┘
```

---

## 🔧 Implementation Plan

### Phase 1: Settings Revamp
1. ✅ Merge settings categories (4 instead of 6)
2. ✅ Add Import/Export Settings buttons
3. ✅ Implement settings JSON export
4. ✅ Implement settings JSON import with validation
5. ✅ Add "Show Restore Addon Only" toggle in Experimental

### Phase 2: Restore Addon Only
1. ✅ Add backend function `RestoreAddonOnly`
2. ✅ Add frontend context menu item (conditional)
3. ✅ Add confirmation dialog
4. ✅ Test with various scenarios

### Phase 3: Polish
1. ✅ Add tooltips for experimental features
2. ✅ Add warning icons for experimental toggles
3. ✅ Update documentation

---

## 🎨 UI Mockups

### Settings Page (General)
```
┌─────────────────────────────────────────────────────────┐
│  SETTINGS                                               │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ┌──────────────┐  ┌─────────────────────────────────┐ │
│  │ 🎨 General   │  │  Theme                          │ │
│  ├──────────────┤  │  Choose your color scheme       │ │
│  │ ⚙️ Behavior  │  │  ┌──────┐ ┌──────┐ ┌──────┐   │ │
│  ├──────────────┤  │  │ Dark │ │ Auto │ │Light │   │ │
│  │ 📦 Sessions  │  │  └──────┘ └──────┘ └──────┘   │ │
│  ├──────────────┤  │                                 │ │
│  │ 🧪 Experim.  │  │  ☑ Remember Last Tab           │ │
│  └──────────────┘  │                                 │ │
│                    │  ─────────────────────────────  │ │
│                    │                                 │ │
│                    │  Settings Management            │ │
│                    │  ┌──────────────┐              │ │
│                    │  │ Import       │              │ │
│                    │  └──────────────┘              │ │
│                    │  ┌──────────────┐              │ │
│                    │  │ Export       │              │ │
│                    │  └──────────────┘              │ │
│                    │  ┌──────────────┐              │ │
│                    │  │ Reset All    │              │ │
│                    │  └──────────────┘              │ │
│                    └─────────────────────────────────┘ │
└─────────────────────────────────────────────────────────┘
```

### Settings Page (Experimental)
```
┌─────────────────────────────────────────────────────────┐
│  SETTINGS                                               │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ┌──────────────┐  ┌─────────────────────────────────┐ │
│  │ 🎨 General   │  │  Experimental Features          │ │
│  ├──────────────┤  │  ⚠️ Use at your own risk        │ │
│  │ ⚙️ Behavior  │  │                                 │ │
│  ├──────────────┤  │  ☐ Show Restore Addon Only     │ │
│  │ 📦 Sessions  │  │    Enable "Restore Addon Only"  │ │
│  ├──────────────┤  │    in session context menu      │ │
│  │ 🧪 Experim.  │◄─│                                 │ │
│  └──────────────┘  │  ☐ Debug Mode                   │ │
│                    │    Show debug logs              │ │
│                    │                                 │ │
│                    │  ☐ Enable Beta Features         │ │
│                    │    Unlock experimental features │ │
│                    └─────────────────────────────────┘ │
└─────────────────────────────────────────────────────────┘
```

---

## 🧪 Testing Scenarios

### Restore Addon Only
1. ✅ App with addon paths + session with _addons → Show menu item
2. ✅ App without addon paths → Hide menu item
3. ✅ Session without _addons folder → Hide menu item
4. ✅ Setting disabled → Hide menu item
5. ✅ Restore only addons, verify data folder untouched
6. ✅ Handle missing addon paths gracefully

### Import/Export Settings
1. ✅ Export settings to JSON
2. ✅ Import valid JSON
3. ✅ Import invalid JSON → Show error
4. ✅ Import with missing fields → Use defaults
5. ✅ Preview changes before import

---

## 📝 Notes

- **Backward Compatibility**: Old settings format should still work
- **Migration**: Auto-migrate old settings to new structure
- **Validation**: Validate all imported settings before applying
- **Backup**: Always backup current settings before import

---

## ✅ Acceptance Criteria

- [ ] Settings page has 4 categories (General, Behavior, Sessions, Experimental)
- [ ] Import/Export settings works correctly
- [ ] "Restore Addon Only" appears in context menu when conditions met
- [ ] "Restore Addon Only" restores ONLY addon folders
- [ ] All existing settings still work after revamp
- [ ] Settings are persisted correctly
- [ ] UI is clean and intuitive

---

## 🚀 Future Enhancements

1. **Cloud Sync Settings** - Sync settings across devices
2. **Settings Profiles** - Multiple setting profiles (Work, Personal, etc.)
3. **Scheduled Auto-backups** - Automatic backup on schedule
4. **Backup Compression** - Compress backups to save space
5. **Backup Encryption** - Encrypt sensitive backups

---

## 📚 Related Files

- `frontend/src/lib/SettingsTab.svelte` - Settings UI
- `frontend/src/lib/SessionsTab.svelte` - Context menu
- `frontend/src/lib/stores/settings.js` - Settings store
- `internal/config/config.go` - Settings backend
- `app.go` - RestoreAddonOnly function

---

**End of Proposal**
