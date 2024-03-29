package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gocraft/dbr/v2"
	"github.com/rzknugraha/zorro-mark/helpers"
	"github.com/rzknugraha/zorro-mark/infrastructures"
	"github.com/rzknugraha/zorro-mark/models"
	"github.com/rzknugraha/zorro-mark/repositories"
	"github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v3"
)

// IDocumentService is
type IDocumentService interface {
	GetDocumentUser(ctx context.Context, filter models.DocumentUserFilter, page helpers.PageReq) (res *helpers.Paginate, err error)
	UpdateDocumentAttributte(ctx context.Context, filter models.UpdateDocReq, userData models.Shortuser) (res *helpers.JSONResponse, err error)
	GetSingleDocByUser(ctx context.Context, docID int, userData models.Shortuser) (res *helpers.JSONResponse, err error)
	GetActivityDoc(ctx context.Context, docID int, userData models.Shortuser) (res *helpers.JSONResponse, err error)
	SaveDraft(ctx context.Context, userData models.Shortuser, dataReq models.DocumentUser) (res *helpers.JSONResponse, err error)
	SendSign(ctx context.Context, userData models.Shortuser, dataReq models.DocumentUserSendSigning, userTarget int) (res *helpers.JSONResponse, err error)
	SaveDraftMultiple(ctx context.Context, userData models.Shortuser, dataReq models.DocumentUserMultiple) (res *helpers.JSONResponse, err error)
	CountDocByUser(ctx context.Context, userData models.Shortuser) (res *helpers.JSONResponse, err error)
	SendSignMultiple(ctx context.Context, userData models.Shortuser, dataReq models.DocumentUserSendSigningMultiple) (res *helpers.JSONResponse, err error)
}

// DocumentService is
type DocumentService struct {
	DocumentRepository         repositories.IDocumentsRepository
	DocumentUserRepository     repositories.IDocumentUserRepository
	DocumentActivityRepository repositories.IDocumentActivityRepository
	CommentRepository          repositories.ICommentRepository
	DB                         infrastructures.ISQLConnection
}

// InitDocumentService init
func InitDocumentService() *DocumentService {
	documentRepositories := new(repositories.DocumentsRepository)
	documentRepositories.DB = &infrastructures.SQLConnection{}

	documentUserRepositories := new(repositories.DocumentUserRepository)
	documentUserRepositories.DB = &infrastructures.SQLConnection{}

	documentActivityRepositories := new(repositories.DocumentActivityRepository)
	documentActivityRepositories.DB = &infrastructures.SQLConnection{}

	commentRepositories := new(repositories.CommentRepository)
	commentRepositories.DB = &infrastructures.SQLConnection{}

	DocumentService := new(DocumentService)
	DocumentService.DocumentRepository = documentRepositories
	DocumentService.DocumentUserRepository = documentUserRepositories
	DocumentService.DocumentActivityRepository = documentActivityRepositories
	DocumentService.CommentRepository = commentRepositories

	return DocumentService
}

// GetDocumentUser is
func (s *DocumentService) GetDocumentUser(ctx context.Context, filter models.DocumentUserFilter, page helpers.PageReq) (res *helpers.Paginate, err error) {
	if page.Limit <= 0 {
		page.Limit = 5
	}

	if page.Page <= 0 {
		page.Page = 1
	}

	// init condition
	condition := map[string]interface{}{
		"document_user.status":  1,
		"document_user.user_id": filter.UserID,
	}

	sorting := "DESC"

	if filter.Starred > 0 {
		condition["document_user.starred"] = filter.Starred
	}

	if filter.Signed > 0 {
		condition["documents.signed"] = filter.Signed
	} else {
		condition["documents.signed"] = 0
	}

	if filter.Signing > 0 {
		condition["document_user.signing"] = filter.Signing
	}

	if filter.Shared > 0 {
		condition["document_user.shared"] = filter.Shared
	}
	if filter.FileName != "" {
		condition["documents.file_name"] = filter.FileName
	}
	if filter.Sort != "" {
		sorting = filter.Sort
	}

	dataDocs, count, err := s.DocumentUserRepository.GetDocByUser(ctx, condition, page, sorting)

	pages := helpers.NewPaginate(page.Page, page.Limit, count)

	if page.Page > pages.PageCount {
		res = &helpers.Paginate{
			Code:    4400,
			Message: fmt.Sprintf("there just have %d page", pages.PageCount),
			Error:   "true",
			Data:    nil,
		}
		return
	}
	res = helpers.WrapPaginate(pages, dataDocs)
	return
}

