package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	api "github.com/djedjethai/generation/api/v1/keyvalue"
	"github.com/djedjethai/generation/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	addr := flag.String("addr", ":8400", "service address")

	clientTLSConfig, err := config.SetupTLSConfig(config.TLSConfig{
		CertFile: config.ClientCertFile,
		KeyFile:  config.ClientKeyFile,
		CAFile:   config.CAFile,
	})

	clientCreds := credentials.NewTLS(clientTLSConfig)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(clientCreds),
		grpc.WithBlock(),
		grpc.WithTimeout(time.Second),
	}

	flag.Parse()
	conn, err := grpc.Dial(*addr, opts...)
	if err != nil {
		log.Fatal(err)
	}
	client := api.NewKeyValueClient(conn)
	ctx := context.Background()
	res, err := client.GetServers(ctx, &api.GetServersRequest{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("servers:")
	for _, server := range res.Servers {
		fmt.Printf("\t- %v\n", server)
	}
}
