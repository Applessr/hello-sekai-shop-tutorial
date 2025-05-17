package grpccon

import (
	"errors"
	"log"
	"net"

	"github.com/Applessr/hello-sekai-shop-tutorial/config"
	authPb "github.com/Applessr/hello-sekai-shop-tutorial/modules/auth/authPb"
	inventoryPb "github.com/Applessr/hello-sekai-shop-tutorial/modules/inventory/inventoryPb"
	itemPb "github.com/Applessr/hello-sekai-shop-tutorial/modules/item/itemPb"
	playerPb "github.com/Applessr/hello-sekai-shop-tutorial/modules/player/playerPb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type (
	GrpcClientFactoryHandler interface {
		Auth() authPb.AuthGrpcServiceClient
		Inventory() inventoryPb.InventoryGrpcServiceClient
		Item() itemPb.ItemGrpcServiceClient
		Player() playerPb.PlayerGrpcServiceClient
	}

	grpcClientFactory struct {
		client *grpc.ClientConn
	}

	grpcAuth struct {
	}
)

func (g *grpcClientFactory) Auth() authPb.AuthGrpcServiceClient {
	return authPb.NewAuthGrpcServiceClient(g.client)
}

func (g *grpcClientFactory) Inventory() inventoryPb.InventoryGrpcServiceClient {
	return inventoryPb.NewInventoryGrpcServiceClient(g.client)
}

func (g *grpcClientFactory) Item() itemPb.ItemGrpcServiceClient {
	return itemPb.NewItemGrpcServiceClient(g.client)
}

func (g *grpcClientFactory) Player() playerPb.PlayerGrpcServiceClient {
	return playerPb.NewPlayerGrpcServiceClient(g.client)
}

func NewGrpcClient(host string) (GrpcClientFactoryHandler, error) {
	opts := make([]grpc.DialOption, 0)

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	clientConn, err := grpc.Dial(host, opts...)
	if err != nil {
		log.Printf("Error:grpc client connection failed %s", err)
		return nil, errors.New("error: grpc client connection failed")
	}

	return &grpcClientFactory{
		client: clientConn,
	}, nil
}

func NewGrpcServer(cfg *config.Jwt, host string) (*grpc.Server, net.Listener) {
	opts := make([]grpc.ServerOption, 0)

	grpcServer := grpc.NewServer(opts...)

	lis, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	return grpcServer, lis
}
