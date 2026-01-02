# SurfManager v2.0 - UI Improvements Proposal

## Overview

Proposal perbaikan UI berdasarkan analisis mendalam terhadap implementasi saat ini. Fokus utama:
1. **Solarized Theme** - Tema baru selain dark mode
2. **Confirmation Actions** - Konsistensi modal konfirmasi
3. **Custom Settings** - Sistem pengaturan enable/disable fitur

---

## 1. Theme System (Solarized + Dark)

### Color Palettes

```css
/* Dark Theme (Current - Default) */
:root[data-theme="dark"] {
  --bg-base:      #0a0a0a;
  --bg-card:      #141414;
  --bg-hover:     #1a1a1a;
  --bg-elevated:  #1f1f1f;
  
  --text-primary:   #ffffff;
  --text-secondary: #888888;
  --text-muted:     #555555;
  
  --border:       #252525;
  --border-hover: #333333;
  
  --primary:       #0d7377;
  --primary-light: #14ffec;
  --primary-dim:   #0d7377/20;
  
  --success: #22c55e;
  --warning: #f59e0b;
  --danger:  #ef4444;
  --info:    #3b82f6;
}

/* Solarized Dark Theme */
:root[data-theme="solarized"] {
  --bg-base:      #002b36;  /* base03 */
  --bg-card:      #073642;  /* base02 */
  --bg-hover:     #094552;
  --bg-elevated:  #0a4f5c;
  
  --text-primary:   #fdf6e3;  /* base3 */
  --text-secondary: #93a1a1;  /* base1 */
  --text-muted:     #657b83;  /* base00 */
  
  --border:       #094552;
  --border-hover: #0a5a68;
  
  --primary:       #2aa198;  /* cyan */
  --primary-light: #35c9be;
  --primary-dim:   #2aa198/20;
  
  --success: #859900;  /* green */
  --warning: #b58900;  /* yellow */
  --danger:  #dc322f;  /* red */
  --info:    #268bd2;  /* blue */
}

/* Solarized Light Theme (Optional Future) */
:root[data-theme="solarized-light"] {
  --bg-base:      #fdf6e3;  /* base3 */
  --bg-card:      #eee8d5;  /* base2 */
  --bg-hover:     #e6e0cc;
  --bg-elevated:  #ddd7c3;
  
  --text-primary:   #073642;  /* base02 */
  --text-secondary: #586e75;  /* base01 */
  --text-muted:     #93a1a1;  /* base1 */
  
  --border:       #d3cdb9;
  --border-hover: #c9c3af;
  
  --primary:       #2aa198;
  --primary-light: #1a8a82;
  --primary-dim:   #2aa198/20;
  
  --success: #859900;
  --warning: #b58900;
  --danger:  #dc322f;
  --info:    #268bd2;
}
```

### Theme Store Implementation

```javascript
// frontend/src/lib/stores/theme.js
import { writable } from 'svelte/store';

const STORAGE_KEY = 'surfmanager-theme';
const DEFAULT_THEME = 'dark';

function createThemeStore() {
  const stored = localStorage.getItem(STORAGE_KEY) || DEFAULT_THEME;
  const { subscribe, set, update } = writable(stored);

  return {
    subscribe,
    set: (theme) => {
      localStorage.setItem(STORAGE_KEY, theme);
      document.documentElement.setAttribute('data-theme', theme);
      set(theme);
    },
    toggle: () => {
      update(current => {
        const themes = ['dark', 'solarized', 'solarized-light'];
        const next = themes[(themes.indexOf(current) + 1) % themes.length];
        localStorage.setItem(STORAGE_KEY, next);
        document.documentElement.setAttribute('data-theme', next);
        return next;
      });
    },
    init: () => {
      document.documentElement.setAttribute('data-theme', stored);
    }
  };
}

export const theme = createThemeStore();
```

### Theme Selector Component

