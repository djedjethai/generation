package config

import (
	// "crypto/tls"
	// "fmt"
	// "net"

	"github.com/djedjethai/generation/internal/deleter"
	"github.com/djedjethai/generation/internal/getter"
	// "github.com/djedjethai/generation/internal/observability"
	"github.com/djedjethai/generation/internal/setter"
	// "github.com/djedjethai/generation/internal/storage"
)

type Services struct {
	Setter  setter.Setter
	Getter  getter.Getter
	Deleter deleter.Deleter
}

// type GetServerer interface {
// 	GetServers() ([]*pb.Server, error)
// }

type PostgresDBParams struct {
	DbName   string
	Host     string
	User     string
	Password string
}
