package task

import ("techno/internal/repository"
	def "techno/internal/service"
)

var _ def.TaskService = (*service)(nil)

type service struct {
	taskRepository repository.TaskRepository
}

func NewService(taskRepository repository.TaskRepository) *service {
	return &service{
		taskRepository: taskRepository,
	}
}