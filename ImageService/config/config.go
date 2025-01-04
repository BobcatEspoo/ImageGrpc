package config

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
)

var (
	UploadDownloadLimiter = make(chan struct{}, 10)
	ListFilesLimiter      = make(chan struct{}, 100)
)

func RateLimitingInterceptor(limiter chan struct{}) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		select {
		case limiter <- struct{}{}:
			defer func() { <-limiter }()
			return handler(ctx, req)
		default:
			return nil, fmt.Errorf("rate limit exceeded, please try again later")
		}
	}
}

func AccessToDB() (*sql.DB, error) {
	dbURL := "postgres://user:password@db:5432/photo_service?sslmode=disable"
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("error in access to db: %v", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database is anavailable: %v", err)
	}
	log.Println("Fine access to db")
	return db, nil
}
