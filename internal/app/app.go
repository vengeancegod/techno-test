package app

import (
	"context"
	"techno/internal/config"

	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
)

type App struct {
	serviceProvider *serviceProvider
	rootCmd         *cobra.Command
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}
	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	return a.runCLI()
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initLogger,
		a.initCLI,
		//a.initWorker,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	loggerCfg := a.serviceProvider.LoggerConfig()
	if err := loggerCfg.Initialize(); err != nil {
		return err
	}
	return nil
}

func (a *App) initConfig(_ context.Context) error {
	err := config.Load(".env")
	if err != nil {
		return err
	}
	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initCLI(ctx context.Context) error {
	a.rootCmd = a.serviceProvider.RootCmd(ctx)
	return nil
}

func (a *App) runCLI() error {
	log.Info().Msg("Starting CLI")
	return a.rootCmd.Execute()
}

func (a *App) Stop(ctx context.Context) error {

	a.serviceProvider.Close()
	return nil
}
