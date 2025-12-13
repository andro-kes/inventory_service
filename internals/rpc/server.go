package rpc

import (
	"context"

	pb "github.com/andro-kes/inventory_service/proto"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InventoryService struct {
	pb.UnimplementedInventoryServiceServer
}

func NewInventoryService(ctx context.Context, pool *pgxpool.Pool) *InventoryService{
	return &InventoryService{}
}

func (is *InventoryService) CreateProduct(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {

}

func (is *InventoryService) DeleteProduct(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {

}

func (is *InventoryService) ListProducts(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {

}

func (is *InventoryService) UpdateProduct(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {

}

func (is *InventoryService) GetProduct(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {

}