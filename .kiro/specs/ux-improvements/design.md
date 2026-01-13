# Design Document

## Overview

This document describes the technical design for implementing UX improvements in SurfManager v2.0. The improvements span across Reset Tab, Session Tab, and Settings Tab, focusing on better user experience through alphabetical sorting, state persistence, enhanced restore workflows, and cleanup of unused features.

## Architecture

The implementation follows the existing Wails + Svelte architecture:
- Frontend: Svelte components with reactive stores
- Backend: Go functions exposed via Wails bindings
- State: localStorage for frontend persistence, file system for backend data

### Component Interaction

```
┌─────────────────────────────────────────────────────────────┐
│                      Frontend (Svelte)                       │
├─────────────────────────────────────────────────────────────┤
│  ResetTab.svelte    SessionsTab.svelte    SettingsTab.svelte │
│       │                    │                     │           │
│       └────────────────────┼─────────────────────┘           │
│                            │                                 │
│                    stores/settings.js                        │
│                    (localStorage persistence)                │
└─────────────────────────────────────────────────────────────┘
                             │
                      Wails Bindings
                             │
┌─────────────────────────────────────────────────────────────┐
│                       Backend (Go)                           │
├─────────────────────────────────────────────────────────────┤
│  app.go (binding layer)                                      │
│       │                                                      │
│       ├── internal/apps/loader.go                            │
│       ├── internal/backup/backup.go                          │
│       └── internal/config/config.go                          │
└─────────────────────────────────────────────────────────────┘
```

## Components and Interfaces

### 1. Alphabetical Sorting (Frontend)

**Location**: `frontend/src/lib/ResetTab.svelte`, `frontend/src/lib/SessionsTab.svelte`

**Implementation**:
```javascript
// Sort apps alphabetically by display_name
apps = apps.sort((a, b) => 
  a.display_name.localeCompare(b.display_name)
);
```

### 2. Remember Last Selected App (Frontend)

**Location**: `frontend/src/lib/stores/settings.js`

**New Properties**:
```javascript
const defaultSettings = {
  // ... existing settings
  lastSelectedAppReset: '',    // Last selected app in Reset Tab
  lastSelectedAppSession: '',  // Last selected app filter in Session Tab
};
```

**Component Updates**:
- ResetTab: Load `lastSelectedAppReset` on mount, save on app selection
- SessionsTab: Load `lastSelectedAppSession` on mount, save on filter dropdown change
- Both tabs persist selection across tab navigation (not just on app close)

### 3. Launch App from Context Menu (Frontend)

**Location**: `frontend/src/lib/SessionsTab.svelte`

**Context Menu Addition**:
```javascript
// Add to context menu items
<button on:click={() => handleLaunchApp(contextMenu.session)}>
  <Play size={14} />
  Launch App
</button>
```

**Handler**:
```javascript
async function handleLaunchApp(session) {
  try {
    await LaunchApp(session.app);
  } catch (e) {
    // Show error toast
  }
}
```

### 4. Post-Restore Launch Prompt (Frontend)

**Location**: `frontend/src/lib/SessionsTab.svelte`

**New State**:
```javascript
let showLaunchPrompt = false;
let launchPromptApp = '';
```

**Modal Component**:
```svelte
{#if showLaunchPrompt}
  <div class="modal">
    <h3>Restore Complete!</h3>
    <p>Would you like to launch {launchPromptApp}?</p>
    <button on:click={handleLaunchFromPrompt}>Launch App</button>
    <button on:click={() => showLaunchPrompt = false}>Close</button>
  </div>
{/if}
```

### 5. Set Active on Addon/Account Restore (Backend + Frontend)

**Backend Location**: `app.go`

**Modification to `RestoreAddonOnly`**:
```go
func (a *App) RestoreAddonOnly(appKey, sessionName string, skipClose bool) error {
    // ... existing restore logic
    
    // Set as active session after successful restore
    if err == nil {
        a.backup.SetActiveSession(appKey, sessionName)
    }
    
    return err
}
```

