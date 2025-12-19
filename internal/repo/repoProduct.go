package repo

import (
	"context"

	pb "github.com/andro-kes/inventory_service/proto"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepo interface {
	Create(ctx context.Context, p *pb.Product) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, pageSize int32, filter, orderBy string) ([]*pb.Product, error)
	Update(ctx context.Context, p *pb.Product) (*pb.Product, error)
	Get(ctx context.Context, id string) (*pb.Product, error)
}

type productRepo struct {
	Pool *pgxpool.Pool
}

func (pr *productRepo) Create(ctx context.Context, p *pb.Product) error

func (pr *productRepo) Delete(ctx context.Context, id string) error

func (pr *productRepo) List(ctx context.Context, pageSize int32, filter, orderBy string) ([]*pb.Product, error)

func (pr *productRepo) Update(ctx context.Context, p *pb.Product) (*pb.Product, error)

func (pr *productRepo) Get(ctx context.Context, id string) (*pb.Product, error)
