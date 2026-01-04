# SurfManager - Virus Detection False Positive Explanation

## Summary

SurfManager.exe di-flag oleh beberapa antivirus/sandbox sebagai suspicious, padahal ini adalah **FALSE POSITIVE**. Dokumen ini menjelaskan kenapa detection terjadi dan kenapa app ini aman.

---

## Detection Details

### Sigma Rule Triggered:
> **"Suspicious Loading of Dbgcore/Dbghelp DLLs from Uncommon Location"**
> 
> Detects loading of dbgcore.dll or dbghelp.dll from uncommon locations such as user directories. These DLLs contain the MiniDumpWriteDump function, which can be abused for credential dumping purposes or in some cases for evading EDR/AV detection by suspending processes.

### VirusTotal Results:
- **Detection Score:** 2/72 (Very Low - False Positive)
- **Flagged by:** DeepInstinct (MALICIOUS), Trapmine (Malicious.moderate.ml.score)
- **Community Score:** Positive
- **Sandbox Results:** No actual malware behavior detected

**Note:** 2/72 = 2.7% detection rate. Industry standard considers <5% as likely false positive, especially when flagged by ML-based engines (DeepInstinct, Trapmine) rather than signature-based engines.

---

## Why DeepInstinct & Trapmine Flag This

### About These Engines:

| Engine | Type | Known For |
|--------|------|-----------|
| **DeepInstinct** | ML/AI-based | High false positive rate on Go/Wails apps |
| **Trapmine** | ML-based | Behavioral heuristics, flags process management |

### Why ML Engines Flag Wails/Go Apps:

1. **Go binaries are large** (~10MB) - ML sees this as suspicious
2. **Static linking** - Go includes all dependencies, looks "packed"
3. **WebView2 behavior** - Loading debug DLLs triggers heuristics
4. **Process management** - Killing processes looks like malware

### Common False Positives:

Many legitimate Go/Wails apps get flagged:
- Wails framework apps
- Electron apps
- Any app using WebView2
- Backup software with process management

---

## Why This Detection Happens (False Positive)

### 1. Wails Framework Uses WebView2

SurfManager dibangun dengan **Wails** framework yang menggunakan **Microsoft WebView2** untuk rendering UI. WebView2 adalah komponen resmi Microsoft yang:

- Berbasis Chromium/Edge
- Load berbagai system DLLs termasuk debug-related DLLs
- Digunakan oleh banyak aplikasi legitimate (Microsoft Teams, VS Code, dll)

```
Wails App
    └── WebView2 (Microsoft Edge Runtime)
            └── Loads dbghelp.dll, dbgcore.dll (normal behavior)
```

### 2. Process Management Features

SurfManager memiliki fitur untuk:
- **Close/Kill processes** sebelum backup/restore (untuk release file locks)
- **Check if app is running**

Fitur ini menggunakan library `gopsutil` yang legitimate, tapi behavior-nya mirip dengan malware yang juga perlu manage processes.

```go
// Legitimate code untuk close app sebelum backup
func (k *Killer) SmartClose(appName string, processNames []string) error {
    // Gracefully close app, then force kill if needed
}
```

### 3. Network Communication

Detected IPs adalah **Google/Microsoft services**:
- `173.194.193.94` - Google (untuk WebView2 updates/telemetry)
- `74.125.201.95` - Google
- `edge.microsoft.com` - Microsoft Edge/WebView2
- `config.edge.skype.com` - Microsoft services

**Ini semua legitimate Microsoft/Google endpoints**, bukan C2 servers!

### 4. Memory Pattern URLs

URLs yang terdeteksi di memory adalah dari:
- **WebView2 cache** - Normal browser behavior
- **PKI/Certificate URLs** (`c.pki.goog`) - Certificate validation
- **Microsoft telemetry** - WebView2 standard behavior

---

## Evidence This Is NOT Malware

### ✅ Open Source
- Source code 100% available di GitHub
- Anyone can audit the code
- No obfuscated or hidden code

### ✅ Very Low Detection Rate
- VirusTotal: **2/72 detections** (2.7% - likely false positive)
- Only ML-based engines flagged (DeepInstinct, Trapmine)
- No signature-based engines detected malware
- ML engines are known for high false positive rates on legitimate apps

### ✅ Legitimate Use Case
- App untuk backup/restore VSCode/Cursor sessions
- Process management needed to release file locks
- Network traffic is to Microsoft/Google (WebView2)

### ✅ Signed Dependencies
- Uses official Microsoft WebView2
- Uses well-known Go libraries (gopsutil, wails)

---

## Technical Breakdown

### DLLs Loaded (Normal for WebView2 apps):

| DLL | Purpose | Why Loaded |
|-----|---------|------------|
| `dbghelp.dll` | Debug helper | WebView2/Chromium crash reporting |
| `dbgcore.dll` | Debug core | WebView2/Chromium diagnostics |
| `webview2loader.dll` | WebView2 | UI rendering |

### Process Management Code:

```go
// File: internal/process/killer.go
// Purpose: Close apps gracefully before backup/restore

func (k *Killer) KillProcess(pid int) error {
    proc, err := process.NewProcess(int32(pid))
    // ... graceful termination, then force kill
}
```

**This is NOT credential dumping!** It's just closing apps so we can backup their data files.

---

## MITRE ATT&CK Explanation

| Detected Technique | Our Legitimate Use |
|-------------------|-------------------|
| T1059 (Command Execution) | Running backup/restore operations |
| T1027 (Obfuscation) | WebView2 JavaScript bundling (normal) |
| T1036 (Masquerading) | FALSE - we don't masquerade |
| T1082 (System Discovery) | Checking if apps are installed/running |
| T1071 (App Layer Protocol) | HTTPS to Google/Microsoft (WebView2) |

---

## How to Verify Safety

### 1. Check Source Code
```
https://github.com/[repo]/SurfManager
```
All code is open source and auditable.

### 2. Build From Source
```bash
git clone [repo]
cd SurfManager
wails build
```
You can build the exact same binary yourself.

### 3. Check Network Traffic
All network traffic goes to:
- `*.google.com` - WebView2
- `*.microsoft.com` - WebView2/Edge

No suspicious C2 servers.

---

## Recommendations for Users

### If Your Antivirus Flags This:

1. **Add Exception** - Whitelist `SurfManager.exe` in your antivirus
2. **Check VirusTotal** - Confirm 0 detections on actual malware engines
3. **Build From Source** - If paranoid, build the app yourself from GitHub

### For Antivirus Vendors:

This is a **false positive** caused by:
- WebView2 loading debug DLLs (normal Chromium behavior)
- Process management features (legitimate for backup software)
- Network traffic to Google/Microsoft (WebView2 telemetry)

Please whitelist this application or improve heuristics to not flag WebView2-based apps.

---

## Conclusion

SurfManager is **NOT malware**. The detection is a **false positive** caused by:

1. **WebView2** loading debug DLLs (standard Chromium behavior)
2. **Process management** features that look suspicious but are legitimate
3. **Network traffic** to Google/Microsoft services (WebView2)

The app is:
- ✅ Open source
- ✅ 2/72 VirusTotal detections (very low, ML-based only)
- ✅ Uses official Microsoft WebView2
- ✅ All network traffic to legitimate endpoints
- ✅ No signature-based engine detected malware

---

**Author:** SurfManager Team  
**Date:** 2026-01-04  
**Status:** False Positive - Safe to Use