//UpdateDocumentAttributte update document only attribute
func (s *DocumentService) UpdateDocumentAttributte(ctx context.Context, filter models.UpdateDocReq, userData models.Shortuser) (res *helpers.JSONResponse, err error) {

	payload := map[string]interface{}{
		filter.FieldType: filter.FieldValue,
	}

	tx, err := s.DocumentRepository.Tx()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  filter,
		}).Error("[Service UpdateDocumentAttributte] error create tx")
		return
	}
	defer tx.RollbackUnlessCommitted()

	var affect int64

	if filter.FieldType == "signed" {
		condition := map[string]interface{}{
			"id": filter.DocumentID,
		}
		affect, err = s.DocumentRepository.UpdateDoc(ctx, tx, condition, payload)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"code":  5500,
				"error": err,
				"data":  filter,
			}).Error("[Service UpdateDocumentAttributte] error update documents")
			return
		}
	} else {
		condition := map[string]interface{}{
			"user_id":     filter.UserID,
			"document_id": filter.DocumentID,
		}
		affect, err = s.DocumentUserRepository.UpdateDocUsers(ctx, tx, condition, payload)

		if err != nil {
			logrus.WithFields(logrus.Fields{
				"code":  5500,
				"error": err,
				"data":  filter,
			}).Error("[Service UpdateDocumentAttributte] error update document users")
			return
		}
	}
	var actvity models.DocumentActivity

	actvity.UserID = filter.UserID
	actvity.DocumentID = filter.DocumentID
	actvity.Name = userData.Name
	actvity.NIP = userData.Nip
	actvity.Status = 1

	switch filter.FieldType {
	case "starred":
		if filter.FieldValue == 1 {
			actvity.Message = "Document has been starred"
			actvity.Type = "starred"
		} else {
			actvity.Message = "Document has been unstarred"
			actvity.Type = "starred"
		}
	case "signed":
		actvity.Message = "Document has been signed"
		actvity.Type = "signed"

	case "status":
		if filter.FieldValue == 1 {
			actvity.Message = "Document has been restore"
			actvity.Type = "status"
		} else {
			actvity.Message = "Document has been deleted"
			actvity.Type = "status"
		}
	case "shared":
		if filter.FieldValue == 1 {
			actvity.Message = "Document has been shared"
			actvity.Type = "shared"
		} else {
			actvity.Message = "Document has been unshared"
			actvity.Type = "shared"
		}
	default:
		actvity.Message = "Document oh document"
		actvity.Type = "unlisted"
	}

	_, err = s.DocumentActivityRepository.StoreDocumentActivity(ctx, tx, actvity)

	tx.Commit()

	var response *helpers.JSONResponse
	if affect > 0 {

		response = &helpers.JSONResponse{
			Code:    2200,
			Message: "Success",
			Data:    nil,
		}
	} else {
		response = &helpers.JSONResponse{
			Code:    4400,
			Message: "Failed to update or Document not found",
			Data:    nil,
		}
	}
	return response, nil
}

//GetSingleDocByUser get single document
func (s *DocumentService) GetSingleDocByUser(ctx context.Context, docID int, userData models.Shortuser) (res *helpers.JSONResponse, err error) {

	found := 1
	data, err := s.DocumentUserRepository.GetSingleDocByUser(ctx, userData.ID, docID)
	if err != nil {
		if err == dbr.ErrNotFound {
			found = 0
		} else {
			logrus.WithFields(logrus.Fields{
				"code":  5500,
				"error": err,
				"data":  docID,
			}).Error("[Service GetSingleDocByUser] error get document users")
			return
		}

	}

	var response *helpers.JSONResponse

	if found == 1 {
		response = &helpers.JSONResponse{
			Code:    2200,
			Message: "Success",
			Data:    data,
		}

	} else {
		response = &helpers.JSONResponse{
			Code:    4400,
			Message: "Not Found",
			Data:    nil,
		}
	}
	return response, nil
}

//GetActivityDoc get single document
func (s *DocumentService) GetActivityDoc(ctx context.Context, docID int, userData models.Shortuser) (res *helpers.JSONResponse, err error) {

	data, err := s.DocumentActivityRepository.GetActivityByDocID(ctx, docID)
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  docID,
		}).Error("[Service GetSingleDocByUser] error get document users")
		return

	}

	var response *helpers.JSONResponse

	response = &helpers.JSONResponse{
		Code:    2200,
		Message: "Success",
		Data:    data,
	}

	return response, nil
}

