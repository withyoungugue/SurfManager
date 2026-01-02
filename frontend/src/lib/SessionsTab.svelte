<script>
  import { onMount } from 'svelte';
  import { Plus, RefreshCw, Search, CheckSquare, FolderOpen, Trash2 } from 'lucide-svelte';
  import { GetActiveApps, GetAllSessions, CreateBackup, RestoreBackup, DeleteSession, SetActiveSession, OpenSessionFolder, CountAutoBackups } from '../../wailsjs/go/main/App.js';
  import { confirm } from './ConfirmModal.svelte';
  import { settings } from './stores/settings.js';

  export let logs = [];

  let apps = [];
  let sessions = [];
  let filter = 'all';
  let search = '';
  let showAuto = false;
  let autoBackupCount = 0;
  let loading = false;
  let selectedSessions = new Set();

  let showNewDialog = false;
  let newBackupApp = '';
  let newBackupName = '';

  // Context menu state
  let contextMenu = { show: false, x: 0, y: 0, session: null };

  $: skipCloseApp = $settings.skipCloseApp;

  onMount(() => {
    showAuto = $settings.showAutoBackups;
    filter = $settings.defaultSessionFilter;
    loadData();

    // Close context menu on click outside
    const handleClick = () => contextMenu.show = false;
    window.addEventListener('click', handleClick);
    return () => window.removeEventListener('click', handleClick);
  });

  async function loadData() {
    loading = true;
    try {
      apps = (await GetActiveApps()) || [];
      sessions = (await GetAllSessions(showAuto)) || [];
      autoBackupCount = await CountAutoBackups();
    } catch (e) {
      log(`Error: ${e}`);
    }
    loading = false;
  }

  function log(msg) {
    const time = new Date().toLocaleTimeString();
    logs = [...logs.slice(-99), `[${time}] ${msg}`];
  }

  $: filteredSessions = (sessions || []).filter(s => {
    if (filter !== 'all' && s.app.toLowerCase() !== filter.toLowerCase()) return false;
    if (search && !s.name.toLowerCase().includes(search.toLowerCase())) return false;
    if (showAuto && !s.is_auto) return false;
    if (!showAuto && s.is_auto) return false;
    return true;
  });

  function openNewDialog() {
    if (apps.length === 0) {
      alert('No apps configured. Add apps in Config tab first.');
      return;
    }
    newBackupApp = apps[0]?.app_name || '';
    newBackupName = '';
    showNewDialog = true;
  }

  async function handleCreateBackup() {
    if (!newBackupApp || !newBackupName.trim()) {
      alert('Please select app and enter session name');
      return;
    }

    log(`Creating backup: ${newBackupApp}/${newBackupName}...`);
    try {
      await CreateBackup(newBackupApp, newBackupName.trim(), skipCloseApp);
      log(`Backup created: ${newBackupName}`);
      showNewDialog = false;
      await loadData();
    } catch (e) {
      log(`Error: ${e}`);
      alert(`Error: ${e}`);
    }
  }

  async function handleRestore(session) {
    if ($settings.confirmBeforeRestore) {
      const confirmed = await confirm.restore(session.name);
      if (!confirmed) return;
    }
    
    log(`Restoring ${session.name}...`);
    try {
      await RestoreBackup(session.app, session.name, skipCloseApp);
      log(`Restored: ${session.name}`);
      await loadData();
    } catch (e) {
      log(`Error: ${e}`);
    }
  }

  async function handleDelete(session) {
    if ($settings.confirmBeforeDelete) {
      const confirmed = await confirm.delete(session.name);
      if (!confirmed) return;
    }
    
    try {
      await DeleteSession(session.app, session.name);
      log(`Deleted: ${session.name}`);
      selectedSessions.delete(`${session.app}/${session.name}`);
      selectedSessions = selectedSessions;
      await loadData();
    } catch (e) {
      log(`Error: ${e}`);
    }
  }

  async function handleSetActive(session) {
    try {
      await SetActiveSession(session.app, session.name);
      log(`Set active: ${session.name}`);
      await loadData();
    } catch (e) {
      log(`Error: ${e}`);
    }
  }

  async function handleOpenFolder(session) {
    try {
      await OpenSessionFolder(session.app, session.name);
    } catch (e) {
      log(`Error: ${e}`);
    }
  }

  function formatSize(bytes) {
    if (!bytes) return '0 B';
    if (bytes < 1024) return `${bytes} B`;
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
    return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
  }

  function formatDate(dateStr) {
    if (!dateStr) return '-';
    const date = new Date(dateStr);
    return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  }

  // Selection handling with CTRL+Click
  function handleRowClick(event, session) {
    const key = `${session.app}/${session.name}`;
    
    if (event.ctrlKey || event.metaKey) {
      // CTRL+Click: toggle selection
      if (selectedSessions.has(key)) {
        selectedSessions.delete(key);
      } else {
        selectedSessions.add(key);
      }
      selectedSessions = selectedSessions;
    }
  }

  function selectAll() {
    if (selectedSessions.size === filteredSessions.length) {
      // Deselect all
      selectedSessions = new Set();
    } else {
      // Select all
      selectedSessions = new Set(filteredSessions.map(s => `${s.app}/${s.name}`));
    }
  }

  async function deleteSelected() {
    if (selectedSessions.size === 0) return;
    
    if ($settings.confirmBeforeDelete) {
      const confirmed = await confirm.bulkDelete(selectedSessions.size);
      if (!confirmed) return;
    }

    for (const key of selectedSessions) {
      const [app, name] = key.split('/');
      try {
        await DeleteSession(app, name);
        log(`Deleted: ${name}`);
      } catch (e) {
        log(`Error deleting ${name}: ${e}`);
      }
    }
    selectedSessions = new Set();
    await loadData();
  }

  // Context menu
  function handleContextMenu(event, session) {
    event.preventDefault();
    contextMenu = {
      show: true,
      x: event.clientX,
      y: event.clientY,
      session
    };
  }

  function closeContextMenu() {
    contextMenu.show = false;
  }
