package main

import (
	"ImageGrpc/config"
	"ImageGrpc/internal/service"
	pb "ImageGrpc/proto"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Ошибка при прослушивании порта:%v", err)
	}
	defer listener.Close()
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			config.RateLimitingInterceptor(config.UploadDownloadLimiter),
			config.RateLimitingInterceptor(config.ListFilesLimiter),
		))
	fileServer := service.NewFileServer()
	pb.RegisterFileServiceServer(server, fileServer)
	log.Println("Запуск gRPC сервера...")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
