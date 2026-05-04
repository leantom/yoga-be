package firestoredb

import "cloud.google.com/go/firestore"

type Registry struct {
	client *firestore.Client
}

func NewRegistry(client *firestore.Client) Registry {
	return Registry{client: client}
}

func (r Registry) Repository(collection string) Repository {
	return Repository{collection: r.client.Collection(collection)}
}
