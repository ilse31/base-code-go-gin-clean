package wire

import (
	"base-code-go-gin-clean/internal/config"
)

func ProvideConfig() (*config.Config, error) {
	return config.Load()
}

func ProvideDB(cfg *config.Config) (*config.DB, error) {
	return config.NewDB(cfg)
}
