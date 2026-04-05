package store

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres_test port=5434 sslmode=disable")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	err = Migrate(db, "../../migrations/")
	if err != nil {
		t.Fatalf("failed to migrate test database: %v", err)
	}
	// 清空所有表数据，保证每个测试从干净的状态开始，CASCADE 同时清除外键关联的数据
	_, err = db.Exec("TRUNCATE TABLE workouts, workout_entries CASCADE")
	if err != nil {
		t.Fatalf("failed to truncate test database: %v", err)
	}
	return db
}

func TestCreateWorkout(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewPostgresWorkoutStore(db)
	tests := []struct {
		name    string
		workout *Workout
		wantErr bool
	}{
		{
			name: "valid workout",
			workout: &Workout{
				Title:           "Test Workout",
				Description:     "Test Description",
				DurationMinutes: 30,
				CaloriesBurned:  300,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "Test Exercise",
						Sets:         3,
						Reps:         IntPtr(10),
						Weight:       Float64Ptr(0),
						Notes:        "Test Notes",
						OrderIndex:   1,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "workout with invalid entries",
			workout: &Workout{
				Title:           "full body",
				Description:     "complete workout",
				DurationMinutes: 90,
				CaloriesBurned:  500,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "Plank",
						Sets:         3,
						Reps:         IntPtr(60),
						Notes:        "keep form",
						OrderIndex:   1,
					},
					{
						ExerciseName:    "squats",
						Sets:            4,
						Reps:            IntPtr(12),
						DurationSeconds: IntPtr(60),
						Weight:          Float64Ptr(185.0),
						Notes:           "full depth",
						OrderIndex:      2,
					},
				},
			},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			createdWorkout, err := store.CreateWorkout(test.workout)
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.workout.Title, createdWorkout.Title)
			assert.Equal(t, test.workout.Description, createdWorkout.Description)
			assert.Equal(t, test.workout.DurationMinutes, createdWorkout.DurationMinutes)
			assert.Equal(t, test.workout.CaloriesBurned, createdWorkout.CaloriesBurned)
			assert.Equal(t, len(test.workout.Entries), len(createdWorkout.Entries))
			retrived, err := store.GetWorkoutByID(int64(createdWorkout.ID))
			require.NoError(t, err)
			assert.Equal(t, createdWorkout.ID, retrived.ID)
			assert.Equal(t, len(test.workout.Entries), len(retrived.Entries))
			for i, entry := range test.workout.Entries {
				assert.Equal(t, entry.ExerciseName, retrived.Entries[i].ExerciseName)
				assert.Equal(t, entry.Sets, retrived.Entries[i].Sets)
				assert.Equal(t, entry.Reps, retrived.Entries[i].Reps)
				assert.Equal(t, entry.DurationSeconds, retrived.Entries[i].DurationSeconds)
				assert.Equal(t, entry.Weight, retrived.Entries[i].Weight)
				assert.Equal(t, entry.Notes, retrived.Entries[i].Notes)
				assert.Equal(t, entry.OrderIndex, retrived.Entries[i].OrderIndex)
			}
		})
	}
}

func IntPtr(i int) *int {
	return &i
}

func Float64Ptr(f float64) *float64 {
	return &f
}
