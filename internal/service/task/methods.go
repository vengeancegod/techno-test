package task

import (
	"context"
	"fmt"
	"strings"
	"techno/internal/model"
)

func (s *service) CreateTask(ctx context.Context, task *model.Task) error {
	task.Title = strings.TrimSpace(task.Title)
	task.Description = strings.TrimSpace(task.Description)

	if err := s.taskRepository.CreateTask(ctx, task); err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	return nil
}

func (s *service) GetByID(ctx context.Context, id int) (*model.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid task id: %d", id)
	}

	task, err := s.taskRepository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return task, nil
}

func (s *service) GetAll(ctx context.Context) ([]*model.Task, error) {
	tasks, err := s.taskRepository.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all tasks: %w", err)
	}

	if tasks == nil {
		return []*model.Task{}, nil
	}

	return tasks, nil
}

func (s *service) GetByStatus(ctx context.Context, status model.TaskStatus) ([]*model.Task, error) {
	tasks, err := s.taskRepository.GetByStatus(ctx, status)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks by status: %w", err)
	}

	if tasks == nil {
		return []*model.Task{}, nil
	}

	return tasks, nil
}

func (s *service) UpdateTask(ctx context.Context, task *model.Task) error {
	if task.ID <= 0 {
		return fmt.Errorf("invalid task id: %d", task.ID)
	}

	existingTask, err := s.taskRepository.GetByID(ctx, task.ID)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}

	task.Title = strings.TrimSpace(task.Title)
	task.Description = strings.TrimSpace(task.Description)

	task.CreatedAt = existingTask.CreatedAt

	if err := s.taskRepository.UpdateTask(ctx, task); err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

func (s *service) DeleteTask(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid task id: %d", id)
	}

	_, err := s.taskRepository.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}

	if err := s.taskRepository.DeleteTask(ctx, id); err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}
