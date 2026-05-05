# Yoga Backend

Go backend for a yoga learning app, designed for Google Cloud Run with Firestore as the primary database.

## Architecture

This project uses an Orion-style layered layout:

- `cmd/api`: executable entrypoint and graceful Cloud Run shutdown.
- `internal/transport/rest`: HTTP routing, request parsing, response formatting.
- `internal/domain`: typed ERD models for users, yoga content, subscriptions, payments, progress, favorites, and reviews.
- `internal/adapter/firestoredb`: Firestore client and repository adapter.
- `internal/config`: environment-based runtime configuration.

Firestore collections:

- `users`
- `yoga_categories`
- `yoga_exercises`
- `yoga_programs`
- `program_exercises`
- `user_progress`
- `favorites`
- `subscription_plans`
- `user_subscriptions`
- `payments`
- `reviews`

## Local Run

```bash
cp .env.example .env
export GOOGLE_CLOUD_PROJECT=your-gcp-project-id
gcloud auth application-default login
go run ./cmd/api
```

Health check:

```bash
curl http://localhost:8080/healthz
```

Example create request:

```bash
curl -X POST http://localhost:8080/v1/categories/ \
  -H 'Content-Type: application/json' \
  -d '{"name":"Morning Flow","slug":"morning-flow","order":1,"isActive":true}'
```

Example indexed list request:

```bash
curl 'http://localhost:8080/v1/exercises/?categoryId=cat_1&isActive=true&sort=title:asc'
```

## API

Every collection supports:

- `GET /v1/{resource}/?limit=50`
- `POST /v1/{resource}/`
- `GET /v1/{resource}/{id}`
- `PATCH /v1/{resource}/{id}`
- `DELETE /v1/{resource}/{id}`

Resources:

- `/v1/users`
- `/v1/categories`
- `/v1/exercises`
- `/v1/programs`
- `/v1/program-exercises`
- `/v1/progress`
- `/v1/favorites`
- `/v1/subscription-plans`
- `/v1/subscriptions`
- `/v1/payments`
- `/v1/reviews`

The API stores `_id` inside each Firestore document. If `_id` is not supplied on create, the server generates a UUID.

## Firestore Indexes

Composite indexes are defined in `firestore.indexes.json` for the query patterns used by the API, including category/content filters, program day ordering, user progress, favorites, subscription status, payments, and reviews.

Deploy indexes:

```bash
firebase deploy --only firestore:indexes
```

The included `firestore.rules` denies direct client access. The Cloud Run service should access Firestore through its service account.

## Cloud Run Deployment

### Repository source deploy

Cloud Run repository deploys use Google Buildpacks. This repo keeps the Go
entrypoint under `cmd/api`, so `project.toml` sets `GOOGLE_BUILDABLE=./cmd/api`
for those builds.

Runtime environment:

- `FIRESTORE_DATABASE`, optional, defaults to `(default)`
- `PORT`, supplied by Cloud Run

The Cloud Run service account needs Firestore permissions, usually
`roles/datastore.user`.

### Cloud Build config

Create an Artifact Registry repository once:

```bash
gcloud artifacts repositories create cloud-run \
  --repository-format=docker \
  --location=asia-southeast1
```

Deploy with Cloud Build:

```bash
gcloud builds submit \
  --config cloudbuild.yaml \
  --substitutions _REGION=asia-southeast1,_SERVICE=yoga-api
```

Runtime environment:

- `GOOGLE_CLOUD_PROJECT`, optional on Cloud Run because the service can fall back to the GCP metadata server
- `FIRESTORE_DATABASE`, optional, defaults to `(default)`
- `PORT`, supplied by Cloud Run

The Cloud Run service account needs Firestore permissions, usually
`roles/datastore.user`.

## Verify

```bash
go test ./...
docker build -t yoga-api:local .
```
