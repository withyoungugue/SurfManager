<script>
  import { onMount } from 'svelte';
  import { Plus, Trash2, FolderOpen, RefreshCw, X, Check, FileJson, ToggleLeft, Edit } from 'lucide-svelte';
  import { GetApps, CheckAppInstalled, SaveApp, DeleteApp, ToggleApp, OpenConfigFolder, SelectFile, SelectFolder, SelectFolderFromHome, SelectExeFromLocalPrograms, SelectFolderFromRoaming, OpenFolder } from '../../wailsjs/go/main/App.js';
  import { confirm } from './ConfirmModal.svelte';
  import { settings } from './stores/settings.js';

  export let logs = [];
  export let showAddDialog = false;

  let apps = [];
  let loading = false;
  
  // Context menu state
  let contextMenu = { show: false, x: 0, y: 0, app: null };
  
  // New app form
  let newApp = {
    app_name: '',
    display_name: '',
    exe_path: '',
    data_path: '',
    app_type: 'vscode', // 'vscode' or 'custom'
    backup_items: [],
    addon_paths: []
  };

  // VSCode preset backup items
  const vscodePresetItems = [
    { path: 'User', description: 'Settings & keybindings', optional: false, enabled: true },
    { path: 'Local Storage', description: 'Extension data', optional: true, enabled: true },
    { path: 'Session Storage', description: 'Session data', optional: true, enabled: true },
    { path: 'Preferences', description: 'App preferences', optional: true, enabled: true },
    { path: 'Local State', description: 'Local state', optional: true, enabled: true },
    { path: 'Cache', description: 'Cache files (large)', optional: true, enabled: false },
    { path: 'CachedData', description: 'Cached data', optional: true, enabled: false },
    { path: 'CachedExtensions', description: 'Cached extensions', optional: true, enabled: false },
  ];

  // Custom backup item input
  let customItemPath = '';
  let customItemDesc = '';

  onMount(loadApps);

  onMount(() => {
    // Close context menu on click outside
    const handleClick = () => contextMenu.show = false;
    window.addEventListener('click', handleClick);
    return () => window.removeEventListener('click', handleClick);
  });

  async function loadApps() {
    loading = true;
    try {
      const allApps = (await GetApps()) || [];
      apps = await Promise.all(allApps.map(async (app) => {
        const installed = await CheckAppInstalled(app.app_name);
        return { ...app, installed };
      }));
    } catch (e) {
      log(`Error: ${e}`);
    }
    loading = false;
  }

  function log(msg) {
    const time = new Date().toLocaleTimeString();
    logs = [...logs.slice(-99), `[${time}] ${msg}`];
  }

  async function handleToggle(app) {
    try {
      await ToggleApp(app.app_name);
      log(`Toggled ${app.display_name}: ${app.active ? 'Inactive' : 'Active'}`);
      await loadApps();
    } catch (e) {
      log(`Error: ${e}`);
    }
  }

  async function handleDelete(app) {
    if ($settings.confirmBeforeDelete) {
      const confirmed = await confirm.delete(app.display_name);
      if (!confirmed) return;
    }
    
    try {
      await DeleteApp(app.app_name);
      log(`Deleted: ${app.display_name}`);
      await loadApps();
    } catch (e) {
      log(`Error: ${e}`);
    }
  }

  async function handleOpenFolder() {
    try {
      await OpenConfigFolder();
    } catch (e) {
      log(`Error: ${e}`);
    }
  }

  // Context menu handlers
  function handleContextMenu(event, app) {
    event.preventDefault();
    contextMenu = {
      show: true,
      x: event.clientX,
      y: event.clientY,
      app
    };
  }

  function closeContextMenu() {
    contextMenu.show = false;
  }

  async function handleOpenAppConfig(app) {
    try {
      // Open the config folder and let user find the JSON
      await OpenConfigFolder();
      log(`Opened config folder for ${app.display_name}`);
    } catch (e) {
      log(`Error: ${e}`);
    }
  }

  async function handleSetActive(app) {
    if (!app.active) {
      await handleToggle(app);
    }
  }

  async function handleSetInactive(app) {
    if (app.active) {
      await handleToggle(app);
    }
  }

  async function selectExe() {
    try {
      const path = await SelectExeFromLocalPrograms('Select Executable');
      if (path) {
        newApp.exe_path = path;
        const fileName = path.split(/[/\\]/).pop();
        const appName = fileName.replace('.exe', '');
        newApp.app_name = appName.toLowerCase();
        newApp.display_name = appName.charAt(0).toUpperCase() + appName.slice(1);
      }
    } catch (e) {
      log(`Error: ${e}`);
    }
  }

  async function selectDataFolder() {
    try {
      const path = await SelectFolderFromRoaming('Select Data Folder');
      if (path) {
        newApp.data_path = path;
      }
    } catch (e) {
      log(`Error: ${e}`);
    }
  }

  async function addAddonFolder() {
    try {
      const path = await SelectFolderFromHome('Select Additional Folder');
      if (path && !newApp.addon_paths.includes(path)) {
        newApp.addon_paths = [...newApp.addon_paths, path];
      }
    } catch (e) {
      log(`Error: ${e}`);
    }
  }

  function removeAddonFolder(path) {
    newApp.addon_paths = newApp.addon_paths.filter(p => p !== path);
  }

  function setAppType(type) {
    newApp.app_type = type;
    if (type === 'vscode') {
      // Reset to preset items
      newApp.backup_items = vscodePresetItems.map(item => ({ ...item }));
    } else {
      // Clear for custom
      newApp.backup_items = [];
    }
  }

  function toggleBackupItem(index) {
    newApp.backup_items[index].enabled = !newApp.backup_items[index].enabled;
    newApp.backup_items = newApp.backup_items;
  }

  function addCustomBackupItem() {
    if (!customItemPath.trim()) return;
    
    newApp.backup_items = [...newApp.backup_items, {
      path: customItemPath.trim(),
      description: customItemDesc.trim() || 'Custom item',
      optional: true,
      enabled: true
    }];
    
    customItemPath = '';
    customItemDesc = '';
  }

  function removeBackupItem(index) {
    newApp.backup_items = newApp.backup_items.filter((_, i) => i !== index);
  }

  async function saveNewApp() {
    if (!newApp.app_name || !newApp.exe_path || !newApp.data_path) {
      alert('Please fill in all required fields');
      return;
    }

    // Filter enabled backup items
    const enabledItems = newApp.backup_items
      .filter(item => item.enabled)
      .map(item => ({
        type: item.path.includes('.') ? 'file' : 'folder',
        path: item.path,
        description: item.description,
        optional: item.optional
      }));

    if (enabledItems.length === 0) {
      alert('Please select at least one backup item');
      return;
    }

    const config = {
      app_name: newApp.app_name,
      display_name: newApp.display_name || newApp.app_name,
      version: '1.0',
      active: true,
      description: `${newApp.display_name} - Managed by SurfManager`,
      paths: {
        data_paths: [newApp.data_path],
        exe_paths: [newApp.exe_path],
        reset_folder: newApp.data_path
      },
      backup_items: enabledItems,
      addon_backup_paths: newApp.addon_paths
    };

    try {
      await SaveApp(config);
      log(`Added: ${config.display_name}`);
      showAddDialog = false;
      resetNewApp();
      await loadApps();
    } catch (e) {
      log(`Error: ${e}`);
    }
  }

  function resetNewApp() {
    newApp = {
      app_name: '',
      display_name: '',
      exe_path: '',
      data_path: '',
      app_type: 'vscode',
      backup_items: vscodePresetItems.map(item => ({ ...item })),
      addon_paths: []
    };
    customItemPath = '';
    customItemDesc = '';
  }

  // Initialize backup items when dialog opens
  $: if (showAddDialog && newApp.backup_items.length === 0) {
    newApp.backup_items = vscodePresetItems.map(item => ({ ...item }));
  }
