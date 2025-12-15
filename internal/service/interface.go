package service

import (
	"context"
	"techno/internal/model"
)

type TaskService interface {
	CreateTask(ctx context.Context, task *model.Task) error
	GetByID(ctx context.Context, id int) (*model.Task, error)
	GetAll(ctx context.Context) ([]*model.Task, error)
	GetByStatus(ctx context.Context, status model.TaskStatus) ([]*model.Task, error)
	UpdateTask(ctx context.Context, task *model.Task) error
	DeleteTask(ctx context.Context, id int) error
}
