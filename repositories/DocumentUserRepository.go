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
	GetDocByUser(ctx context.Context, conditon map[string]interface{}, paging helpers.PageReq, sorting string) (dataDocs []models.DocumentUserJoinDoc, total int, err error)
	UpdateDocUsers(ctx context.Context, db *dbr.Tx, Condition map[string]interface{}, Payload map[string]interface{}) (affect int64, err error)
	GetSingleDocByUser(ctx context.Context, userID int, documentID int) (dataDoc models.DocumentUserJoinDoc, err error)
	CountDocByUser(ctx context.Context, userID int) (count models.CountDocUser, err error)
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
			"created_at",
			"status",
			"x_axis",
			"y_axis",
			"width",
			"height",
			"page",
			"image",
			"tampilan",
			"source_user_id",
			"source_username").
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
func (r *DocumentUserRepository) GetDocByUser(ctx context.Context, conditon map[string]interface{}, paging helpers.PageReq, sorting string) (dataDocs []models.DocumentUserJoinDoc, total int, err error) {

	db := r.DB.EsignRead()

	q := db.Select(
		"document_user.id",
		"document_user.document_id",
		"document_user.user_id",
		"document_user.starred",
		"document_user.shared",
		"document_user.signing",
		"document_user.labels",
		"documents.signed",
		"document_user.created_at",
		"document_user.updated_at",
		"document_user.status",
		"documents.file_name",
		"documents.path",
		"document_user.source_user_id",
		"document_user.source_username",
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
	if sorting == "ASC" {
		q.OrderAsc("document_user.id")
	} else {
		q.OrderDesc("document_user.id")
	}

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

// UpdateDocUsers func
func (r *DocumentUserRepository) UpdateDocUsers(ctx context.Context, db *dbr.Tx, Condition map[string]interface{}, Payload map[string]interface{}) (affect int64, err error) {
	span, _ := apm.StartSpan(ctx, "UpdateClient", "DocumentUserRepository")
	defer span.End()

	up := db.Update("document_user")

	for key, val := range Condition {
		up.Where(key+" = ?", val)
	}

	result, err := up.SetMap(Payload).ExecContext(ctx)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  Condition,
		}).Error("[REPO UpdateDocUsers] error update")
	}
	affect, _ = result.RowsAffected()

	return
}

//GetSingleDocByUser get document spesific user
func (r *DocumentUserRepository) GetSingleDocByUser(ctx context.Context, userID int, documentID int) (dataDoc models.DocumentUserJoinDoc, err error) {

	db := r.DB.EsignRead()

	err = db.Select(
		"document_user.id",
		"document_user.document_id",
		"document_user.user_id",
		"document_user.starred",
		"document_user.shared",
		"document_user.signing",
		"document_user.labels",
		"documents.signed",
		"document_user.created_at",
		"document_user.updated_at",
		"document_user.status",
		"documents.file_name",
		"documents.path",
		"document_user.x_axis",
		"document_user.y_axis",
		"document_user.width",
		"document_user.height",
		"document_user.page",
		"document_user.image",
		"document_user.tampilan",
		"document_user.source_user_id",
		"document_user.source_username",
	).
		From("document_user").
		LeftJoin("documents", "document_user.document_id = documents.id").
		Where("document_user.document_id = ?", documentID).
		Where("document_user.user_id = ?", userID).
		LoadOneContext(ctx, &dataDoc)

	if err != nil {
		if err != dbr.ErrNotFound {
			logrus.WithFields(logrus.Fields{
				"code":  5500,
				"error": err,
				"data":  userID,
			}).Error("[REPO GetSingleDocByUser] error get data")
		}
		return

	}

	return
}

//CountDocByUser get document spesific user
func (r *DocumentUserRepository) CountDocByUser(ctx context.Context, userID int) (count models.CountDocUser, err error) {

	db := r.DB.EsignRead()

	query := fmt.Sprintf("select user_id, sum(signing) as signing ,count(status) as upload,sum(shared) as shared , sum(labels) as draft,sum(starred) as starred from document_user where user_id = %d and status  >= 1  GROUP BY user_id;", userID)

	q := db.SelectBySql(query)

	_, err = q.LoadContext(ctx, &count)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  userID,
		}).Error("[REPO CountDocByUser] error get doc from DB")

		return
	}

	return
}
