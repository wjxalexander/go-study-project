package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/jingxinwangdev/go-prject/internal/app"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	router := chi.NewRouter()
	router.Get("/health", app.HealthCheckHandler)
	router.Get("/workouts/{id}", app.WorkoutHandler.HandleGetWorkoutByID)
	router.Post("/workouts", app.WorkoutHandler.HandleCreateWorkout)
	return router
}
