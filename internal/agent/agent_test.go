package agent

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sync"
	"testing"
	"time"

	api "github.com/djedjethai/generation/api/v1/keyvalue"
	"github.com/djedjethai/generation/internal/config"
	"github.com/djedjethai/generation/internal/observability"
	"github.com/stretchr/testify/require"
	"github.com/travisjeffery/go-dynaport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func TestAgent(t *testing.T) {
	serverTLSConfig, err := config.SetupTLSConfig(config.TLSConfig{
		CertFile:      config.ServerCertFile,
		KeyFile:       config.ServerKeyFile,
		CAFile:        config.CAFile,
		Server:        true,
		ServerAddress: "127.0.0.1",
	})
	require.NoError(t, err)
	peerTLSConfig, err := config.SetupTLSConfig(config.TLSConfig{
		// CertFile:      config.RootClientCertFile,
		CertFile: config.ClientCertFile,
		// KeyFile:       config.RootClientKeyFile,
		KeyFile:       config.ClientKeyFile,
		CAFile:        config.CAFile,
		Server:        false,
		ServerAddress: "127.0.0.1",
	})
	require.NoError(t, err)
	// var agents []*agent.Agent
	var agents []*Agent
	for i := 0; i < 3; i++ {
		ports := dynaport.Get(2)
		bindAddr := fmt.Sprintf("%s:%d", "127.0.0.1", ports[0])
		rpcPort := ports[1]
		dataDir, err := ioutil.TempDir("", "agent-test-log")
		require.NoError(t, err)
		var startJoinAddrs []string
		if i != 0 {
			startJoinAddrs = append(
				startJoinAddrs,
				agents[0].config.BindAddr,
			)
		}

		// set the configs
		config := Config{
			NodeName:       fmt.Sprintf("%d", i),
			StartJoinAddrs: startJoinAddrs,
			BindAddr:       bindAddr,
			PortGRPC:       rpcPort,
			DataDir:        dataDir,
			// ACLModelFile:    config.ACLModelFile,
			// ACLPolicyFile:   config.ACLPolicyFile,
			ServerTLSConfig: serverTLSConfig,
			PeerTLSConfig:   peerTLSConfig,
			Bootstrap:       i == 0,
			//
			FileLoggerActive: false,
			DBLoggerActive:   false,
			Shards:           3,
			ItemsPerShard:    10,
			Protocol:         "grpc",
			IsTracing:        false,
			IsMetrics:        false,
			JaegerEndpoint:   "",
			Observability:    &observability.Observability{},
		}

		logger := observability.NewSrvLogger("debug")
		config.Observability.Logger = logger

		agent, err := New(config)
		require.NoError(t, err)
		agents = append(agents, agent)
	}
	defer func() {
		for _, agent := range agents {
			err := agent.Shutdown()
			require.NoError(t, err)
			require.NoError(t,
				os.RemoveAll(agent.config.DataDir),
			)
		}
	}()

	time.Sleep(3 * time.Second)
	leaderClient := client(t, agents[0], peerTLSConfig)

	// put 2 values(from leader)
	_, err = leaderClient.Put(
		context.Background(),
		&api.PutRequest{
			Records: &api.Records{
				Key:   "key1",
				Value: "value1",
			},
		},
	)
	require.NoError(t, err)

	_, err = leaderClient.Put(
		context.Background(),
		&api.PutRequest{
			Records: &api.Records{
				Key:   "key2",
				Value: "value2",
			},
		},
	)
	require.NoError(t, err)

	// get one value(from leader)
	consume, err := leaderClient.Get(
		context.Background(),
		&api.GetRequest{
			Key: "key1",
		},
	)
	require.NoError(t, err)
	require.Equal(t, consume.Value, "value1")

	// get an invalid key
	consume, err = leaderClient.Get(
		context.Background(),
		&api.GetRequest{
			Key: "key3",
		},
	)
	// TODO set an GRPC err..... ???
	require.Error(t, err, "need to set a grpc err")
	require.Nil(t, consume)

	// wait for propagation
	time.Sleep(3 * time.Second)

	// get all keys(from follower)
	followerClient := client(t, agents[1], peerTLSConfig)
	response, err := followerClient.GetKeys(
		context.Background(),
		&api.GetKeysRequest{},
	)
	require.NoError(t, err)
	require.Equal(t, len(response.Keys), 2)

	// get all keysvalues(from follower)
	keysvalues, err := followerClient.GetKeysValuesStream(
		context.Background(),
		&api.Empty{},
	)

	var keys []string
	var values []string
	var strErr = false
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		for {
			select {
			case <-keysvalues.Context().Done():
				return
			default:
				res, err := keysvalues.Recv()
				if errors.Is(err, io.EOF) {
					_ = keysvalues.CloseSend()
					return
				}
				if err != nil {
					strErr = true
					_ = keysvalues.CloseSend()
					return
				}
				keys = append(keys, res.Records.Key)
				values = append(values, res.Records.Value)
				wg.Done()
			}
		}
	}()
	wg.Wait()
	require.Equal(t, false, strErr)
	require.Equal(t, len(keys), 2)
	require.Equal(t, len(values), 2)

	//  delete a key(from leader)
	_, err = leaderClient.Delete(
		context.Background(),
		&api.DeleteRequest{
			Key: "key1",
		},
	)
	require.NoError(t, err)

	// wait for propagation
	time.Sleep(3 * time.Second)

	// get all keys(from leader), make sure the deleted one is gone
	response, err = followerClient.GetKeys(
		context.Background(),
		&api.GetKeysRequest{},
	)
	require.NoError(t, err)
	require.Equal(t, len(response.Keys), 1)
}

func client(t *testing.T, agent *Agent, tlsConfig *tls.Config) api.KeyValueClient {
	tlsCreds := credentials.NewTLS(tlsConfig)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(tlsCreds)}
	rpcAddr, err := agent.config.RPCAddr()
	require.NoError(t, err)

	conn, err := grpc.Dial(fmt.Sprintf(
		"%s",
		rpcAddr,
	), opts...)
	require.NoError(t, err)

	client := api.NewKeyValueClient(conn)
	return client
}
