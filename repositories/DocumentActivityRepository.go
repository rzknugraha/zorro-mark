package repositories

import (

	//"github.com/afex/hystrix-go/hystrix"

	"context"
	"time"

	dbr "github.com/gocraft/dbr/v2"
	"github.com/rzknugraha/zorro-mark/infrastructures"
	"github.com/rzknugraha/zorro-mark/models"
	"github.com/sirupsen/logrus"
	"go.elastic.co/apm"
	"gopkg.in/guregu/null.v3"
)

// IDocumentActivityRepository is
type IDocumentActivityRepository interface {
	StoreDocumentActivity(ctx context.Context, db *dbr.Tx, doc models.DocumentActivity) (idDocUser int64, err error)
	GetActivityByDocID(ctx context.Context, documentID int) (activity []models.DocumentActivity, err error)
}

// DocumentActivityRepository is
type DocumentActivityRepository struct {
	DB infrastructures.ISQLConnection
}

// StoreDocumentActivity store agent type data to database
func (r *DocumentActivityRepository) StoreDocumentActivity(ctx context.Context, db *dbr.Tx, doc models.DocumentActivity) (idDocUser int64, err error) {
	span, _ := apm.StartSpan(ctx, "StoreDocumentActivity", "DocumentUserRepository")
	defer span.End()

	TimeNow := time.Now()
	time := TimeNow.Format("2006-01-02 15:04:05")

	doc.CreatedAt = null.NewString(time, true)
	res, err := db.InsertInto("document_activity").
		Columns(
			"document_id",
			"user_id",
			"type",
			"message",
			"name",
			"nip",
			"created_at",
			"status").
		Record(&doc).ExecContext(ctx)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  doc,
		}).Error("[REPO StoreDocumentActivity] error store DB")
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

//GetActivityByDocID get document spesific user
func (r *DocumentActivityRepository) GetActivityByDocID(ctx context.Context, documentID int) (activity []models.DocumentActivity, err error) {

	db := r.DB.EsignRead()

	_, err = db.Select("*").
		From("document_activity").
		Where("document_id = ?", documentID).
		LoadContext(ctx, &activity)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  documentID,
		}).Error("[REPO GetSingleDocByUser] error get data")

		return

	}

	return
}
