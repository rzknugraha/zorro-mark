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
	"github.com/rzknugraha/zorro-mark/infrastructures"
	"github.com/rzknugraha/zorro-mark/models"
	"github.com/rzknugraha/zorro-mark/repositories"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// IUploadService is
type IUploadService interface {
	StoreFile(ctx context.Context, file multipart.File, oldName string, IDUser int) (Response *helpers.JSONResponse, err error)
}

// UploadService is
type UploadService struct {
	DocumentRepository     repositories.IDocumentsRepository
	DocumentUserRepository repositories.IDocumentUserRepository
	DB                     infrastructures.ISQLConnection
}

// InitUploadService init
func InitUploadService() *UploadService {
	documentRepositories := new(repositories.DocumentsRepository)
	documentRepositories.DB = &infrastructures.SQLConnection{}

	documentUserRepositories := new(repositories.DocumentUserRepository)
	documentUserRepositories.DB = &infrastructures.SQLConnection{}

	UploadService := new(UploadService)
	UploadService.DocumentRepository = documentRepositories
	UploadService.DocumentUserRepository = documentUserRepositories

	return UploadService
}

// StoreFile is
func (s *UploadService) StoreFile(ctx context.Context, file multipart.File, oldName string, IDUser int) (Response *helpers.JSONResponse, err error) {

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

	tx, err := s.DocumentRepository.Tx()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  fileResp,
		}).Error("[Service StoreFile] error create tx")
		return
	}
	fmt.Println("siniini")
	defer tx.RollbackUnlessCommitted()
	TimeNow := time.Now()
	dataDoc := models.Documents{
		CreatedBy: IDUser,
		FileName:  oldName,
		Path:      fullPath,
		Status:    1,
		CreatedAt: TimeNow.Format("2006-01-02 15:04:05"),
	}

	count, err := s.DocumentRepository.CountNameByUserID(ctx, dataDoc)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  fileResp,
		}).Error("[Service StoreFile] error count filename doc")
		return
	}
	if count > 0 {
		dataDoc.FileName = fmt.Sprintf("%s (%v)", dataDoc.FileName, count)
	}
	idDoc, err := s.DocumentRepository.StoreDocuments(ctx, tx, dataDoc)
	if err != nil || idDoc == 0 {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  fileResp,
		}).Error("[Service StoreFile] error store doc")
		return
	}

	docUser := models.DocumentUser{
		DocumentID: int(idDoc),
		UserID:     IDUser,
		Status:     1,
		CreatedAt:  TimeNow.Format("2006-01-02 15:04:05"),
	}

	idStore, err := s.DocumentUserRepository.StoreDocumentUser(ctx, tx, docUser)
	if err != nil || idStore == 0 {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  fileResp,
		}).Error("[Service StoreFile] error store doc user")
		return
	}

	tx.Commit()

	return &helpers.JSONResponse{
		Code:    2200,
		Message: "Success",
		Data:    fileResp,
	}, nil
}
