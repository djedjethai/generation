package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/djedjethai/clientGeneration0/pkg/config"
	pb "github.com/djedjethai/clientGeneration0/pkg/proto/keyvalue"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"io/ioutil"
	"net/http"
	"os"

	// "strings"
	"time"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	clientTLSConfig, err := config.SetupTLSConfig(config.TLSConfig{
		CertFile: config.ClientCertFile,
		KeyFile:  config.ClientKeyFile,
		CAFile:   config.CAFile,
	})
	if err != nil {
		log.Fatal("Error set the client: ", err)
	}
	clientCreds := credentials.NewTLS(clientTLSConfig)

	log.Println("see cert certif path: ", config.ClientCertFile)
	log.Println("see cert certif path1: ", config.ClientKeyFile)
	log.Println("see cert certif path2: ", config.CAFile)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(clientCreds),
		grpc.WithBlock(),
		grpc.WithTimeout(time.Second),
	}

	// var action, key, value string

	// if len(os.Args) > 2 {
	// 	action, key = os.Args[1], os.Args[2]
	// 	value = strings.Join(os.Args[3:], "")
	// }

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	type data struct {
		Action string `json:"action"`
		Key    string `json:"key"`
		Value  string `json:"value"`
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatalln(err)
		}

		dt := data{}
		err = json.Unmarshal(body, &dt)
		if err != nil {
			log.Fatal("err unmarshal: ", err)
		}

		log.Println("see the body: ", dt)

		conn, err := grpc.DialContext(ctx, "generation:8400", opts...)
		// conn, err := grpc.DialContext(ctx, ":8400", opts...)
		// conn, err := grpc.DialContext(ctx, "192.168.59.106:8400", opts...)
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}

		client := pb.NewKeyValueClient(conn)

		switch dt.Action {
		case "get":
			r, err := client.Get(ctx, &pb.GetRequest{Key: dt.Key})
			if err != nil {
				log.Fatalf("could not get value for key %s: %v", dt.Key, err)
			}

			log.Printf("Get %s returns: %s", dt.Key, r.Value)
		case "put":
			_, err := client.Put(ctx, &pb.PutRequest{
				Records: &pb.Records{Key: dt.Key, Value: dt.Value},
			})
			if err != nil {
				log.Fatalf("could not set value for key %s: %v", dt.Key, err)
			}

			log.Printf("Put: %s", dt.Key)
		case "delete":
			_, err := client.Delete(ctx, &pb.DeleteRequest{Key: dt.Key})
			if err != nil {
				log.Fatalf("could not delete value for key %s: %v", dt.Key, err)
			}

			log.Printf("Deleted: %s", dt.Key)
		case "getkeys":
			r, err := client.GetKeys(ctx, &pb.GetKeysRequest{})
			if err != nil {
				log.Fatalf("could not get keys: %v", err)
			}

			log.Printf("GetKeys returns: %s", r.Keys)
		case "getkeysvalues":
			stream, err := client.GetKeysValuesStream(ctx, &pb.Empty{})
			if err != nil {
				log.Fatalf("could not get keys: %v", err)
			}

			for {
				select {
				case <-stream.Context().Done():
					os.Exit(0)
				default:
					// Recieve on the stream
					res, err := stream.Recv()
					if errors.Is(err, io.EOF) {
						os.Exit(0)
					}
					if err != nil {
						os.Exit(1)
					}
					fmt.Println("The ressssult: ", res.Records)
				}
			}

		default:
			log.Fatalf("Syntax: go run [get|put|delete|getkeys] key value ...")
		}
	})

	// Run the web server.
	log.Println("start producer-api on port 3000 !!")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