```svelte
<!-- frontend/src/lib/ThemeSelector.svelte -->
<script>
  import { theme } from './stores/theme.js';
  import { Sun, Moon, Palette } from 'lucide-svelte';

  const themes = [
    { id: 'dark', label: 'Dark', icon: Moon },
    { id: 'solarized', label: 'Solarized Dark', icon: Palette },
    { id: 'solarized-light', label: 'Solarized Light', icon: Sun },
  ];
</script>

<div class="flex gap-1 p-1 bg-[var(--bg-card)] rounded-lg border border-[var(--border)]">
  {#each themes as t}
    <button
      class="p-2 rounded-md transition-all {$theme === t.id 
        ? 'bg-[var(--primary)] text-white' 
        : 'hover:bg-[var(--bg-hover)] text-[var(--text-secondary)]'}"
      on:click={() => theme.set(t.id)}
      title={t.label}
    >
      <svelte:component this={t.icon} size={16} />
    </button>
  {/each}
</div>
```

---

## 2. Confirmation Actions System

### Problem
Saat ini menggunakan `confirm()` browser yang jelek dan tidak konsisten dengan design.

### Solution
Gunakan `ConfirmModal` yang sudah ada untuk semua konfirmasi.

### Locations to Fix

| File | Line | Current | Fix |
|------|------|---------|-----|
| ResetTab.svelte | ~60 | `confirm('Reset?')` | `ConfirmModal.confirm.danger()` |
| SessionsTab.svelte | ~95 | `confirm('Delete?')` | `ConfirmModal.confirm.danger()` |
| SessionsTab.svelte | ~105 | `confirm('Restore?')` | `ConfirmModal.confirm()` |
| ConfigTab.svelte | ~75 | `confirm('Delete?')` | `ConfirmModal.confirm.danger()` |
| NotesTab.svelte | ~42 | `confirm('Discard?')` | `ConfirmModal.confirm()` |

### Enhanced ConfirmModal API

```svelte
<!-- Tambahan helper methods -->
<script context="module">
  // Existing...
  
  // Quick helpers
  confirm.delete = (itemName) => confirm({
    title: 'Delete Confirmation',
    message: `Are you sure you want to delete "${itemName}"? This cannot be undone.`,
    confirmText: 'Delete',
    danger: true
  });

  confirm.reset = (appName) => confirm({
    title: 'Reset Confirmation', 
    message: `Reset all data for ${appName}? A backup will be created automatically.`,
    confirmText: 'Reset',
    danger: true
  });

  confirm.restore = (sessionName) => confirm({
    title: 'Restore Session',
    message: `Restore "${sessionName}"? Current data will be replaced.`,
    confirmText: 'Restore'
  });

  confirm.discard = () => confirm({
    title: 'Unsaved Changes',
    message: 'You have unsaved changes. Discard them?',
    confirmText: 'Discard',
    danger: true
  });
</script>
```

### Usage Example

```svelte
<!-- Before -->
async function handleDelete(session) {
  if (!confirm(`Delete "${session.name}"?`)) return;
  // ...
}

<!-- After -->
import { confirm } from './ConfirmModal.svelte';

async function handleDelete(session) {
  if (!await confirm.delete(session.name)) return;
  // ...
}
```

---

## 3. Custom Settings System

### Settings Structure

```javascript
// frontend/src/lib/stores/settings.js
import { writable } from 'svelte/store';

const STORAGE_KEY = 'surfmanager-settings';

const defaultSettings = {
  // Appearance
  theme: 'dark',
  sidebarExpanded: false,
  
  // Behavior
  autoBackup: true,
  confirmBeforeReset: true,
  confirmBeforeDelete: true,
  confirmBeforeRestore: true,
  
  // Sessions
  showAutoBackups: false,
  defaultSessionFilter: 'all',
  
  // Notes
  autoSaveNotes: true,
  autoSaveDelay: 3000, // ms
  
  // Advanced
  logRetention: 100,
  showDebugInfo: false,
  
  // Window
  rememberWindowSize: true,
  rememberLastTab: true,
  lastActiveTab: 'reset',
};

function createSettingsStore() {
  const stored = JSON.parse(localStorage.getItem(STORAGE_KEY) || '{}');
  const initial = { ...defaultSettings, ...stored };
  const { subscribe, set, update } = writable(initial);

  return {
    subscribe,
    set: (settings) => {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(settings));
      set(settings);
    },
    update: (key, value) => {
      update(s => {
        const newSettings = { ...s, [key]: value };
        localStorage.setItem(STORAGE_KEY, JSON.stringify(newSettings));
        return newSettings;
      });
    },
    reset: () => {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(defaultSettings));
      set(defaultSettings);
    }
  };
}

export const settings = createSettingsStore();
```

