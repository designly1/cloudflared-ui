package systemd

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/coreos/go-systemd/v22/sdjournal"
)

// LogEntry represents a single log entry from journald
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
	Priority  string    `json:"priority"`
}

// StreamLogs streams logs from the cloudflared service via journald
func (s *SystemdService) StreamLogs(ctx context.Context, writer io.Writer, follow bool) error {
	journal, err := sdjournal.NewJournal()
	if err != nil {
		return fmt.Errorf("failed to open journal: %w", err)
	}
	defer journal.Close()

	// Add match for cloudflared service
	if err := journal.AddMatch("_SYSTEMD_UNIT=" + serviceName); err != nil {
		return fmt.Errorf("failed to add journal match: %w", err)
	}

	// Seek to the tail (last entries)
	if err := journal.SeekTail(); err != nil {
		return fmt.Errorf("failed to seek to tail: %w", err)
	}

	// Move back a bit to show recent logs
	for i := 0; i < 100; i++ {
		if _, err := journal.Previous(); err != nil {
			break
		}
	}

	// Read logs
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			n, err := journal.Next()
			if err != nil {
				return fmt.Errorf("failed to read next entry: %w", err)
			}

			if n == 0 {
				if !follow {
					return nil
				}
				// Wait for new entries
				journal.Wait(time.Second)
				continue
			}

			entry, err := journal.GetEntry()
			if err != nil {
				continue
			}

			message := entry.Fields[sdjournal.SD_JOURNAL_FIELD_MESSAGE]
			if message != "" {
				timestamp := time.Unix(0, int64(entry.RealtimeTimestamp)*1000)
				logLine := fmt.Sprintf("[%s] %s\n", 
					timestamp.Format("2006-01-02 15:04:05"), 
					message)
				
				if _, err := writer.Write([]byte(logLine)); err != nil {
					return fmt.Errorf("failed to write log: %w", err)
				}
			}
		}
	}
}

// GetRecentLogs retrieves the last N log entries
func (s *SystemdService) GetRecentLogs(count int) ([]LogEntry, error) {
	journal, err := sdjournal.NewJournal()
	if err != nil {
		return nil, fmt.Errorf("failed to open journal: %w", err)
	}
	defer journal.Close()

	// Add match for cloudflared service
	if err := journal.AddMatch("_SYSTEMD_UNIT=" + serviceName); err != nil {
		return nil, fmt.Errorf("failed to add journal match: %w", err)
	}

	// Seek to the tail
	if err := journal.SeekTail(); err != nil {
		return nil, fmt.Errorf("failed to seek to tail: %w", err)
	}

	// Collect entries
	var entries []LogEntry
	for i := 0; i < count; i++ {
		n, err := journal.Previous()
		if err != nil || n == 0 {
			break
		}

		entry, err := journal.GetEntry()
		if err != nil {
			continue
		}

		message := entry.Fields[sdjournal.SD_JOURNAL_FIELD_MESSAGE]
		priority := entry.Fields[sdjournal.SD_JOURNAL_FIELD_PRIORITY]
		
		if message != "" {
			timestamp := time.Unix(0, int64(entry.RealtimeTimestamp)*1000)
			entries = append(entries, LogEntry{
				Timestamp: timestamp,
				Message:   message,
				Priority:  priority,
			})
		}
	}

	// Reverse to get chronological order
	for i, j := 0, len(entries)-1; i < j; i, j = i+1, j-1 {
		entries[i], entries[j] = entries[j], entries[i]
	}

	return entries, nil
}

