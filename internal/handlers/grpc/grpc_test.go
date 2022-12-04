package grpc

import (
	"context"
	// "fmt"
	"net"
	"reflect"
	"testing"

	pb "github.com/djedjethai/generation/api/v1/keyvalue"
	"github.com/djedjethai/generation/internal/config"
	"github.com/djedjethai/generation/internal/deleter"
	"github.com/djedjethai/generation/internal/getter"
	lgr "github.com/djedjethai/generation/internal/logger"
	"github.com/djedjethai/generation/internal/observability"
	"github.com/djedjethai/generation/internal/setter"
	"github.com/djedjethai/generation/internal/storage"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	gglGrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func setupTest(t *testing.T) (pb.KeyValueClient, func()) {
	t.Helper()
	// s := gglGrpc.NewServer()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	// set tls for the client
	clientTLSConfig, err := config.SetupTLSConfig(config.TLSConfig{
		CertFile: config.ClientCertFile,
		KeyFile:  config.ClientKeyFile,
		CAFile:   config.CAFile,
	})
	require.NoError(t, err)
	clientCreds := credentials.NewTLS(clientTLSConfig)

	// clientOptions := []gglGrpc.DialOption{gglGrpc.WithInsecure()}
	// cc, err := gglGrpc.Dial(l.Addr().String(), clientOptions...)
	cc, err := gglGrpc.Dial(
		l.Addr().String(),
		grpc.WithTransportCredentials(clientCreds),
	)
	require.NoError(t, err)

	// set service
	obs := observability.Observability{}
	shardedMap := storage.NewShardedMap(2, 10, obs)
	setSrv := setter.NewSetter(shardedMap, obs)
	getSrv := getter.NewGetter(shardedMap, obs)
	delSrv := deleter.NewDeleter(shardedMap, obs)
	postgresConfig := config.PostgresDBParams{}
	srv := config.Services{setSrv, getSrv, delSrv}
	loggerFacade, err := lgr.NewLoggerFacade(srv, false, postgresConfig)
	require.NoError(t, err)

	serverTLSConfig, err := config.SetupTLSConfig(config.TLSConfig{
		CertFile:      config.ServerCertFile,
		KeyFile:       config.ServerKeyFile,
		CAFile:        config.CAFile,
		ServerAddress: l.Addr().String(),
	})
	require.NoError(t, err)
	serverCreds := credentials.NewTLS(serverTLSConfig)
	require.NoError(t, err)

	s, err := NewGRPCServer(&srv, loggerFacade, gglGrpc.Creds(serverCreds))
	require.NoError(t, err)

	go func() {
		s.Serve(l)
	}()

	clientMocked := pb.NewKeyValueClient(cc)
	return clientMocked, func() {
		s.Stop()
		cc.Close()
		l.Close()
	}
}

func TestSet(t *testing.T) {
	cl, teardown := setupTest(t)
	defer teardown()

	ctx := context.Background()

	want := &pb.PutResponse{}

	resp, err := cl.Put(ctx, &pb.PutRequest{
		Records: &pb.Records{
			Key:   "key",
			Value: "value",
		},
	})
	require.NoError(t, err)

	require.Equal(t, reflect.TypeOf(want), reflect.TypeOf(resp))
}

func TestGet(t *testing.T) {
	cl, teardown := setupTest(t)
	defer teardown()

	ctx := context.Background()

	_, err := cl.Put(ctx, &pb.PutRequest{
		Records: &pb.Records{
			Key:   "key",
			Value: "value",
		},
	})
	require.NoError(t, err)

	want := &pb.GetResponse{
		Value: "value",
	}

	resp, err := cl.Get(ctx, &pb.GetRequest{
		Key: "key",
	})
	require.NoError(t, err)

	require.Equal(t, want.Value, resp.Value)
}

func TestGetKeys(t *testing.T) {
	cl, teardown := setupTest(t)
	defer teardown()

	ctx := context.Background()

	_, err := cl.Put(ctx, &pb.PutRequest{
		Records: &pb.Records{
			Key:   "key",
			Value: "value",
		},
	})
	require.NoError(t, err)

	_, err = cl.Put(ctx, &pb.PutRequest{
		Records: &pb.Records{
			Key:   "key1",
			Value: "value1",
		},
	})
	require.NoError(t, err)

	want := &pb.GetKeysResponse{
		Keys: []string{"key", "key1"},
	}

	resp, err := cl.GetKeys(ctx, &pb.GetKeysRequest{})
	require.NoError(t, err)

	require.Equal(t, len(want.Keys), len(resp.Keys))
}

func TestDelete(t *testing.T) {
	cl, teardown := setupTest(t)
	defer teardown()

	ctx := context.Background()

	_, err := cl.Put(ctx, &pb.PutRequest{
		Records: &pb.Records{
			Key:   "key",
			Value: "value",
		},
	})
	require.NoError(t, err)

	want := &pb.DeleteResponse{}

	resp, err := cl.Delete(ctx, &pb.DeleteRequest{
		Key: "key",
	})
	require.NoError(t, err)

	resp1, err := cl.GetKeys(ctx, &pb.GetKeysRequest{})
	require.NoError(t, err)

	require.Equal(t, reflect.TypeOf(want), reflect.TypeOf(resp))
	require.Equal(t, len(resp1.Keys), 0)
}
