package repositories

import (

	//"github.com/afex/hystrix-go/hystrix"

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

// StoreUser store agent type data to database
func (r *UserRepository) StoreUser(user models.User) (count int64, err error) {

	db := r.DB.EsignWrite()

	tx, _ := db.Begin()

	defer tx.RollbackUnlessCommitted()
	user.CreatedAt = time.Now().Local()
	user.UpdatedAt = time.Now().Local()

	res, err := tx.InsertInto("users").
		Columns("id_dpr", "nama", "ktp", "nama_jabatan", "nama_satker", "status", "created_at", "updated_at", "nip", "id_satker", "id_subsatker", "password").
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
