package history

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

const maxEntries = 500

type Entry struct {
	Timestamp time.Time         `json:"timestamp"`
	Method    string            `json:"method"`
	URL       string            `json:"url"`
	Headers   map[string]string `json:"headers,omitempty"`
	Body      string            `json:"body,omitempty"`
}

type Store struct {
	Entries []Entry `json:"entries"`
	path    string
}

func New() *Store {
	path, err := storePath()
	if err != nil {
		return &Store{}
	}
	s := &Store{path: path}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return s
	}
	if err != nil {
		return s
	}
	_ = json.Unmarshal(data, s)
	return s
}

func (s *Store) Add(e Entry) {
	s.Entries = append([]Entry{e}, s.Entries...)
	if len(s.Entries) > maxEntries {
		s.Entries = s.Entries[:maxEntries]
	}
	_ = s.save()
}

func (s *Store) Len() int { return len(s.Entries) }

func (s *Store) Get(i int) Entry { return s.Entries[i] }

func (s *Store) save() error {
	if s.path == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o600)
}

func storePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".curly", "history.json"), nil
}