</script>

<div class="h-full flex flex-col gap-4 animate-fadeIn">
  <!-- Header -->
  <div class="flex items-center justify-between flex-wrap gap-4">
    <div class="flex items-center gap-3">
      <h2 class="text-xl font-semibold text-[var(--text-primary)]">Sessions</h2>
      <button 
        class="px-4 py-2 rounded-lg font-medium bg-[var(--primary)] hover:bg-[var(--primary-light)] hover:text-black text-white transition-all flex items-center gap-2 text-sm"
        on:click={openNewDialog}
      >
        <Plus size={16} />
        New Backup
      </button>
    </div>
    
    <div class="flex items-center gap-3 flex-wrap">
      <div class="relative">
        <Search size={16} class="absolute left-3 top-1/2 -translate-y-1/2 text-[var(--text-muted)]" />
        <input
          type="text"
          placeholder="Search..."
          class="bg-[var(--bg-hover)] border border-[var(--border)] rounded-lg pl-9 pr-3 py-1.5 text-sm text-[var(--text-primary)] placeholder-[var(--text-muted)] focus:border-[var(--primary)] focus:outline-none w-36"
          bind:value={search}
        />
      </div>

      <button
        class="px-3 py-1.5 rounded-lg text-sm font-medium transition-all border
               {showAuto 
                 ? 'bg-[var(--warning)]/20 text-[var(--warning)] border-[var(--warning)]' 
                 : 'bg-[var(--bg-hover)] text-[var(--text-muted)] border-[var(--border)]'}"
        on:click={() => { showAuto = !showAuto; loadData(); }}
      >
        Auto-Backup ({autoBackupCount})
      </button>

      <button 
        class="p-2 rounded-lg text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-hover)] transition-colors"
        on:click={loadData} 
        title="Refresh"
      >
        <RefreshCw size={18} class={loading ? 'animate-spin' : ''} />
      </button>
    </div>
  </div>

  <!-- Selection Actions Bar -->
  <div class="flex items-center gap-3">
    <select 
      class="bg-[var(--bg-hover)] border border-[var(--border)] rounded-lg px-3 py-1.5 text-sm text-[var(--text-primary)] focus:border-[var(--primary)] focus:outline-none"
      bind:value={filter}
    >
      <option value="all">All Apps</option>
      {#each apps as app}
        <option value={app.app_name}>{app.display_name}</option>
      {/each}
    </select>

    <button 
      class="px-3 py-1.5 rounded-lg text-sm font-medium bg-[var(--bg-hover)] hover:bg-[var(--border)] border border-[var(--border)] text-[var(--text-secondary)] transition-all flex items-center gap-2"
      on:click={selectAll}
      title="Select All (or CTRL+Click rows)"
    >
      <CheckSquare size={14} />
      {selectedSessions.size === filteredSessions.length && filteredSessions.length > 0 ? 'Deselect All' : 'Select All'}
    </button>
    
    {#if selectedSessions.size > 0}
      <button 
        class="px-3 py-1.5 rounded-lg text-sm font-medium bg-[var(--danger)]/20 hover:bg-[var(--danger)]/30 text-[var(--danger)] border border-[var(--danger)]/30 transition-all"
        on:click={deleteSelected}
      >
        Delete Selected ({selectedSessions.size})
      </button>
    {/if}

    <span class="text-sm text-[var(--text-muted)] ml-auto">
      {filteredSessions.length} session(s)
    </span>
  </div>

  <!-- Sessions Table -->
  <div class="flex-1 bg-[var(--bg-elevated)] rounded-xl border border-[var(--border)] overflow-hidden">
    <div class="overflow-auto h-full">
      <table class="w-full">
        <thead class="bg-[var(--bg-card)] sticky top-0">
          <tr class="text-left text-xs text-[var(--text-secondary)]">
            <th class="p-3 w-12">#</th>
            <th class="p-3">App</th>
            <th class="p-3">Session Name</th>
            <th class="p-3">Size</th>
            <th class="p-3">Created</th>
            <th class="p-3">Status</th>
            <th class="p-3 text-right">Actions</th>
          </tr>
        </thead>
        <tbody>
          {#if filteredSessions.length === 0}
            <tr>
              <td colspan="7" class="p-8 text-center text-[var(--text-muted)]">
                {showAuto ? 'No auto-backups found' : 'No sessions found'}
              </td>
            </tr>
          {/if}
          {#each filteredSessions as session, index}
            {@const key = `${session.app}/${session.name}`}
            {@const isSelected = selectedSessions.has(key)}
            <tr 
              class="border-b border-[var(--border)] transition-colors cursor-pointer
                     {isSelected ? 'bg-[var(--primary-dim)]' : 'hover:bg-[var(--bg-card)]'}"
              on:click={(e) => handleRowClick(e, session)}
              on:contextmenu={(e) => handleContextMenu(e, session)}
            >
              <td class="p-3 text-[var(--text-muted)] text-sm">{index + 1}</td>
              <td class="p-3">
                <span class="font-medium text-[var(--text-primary)] capitalize">{session.app}</span>
              </td>
              <td class="p-3 text-[var(--text-secondary)]">{session.name}</td>
              <td class="p-3 text-[var(--text-muted)] text-sm">{formatSize(session.size)}</td>
              <td class="p-3 text-[var(--text-muted)] text-sm">{formatDate(session.created)}</td>
              <td class="p-3">
                {#if session.is_auto}
                  <span class="px-2 py-0.5 rounded-full text-xs font-medium bg-[var(--warning)]/20 text-[var(--warning)]">Auto</span>
                {:else if session.is_active}
                  <span class="px-2 py-0.5 rounded-full text-xs font-medium bg-[var(--success)]/20 text-[var(--success)]">Active</span>
                {:else}
                  <span class="px-2 py-0.5 rounded-full text-xs font-medium bg-[var(--bg-hover)] text-[var(--text-muted)]">Ready</span>
                {/if}
              </td>
              <td class="p-3">
                <div class="flex items-center justify-end gap-2">
                  <button
                    class="px-2 py-1 rounded text-xs font-medium text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-hover)] transition-colors flex items-center gap-1"
                    on:click|stopPropagation={() => handleOpenFolder(session)}
                  >
                    <FolderOpen size={12} />
                    Folder
                  </button>
                  <button
                    class="px-2 py-1 rounded text-xs font-medium text-[var(--text-secondary)] hover:text-[var(--danger)] hover:bg-[var(--bg-hover)] transition-colors flex items-center gap-1"
                    on:click|stopPropagation={() => handleDelete(session)}
                  >
                    <Trash2 size={12} />
                    Delete
                  </button>
                </div>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  </div>

  <!-- Help text -->
  <p class="text-xs text-[var(--text-muted)]">
    💡 Tip: CTRL+Click to select multiple rows, or right-click for quick actions
  </p>
</div>

<!-- Context Menu -->
{#if contextMenu.show}
  <div 
    class="fixed bg-[var(--bg-card)] border border-[var(--border)] rounded-lg shadow-xl py-1 z-50 min-w-[160px]"
    style="left: {contextMenu.x}px; top: {contextMenu.y}px"
    on:click|stopPropagation={closeContextMenu}
  >
    <button
      class="w-full px-4 py-2 text-left text-sm text-[var(--text-secondary)] hover:bg-[var(--bg-hover)] hover:text-[var(--primary)] transition-colors"
      on:click={() => { handleRestore(contextMenu.session); closeContextMenu(); }}
    >
      Restore Session
    </button>
    {#if !contextMenu.session?.is_auto}
      <button
        class="w-full px-4 py-2 text-left text-sm text-[var(--text-secondary)] hover:bg-[var(--bg-hover)] hover:text-[var(--success)] transition-colors"
        on:click={() => { handleSetActive(contextMenu.session); closeContextMenu(); }}
      >
        Set as Active
      </button>
    {/if}
    <button
      class="w-full px-4 py-2 text-left text-sm text-[var(--text-secondary)] hover:bg-[var(--bg-hover)] hover:text-[var(--text-primary)] transition-colors"
      on:click={() => { handleOpenFolder(contextMenu.session); closeContextMenu(); }}
    >
      Open Folder
    </button>
    <div class="border-t border-[var(--border)] my-1"></div>
    <button
      class="w-full px-4 py-2 text-left text-sm text-[var(--danger)] hover:bg-[var(--danger)]/10 transition-colors"
      on:click={() => { handleDelete(contextMenu.session); closeContextMenu(); }}
    >
      Delete Session
    </button>
  </div>
{/if}

<!-- New Backup Dialog -->
{#if showNewDialog}
  <div class="fixed inset-0 bg-black/60 flex items-center justify-center z-50" on:click|self={() => showNewDialog = false}>
    <div class="bg-[var(--bg-card)] rounded-xl border border-[var(--border)] w-full max-w-md p-6 animate-fadeIn">
      <h3 class="text-lg font-semibold text-[var(--text-primary)] mb-4">Create New Backup</h3>
      
      <div class="space-y-4">
        <div>
          <label class="block text-sm text-[var(--text-secondary)] mb-1">Application</label>
          <select 
            class="w-full bg-[var(--bg-hover)] border border-[var(--border)] rounded-lg px-3 py-2 text-[var(--text-primary)] focus:border-[var(--primary)] focus:outline-none"
            bind:value={newBackupApp}
          >
            {#each apps as app}
              <option value={app.app_name}>{app.display_name}</option>
            {/each}
          </select>
        </div>

        <div>
          <label class="block text-sm text-[var(--text-secondary)] mb-1">Session Name</label>
          <input
            type="text"
            class="w-full bg-[var(--bg-hover)] border border-[var(--border)] rounded-lg px-3 py-2 text-[var(--text-primary)] placeholder-[var(--text-muted)] focus:border-[var(--primary)] focus:outline-none"
            placeholder="e.g., work-main, personal"
            bind:value={newBackupName}
          />
        </div>
      </div>

      <div class="flex justify-end gap-3 mt-6">
        <button 
          class="px-4 py-2 rounded-lg font-medium bg-[var(--bg-hover)] hover:bg-[var(--border)] border border-[var(--border)] text-[var(--text-secondary)] transition-all"
          on:click={() => showNewDialog = false}
        >
          Cancel
        </button>
        <button 
          class="px-4 py-2 rounded-lg font-medium bg-[var(--primary)] hover:bg-[var(--primary-light)] hover:text-black text-white transition-all"
          on:click={handleCreateBackup}
        >
          Create
        </button>
      </div>
    </div>
  </div>
{/if}
