package cmd

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/proullon/ramsql/engine/log"

	// mysql driver
	_ "github.com/proullon/ramsql/driver"

	"go.smartmachine.io/go-grpc-api/pkg/protocol/grpc"
	"go.smartmachine.io/go-grpc-api/pkg/service/v1"
)

// Config is configuration for Server
type Config struct {
	// gRPC server start parameters section
	// gRPC is TCP port to listen by gRPC server
	GRPCPort string
}

// RunServer runs gRPC server and HTTP gateway
func RunServer() error {
	ctx := context.Background()

	// get configuration
	var cfg Config
	flag.StringVar(&cfg.GRPCPort, "grpc-port", "", "gRPC port to bind")

	flag.Parse()

	if len(cfg.GRPCPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server: '%s'", cfg.GRPCPort)
	}

	db, err := sql.Open("ramsql", "ToDoSchema")
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	// setup schema
	res, err := db.Exec(`
CREATE TABLE ToDo (
  ID BIGSERIAL PRIMARY KEY,
  Title varchar(200),
  Description varchar(1024),
  Reminder timestamp
);
`)

	if err != nil {
		return fmt.Errorf("failed to create schema: %v", err)
	}

	log.Info("created schema: %v", res)


	v1API := v1.NewToDoServiceServer(db)

	return grpc.RunServer(ctx, v1API, cfg.GRPCPort)
}