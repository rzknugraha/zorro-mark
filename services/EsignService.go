package services

import (
	"context"
	"net/http"
	"os"
	"time"

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

	var actvity models.DocumentActivity

	actvity.UserID = dataUser.ID
	actvity.DocumentID = dataSign.DocumentID
	actvity.Name = dataUser.Name
	actvity.NIP = dataUser.Nip

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

	res, err := s.EsignRepository.PostEsign(ctx, dataSign)
	if err != nil {

		actvity.Status = 0
		actvity.Message = err.Error()
		actvity.Type = "error-signed"

		_, err = s.DocumentActivityRepository.StoreDocumentActivity(ctx, tx, actvity)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"code":  5500,
				"error": err,
				"data":  nil,
			}).Error("[Service PostSign] error StoreDocumentActivity")
			return nil, err
		}
		tx.Commit()

		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		response = &helpers.JSONResponse{
			Code:    4400,
			Message: res.Error,
			Data:    nil,
		}

		actvity.Status = 0
		actvity.Message = res.Error
		actvity.Type = "error-signed"

		_, err = s.DocumentActivityRepository.StoreDocumentActivity(ctx, tx, actvity)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"code":  5500,
				"error": err,
				"data":  nil,
			}).Error("[Service PostSign] error StoreDocumentActivity")
			return nil, err
		}
		tx.Commit()
	} else {
		response = &helpers.JSONResponse{
			Code:    2200,
			Message: "success",
			Data:    nil,
		}

		TimeNow := time.Now()

		condition := map[string]interface{}{
			"id": dataSign.DocumentID,
		}
		payload := map[string]interface{}{
			"signed":     1,
			"path":       res.PathFile,
			"updated_at": TimeNow.Format("2006-01-02 15:04:05"),
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

		actvity.Status = 1
		actvity.Message = "Document has been signed"
		actvity.Type = "signed"

		_, err = s.DocumentActivityRepository.StoreDocumentActivity(ctx, tx, actvity)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"code":  5500,
				"error": err,
				"data":  nil,
			}).Error("[Service PostSign] error StoreDocumentActivity")
			return nil, err
		}

		condition1 := map[string]interface{}{
			"document_id": dataSign.DocumentID,
		}
		payload1 := map[string]interface{}{
			"status":     1,
			"updated_at": TimeNow.Format("2006-01-02 15:04:05"),
		}

		_, err = s.DocumentUserRepository.UpdateDocUsers(ctx, tx, condition1, payload1)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"code":  5500,
				"error": err,
				"data":  nil,
			}).Error("[Service PostSign] error UpdateDocUsers")
			return nil, err
		}
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
