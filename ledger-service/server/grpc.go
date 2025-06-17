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

func (s *servo) PostTransfer(ctx context.Context, req *pb.PostTransferRequest) (*pb.Empty, error) {
	currency := domain.Currency(req.GetCurrency())

	err := s.svc.PostTransfer(ctx, req.GetExternalID(), req.GetFromAccountID(), req.GetToAccountID(), req.GetAmount(), currency)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *servo) GetBalanceCents(ctx context.Context, req *pb.GetBalanceCentsRequest) (*pb.GetBalanceCentsResponse, error) {
	balanceCents, err := s.svc.GetBalanceCents(ctx, req.GetAccountID())
	if err != nil {
		return nil, err
	}

	return &pb.GetBalanceCentsResponse{
		BalanceCents: balanceCents,
	}, nil
}
