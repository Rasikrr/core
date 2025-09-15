package orders

import (
	"context"

	"github.com/Rasikrr/core/log"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)
import pb "github.com/Rasikrr/core/pkg/api/proto/orders"

type server struct {
	pb.UnimplementedOrdersServer
}

func NewServer(grpcServer *grpc.Server) {
	pb.RegisterOrdersServer(grpcServer, newGrpcServer())
}

func newGrpcServer() *server {
	return &server{}
}

func (s *server) CreateOrder(ctx context.Context, in *pb.Order) (*emptypb.Empty, error) {
	log.Infof(ctx, "CreateOrder GRPC called with %v", in)
	return &emptypb.Empty{}, nil
}
