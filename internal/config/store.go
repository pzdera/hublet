package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Store struct {
	mu         sync.RWMutex
	dataDir    string
	configPath string
	config     Config
}

func NewStore(dataDir string) (*Store, error) {
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return nil, fmt.Errorf("create data directory: %w", err)
	}

	store := &Store{
		dataDir:    dataDir,
		configPath: filepath.Join(dataDir, "config.json"),
	}

	if err := store.load(); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *Store) Get() Config {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return clone(s.config)
}

func (s *Store) Save(cfg Config) error {
	if err := Validate(cfg); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.backupCurrent(); err != nil {
		return err
	}

	if err := writeAtomic(s.configPath, cfg); err != nil {
		return err
	}

	s.config = clone(cfg)

	return s.pruneBackups(20)
}

func (s *Store) load() error {
	raw, err := os.ReadFile(s.configPath)

	if os.IsNotExist(err) {
		cfg := Default()

		if err := writeAtomic(s.configPath, cfg); err != nil {
			return err
		}

		s.config = cfg
		return nil
	}

	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	var cfg Config

	if err := json.Unmarshal(raw, &cfg); err != nil {
		return fmt.Errorf("decode config: %w", err)
	}

	if err := Validate(cfg); err != nil {
		return fmt.Errorf("validate config: %w", err)
	}

	s.config = cfg

	return nil
}

func (s *Store) backupCurrent() error {
	raw, err := os.ReadFile(s.configPath)

	if os.IsNotExist(err) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("read config for backup: %w", err)
	}

	backupDir := filepath.Join(s.dataDir, "backups")

	if err := os.MkdirAll(backupDir, 0o755); err != nil {
		return fmt.Errorf("create backup directory: %w", err)
	}

	filename := "config-" + time.Now().UTC().Format("20060102T150405.000000000Z") + ".json"
	path := filepath.Join(backupDir, filename)

	if err := os.WriteFile(path, raw, 0o644); err != nil {
		return fmt.Errorf("write config backup: %w", err)
	}

	return nil
}

func (s *Store) pruneBackups(keep int) error {
	backupDir := filepath.Join(s.dataDir, "backups")

	entries, err := os.ReadDir(backupDir)

	if os.IsNotExist(err) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("read backup directory: %w", err)
	}

	if len(entries) <= keep {
		return nil
	}

	for _, entry := range entries[:len(entries)-keep] {
		if err := os.Remove(filepath.Join(backupDir, entry.Name())); err != nil {
			return fmt.Errorf("remove old backup: %w", err)
		}
	}

	return nil
}

func writeAtomic(path string, cfg Config) error {
	raw, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("encode config: %w", err)
	}

	raw = append(raw, '\n')

	temp, err := os.CreateTemp(filepath.Dir(path), ".config-*.tmp")
	if err != nil {
		return fmt.Errorf("create temporary config: %w", err)
	}

	tempPath := temp.Name()
	defer os.Remove(tempPath)

	if _, err := temp.Write(raw); err != nil {
		temp.Close()
		return fmt.Errorf("write temporary config: %w", err)
	}

	if err := temp.Sync(); err != nil {
		temp.Close()
		return fmt.Errorf("sync temporary config: %w", err)
	}

	if err := temp.Close(); err != nil {
		return fmt.Errorf("close temporary config: %w", err)
	}

	if err := os.Chmod(tempPath, 0o644); err != nil {
		return fmt.Errorf("set config permissions: %w", err)
	}

	if err := os.Rename(tempPath, path); err != nil {
		return fmt.Errorf("replace config: %w", err)
	}

	return nil
}

func clone(cfg Config) Config {
	raw, _ := json.Marshal(cfg)

	var result Config
	_ = json.Unmarshal(raw, &result)

	return result
}
