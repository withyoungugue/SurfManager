<script>
  import { settings } from './stores/settings.js';
  import { theme } from './stores/theme.js';
  import { 
    Palette, Shield, Database, FileText, 
    Monitor, RotateCcw, ChevronRight 
  } from 'lucide-svelte';
  import ThemeSelector from './ThemeSelector.svelte';
  import SettingToggle from './SettingToggle.svelte';

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

  function updateSetting(key, value) {
    settings.update(key, value);
  }
</script>

<div class="h-full flex animate-fadeIn gap-4">
  <!-- Sidebar -->
  <div class="w-48 bg-[var(--bg-elevated)] rounded-xl border border-[var(--border)] p-2 flex flex-col gap-1">
    {#each sections as section}
      <button
        class="w-full flex items-center gap-2 px-3 py-2.5 rounded-lg text-left transition-all
          {activeSection === section.id 
            ? 'bg-[var(--primary)] text-white' 
            : 'hover:bg-[var(--bg-hover)] text-[var(--text-secondary)]'}"
        on:click={() => activeSection = section.id}
      >
        <svelte:component this={section.icon} size={16} />
        <span class="text-sm font-medium">{section.label}</span>
        <ChevronRight size={14} class="ml-auto opacity-50" />
      </button>
    {/each}
  </div>

  <!-- Content -->
  <div class="flex-1 overflow-y-auto">
    {#if activeSection === 'appearance'}
      <h2 class="text-lg font-semibold mb-4 text-[var(--text-primary)]">Appearance</h2>
      
      <div class="space-y-3">
        <div class="flex items-center justify-between p-4 bg-[var(--bg-card)] rounded-lg border border-[var(--border)]">
          <div>
            <p class="font-medium text-[var(--text-primary)]">Theme</p>
            <p class="text-sm text-[var(--text-secondary)]">Choose your preferred color scheme</p>
          </div>
          <ThemeSelector />
        </div>
      </div>

    {:else if activeSection === 'behavior'}
      <h2 class="text-lg font-semibold mb-4 text-[var(--text-primary)]">Behavior</h2>
      
      <div class="space-y-3">
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

        <SettingToggle
          label="Keep App Running"
          description="Don't close target app before reset/backup/restore operations (enable this to skip closing)"
          checked={$settings.skipCloseApp}
          on:change={() => toggle('skipCloseApp')}
        />
      </div>

    {:else if activeSection === 'sessions'}
      <h2 class="text-lg font-semibold mb-4 text-[var(--text-primary)]">Sessions</h2>
      
      <div class="space-y-3">
        <SettingToggle
          label="Show Auto Backups"
          description="Include auto-generated backups in session list by default"
          checked={$settings.showAutoBackups}
          on:change={() => toggle('showAutoBackups')}
        />

        <div class="p-4 bg-[var(--bg-card)] rounded-lg border border-[var(--border)]">
          <p class="font-medium mb-2 text-[var(--text-primary)]">Default Filter</p>
          <p class="text-sm text-[var(--text-secondary)] mb-3">Default app filter when opening Sessions tab</p>
          <select 
            class="w-full bg-[var(--bg-hover)] border border-[var(--border)] rounded-lg px-3 py-2 text-[var(--text-primary)] focus:outline-none focus:border-[var(--primary)]"
            bind:value={$settings.defaultSessionFilter}
            on:change={(e) => updateSetting('defaultSessionFilter', e.target.value)}
          >
            <option value="all">All Apps</option>
            <option value="active">Active Only</option>
          </select>
        </div>
      </div>

    {:else if activeSection === 'notes'}
      <h2 class="text-lg font-semibold mb-4 text-[var(--text-primary)]">Notes</h2>
      
      <div class="space-y-3">
        <SettingToggle
          label="Auto Save"
          description="Automatically save notes while typing"
          checked={$settings.autoSaveNotes}
          on:change={() => toggle('autoSaveNotes')}
        />

        <div class="p-4 bg-[var(--bg-card)] rounded-lg border border-[var(--border)]">
          <p class="font-medium mb-2 text-[var(--text-primary)]">Auto Save Delay</p>
          <p class="text-sm text-[var(--text-secondary)] mb-3">
            Wait time before auto-saving (in seconds)
          </p>
          <input 
            type="range" 
            min="1" 
            max="10" 
            value={$settings.autoSaveDelay / 1000}
            on:input={(e) => updateSetting('autoSaveDelay', parseInt(e.target.value) * 1000)}
            class="w-full accent-[var(--primary)]"
          />
          <p class="text-sm text-center mt-2 text-[var(--text-secondary)]">{$settings.autoSaveDelay / 1000}s</p>
        </div>
      </div>

    {:else if activeSection === 'advanced'}
      <h2 class="text-lg font-semibold mb-4 text-[var(--text-primary)]">Advanced</h2>
      
      <div class="space-y-3">
        <div class="p-4 bg-[var(--bg-card)] rounded-lg border border-[var(--border)]">
          <p class="font-medium mb-2 text-[var(--text-primary)]">Log Retention</p>
          <p class="text-sm text-[var(--text-secondary)] mb-3">
            Maximum number of log entries to keep
          </p>
          <input 
            type="number" 
            min="10" 
            max="1000" 
            value={$settings.logRetention}
            on:change={(e) => updateSetting('logRetention', parseInt(e.target.value))}
            class="w-full bg-[var(--bg-hover)] border border-[var(--border)] rounded-lg px-3 py-2 text-[var(--text-primary)] focus:outline-none focus:border-[var(--primary)]"
          />
        </div>

        <SettingToggle
          label="Remember Last Tab"
          description="Open last active tab on startup"
          checked={$settings.rememberLastTab}
          on:change={() => toggle('rememberLastTab')}
        />

        <button
          class="w-full p-3 bg-[var(--danger)]/10 text-[var(--danger)] rounded-lg 
            hover:bg-[var(--danger)]/20 transition-colors flex items-center justify-center gap-2 border border-[var(--danger)]/30"
          on:click={() => { if(confirm('Reset all settings to default?')) settings.reset(); }}
        >
          <RotateCcw size={16} />
          Reset All Settings
        </button>
      </div>
    {/if}
  </div>
</div>
