package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rzknugraha/zorro-mark/helpers"
	"github.com/rzknugraha/zorro-mark/infrastructures"
	"github.com/rzknugraha/zorro-mark/models"
	"github.com/rzknugraha/zorro-mark/repositories"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// IEsignService is
type IEsignService interface {
	PostSign(ctx context.Context, dataSign models.EsignReq, dataUser models.Shortuser) (Response *helpers.JSONResponse, err error)
	PostSignMultiple(ctx context.Context, dataSign models.EsignMutipleReq, dataUser models.Shortuser) (response *helpers.JSONResponse, err error)
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

		fmt.Println("condition")
		fmt.Println(condition)
		fmt.Println("payload")
		fmt.Println(payload)

		_, err1 = s.DocumentRepository.UpdateDoc(ctx, tx, condition, payload)
		if err1 != nil {
			logrus.WithFields(logrus.Fields{
				"code":  5500,
				"error": err,
				"data":  nil,
			}).Error("[Service PostSign] error update doc signed")
			return nil, err1
		}

		fmt.Println("selesai update doc")

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

// PostSignMultiple is
func (s *EsignService) PostSignMultiple(ctx context.Context, dataSign models.EsignMutipleReq, dataUser models.Shortuser) (response *helpers.JSONResponse, err error) {

	//check passpharse

	esignCheck := models.EsignReq{

		NIK:        dataUser.IdentityNO,
		Tampilan:   "invisible",
		Passphrase: dataSign.Passphrase,
		Image:      false,
		FilePath:   viper.GetString("esign.dummy"),
	}

	rescheck, err := s.EsignRepository.PostEsign(ctx, esignCheck)
	if err != nil {
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"code":  5500,
				"error": err,
				"data":  nil,
			}).Error("[Service PostSignMultiple] error Post dummy pdf")
			return nil, err
		}
	}

	if rescheck.StatusCode != 200 {
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"code":  5500,
				"error": err,
				"data":  nil,
			}).Error("[Service PostSignMultiple] error Post dummy pdf not 200")
			return nil, err
		}
	}

	docIDs := strings.Split(dataSign.DocumentID, ",")

	singleSign := models.EsignReq{
		NIK:        dataUser.IdentityNO,
		Passphrase: dataSign.Passphrase,
	}

	fatalErrors := make(chan error)
	wgDone := make(chan bool)

	newctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	for _, docID := range docIDs {
		wg.Add(1)
		go func(docIDInside string) {
			defer wg.Done()

			select {
			case <-newctx.Done():
				return // Error somewhere, terminate
			default: // Default is must to avoid blocking
			}
			intDocID, err := strconv.Atoi(docIDInside)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"code":  5500,
					"error": err,
					"data":  docIDInside,
				}).Error("[Service PostSignMultiple] error convert to int docID")
				fatalErrors <- err
				cancel()
				return
			}
			getDoc, err := s.DocumentUserRepository.GetSingleDocByUser(ctx, dataUser.ID, intDocID)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"code":  5500,
					"error": err,
					"data":  intDocID,
				}).Error("[Service PostSignMultiple] error get document by get ID")
				fatalErrors <- err
				cancel()
				return
			}

			singleSign.DocumentID = intDocID
			singleSign.FilePath = getDoc.Path
			singleSign.Page = getDoc.Page
			singleSign.Tampilan = getDoc.Tampilan.ValueOrZero()
			singleSign.ImagePath = dataUser.SignFile
			singleSign.Height = getDoc.Height
			singleSign.Width = getDoc.Width
			singleSign.XAxis = getDoc.XAxis
			singleSign.YAxis = getDoc.YAxis

			result, err := s.PostSign(ctx, singleSign, dataUser)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"code":  5500,
					"error": err,
					"data":  singleSign,
				}).Error("[Service PostSignMultiple] error when multiple sign")
				fatalErrors <- err
				cancel()
				return

			}

			if result == nil {
				err1 := errors.New("not 2200")
				logrus.WithFields(logrus.Fields{
					"code":  5500,
					"error": err1,
					"data":  result,
				}).Error("[Service PostSignMultiple] error not 2200")
				fatalErrors <- err1
				cancel()
				return

			}

			fmt.Println(result)
		}(docID)
	}

	go func() {
		wg.Wait()
		close(wgDone)
	}()

	select {
	case <-wgDone:
		// carry on
		break
	case err := <-fatalErrors:
		// close(fatalErrors)
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  singleSign,
		}).Error("[Service PostSignMultiple] error when multiple sign")
		response := &helpers.JSONResponse{
			Code:    5500,
			Message: err.Error(),
			Data:    nil,
		}

		return response, err
	}

	fmt.Println("ga masuk error")
	response = &helpers.JSONResponse{
		Code:    2200,
		Message: "success",
		Data:    nil,
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
