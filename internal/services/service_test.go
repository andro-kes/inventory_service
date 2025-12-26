package services

import (
	"context"
	"testing"
	"time"

	pb "github.com/andro-kes/inventory_service/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TestRepo struct {
	Storage map[string]any
	Err error
}

func (r *TestRepo) Create(ctx context.Context, p *pb.Product) (*pb.Product, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	r.Storage[p.Id] = p
	return p, nil
}

func (r *TestRepo) Delete(ctx context.Context, id string) error {
	if r.Err != nil {
		return r.Err
	}
	
	if _, ok := r.Storage[id]; ok {
		delete(r.Storage, id)
		return nil
	} else {
		return assert.AnError
	}
}

func (r *TestRepo) Get(ctx context.Context, id string) (*pb.Product, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	
	if _, ok := r.Storage[id]; !ok {
		return nil, assert.AnError
	}

	return r.Storage[id].(*pb.Product), nil
}

func (r *TestRepo) List(ctx context.Context, prevSize, pageSize int32, filter, orderBy string) ([]*pb.Product, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	if len(r.Storage) == 0 {
		return nil, assert.AnError
	}

	p := make([]*pb.Product, 0, len(r.Storage))
	for _, v := range r.Storage {
		p = append(p, v.(*pb.Product))
	}

	return p, nil
}

func (r *TestRepo) Update(ctx context.Context, p *pb.Product, mask *fieldmaskpb.FieldMask) (*pb.Product, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	if _, ok := r.Storage[p.Id]; ok {
		switch {
		case mask == nil:
			r.Storage[p.Id] = p
		case mask.Paths[0] == "name":
			r.Storage[p.Id].(*pb.Product).Name = p.Name
		}
	}

	return p, nil
}

func NewTestService(err error) *ProductService {
	repo := &TestRepo{
		Storage: make(map[string]any),
		Err: nil,
	}

	return &ProductService{
		Repo: repo,
	}
}

var testProduct = &pb.Product{
	Id:          "1",
	Name:        "test",
	Description: "test",
	Price:       100,
	Quantity:    2,
	Tags:        []string{"test"},
	Available:   true,
	CreatedAt:   timestamppb.New(time.Now()),
	UpdatedAt:   timestamppb.New(time.Now()),
}

func TestCreate(t *testing.T) {
	service := NewTestService(nil)

	p, err := service.Create(t.Context(), testProduct)
	assert.NoError(t, err)

	assert.Equal(t, p.Id, testProduct.Id)
	assert.Equal(t, p.Name, testProduct.Name)
	assert.Equal(t, p.Description, testProduct.Description)
	assert.Equal(t, p.Price, testProduct.Price)
	assert.Equal(t, p.Quantity, testProduct.Quantity)
}