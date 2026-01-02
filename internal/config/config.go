// Package config provides cross-platform configuration and path resolution for SurfManager.
// Supports: Windows, Linux, macOS
// Handles: Platform detection, user directories, app data locations
package config

import (
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

// Platform constants
const (
	PlatformWindows = "windows"
	PlatformLinux   = "linux"
	PlatformMacOS   = "darwin"
)

// PlatformPaths holds cross-platform path configuration
type PlatformPaths struct {
	HomeDir   string
	AppData   string
	AppConfig string
	Documents string
	Downloads string
	Desktop   string
}

// Manager handles configuration and path resolution
type Manager struct {
	platform    string
	homeDir     string
	currentUser string
	paths       *PlatformPaths
}

var (
	instance *Manager
	once     sync.Once
)

// GetManager returns the singleton Manager instance
func GetManager() *Manager {
	once.Do(func() {
		instance = newManager()
	})
	return instance
}

// newManager creates a new Manager instance
func newManager() *Manager {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = ""
	}

	currentUser := os.Getenv("USER")
	if currentUser == "" {
		currentUser = os.Getenv("USERNAME")
	}
	if currentUser == "" {
		currentUser = "user"
	}

	m := &Manager{
		platform:    runtime.GOOS,
		homeDir:     homeDir,
		currentUser: currentUser,
	}

	m.paths = m.getPlatformPaths("")
	return m
}

// GetPlatform returns the current platform (windows, linux, darwin)
func (m *Manager) GetPlatform() string {
	return m.platform
}

// IsWindows returns true if running on Windows
func (m *Manager) IsWindows() bool {
	return m.platform == PlatformWindows
}

// IsLinux returns true if running on Linux
func (m *Manager) IsLinux() bool {
	return m.platform == PlatformLinux
}

// IsMacOS returns true if running on macOS
func (m *Manager) IsMacOS() bool {
	return m.platform == PlatformMacOS
}

// GetHomeDir returns the current user's home directory
func (m *Manager) GetHomeDir() string {
	return m.homeDir
}

// GetCurrentUser returns the current username
func (m *Manager) GetCurrentUser() string {
	return m.currentUser
}

// getUserHome returns the home directory for a specific user
func (m *Manager) getUserHome(username string) string {
	if username == "" || username == m.currentUser {
		return m.homeDir
	}

	// For other users, we need to construct the path based on platform conventions
	switch m.platform {
	case PlatformWindows:
		// Try to get the system drive from environment or use the current user's home as reference
		systemDrive := os.Getenv("SystemDrive")
		if systemDrive == "" {
			// Extract drive letter from current user's home directory
			if len(m.homeDir) >= 2 && m.homeDir[1] == ':' {
				systemDrive = m.homeDir[:2]
			} else {
				systemDrive = "C:"
			}
		}
		return filepath.Join(systemDrive, "Users", username)
	case PlatformMacOS:
		return filepath.Join("/Users", username)
	default: // Linux
		return filepath.Join("/home", username)
	}
}

// getPlatformPaths returns platform-specific paths for a user
func (m *Manager) getPlatformPaths(username string) *PlatformPaths {
	home := m.getUserHome(username)

	switch m.platform {
	case PlatformWindows:
		return m.getWindowsPaths(home)
	case PlatformMacOS:
		return m.getMacOSPaths(home)
	default: // Linux
		return m.getLinuxPaths(home)
	}
}

// getWindowsPaths returns Windows-specific paths
func (m *Manager) getWindowsPaths(home string) *PlatformPaths {
	return &PlatformPaths{
		HomeDir:   home,
		AppData:   filepath.Join(home, "AppData", "Local"),
		AppConfig: filepath.Join(home, "AppData", "Roaming"),
		Documents: filepath.Join(home, "Documents"),
		Downloads: filepath.Join(home, "Downloads"),
		Desktop:   filepath.Join(home, "Desktop"),
	}
}

// getLinuxPaths returns Linux-specific paths
func (m *Manager) getLinuxPaths(home string) *PlatformPaths {
	return &PlatformPaths{
		HomeDir:   home,
		AppData:   filepath.Join(home, ".local", "share"),
		AppConfig: filepath.Join(home, ".config"),
		Documents: filepath.Join(home, "Documents"),
		Downloads: filepath.Join(home, "Downloads"),
		Desktop:   filepath.Join(home, "Desktop"),
	}
}