### Settings Tab Component

```svelte
<!-- frontend/src/lib/SettingsTab.svelte -->
<script>
  import { settings } from './stores/settings.js';
  import { theme } from './stores/theme.js';
  import { 
    Palette, Shield, Database, FileText, 
    Monitor, RotateCcw, ChevronRight 
  } from 'lucide-svelte';
  import ThemeSelector from './ThemeSelector.svelte';

  let activeSection = 'appearance';

  const sections = [
    { id: 'appearance', label: 'Appearance', icon: Palette },
    { id: 'behavior', label: 'Behavior', icon: Shield },
    { id: 'sessions', label: 'Sessions', icon: Database },
    { id: 'notes', label: 'Notes', icon: FileText },
    { id: 'advanced', label: 'Advanced', icon: Monitor },
  ];

  function toggle(key) {
    settings.update(key, !$settings[key]);
  }
</script>

<div class="h-full flex animate-fadeIn">
  <!-- Sidebar -->
  <div class="w-48 border-r border-[var(--border)] p-2">
    {#each sections as section}
      <button
        class="w-full flex items-center gap-2 px-3 py-2 rounded-lg text-left transition-all
          {activeSection === section.id 
            ? 'bg-[var(--primary)] text-white' 
            : 'hover:bg-[var(--bg-hover)] text-[var(--text-secondary)]'}"
        on:click={() => activeSection = section.id}
      >
        <svelte:component this={section.icon} size={16} />
        <span class="text-sm">{section.label}</span>
        <ChevronRight size={14} class="ml-auto opacity-50" />
      </button>
    {/each}
  </div>

  <!-- Content -->
  <div class="flex-1 p-6 overflow-y-auto">
    {#if activeSection === 'appearance'}
      <h2 class="text-lg font-semibold mb-4">Appearance</h2>
      
      <div class="space-y-4">
        <div class="flex items-center justify-between p-4 bg-[var(--bg-card)] rounded-lg">
          <div>
            <p class="font-medium">Theme</p>
            <p class="text-sm text-[var(--text-secondary)]">Choose your preferred color scheme</p>
          </div>
          <ThemeSelector />
        </div>

        <SettingToggle
          label="Expanded Sidebar"
          description="Show labels in sidebar navigation"
          checked={$settings.sidebarExpanded}
          on:change={() => toggle('sidebarExpanded')}
        />
      </div>

    {:else if activeSection === 'behavior'}
      <h2 class="text-lg font-semibold mb-4">Behavior</h2>
      
      <div class="space-y-4">
        <SettingToggle
          label="Auto Backup"
          description="Automatically backup before reset"
          checked={$settings.autoBackup}
          on:change={() => toggle('autoBackup')}
        />

        <SettingToggle
          label="Confirm Before Reset"
          description="Show confirmation dialog before resetting app data"
          checked={$settings.confirmBeforeReset}
          on:change={() => toggle('confirmBeforeReset')}
        />

        <SettingToggle
          label="Confirm Before Delete"
          description="Show confirmation dialog before deleting sessions"
          checked={$settings.confirmBeforeDelete}
          on:change={() => toggle('confirmBeforeDelete')}
        />

        <SettingToggle
          label="Confirm Before Restore"
          description="Show confirmation dialog before restoring sessions"
          checked={$settings.confirmBeforeRestore}
          on:change={() => toggle('confirmBeforeRestore')}
        />
      </div>

    {:else if activeSection === 'sessions'}
      <h2 class="text-lg font-semibold mb-4">Sessions</h2>
      
      <div class="space-y-4">
        <SettingToggle
          label="Show Auto Backups"
          description="Include auto-generated backups in session list"
          checked={$settings.showAutoBackups}
          on:change={() => toggle('showAutoBackups')}
        />

        <div class="p-4 bg-[var(--bg-card)] rounded-lg">
          <p class="font-medium mb-2">Default Filter</p>
          <select 
            class="w-full bg-[var(--bg-hover)] border border-[var(--border)] rounded-lg px-3 py-2"
            bind:value={$settings.defaultSessionFilter}
            on:change={(e) => settings.update('defaultSessionFilter', e.target.value)}
          >
            <option value="all">All Apps</option>
            <option value="active">Active Only</option>
          </select>
        </div>
      </div>

    {:else if activeSection === 'notes'}
      <h2 class="text-lg font-semibold mb-4">Notes</h2>
      
      <div class="space-y-4">
        <SettingToggle
          label="Auto Save"
          description="Automatically save notes while typing"
          checked={$settings.autoSaveNotes}
          on:change={() => toggle('autoSaveNotes')}
        />

        <div class="p-4 bg-[var(--bg-card)] rounded-lg">
          <p class="font-medium mb-2">Auto Save Delay</p>
          <p class="text-sm text-[var(--text-secondary)] mb-2">
            Wait time before auto-saving (in seconds)
          </p>
          <input 
            type="range" 
            min="1" 
            max="10" 
            value={$settings.autoSaveDelay / 1000}
            on:change={(e) => settings.update('autoSaveDelay', e.target.value * 1000)}
            class="w-full"
          />
          <p class="text-sm text-center mt-1">{$settings.autoSaveDelay / 1000}s</p>
        </div>
      </div>

    {:else if activeSection === 'advanced'}
      <h2 class="text-lg font-semibold mb-4">Advanced</h2>
      
      <div class="space-y-4">
        <div class="p-4 bg-[var(--bg-card)] rounded-lg">
          <p class="font-medium mb-2">Log Retention</p>
          <p class="text-sm text-[var(--text-secondary)] mb-2">
            Maximum number of log entries to keep
          </p>
          <input 
            type="number" 
            min="10" 
            max="1000" 
            value={$settings.logRetention}
            on:change={(e) => settings.update('logRetention', parseInt(e.target.value))}
            class="w-full bg-[var(--bg-hover)] border border-[var(--border)] rounded-lg px-3 py-2"
          />
        </div>

        <SettingToggle
          label="Show Debug Info"
          description="Display additional debugging information"
          checked={$settings.showDebugInfo}
          on:change={() => toggle('showDebugInfo')}
        />

        <SettingToggle
          label="Remember Window Size"
          description="Restore window size on startup"
          checked={$settings.rememberWindowSize}
          on:change={() => toggle('rememberWindowSize')}
        />

        <SettingToggle
          label="Remember Last Tab"
          description="Open last active tab on startup"
          checked={$settings.rememberLastTab}
          on:change={() => toggle('rememberLastTab')}
        />

        <button
          class="w-full p-3 bg-[var(--danger)]/10 text-[var(--danger)] rounded-lg 
            hover:bg-[var(--danger)]/20 transition-colors flex items-center justify-center gap-2"
          on:click={() => settings.reset()}
        >
          <RotateCcw size={16} />
          Reset All Settings
        </button>
      </div>
    {/if}
  </div>
</div>
```

