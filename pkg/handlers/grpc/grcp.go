package grpc

import (
	"context"

	pb "github.com/djedjethai/generation/api/v1/keyvalue"
	"github.com/djedjethai/generation/pkg/deleter"
	"github.com/djedjethai/generation/pkg/getter"
	"github.com/djedjethai/generation/pkg/logger"
	"github.com/djedjethai/generation/pkg/setter"
)

type Server struct {
	pb.UnimplementedKeyValueServer
	SetSrv       setter.Setter
	GetSrv       getter.Getter
	DelSrv       deleter.Deleter
	LoggerFacade *logger.LoggerFacade
}

func (s *Server) Put(ctx context.Context, r *pb.PutRequest) (*pb.PutResponse, error) {
	// TODO

	err := s.SetSrv.Set(ctx, r.Key, []byte(r.Value))
	if err == nil {
		s.LoggerFacade.WriteSet(string(r.Key), string(r.Value))
	}

	return &pb.PutResponse{}, err
}

func (s *Server) Get(ctx context.Context, r *pb.GetRequest) (*pb.GetResponse, error) {

	value, err := s.GetSrv.Get(ctx, r.Key)

	return &pb.GetResponse{Value: value.(string)}, err
}

func (s *Server) GetKeys(ctx context.Context, r *pb.GetKeysRequest) (*pb.GetKeysResponse, error) {
	keys := s.GetSrv.GetKeys(ctx)

	return &pb.GetKeysResponse{Keys: keys}, nil
}

func (s *Server) Delete(ctx context.Context, r *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	err := s.DelSrv.Delete(ctx, r.Key)
	if err == nil {
		s.LoggerFacade.WriteDelete(r.Key)
	}

	return &pb.DeleteResponse{}, err
}
