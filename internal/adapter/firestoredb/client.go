package firestoredb

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/quangho/yoga-be/internal/config"
)

func NewClient(ctx context.Context, cfg config.Config) (*firestore.Client, error) {
	if cfg.FirestoreDatabase == "" || cfg.FirestoreDatabase == "(default)" {
		return firestore.NewClient(ctx, cfg.ProjectID)
	}
	return firestore.NewClientWithDatabase(ctx, cfg.ProjectID, cfg.FirestoreDatabase)
}