//SaveDraft save draft single document
func (s *DocumentService) SaveDraft(ctx context.Context, userData models.Shortuser, dataReq models.DocumentUser) (res *helpers.JSONResponse, err error) {

	var affect int64

	TimeNow := time.Now()
	payload := map[string]interface{}{
		"tampilan":   dataReq.Tampilan,
		"page":       dataReq.Page,
		"image":      dataReq.Image,
		"x_axis":     dataReq.XAxis,
		"y_axis":     dataReq.YAxis,
		"width":      dataReq.Width,
		"height":     dataReq.Height,
		"labels":     1,
		"updated_at": TimeNow.Format("2006-01-02 15:04:05"),
	}

	condition := map[string]interface{}{
		"user_id":     userData.ID,
		"document_id": dataReq.DocumentID,
	}

	tx, err := s.DocumentRepository.Tx()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  dataReq,
		}).Error("[Service SaveDraft] error create tx")
		return
	}

	defer tx.RollbackUnlessCommitted()
	fmt.Println("condition")
	fmt.Println(condition)
	fmt.Println("payload")
	fmt.Println(payload)
	affect, err = s.DocumentUserRepository.UpdateDocUsers(ctx, tx, condition, payload)
	if err != nil {

		return
	}

	var response *helpers.JSONResponse
	if affect > 0 {

		response = &helpers.JSONResponse{
			Code:    2200,
			Message: "Success",
			Data:    nil,
		}
	} else {
		response = &helpers.JSONResponse{
			Code:    4400,
			Message: "Failed to save draft",
			Data:    nil,
		}
	}

	tx.Commit()

	return response, nil

}

//SendSign save draft single document
func (s *DocumentService) SendSign(ctx context.Context, userData models.Shortuser, dataReq models.DocumentUserSendSigning, userTarget int) (res *helpers.JSONResponse, err error) {

	var affect int64

	TimeNow := time.Now()
	payload := map[string]interface{}{
		"tampilan":   dataReq.Tampilan,
		"page":       dataReq.Page,
		"image":      dataReq.Image,
		"x_axis":     dataReq.XAxis,
		"y_axis":     dataReq.YAxis,
		"width":      dataReq.Width,
		"height":     dataReq.Height,
		"signing":    1,
		"status":     2,
		"updated_at": TimeNow.Format("2006-01-02 15:04:05"),
	}

	condition := map[string]interface{}{
		"user_id":     userData.ID,
		"document_id": dataReq.DocumentID,
	}

	tx, err := s.DocumentRepository.Tx()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  dataReq,
		}).Error("[Service SaveDraft] error create tx")
		return
	}

	defer tx.RollbackUnlessCommitted()

	affect, err = s.DocumentUserRepository.UpdateDocUsers(ctx, tx, condition, payload)
	if err != nil {

		return
	}

	var response *helpers.JSONResponse
	if affect > 0 {

		response = &helpers.JSONResponse{
			Code:    2200,
			Message: "Success",
			Data:    nil,
		}
	} else {
		response = &helpers.JSONResponse{
			Code:    4400,
			Message: "Failed to save signing",
			Data:    nil,
		}
	}

	var docUser models.DocumentUser

	docUser.ID = dataReq.ID
	docUser.DocumentID = dataReq.DocumentID
	docUser.Starred = dataReq.Starred
	docUser.Shared = dataReq.Shared
	docUser.Signing = dataReq.Signing
	docUser.Labels = dataReq.Labels
	docUser.CreatedAt = dataReq.CreatedAt
	docUser.UpdatedAt = dataReq.UpdatedAt
	docUser.Tampilan = dataReq.Tampilan
	docUser.Page = dataReq.Page
	docUser.Image = dataReq.Image
	docUser.XAxis = dataReq.XAxis
	docUser.YAxis = dataReq.YAxis
	docUser.Width = dataReq.Width
	docUser.Height = dataReq.Height
	docUser.UserID = userTarget
	docUser.Signing = 1
	docUser.Status = 1
	docUser.CreatedAt = TimeNow.Format("2006-01-02 15:04:05")
	docUser.SourceUserID = null.IntFrom(int64(userData.ID))
	docUser.SourceUsername = null.StringFrom(userData.Name)

	idStore, err := s.DocumentUserRepository.StoreDocumentUser(ctx, tx, docUser)
	if err != nil || idStore == 0 {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  dataReq,
		}).Error("[Service SendSign] error store doc user")
		return
	}
	comment := models.Comment{
		IDDocument: dataReq.DocumentID,
		IDUser:     dataReq.UserID,
		NameUser:   userData.Name,
		Comment:    dataReq.Comment.String,
	}
	//store comment
	_, err = s.CommentRepository.StoreComment(ctx, tx, comment)

	tx.Commit()

	return response, nil

}

