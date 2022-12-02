package grpc

import (
	"context"

	pb "github.com/djedjethai/generation/api/v1/keyvalue"
	"github.com/djedjethai/generation/pkg/config"
	"github.com/djedjethai/generation/pkg/logger"
)

type Server struct {
	pb.UnimplementedKeyValueServer
	Services     config.Services
	LoggerFacade *logger.LoggerFacade
}

func (s *Server) Put(ctx context.Context, r *pb.PutRequest) (*pb.PutResponse, error) {
	// TODO

	err := s.Services.Setter.Set(ctx, r.Key, []byte(r.Value))
	if err == nil {
		s.LoggerFacade.WriteSet(string(r.Key), string(r.Value))
	}

	return &pb.PutResponse{}, err
}

func (s *Server) Get(ctx context.Context, r *pb.GetRequest) (*pb.GetResponse, error) {

	value, err := s.Services.Getter.Get(ctx, r.Key)

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
