package main

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Addr     string
	Targets  []string
	Workers  int
	Interval time.Duration
	Timeout  time.Duration
}

// loadConfig reads configuration from flags, falling back to GOWATCH_*
// environment variables and then to defaults.
func loadConfig(args []string, getenv func(string) string) (Config, error) {
	fs := flag.NewFlagSet("gowatch", flag.ContinueOnError)

	addr := fs.String("addr", envOr(getenv, "GOWATCH_ADDR", "127.0.0.1:8080"),
		"address the HTTP server listens on")
	targets := fs.String("targets", getenv("GOWATCH_TARGETS"),
		"comma-separated list of URLs to monitor")
	workers := fs.Int("workers", envInt(getenv, "GOWATCH_WORKERS", 4),
		"number of concurrent checks")
	interval := fs.Duration("interval", envDuration(getenv, "GOWATCH_INTERVAL", 30*time.Second),
		"time between check rounds")
	timeout := fs.Duration("timeout", envDuration(getenv, "GOWATCH_TIMEOUT", 5*time.Second),
		"timeout for each check")

	if err := fs.Parse(args); err != nil {
		return Config{}, err
	}

	cfg := Config{
		Addr:     *addr,
		Targets:  splitTargets(*targets),
		Workers:  *workers,
		Interval: *interval,
		Timeout:  *timeout,
	}

	return cfg, cfg.validate()
}

func (c Config) validate() error {
	var errs []error

	if len(c.Targets) == 0 {
		errs = append(errs, errors.New("at least one target is required (-targets or GOWATCH_TARGETS)"))
	}

	for _, target := range c.Targets {
		u, err := url.Parse(target)
		if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
			errs = append(errs, fmt.Errorf("invalid target %q: must be an http(s) URL", target))
		}
	}

	if c.Workers < 1 {
		errs = append(errs, errors.New("workers must be at least 1"))
	}

	if c.Interval <= 0 {
		errs = append(errs, errors.New("interval must be positive"))
	}

	if c.Timeout <= 0 {
		errs = append(errs, errors.New("timeout must be positive"))
	}

	return errors.Join(errs...)
}

func splitTargets(raw string) []string {
	var targets []string

	for _, part := range strings.Split(raw, ",") {
		if part = strings.TrimSpace(part); part != "" {
			targets = append(targets, part)
		}
	}

	return targets
}

func envOr(getenv func(string) string, key, fallback string) string {
	if value := getenv(key); value != "" {
		return value
	}

	return fallback
}

func envInt(getenv func(string) string, key string, fallback int) int {
	if value, err := strconv.Atoi(getenv(key)); err == nil {
		return value
	}

	return fallback
}

func envDuration(getenv func(string) string, key string, fallback time.Duration) time.Duration {
	if value, err := time.ParseDuration(getenv(key)); err == nil {
		return value
	}

	return fallback
}