### SettingToggle Component

```svelte
<!-- frontend/src/lib/SettingToggle.svelte -->
<script>
  export let label = '';
  export let description = '';
  export let checked = false;
</script>

<label class="flex items-center justify-between p-4 bg-[var(--bg-card)] rounded-lg cursor-pointer hover:bg-[var(--bg-hover)] transition-colors">
  <div>
    <p class="font-medium">{label}</p>
    {#if description}
      <p class="text-sm text-[var(--text-secondary)]">{description}</p>
    {/if}
  </div>
  
  <div class="relative">
    <input 
      type="checkbox" 
      {checked}
      on:change
      class="sr-only peer"
    />
    <div class="w-11 h-6 bg-[var(--bg-hover)] rounded-full peer 
      peer-checked:bg-[var(--primary)] transition-colors"></div>
    <div class="absolute left-1 top-1 w-4 h-4 bg-white rounded-full transition-transform
      peer-checked:translate-x-5"></div>
  </div>
</label>
```

---

## 4. Updated Navigation (Add Settings Tab)

```svelte
<!-- App.svelte - Updated tabs -->
const tabs = [
  { id: 'reset', label: 'Reset', icon: RotateCcw },
  { id: 'sessions', label: 'Sessions', icon: Database },
  { id: 'config', label: 'Config', icon: Settings },
  { id: 'notes', label: 'Notes', icon: FileText },
  { id: 'settings', label: 'Settings', icon: Sliders },  // NEW
  { id: 'about', label: 'About', icon: Info },
];
```