// getMacOSPaths returns macOS-specific paths
func (m *Manager) getMacOSPaths(home string) *PlatformPaths {
	return &PlatformPaths{
		HomeDir:   home,
		AppData:   filepath.Join(home, "Library", "Application Support"),
		AppConfig: filepath.Join(home, "Library", "Application Support"),
		Documents: filepath.Join(home, "Documents"),
		Downloads: filepath.Join(home, "Downloads"),
		Desktop:   filepath.Join(home, "Desktop"),
	}
}

// GetPlatformPaths returns the platform paths for the current user or specified user
func (m *Manager) GetPlatformPaths(username string) *PlatformPaths {
	if username == "" {
		return m.paths
	}
	return m.getPlatformPaths(username)
}

// GetDocumentsDir returns the user's Documents directory
func (m *Manager) GetDocumentsDir() string {
	return m.paths.Documents
}

// GetDocumentsDirForUser returns the Documents directory for a specific user
func (m *Manager) GetDocumentsDirForUser(username string) string {
	paths := m.getPlatformPaths(username)
	return paths.Documents
}

// GetAppDataDir returns the application data directory for an app
// Windows: AppData/Local/{appName}
// Linux: ~/.local/share/{appName}
// macOS: ~/Library/Application Support/{appName}
func (m *Manager) GetAppDataDir(appName string) string {
	return filepath.Join(m.paths.AppData, appName)
}

// GetAppDataDirForUser returns the app data directory for a specific user
func (m *Manager) GetAppDataDirForUser(appName, username string) string {
	paths := m.getPlatformPaths(username)
	return filepath.Join(paths.AppData, appName)
}

// GetAppConfigDir returns the application config directory for an app
// Windows: AppData/Roaming/{appName}
// Linux: ~/.config/{appName}
// macOS: ~/Library/Application Support/{appName}
func (m *Manager) GetAppConfigDir(appName string) string {
	return filepath.Join(m.paths.AppConfig, appName)
}

// GetAppConfigDirForUser returns the app config directory for a specific user
func (m *Manager) GetAppConfigDirForUser(appName, username string) string {
	paths := m.getPlatformPaths(username)
	return filepath.Join(paths.AppConfig, appName)
}

// GetSurfManagerDir returns the SurfManager root directory in Documents
func (m *Manager) GetSurfManagerDir() string {
	return filepath.Join(m.paths.Documents, "SurfManager")
}

// GetSurfManagerDirForUser returns the SurfManager directory for a specific user
func (m *Manager) GetSurfManagerDirForUser(username string) string {
	paths := m.getPlatformPaths(username)
	return filepath.Join(paths.Documents, "SurfManager")
}

// GetBackupDir returns the backup directory: Documents/SurfManager/backup
func (m *Manager) GetBackupDir() string {
	return filepath.Join(m.GetSurfManagerDir(), "backup")
}

// GetBackupDirForUser returns the backup directory for a specific user
func (m *Manager) GetBackupDirForUser(username string) string {
	return filepath.Join(m.GetSurfManagerDirForUser(username), "backup")
}

// GetAutoBackupDir returns the auto-backup directory: Documents/SurfManager/auto-backups
func (m *Manager) GetAutoBackupDir() string {
	return filepath.Join(m.GetSurfManagerDir(), "auto-backups")
}

// GetAutoBackupDirForUser returns the auto-backup directory for a specific user
func (m *Manager) GetAutoBackupDirForUser(username string) string {
	return filepath.Join(m.GetSurfManagerDirForUser(username), "auto-backups")
}

// GetNotesDir returns the notes directory: Documents/SurfManager/notes
func (m *Manager) GetNotesDir() string {
	return filepath.Join(m.GetSurfManagerDir(), "notes")
}

// GetNotesDirForUser returns the notes directory for a specific user
func (m *Manager) GetNotesDirForUser(username string) string {
	return filepath.Join(m.GetSurfManagerDirForUser(username), "notes")
}

// GetConfigDir returns the SurfManager config directory: ~/.surfmanager/AppConfigs
func (m *Manager) GetConfigDir() string {
	return filepath.Join(m.homeDir, ".surfmanager", "AppConfigs")
}

// GetConfigDirForUser returns the config directory for a specific user
func (m *Manager) GetConfigDirForUser(username string) string {
	home := m.getUserHome(username)
	return filepath.Join(home, ".surfmanager", "AppConfigs")
}

