# Implementation Plan: UX Improvements v2.1.0

## Overview

This implementation plan covers UX improvements for SurfManager including alphabetical sorting, state persistence, enhanced restore workflows, cleanup of unused features, and bulk session management.

## Tasks

- [x] 1. Update Settings Store
  - [x] 1.1 Add new properties for last selected app persistence
    - Add `lastSelectedAppReset` and `lastSelectedAppSession` to defaultSettings
    - Remove `rememberLastTab`, `lastActiveTab`, `debugMode`, `skipDataFolder` properties
    - Update exportSettings and importSettings to handle new/removed properties
    - _Requirements: 2.1, 2.3, 7.2, 9.2_

- [x] 2. Implement Alphabetical Sorting
  - [x] 2.1 Sort apps in ResetTab
    - Sort apps array by display_name using localeCompare after loading
    - _Requirements: 1.1_
  - [x] 2.2 Sort apps in SessionsTab
    - Sort apps array by display_name using localeCompare after loading
    - _Requirements: 1.2_

- [x] 3. Implement Remember Last Selected App
  - [x] 3.1 Update ResetTab to persist and restore selection
    - On app selection, save to settings store via `lastSelectedAppReset`
    - On mount, restore selection from settings store
    - Fallback to first app if saved app doesn't exist
    - _Requirements: 2.1, 2.2, 2.5_
  - [x] 3.2 Update SessionsTab to persist and restore filter selection
    - On filter change, save to settings store via `lastSelectedAppSession`
    - On mount, restore filter from settings store
    - Fallback to 'all' if saved app doesn't exist
    - _Requirements: 2.3, 2.4, 2.5_

- [x] 4. Add Launch App to Session Context Menu
  - [x] 4.1 Add Launch App option to context menu
    - Add Play icon and "Launch App" button to context menu
    - Import LaunchApp from wailsjs bindings
    - _Requirements: 3.1_
  - [x] 4.2 Implement handleLaunchApp function
    - Call LaunchApp with session.app
    - Show error toast on failure
    - _Requirements: 3.2, 3.3_

- [x] 5. Implement Post-Restore Launch Prompt
  - [x] 5.1 Create launch prompt modal state and UI
    - Add showLaunchPrompt and launchPromptApp state variables
    - Create modal with "Launch App" and "Close" buttons
    - Handle Escape key and click outside to close
    - _Requirements: 4.1, 4.3, 4.4_
  - [x] 5.2 Show prompt after successful restore operations
    - Modify handleRestore to show prompt on success
    - Modify handleRestoreAccountOnly to show prompt on success
    - Modify handleRestoreAddonOnly to show prompt on success
    - _Requirements: 4.1, 4.2_

- [x] 6. Set Active Session on Addon/Account Restore
  - [x] 6.1 Update RestoreAddonOnly in app.go
    - Call SetActiveSession after successful restore
    - _Requirements: 5.1_
  - [x] 6.2 Update RestoreAccountOnly in app.go
    - Call SetActiveSession after successful restore
    - _Requirements: 5.2_

- [x] 7. Auto-Generate New ID After Restore
  - [x] 7.1 Update RestoreBackup in app.go
    - Call GenerateNewID after successful restore
    - Update progress message with ID count
    - _Requirements: 6.1, 6.4, 6.5_
  - [x] 7.2 Update RestoreAddonOnly in app.go
    - Call GenerateNewID after successful restore
    - Update progress message with ID count
    - _Requirements: 6.2, 6.4, 6.5_
  - [x] 7.3 Update RestoreAccountOnly in app.go
    - Call GenerateNewID after successful restore
    - Update progress message with ID count
    - _Requirements: 6.3, 6.4, 6.5_

- [x] 8. Remove Unused Settings from UI
  - [x] 8.1 Remove Remember Last Tab toggle from SettingsTab
    - Remove the SettingToggle component for rememberLastTab in General section
    - _Requirements: 7.1_
  - [x] 8.2 Remove Debug Mode toggle from SettingsTab
    - Remove the SettingToggle component for debugMode in Experimental section
    - _Requirements: 9.1_
  - [x] 8.3 Remove Skip Data Folder toggle from SettingsTab
    - Remove the SettingToggle component for skipDataFolder in Experimental section
    - _Requirements: 10.1_

- [x] 9. Implement Bulk Session Management
  - [x] 9.1 Add ClearAllSessions function to backend
    - Create new function in app.go that deletes all sessions for all apps
    - _Requirements: 8.4_
  - [x] 9.2 Add BackupAllSessions function to backend
    - Create new function in app.go that creates zip archive of backup folder
    - _Requirements: 8.5_
  - [x] 9.3 Add bulk management buttons to SettingsTab
    - Create grid layout with "Backup All Sessions" and "Clear All Sessions" buttons
    - Add confirmation modal for Clear All Sessions
    - Show toast notifications on completion
    - _Requirements: 8.1, 8.2, 8.3, 8.6_

- [x] 10. Remove Skip Data Folder from Backend
  - [x] 10.1 Update AppConfig struct in loader.go
    - Remove SkipDataFolder field or mark as ignored
    - Ensure backward compatibility with existing configs
    - _Requirements: 10.2, 10.4_
  - [x] 10.2 Remove skip data folder logic from app.go
    - Remove skipDataFolder checks in ResetApp
    - Remove skipDataFolder checks in CreateBackup
    - _Requirements: 10.3_

- [x] 11. Reorganize Reset Tab Layout
  - [x] 11.1 Update app row layout in ResetTab
    - Move app name to left side
    - Add badges for session count, auto-backup count, addon count on right side
    - Load counts for each app in the list
    - _Requirements: 11.1, 11.2_

- [x] 12. Checkpoint - Test All Changes
  - Ensure all tests pass, ask the user if questions arise.
  - Verify alphabetical sorting works
  - Verify selection persistence works across tab navigation
  - Verify restore operations set active session and generate new IDs
  - Verify removed settings are no longer visible

- [x] 13. Update Version and Changelog
  - [x] 13.1 Update version to 2.1.0
    - Update version in frontend/src/lib/version.js
    - Update version in frontend/package.json
    - _Requirements: N/A_
  - [x] 13.2 Create CHANGELOG for v2.1.0
    - Create or update CHANGELOG.md with all changes
    - Include: UX improvements, new features, removed features, bug fixes
    - _Requirements: N/A_

- [x] 14. Final Checkpoint
  - Ensure all tests pass, ask the user if questions arise.
  - Verify all requirements are implemented
  - Verify changelog is complete

## Notes

- Tasks are ordered to minimize dependencies between changes
- Backend changes (Go) should be done before frontend changes that depend on them
- Settings store changes should be done first as other components depend on it
- Property tests can be added as optional sub-tasks if needed
