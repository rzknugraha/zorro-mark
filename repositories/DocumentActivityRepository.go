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
)

// IDocumentActivityRepository is
type IDocumentActivityRepository interface {
	StoreDocumentActivity(ctx context.Context, db *dbr.Tx, doc models.DocumentActivity) (idDocUser int64, err error)
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
	doc.CreatedAt = TimeNow.Format("2006-01-02 15:04:05")

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
