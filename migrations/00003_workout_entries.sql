-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS workout_entries (
  -- workouts 和 workout_entries 的关系是一对多 (1:N) 的父子关系。
  -- 所有动作都属于 workout_id 的workout
  id BIGSERIAL PRIMARY KEY,
  workout_id BIGINT NOT NULL REFERENCES workouts(id) ON DELETE CASCADE,
  -- ON DELETE CASCADE： 如果我删除了这个训练workout（比如删除了“胸肌训练”），那么在这个训练里的所有动作（卧推、飞鸟等）也会自动被删除。
  exercise_name VARCHAR(255) NOT NULL,
  sets INTEGER NOT NULL,
  reps INTEGER,
  duration_seconds INTEGER,
  weight DECIMAL(5, 2),
  notes TEXT,
  order_index INTEGER NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT valid_workout_entry CHECK (
    (reps IS NOT NULL OR duration_seconds IS NOT NULL) AND
    (reps IS NULL OR duration_seconds IS NULL)
  )
  -- 这段 CHECK 约束确保了 reps（次数）和 duration_seconds（持续时间）二者必选其一，且只能选其一（异或逻辑）。
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE workout_entries;
-- +goose StatementEnd