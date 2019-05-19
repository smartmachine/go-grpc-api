package cmd

import (
	"context"
	"flag"
	"fmt"
	"github.com/jinzhu/gorm"
	apiv1 "go.smartmachine.io/go-grpc-api/pkg/api/v1"
	"go.smartmachine.io/go-grpc-api/pkg/logger"
	"go.smartmachine.io/go-grpc-api/pkg/protocol/rest"

	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"go.smartmachine.io/go-grpc-api/pkg/protocol/grpc"
	servicev1 "go.smartmachine.io/go-grpc-api/pkg/service/v1"
)

// Config is configuration for Server
type Config struct {
	// gRPC server start parameters section
	// gRPC is TCP port to listen by gRPC server
	GRPCPort string
	HTTPPort string

	// Log parameters section
	// LogLevel is global log level: Debug(-1), Info(0), Warn(1), Error(2), DPanic(3), Panic(4), Fatal(5)
	LogLevel int
	// LogTimeFormat is print time format for logger e.g. 2006-01-02T15:04:05Z07:00
	LogTimeFormat string
}

// RunServer runs gRPC server and HTTP gateway
func RunServer() error {
	ctx := context.Background()

	// get configuration
	var cfg Config
	flag.StringVar(&cfg.GRPCPort, "grpc-port", "1234", "gRPC port to bind")
	flag.StringVar(&cfg.HTTPPort, "http-port", "8080", "HTTP port to bind")
	flag.IntVar(&cfg.LogLevel, "log-level", 0, "Global log level")
	flag.StringVar(&cfg.LogTimeFormat, "log-time-format", "2006-01-02T15:04:05.0000Z07:00",
		"Print time format for logger e.g. 2006-01-02T15:04:05Z07:00")
	flag.Parse()

	if len(cfg.GRPCPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server: '%s'", cfg.GRPCPort)
	}
	if len(cfg.HTTPPort) == 0 {
		return fmt.Errorf("invalid TCP port for HTTP gateway: '%s'", cfg.HTTPPort)
	}

	if err := logger.Init(cfg.LogLevel, cfg.LogTimeFormat); err != nil {
		return fmt.Errorf("failed to initialize logger: %v", err)
	}

	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	err = db.AutoMigrate(&apiv1.ToDoORM{}).Error
	if err != nil {
		return fmt.Errorf("failed to create schema: %v", err)
	}


	logger.Log.Info("created schema: ToDoORM")

	v1API := servicev1.NewToDoServiceServer(db)

	// run HTTP gateway
	go func() {
		_ = rest.RunServer(ctx, cfg.GRPCPort, cfg.HTTPPort)
	}()

	return grpc.RunServer(ctx, v1API, cfg.GRPCPort)
}
