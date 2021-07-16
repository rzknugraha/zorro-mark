package repositories

import (

	//"github.com/afex/hystrix-go/hystrix"

	"context"
	"errors"
	"time"

	"github.com/rzknugraha/zorro-mark/infrastructures"
	"github.com/rzknugraha/zorro-mark/models"
	"github.com/sirupsen/logrus"
)

// IUserRepository is
type IUserRepository interface {
	GetUserByIDDPR(IDDpr int) (user models.User, err error)
	StoreUser(user models.User) (count int64, err error)
	UpdateUserByID(id int, user models.User) (err error)
	Login(l models.Login) (auth bool, user models.User, err error)
	GetUserByNIPstore(nip string) (user models.User, err error)
	FindOneUser(ctx context.Context, Condition map[string]interface{}) (User models.User, err error)
}

// UserRepository is
type UserRepository struct {
	DB infrastructures.ISQLConnection
}

// GetUserByIDDPR store agent type data to database
func (r *UserRepository) GetUserByIDDPR(IDDpr int) (user models.User, err error) {
	db := r.DB.EsignRead()

	tx, _ := db.Begin()

	defer tx.RollbackUnlessCommitted()
	err = tx.Select("*").From("users").Where("id_dpr = ?", IDDpr).LoadOne(&user)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"code":  5500,
			"error": err,
			"data":  IDDpr,
		}).Error("[REPO GetUserByIDDPR] error get from DB")
		tx.Rollback()
		return
	}
	tx.Commit()
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
