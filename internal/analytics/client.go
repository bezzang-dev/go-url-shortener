package analytics

import (
	"context"
	"log"
	"time"

	pb "github.com/bezzang-dev/go-url-shortener/proto/analytics/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	service pb.AnalyticsServiceClient
	conn *grpc.ClientConn
}

func NewClient(address string) (*Client, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	svc := pb.NewAnalyticsServiceClient(conn)

	return &Client{
		service: svc,
		conn: conn,
	}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) LogAccessAsync(shortCode, ip, userAgent string) {
	go func ()  {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := c.service.LogAccess(ctx, &pb.AccessLogRequest{
			ShortCode: shortCode,
			ClientIp: ip,
			UserAgent: userAgent,
			AccessedAt: time.Now().Unix(),
		})

		if err != nil {
			log.Printf("Analytics Server Error: %v", err)
		}
	}()
}