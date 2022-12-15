package storage

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	api "github.com/djedjethai/generation/api/v1/keyvalue"
	"github.com/djedjethai/generation/internal/models"
	"github.com/djedjethai/generation/internal/observability"
	"github.com/djedjethai/generation/internal/raftlog"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"

	"google.golang.org/protobuf/proto"
)

type DistributedStorage struct {
	logConfig raftlog.Config
	config    Config
	log       *raftlog.Log
	sm        *ShardedMap
	raft      *raft.Raft
}

func NewDistributedStorage(dataDir string, conf Config, nShard, maxLgt int, observ *observability.Observability) (*DistributedStorage, error) {
	l := &DistributedStorage{
		logConfig: raftlog.Config{},
		config:    conf,
	}

	if err := l.setupShardedMap(nShard, maxLgt, observ); err != nil {
		return nil, err
	}
	if err := l.setupRaft(dataDir); err != nil {
		return nil, err
	}
	return l, nil
}

func (l *DistributedStorage) setupShardedMap(nShard, maxLgt int, observ *observability.Observability) error {
	if nShard < 1 || maxLgt < 1 {
		return errors.New("Storage needs some rooms")
	}
	nsm := NewShardedMap(nShard, maxLgt, observ)
	l.sm = &nsm
	return nil
}

func (l *DistributedStorage) setupRaft(dataDir string) error {
	fsm := &fsm{sm: l.sm}
	logDir := filepath.Join(dataDir, "raft", "log")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}
	logConfig := l.logConfig
	logConfig.Segment.InitialOffset = 1
	logStore, err := newLogStore(logDir, logConfig)
	if err != nil {
		return err
	}

	stableStore, err := raftboltdb.NewBoltStore(
		filepath.Join(dataDir, "raft", "stable"),
	)
	if err != nil {
		return err
	}
	retain := 1
	snapshotStore, err := raft.NewFileSnapshotStore(
		filepath.Join(dataDir, "raft"),
		retain,
		os.Stderr,
	)
	if err != nil {
		return err
	}
	maxPool := 5
	timeout := 10 * time.Second
	transport := raft.NewNetworkTransport(
		l.config.Raft.StreamLayer,
		maxPool,
		timeout,
		os.Stderr,
	)
	config := raft.DefaultConfig()
	config.LocalID = l.config.Raft.LocalID
	if l.config.Raft.HeartbeatTimeout != 0 {
		config.HeartbeatTimeout = l.config.Raft.HeartbeatTimeout
	}
	if l.config.Raft.ElectionTimeout != 0 {
		config.ElectionTimeout = l.config.Raft.ElectionTimeout
	}
	if l.config.Raft.LeaderLeaseTimeout != 0 {
		config.LeaderLeaseTimeout = l.config.Raft.LeaderLeaseTimeout
	}
	if l.config.Raft.CommitTimeout != 0 {
		config.CommitTimeout = l.config.Raft.CommitTimeout
	}
	l.raft, err = raft.NewRaft(
		config,
		fsm,
		logStore,
		stableStore,
		snapshotStore,
		transport,
	)
	if err != nil {
		return err
	}
	hasState, err := raft.HasExistingState(
		logStore,
		stableStore,
		snapshotStore,
	)
	if err != nil {
		return err
	}

	if l.config.Raft.Bootstrap && !hasState {
		configA := raft.Configuration{
			Servers: []raft.Server{{
				ID:      config.LocalID,
				Address: transport.LocalAddr(),
			}},
		}
		err = l.raft.BootstrapCluster(configA).Error()
	}
	return err
	// return nil
}

