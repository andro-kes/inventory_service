package services

import (
	"context"

	"github.com/andro-kes/inventory_service/internal/repo"
	pb "github.com/andro-kes/inventory_service/proto"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type ProductService struct {
	Repo repo.ProductRepo
}

func NewProductService(ctx context.Context, pool *pgxpool.Pool) *ProductService {
	return &ProductService{
		Repo: repo.NewProductRepo(ctx, pool),
	}
}

func (ps *ProductService) Create(ctx context.Context, p *pb.Product) (*pb.Product, error) {
	id := uuid.NewString()
	p.Id = id

	return ps.Repo.Create(ctx, p)
}

func (ps *ProductService) Delete(ctx context.Context, id string) error {
	return ps.Repo.Delete(ctx, id)
}

func (ps *ProductService) List(ctx context.Context, prevSize, pageSize int32, filter, orderBy string) ([]*pb.Product, error) {
	return ps.Repo.List(ctx, prevSize, pageSize, filter, orderBy)
}

func (ps *ProductService) Update(ctx context.Context, p *pb.Product, mask *fieldmaskpb.FieldMask) (*pb.Product, error) {
	return ps.Repo.Update(ctx, p, mask)
}

func (ps *ProductService) Get(ctx context.Context, id string) (*pb.Product, error) {
	return ps.Repo.Get(ctx, id)
}
