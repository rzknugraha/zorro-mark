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

// IDocumentUserRepository is
type IDocumentUserRepository interface {
	StoreDocumentUser(ctx context.Context, db *dbr.Tx, doc models.DocumentUser) (idDocUser int64, err error)
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