// should have Put/Get/Delete
func (l *DistributedStorage) Set(ctx context.Context, key string, value interface{}) error {
	_, err := l.apply(
		SetRequestType,
		&api.Records{
			Key:   key,
			Value: value.(string),
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (l *DistributedStorage) Get(ctx context.Context, key string) (interface{}, error) {
	res, err := l.apply(
		GetRequestType,
		&api.Records{
			Key: key,
		},
	)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (l *DistributedStorage) Delete(ctx context.Context, key string, sh *Shard) error {
	_, err := l.apply(
		DeleteRequestType,
		&api.Records{
			Key: key,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

// no replication to followers nodes
// will be access directly from followers nodes
func (l *DistributedStorage) Keys(ctx context.Context) []string {
	return l.sm.Keys(ctx)
}

// no replication to followers nodes
// will be access directly from followers nodes
func (l *DistributedStorage) KeysValues(ctx context.Context, ch chan models.KeysValues) error {
	return l.sm.KeysValues(ctx, ch)
}

// for testing purpose only
func (l *DistributedStorage) Read(ctx context.Context, key string) (string, error) {
	val, err := l.sm.Get(ctx, key)
	if err != nil {
		return "", err
	}
	return val.(string), nil
}

// apply will switch on the RequestType(Put/Get/Delete)
func (l *DistributedStorage) apply(reqType RequestType, req proto.Message) (interface{}, error) {
	var buf bytes.Buffer
	_, err := buf.Write([]byte{byte(reqType)})
	if err != nil {
		return nil, err
	}

	b, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}

	timeout := 10 * time.Second
	future := l.raft.Apply(buf.Bytes(), timeout)
	if future.Error() != nil {
		return nil, future.Error()
	}

	res := future.Response()
	if err, ok := res.(error); ok {
		return nil, err
	}

	return res, nil
}

var _ raft.FSM = (*fsm)(nil)

type fsm struct {
	sm *ShardedMap
}

type RequestType uint8

const (
	SetRequestType RequestType = iota
	GetRequestType
	DeleteRequestType
	AppendRequestType
)

// will switch on reqType(Put/Get/Delete)
func (l *fsm) Apply(record *raft.Log) interface{} {
	buf := record.Data
	reqType := RequestType(buf[0])
	switch reqType {
	case SetRequestType:
		return l.applySet(buf[1:])
	case GetRequestType:
		return l.applyGet(buf[1:])
	case DeleteRequestType:
		return l.applyDelete(buf[1:])
	}
	return nil
}

// will have applyPut(records)/applyGet(key)/applyDelete(key)
func (l *fsm) applySet(b []byte) interface{} {
	// var req api.PutRequest
	var req api.Records
	err := proto.Unmarshal(b, &req)
	if err != nil {
		return err
	}

	ctx := context.Background()

	// fmt.Println("see in applySet: ", req.Records)
	// err = l.sm.Set(ctx, req.Records.Key, req.Records.Value)
	err = l.sm.Set(ctx, req.Key, req.Value)
	if err != nil {
		return err
	}

	// will return the expected response from the method executed on the storage
	return err
}

func (l *fsm) applyGet(b []byte) interface{} {
	var req api.GetRequest
	err := proto.Unmarshal(b, &req)
	if err != nil {
		fmt.Println("seee the get err: ", err)
		return err
	}

	ctx := context.Background()

	ndVal, err := l.sm.Get(ctx, req.Key)
	if err != nil {
		return err
	}

	// will return the expected response from the method executed on the storage
	// return ndVal, err
	// TODO here I cut the err .....
	return ndVal
}

func (l *fsm) applyDelete(b []byte) interface{} {
	var req api.DeleteRequest
	err := proto.Unmarshal(b, &req)
	if err != nil {
		return err
	}

	ctx := context.Background()

	// TODO see here the story of the "nil" shardedMap.....
	err = l.sm.Delete(ctx, req.Key, nil)
	if err != nil {
		return err
	}

	// will return the expected response from the method executed on the storage
	return err
}

// will read all the storage and snapshot it
// should snapshot to the db...
func (l *fsm) Snapshot() (raft.FSMSnapshot, error) {
	ctx := context.Background()
	var ch = make(chan models.KeysValues)

	go func() {
		err := l.sm.KeysValues(ctx, ch)
		if err != nil {
			return
		}
	}()

	r := strings.NewReader("")
	buf := new(bytes.Buffer)
	for d := range ch {

		b, err := proto.Marshal(&api.Records{
			Key:   d.Key,
			Value: d.Value,
		})
		if err != nil {
			fmt.Println("err marshaling in snapshot()")
		}
		buf.Write(b)
		buf.WriteString("\n")
	}
	r.Read(buf.Bytes())

	return &snapshot{reader: r}, nil
}

var _ raft.FSMSnapshot = (*snapshot)(nil)

type snapshot struct {
	reader io.Reader
}

// a way to persist the fsm to the disk, memory, s3 or whatever...
// if a node goes down, when back up it will use this snapshot to reset its fsm state
func (s *snapshot) Persist(sink raft.SnapshotSink) error {
	if _, err := io.Copy(sink, s.reader); err != nil {
		_ = sink.Cancel()
		return err
	}
	return sink.Close()
}

func (s *snapshot) Release() {}

// here read the snapshot from r io.ReadCloser(how come ??),
// but in our case it will be the snapshot of the storage(which is managed by the fsm)
// and reset the storage(executing the snapshot records) of the newly up node
func (l *fsm) Restore(r io.ReadCloser) error {

	ctx := context.Background()

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		bts := scanner.Bytes()
		dt := &api.Records{}
		err := proto.Unmarshal(bts, dt)
		if err != nil {
			return err
		}
		err = l.sm.Set(ctx, dt.Key, dt.Value)
		if err != nil {
			return err
		}
	}
	return nil
}

// raft logger
var _ raft.LogStore = (*logStore)(nil)

type logStore struct {
	*raftlog.Log
}

func newLogStore(dir string, c raftlog.Config) (*logStore, error) {
	log, err := raftlog.NewLog(dir, c)
	if err != nil {
		return nil, err
	}
	return &logStore{log}, nil
}

func (l *logStore) FirstIndex() (uint64, error) {
	return l.LowestOffset()
}

func (l *logStore) LastIndex() (uint64, error) {
	off, err := l.HighestOffset()
	return off, err
}

func (l *logStore) GetLog(index uint64, out *raft.Log) error {
	in, err := l.Read(index)
	if err != nil {
		return err
	}
	out.Data = in.Value
	out.Index = in.Offset
	out.Type = raft.LogType(in.Type)
	out.Term = in.Term
	return nil
}

func (l *logStore) StoreLog(record *raft.Log) error {
	return l.StoreLogs([]*raft.Log{record})
}
func (l *logStore) StoreLogs(records []*raft.Log) error {

	for _, record := range records {
		if _, err := l.Append(&raftlog.Record{
			Value: record.Data,
			Term:  record.Term,
			Type:  uint32(record.Type),
		}); err != nil {
			return err
		}
	}

	return nil
}

func (l *logStore) DeleteRange(min, max uint64) error {
	return l.Truncate(max)
}

// stream
var _ raft.StreamLayer = (*StreamLayer)(nil)

type StreamLayer struct {
	ln              net.Listener
	serverTLSConfig *tls.Config
	peerTLSConfig   *tls.Config
}

func NewStreamLayer(ln net.Listener, serverTLSConfig *tls.Config, peerTLSConfig *tls.Config) *StreamLayer {
	return &StreamLayer{
		ln:              ln,
		serverTLSConfig: serverTLSConfig,
		peerTLSConfig:   peerTLSConfig,
	}
}

const RaftRPC = 1

func (s *StreamLayer) Dial(addr raft.ServerAddress, timeout time.Duration) (net.Conn, error) {
	dialer := &net.Dialer{Timeout: timeout}
	var conn, err = dialer.Dial("tcp", string(addr))
	if err != nil {
		return nil, err
	}
	// identify to mux this is a raft rpc
	_, err = conn.Write([]byte{byte(RaftRPC)})
	if err != nil {
		return nil, err
	}
	if s.peerTLSConfig != nil {
		conn = tls.Client(conn, s.peerTLSConfig)
	}
	return conn, err
}

func (s *StreamLayer) Accept() (net.Conn, error) {
	conn, err := s.ln.Accept()
	if err != nil {
		return nil, err
	}
	b := make([]byte, 1)
	_, err = conn.Read(b)
	if err != nil {
		return nil, err
	}
	if bytes.Compare([]byte{byte(RaftRPC)}, b) != 0 {
		return nil, fmt.Errorf("not a raft rpc")
	}
	if s.serverTLSConfig != nil {
		return tls.Server(conn, s.serverTLSConfig), nil
	}
	return conn, nil
}

func (s *StreamLayer) Close() error {
	return s.ln.Close()
}

func (s *StreamLayer) Addr() net.Addr {
	return s.ln.Addr()
}

func (l *DistributedStorage) Join(id, addr string) error {
	configFuture := l.raft.GetConfiguration()
	if err := configFuture.Error(); err != nil {
		return err
	}
	serverID := raft.ServerID(id)
	serverAddr := raft.ServerAddress(addr)
	for _, srv := range configFuture.Configuration().Servers {
		if srv.ID == serverID || srv.Address == serverAddr {
			if srv.ID == serverID && srv.Address == serverAddr {
				// server has already joined
				return nil
			}
			// remove the existing server
			removeFuture := l.raft.RemoveServer(serverID, 0, 0)
			if err := removeFuture.Error(); err != nil {
				return err
			}
		}
	}
	addFuture := l.raft.AddVoter(serverID, serverAddr, 0, 0)
	if err := addFuture.Error(); err != nil {
		return err
	}
	return nil
}

func (l *DistributedStorage) Leave(id string) error {
	removeFuture := l.raft.RemoveServer(raft.ServerID(id), 0, 0)
	return removeFuture.Error()
}

func (l *DistributedStorage) WaitForLeader(timeout time.Duration) error {
	fmt.Println("set new leader: ", l.raft.Leader())
	timeoutc := time.After(timeout)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-timeoutc:
			return fmt.Errorf("timed out")
		case <-ticker.C:
			if l := l.raft.Leader(); l != "" {
				return nil
			}
		}
	}
}

func (l *DistributedStorage) Close() error {
	f := l.raft.Shutdown()
	if err := f.Error(); err != nil {
		return err
	}
	return nil
	// return l.log.Close()
}