// GetSessionDir returns the session directory: Documents/SurfManager/backup/sessions
func (m *Manager) GetSessionDir() string {
	return filepath.Join(m.GetBackupDir(), "sessions")
}

// GetSessionDirForUser returns the session directory for a specific user
func (m *Manager) GetSessionDirForUser(username string) string {
	return filepath.Join(m.GetBackupDirForUser(username), "sessions")
}

// GetSessionConfigFile returns the path to sessions.json
func (m *Manager) GetSessionConfigFile() string {
	return filepath.Join(m.GetSurfManagerDir(), "sessions.json")
}

// ExpandUserPath expands ~ notation in paths
func (m *Manager) ExpandUserPath(path string) string {
	if len(path) == 0 {
		return path
	}

	if path[0] == '~' {
		if len(path) == 1 {
			return m.homeDir
		}
		if path[1] == '/' || path[1] == '\\' {
			return filepath.Join(m.homeDir, path[2:])
		}
	}

	return path
}

// ExpandUserPathForUser expands ~ notation for a specific user
func (m *Manager) ExpandUserPathForUser(path, username string) string {
	if len(path) == 0 {
		return path
	}

	home := m.getUserHome(username)

	if path[0] == '~' {
		if len(path) == 1 {
			return home
		}
		if path[1] == '/' || path[1] == '\\' {
			return filepath.Join(home, path[2:])
		}
	}

	return path
}

// NormalizePath normalizes a path for the current platform
func (m *Manager) NormalizePath(path string) string {
	// Clean the path
	path = filepath.Clean(path)

	// Convert to absolute if relative
	if !filepath.IsAbs(path) {
		cwd, err := os.Getwd()
		if err == nil {
			path = filepath.Join(cwd, path)
		}
	}

	return path
}

// EnsureDir creates a directory if it doesn't exist
func (m *Manager) EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// EnsureSurfManagerDirs creates all SurfManager directories
func (m *Manager) EnsureSurfManagerDirs() error {
	dirs := []string{
		m.GetSurfManagerDir(),
		m.GetBackupDir(),
		m.GetAutoBackupDir(),
		m.GetNotesDir(),
		m.GetSessionDir(),
		m.GetConfigDir(),
	}

	for _, dir := range dirs {
		if err := m.EnsureDir(dir); err != nil {
			return err
		}
	}

	return nil
}

// GetPathSeparator returns the OS-specific path separator
func (m *Manager) GetPathSeparator() string {
	return string(filepath.Separator)
}

// GetExecutableExtension returns the executable extension for the platform
func (m *Manager) GetExecutableExtension() string {
	if m.IsWindows() {
		return ".exe"
	}
	return ""
}

// GetPlatformInfo returns platform information
func (m *Manager) GetPlatformInfo() map[string]string {
	return map[string]string{
		"platform":  m.platform,
		"arch":      runtime.GOARCH,
		"goVersion": runtime.Version(),
		"homeDir":   m.homeDir,
		"user":      m.currentUser,
	}
}

// Convenience functions for package-level access

// GetDocumentsDir returns the user's Documents directory
func GetDocumentsDir() string {
	return GetManager().GetDocumentsDir()
}

// GetAppDataDir returns the application data directory for an app
func GetAppDataDir(appName string) string {
	return GetManager().GetAppDataDir(appName)
}

// GetBackupDir returns the backup directory
func GetBackupDir() string {
	return GetManager().GetBackupDir()
}

// GetAutoBackupDir returns the auto-backup directory
func GetAutoBackupDir() string {
	return GetManager().GetAutoBackupDir()
}

// GetNotesDir returns the notes directory
func GetNotesDir() string {
	return GetManager().GetNotesDir()
}

// GetConfigDir returns the config directory
func GetConfigDir() string {
	return GetManager().GetConfigDir()
}

// IsWindows returns true if running on Windows
func IsWindows() bool {
	return GetManager().IsWindows()
}

// IsLinux returns true if running on Linux
func IsLinux() bool {
	return GetManager().IsLinux()
}

// IsMacOS returns true if running on macOS
func IsMacOS() bool {
	return GetManager().IsMacOS()
}

// ExpandUserPath expands ~ notation in paths
func ExpandUserPath(path string) string {
	return GetManager().ExpandUserPath(path)
}
