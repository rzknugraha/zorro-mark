package services

import (
	"context"
	"net/http"
	"os"

	"github.com/rzknugraha/zorro-mark/helpers"
	"github.com/rzknugraha/zorro-mark/infrastructures"
	"github.com/rzknugraha/zorro-mark/models"
	"github.com/rzknugraha/zorro-mark/repositories"
	"github.com/sirupsen/logrus"
)

// IEsignService is
type IEsignService interface {
	PostSign(ctx context.Context, dataSign models.EsignReq, dataUser models.Shortuser) (Response *helpers.JSONResponse, err error)
}

// EsignService is
type EsignService struct {
	DocumentRepository         repositories.IDocumentsRepository
	DocumentUserRepository     repositories.IDocumentUserRepository
	DocumentActivityRepository repositories.IDocumentActivityRepository
	EsignRepository            repositories.IEsignRepository
	DB                         infrastructures.ISQLConnection
}

// InitEsignService init
func InitEsignService() *EsignService {
	documentRepositories := new(repositories.DocumentsRepository)
	documentRepositories.DB = &infrastructures.SQLConnection{}

	documentUserRepositories := new(repositories.DocumentUserRepository)
	documentUserRepositories.DB = &infrastructures.SQLConnection{}

	esignRepositories := new(repositories.EsignRepository)

	documentActivityRepositories := new(repositories.DocumentActivityRepository)
	documentActivityRepositories.DB = &infrastructures.SQLConnection{}

	EsignService := new(EsignService)
	EsignService.DocumentRepository = documentRepositories
	EsignService.DocumentUserRepository = documentUserRepositories
	EsignService.EsignRepository = esignRepositories
	EsignService.DocumentActivityRepository = documentActivityRepositories

	return EsignService
}

// PostSign is
func (s *EsignService) PostSign(ctx context.Context, dataSign models.EsignReq, dataUser models.Shortuser) (response *helpers.JSONResponse, err error) {

	res, err := s.EsignRepository.PostEsign(ctx, dataSign)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		response = &helpers.JSONResponse{
			Code:    4400,
			Message: res.Error,
			Data:    nil,
		}
	} else {
		response = &helpers.JSONResponse{
			Code:    2200,
			Message: "success",
			Data:    nil,
		}

		tx, err1 := s.DocumentRepository.Tx()
		if err1 != nil {
			logrus.WithFields(logrus.Fields{
				"code":  5500,
				"error": err,
				"data":  nil,
			}).Error("[Service PostSign] error init tx")
			return nil, err1
		}
		defer tx.RollbackUnlessCommitted()

		condition := map[string]interface{}{
			"id": dataSign.DocumentID,
		}
		payload := map[string]interface{}{
			"signed": 1,
			"path":   res.PathFile,
		}

		_, err1 = s.DocumentRepository.UpdateDoc(ctx, tx, condition, payload)
		if err1 != nil {
			logrus.WithFields(logrus.Fields{
				"code":  5500,
				"error": err,
				"data":  nil,
			}).Error("[Service PostSign] error update doc signed")
			return nil, err1
		}

		var actvity models.DocumentActivity

		actvity.UserID = dataUser.ID
		actvity.DocumentID = dataSign.DocumentID
		actvity.Name = dataUser.Name
		actvity.NIP = dataUser.Nip
		actvity.Status = 1
		actvity.Message = "Document has been signed"
		actvity.Type = "signed"

		_, err = s.DocumentActivityRepository.StoreDocumentActivity(ctx, tx, actvity)

		tx.Commit()

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
