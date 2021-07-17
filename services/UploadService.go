package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"github.com/rzknugraha/zorro-mark/helpers"
	"github.com/rzknugraha/zorro-mark/models"
	"github.com/spf13/viper"
)

// IUploadService is
type IUploadService interface {
	StoreFile(ctx context.Context, file multipart.File, oldName string) (Response *helpers.JSONResponse, err error)
}

// UploadService is
type UploadService struct {
}

// InitUploadService init
func InitUploadService() *UploadService {

	UploadService := new(UploadService)

	return UploadService
}

// StoreFile is
func (s *UploadService) StoreFile(ctx context.Context, file multipart.File, oldName string) (Response *helpers.JSONResponse, err error) {

	trimSpace := strings.ReplaceAll(oldName, " ", "")
	path := viper.GetString("storage.path")

	fileResp := models.UploadResp{
		FileName: fmt.Sprintf("%d-%s", time.Now().UnixNano(), trimSpace),
	}

	fullPath := fmt.Sprintf("/%s/%s", path, fileResp.FileName)

	dst, err := os.Create("." + fullPath)
	if err != nil {
		return

	}

	defer dst.Close()

	// Copy the uploaded file to the filesystem
	// at the specified destination
	_, err = io.Copy(dst, file)
	if err != nil {
		return
	}

	return &helpers.JSONResponse{
		Code:    2200,
		Message: "Success",
		Data:    fileResp,
	}, nil
}
