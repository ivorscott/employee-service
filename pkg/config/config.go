// Package config configures application variables.
package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ardanlabs/conf"
)

// AppConfig contains application configuration with good defaults.
type AppConfig struct {
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
		Port       int    `conf:"default:5432,noprint"`
		Name       string `conf:"default:employee,noprint"`
		DisableTLS bool   `conf:"default:false"`
	}
	TestDB struct {
		User       string `conf:"default:postgres,noprint"`
		Password   string `conf:"default:postgres,noprint"`
		Host       string `conf:"default:localhost,noprint"`
		Port       int    `conf:"default:5442,noprint"`
		Name       string `conf:"default:employee_test,noprint"`
		DisableTLS bool   `conf:"default:true"`
	}
}

// NewAppConfig creates a new AppConfig for the application.
func NewAppConfig() (*AppConfig, error) {
	var cfg AppConfig

	if err := conf.Parse(os.Args[1:], "API", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("API", &cfg)
			if err != nil {
				return nil, fmt.Errorf("error generating config usage: %w", err)
			}
			fmt.Println(usage)
			return nil, nil
		}
		return nil, fmt.Errorf("error parsing config: %w", err)
	}
	return &cfg, nil
}
