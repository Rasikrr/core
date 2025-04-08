package application

import (
	"context"
	"github.com/Rasikrr/core/brokers/nats"
	"github.com/Rasikrr/core/config"
	"github.com/Rasikrr/core/database"
	coreGrpc "github.com/Rasikrr/core/grpc"
	"github.com/Rasikrr/core/http"
	"github.com/Rasikrr/core/redis"

	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	name       string
	config     *config.Config
	redis      redis.Cache
	postgres   *database.Postgres
	httpServer *http.Server
	grpcServer *coreGrpc.Server

	publisher          nats.Publisher
	subscriber         nats.Subscriber
	subscriberHandlers []nats.SubscriberHandler

	starters Starters
	closers  Closers
}

func NewApp(ctx context.Context) *App {
	cfg, err := config.Parse()
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}
	return NewAppWithConfig(ctx, &cfg)
}

func NewAppWithConfig(ctx context.Context, cfg *config.Config) *App {
	app := &App{
		name:   cfg.Name(),
		config: cfg,
	}
	if err := app.initPostgres(ctx); err != nil {
		log.Fatalf("failed to init postgres: %v", err)
	}
	if err := app.initRedis(ctx); err != nil {
		log.Fatalf("failed to init redis: %v", err)
	}
	if err := app.initGRPC(ctx); err != nil {
		log.Fatalf("failed to init grpc: %v", err)
	}
	if err := app.initHTTP(ctx); err != nil {
		log.Fatalf("failed to init http: %v", err)
	}
	if err := app.initNats(ctx); err != nil {
		log.Fatalf("failed to init nats: %v", err)
	}

	return app
}

func (a *App) Start(ctx context.Context) error {
	a.initSubscribers(ctx)

	stopChan := make(chan struct{})
	go a.GracefulShutdown(ctx, stopChan)
	if err := a.start(ctx); err != nil {
		return err
	}
	<-stopChan

	return nil
}

func (a *App) start(ctx context.Context) error {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("panic recovered: %v", err)
		}
	}()
	g := errgroup.Group{}
	for _, s := range a.starters.starters {
		g.Go(func() error {
			return s.Start(ctx)
		})
	}
	return g.Wait()
}

func (a *App) Close(ctx context.Context) error {
	for _, c := range a.closers.closers {
		if err := c.Close(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) GracefulShutdown(ctx context.Context, stopChan chan struct{}) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-sigChan

	if err := a.Close(ctx); err != nil {
		log.Printf("error while closing app: %v", err)
	}
	close(stopChan)
}

func (a *App) GrpcServer() *coreGrpc.Server {
	if a.grpcServer == nil {
		log.Fatalf("grpc server is not initialized or not required, please check your config")
	}
	return a.grpcServer
}

func (a *App) Postgres() *database.Postgres {
	if a.postgres == nil {
		log.Fatalf("postgres is not initialized or not required. please check your config")
	}
	return a.postgres
}

func (a *App) HTTPServer() *http.Server {
	if a.httpServer == nil {
		log.Fatalf("http server is not initialized or not required. please check your config")
	}
	return a.httpServer
}

func (a *App) Redis() redis.Cache {
	if a.redis == nil {
		log.Fatalf("redis is not initialized or not required. please check your config")
	}
	return a.redis
}

func (a *App) Config() *config.Config {
	return a.config
}

func (a *App) NATSPublisher() nats.Publisher {
	if a.publisher == nil {
		log.Fatalf("nats is not initialized or not required. please check your config")
	}
	return a.publisher
}

func (a *App) NATSSubscriber() nats.Subscriber {
	if a.subscriber == nil {
		log.Fatalf("nats is not initialized or not required. please check your config")
	}
	return a.subscriber
}

func (a *App) WithSubscribers(handlers ...nats.SubscriberHandler) {
	a.subscriberHandlers = append(a.subscriberHandlers, handlers...)
}

func (a *App) initSubscribers(_ context.Context) {
	if a.Config().NATS.Required && a.subscriber != nil {
		a.subscriber.WithHandlers(a.subscriberHandlers...)
	}
}
