package domain

import "time"

type User struct {
	ID           string    `json:"_id" firestore:"_id"`
	FullName     string    `json:"fullName" firestore:"fullName"`
	Email        string    `json:"email" firestore:"email"`
	Phone        string    `json:"phone,omitempty" firestore:"phone,omitempty"`
	AvatarURL    string    `json:"avatarUrl,omitempty" firestore:"avatarUrl,omitempty"`
	AuthProvider string    `json:"authProvider" firestore:"authProvider"`
	Status       string    `json:"status" firestore:"status"`
	CreatedAt    time.Time `json:"createdAt" firestore:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt" firestore:"updatedAt"`
}

type YogaCategory struct {
	ID          string `json:"_id" firestore:"_id"`
	Name        string `json:"name" firestore:"name"`
	Slug        string `json:"slug" firestore:"slug"`
	Description string `json:"description,omitempty" firestore:"description,omitempty"`
	ImageURL    string `json:"imageUrl,omitempty" firestore:"imageUrl,omitempty"`
	Order       int    `json:"order" firestore:"order"`
	IsActive    bool   `json:"isActive" firestore:"isActive"`
}

type YogaExercise struct {
	ID              string   `json:"_id" firestore:"_id"`
	CategoryID      string   `json:"categoryId" firestore:"categoryId"`
	Title           string   `json:"title" firestore:"title"`
	Slug            string   `json:"slug" firestore:"slug"`
	Description     string   `json:"description,omitempty" firestore:"description,omitempty"`
	VideoURL        string   `json:"videoUrl,omitempty" firestore:"videoUrl,omitempty"`
	ThumbnailURL    string   `json:"thumbnailUrl,omitempty" firestore:"thumbnailUrl,omitempty"`
	Level           string   `json:"level" firestore:"level"`
	DurationSeconds int      `json:"durationSeconds" firestore:"durationSeconds"`
	BodyParts       []string `json:"bodyParts,omitempty" firestore:"bodyParts,omitempty"`
	Benefits        []string `json:"benefits,omitempty" firestore:"benefits,omitempty"`
	IsPremium       bool     `json:"isPremium" firestore:"isPremium"`
	IsActive        bool     `json:"isActive" firestore:"isActive"`
}

type YogaProgram struct {
	ID               string `json:"_id" firestore:"_id"`
	CategoryID       string `json:"categoryId" firestore:"categoryId"`
	Title            string `json:"title" firestore:"title"`
	Slug             string `json:"slug" firestore:"slug"`
	Description      string `json:"description,omitempty" firestore:"description,omitempty"`
	Level            string `json:"level" firestore:"level"`
	TotalDays        int    `json:"totalDays" firestore:"totalDays"`
	EstimatedMinutes int    `json:"estimatedMinutes" firestore:"estimatedMinutes"`
	ThumbnailURL     string `json:"thumbnailUrl,omitempty" firestore:"thumbnailUrl,omitempty"`
	IsPremium        bool   `json:"isPremium" firestore:"isPremium"`
	IsActive         bool   `json:"isActive" firestore:"isActive"`
}

type ProgramExercise struct {
	ID              string `json:"_id" firestore:"_id"`
	ProgramID       string `json:"programId" firestore:"programId"`
	ExerciseID      string `json:"exerciseId" firestore:"exerciseId"`
	DayNumber       int    `json:"dayNumber" firestore:"dayNumber"`
	Order           int    `json:"order" firestore:"order"`
	DurationSeconds int    `json:"durationSeconds" firestore:"durationSeconds"`
}

type UserProgress struct {
	ID              string     `json:"_id" firestore:"_id"`
	UserID          string     `json:"userId" firestore:"userId"`
	ExerciseID      string     `json:"exerciseId" firestore:"exerciseId"`
	ProgramID       string     `json:"programId,omitempty" firestore:"programId,omitempty"`
	ProgressPercent int        `json:"progressPercent" firestore:"progressPercent"`
	WatchedSeconds  int        `json:"watchedSeconds" firestore:"watchedSeconds"`
	IsCompleted     bool       `json:"isCompleted" firestore:"isCompleted"`
	CompletedAt     *time.Time `json:"completedAt,omitempty" firestore:"completedAt,omitempty"`
	LastWatchedAt   time.Time  `json:"lastWatchedAt" firestore:"lastWatchedAt"`
}

type Favorite struct {
	ID         string    `json:"_id" firestore:"_id"`
	UserID     string    `json:"userId" firestore:"userId"`
	ExerciseID string    `json:"exerciseId" firestore:"exerciseId"`
	CreatedAt  time.Time `json:"createdAt" firestore:"createdAt"`
}

type SubscriptionPlan struct {
	ID           string   `json:"_id" firestore:"_id"`
	Name         string   `json:"name" firestore:"name"`
	Code         string   `json:"code" firestore:"code"`
	Price        float64  `json:"price" firestore:"price"`
	Currency     string   `json:"currency" firestore:"currency"`
	DurationDays int      `json:"durationDays" firestore:"durationDays"`
	Features     []string `json:"features,omitempty" firestore:"features,omitempty"`
	IsActive     bool     `json:"isActive" firestore:"isActive"`
}

type UserSubscription struct {
	ID        string    `json:"_id" firestore:"_id"`
	UserID    string    `json:"userId" firestore:"userId"`
	PlanID    string    `json:"planId" firestore:"planId"`
	StartDate time.Time `json:"startDate" firestore:"startDate"`
	EndDate   time.Time `json:"endDate" firestore:"endDate"`
	Status    string    `json:"status" firestore:"status"`
	AutoRenew bool      `json:"autoRenew" firestore:"autoRenew"`
}

type Payment struct {
	ID              string     `json:"_id" firestore:"_id"`
	UserID          string     `json:"userId" firestore:"userId"`
	SubscriptionID  string     `json:"subscriptionId" firestore:"subscriptionId"`
	PlanID          string     `json:"planId" firestore:"planId"`
	Amount          float64    `json:"amount" firestore:"amount"`
	Currency        string     `json:"currency" firestore:"currency"`
	Method          string     `json:"method" firestore:"method"`
	Status          string     `json:"status" firestore:"status"`
	TransactionCode string     `json:"transactionCode,omitempty" firestore:"transactionCode,omitempty"`
	PaidAt          *time.Time `json:"paidAt,omitempty" firestore:"paidAt,omitempty"`
}

type Review struct {
	ID         string    `json:"_id" firestore:"_id"`
	UserID     string    `json:"userId" firestore:"userId"`
	ExerciseID string    `json:"exerciseId" firestore:"exerciseId"`
	Rating     int       `json:"rating" firestore:"rating"`
	Comment    string    `json:"comment,omitempty" firestore:"comment,omitempty"`
	CreatedAt  time.Time `json:"createdAt" firestore:"createdAt"`
}