**Modification to `RestoreAccountOnly`**:
```go
func (a *App) RestoreAccountOnly(appKey, sessionName string) error {
    // ... existing restore logic
    
    // Set as active session after successful restore
    if err == nil {
        a.backup.SetActiveSession(appKey, sessionName)
    }
    
    return err
}
```

### 6. Auto-Generate New ID After Restore (Backend + Frontend)

**Backend Location**: `app.go`

**Modification to all restore functions**:
```go
func (a *App) RestoreBackup(appKey, sessionName string, skipClose bool) error {
    // ... existing restore logic
    
    if err == nil {
        // Generate new IDs after restore
        count, _ := a.GenerateNewID(appKey)
        wailsRuntime.EventsEmit(a.ctx, "progress", map[string]interface{}{
            "percent": 100,
            "message": fmt.Sprintf("Restore complete! Updated %d ID(s)", count),
        })
    }
    
    return err
}
```

### 7. Remove Unused Settings (Frontend)

**Location**: `frontend/src/lib/stores/settings.js`, `frontend/src/lib/SettingsTab.svelte`

**Settings Store Changes**:
```javascript
const defaultSettings = {
  // REMOVE these:
  // rememberLastTab: true,
  // lastActiveTab: 'reset',
  // debugMode: false,
  // skipDataFolder: false,
  
  // ... keep other settings
};
```

**SettingsTab Changes**:
- Remove "Remember Last Tab" toggle from General section
- Remove "Debug Mode" toggle from Experimental section
- Remove "Skip Data Folder" toggle from Experimental section

### 8. Bulk Session Management (Frontend + Backend)

**Backend Location**: `app.go`

**New Functions**:
```go
// ClearAllSessions deletes all backup sessions for all apps
func (a *App) ClearAllSessions() error {
    apps := a.GetActiveApps()
    for _, app := range apps {
        sessions, _ := a.backup.GetSessions(app.AppName, false)
        for _, session := range sessions {
            a.backup.DeleteSession(app.AppName, session.Name)
        }
    }
    return nil
}

// BackupAllSessions creates a zip archive of all sessions
func (a *App) BackupAllSessions() (string, error) {
    // Create zip archive of backup folder
    // Return path to created archive
}
```

**Frontend Location**: `frontend/src/lib/SettingsTab.svelte`

**New UI Section**:
```svelte
<div class="grid grid-cols-2 gap-2">
  <button on:click={handleBackupAllSessions}>
    <Download size={16} />
    Backup All Sessions
  </button>
  <button on:click={handleClearAllSessions}>
    <Trash2 size={16} />
    Clear All Sessions
  </button>
</div>
```

### 9. Remove Skip Data Folder Feature (Backend)

**Backend Location**: `internal/apps/loader.go`

**Changes**:
- Remove `SkipDataFolder` field from `AppConfig` struct
- Update JSON unmarshaling to ignore the field if present in existing configs

### 10. Reorganize Reset Tab Layout (Frontend)

**Location**: `frontend/src/lib/ResetTab.svelte`

**New App Row Layout**:
```svelte
<button class="app-row">
  <span class="status-dot"></span>
  <span class="app-name">{app.display_name}</span>
  <div class="badges">
    <span class="badge">{sessionCount} sessions</span>
    <span class="badge">{autoBackupCount} auto</span>
    {#if addonCount > 0}
      <span class="badge">{addonCount} addons</span>
    {/if}
  </div>
</button>
```

## Data Models

### Settings Store (Updated)

```javascript
const defaultSettings = {
  // General
  theme: 'dark',
  
  // Remember Selected Apps (NEW)
  lastSelectedAppReset: '',
  lastSelectedAppSession: '',
  
  // Behavior
  confirmBeforeReset: true,
  confirmBeforeDelete: true,
  confirmBeforeRestore: true,
  autoBackup: true,
  skipCloseApp: false,
  
  // Sessions
  showAutoBackups: false,
  defaultSessionFilter: 'all',
  maxAutoBackups: 5,
  
  // Experimental
  showRestoreAddonOnly: false,
  experimentalRestoreAccountOnly: false,
  
  // REMOVED:
  // rememberLastTab, lastActiveTab, debugMode, skipDataFolder
};
```

