# Requirements Document

## Introduction

This document outlines UX improvements for SurfManager v2.0, focusing on better user experience in the Reset Tab and Session Tab. The improvements include alphabetical sorting, state persistence, and enhanced restore workflows.

## Glossary

- **Reset_Tab**: The UI tab for resetting application data
- **Session_Tab**: The UI tab for managing backup sessions
- **App_Selector**: The dropdown/list component for selecting applications
- **Active_Session**: The currently active backup session for an application
- **Restore_Operation**: The process of restoring a backup session to an application

## Requirements

### Requirement 1: Alphabetical App Sorting

**User Story:** As a user, I want the apps list to be sorted alphabetically, so that I can quickly find the app I'm looking for.

#### Acceptance Criteria

1. WHEN the Reset_Tab loads the apps list, THE App_Selector SHALL display apps sorted alphabetically by display name
2. WHEN the Session_Tab loads the apps list, THE App_Selector SHALL display apps sorted alphabetically by display name
3. WHEN new apps are added to the configuration, THE App_Selector SHALL maintain alphabetical order

### Requirement 2: Remember Last Selected App

**User Story:** As a user, I want the app I selected to be remembered when I navigate between tabs, so that I don't have to re-select it every time.

#### Acceptance Criteria

1. WHEN a user selects an app in the Reset_Tab, THE System SHALL persist the selection in local storage
2. WHEN a user navigates away from Reset_Tab and returns, THE Reset_Tab SHALL restore the previously selected app
3. WHEN a user selects an app in the Session_Tab, THE System SHALL persist the selection in local storage
4. WHEN a user navigates away from Session_Tab and returns, THE Session_Tab SHALL restore the previously selected app
5. WHEN the persisted app no longer exists in the configuration, THE System SHALL fallback to the first app in the sorted list

### Requirement 3: Launch App from Session Context Menu

**User Story:** As a user, I want to launch an application directly from the session right-click menu, so that I can quickly open the app after managing sessions.

#### Acceptance Criteria

1. WHEN a user right-clicks on a session row, THE Context_Menu SHALL include a "Launch App" option
2. WHEN a user clicks "Launch App" from the context menu, THE System SHALL launch the associated application executable
3. IF the application executable is not found, THEN THE System SHALL display an error toast notification

### Requirement 4: Post-Restore Launch Prompt

**User Story:** As a user, I want to be prompted to launch the application after a successful restore, so that I can quickly continue working.

#### Acceptance Criteria

1. WHEN a restore operation completes successfully, THE System SHALL display a confirmation modal with "Launch App" and "Close" options
2. WHEN the user clicks "Launch App" in the modal, THE System SHALL launch the application and close the modal
3. WHEN the user clicks "Close" in the modal, THE System SHALL close the modal without launching the app
4. WHEN the user presses Escape or clicks outside the modal, THE System SHALL close the modal without launching the app

### Requirement 5: Set Active Session on Addon/Account Restore

**User Story:** As a user, I want the session to be marked as active when I restore only addons or account, so that I can track which session is currently in use.

#### Acceptance Criteria

1. WHEN a user performs "Restore Addon Only" operation successfully, THE System SHALL set that session as the active session
2. WHEN a user performs "Restore Account Only" operation successfully, THE System SHALL set that session as the active session
3. WHEN the active session is set, THE Session_Tab SHALL visually indicate the active session in the list

### Requirement 6: Auto-Generate New ID After Restore

**User Story:** As a user, I want the application's machine ID to be regenerated after any restore operation, so that the restored session has a fresh identity.

#### Acceptance Criteria

1. WHEN a full session restore completes successfully, THE System SHALL automatically generate new machine IDs for the application
2. WHEN an addon-only restore completes successfully, THE System SHALL automatically generate new machine IDs for the application
3. WHEN an account-only restore completes successfully, THE System SHALL automatically generate new machine IDs for the application
4. WHEN new IDs are generated, THE System SHALL update all relevant JSON files (machineId, telemetry.machineId, sessionId)
5. WHEN new IDs are generated, THE System SHALL display the count of updated files in the progress message

### Requirement 7: Remove Unused Remember Last Tab Setting

**User Story:** As a developer, I want to remove the unused "Remember Last Tab" setting, so that the settings UI is cleaner and doesn't show non-functional options.

#### Acceptance Criteria

1. THE Settings_Tab SHALL NOT display the "Remember Last Tab" toggle option
2. THE System SHALL remove the rememberLastTab property from the settings store
3. THE System SHALL remove any related lastActiveTab persistence logic

### Requirement 8: Bulk Session Management in Settings

**User Story:** As a user, I want to backup all sessions or clear all sessions from the Settings tab, so that I can manage my data more efficiently.

#### Acceptance Criteria

1. THE Settings_Tab SHALL display a "Backup All Sessions" button in a grid layout
2. THE Settings_Tab SHALL display a "Clear All Sessions" button in a grid layout
3. WHEN a user clicks "Clear All Sessions", THE System SHALL display a confirmation modal before proceeding
4. WHEN the user confirms the clear operation, THE System SHALL delete all backup sessions for all apps
5. WHEN a user clicks "Backup All Sessions", THE System SHALL create a backup archive of all sessions
6. WHEN bulk operations complete, THE System SHALL display a success or error toast notification

### Requirement 9: Remove Unused Debug Mode Setting

**User Story:** As a developer, I want to remove the unused "Debug Mode" setting from Experimental section, so that the settings UI doesn't show non-functional options.

#### Acceptance Criteria

1. THE Settings_Tab SHALL NOT display the "Debug Mode" toggle in the Experimental section
2. THE System SHALL remove the debugMode property from the settings store
3. THE System SHALL remove any related debug mode logic

### Requirement 10: Remove Skip Data Folder Feature

**User Story:** As a developer, I want to remove the "Skip Data Folder" feature from app configuration, so that the UI is cleaner since "Restore Account Only" and "Restore Addon Only" already provide this functionality.

#### Acceptance Criteria

1. THE Config_Tab SHALL NOT display the "Skip Data Folder" checkbox option
2. THE System SHALL remove the skipDataFolder property from AppConfig
3. THE System SHALL remove any related skip data folder logic in backup/restore operations
4. WHEN loading existing app configs with skipDataFolder, THE System SHALL ignore the property

### Requirement 11: Reorganize Reset Tab App Info Layout

**User Story:** As a user, I want to see session count, auto-backup count, and addon info next to the app name, so that I can quickly see relevant information without scrolling.

#### Acceptance Criteria

1. THE Reset_Tab SHALL display the app name on the left side of the app row
2. THE Reset_Tab SHALL display session count, auto-backup count, and addon count on the right side of the app row
3. THE info badges SHALL be displayed inline after the app name, aligned to the right edge
4. THE layout SHALL be responsive and maintain proper spacing on different screen sizes
