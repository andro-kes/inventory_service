package repo

import (
	"context"
	"time"

	"github.com/andro-kes/inventory_service/internal/repo/builder"
	pb "github.com/andro-kes/inventory_service/proto"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type ProductRepo interface {
	Create(ctx context.Context, p *pb.Product) (*pb.Product, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, prevSize, pageSize int32, filter, orderBy string) ([]*pb.Product, error)
	Update(ctx context.Context, p *pb.Product, mask *fieldmaskpb.FieldMask) (*pb.Product, error)
	Get(ctx context.Context, id string) (*pb.Product, error)
}

type productRepo struct {
	Pool *pgxpool.Pool
}

func NewProductRepo(ctx context.Context, pool *pgxpool.Pool) ProductRepo {
	return &productRepo{
		Pool: pool,
	}
}

func (pr *productRepo) Create(ctx context.Context, p *pb.Product) (*pb.Product, error) {
	sql, args := builder.NewSQLBuilder().
	Insert("products").
	Columns("id", "description", "price", "quantity", "tags", "available", "created_at", "updated_at").
	Values(p.Id, p.Description, p.Price, p.Quantity, p.Tags, p.Available, time.Now(), time.Now()).
	Returning("id", "description", "price", "quantity", "tags", "available", "created_at", "updated_at").
	Build()

	tx, err := pr.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var created *pb.Product
	row := tx.QueryRow(ctx, sql, args...)
	err = row.Scan(&created)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return created, nil
}

func (pr *productRepo) Delete(ctx context.Context, id string) error {
	sql, args := builder.NewSQLBuilder().
		Delete().From("products").Where("id = ?", id).Build()

	tx, err := pr.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (pr *productRepo) List(ctx context.Context, prevSize, pageSize int32, filter, orderBy string) ([]*pb.Product, error) {
	ob := "created_at DESC"
	switch orderBy {
	case "price", "price DESC", "price ASC",
		"created_at", "created_at DESC", "created_at ASC":
		ob = orderBy
	}

	b := builder.NewSQLBuilder().
		Select("id", "description", "price", "quantity", "tags", "avaible", "created_at", "updated_at").
		From("products").
		Where("quantity > 0").
		Where("avaible = true").
		OrderBy(ob).
		Offset(int(prevSize)).
		Limit(int(pageSize))

	if filter != "" {
		b.Where("tags @> ARRAY[?]::text[]", filter)
	}

	sql, args := b.Build()

	rows, err := pr.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]*pb.Product, 0, pageSize)
	for rows.Next() {
		var p pb.Product
		if err := rows.Scan(
			&p.Id, &p.Description, &p.Price, &p.Quantity,
			&p.Tags, &p.Available, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		products = append(products, &p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (pr *productRepo) Update(ctx context.Context, p *pb.Product, mask *fieldmaskpb.FieldMask) (*pb.Product, error) {
    b := builder.NewSQLBuilder().
        Update("products").
        Where("id = ?", p.GetId()).
        Returning("id", "description", "price", "quantity", "tags", "avaible", "created_at", "updated_at")

    for _, path := range mask.GetPaths() {
        switch path {
        case "description":
            b.Set("description = ?", p.GetDescription())
        case "price":
            b.Set("price = ?", p.GetPrice())
        case "quantity":
            b.Set("quantity = ?", p.GetQuantity())
        case "tags":
            b.Set("tags = ?", p.GetTags())
        case "avaible":
            b.Set("avaible = ?", p.GetAvailable())
        default:
            return nil, status.Errorf(codes.InvalidArgument, "unknown field in update_mask: %s", path)
        }
    }

	b.Set("updated_at = ?", time.Now())
    sql, args := b.Build()
    row := pr.Pool.QueryRow(ctx, sql, args...)

    var res pb.Product
    if err := row.Scan(
        &res.Id, &res.Description, &res.Price, &res.Quantity,
        &res.Tags, &res.Available, &res.CreatedAt, &res.UpdatedAt,
    ); err != nil {
        return nil, status.Errorf(codes.Internal, "update failed: %v", err)
    }

    return &res, nil
}

func (pr *productRepo) Get(ctx context.Context, id string) (*pb.Product, error) {
	sql, args := builder.NewSQLBuilder(). 
	Select("id", "description", "price", "quantity", "tags", "avaible", "created_at", "updated_at"). 
	From("products"). 
	Where("id = ?", id).Build()

	var p pb.Product
	err := pr.Pool.QueryRow(ctx, sql, args...).Scan(
		&p.Id, &p.Description, &p.Price, &p.Quantity,
		&p.Tags, &p.Available, &p.CreatedAt, &p.UpdatedAt,
	)

	return &p, err
}
