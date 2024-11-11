package main

import (
	"context"
	"fmt"
	"log"
	"log-service/data"
	"log-service/logs"
	"net"

	"google.golang.org/grpc"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	err := l.Models.LogEntry.Insert(logEntry)

	if err != nil {
		response := &logs.LogResponse{Result: "failed"}

		return response, err
	}

	response := &logs.LogResponse{
		Result: "logged",
	}

	return response, nil
}

func (app *Config) gRPCListen() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))

	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	server := grpc.NewServer()

	logs.RegisterLogServiceServer(server, &LogServer{Models: app.Models})

	log.Printf("gRPC Server started on port %s", gRpcPort)

	if err = server.Serve(listener); err != nil {
		log.Fatalf("Failed to start server for gRPC: %v", err)
	}
}
