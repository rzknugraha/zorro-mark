package repositories

import (

	//"github.com/afex/hystrix-go/hystrix"

	"context"

	dbr "github.com/gocraft/dbr/v2"
	"github.com/rzknugraha/zorro-mark/infrastructures"
	"github.com/rzknugraha/zorro-mark/models"
	"github.com/sirupsen/logrus"
	"go.elastic.co/apm"
)

// ICommentRepository is
type ICommentRepository interface {
	Tx() (tx *dbr.Tx, err error)
	StoreComment(ctx context.Context, db *dbr.Tx, doc models.Comment) (idComment int64, err error)
	GetCommentByDocID(ctx context.Context, conditon map[string]interface{}) (comments []models.Comment, err error)
}

// CommentRepository is
type CommentRepository struct {
	DB infrastructures.ISQLConnection
}

// Tx init a new transaction
func (r *CommentRepository) Tx() (tx *dbr.Tx, err error) {
	db := r.DB.EsignWrite()
	tx, err = db.Begin()

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  nil,
		}).Error("Error begin transaction")
	}

	return
}

// StoreComment store agent type data to database
func (r *CommentRepository) StoreComment(ctx context.Context, db *dbr.Tx, doc models.Comment) (idComment int64, err error) {
	span, _ := apm.StartSpan(ctx, "StoreComment", "CommentRepository")
	defer span.End()

	res, err := db.InsertInto("document_comment").
		Columns(
			"id_document",
			"id_user",
			"name_user",
			"comment",
			"created_at").
		Record(&doc).ExecContext(ctx)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  doc,
		}).Error("[REPO StoreComment] error store DB")
		return
	}

	idComment, err = res.LastInsertId()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  doc,
		}).Error("[REPO StoreComment] error get last ID")
		return
	}

	return
}

//GetCommentByDocID get document spesific user
func (r *CommentRepository) GetCommentByDocID(ctx context.Context, conditon map[string]interface{}) (comments []models.Comment, err error) {

	db := r.DB.EsignRead()

	q := db.Select("*").
		From("document_comment")

	for key, val := range conditon {

		q.Where(key+" = ?", val)

	}

	_, err = q.LoadContext(ctx, &comments)
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
