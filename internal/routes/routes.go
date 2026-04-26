package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/jingxinwangdev/go-prject/internal/app"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	router := chi.NewRouter()

	router.Group(func(r chi.Router) {
		r.Use(app.Middleware.Authenticate)
		r.Get("/workouts/{id}", app.Middleware.RequireAuthenticatedUser(app.WorkoutHandler.HandleGetWorkoutByID))

		r.Post("/workouts", app.Middleware.RequireAuthenticatedUser(app.WorkoutHandler.HandleCreateWorkout))
		r.Put("/workouts/{id}", app.Middleware.RequireAuthenticatedUser(app.WorkoutHandler.HandleUpdateWorkout))
		r.Delete("/workouts/{id}", app.Middleware.RequireAuthenticatedUser(app.WorkoutHandler.HandleDeleteWorkout))
	})
	router.Get("/health", app.HealthCheckHandler)

	router.Post("/users", app.UserHandler.HandleRegisterUser)
	router.Post("/tokens/authentication", app.TokenHandler.HandleCreateToken)
	return router
}
