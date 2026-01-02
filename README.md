# SurfManager v2.0

> **Advanced Session & Data Manager for Development Tools**

[![Version](https://img.shields.io/badge/version-2.0.0-brightgreen.svg)](https://github.com/risunCode/SurfManager)
[![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20macOS%20%7C%20Linux-blue.svg)](https://github.com/risunCode/SurfManager)
[![Go](https://img.shields.io/badge/go-1.21+-00ADD8.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

---

## 👋 Welcome to SurfManager!

**SurfManager** is a modern solution for managing session data of development tools like VS Code, Cursor, Windsurf, and similar applications. Built with Go + Wails + Svelte for blazing fast performance and a beautiful native UI.

Perfect for developers who need to:
- 🔄 Switch between multiple accounts/profiles effortlessly
- 💾 Backup workspace settings before experimenting
- 🚀 Maintain organized development workflows
- 🛡️ Have a safety net for important configurations

---

## ✨ Features

| Feature | Description |
|---------|-------------|
| **📱 Session Management** | Backup, restore, and manage multiple app sessions |
| **🔄 Account Switching** | Switch between different accounts in seconds |
| **�️ Smairt App Close** | Auto-close running apps before operations (optional) |
| **� Proogress Tracking** | Real-time progress bars for all operations |
| **� Seaarch & Filter** | Quick search through sessions and auto-backups |
| **� Auto-eBackup** | Automatic backup before reset operations |
| **🎨 Theme System** | Dark, Solarized Dark, and Solarized Light themes |
| **⚙️ Customizable Settings** | Persistent settings for personalized experience |

### 🚀 What's New in v2.0

- **Complete Rewrite** - Go + Wails + Svelte (from Python + PyQt6)
- **3x Smaller Binary** - ~15MB (from 40+MB)
- **2x Faster Startup** - <0.5s (from ~1s)
- **Modern UI** - JetBrains Mono font, realtime clock, custom modals
- **Theme System** - 3 beautiful themes with persistence
- **Skip Close App** - Perform operations without closing target app
- **Customizable Backups** - Choose exactly what to backup
- **Additional Folders** - Backup extra directories (e.g., ~/.aws, ~/.ssh)

---

## 📁 App Data Locations

SurfManager manages app data stored in platform-specific locations:

| Platform | App Data (Config) | App Data (Local) | Example Apps |
|----------|-------------------|------------------|--------------|
| **Windows** | `%APPDATA%` (`C:\Users\{user}\AppData\Roaming`) | `%LOCALAPPDATA%` (`C:\Users\{user}\AppData\Local`) | `Roaming\Code`, `Roaming\Cursor` |
| **macOS** | `~/Library/Application Support` | `~/Library/Application Support` | `Application Support/Code` |
| **Linux** | `~/.config` | `~/.local/share` | `~/.config/Code`, `~/.config/Cursor` |

### SurfManager Storage Locations

| Data | Windows | macOS | Linux |
|------|---------|-------|-------|
| **Backups** | `Documents\SurfManager\backup` | `~/Documents/SurfManager/backup` | `~/Documents/SurfManager/backup` |
| **Auto-Backups** | `Documents\SurfManager\auto-backups` | `~/Documents/SurfManager/auto-backups` | `~/Documents/SurfManager/auto-backups` |
| **Notes** | `Documents\SurfManager\notes` | `~/Documents/SurfManager/notes` | `~/Documents/SurfManager/notes` |
| **App Configs** | `~\.surfmanager\AppConfigs` | `~/.surfmanager/AppConfigs` | `~/.surfmanager/AppConfigs` |

---

## 📖 How to Use

### Quick Start

**Step 1: Setup First Account**
1. Login to your IDE (VS Code/Cursor/Windsurf)
2. Configure your workspace, install extensions
3. Open SurfManager → Sessions → New Backup
4. Enter session name (e.g., "work-account")

**Step 2: Add More Accounts**
1. Go to Reset tab → Click Reset on your app
2. Login with different credentials in your IDE
3. Create another backup (e.g., "personal-account")

**Step 3: Switch Between Accounts**
1. Go to Sessions tab
2. Right-click session → Restore
3. Launch your IDE - you're logged in! 🎉

### Tips

- **Right-click** anywhere for context menus
- **CTRL+Click** rows to select multiple items
- **Enable "Skip Close App"** if you want to backup while app is running
- **Use descriptive names** like "work-main", "personal-dev"

---

## ⚠️ Limitations

### Platform-Specific

**Windows User Isolation**
Sessions are tied to the Windows user account. You cannot transfer backups between different Windows users due to encryption.

- ✅ Switch accounts on the same Windows user
- ❌ Copy backups to another Windows user
- ✅ Each Windows user creates their own backups

---

## � Installation

### Download Pre-built Release (Recommended)

1. Visit [Releases Page](https://github.com/risunCode/SurfManager/releases)
2. Download for your platform:
   - Windows: `SurfManager-windows-amd64.exe`
   - macOS: `SurfManager-darwin-amd64` (needs testing)
   - Linux: `SurfManager-linux-amd64` (needs testing)
3. Run directly - no installation required!

---

## � Building from Source

### Prerequisites (All Platforms)

- **Go 1.21+** - [Download](https://golang.org/dl/)
- **Node.js 18+** - [Download](https://nodejs.org/)
- **Wails CLI** - Install with: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

### Windows

```powershell
# Install prerequisites
# 1. Install Go from https://golang.org/dl/
# 2. Install Node.js from https://nodejs.org/
# 3. Install Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Clone and build
git clone https://github.com/risunCode/SurfManager.git
cd SurfManager
cd frontend && npm install && cd ..

# Development mode
wails dev

# Build for production
wails build
# Output: build/bin/SurfManager.exe
```

### macOS

```bash
# Install prerequisites
# Using Homebrew:
brew install go node

# Install Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Install Xcode Command Line Tools (required for CGO)
xcode-select --install

# Clone and build
git clone https://github.com/risunCode/SurfManager.git
cd SurfManager
cd frontend && npm install && cd ..

# Development mode
wails dev

# Build for production
wails build
# Output: build/bin/SurfManager.app
```

### Linux (Ubuntu/Debian)

```bash
# Install prerequisites
sudo apt update
sudo apt install -y golang nodejs npm build-essential libgtk-3-dev libwebkit2gtk-4.0-dev

# Install Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Add Go bin to PATH (add to ~/.bashrc for persistence)
export PATH=$PATH:$(go env GOPATH)/bin

# Clone and build
git clone https://github.com/risunCode/SurfManager.git
cd SurfManager
cd frontend && npm install && cd ..

# Development mode
wails dev

# Build for production
wails build
# Output: build/bin/SurfManager
```

### Linux (Fedora/RHEL)

```bash
# Install prerequisites
sudo dnf install -y golang nodejs npm gcc gtk3-devel webkit2gtk3-devel

# Install Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Add Go bin to PATH
export PATH=$PATH:$(go env GOPATH)/bin

# Clone and build (same as Ubuntu)
git clone https://github.com/risunCode/SurfManager.git
cd SurfManager
cd frontend && npm install && cd ..
wails build
```

### Linux (Arch)

```bash
# Install prerequisites
sudo pacman -S go nodejs npm base-devel gtk3 webkit2gtk

# Install Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Add Go bin to PATH
export PATH=$PATH:$(go env GOPATH)/bin

# Clone and build
git clone https://github.com/risunCode/SurfManager.git
cd SurfManager
cd frontend && npm install && cd ..
wails build
```

### Verify Wails Installation

```bash
wails doctor
```

This will check if all dependencies are installed correctly.

---

## 🆘 Help Wanted: Linux & macOS Compatibility

**We need your help!** SurfManager is primarily developed and tested on Windows. We need contributors to help with:

### Linux
- [ ] Test app data path detection (`~/.config`, `~/.local/share`)
- [ ] Test VSCode/Cursor data locations
- [ ] Verify file dialogs work correctly
- [ ] Test process detection and termination
- [ ] Package for different distributions (AppImage, Flatpak, Snap)

### macOS
- [ ] Test app data path detection (`~/Library/Application Support`)
- [ ] Test VSCode/Cursor data locations  
- [ ] Verify `.app` bundle selection works
- [ ] Test process detection and termination
- [ ] Code signing and notarization

### How to Contribute

1. Fork the repository
2. Test on your Linux/macOS machine
3. Report issues with detailed logs
4. Submit PRs for fixes

**Even just testing and reporting issues helps a lot!** 🙏
 
---

## 🤝 Contributing

We welcome contributions! Whether it's bug reports, feature suggestions, or code contributions.

```bash
git checkout -b feature/awesome-feature
git commit -m 'Add awesome feature'
git push origin feature/awesome-feature
# Open a Pull Request
```

---

## 📄 License

SurfManager is open-source under the MIT License.

---

## 🙏 Credits

**Built with ❤️ by risunCode**

**Technologies:** Go, Wails v2, Svelte, TailwindCSS, Lucide Icons

---

<div align="center">

**SurfManager v2.0**

*Making development workflows smoother, one session at a time*

[⭐ Star on GitHub](https://github.com/risunCode/SurfManager) | [🐛 Report Issues](https://github.com/risunCode/SurfManager/issues) | [💡 Suggest Features](https://github.com/risunCode/SurfManager/discussions)

</div>
