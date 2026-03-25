package store

import "database/sql"

type Workout struct {
	ID              int64          `json:"id"`
	Title           string         `json:"title"`
	Description     string         `json:"description"`
	DurationMinutes int            `json:"duration_minutes"`
	CaloriesBurned  int            `json:"calories_burned"`
	Entries         []WorkoutEntry `json:"entries"`
}

type WorkoutEntry struct {
	ID              int64    `json:"id"`
	ExerciseName    string   `json:"exercise_name"`
	Sets            int      `json:"sets"`
	Reps            *int     `json:"reps"` //在 Web API 和数据库映射中，当一个字段是可选的（Optional）或者允许数据库存为 NULL 时
	DurationSeconds *int     `json:"duration_seconds"`
	Weight          *float64 `json:"weight"`
	Notes           *string  `json:"notes"`
	OrderIndex      int      `json:"order_index"`
}

type PostgresWorkoutStore struct {
	db *sql.DB
}

func NewPostgresWorkoutStore(db *sql.DB) *PostgresWorkoutStore {
	return &PostgresWorkoutStore{db: db}
}

// interface: define a collection of methods signatures
// purpose: decouple the specific database implementation from the business logic
type WorkoutStore interface {
	CreateWorkout(workout *Workout) (*Workout, error)
	GetWorkoutByID(id int64) (*Workout, error)
	UpdateWorkout(workout *Workout) error
	DeleteWorkout(id int64) error
}

func (pg *PostgresWorkoutStore) CreateWorkout(workout *Workout) (*Workout, error) {
	// transactions
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}
	/*
	* defer 关键字会将其后的函数调用推迟到当前函数执行结束前（即 return 之前）才执行。
	* 无论函数是正常返回，还是因为发生 panic 而异常退出，defer 的代码都会被执行。
	* 情况 A：发生错误（Early Return）
	* 结果：函数返回前，defer 触发 tx.Rollback()。事务被回滚，数据库撤销所有未提交的更改，保证数据安全。
	* 情况 B：发生 Panic
	* 结果：defer 仍然会执行 tx.Rollback()，释放数据库连接，避免连接一直被占用。
	 */
	defer tx.Rollback()
	query := `
		INSERT INTO workouts (title, description, duration_minutes, calories_burned)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	// 插入 workout 并且返回 id
	err = tx.QueryRow(query, workout.Title, workout.Description, workout.DurationMinutes, workout.CaloriesBurned).Scan(&workout.ID)
	if err != nil {
		return nil, err
	}
	for _, entry := range workout.Entries {
		query = `
			INSERT INTO workout_entries (workout_id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id
		`
		err = tx.QueryRow(query, workout.ID, entry.ExerciseName, entry.Sets, entry.Reps, entry.DurationSeconds, entry.Weight, entry.Notes, entry.OrderIndex).Scan(&entry.ID)
		if err != nil {
			return nil, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return workout, nil
}

func (pg *PostgresWorkoutStore) GetWorkoutByID(id int64) (*Workout, error) {
	query := `
		SELECT id, title, description, duration_minutes, calories_burned
		FROM workouts
		WHERE id = $1
	`
	workout := &Workout{}
	err := pg.db.QueryRow(query, id).Scan(&workout.ID, &workout.Title, &workout.Description, &workout.DurationMinutes, &workout.CaloriesBurned)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	entryQuery := `
	SELECT id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index
	FROM workout_entries
	WHERE workout_id = $1
	ORDER BY order_index
  `

	rows, err := pg.db.Query(entryQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var entry WorkoutEntry
		err = rows.Scan(
			&entry.ID,
			&entry.ExerciseName,
			&entry.Sets,
			&entry.Reps,
			&entry.DurationSeconds,
			&entry.Weight,
			&entry.Notes,
			&entry.OrderIndex,
		)
		if err != nil {
			return nil, err
		}
		workout.Entries = append(workout.Entries, entry)
	}
	return workout, nil
}

func (pg *PostgresWorkoutStore) UpdateWorkout(workout *Workout) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	// rollback the transaction if it fails
	defer tx.Rollback()

	query := `
		UPDATE workouts SET title = $1, description = $2, duration_minutes = $3, calories_burned = $4
		WHERE id = $5
	`
	result, err := tx.Exec(query, workout.Title, workout.Description, workout.DurationMinutes, workout.CaloriesBurned, workout.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	// patch update:
	_, err = tx.Exec(`
		DELETE FROM workout_entries WHERE workout_id = $1
	`, workout.ID)
	if err != nil {
		return err
	}
	for _, entry := range workout.Entries {
		query := `
			INSERT INTO workout_entries (workout_id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id
		`
		err = tx.QueryRow(query,
			workout.ID,
			entry.ExerciseName,
			entry.Sets,
			entry.Reps,
			entry.DurationSeconds,
			entry.Weight,
			entry.Notes,
			entry.OrderIndex,
		).Scan(&entry.ID)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (pg *PostgresWorkoutStore) DeleteWorkout(id int64) error {
	query := `
		DELETE FROM workouts WHERE id = $1
	`
	result, err := pg.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	query = `
		DELETE FROM workout_entries WHERE workout_id = $1
	`
	_, err = pg.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