</script>

<div class="h-full flex flex-col gap-4 animate-fadeIn">
  <div class="flex items-center justify-between">
    <h2 class="text-xl font-semibold text-[var(--text-primary)]">App Configuration</h2>
    <div class="flex items-center gap-3">
      <button 
        class="px-4 py-2 rounded-lg font-medium bg-[var(--primary)] hover:bg-[var(--primary-light)] hover:text-black text-white transition-all flex items-center gap-2"
        on:click={() => showAddDialog = true}
      >
        <Plus size={16} />
        Add App
      </button>
      <button 
        class="px-4 py-2 rounded-lg font-medium bg-[var(--bg-hover)] hover:bg-[var(--border)] border border-[var(--border)] text-[var(--text-secondary)] transition-all flex items-center gap-2"
        on:click={handleOpenFolder}
      >
        <FolderOpen size={16} />
        Open Folder
      </button>
      <button 
        class="p-2 rounded-lg text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-hover)] transition-colors"
        on:click={loadApps} 
        title="Refresh"
      >
        <RefreshCw size={18} class={loading ? 'animate-spin' : ''} />
      </button>
    </div>
  </div>

  <div class="flex-1 bg-[var(--bg-elevated)] rounded-xl border border-[var(--border)] overflow-hidden">
    <div class="overflow-auto h-full">
      <table class="w-full">
        <thead class="bg-[var(--bg-card)] sticky top-0">
          <tr class="text-left text-sm text-[var(--text-secondary)]">
            <th class="p-3 w-10">#</th>
            <th class="p-3">App Name</th>
            <th class="p-3">Status</th>
            <th class="p-3 w-48">Actions</th>
          </tr>
        </thead>
        <tbody>
          {#if apps.length === 0}
            <tr>
              <td colspan="4" class="p-8 text-center text-[var(--text-muted)]">
                No apps configured. Click "Add App" to get started.
              </td>
            </tr>
          {/if}
          {#each apps as app, i}
            <tr 
              class="border-b border-[var(--border)] hover:bg-[var(--bg-card)] transition-colors cursor-pointer"
              on:contextmenu={(e) => handleContextMenu(e, app)}
            >
              <td class="p-3 text-[var(--text-muted)]">{i + 1}</td>
              <td class="p-3">
                <span class="font-medium text-[var(--text-primary)]">{app.display_name}</span>
                <span class="text-xs text-[var(--text-muted)] ml-2">({app.app_name})</span>
              </td>
              <td class="p-3">
                <div class="flex gap-1">
                  <span class="px-2 py-0.5 rounded-full text-xs font-medium
                              {app.installed ? 'bg-[var(--success)]/20 text-[var(--success)]' : 'bg-[var(--danger)]/20 text-[var(--danger)]'}">
                    {app.installed ? 'Installed' : 'Not Found'}
                  </span>
                  {#if !app.active}
                    <span class="px-2 py-0.5 rounded-full text-xs font-medium bg-[var(--warning)]/20 text-[var(--warning)]">
                      Inactive
                    </span>
                  {/if}
                </div>
              </td>
              <td class="p-3">
                <div class="flex items-center gap-2">
                  <button
                    class="px-3 py-1 rounded text-xs font-medium transition-all
                           {app.active 
                             ? 'bg-[var(--success)]/20 text-[var(--success)] hover:bg-[var(--success)]/30' 
                             : 'bg-[var(--bg-hover)] text-[var(--text-muted)] hover:bg-[var(--border)]'}"
                    on:click={() => handleToggle(app)}
                  >
                    {app.active ? 'Active' : 'Inactive'}
                  </button>
                  <button
                    class="p-1.5 rounded text-[var(--text-secondary)] hover:text-[var(--danger)] hover:bg-[var(--bg-hover)] transition-colors"
                    on:click={() => handleDelete(app)}
                    title="Delete"
                  >
                    <Trash2 size={14} />
                  </button>
                </div>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  </div>

  <div class="text-sm text-[var(--text-muted)]">
    <p>Configs stored in: <code class="text-[var(--text-secondary)] bg-[var(--bg-hover)] px-2 py-0.5 rounded">~/.surfmanager/AppConfigs/</code></p>
  </div>
</div>

<!-- Context Menu -->
{#if contextMenu.show}
  <div 
    class="fixed bg-[var(--bg-card)] border border-[var(--border)] rounded-lg shadow-xl py-1 z-50 min-w-[180px]"
    style="left: {contextMenu.x}px; top: {contextMenu.y}px"
    on:click|stopPropagation={closeContextMenu}
  >
    <div class="px-3 py-2 border-b border-[var(--border)]">
      <span class="text-sm font-medium text-[var(--text-primary)]">{contextMenu.app?.display_name}</span>
    </div>
    
    <button
      class="w-full px-3 py-2 text-left text-sm text-[var(--text-secondary)] hover:bg-[var(--bg-hover)] hover:text-[var(--success)] transition-colors flex items-center gap-2"
      on:click={() => { handleSetActive(contextMenu.app); closeContextMenu(); }}
    >
      <Check size={14} />
      Set Active
    </button>
    <button
      class="w-full px-3 py-2 text-left text-sm text-[var(--text-secondary)] hover:bg-[var(--bg-hover)] hover:text-[var(--warning)] transition-colors flex items-center gap-2"
      on:click={() => { handleSetInactive(contextMenu.app); closeContextMenu(); }}
    >
      <ToggleLeft size={14} />
      Set Inactive
    </button>
    
    <div class="border-t border-[var(--border)] my-1"></div>
    
    <button
      class="w-full px-3 py-2 text-left text-sm text-[var(--text-secondary)] hover:bg-[var(--bg-hover)] hover:text-[var(--text-primary)] transition-colors flex items-center gap-2"
      on:click={() => { handleOpenAppConfig(contextMenu.app); closeContextMenu(); }}
    >
      <FileJson size={14} />
      Edit JSON Config
    </button>
    <button
      class="w-full px-3 py-2 text-left text-sm text-[var(--text-secondary)] hover:bg-[var(--bg-hover)] hover:text-[var(--text-primary)] transition-colors flex items-center gap-2"
      on:click={() => { handleOpenFolder(); closeContextMenu(); }}
    >
      <FolderOpen size={14} />
      Open Config Folder
    </button>
    
    <div class="border-t border-[var(--border)] my-1"></div>
    
    <button
      class="w-full px-3 py-2 text-left text-sm text-[var(--danger)] hover:bg-[var(--danger)]/10 transition-colors flex items-center gap-2"
      on:click={() => { handleDelete(contextMenu.app); closeContextMenu(); }}
    >
      <Trash2 size={14} />
      Delete App
    </button>
  </div>
{/if}

<!-- Add App Dialog -->
{#if showAddDialog}
  <div class="fixed inset-0 bg-black/60 flex items-center justify-center z-50" on:click|self={() => { showAddDialog = false; resetNewApp(); }}>
    <div class="bg-[var(--bg-card)] rounded-xl border border-[var(--border)] w-full max-w-4xl p-6 animate-fadeIn max-h-[90vh] overflow-auto">
      <div class="flex items-center justify-between mb-4">
        <h3 class="text-lg font-semibold text-[var(--text-primary)]">Add Application</h3>
        <button 
          class="p-1 rounded hover:bg-[var(--bg-hover)] text-[var(--text-muted)]"
          on:click={() => { showAddDialog = false; resetNewApp(); }}
        >
          <X size={20} />
        </button>
      </div>
      
      <!-- Grid Layout -->
      <div class="grid grid-cols-2 gap-6">
        <!-- Left Column -->
        <div class="space-y-4">
          <!-- App Name -->
          <div>
            <label class="block text-sm text-[var(--text-secondary)] mb-1">App Name</label>
            <input
              type="text"
              class="w-full bg-[var(--bg-elevated)] border border-[var(--border)] rounded-lg px-3 py-2 text-[var(--text-primary)] placeholder-[var(--text-muted)]"
              bind:value={newApp.display_name}
              placeholder="Auto-filled from .exe"
            />
          </div>

          <!-- App Type Toggle -->
          <div>
            <label class="block text-sm text-[var(--text-secondary)] mb-2">App Type</label>
            <div class="grid grid-cols-2 gap-2">
              <button
                class="p-2.5 rounded-lg border-2 text-left transition-all
                       {newApp.app_type === 'vscode' 
                         ? 'border-[var(--primary)] bg-[var(--primary-dim)]' 
                         : 'border-[var(--border)] hover:border-[var(--border-hover)]'}"
                on:click={() => setAppType('vscode')}
              >
                <div class="flex items-center gap-2">
                  <div class="w-3 h-3 rounded-full border-2 flex items-center justify-center
                              {newApp.app_type === 'vscode' ? 'border-[var(--primary)]' : 'border-[var(--text-muted)]'}">
                    {#if newApp.app_type === 'vscode'}
                      <div class="w-1.5 h-1.5 rounded-full bg-[var(--primary)]"></div>
                    {/if}
                  </div>
                  <span class="text-sm font-medium text-[var(--text-primary)]">VSCode Preset</span>
                </div>
              </button>
              <button
                class="p-2.5 rounded-lg border-2 text-left transition-all
                       {newApp.app_type === 'custom' 
                         ? 'border-[var(--primary)] bg-[var(--primary-dim)]' 
                         : 'border-[var(--border)] hover:border-[var(--border-hover)]'}"
                on:click={() => setAppType('custom')}
              >
                <div class="flex items-center gap-2">
                  <div class="w-3 h-3 rounded-full border-2 flex items-center justify-center
                              {newApp.app_type === 'custom' ? 'border-[var(--primary)]' : 'border-[var(--text-muted)]'}">
                    {#if newApp.app_type === 'custom'}
                      <div class="w-1.5 h-1.5 rounded-full bg-[var(--primary)]"></div>
                    {/if}
                  </div>
                  <span class="text-sm font-medium text-[var(--text-primary)]">Custom</span>
                </div>
              </button>
            </div>
          </div>

          <!-- Executable -->
          <div>
            <label class="block text-sm text-[var(--text-secondary)] mb-1">Executable *</label>
            <div class="flex gap-2">
              <input
                type="text"
                class="flex-1 bg-[var(--bg-hover)] border border-[var(--border)] rounded-lg px-3 py-2 text-[var(--text-primary)] placeholder-[var(--text-muted)] text-sm"
                bind:value={newApp.exe_path}
                placeholder="Select .exe file"
                readonly
              />
              <button 
                class="px-3 py-2 rounded-lg font-medium bg-[var(--bg-hover)] hover:bg-[var(--border)] border border-[var(--border)] text-[var(--text-secondary)] transition-all text-sm"
                on:click={selectExe}
              >
                Browse
              </button>
            </div>
          </div>

          <!-- Data Folder -->
          <div>
            <label class="block text-sm text-[var(--text-secondary)] mb-1">Data Folder *</label>
            <div class="flex gap-2">
              <input
                type="text"
                class="flex-1 bg-[var(--bg-hover)] border border-[var(--border)] rounded-lg px-3 py-2 text-[var(--text-primary)] placeholder-[var(--text-muted)] text-sm"
                bind:value={newApp.data_path}
                placeholder="AppData/Roaming/..."
                readonly
              />
              <button 
                class="px-3 py-2 rounded-lg font-medium bg-[var(--bg-hover)] hover:bg-[var(--border)] border border-[var(--border)] text-[var(--text-secondary)] transition-all text-sm"
                on:click={selectDataFolder}
              >
                Browse
              </button>
            </div>
          </div>

          <!-- Additional Folders -->
          <div>
            <label class="block text-sm text-[var(--text-secondary)] mb-1">📁 Additional Folders</label>
            <p class="text-xs text-[var(--text-muted)] mb-2">Also backed up & restored (e.g., ~/.aws)</p>
            <div class="space-y-1.5 max-h-24 overflow-auto">
              {#each newApp.addon_paths as path}
                <div class="flex items-center gap-2 bg-[var(--bg-hover)] rounded px-2 py-1.5 text-xs">
                  <span class="flex-1 truncate text-[var(--text-secondary)] font-mono">{path}</span>
                  <button class="text-[var(--danger)] hover:text-[var(--danger)]/80" on:click={() => removeAddonFolder(path)}>
                    <X size={14} />
                  </button>
                </div>
              {/each}
            </div>
            <button 
              class="w-full mt-2 px-3 py-1.5 rounded-lg text-sm bg-[var(--bg-hover)] hover:bg-[var(--border)] border border-[var(--border)] border-dashed text-[var(--text-secondary)] transition-all"
              on:click={addAddonFolder}
            >
              + Add Folder
            </button>
          </div>
        </div>

        <!-- Right Column: Backup Items -->
        <div>
          <label class="block text-sm text-[var(--text-secondary)] mb-2">📦 Backup Items</label>
          <div class="bg-[var(--bg-elevated)] rounded-lg border border-[var(--border)] p-3 h-[320px] overflow-auto">
            {#if newApp.backup_items.length === 0}
              <p class="text-sm text-[var(--text-muted)] text-center py-4">No backup items. Add below.</p>
            {/if}
            <div class="space-y-1.5">
              {#each newApp.backup_items as item, index}
                <div 
                  class="flex items-center gap-2 p-2 rounded hover:bg-[var(--bg-hover)] transition-colors
                         {item.enabled ? '' : 'opacity-50'}"
                >
                  <button
                    class="w-4 h-4 rounded border flex items-center justify-center transition-all flex-shrink-0
                           {item.enabled 
                             ? 'bg-[var(--primary)] border-[var(--primary)]' 
                             : 'border-[var(--border)] hover:border-[var(--text-muted)]'}"
                    on:click={() => toggleBackupItem(index)}
                  >
                    {#if item.enabled}
                      <Check size={10} class="text-white" />
                    {/if}
                  </button>
                  <div class="flex-1 min-w-0">
                    <span class="text-xs text-[var(--text-primary)] font-mono">{item.path}</span>
                    <span class="text-[10px] text-[var(--text-muted)] ml-1">{item.description}</span>
                  </div>
                  <span class="text-[9px] px-1 py-0.5 rounded flex-shrink-0 {item.optional ? 'bg-[var(--bg-hover)] text-[var(--text-muted)]' : 'bg-[var(--warning)]/20 text-[var(--warning)]'}">
                    {item.optional ? 'Opt' : 'Req'}
                  </span>
                  {#if newApp.app_type === 'custom' || !vscodePresetItems.find(p => p.path === item.path)}
                    <button
                      class="p-0.5 rounded text-[var(--text-muted)] hover:text-[var(--danger)] flex-shrink-0"
                      on:click={() => removeBackupItem(index)}
                    >
                      <X size={12} />
                    </button>
                  {/if}
                </div>
              {/each}
            </div>
            
            <!-- Add Custom Item -->
            <div class="flex gap-2 mt-3 pt-3 border-t border-[var(--border)]">
              <input
                type="text"
                class="flex-1 bg-[var(--bg-hover)] border border-[var(--border)] rounded px-2 py-1 text-xs text-[var(--text-primary)] placeholder-[var(--text-muted)]"
                placeholder="Folder/file"
                bind:value={customItemPath}
              />
              <input
                type="text"
                class="flex-1 bg-[var(--bg-hover)] border border-[var(--border)] rounded px-2 py-1 text-xs text-[var(--text-primary)] placeholder-[var(--text-muted)]"
                placeholder="Description"
                bind:value={customItemDesc}
              />
              <button 
                class="px-2 py-1 rounded text-xs font-medium bg-[var(--primary)] hover:bg-[var(--primary-light)] hover:text-black text-white transition-all"
                on:click={addCustomBackupItem}
              >
                Add
              </button>
            </div>
          </div>
        </div>
      </div>

      <div class="flex justify-end gap-3 mt-6 pt-4 border-t border-[var(--border)]">
        <button 
          class="px-4 py-2 rounded-lg font-medium bg-[var(--bg-hover)] hover:bg-[var(--border)] border border-[var(--border)] text-[var(--text-secondary)] transition-all"
          on:click={() => { showAddDialog = false; resetNewApp(); }}
        >
          Cancel
        </button>
        <button 
          class="px-4 py-2 rounded-lg font-medium bg-[var(--primary)] hover:bg-[var(--primary-light)] hover:text-black text-white transition-all"
          on:click={saveNewApp}
        >
          Save
        </button>
      </div>
    </div>
  </div>
{/if}
