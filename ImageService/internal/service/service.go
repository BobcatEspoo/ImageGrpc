package service

import (
	"ImageGrpc/config"
	pb "ImageGrpc/proto"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type FileServer struct {
	pb.UnimplementedFileServiceServer
	MainDB *sql.DB
}

func NewFileServer() *FileServer {
	db, err := config.AccessToDB()
	if err != nil {
		log.Fatal(err)
	}
	return &FileServer{MainDB: db}
}

func (f *FileServer) UploadFile(ctx context.Context, req *pb.UploadRequest) (*pb.UploadResponse, error) {
	query := "INSERT INTO files(file_name, file_data, uploaded_at, updated_at) VALUES ($1, $2, $3, $4)"
	_, err := f.MainDB.Exec(query, req.FileName, req.File, time.Now().Unix(), time.Now().Unix())
	if err != nil {
		return &pb.UploadResponse{Status: fmt.Sprint("have error: ", err)}, err
	}
	return &pb.UploadResponse{Status: "fine upload"}, nil
}
func (f *FileServer) GetFile(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	result := []byte{}
	query := "SELECT file_data FROM files WHERE file_name = $1"
	err := f.MainDB.QueryRow(query, req.FileName).Scan(&result)
	if err != nil {
		return nil, err
	}
	return &pb.GetResponse{File: result}, nil
}
func (f *FileServer) AllFIle(ctx context.Context, req *pb.AllRequest) (*pb.AllResponse, error) {
	result := &pb.AllResponse{}
	query := "SELECT file_name, uploaded_at, updated_at  FROM files"
	rows, err := f.MainDB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&result.FileName, &result.CreationDate, &result.UpdateDate)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}
func (f *FileServer) mustEmbedUnimplementedFileServiceServer() {}
