package timer

import (
	"context"
	"fmt"
	"log"
	"techno/internal/config/logger"
	"techno/internal/model"
	"techno/internal/service"
	"time"

	"github.com/rs/zerolog"
)

type TaskCleaner struct {
	taskService service.TaskService
	interval    time.Duration
	stopChan    chan struct{}
	log         zerolog.Logger
}

func NewTaskCleaner(taskService service.TaskService, interval time.Duration) *TaskCleaner {
	return &TaskCleaner{
		taskService: taskService,
		interval:    interval,
		stopChan:    make(chan struct{}),
		log:         logger.GetLogger("timer.task_cleaner"),
	}
}

func (tc *TaskCleaner) Start(ctx context.Context) {
	ticker := time.NewTicker(tc.interval)
	defer ticker.Stop()

	log.Printf("Task cleaner started (interval: %v)", tc.interval)

	tc.cleanCompletedTasks(ctx)

	for {
		select {
		case <-ticker.C:
			tc.log.Debug().Msg("task cleaner tick")
			cleanCtx := context.Background()
			tc.cleanCompletedTasks(cleanCtx)

		case <-tc.stopChan:
			tc.log.Info().Msg("task cleaner stopped by Stop()")
			return

		case <-ctx.Done():
			tc.log.Info().Msg("task cleaner stopped by context")
			return
		}
	}
}

func (tc *TaskCleaner) Stop() {
	close(tc.stopChan)
}

func (tc *TaskCleaner) cleanCompletedTasks(ctx context.Context) {
	start := time.Now()
	tasks, err := tc.taskService.GetByStatus(ctx, model.Closed)
	if err != nil {
		tc.log.Error().
			Err(err).
			Dur("duration", time.Since(start)).
			Msg("Error getting completed tasks")
		return
	}

	tc.log.Info().
		Int("count", len(tasks)).
		Dur("duration", time.Since(start)).
		Msg("found completed tasks")

	if len(tasks) == 0 {
		log.Println("No completed tasks to clean")
		return
	}

	fmt.Printf("%-5s %-50s %-20s\n", "ID", "Title", "Completed At")

	for _, task := range tasks {
		title := task.Title
		if len(title) > 50 {
			title = title[:47] + "..."
		}
		fmt.Printf("%-5d %-50s %-20s\n", task.ID, title, task.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	deletedCount := 0
	for _, task := range tasks {
		if err := tc.taskService.DeleteTask(ctx, task.ID); err != nil {
			tc.log.Error().
				Err(err).
				Int("task_id", task.ID).
				Msg("Failed to delete task")
			continue
		}
		deletedCount++

		tc.log.Info().
			Int("task_id", task.ID).
			Str("title", task.Title).
			Msg("task deleted successfull")
	}

	log.Printf("Successfull deleted %d completed task(s)\n", deletedCount)
}
