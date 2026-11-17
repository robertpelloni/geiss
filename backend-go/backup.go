package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type BackupInfo struct {
	Filename  string    `json:"filename"`
	SizeBytes int64     `json:"size_bytes"`
	CreatedAt time.Time `json:"created_at"`
}

type BackupResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Backups []BackupInfo `json:"backups"`
}

// createDatabaseBackup safely snapshots jules.db to the backups/ folder with a timestamped filename.
func createDatabaseBackup(dbFile, backupDir string) (*BackupInfo, error) {
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create backup dir: %w", err)
	}

	src, err := os.Open(dbFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open source db: %w", err)
	}
	defer src.Close()

	timestamp := time.Now().UTC().Format("20060102_150405")
	backupFilename := fmt.Sprintf("jules_backup_%s.db", timestamp)
	destPath := filepath.Join(backupDir, backupFilename)

	dest, err := os.Create(destPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create backup file: %w", err)
	}
	defer dest.Close()

	bytesWritten, err := io.Copy(dest, src)
	if err != nil {
		return nil, fmt.Errorf("failed to copy database: %w", err)
	}

	return &BackupInfo{
		Filename:  backupFilename,
		SizeBytes: bytesWritten,
		CreatedAt: time.Now().UTC(),
	}, nil
}

// createStateSnapshot records the current submodule state as JSON. Why:
// the backend no longer has a database (the gorm/sqlite layer was removed
// as dead code), so a "database backup" had nothing to copy and always
// failed with HTTP 500. The state worth snapshotting now is submodule sync.
func createStateSnapshot(backupDir string) (*BackupInfo, error) {
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create backup dir: %w", err)
	}
	statuses, err := cachedSubmodules()
	if err != nil {
		return nil, fmt.Errorf("failed to read submodule state: %w", err)
	}
	payload, err := json.MarshalIndent(map[string]any{
		"created_at": time.Now().UTC().Format(time.RFC3339),
		"submodules": statuses,
	}, "", "  ")
	if err != nil {
		return nil, err
	}
	filename := fmt.Sprintf("state_snapshot_%s.json", time.Now().UTC().Format("20060102_150405"))
	if err := os.WriteFile(filepath.Join(backupDir, filename), payload, 0644); err != nil {
		return nil, fmt.Errorf("failed to write snapshot: %w", err)
	}
	return &BackupInfo{Filename: filename, SizeBytes: int64(len(payload)), CreatedAt: time.Now().UTC()}, nil
}

// listBackups lists all existing backups in the backup directory, sorted newest first.
func listBackups(backupDir string) ([]BackupInfo, error) {
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(backupDir)
	if err != nil {
		return nil, err
	}

	// Empty slice, not nil: JSON must encode [] (frontend expects an array).
	backups := []BackupInfo{}
	for _, entry := range entries {
		ext := filepath.Ext(entry.Name())
		if entry.IsDir() || (ext != ".db" && ext != ".json") {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		backups = append(backups, BackupInfo{
			Filename:  entry.Name(),
			SizeBytes: info.Size(),
			CreatedAt: info.ModTime().UTC(),
		})
	}

	sort.Slice(backups, func(i, j int) bool {
		return backups[i].CreatedAt.After(backups[j].CreatedAt)
	})

	return backups, nil
}

var (
	backupDBPath  = "jules.db"
	backupDirPath = "backups"
)

// databaseBackupHandler handles GET (list backups) and POST (trigger new backup).
func databaseBackupHandler(w http.ResponseWriter, r *http.Request) {
	dbPath := backupDBPath
	backupDir := backupDirPath

	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		backups, err := listBackups(backupDir)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(BackupResponse{
			Success: true,
			Message: "Backups retrieved successfully.",
			Backups: backups,
		})
		return
	}

	if r.Method == http.MethodPost {
		// Legacy path only when a real jules.db exists; otherwise snapshot state.
		create := createStateSnapshot
		if _, err := os.Stat(dbPath); err == nil {
			create = func(dir string) (*BackupInfo, error) { return createDatabaseBackup(dbPath, dir) }
		}
		info, err := create(backupDir)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(BackupResponse{
				Success: false,
				Message: fmt.Sprintf("Backup failed: %v", err),
			})
			return
		}

		backups, _ := listBackups(backupDir)
		json.NewEncoder(w).Encode(BackupResponse{
			Success: true,
			Message: fmt.Sprintf("Snapshot created: %s (%d bytes)", info.Filename, info.SizeBytes),
			Backups: backups,
		})
		return
	}

	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}
