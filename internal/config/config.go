package config

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/compute/metadata"
)

type Config struct {
	Port              string
	ProjectID         string
	FirestoreDatabase string
}

var metadataProjectID = lookupMetadataProjectID

func Load() (Config, error) {
	projectID, err := resolveProjectID()
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		Port:              env("PORT", "8080"),
		ProjectID:         projectID,
		FirestoreDatabase: env("FIRESTORE_DATABASE", "(default)"),
	}
	return cfg, nil
}

func resolveProjectID() (string, error) {
	if projectID := firstNonEmpty(
		os.Getenv("GOOGLE_CLOUD_PROJECT"),
		os.Getenv("GCP_PROJECT"),
		os.Getenv("PROJECT_ID"),
	); projectID != "" {
		return projectID, nil
	}

	projectID, err := metadataProjectID()
	if err == nil && projectID != "" {
		return projectID, nil
	}
	if err != nil {
		return "", fmt.Errorf("GOOGLE_CLOUD_PROJECT, GCP_PROJECT, or PROJECT_ID is required; metadata lookup failed: %w", err)
	}
	return "", fmt.Errorf("GOOGLE_CLOUD_PROJECT, GCP_PROJECT, or PROJECT_ID is required")
}

func lookupMetadataProjectID() (string, error) {
	client := metadata.NewClient(&http.Client{Timeout: time.Second})
	return client.ProjectID()
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
