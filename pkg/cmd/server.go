package cmd

import (
	"context"
	"flag"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/proullon/ramsql/engine/log"
	v12 "go.smartmachine.io/go-grpc-api/pkg/api/v1"
	"go.smartmachine.io/go-grpc-api/pkg/protocol/rest"

	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"go.smartmachine.io/go-grpc-api/pkg/protocol/grpc"
	"go.smartmachine.io/go-grpc-api/pkg/service/v1"
)

// Config is configuration for Server
type Config struct {
	// gRPC server start parameters section
	// gRPC is TCP port to listen by gRPC server
	GRPCPort string
	HTTPPort string
}

// RunServer runs gRPC server and HTTP gateway
func RunServer() error {
	ctx := context.Background()

	// get configuration
	var cfg Config
	flag.StringVar(&cfg.GRPCPort, "grpc-port", "1234", "gRPC port to bind")
	flag.StringVar(&cfg.HTTPPort, "http-port", "8080", "HTTP port to bind")

	flag.Parse()

	if len(cfg.GRPCPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server: '%s'", cfg.GRPCPort)
	}
	if len(cfg.HTTPPort) == 0 {
		return fmt.Errorf("invalid TCP port for HTTP gateway: '%s'", cfg.HTTPPort)
	}

	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	db.AutoMigrate(&v12.ToDoORM{})


	if err != nil {
		return fmt.Errorf("failed to create schema: %v", err)
	}

	log.Info("created schema: ToDoORM")


	v1API := v1.NewToDoServiceServer(db)

	// run HTTP gateway
	go func() {
		_ = rest.RunServer(ctx, cfg.GRPCPort, cfg.HTTPPort)
	}()


	return grpc.RunServer(ctx, v1API, cfg.GRPCPort)
}