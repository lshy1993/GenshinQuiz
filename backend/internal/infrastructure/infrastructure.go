package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/hibiken/asynq"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"

	"genshin-quiz-backend/internal/config"
)

// Infrastructure holds all the infrastructure dependencies
type Infrastructure struct {
	DB          *sql.DB
	Redis       *redis.Client
	Logger      *log.Logger
	TaskClient  *asynq.Client
	TaskServer  *asynq.Server
	Config      *config.Config
}

// NewInfrastructure creates a new infrastructure instance
func NewInfrastructure(cfg *config.Config) (*Infrastructure, error) {
	// Initialize logger
	logger, err := initLogger(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	// Initialize Sentry
	if err := initSentry(cfg); err != nil {
		logger.Printf("failed to initialize sentry: %v", err)
	}

	// Initialize database
	db, err := initDatabase(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize Redis
	redisClient, err := initRedis(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize redis: %w", err)
	}

	// Initialize Asynq client and server
	taskClient, taskServer := initTaskQueue(cfg, logger)

	return &Infrastructure{
		DB:         db,
		Redis:      redisClient,
		Logger:     logger,
		TaskClient: taskClient,
		TaskServer: taskServer,
		Config:     cfg,
	}, nil
}

// Close closes all infrastructure connections
func (i *Infrastructure) Close() error {
	var errs []error

	if i.DB != nil {
		if err := i.DB.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if i.Redis != nil {
		if err := i.Redis.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if i.TaskClient != nil {
		if err := i.TaskClient.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if i.TaskServer != nil {
		i.TaskServer.Shutdown()
	}

	if i.Logger != nil {
		_ = i.Logger.Sync()
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing infrastructure: %v", errs)
	}

	return nil
}

func initLogger(cfg *config.Config) (*log.Logger, error) {
	logger := log.New(log.Writer(), "[genshin-quiz] ", log.LstdFlags|log.Lshortfile)
	return logger, nil
}

func initSentry(cfg *config.Config) error {
	if cfg.SentryDSN == "" {
		return nil
	}

	return sentry.Init(sentry.ClientOptions{
		Dsn:         cfg.SentryDSN,
		Environment: cfg.Environment,
		Release:     cfg.Version,
	})
}

func initDatabase(cfg *config.Config, logger *log.Logger) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	logger.Printf("database connection established - host: %s, database: %s", cfg.Database.Host, cfg.Database.Name)

	return db, nil
}

func initRedis(cfg *config.Config, logger *log.Logger) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.Addr,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	logger.Printf("redis connection established - addr: %s, db: %d", cfg.Redis.Addr, cfg.Redis.DB)

	return client, nil
}

func initTaskQueue(cfg *config.Config, logger *log.Logger) (*asynq.Client, *asynq.Server) {
	redisOpt := asynq.RedisClientOpt{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	}

	client := asynq.NewClient(redisOpt)

	server := asynq.NewServer(redisOpt, asynq.Config{
		Concurrency: cfg.Worker.Concurrency,
		Queues: map[string]int{
			"critical": 6,
			"default":  3,
			"low":      1,
		},
		RetryDelayFunc: func(n int, e error, t *asynq.Task) time.Duration {
			return time.Duration(n) * time.Second
		},
		Logger: NewZapAsynqLogger(logger),
	})

	return client, server
}

// ZapAsynqLogger wraps zap.Logger to implement asynq.Logger interface
type ZapAsynqLogger struct {
	logger *zap.Logger
}

func NewZapAsynqLogger(logger *zap.Logger) *ZapAsynqLogger {
	return &ZapAsynqLogger{logger: logger}
}

func (l *ZapAsynqLogger) Debug(args ...interface{}) {
	l.logger.Sugar().Debug(args...)
}

func (l *ZapAsynqLogger) Info(args ...interface{}) {
	l.logger.Sugar().Info(args...)
}

func (l *ZapAsynqLogger) Warn(args ...interface{}) {
	l.logger.Sugar().Warn(args...)
}

func (l *ZapAsynqLogger) Error(args ...interface{}) {
	l.logger.Sugar().Error(args...)
}

func (l *ZapAsynqLogger) Fatal(args ...interface{}) {
	l.logger.Sugar().Fatal(args...)
}