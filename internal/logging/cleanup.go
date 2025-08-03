package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// CleanupConfig controls log cleanup behavior
type CleanupConfig struct {
	MaxEpisodeFiles int    // Maximum episode log files to keep
	LogDir          string // Directory containing logs
	DryRun          bool   // If true, only print what would be deleted
}

// DefaultCleanupConfig returns sensible defaults
func DefaultCleanupConfig() CleanupConfig {
	return CleanupConfig{
		MaxEpisodeFiles: 50, // Keep only last 50 episodes
		LogDir:          "logs",
		DryRun:          false,
	}
}

// CleanupOldLogs removes old episode log files, keeping only the most recent ones
func CleanupOldLogs(config CleanupConfig) error {
	// Find all episode log files
	pattern := filepath.Join(config.LogDir, "episode_*.json")
	episodeFiles, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("failed to find episode files: %v", err)
	}

	// Also check for events files
	eventsPattern := filepath.Join(config.LogDir, "events_*.json")
	eventsFiles, err := filepath.Glob(eventsPattern)
	if err != nil {
		return fmt.Errorf("failed to find events files: %v", err)
	}

	allFiles := append(episodeFiles, eventsFiles...)
	
	if len(allFiles) <= config.MaxEpisodeFiles {
		fmt.Printf("Found %d log files, within limit of %d. No cleanup needed.\n", 
			len(allFiles), config.MaxEpisodeFiles)
		return nil
	}

	// Sort files by modification time (oldest first)
	sort.Slice(allFiles, func(i, j int) bool {
		infoI, errI := os.Stat(allFiles[i])
		infoJ, errJ := os.Stat(allFiles[j])
		
		if errI != nil || errJ != nil {
			// Fallback to filename sort if stat fails
			return allFiles[i] < allFiles[j]
		}
		
		return infoI.ModTime().Before(infoJ.ModTime())
	})

	// Calculate how many files to delete
	filesToDelete := len(allFiles) - config.MaxEpisodeFiles
	if filesToDelete <= 0 {
		return nil
	}

	fmt.Printf("Found %d log files, keeping %d newest, deleting %d oldest\n", 
		len(allFiles), config.MaxEpisodeFiles, filesToDelete)

	// Delete oldest files
	var deletedCount int
	var totalSizeDeleted int64
	
	for i := 0; i < filesToDelete; i++ {
		file := allFiles[i]
		
		// Get file size before deletion
		if info, err := os.Stat(file); err == nil {
			totalSizeDeleted += info.Size()
		}
		
		if config.DryRun {
			fmt.Printf("Would delete: %s\n", file)
		} else {
			if err := os.Remove(file); err != nil {
				fmt.Printf("Warning: failed to delete %s: %v\n", file, err)
			} else {
				deletedCount++
			}
		}
	}

	if config.DryRun {
		fmt.Printf("Dry run complete. Would delete %d files (~%.1f MB)\n", 
			filesToDelete, float64(totalSizeDeleted)/(1024*1024))
	} else {
		fmt.Printf("Cleanup complete. Deleted %d files (%.1f MB freed)\n", 
			deletedCount, float64(totalSizeDeleted)/(1024*1024))
	}

	return nil
}

// GetLogDirSize returns the total size of the log directory
func GetLogDirSize(logDir string) (int64, int, error) {
	var totalSize int64
	var fileCount int

	err := filepath.Walk(logDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".json") {
			totalSize += info.Size()
			fileCount++
		}
		
		return nil
	})

	return totalSize, fileCount, err
}

// AutoCleanup performs automatic cleanup if log directory exceeds thresholds
func AutoCleanup(logDir string, maxSizeMB float64, maxFiles int) error {
	size, count, err := GetLogDirSize(logDir)
	if err != nil {
		return fmt.Errorf("failed to get log directory size: %v", err)
	}

	sizeMB := float64(size) / (1024 * 1024)
	
	fmt.Printf("Log directory: %.1f MB, %d files\n", sizeMB, count)

	// Check if cleanup is needed
	needsCleanup := sizeMB > maxSizeMB || count > maxFiles
	
	if !needsCleanup {
		fmt.Printf("Log directory within limits (%.1f/%.1f MB, %d/%d files)\n", 
			sizeMB, maxSizeMB, count, maxFiles)
		return nil
	}

	fmt.Printf("Log directory exceeds limits, performing cleanup...\n")
	
	config := CleanupConfig{
		MaxEpisodeFiles: maxFiles / 2, // Keep half the limit
		LogDir:          logDir,
		DryRun:          false,
	}

	return CleanupOldLogs(config)
}