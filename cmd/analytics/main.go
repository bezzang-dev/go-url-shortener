package main

import (
	"context"
	"log"
	"net"

	"github.com/bezzang-dev/go-url-shortener/internal/repository"
	"github.com/bezzang-dev/go-url-shortener/internal/service"
	"github.com/bezzang-dev/go-url-shortener/internal/domain"
	pb "github.com/bezzang-dev/go-url-shortener/proto/analytics/v1"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type server struct {
	pb.UnimplementedAnalyticsServiceServer
	service *service.LogService
}

func (s *server) LogAccess(ctx context.Context, req *pb.AccessLogRequest) (*pb.AccessLogResponse, error) {
	log.Printf("ShortCode: %s | IP: %s | Time: %d", req.ShortCode, req.ClientIp, req.AccessedAt)
	err := s.service.RecordAccess(ctx, req.ShortCode, req.ClientIp, req.UserAgent, req.AccessedAt)
	if err != nil {
		return &pb.AccessLogResponse{Success: false}, err
	}
	return &pb.AccessLogResponse{Success: true}, nil
}

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=shortener port=5432 sslmode=disable TimeZone=Asia/Seoul"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := db.AutoMigrate(&domain.AccessLog{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	repo := repository.NewLogRepository(db)
	logService := service.NewLogService(repo)
	srv := &server{service: logService}
	grpcServer := grpc.NewServer()

	pb.RegisterAnalyticsServiceServer(grpcServer, srv)

	log.Printf("Analytics gRPC Server listening at %v", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}