package main

import (
	"testing"
	"time"
)

func fakeEnv(values map[string]string) func(string) string {
	return func(key string) string { return values[key] }
}

func TestLoadConfigFromEnv(t *testing.T) {
	cfg, err := loadConfig(nil, fakeEnv(map[string]string{
		"GOWATCH_TARGETS":  "https://a.example, https://b.example ,",
		"GOWATCH_INTERVAL": "1m",
		"GOWATCH_WORKERS":  "8",
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(cfg.Targets) != 2 || cfg.Targets[1] != "https://b.example" {
		t.Errorf("Targets = %q", cfg.Targets)
	}

	if cfg.Interval != time.Minute || cfg.Workers != 8 {
		t.Errorf("Interval = %s, Workers = %d", cfg.Interval, cfg.Workers)
	}

	if cfg.Addr != "127.0.0.1:8080" || cfg.Timeout != 5*time.Second {
		t.Errorf("defaults not applied: %+v", cfg)
	}
}

func TestLoadConfigFlagsOverrideEnv(t *testing.T) {
	cfg, err := loadConfig(
		[]string{"-targets", "http://flag.example", "-addr", ":9000"},
		fakeEnv(map[string]string{"GOWATCH_TARGETS": "http://env.example"}),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Targets[0] != "http://flag.example" || cfg.Addr != ":9000" {
		t.Errorf("flags did not override env: %+v", cfg)
	}
}

func TestLoadConfigInvalid(t *testing.T) {
	tests := map[string][]string{
		"no targets":       {},
		"bad scheme":       {"-targets", "ftp://example.com"},
		"no host":          {"-targets", "http://"},
		"zero workers":     {"-targets", "http://a.example", "-workers", "0"},
		"zero interval":    {"-targets", "http://a.example", "-interval", "0s"},
		"negative timeout": {"-targets", "http://a.example", "-timeout", "-1s"},
	}

	for name, args := range tests {
		t.Run(name, func(t *testing.T) {
			if _, err := loadConfig(args, fakeEnv(nil)); err == nil {
				t.Error("expected an error")
			}
		})
	}
}
