package config

import (
	"os"
	"path/filepath"
)

var (
	CAFile         = configFile("ca.pem")
	ServerCertFile = configFile("server.pem")
	ServerKeyFile  = configFile("server-key.pem")
	ClientCertFile = configFile("client.pem")
	ClientKeyFile  = configFile("client-key.pem")
	// RootClientCertFile   = configFile("root-client.pem")
	// RootClientKeyFile    = configFile("root-client-key.pem")
	// NobodyClientCertFile = configFile("nobody-client.pem")
	// NobodyClientKeyFile  = configFile("nobody-client-key.pem")
)

func configFile(filename string) string {
	if dir := os.Getenv("CONFIG_DIR"); dir != "" {
		return filepath.Join(dir, filename)
	}

	// for development purpose
	// dir, err := os.Getwd()
	// if err != nil {
	// 	panic(err)
	// }

	// if filepath.Base(dir) == "cmd" {
	// 	return filepath.Join(dir, "../", ".generation", filename)
	// } else if filepath.Base(dir) == "grpc" {
	// 	return filepath.Join(dir, "../../..", ".generation", filename)
	// } else {
	// 	return filepath.Join(dir, "../..", ".generation", filename)
	// }

	// move to that for prod
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(homeDir, ".generation", filename)

}
