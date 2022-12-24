package main

import (
	"context"
	"flag"
	"fmt"
	api "github.com/djedjethai/generation/api/v1/keyvalue"
	"google.golang.org/grpc"
	"log"
)

func main() {
	addr := flag.String("addr", ":8400", "service address")
	flag.Parse()
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
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
