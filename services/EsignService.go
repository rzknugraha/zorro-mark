package services

import (
	"context"
	"net/http"
	"os"

	"github.com/rzknugraha/zorro-mark/helpers"
	"github.com/rzknugraha/zorro-mark/infrastructures"
	"github.com/rzknugraha/zorro-mark/models"
	"github.com/rzknugraha/zorro-mark/repositories"
)

// IEsignService is
type IEsignService interface {
	PostSign(ctx context.Context, dataSign models.EsignReq) (Response *helpers.JSONResponse, err error)
}

// EsignService is
type EsignService struct {
	DocumentRepository     repositories.IDocumentsRepository
	DocumentUserRepository repositories.IDocumentUserRepository
	EsignRepository        repositories.IEsignRepository
	DB                     infrastructures.ISQLConnection
}

// InitEsignService init
func InitEsignService() *EsignService {
	documentRepositories := new(repositories.DocumentsRepository)
	documentRepositories.DB = &infrastructures.SQLConnection{}

	documentUserRepositories := new(repositories.DocumentUserRepository)
	documentUserRepositories.DB = &infrastructures.SQLConnection{}

	esignRepositories := new(repositories.EsignRepository)

	EsignService := new(EsignService)
	EsignService.DocumentRepository = documentRepositories
	EsignService.DocumentUserRepository = documentUserRepositories
	EsignService.EsignRepository = esignRepositories

	return EsignService
}

// PostSign is
func (s *EsignService) PostSign(ctx context.Context, dataSign models.EsignReq) (response *helpers.JSONResponse, err error) {

	res, err := s.EsignRepository.PostEsign(ctx, dataSign)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		response = &helpers.JSONResponse{
			Code:    4400,
			Message: res.Message,
			Data:    nil,
		}
	} else {
		response = &helpers.JSONResponse{
			Code:    2200,
			Message: res.Message,
			Data:    nil,
		}
	}
	return
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	return r
}
