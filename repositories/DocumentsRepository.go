package repositories

import (

	//"github.com/afex/hystrix-go/hystrix"

	"context"
	"fmt"

	dbr "github.com/gocraft/dbr/v2"
	"github.com/rzknugraha/zorro-mark/infrastructures"
	"github.com/rzknugraha/zorro-mark/models"
	"github.com/sirupsen/logrus"
	"go.elastic.co/apm"
)

// IDocumentsRepository is
type IDocumentsRepository interface {
	Tx() (tx *dbr.Tx, err error)
	StoreDocuments(ctx context.Context, db *dbr.Tx, doc models.Documents) (idDoc int64, err error)
	CountNameByUserID(ctx context.Context, doc models.Documents) (totalID int64, err error)
	UpdateDoc(ctx context.Context, db *dbr.Tx, Condition map[string]interface{}, Payload map[string]interface{}) (affect int64, err error)
	GetDoc(ctx context.Context, conditon map[string]interface{}) (doc models.Documents, err error)
}

// DocumentsRepository is
type DocumentsRepository struct {
	DB infrastructures.ISQLConnection
}

// Tx init a new transaction
func (r *DocumentsRepository) Tx() (tx *dbr.Tx, err error) {
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

// StoreDocuments store agent type data to database
func (r *DocumentsRepository) StoreDocuments(ctx context.Context, db *dbr.Tx, doc models.Documents) (idDoc int64, err error) {
	span, _ := apm.StartSpan(ctx, "StoreDocuments", "DocumentsRepository")
	defer span.End()

	res, err := db.InsertInto("documents").
		Columns(
			"created_by",
			"file_name",
			"path",
			"created_at",
			"status").
		Record(&doc).ExecContext(ctx)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  doc,
		}).Error("[REPO StoreDocuments] error store DB")
		return
	}

	idDoc, err = res.LastInsertId()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  doc,
		}).Error("[REPO StoreDocuments] error get last ID")
		return
	}

	return
}

//CountNameByUserID count same filename by user ID
func (r *DocumentsRepository) CountNameByUserID(ctx context.Context, doc models.Documents) (totalID int64, err error) {

	db := r.DB.EsignRead()

	_, err = db.Select("count(*)").From("documents").
		Where("created_by = ?", doc.CreatedBy).
		Where("file_name like ?", doc.FileName+"%").
		LoadContext(ctx, &totalID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  doc,
		}).Error("[REPO CountNameByUserID] error count filename from DB")

		return
	}

	return
}

//GetDoc get document
func (r *DocumentsRepository) GetDoc(ctx context.Context, conditon map[string]interface{}) (doc models.Documents, err error) {

	db := r.DB.EsignRead()

	q := db.Select("*").From("documents")

	for key, val := range conditon {
		q.Where(key+" = ?", val)
	}

	err = q.LoadOneContext(ctx, &doc)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  conditon,
		}).Error("[REPO GetDocByUserID]error count filename from DB")

		return
	}

	return
}

// UpdateDoc func
func (r *DocumentsRepository) UpdateDoc(ctx context.Context, db *dbr.Tx, Condition map[string]interface{}, Payload map[string]interface{}) (affect int64, err error) {
	span, _ := apm.StartSpan(ctx, "UpdateDoc", "DocumentsRepository")
	defer span.End()

	up := db.Update("documents")

	for key, val := range Condition {
		up.Where(key+" = ?", val)
	}
	fmt.Println("up")
	fmt.Println(up)
	result, err := up.SetMap(Payload).ExecContext(ctx)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  Condition,
		}).Error("[REPO UpdateDoc] error update")
	}
	affect, _ = result.RowsAffected()

	return
}
