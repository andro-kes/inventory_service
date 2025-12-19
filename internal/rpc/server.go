package rpc

import (
	"context"

	"github.com/andro-kes/inventory_service/internal/inverr"
	"github.com/andro-kes/inventory_service/internal/services"
	pb "github.com/andro-kes/inventory_service/proto"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InventoryService struct {
	pb.UnimplementedInventoryServiceServer
	ProductService *services.ProductService
}

func NewInventoryService(ctx context.Context, pool *pgxpool.Pool) *InventoryService {
	return &InventoryService{
		ProductService: services.NewProductService(ctx, pool),
	}
}

func (is *InventoryService) CreateProduct(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	err := is.ProductService.Create(ctx, req.Product)
	if err != nil {
		return nil, inverr.CreateProductError
	}

	var resp *pb.CreateResponse
	resp.Product = req.Product
	return resp, nil
}

func (is *InventoryService) DeleteProduct(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	var resp *pb.DeleteResponse

	err := is.ProductService.Delete(ctx, req.Id)
	if err != nil {
		resp.Success = false
		return resp, inverr.DeleteProductError
	}

	resp.Success = true
	return resp, nil
}

func (is *InventoryService) ListProducts(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	var resp *pb.ListResponse

	products, err := is.ProductService.List(ctx, req.PageSize, req.Filter, req.OrderBy)
	if err != nil {
		return nil, inverr.ListProductsError
	}

	resp.Products = products
	return resp, nil
}

func (is *InventoryService) UpdateProduct(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	var resp *pb.UpdateResponse

	product, err := is.ProductService.Update(ctx, req.Product)
	if err != nil {
		return nil, err
	}

	resp.Product = product

	return resp, nil
}

func (is *InventoryService) GetProduct(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	var resp *pb.GetResponse

	product, err := is.ProductService.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	resp.Product = product

	return resp, nil
}
