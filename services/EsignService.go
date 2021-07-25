package services

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

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
func (s *EsignService) PostSign(ctx context.Context, dataSign models.EsignReq) (Response *helpers.JSONResponse, err error) {

	//prepare the reader instances to encode
	values := map[string]io.Reader{
		"file":       mustOpen(dataSign.FilePath), // lets assume its this file
		"nik":        strings.NewReader(dataSign.NIK),
		"passphrase": strings.NewReader(dataSign.Passphrase),
		"image":      strings.NewReader(dataSign.Image),
		"linkQR":     strings.NewReader(dataSign.LinkQR),
		"tampilan":   strings.NewReader(dataSign.Tampilan),
	}

	fmt.Println("values")
	fmt.Println(values)

	err = s.EsignRepository.PostEsign(ctx, values)
	if err != nil {
		return nil, err
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
