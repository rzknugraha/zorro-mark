package services

import (
	"context"
	"time"

	"github.com/gocraft/dbr/v2"
	"github.com/rzknugraha/zorro-mark/helpers"
	"github.com/rzknugraha/zorro-mark/infrastructures"
	"github.com/rzknugraha/zorro-mark/models"
	"github.com/rzknugraha/zorro-mark/repositories"
	"github.com/sirupsen/logrus"
)

// ICommentService is
type ICommentService interface {
	GetCommentByDocID(ctx context.Context, docID int) (res *helpers.JSONResponse, err error)
	StoreComment(ctx context.Context, userID int, dataComments models.Comment) (res *helpers.JSONResponse, err error)
}

// CommentService is
type CommentService struct {
	CommentRepository   repositories.ICommentRepository
	DocumentsRepository repositories.IDocumentsRepository
	DB                  infrastructures.ISQLConnection
}

// InitCommentService init
func InitCommentService() *CommentService {
	commentRepositories := new(repositories.CommentRepository)
	commentRepositories.DB = &infrastructures.SQLConnection{}

	documentsRepositories := new(repositories.DocumentsRepository)
	documentsRepositories.DB = &infrastructures.SQLConnection{}

	CommentService := new(CommentService)
	CommentService.CommentRepository = commentRepositories
	CommentService.DocumentsRepository = documentsRepositories

	return CommentService
}

// GetCommentByDocID is
func (s *CommentService) GetCommentByDocID(ctx context.Context, docID int) (res *helpers.JSONResponse, err error) {
	// init conditionDoc
	conditionDoc := map[string]interface{}{
		"id": docID,
	}

	_, err = s.DocumentsRepository.GetDoc(ctx, conditionDoc)
	if err != nil {
		if err == dbr.ErrNotFound {
			res = &helpers.JSONResponse{
				Code:    4400,
				Message: "Document not found for comment",
				Data:    nil,
			}
		} else {
			res = &helpers.JSONResponse{
				Code:    5500,
				Message: "Error Get Comment",
				Data:    nil,
			}
		}
		return
	}

	// init condition
	condition := map[string]interface{}{
		"id_document": docID,
	}

	dataComments, err := s.CommentRepository.GetCommentByDocID(ctx, condition)
	if err != nil {
		res = &helpers.JSONResponse{
			Code:    5500,
			Message: "Error Get Comment",
			Data:    dataComments,
		}

		return
	}
	if len(dataComments) > 0 {
		res = &helpers.JSONResponse{
			Code:    2200,
			Message: "Success",
			Data:    dataComments,
		}
	} else {
		res = &helpers.JSONResponse{
			Code:    4400,
			Message: "Not Found",
			Data:    dataComments,
		}
	}

	return
}

// StoreComment to store comment from user post
func (s *CommentService) StoreComment(ctx context.Context, userID int, dataComments models.Comment) (res *helpers.JSONResponse, err error) {

	dataComments.IDUser = userID
	TimeNow := time.Now()
	dataComments.CreatedAt = TimeNow.Format("2006-01-02 15:04:05")

	// init conditionDoc
	conditionDoc := map[string]interface{}{
		"id": dataComments.IDDocument,
	}

	_, err = s.DocumentsRepository.GetDoc(ctx, conditionDoc)
	if err != nil {
		if err == dbr.ErrNotFound {
			res = &helpers.JSONResponse{
				Code:    4400,
				Message: "Document not found for comment",
				Data:    nil,
			}
		} else {
			res = &helpers.JSONResponse{
				Code:    5500,
				Message: "Error Get Comment",
				Data:    nil,
			}
		}
		return
	}

	tx, err := s.CommentRepository.Tx()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  dataComments,
		}).Error("[Service StoreComment] error create tx")
		return
	}

	defer tx.RollbackUnlessCommitted()

	_, err = s.CommentRepository.StoreComment(ctx, tx, dataComments)
	if err != nil {
		res = &helpers.JSONResponse{
			Code:    5500,
			Message: "Error When Store",
			Data:    dataComments,
		}

		return
	}

	res = &helpers.JSONResponse{
		Code:    2200,
		Message: "Success",
		Data:    nil,
	}

	tx.Commit()

	return
}
