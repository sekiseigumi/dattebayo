package logger

import (
	"fmt"
	"sync"
	"time"
)

type LogEntry struct {
	Timestamp time.Time
	Source    string
	Message   string
}

type Logger struct {
	entries []LogEntry
	mu      sync.RWMutex
}

func NewLogger() *Logger {
	return &Logger{
		entries: make([]LogEntry, 0),
	}
}

func (l *Logger) Log(source, message string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	entry := LogEntry{
		Timestamp: time.Now(),
		Source:    source,
		Message:   message,
	}

	l.entries = append(l.entries, entry)
}

func (l *Logger) GetEntries() []LogEntry {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return append([]LogEntry(nil), l.entries...)
}

func (l *Logger) FormatEntry(entry LogEntry) string {
	return fmt.Sprintf("[%s] %s: %s", entry.Timestamp.Format("15:04:05"), entry.Source, entry.Message)
}
