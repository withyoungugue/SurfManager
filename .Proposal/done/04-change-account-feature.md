# Proposal: Change Account Feature (Quick Session Switch)

## Overview

Fitur baru untuk switch akun dengan cepat tanpa full restore. Hanya replace file-file yang berkaitan dengan login/session saja.

## Problem Statement

Full restore itu:
- Lama (copy semua data)
- Overkill kalau cuma mau ganti akun
- Butuh close app dulu

User mau cara cepat untuk switch akun di VSCode/Cursor **TANPA harus close app**.

## Proposed Solution

### File yang akan di-replace:

```
User/globalStorage/
└── state.vscdb                      # SQLite database untuk auth state
```

**Cuma 1 file!** File ini menyimpan semua auth tokens dan session data yang diperlukan untuk login.

---

## 🔬 Research: Bisa Gak Tanpa Close App?

### Ide Awal: Suspend/Freeze Process

Idenya adalah "freeze" process sementara (suspend), replace file, lalu resume.

**Hasil Research:**
> "When a process is suspended, the locks it has on the Dlls it references are not freed."
> — [SuperUser](https://superuser.com/questions/1056944/windows-10-processes-stuck-in-suspended-state)

**Kesimpulan:** ❌ **TIDAK BISA** - Suspend process TIDAK melepas file locks!

Windows API `NtSuspendProcess` hanya menghentikan eksekusi thread, tapi file handles tetap terbuka dan locked.

---

### Ide Alternatif 1: Force Close File Handles

Ada cara untuk close file handle di process lain menggunakan `DuplicateHandle` dengan flag `DUPLICATE_CLOSE_SOURCE`.

**Cara kerja:**
```
1. Enumerate semua handles di target process (NtQuerySystemInformation)
2. Find handle yang match dengan file kita
3. DuplicateHandle dengan DUPLICATE_CLOSE_SOURCE untuk close handle di process lain
4. Replace file
5. ??? App akan crash/error karena handle-nya hilang
```

**Masalah:**
- ⚠️ **SANGAT BERBAHAYA** - App bisa crash atau corrupt data
- ⚠️ Butuh admin privileges
- ⚠️ SQLite akan error karena connection-nya tiba-tiba invalid
- ⚠️ Chromium (yang handle Cookies) bisa crash

**Kesimpulan:** ❌ **TIDAK RECOMMENDED** - Terlalu risky, bisa corrupt data

---

### Ide Alternatif 2: Copy File Tanpa Lock (SQLite Immutable Mode)

Untuk **BACA** SQLite yang locked, bisa pakai URI dengan `?immutable=1`:
```
file:state.vscdb?immutable=1
```

**Tapi ini cuma untuk BACA, bukan WRITE/REPLACE!**

**Kesimpulan:** ❌ **TIDAK APPLICABLE** - Kita butuh replace file, bukan baca

---

### Ide Alternatif 3: Shadow Copy (VSS)

Windows Volume Shadow Copy Service bisa bikin snapshot volume.

**Masalah:**
- Butuh admin privileges
- Complex implementation
- Overkill untuk use case ini
- Tetap tidak bisa WRITE ke file yang locked

**Kesimpulan:** ❌ **OVERKILL**

---

## 📊 Summary Research

| Approach | Bisa? | Alasan |
|----------|-------|--------|
| Suspend Process | ❌ | File locks TIDAK dilepas saat suspend |
| Force Close Handles | ❌ | Berbahaya, bisa crash/corrupt |
| SQLite Immutable | ❌ | Hanya untuk read, bukan write |
| Shadow Copy (VSS) | ❌ | Overkill, tetap tidak bisa write |
| **Close App First** | ✅ | **Satu-satunya cara yang safe** |

---

## 🎯 Final Recommendation

**App HARUS di-close dulu.** Tidak ada cara yang safe untuk replace file yang locked tanpa close app.

### Kenapa?

1. **SQLite `state.vscdb`** - Database ini di-lock exclusive oleh VSCode/Cursor. Kalau kita force close handle-nya, SQLite connection di app akan invalid dan bisa corrupt database.

2. **Chromium `Cookies`** - File ini di-manage oleh Chromium engine. Force close bisa bikin browser component crash.

3. **Data Integrity** - Kalau app masih running dan kita replace file, app bisa overwrite lagi dengan data lama saat save.

---

## ✅ Implementation Plan (Safe Approach)

### Phase 1: Backend (Go)

#### 1.1 Add `ChangeAccount()` function di `app.go`

```go
// ChangeAccount performs a quick account switch by replacing only auth-related files
func (a *App) ChangeAccount(appKey, sessionName string) error {
    // 1. Get app config
    // 2. Close app (REQUIRED - tidak ada opsi skip)
    // 3. Call backup.QuickRestore() untuk replace file-file tertentu
    // 4. Return success
}
```

#### 1.2 Add `QuickRestore()` function di `internal/backup/backup.go`

```go
var quickRestoreFile = "User/globalStorage/state.vscdb"

// QuickRestore restores only the auth state file for quick account switch
func (m *Manager) QuickRestore(appKey, sessionName, targetPath string, progressCb ProgressCallback) error {
    // 1. Validate backup exists
    // 2. Delete existing state.vscdb in target
    // 3. Copy state.vscdb from backup to target
    // 4. Return success
}
```

### Phase 2: Frontend (Svelte)

Add "Change Account" option di context menu dengan handler yang SELALU close app dulu.

### UX Flow

```
User right-click session
    ↓
Click "Change Account"
    ↓
Confirm dialog: "This will close [App] and switch account. Continue?"
    ↓
App closes automatically
    ↓
Replace 1 file only (SUPER FAST!)
    ↓
Done! User reopens app with new account
```

---

## 💡 Benefit vs Full Restore

| Aspect | Full Restore | Change Account |
|--------|--------------|----------------|
| Files copied | ALL (100+ MB) | **1 file** (~2-5 MB) |
| Time | 10-30 seconds | **< 1 second** |
| Extensions | Reset | Preserved |
| Settings | Reset | Preserved |
| Workspaces | Reset | Preserved |

**Change Account SUPER cepat karena hanya replace 1 file!**

---

## Questions for Review

1. **Apakah acceptable harus close app?**
   - Berdasarkan research, ini satu-satunya cara yang safe

2. **Mau coba approach yang risky (force close handles)?**
   - Bisa diimplementasi tapi dengan warning besar bahwa bisa corrupt data

---

**Status:** ✅ IMPLEMENTED (Experimental)  
**Author:** Kiro  
**Date:** 2026-01-04

## Implementation Summary

Files changed:
- `app.go` - Added `RestoreAccountOnly()` function
- `internal/backup/backup.go` - Added `RestoreAccountOnly()` function  
- `frontend/src/lib/SessionsTab.svelte` - Added context menu option "Restore Account Only"

---

## ⚠️ Experimental Feature

Fitur "Restore Account Only" adalah **EXPERIMENTAL** dan perlu diaktifkan secara manual di Settings.

### Kenapa Experimental?

1. **Belum fully tested** - Perlu testing lebih lanjut di berbagai skenario
2. **Potential issues** - Mungkin ada edge cases yang belum ter-handle
3. **User awareness** - User perlu tau bahwa ini fitur experimental

### How to Enable

1. Buka **Settings** tab
2. Scroll ke section **Experimental Features**
3. Enable toggle **"Restore Account Only"**
4. Fitur akan muncul di context menu Sessions

### Settings Implementation

```javascript
// Di settings store, tambah:
experimentalRestoreAccountOnly: false

// Di SessionsTab, check setting sebelum show menu:
{#if $settings.experimentalRestoreAccountOnly}
  <button on:click={handleRestoreAccountOnly}>
    👤 Restore Account Only
  </button>
{/if}
```

### Files to Update

| File | Change |
|------|--------|
| `frontend/src/lib/stores/settings.js` | Add `experimentalRestoreAccountOnly: false` |
| `frontend/src/lib/SettingsTab.svelte` | Add toggle for experimental feature |
| `frontend/src/lib/SessionsTab.svelte` | Wrap menu item with setting check |
