package app

import (
	"context"
	"log"
	"techno/internal/config"
	"time"
)

type WorkerApp struct {
	serviceProvider *serviceProvider
	cleanerCtx      context.Context
	cleanerCancel   context.CancelFunc
}

func NewWorkerApp(ctx context.Context) (*WorkerApp, error) {
	a := &WorkerApp{}
	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (a *WorkerApp) Run() error {
	return a.runWorker()
}

func (a *WorkerApp) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initLogger,
		a.initWorker,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *WorkerApp) initConfig(_ context.Context) error {
	return config.Load(".env")
}

func (a *WorkerApp) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *WorkerApp) initWorker(ctx context.Context) error {
	a.cleanerCtx, a.cleanerCancel = context.WithCancel(context.Background())
	cleaner := a.serviceProvider.TaskCleaner(a.cleanerCtx)
	go cleaner.Start(a.cleanerCtx)
	return nil
}

func (a *WorkerApp) runWorker() error {
	log.Println("Starting Worker")

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-a.cleanerCtx.Done():
			return nil
		case <-ticker.C:

		}
	}
}

func (a *WorkerApp) initLogger(_ context.Context) error {
	loggerCfg := a.serviceProvider.LoggerConfig()
	if err := loggerCfg.Initialize(); err != nil {
		return err
	}
	return nil
}

func (a *WorkerApp) Stop(ctx context.Context) error {
	if a.cleanerCancel != nil {
		a.cleanerCancel()
	}

	if a.serviceProvider.taskCleaner != nil {
		a.serviceProvider.taskCleaner.Stop()
	}

	a.serviceProvider.Close()
	return nil
}
