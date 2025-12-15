package task

import (
	"context"
	"errors"
	"fmt"
	"techno/internal/config/logger"
	"techno/internal/model"
	rep "techno/internal/repository"
	"time"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

var _ rep.TaskRepository = (*repository)(nil)

type repository struct {
	pool *pgxpool.Pool
	log  zerolog.Logger
}

func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{
		pool: pool,
		log:  logger.GetLogger("repository.task"),
	}
}

func (r *repository) CreateTask(ctx context.Context, task *model.Task) error {
	start := time.Now()

	r.log.Info().
		Str("title", task.Title).
		Int("desc_len", len(task.Description)).
		Msg("Creating task")
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := "INSERT INTO tasks (title, description) VALUES ($1, $2) RETURNING id, status, created_at"
	err = tx.QueryRow(ctx, query, task.Title, task.Description).Scan(&task.ID, &task.Status, &task.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed created task: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		r.log.Error().Err(err).Int("task_id", task.ID).Msg("failed to commit transaction")
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	r.log.Info().
		Int("task_id", task.ID).
		Str("title", task.Title).
		Str("status", task.Status.StringStatus()).
		Dur("duration", time.Since(start)).
		Msg("Task created successfull")
	return nil
}

func (r *repository) GetByID(ctx context.Context, id int) (*model.Task, error) {
	query := "SELECT * FROM tasks WHERE id = $1"

	var task model.Task
	err := r.pool.QueryRow(ctx, query, id).Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("task with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	return &task, nil
}

func (r *repository) GetAll(ctx context.Context) ([]*model.Task, error) {
	rows, err := r.pool.Query(ctx, "SELECT * FROM tasks ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*model.Task
	for rows.Next() {
		task := &model.Task{}
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error: %w", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *repository) GetByStatus(ctx context.Context, status model.TaskStatus) ([]*model.Task, error) {
	query := "SELECT * FROM tasks WHERE status = $1 ORDER BY created_at DESC"
	start := time.Now()
	rows, err := r.pool.Query(ctx, query, status)
	if err != nil {
		return nil, fmt.Errorf("failed get task by status: %w", err)
	}

	defer rows.Close()

	var tasks []*model.Task
	for rows.Next() {
		task := &model.Task{}
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed scan task: %w", err)
		}
		tasks = append(tasks, task)
	}
	r.log.Info().
		Str("status", string(status)).
		Int("count", len(tasks)).
		Dur("duration", time.Since(start)).
		Msg("Retrieved tasks by status")

	return tasks, nil
}

func (r *repository) UpdateTask(ctx context.Context, task *model.Task) error {
	start := time.Now()

	r.log.Info().
		Int("task_id", task.ID).
		Str("new_title", task.Title).
		Msg("Updating task")
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := "UPDATE tasks SET title = $1, description = $2, status = $3 WHERE id = $4"

	result, err := tx.Exec(ctx, query, task.Title, task.Description, task.Status, task.ID)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("task with id %d not found", task.ID)
	}

	if err := tx.Commit(ctx); err != nil {
		r.log.Error().Err(err).Int("task_id", task.ID).Msg("failed to commit transaction")
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	r.log.Info().
		Int("task_id", task.ID).
		Str("title", task.Title).
		Str("status", task.Status.StringStatus()).
		Dur("duration", time.Since(start)).
		Msg("Task updated successfull")
	return nil
}

func (r *repository) DeleteTask(ctx context.Context, id int) error {
	start := time.Now()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := "DELETE FROM tasks WHERE id = $1"

	result, err := tx.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("task with id %d not found", id)
	}

	if err := tx.Commit(ctx); err != nil {
		r.log.Error().Err(err).Int("task_id", id).Msg("Failed to commit transaction")
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	r.log.Info().
		Int("task_id", id).
		Int64("rows_affected", rowsAffected).
		Dur("duration", time.Since(start)).
		Msg("Task deleted successfully")
	return nil
}
