package repositories

import (

	//"github.com/afex/hystrix-go/hystrix"

	"context"
	"fmt"

	dbr "github.com/gocraft/dbr/v2"
	"github.com/rzknugraha/zorro-mark/helpers"
	"github.com/rzknugraha/zorro-mark/infrastructures"
	"github.com/rzknugraha/zorro-mark/models"
	"github.com/sirupsen/logrus"
	"go.elastic.co/apm"
)

// IDocumentUserRepository is
type IDocumentUserRepository interface {
	StoreDocumentUser(ctx context.Context, db *dbr.Tx, doc models.DocumentUser) (idDocUser int64, err error)
	GetDocByUser(ctx context.Context, conditon map[string]interface{}, paging helpers.PageReq) (dataDocs []models.DocumentUserJoinDoc, total int, err error)
}

// DocumentUserRepository is
type DocumentUserRepository struct {
	DB infrastructures.ISQLConnection
}

// StoreDocumentUser store agent type data to database
func (r *DocumentUserRepository) StoreDocumentUser(ctx context.Context, db *dbr.Tx, doc models.DocumentUser) (idDocUser int64, err error) {
	span, _ := apm.StartSpan(ctx, "StoreDocumentUser", "DocumentUserRepository")
	defer span.End()

	res, err := db.InsertInto("document_user").
		Columns(
			"document_id",
			"user_id",
			"starred",
			"shared",
			"signing",
			"labels",
			"signed",
			"created_at",
			"status").
		Record(&doc).ExecContext(ctx)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  doc,
		}).Error("[REPO StoreDocumentUser] error store DB")
		return
	}

	idDocUser, err = res.LastInsertId()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  doc,
		}).Error("[REPO StoreDocumentUser] error get last ID")
		return
	}

	return
}

//GetDocByUser get document spesific user
func (r *DocumentUserRepository) GetDocByUser(ctx context.Context, conditon map[string]interface{}, paging helpers.PageReq) (dataDocs []models.DocumentUserJoinDoc, total int, err error) {

	db := r.DB.EsignRead()

	q := db.Select(
		"document_user.id",
		"document_user.document_id",
		"document_user.user_id",
		"document_user.starred",
		"document_user.shared",
		"document_user.signing",
		"document_user.labels",
		"document_user.signed",
		"document_user.created_at",
		"document_user.updated_at",
		"document_user.status",
		"documents.file_name",
		"documents.path",
	).
		From("document_user").
		LeftJoin("documents", "document_user.document_id = documents.id")

	for key, val := range conditon {
		if key == "documents.file_name" {
			stringVal := fmt.Sprintf("%s", val)
			q.Where(key+" like ?", "%"+stringVal+"%")
		} else {
			q.Where(key+" = ?", val)
		}
	}

	var countData []models.DocumentUserJoinDoc
	total, err = q.Load(&countData)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  countData,
		}).Error("[REPO GetDocByUser]error when count data from DB")

		return
	}

	q.Limit(uint64(paging.Limit))
	q.Offset(uint64(paging.Offset))

	_, err = q.LoadContext(ctx, &dataDocs)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  conditon,
		}).Error("[REPO GetDocByUser]error get doc from DB")

		return
	}

	return
}
