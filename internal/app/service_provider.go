package app

import (
	"context"
	"log"
	"techno/internal/cli"
	"techno/internal/config/db"
	"techno/internal/config/logger"
	infra "techno/internal/db"
	"techno/internal/repository"
	taskRepo "techno/internal/repository/task"
	"techno/internal/service"
	taskService "techno/internal/service/task"
	"techno/internal/timer"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"
)

type serviceProvider struct {
	dbConfig     db.DBConfig
	loggerConfig logger.LoggerConfig
	db           *pgxpool.Pool

	taskRepository repository.TaskRepository
	taskService    service.TaskService
	taskCleaner    *timer.TaskCleaner
	taskCommands   *cli.TaskCommands
	rootCmd        *cobra.Command
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) LoggerConfig() logger.LoggerConfig {
	if s.loggerConfig == nil {
		cfg, err := logger.NewLoggerConfig()
		if err != nil {
			log.Fatal("failed to get log config: %s", err.Error())
		}
		s.loggerConfig = cfg
	}
	return s.loggerConfig
}

func (s *serviceProvider) DBConfig() db.DBConfig {
	if s.dbConfig == nil {
		cfg, err := db.NewDBConfig()
		if err != nil {
			log.Fatalf("failed to get db config: %s", err.Error())
		}
		s.dbConfig = cfg
	}
	return s.dbConfig
}

func (s *serviceProvider) DB(ctx context.Context) *pgxpool.Pool {
	if s.db == nil {
		pool, err := infra.InitDB(s.DBConfig())
		if err != nil {
			log.Fatalf("failed to connect to database: %s", err.Error())
		}
		s.db = pool
	}
	return s.db
}

func (s *serviceProvider) TaskRepository(ctx context.Context) repository.TaskRepository {
	if s.taskRepository == nil {
		s.taskRepository = taskRepo.NewRepository(s.DB(ctx))
	}
	return s.taskRepository
}

func (s *serviceProvider) TaskService(ctx context.Context) service.TaskService {
	if s.taskService == nil {
		s.taskService = taskService.NewService(s.TaskRepository(ctx))
	}
	return s.taskService
}

func (s *serviceProvider) TaskCommands(ctx context.Context) *cli.TaskCommands {
	if s.taskCommands == nil {
		s.taskCommands = cli.NewTaskCommands(s.TaskService(ctx))
	}
	return s.taskCommands
}

func (s *serviceProvider) RootCmd(ctx context.Context) *cobra.Command {
	if s.rootCmd == nil {
		s.rootCmd = cli.NewRootCommand()

		s.TaskCommands(ctx).RegisterCommands(s.rootCmd)
	}
	return s.rootCmd
}

func (s *serviceProvider) TaskCleaner(ctx context.Context) *timer.TaskCleaner {
	if s.taskCleaner == nil {
		s.taskCleaner = timer.NewTaskCleaner(s.TaskService(ctx), 30*time.Second) //5*time.Minute)
	}
	return s.taskCleaner
}

func (s *serviceProvider) Close() {
	if s.db != nil {
		s.db.Close()
		log.Println("Database connection closed")
	}
}
