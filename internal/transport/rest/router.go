package rest

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/quangho/yoga-be/internal/adapter/firestoredb"
)

func NewRouter(registry firestoredb.Registry) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{"status": "ok"})
	})

	r.Route("/v1", func(r chi.Router) {
		mountResource(r, registry, resourceConfig{
			Path:        "/users",
			Collection:  "users",
			Filters:     map[string]string{"email": "==", "phone": "==", "authProvider": "==", "status": "=="},
			Orders:      []string{"createdAt", "updatedAt", "fullName", "email"},
			DefaultSort: []string{"createdAt:desc"},
			CreateNow:   []string{"createdAt", "updatedAt"},
			UpdateNow:   []string{"updatedAt"},
		})
		mountResource(r, registry, resourceConfig{
			Path:        "/categories",
			Collection:  "yoga_categories",
			Filters:     map[string]string{"slug": "==", "isActive": "=="},
			Orders:      []string{"order", "name"},
			DefaultSort: []string{"order:asc"},
		})
		mountResource(r, registry, resourceConfig{
			Path:        "/exercises",
			Collection:  "yoga_exercises",
			Filters:     map[string]string{"categoryId": "==", "slug": "==", "level": "==", "isPremium": "==", "isActive": "==", "bodyParts": "array-contains", "benefits": "array-contains"},
			Orders:      []string{"title", "durationSeconds"},
			DefaultSort: []string{"title:asc"},
		})
		mountResource(r, registry, resourceConfig{
			Path:        "/programs",
			Collection:  "yoga_programs",
			Filters:     map[string]string{"categoryId": "==", "slug": "==", "level": "==", "isPremium": "==", "isActive": "=="},
			Orders:      []string{"title", "totalDays", "estimatedMinutes"},
			DefaultSort: []string{"title:asc"},
		})
		mountResource(r, registry, resourceConfig{
			Path:        "/program-exercises",
			Collection:  "program_exercises",
			Filters:     map[string]string{"programId": "==", "exerciseId": "==", "dayNumber": "=="},
			Orders:      []string{"dayNumber", "order"},
			DefaultSort: []string{"dayNumber:asc", "order:asc"},
		})
		mountResource(r, registry, resourceConfig{
			Path:        "/progress",
			Collection:  "user_progress",
			Filters:     map[string]string{"userId": "==", "exerciseId": "==", "programId": "==", "isCompleted": "=="},
			Orders:      []string{"lastWatchedAt", "completedAt", "progressPercent"},
			DefaultSort: []string{"lastWatchedAt:desc"},
			CreateNow:   []string{"lastWatchedAt"},
			UpdateNow:   []string{"lastWatchedAt"},
		})
		mountResource(r, registry, resourceConfig{
			Path:        "/favorites",
			Collection:  "favorites",
			Filters:     map[string]string{"userId": "==", "exerciseId": "=="},
			Orders:      []string{"createdAt"},
			DefaultSort: []string{"createdAt:desc"},
			CreateNow:   []string{"createdAt"},
		})
		mountResource(r, registry, resourceConfig{
			Path:        "/subscription-plans",
			Collection:  "subscription_plans",
			Filters:     map[string]string{"code": "==", "currency": "==", "isActive": "=="},
			Orders:      []string{"price", "durationDays", "name"},
			DefaultSort: []string{"price:asc"},
		})
		mountResource(r, registry, resourceConfig{
			Path:        "/subscriptions",
			Collection:  "user_subscriptions",
			Filters:     map[string]string{"userId": "==", "planId": "==", "status": "==", "autoRenew": "=="},
			Orders:      []string{"startDate", "endDate"},
			DefaultSort: []string{"endDate:desc"},
		})
		mountResource(r, registry, resourceConfig{
			Path:        "/payments",
			Collection:  "payments",
			Filters:     map[string]string{"userId": "==", "subscriptionId": "==", "planId": "==", "method": "==", "status": "==", "transactionCode": "=="},
			Orders:      []string{"paidAt", "amount"},
			DefaultSort: []string{"paidAt:desc"},
		})
		mountResource(r, registry, resourceConfig{
			Path:        "/reviews",
			Collection:  "reviews",
			Filters:     map[string]string{"userId": "==", "exerciseId": "==", "rating": "=="},
			Orders:      []string{"createdAt", "rating"},
			DefaultSort: []string{"createdAt:desc"},
			CreateNow:   []string{"createdAt"},
		})
	})

	return r
}
