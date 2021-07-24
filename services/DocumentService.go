package services

import (
	"context"
	"fmt"

	"github.com/gocraft/dbr/v2"
	"github.com/rzknugraha/zorro-mark/helpers"
	"github.com/rzknugraha/zorro-mark/infrastructures"
	"github.com/rzknugraha/zorro-mark/models"
	"github.com/rzknugraha/zorro-mark/repositories"
	"github.com/sirupsen/logrus"
)

// IDocumentService is
type IDocumentService interface {
	GetDocumentUser(ctx context.Context, filter models.DocumentUserFilter, page helpers.PageReq) (res *helpers.Paginate, err error)
	UpdateDocumentAttributte(ctx context.Context, filter models.UpdateDocReq, userData models.Shortuser) (res *helpers.JSONResponse, err error)
	GetSingleDocByUser(ctx context.Context, docID int, userData models.Shortuser) (res *helpers.JSONResponse, err error)
	GetActivityDoc(ctx context.Context, docID int, userData models.Shortuser) (res *helpers.JSONResponse, err error)
}

// DocumentService is
type DocumentService struct {
	DocumentRepository         repositories.IDocumentsRepository
	DocumentUserRepository     repositories.IDocumentUserRepository
	DocumentActivityRepository repositories.IDocumentActivityRepository
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

	DocumentService := new(DocumentService)
	DocumentService.DocumentRepository = documentRepositories
	DocumentService.DocumentUserRepository = documentUserRepositories
	DocumentService.DocumentActivityRepository = documentActivityRepositories

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

	sorting := "ASC"

	if filter.Starred > 0 {
		condition["document_user.starred"] = filter.Starred
	}

	if filter.Signed > 0 {
		condition["documents.signed"] = filter.Signed
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
			"document_id": filter.DocumentID,
		}
		affect, err = s.DocumentUserRepository.UpdateDocUsers(ctx, tx, condition, payload)
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
