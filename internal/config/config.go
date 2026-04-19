package config

import (
	"bufio"
	"os"
	"strings"
)

const Path = "/etc/togram/config"

type Config struct {
	Token  string
	ChatID string
}

func Load() (*Config, error) {
	f, err := os.Open(Path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, err
	}
	defer f.Close()

	cfg := &Config{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		switch strings.TrimSpace(k) {
		case "token":
			cfg.Token = strings.TrimSpace(v)
		case "chat":
			cfg.ChatID = strings.TrimSpace(v)
		}
	}
	return cfg, scanner.Err()
}
