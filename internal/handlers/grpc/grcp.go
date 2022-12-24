package grpc

import (
	"context"
	"fmt"

	// "fmt"

	pb "github.com/djedjethai/generation/api/v1/keyvalue"
	"github.com/djedjethai/generation/internal/config"
	"github.com/djedjethai/generation/internal/logger"
	"github.com/djedjethai/generation/internal/models"
	"google.golang.org/grpc"
	// "google.golang.org/grpc/status"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// type GetServerer interface {
// 	GetServers() ([]*pb.Server, error)
// }

type Server struct {
	pb.UnimplementedKeyValueServer
	Services     *config.Services
	LoggerFacade *logger.LoggerFacade
}

func NewGRPCServer(services config.Services, loggerFacade *logger.LoggerFacade, opts ...grpc.ServerOption) (*grpc.Server, error) {

	gsrv := grpc.NewServer(opts...)
	// gsrv := grpc.NewServer() // uncomment here for no tls

	// set healthcheck end point
	hsrv := health.NewServer()
	hsrv.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(gsrv, hsrv)

	srv, err := newgrpcserver(&services, loggerFacade)
	if err != nil {
		return nil, err
	}

	pb.RegisterKeyValueServer(gsrv, srv)

	return gsrv, nil
}

func newgrpcserver(services *config.Services, loggerFacade *logger.LoggerFacade) (*Server, error) {
	return &Server{
		Services:     services,
		LoggerFacade: loggerFacade,
	}, nil
}

// func (s *Server) GetServers(ctx context.Context, req *pb.GetServersRequest) (*pb.GetServersResponse, error) {
// 	servers, err := s.GetServerer.GetServers()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &pb.GetServersResponse{Servers: servers}, nil
// }

func (s *Server) Put(ctx context.Context, r *pb.PutRequest) (*pb.PutResponse, error) {

	err := s.Services.Setter.Set(ctx, r.Records.Key, []byte(r.Records.Value))
	if err == nil {
		s.LoggerFacade.WriteSet(string(r.Records.Key), string(r.Records.Value))
	}

	return &pb.PutResponse{}, err
}

func (s *Server) Get(ctx context.Context, r *pb.GetRequest) (*pb.GetResponse, error) {

	value, err := s.Services.Getter.Get(ctx, r.Key)
	if err != nil {
		if err.Error() == "no such key" {
			// return &pb.GetResponse{Value: value.(string)}, status.Error(404, "and now")
			return &pb.GetResponse{Value: value.(string)}, pb.ErrorNoSuchKey{Key: r.Key}.RetErr()
		} else {
			return &pb.GetResponse{Value: value.(string)}, err
		}
	}

	// TODO if implement other types, the type assertion will have to be adapt
	// if value == "" grpc return it as nil
	return &pb.GetResponse{Value: value.(string)}, err
}

func (s *Server) GetKeys(ctx context.Context, r *pb.GetKeysRequest) (*pb.GetKeysResponse, error) {
	keys := s.Services.Getter.GetKeys(ctx)

	return &pb.GetKeysResponse{Keys: keys}, nil
}

func (s *Server) Delete(ctx context.Context, r *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	err := s.Services.Deleter.Delete(ctx, r.Key)
	if err == nil {
		s.LoggerFacade.WriteDelete(r.Key)
	}

	return &pb.DeleteResponse{}, err
}

func (s *Server) GetKeysValuesStream(r *pb.Empty, stream pb.KeyValue_GetKeysValuesStreamServer) error {
	// get keys
	ctx := context.Background()

	kv := make(chan models.KeysValues)

	for {
		select {
		case <-stream.Context().Done():
			return nil
		default:
			go func() {
				err := s.Services.Getter.GetKeysValues(ctx, kv)
				if err != nil {
					fmt.Println("see the GetKeysValuesStream:::::: ", err)
					// return err
				}
			}()

			for v := range kv {
				if err := stream.Send(&pb.GetRecords{
					Records: &pb.Records{
						Key:   v.Key,
						Value: v.Value,
					},
				}); err != nil {
					return err
				}
			}
			return nil
		}
	}
}
