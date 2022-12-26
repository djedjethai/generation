package config

import (
	"github.com/djedjethai/generation/internal/deleter"
	"github.com/djedjethai/generation/internal/getter"
	// "github.com/djedjethai/generation/internal/observability"
	"github.com/djedjethai/generation/internal/setter"
	// "github.com/djedjethai/generation/internal/storage"
	// pb "github.com/djedjethai/generation/api/v1/keyvalue"
)

type Services struct {
	Setter  setter.Setter
	Getter  getter.Getter
	Deleter deleter.Deleter
}

type PostgresDBParams struct {
	DbName   string
	Host     string
	User     string
	Password string
}
