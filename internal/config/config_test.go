package config

import (
	"errors"
	"testing"
)

func TestLoadUsesEnvironmentProjectID(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("FIRESTORE_DATABASE", "custom-db")
	t.Setenv("GOOGLE_CLOUD_PROJECT", "env-project")

	original := metadataProjectID
	defer func() { metadataProjectID = original }()

	metadataProjectID = func() (string, error) {
		t.Fatal("metadata lookup should not run when env project id is set")
		return "", nil
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.ProjectID != "env-project" {
		t.Fatalf("ProjectID = %q, want %q", cfg.ProjectID, "env-project")
	}
	if cfg.Port != "9090" {
		t.Fatalf("Port = %q, want %q", cfg.Port, "9090")
	}
	if cfg.FirestoreDatabase != "custom-db" {
		t.Fatalf("FirestoreDatabase = %q, want %q", cfg.FirestoreDatabase, "custom-db")
	}
}

func TestLoadFallsBackToMetadataProjectID(t *testing.T) {
	t.Setenv("GOOGLE_CLOUD_PROJECT", "")
	t.Setenv("GCP_PROJECT", "")
	t.Setenv("PROJECT_ID", "")

	original := metadataProjectID
	defer func() { metadataProjectID = original }()

	metadataProjectID = func() (string, error) {
		return "metadata-project", nil
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.ProjectID != "metadata-project" {
		t.Fatalf("ProjectID = %q, want %q", cfg.ProjectID, "metadata-project")
	}
	if cfg.Port != "8080" {
		t.Fatalf("Port = %q, want %q", cfg.Port, "8080")
	}
	if cfg.FirestoreDatabase != "(default)" {
		t.Fatalf("FirestoreDatabase = %q, want %q", cfg.FirestoreDatabase, "(default)")
	}
}

func TestLoadReturnsErrorWhenProjectIDUnavailable(t *testing.T) {
	t.Setenv("GOOGLE_CLOUD_PROJECT", "")
	t.Setenv("GCP_PROJECT", "")
	t.Setenv("PROJECT_ID", "")

	original := metadataProjectID
	defer func() { metadataProjectID = original }()

	metadataProjectID = func() (string, error) {
		return "", errors.New("metadata unavailable")
	}

	_, err := Load()
	if err == nil {
		t.Fatal("Load() error = nil, want non-nil")
	}
}
