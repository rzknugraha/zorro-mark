package repositories

import (

	//"github.com/afex/hystrix-go/hystrix"

	"context"
	"errors"
	"time"

	dbr "github.com/gocraft/dbr/v2"
	"github.com/rzknugraha/zorro-mark/infrastructures"
	"github.com/rzknugraha/zorro-mark/models"
	"github.com/sirupsen/logrus"
	"go.elastic.co/apm"
)

// IUserRepository is
type IUserRepository interface {
	GetUserByIDDPR(IDDpr int) (user models.User, err error)
	StoreUser(user models.User) (count int64, err error)
	UpdateUserByID(id int, user models.User) (err error)
	Login(l models.Login) (auth bool, user models.User, err error)
	GetUserByNIPstore(nip string) (user models.User, err error)
	FindOneUser(ctx context.Context, Condition map[string]interface{}) (User models.User, err error)
	UpdateUserCond(ctx context.Context, db *dbr.Tx, Condition map[string]interface{}, Payload map[string]interface{}) (affect int64, err error)
	Tx() (tx *dbr.Tx, err error)
	GetAll(ctx context.Context) (user []models.Shortuser, err error)
}

// UserRepository is
type UserRepository struct {
	DB infrastructures.ISQLConnection
}

// Tx init a new transaction
func (r *UserRepository) Tx() (tx *dbr.Tx, err error) {
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

// GetUserByIDDPR store agent type data to database
func (r *UserRepository) GetUserByIDDPR(IDDpr int) (user models.User, err error) {
	db := r.DB.EsignRead()

	err = db.Select("*").From("users").Where("id_dpr = ?", IDDpr).LoadOne(&user)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  IDDpr,
		}).Error("[REPO GetUserByIDDPR] error get from DB")

		return
	}

	return
}

// StoreUser store user  data to database
func (r *UserRepository) StoreUser(user models.User) (count int64, err error) {

	db := r.DB.EsignWrite()

	tx, _ := db.Begin()

	defer tx.RollbackUnlessCommitted()
	user.CreatedAt = time.Now().Local()
	user.UpdatedAt = time.Now().Local()

	res, err := tx.InsertInto("users").
		Columns("id_dpr", "nama", "ktp", "nama_jabatan", "nama_satker", "status", "created_at", "updated_at", "password",
			"email", "handphone", "role", "provinsi", "avatar", "identity_file", "sign_file", "sr_file", "sn_certificate").
		Record(&user).
		Exec()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  user,
		}).Error("[REPO StoreUser] error get from DB")
		tx.Rollback()
		return
	}
	tx.Commit()
	count, _ = res.RowsAffected()
	return
}

// UpdateUserByID store agent type data to database
func (r *UserRepository) UpdateUserByID(id int, user models.User) (err error) {

	db := r.DB.EsignWrite()

	tx, _ := db.Begin()

	defer tx.RollbackUnlessCommitted()

	rec, err := tx.Update("users").
		Set("nama", user.Nama).
		Set("ktp", user.Ktp).
		Set("nama_jabatan", user.NamaJabatan).
		Set("nama_satker", user.NamaSatker).
		Set("updated_at", time.Now().Local()).
		Set("nip", user.NIP).
		Set("id_satker", user.IDSatker).
		Set("id_subsatker", user.IDSubSatker).
		Set("role", user.Role).
		Where("id = ?", id).Exec()

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  user,
		}).Error("[REPO UpdateUserByID] error update from DB")
		tx.Rollback()
		return
	}
	tx.Commit()
	rows, _ := rec.RowsAffected()
	if rows == 0 {
		return errors.New("No rows affected")
	}

	return
}

// Login find match username and password
func (r *UserRepository) Login(l models.Login) (auth bool, user models.User, err error) {
	db := r.DB.EsignRead()

	defer db.Close()
	_, err = db.Select("*").From("users").
		Where("nip = ?", l.Nip).
		Where("password = ?", l.Password).
		Load(&user)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  l,
		}).Error("[REPO Login] error get from DB")

		return
	}
	if user.ID == 0 {
		auth = false
	} else {
		auth = true
	}
	return
}

// GetUserByNIPstore agent type data to database
func (r *UserRepository) GetUserByNIPstore(nip string) (user models.User, err error) {
	db := r.DB.EsignRead()

	tx, _ := db.Begin()

	defer tx.RollbackUnlessCommitted()
	err = tx.Select("*").From("users").Where("nip = ?", nip).LoadOne(&user)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  nip,
		}).Error("[REPO GetUserByNIP] error get from DB")
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}

// FindOneUser func
func (r *UserRepository) FindOneUser(ctx context.Context, Condition map[string]interface{}) (User models.User, err error) {

	db := r.DB.EsignRead()

	a := db.Select("*").From("users")

	for key, val := range Condition {
		a.Where(key+" = ?", val)
	}

	_, err = a.LoadContext(ctx, &User)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  Condition,
		}).Error("[REPO FindOneUser] error get from DB")
		return
	}

	return
}

// UpdateUserCond func
func (r *UserRepository) UpdateUserCond(ctx context.Context, db *dbr.Tx, Condition map[string]interface{}, Payload map[string]interface{}) (affect int64, err error) {
	span, _ := apm.StartSpan(ctx, "UpdateUserCond", "UserRepository")
	defer span.End()

	up := db.Update("users")

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

// GetAll agent type data to database
func (r *UserRepository) GetAll(ctx context.Context) (user []models.Shortuser, err error) {
	db := r.DB.EsignRead()

	_, err = db.Select("*").From("users").LoadContext(ctx, &user)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": errors.New("Failed Get All User"),
			"data":  nil,
		}).Error("[REPO GetUserByNIP] error get from DB")

		return
	}

	return
}
