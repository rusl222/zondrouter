package app

import (
	"context"
	"fmt"
	"os"

	"github.com/rusl222/zondrouter/internal/config"
	"github.com/rusl222/zondrouter/internal/gateway"
)

type App struct {
	//g *gateway.Gateway
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	if len(os.Args) < 2 {
		fmt.Println("Требуется путь к файлу конфигурации *.toml")
		os.Exit(1)
	}
	err := config.Load(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	return nil
}

func (a *App) Run() error {
	cfg, err := config.NewGatewayConfig()
	if err != nil {
		return err
	}
	return a.runGateway(cfg)
}

type GatewayConfigProvider interface {
	Masters() []config.Line
}

func (a *App) runGateway(provider GatewayConfigProvider) error {

	for _, dir := range provider.Masters() {
		var g gateway.Gateway

		g.Client = dir.Master
		g.Servers = dir.Slave

		go g.Run()

	}

	return nil
}