---

## 5. Implementation Phases

### Phase 1: Theme System (2-3 hours)
- [ ] Create CSS custom properties in style.css
- [ ] Create theme store
- [ ] Create ThemeSelector component
- [ ] Update tailwind.config.js
- [ ] Update all components to use CSS variables

### Phase 2: Confirmation System (1-2 hours)
- [ ] Enhance ConfirmModal with helper methods
- [ ] Replace all `confirm()` calls in ResetTab
- [ ] Replace all `confirm()` calls in SessionsTab
- [ ] Replace all `confirm()` calls in ConfigTab
- [ ] Replace all `confirm()` calls in NotesTab

### Phase 3: Settings System (2-3 hours)
- [ ] Create settings store
- [ ] Create SettingToggle component
- [ ] Create SettingsTab component
- [ ] Add Settings tab to navigation
- [ ] Integrate settings with existing features

### Phase 4: Polish (1-2 hours)
- [ ] Add fadeIn animations to all tabs
- [ ] Add loading states to buttons
- [ ] Test all themes
- [ ] Test all settings
- [ ] Fix any responsive issues

---

## 6. File Changes Summary

### New Files
```
frontend/src/lib/stores/
├── theme.js
└── settings.js

frontend/src/lib/
├── ThemeSelector.svelte
├── SettingsTab.svelte
└── SettingToggle.svelte
```

### Modified Files
```
frontend/src/style.css          # Add CSS variables
frontend/tailwind.config.js     # Update color config
frontend/src/App.svelte         # Add Settings tab
frontend/src/lib/ConfirmModal.svelte  # Add helper methods
frontend/src/lib/ResetTab.svelte      # Use ConfirmModal
frontend/src/lib/SessionsTab.svelte   # Use ConfirmModal
frontend/src/lib/ConfigTab.svelte     # Use ConfirmModal
frontend/src/lib/NotesTab.svelte      # Use ConfirmModal
```

---

## 7. Preview

### Dark Theme (Current)
```
┌─────────────────────────────────────┐
│  ████  SurfManager     [🌙][👤][📦] │
├──────┬──────────────────────────────┤
│ 🔄   │  ┌─────────┐ ┌─────────┐    │
│ 💾   │  │ VSCode  │ │ Cursor  │    │
│ ⚙️   │  │ #141414 │ │ #141414 │    │
│ 📝   │  └─────────┘ └─────────┘    │
│ ⚡   │                              │
│ ℹ️   │  Background: #0a0a0a         │
└──────┴──────────────────────────────┘
```

### Solarized Dark Theme
```
┌─────────────────────────────────────┐
│  ████  SurfManager     [🎨][👤][📦] │
├──────┬──────────────────────────────┤
│ 🔄   │  ┌─────────┐ ┌─────────┐    │
│ 💾   │  │ VSCode  │ │ Cursor  │    │
│ ⚙️   │  │ #073642 │ │ #073642 │    │
│ 📝   │  └─────────┘ └─────────┘    │
│ ⚡   │                              │
│ ℹ️   │  Background: #002b36         │
└──────┴──────────────────────────────┘
```

---

*Modern themes, consistent UX, customizable settings* ✨