### AppConfig (Updated)

```go
type AppConfig struct {
    AppName        string       `json:"app_name"`
    DisplayName    string       `json:"display_name"`
    Version        string       `json:"version"`
    Active         bool         `json:"active"`
    Description    string       `json:"description"`
    Paths          AppPaths     `json:"paths"`
    BackupItems    []BackupItem `json:"backup_items"`
    AddonPaths     []string     `json:"addon_backup_paths"`
    // REMOVED: SkipDataFolder bool
}
```


## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system-essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

### Property 1: Alphabetical Sorting Invariant

*For any* list of apps loaded from configuration, the displayed list SHALL always be sorted alphabetically by display_name using locale-aware comparison.

**Validates: Requirements 1.1, 1.2, 1.3**

### Property 2: Selection Persistence Round-Trip

*For any* app selection made by the user in Reset Tab or Session Tab, saving to localStorage and then loading back SHALL return the same app name.

**Validates: Requirements 2.1, 2.2, 2.3, 2.4**

### Property 3: Active Session After Restore

*For any* successful restore operation (full, addon-only, or account-only), the restored session SHALL be marked as the active session for that app.

**Validates: Requirements 5.1, 5.2**

### Property 4: Auto-Generate ID After Restore

*For any* successful restore operation (full, addon-only, or account-only), the system SHALL generate new machine IDs, and the count of updated files SHALL be greater than or equal to zero.

**Validates: Requirements 6.1, 6.2, 6.3, 6.4**

### Property 5: Clear All Sessions Completeness

*For any* confirmed "Clear All Sessions" operation, after completion, the total session count across all apps SHALL be zero.

**Validates: Requirements 8.4**

### Property 6: Backward Compatibility for skipDataFolder

*For any* existing app config JSON that contains a skipDataFolder property, loading the config SHALL succeed and the property SHALL be ignored (not cause errors).

**Validates: Requirements 10.4**

## Error Handling

### Frontend Error Handling

1. **App Selection Fallback**: If persisted app no longer exists, fallback to first app in sorted list
2. **Launch App Failure**: Display error toast with message from backend
3. **Restore Failure**: Display error modal with option to disable "Skip Close App" if file lock detected
4. **Bulk Operation Failure**: Display error toast and continue with remaining items

### Backend Error Handling

1. **GenerateNewID**: Return count of 0 if no JSON files found, don't fail
2. **ClearAllSessions**: Log errors for individual session deletions but continue
3. **BackupAllSessions**: Return error if zip creation fails

## Testing Strategy

### Unit Tests

Unit tests will focus on specific examples and edge cases:

1. **Sorting Edge Cases**:
   - Empty app list
   - Single app
   - Apps with same display name
   - Apps with special characters in name

2. **Selection Persistence Edge Cases**:
   - Empty localStorage
   - Corrupted localStorage value
   - App no longer exists

3. **UI Component Tests**:
   - Context menu contains "Launch App" option
   - Post-restore modal appears with correct buttons
   - Settings toggles removed from UI

### Property-Based Tests

Property-based tests will verify universal properties using fast-check library:

1. **Alphabetical Sorting Property**:
   - Generate random app lists
   - Verify output is always sorted

2. **Selection Round-Trip Property**:
   - Generate random app names
   - Save and load, verify equality

3. **Active Session Property**:
   - Generate random restore operations
   - Verify session is marked active

4. **ID Generation Property**:
   - Generate random restore scenarios
   - Verify ID count is non-negative

### Test Configuration

- **Framework**: Vitest for frontend, Go testing for backend
- **Property Testing Library**: fast-check (JavaScript)
- **Minimum Iterations**: 100 per property test
- **Tag Format**: `Feature: ux-improvements, Property {number}: {property_text}`
