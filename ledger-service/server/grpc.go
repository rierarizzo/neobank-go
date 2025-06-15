package server

import (
	"context"
	"github.com/rierarizzo/neobank-go/ledger-service/domain"
	pb "github.com/rierarizzo/neobank-go/ledger-service/proto/ledger"
	"github.com/rierarizzo/neobank-go/ledger-service/services"
)

type servo struct {
	pb.UnimplementedLedgerServer
	svc services.LedgerSvc
}

func NewGRPCServer(svc services.LedgerSvc) pb.LedgerServer {
	return &servo{svc: svc}
}

func (s *servo) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	currency := domain.Currency(req.GetCurrency())

	acct, err := s.svc.CreateAccount(ctx, currency)
	if err != nil {
		return nil, err
	}

	return &pb.CreateAccountResponse{
		Id:       acct.ID,
		Currency: string(acct.Currency),
		Status:   string(acct.Status),
	}, nil
}
