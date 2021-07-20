package services

import (
	"context"
	"fmt"

	"github.com/rzknugraha/zorro-mark/helpers"
	"github.com/rzknugraha/zorro-mark/infrastructures"
	"github.com/rzknugraha/zorro-mark/models"
	"github.com/rzknugraha/zorro-mark/repositories"
)

// IDocumentService is
type IDocumentService interface {
	GetDocumentUser(ctx context.Context, filter models.DocumentUserFilter, page helpers.PageReq) (res *helpers.Paginate, err error)
}

// DocumentService is
type DocumentService struct {
	DocumentRepository     repositories.IDocumentsRepository
	DocumentUserRepository repositories.IDocumentUserRepository
	DB                     infrastructures.ISQLConnection
}

// InitDocumentService init
func InitDocumentService() *DocumentService {
	documentRepositories := new(repositories.DocumentsRepository)
	documentRepositories.DB = &infrastructures.SQLConnection{}

	documentUserRepositories := new(repositories.DocumentUserRepository)
	documentUserRepositories.DB = &infrastructures.SQLConnection{}

	DocumentService := new(DocumentService)
	DocumentService.DocumentRepository = documentRepositories
	DocumentService.DocumentUserRepository = documentUserRepositories

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

	dataDocs, count, err := s.DocumentUserRepository.GetDocByUser(ctx, condition, page)

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
