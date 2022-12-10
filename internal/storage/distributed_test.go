package storage

import (
	"fmt"

	// api "github.com/djedjethai/generation/api/v1/keyvalue"
	// "golang.org/x/net/context"

	// // "github.com/djedjethai/proglog/internal/log"
	"io/ioutil"
	"net"
	"os"
	// "reflect"
	"testing"
	"time"

	"github.com/djedjethai/generation/internal/observability"
	"github.com/hashicorp/raft"
	"github.com/stretchr/testify/require"
	"github.com/travisjeffery/go-dynaport"
)

func TestMultipleNodes(t *testing.T) {
	obs := observability.Observability{}
	nbrShards := 2
	itemsShard := 5

	var logs []*DistributedLog
	nodeCount := 3
	ports := dynaport.Get(nodeCount)
	for i := 0; i < nodeCount; i++ {
		dataDir, err := ioutil.TempDir("", "distributed-log-test")
		require.NoError(t, err)
		defer func(dir string) {
			_ = os.RemoveAll(dir)
		}(dataDir)
		ln, err := net.Listen(
			"tcp",
			fmt.Sprintf("127.0.0.1:%d", ports[i]),
		)
		require.NoError(t, err)
		config := Config{}
		config.Raft.StreamLayer = NewStreamLayer(ln, nil, nil)
		config.Raft.LocalID = raft.ServerID(fmt.Sprintf("%d", i))
		config.Raft.HeartbeatTimeout = 50 * time.Millisecond
		config.Raft.ElectionTimeout = 50 * time.Millisecond
		config.Raft.LeaderLeaseTimeout = 50 * time.Millisecond
		config.Raft.CommitTimeout = 5 * time.Millisecond
		if i == 0 {
			config.Raft.Bootstrap = true
		}
		l, err := NewDistributedLog(dataDir, config, nbrShards, itemsShard, &obs)
		// l, err := NewDistributedLog(dataDir, config)
		require.NoError(t, err)
		if i != 0 {
			err = logs[0].Join(
				fmt.Sprintf("%d", i), ln.Addr().String(),
			)
			require.NoError(t, err)
		} else {
			err = l.WaitForLeader(3 * time.Second)
			require.NoError(t, err)
		}
		logs = append(logs, l)
	}

	// records := []*api.Records{
	// 	{Key: "firstKey", Value: "firstValue"},
	// 	{Key: "secondValue", Value: "secondValue"},
	// }

	// ctx := context.Background()
	// for _, record := range records {
	// 	err := logs[0].Set(ctx, record.Key, record.Value)
	// 	require.NoError(t, err)
	// 	require.Eventually(t, func() bool {
	// 		for j := 0; j < nodeCount; j++ {
	// 			got, err := logs[j].Get(ctx, record.Key)
	// 			if err != nil {
	// 				return false
	// 			}
	// 			// record.Offset = off
	// 			if !reflect.DeepEqual(got.(string), record.Value) {
	// 				return false
	// 			}
	// 		}
	// 		return true
	// 	}, 500*time.Millisecond, 50*time.Millisecond)
	// }

	// err := logs[0].Leave("1")
	// require.NoError(t, err)
	// time.Sleep(50 * time.Millisecond)
	// // off, err := logs[0].Append(&api.Record{
	// err = logs[0].Set(ctx, "thirdKey", "thirdValue")
	// require.NoError(t, err)
	// time.Sleep(50 * time.Millisecond)
	// record, err := logs[1].Get(ctx, "thirdKey")
	// require.IsType(t, api.ErrOffsetOutOfRange{}, err)
	// require.Nil(t, record)
	// record, err = logs[2].Get(ctx, "thirdKey")
	// require.NoError(t, err)
	// require.Equal(t, "thirdValue", record.(string))
	// // require.Equal(t, off, record.Offset)
}
