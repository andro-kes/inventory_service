package inverr

import (
	"google.golang.org/grpc/codes"
)

type InvError struct {
	msg string
	grpcCode codes.Code
}

func New(msg string, code codes.Code) *InvError {
	return &InvError{
		msg: msg,
		grpcCode: code,
	}
}

func (ie *InvError) Error() string {
	if ie == nil {
		return "<nil>"
	}
	return ie.msg
}

var (
	InvalidPoolConfig = New("failed to parse config", codes.Internal)
	CreatePoolError = New("failed to create pool", codes.Internal)

	CreateProductError = New("failed to create product", codes.Internal)
	DeleteProductError = New("failed to delete product", codes.Internal)
	ListProductsError = New("failed to list product", codes.Internal)
)