package grpc

import (
	"context"
	// "fmt"
	"net"
	"reflect"
	"testing"

	pb "github.com/djedjethai/generation0/api/proto/keyvalue"
	"github.com/djedjethai/generation0/pkg/config"
	"github.com/djedjethai/generation0/pkg/deleter"
	"github.com/djedjethai/generation0/pkg/getter"
	lgr "github.com/djedjethai/generation0/pkg/logger"
	"github.com/djedjethai/generation0/pkg/observability"
	"github.com/djedjethai/generation0/pkg/setter"
	"github.com/djedjethai/generation0/pkg/storage"
	"github.com/stretchr/testify/require"
	gglGrpc "google.golang.org/grpc"
)

func setupTest(t *testing.T) (pb.KeyValueClient, func()) {
	t.Helper()
	s := gglGrpc.NewServer()
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	clientOptions := []gglGrpc.DialOption{gglGrpc.WithInsecure()}
	cc, err := gglGrpc.Dial(l.Addr().String(), clientOptions...)
	require.NoError(t, err)

	// set service
	obs := observability.Observability{}
	shardedMap := storage.NewShardedMap(2, 10, obs)
	setSrv := setter.NewSetter(shardedMap, obs)
	getSrv := getter.NewGetter(shardedMap, obs)
	delSrv := deleter.NewDeleter(shardedMap, obs)
	postgresConfig := config.PostgresDBParams{}
	loggerFacade, err := lgr.NewLoggerFacade(setSrv, delSrv, false, false, postgresConfig, "encryptionkey")
	require.NoError(t, err)

	pb.RegisterKeyValueServer(s, &Server{
		SetSrv:       setSrv,
		GetSrv:       getSrv,
		DelSrv:       delSrv,
		LoggerFacade: loggerFacade,
	})

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
		Key:   "key",
		Value: "value",
	})
	require.NoError(t, err)

	require.Equal(t, reflect.TypeOf(want), reflect.TypeOf(resp))
}

func TestGet(t *testing.T) {
	cl, teardown := setupTest(t)
	defer teardown()

	ctx := context.Background()

	_, err := cl.Put(ctx, &pb.PutRequest{
		Key:   "key",
		Value: "value",
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
		Key:   "key",
		Value: "value",
	})
	require.NoError(t, err)

	_, err = cl.Put(ctx, &pb.PutRequest{
		Key:   "key1",
		Value: "value1",
	})
	require.NoError(t, err)

	want := &pb.GetKeysResponse{
		Keys: []string{"key", "key1"},
	}

	resp, err := cl.GetKeys(ctx, &pb.GetKeysRequest{})
	require.NoError(t, err)

	require.Equal(t, want.Keys, resp.Keys)
}

func TestDelete(t *testing.T) {
	cl, teardown := setupTest(t)
	defer teardown()

	ctx := context.Background()

	_, err := cl.Put(ctx, &pb.PutRequest{
		Key:   "key",
		Value: "value",
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
