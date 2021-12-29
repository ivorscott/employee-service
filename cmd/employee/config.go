package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ivorscott/employee-service/pkg/trace"
	"github.com/ivorscott/employee-service/res/database"

	"github.com/ardanlabs/conf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const logPath = "./log/out.log"

type appConfig struct {
	Web struct {
		Address         string        `conf:"default:localhost:4000"`
		Debug           string        `conf:"default:localhost:6060"`
		Production      bool          `conf:"default:false"`
		ReadTimeout     time.Duration `conf:"default:5s"`
		WriteTimeout    time.Duration `conf:"default:5s"`
		ShutdownTimeout time.Duration `conf:"default:5s"`
		FrontendAddress string        `conf:"default:https://localhost:3000"`
	}
	DB struct {
		User       string `conf:"default:postgres,noprint"`
		Password   string `conf:"default:postgres,noprint"`
		Host       string `conf:"default:localhost,noprint"`
		Name       string `conf:"default:employee,noprint"`
		DisableTLS bool   `conf:"default:false"`
	}
}

func newLoggerOrPanic() (*zap.Logger, func() error) {
	// fail immediately if we cannot log to file
	if _, err := os.OpenFile(logPath, os.O_RDONLY|os.O_CREATE, 0600); err != nil {
		panic(err)
	}

	// log retention policy
	lj := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    1000,
		MaxBackups: 100,
		MaxAge:     30,
		Compress:   true,
	}

	// integrate lumberjack logger with zap
	ws := zapcore.AddSync(lj)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, ws, zapcore.DebugLevel),        // log to file
		zapcore.NewCore(encoder, os.Stdout, zapcore.DebugLevel), // log to stdout
	)
	logger := zap.New(core, zap.AddCaller())
	return logger, logger.Sync
}

func newAppConfig() (appConfig, error) {
	var cfg appConfig

	if err := conf.Parse(os.Args[1:], "API", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("API", &cfg)
			if err != nil {
				return cfg, fmt.Errorf("error generating config usage: %w", err)
			}
			fmt.Println(usage)
			return cfg, nil
		}
		return cfg, fmt.Errorf("error parsing config: %w", err)
	}
	return cfg, nil
}

func newDBConnection(cfg appConfig) (*database.Repository, func(), error) {
	return database.NewRepository(database.Config{
		User:       cfg.DB.User,
		Host:       cfg.DB.Host,
		Name:       cfg.DB.Name,
		Password:   cfg.DB.Password,
		DisableTLS: cfg.DB.DisableTLS,
	})
}

func newTraceProviderGlobal() (func(ctx context.Context) error, error) {
	prv, err := trace.NewProvider(trace.ProviderConfig{
		JaegerEndpoint: "http://localhost:14268/api/traces",
		ServiceName:    "employee-service",
		ServiceVersion: "1.0.0",
		Environment:    "dev",
		Disabled:       false,
	})
	return prv.Close, err
}
