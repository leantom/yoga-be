package firestoredb

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrNotFound = errors.New("document not found")

type Repository struct {
	collection *firestore.CollectionRef
}

type ListOptions struct {
	Filters []Filter
	OrderBy []Order
	Limit   int
}

type Filter struct {
	Field string
	Op    string
	Value any
}

type Order struct {
	Field     string
	Direction firestore.Direction
}

func (r Repository) Create(ctx context.Context, data map[string]any) (map[string]any, error) {
	id, _ := data["_id"].(string)
	if id == "" {
		id = uuid.NewString()
		data["_id"] = id
	}
	if _, err := r.collection.Doc(id).Create(ctx, data); err != nil {
		return nil, err
	}
	return data, nil
}

func (r Repository) Get(ctx context.Context, id string) (map[string]any, error) {
	doc, err := r.collection.Doc(id).Get(ctx)
	if err != nil {
		if isNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return doc.Data(), nil
}

func (r Repository) Update(ctx context.Context, id string, patch map[string]any) (map[string]any, error) {
	updates := make([]firestore.Update, 0, len(patch))
	for key, value := range patch {
		if key == "_id" {
			continue
		}
		updates = append(updates, firestore.Update{Path: key, Value: value})
	}
	if len(updates) > 0 {
		if _, err := r.collection.Doc(id).Update(ctx, updates); err != nil {
			if isNotFound(err) {
				return nil, ErrNotFound
			}
			return nil, err
		}
	}
	return r.Get(ctx, id)
}

func (r Repository) Delete(ctx context.Context, id string) error {
	if _, err := r.collection.Doc(id).Delete(ctx, firestore.Exists); err != nil {
		if isNotFound(err) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (r Repository) List(ctx context.Context, opts ListOptions) ([]map[string]any, error) {
	q := r.collection.Query
	for _, filter := range opts.Filters {
		q = q.Where(filter.Field, filter.Op, filter.Value)
	}
	for _, order := range opts.OrderBy {
		q = q.OrderBy(order.Field, order.Direction)
	}
	if opts.Limit > 0 {
		q = q.Limit(opts.Limit)
	}

	iter := q.Documents(ctx)
	defer iter.Stop()

	items := make([]map[string]any, 0)
	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, err
		}
		items = append(items, doc.Data())
	}
	return items, nil
}

func isNotFound(err error) bool {
	return status.Code(err) == codes.NotFound
}

func ParseLimit(raw string, fallback int, max int) (int, error) {
	if raw == "" {
		return fallback, nil
	}
	limit, err := strconv.Atoi(raw)
	if err != nil || limit < 1 {
		return 0, fmt.Errorf("limit must be a positive integer")
	}
	if limit > max {
		return max, nil
	}
	return limit, nil
}

func ParseBool(raw string) (bool, error) {
	switch strings.ToLower(raw) {
	case "true", "1", "yes":
		return true, nil
	case "false", "0", "no":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean value %q", raw)
	}
}
