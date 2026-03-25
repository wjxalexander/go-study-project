package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jingxinwangdev/go-prject/internal/store"
	"github.com/jingxinwangdev/go-prject/internal/utils"
)

type WorkoutHandler struct {
	workoutStore store.WorkoutStore
	logger       *log.Logger
}

// store.WorkoutStore 是一个接口（Interface）。
// Interface: 永远用 InterfaceName（不要指针）。
func NewWorkoutHandler(workoutStore store.WorkoutStore, logger *log.Logger) *WorkoutHandler {
	return &WorkoutHandler{
		workoutStore: workoutStore,
		logger:       logger,
	}
}

// w http.ResponseWriter 是 interface —— 接口内部自带指针语义，无需 *。
// r *http.Request 是 struct —— 需要显式用 * 来传引用，避免拷贝且共享状态。
func (wh *WorkoutHandler) HandleGetWorkoutByID(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIdParam(r, "id")
	if err != nil {
		wh.logger.Printf("Error reading id parameter: %v", err)
		utils.WriteJsonResponse(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	workout, err := wh.workoutStore.GetWorkoutByID(workoutID)
	if err != nil {
		wh.logger.Printf("Error getting workout: %v", err)
		utils.WriteJsonResponse(w, http.StatusInternalServerError, utils.Envelope{"error": "Workout not found"})
		return
	}
	utils.WriteJsonResponse(w, http.StatusOK, utils.Envelope{"data": workout})
}

func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		wh.logger.Printf("Error decoding workout: %v", err)
		utils.WriteJsonResponse(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	createdWorkout, err := wh.workoutStore.CreateWorkout(&workout)
	if err != nil {
		wh.logger.Printf("Error creating workout: %v", err)
		utils.WriteJsonResponse(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to create workout"})
		return
	}
	utils.WriteJsonResponse(w, http.StatusCreated, utils.Envelope{"workout": createdWorkout})
}

func (wh *WorkoutHandler) HandleUpdateWorkout(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIdParam(r, "id")
	if err != nil {
		wh.logger.Printf("Error reading id parameter: %v", err)
		utils.WriteJsonResponse(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	existingWorkout, err := wh.workoutStore.GetWorkoutByID(workoutID)
	if err != nil {
		wh.logger.Printf("Error getting workout: %v", err)
		utils.WriteJsonResponse(w, http.StatusInternalServerError, utils.Envelope{"error": "Workout not found"})
		return
	}
	if existingWorkout == nil {
		wh.logger.Printf("Workout not found")
		utils.WriteJsonResponse(w, http.StatusNotFound, utils.Envelope{"error": "Workout not found"})
		return
	}
	var updatedWorkoutRequest struct {
		Title           *string              `json:"title"`
		Description     *string              `json:"description"`
		DurationMinutes *int                 `json:"duration_minutes"`
		CaloriesBurned  *int                 `json:"calories_burned"`
		Entries         []store.WorkoutEntry `json:"entries"`
	}
	err = json.NewDecoder(r.Body).Decode(&updatedWorkoutRequest)
	if err != nil {
		wh.logger.Printf("Error decoding updated workout request: %v", err)
		utils.WriteJsonResponse(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	if updatedWorkoutRequest.Title != nil {
		existingWorkout.Title = *updatedWorkoutRequest.Title
	}
	if updatedWorkoutRequest.Description != nil {
		existingWorkout.Description = *updatedWorkoutRequest.Description
	}
	if updatedWorkoutRequest.DurationMinutes != nil {
		existingWorkout.DurationMinutes = *updatedWorkoutRequest.DurationMinutes
	}
	if updatedWorkoutRequest.CaloriesBurned != nil {
		existingWorkout.CaloriesBurned = *updatedWorkoutRequest.CaloriesBurned
	}
	if updatedWorkoutRequest.Entries != nil {
		existingWorkout.Entries = updatedWorkoutRequest.Entries
	}
	err = wh.workoutStore.UpdateWorkout(existingWorkout)
	if err != nil {
		wh.logger.Printf("Error updating workout: %v", err)
		utils.WriteJsonResponse(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to update workout"})
		return
	}
	utils.WriteJsonResponse(w, http.StatusOK, utils.Envelope{"workout": existingWorkout})
}

func (wh *WorkoutHandler) HandleDeleteWorkout(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIdParam(r, "id")
	if err != nil {
		wh.logger.Printf("Error reading id parameter: %v", err)
		utils.WriteJsonResponse(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	err = wh.workoutStore.DeleteWorkout(workoutID)
	if err != nil {
		wh.logger.Printf("Error deleting workout: %v", err)
		utils.WriteJsonResponse(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to delete workout"})
		return
	}
	utils.WriteJsonResponse(w, http.StatusNoContent, utils.Envelope{"message": "Workout deleted successfully"})
}
