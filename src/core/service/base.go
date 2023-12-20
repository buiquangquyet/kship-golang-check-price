package service

import (
	"check-price/src/common"
	"check-price/src/common/log"
	"context"
	"encoding/json"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type baseService struct {
}

func NewBaseService() *baseService {
	return &baseService{}
}

func (b *baseService) connectGrpc(domain string) *grpc.ClientConn {
	conn, err := grpc.Dial(domain,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(
			grpc_middleware.ChainUnaryClient(
				otelgrpc.UnaryClientInterceptor(),
			),
		),
	)
	if err != nil {
		log.GetLogger().GetZap().Fatalf("connect grpc error, domain:[%s], err:[%s]", domain, err.Error())
	}
	return conn
}

func (b *baseService) grpToIError(ctx context.Context, inputErr error) *common.Error {
	var ierr common.Error
	grpcErr, ok := status.FromError(inputErr)
	if !ok {
		return common.ErrSystemError(ctx, fmt.Sprintf("grpc error convert failed, err:[%s]", inputErr.Error()))
	}

	err := json.Unmarshal([]byte(grpcErr.Message()), &ierr)
	if err != nil {
		return common.ErrSystemError(ctx, fmt.Sprintf("grpc error unmarshal failed with input [%s]", inputErr.Error()))
	}

	return &ierr
}