//SaveDraftMultiple save draft single document
func (s *DocumentService) SaveDraftMultiple(ctx context.Context, userData models.Shortuser, dataReq models.DocumentUserMultiple) (res *helpers.JSONResponse, err error) {

	docIDs := strings.Split(dataReq.DocumentID, ",")

	docDraft := models.DocumentUser{
		UserID:   userData.ID,
		Tampilan: dataReq.Tampilan,
		Page:     dataReq.Page,
		Image:    dataReq.Image,
		XAxis:    dataReq.XAxis,
		YAxis:    dataReq.YAxis,
		Width:    dataReq.Width,
		Height:   dataReq.Height,
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
				}).Error("[Service SaveDraftMultiple] error convert to int docID")
				fatalErrors <- err
				cancel()
				return
			}

			docDraft.DocumentID = intDocID

			result, err := s.SaveDraft(ctx, userData, docDraft)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"code":  5500,
					"error": err,
					"data":  docIDInside,
				}).Error("[Service SaveDraftMultiple] error save draft")
				fatalErrors <- err
				return
			}

			if result == nil {
				err1 := errors.New("not 2200")
				logrus.WithFields(logrus.Fields{
					"code":  5500,
					"error": err1,
					"data":  result,
				}).Error("[Service SaveDraftMultiple] error not 2200")
				fatalErrors <- err1
				cancel()
				return

			}

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
			"data":  docDraft,
		}).Error("[Service SaveDraftMultiple] error when save multiple draft")
		response := &helpers.JSONResponse{
			Code:    5500,
			Message: err.Error(),
			Data:    nil,
		}

		return response, err
	}

	fmt.Println("ga masuk error")
	res = &helpers.JSONResponse{
		Code:    2200,
		Message: "success",
		Data:    nil,
	}
	return

}

//CountDocByUser update document only attribute
func (s *DocumentService) CountDocByUser(ctx context.Context, userData models.Shortuser) (res *helpers.JSONResponse, err error) {

	data, err := s.DocumentUserRepository.CountDocByUser(ctx, userData.ID)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  userData,
		}).Error("[Service CountDocByUser] error get count user")
		return
	}

	res = &helpers.JSONResponse{
		Code:    2200,
		Message: "success",
		Data:    data,
	}

	return
}

//SendSignMultiple send sign multiple document
func (s *DocumentService) SendSignMultiple(ctx context.Context, userData models.Shortuser, dataReq models.DocumentUserSendSigningMultiple) (res *helpers.JSONResponse, err error) {

	documentIDs := strings.Split(dataReq.DocumentID, ",")

	docSign := models.DocumentUserSendSigning{
		Tampilan: dataReq.Tampilan,
		Page:     dataReq.Page,
		Image:    dataReq.Image,
		XAxis:    dataReq.XAxis,
		YAxis:    dataReq.YAxis,
		Width:    dataReq.Width,
		Height:   dataReq.Height,
		Comment:  dataReq.Comment,
	}

	fatalErrors := make(chan error)
	wgDone := make(chan bool)

	newctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	for _, documentID := range documentIDs {
		wg.Add(1)
		go func(documentIDInside string) {
			defer wg.Done()

			select {
			case <-newctx.Done():
				return // Error somewhere, terminate
			default: // Default is must to avoid blocking
			}
			intdocumentID, err := strconv.Atoi(documentIDInside)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"code":  5500,
					"error": err,
					"data":  documentIDInside,
				}).Error("[Service SendSignMultiple] error convert to int documentID")
				fatalErrors <- err
				cancel()
				return
			}

			docSign.DocumentID = intdocumentID

			result, err := s.SendSign(ctx, userData, docSign, dataReq.TargetID)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"code":  5500,
					"error": err,
					"data":  documentIDInside,
				}).Error("[Service SendSignMultiple] error send sign")
				fatalErrors <- err
				return
			}

			if result == nil {
				err1 := errors.New("not 2200")
				logrus.WithFields(logrus.Fields{
					"code":  5500,
					"error": err1,
					"data":  result,
				}).Error("[Service SendSignMultiple] error not 2200")
				fatalErrors <- err1
				cancel()
				return

			}

		}(documentID)
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
			"data":  docSign,
		}).Error("[Service SendSignMultiple] error when save multiple draft")
		response := &helpers.JSONResponse{
			Code:    5500,
			Message: err.Error(),
			Data:    nil,
		}

		return response, err
	}

	fmt.Println("ga masuk error")
	res = &helpers.JSONResponse{
		Code:    2200,
		Message: "success",
		Data:    nil,
	}
	return

}
