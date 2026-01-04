// Settings store with localStorage persistence
import { writable } from 'svelte/store';

const STORAGE_KEY = 'surfmanager-settings';

const defaultSettings = {
  // Appearance
  theme: 'dark',
  
  // Behavior
  autoBackup: true,
  confirmBeforeReset: true,
  confirmBeforeDelete: true,
  confirmBeforeRestore: true,
  skipCloseApp: false, // Don't close app before reset/backup/restore
  
  // Sessions
  showAutoBackups: false,
  defaultSessionFilter: 'all',
  
  // Notes
  autoSaveNotes: false,
  autoSaveDelay: 3000,
  
  // Advanced
  logRetention: 100,
  rememberLastTab: true,
  lastActiveTab: 'reset',
  
  // Experimental Features
  experimentalRestoreAccountOnly: false, // Quick account switch (restore only state.vscdb)
};

function createSettingsStore() {
  // Get stored settings or defaults
  let stored = {};
  if (typeof localStorage !== 'undefined') {
    try {
      stored = JSON.parse(localStorage.getItem(STORAGE_KEY) || '{}');
    } catch (e) {
      stored = {};
    }
  }
  
  const initial = { ...defaultSettings, ...stored };
  const { subscribe, set, update } = writable(initial);

  return {
    subscribe,
    set: (settings) => {
      if (typeof localStorage !== 'undefined') {
        localStorage.setItem(STORAGE_KEY, JSON.stringify(settings));
      }
      set(settings);
    },
    update: (key, value) => {
      update(s => {
        const newSettings = { ...s, [key]: value };
        if (typeof localStorage !== 'undefined') {
          localStorage.setItem(STORAGE_KEY, JSON.stringify(newSettings));
        }
        return newSettings;
      });
    },
    reset: () => {
      if (typeof localStorage !== 'undefined') {
        localStorage.setItem(STORAGE_KEY, JSON.stringify(defaultSettings));
      }
      set(defaultSettings);
    },
    get: (key) => {
      let current;
      subscribe(s => current = s)();
      return current[key];
    }
  };
}

export const settings = createSettingsStore();
