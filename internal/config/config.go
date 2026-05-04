package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port              string
	ProjectID         string
	FirestoreDatabase string
}

func Load() (Config, error) {
	cfg := Config{
		Port:              env("PORT", "8080"),
		ProjectID:         firstNonEmpty(os.Getenv("GOOGLE_CLOUD_PROJECT"), os.Getenv("GCP_PROJECT"), os.Getenv("PROJECT_ID")),
		FirestoreDatabase: env("FIRESTORE_DATABASE", "(default)"),
	}
	if cfg.ProjectID == "" {
		return Config{}, fmt.Errorf("GOOGLE_CLOUD_PROJECT, GCP_PROJECT, or PROJECT_ID is required")
	}
	return cfg, nil
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
