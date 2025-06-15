package main

import (
	"github.com/rierarizzo/neobank-go/ledger-service/config"
	pb "github.com/rierarizzo/neobank-go/ledger-service/proto/ledger"
	"github.com/rierarizzo/neobank-go/ledger-service/server"
	"github.com/rierarizzo/neobank-go/ledger-service/services"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

func main() {
	err := config.ConfigureZap()
	if err != nil {
		panic("Logger configuration failed: " + err.Error())
	}

	zap.L().Info("Initializing Ledger Service")

	db, err := config.ConnectToPostgres()
	if err != nil {
		zap.L().Error("DB connection failed", zap.Error(err))
		return
	}

	ledgerSvc := services.NewLedgerPostgresService(db)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		zap.L().Error("Failed to listen", zap.String("addr", ":50051"), zap.Error(err))
		return
	}

	grpcServer := grpc.NewServer()
	pb.RegisterLedgerServer(grpcServer, server.NewGRPCServer(ledgerSvc))
	reflection.Register(grpcServer)

	zap.L().Info("gRPC server listening on :50051")

	if err := grpcServer.Serve(lis); err != nil {
		zap.L().Fatal("gRPC Serve failed", zap.Error(err))
	}
}
