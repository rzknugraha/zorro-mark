package models

import (
	"time"

	"github.com/guregu/null"
)

//User data struct
type User struct {
	ID            int         `db:"id" json:"id"`
	IDDpr         int         `db:"id_dpr" json:"id_dpr"`
	Nama          string      `db:"nama" json:"nama"`
	Ktp           string      `db:"ktp" json:"ktp"`
	NamaJabatan   string      `db:"nama_jabatan" json:"nama_jabatan"`
	NamaSatker    string      `db:"nama_satker" json:"nama_satker"`
	Status        int         `db:"status" json:"status"`
	CreatedAt     time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time   `db:"updated_at" json:"updated_at"`
	NIP           string      `db:"nip" json:"nip"`
	IDSatker      int         `db:"id_satker" json:"id_satker"`
	IDSubSatker   int         `db:"id_sub_satker" json:"id_sub_satker"`
	Email         null.String `db:"email" json:"email"`
	Handphone     null.String `db:"handphone" json:"handphone"`
	Role          null.String `db:"role" json:"role"`
	Provinsi      null.String `db:"provinsi" json:"provinsi"`
	Avatar        null.String `db:"avatar" json:"avatar"`
	IdentityFile  null.String `db:"identity_file" json:"identity_file"`
	SignFile      null.String `db:"sign_file" json:"sign_file"`
	SRFile        null.String `db:"sr_file" json:"sr_file"`
	SNCertificate null.String `db:"sn_certificate" json:"sn_certificate"`
	Password      string      `db:"password" json:"password"`
}

//UserDPR data struct
type UserDPR struct {
	ID          string `json:"id"`
	Nama        string `json:"nama"`
	Nip         string `json:"nip"`
	KTP         string `json:"KTP"`
	NamaJabatan string `json:"nama_jabatan"`
	IDSatker    string `json:"id_satker"`
	NamaSatker  string `json:"nama_satker"`
	IDSubSatker string `json:"id_subsatker"`
}

//Login for login paylioad
type Login struct {
	Nip      string `json:"nip"`
	Password string `json:"password"`
}

//Shortuser for Shortuser paylioad
type Shortuser struct {
	ID          int    `db:"id" json:"id"`
	Nip         string `json:"nip"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Avatar      string `json:"avatar"`
	IdentityNO  string `db:"identity_no" json:"identity_no"`
	NamaJabatan string `db:"nama_jabatan" json:"nama_jabatan"`
	NamaSatker  string `db:"nama_satker" json:"nama_satker"`
	SignFile    string `db:"sign_file" json:"sign_file"`
}

//ListUser for Shortuser paylioad
type ListUser struct {
	ID          int         `db:"id" json:"id"`
	Nip         null.String `json:"nip"`
	Nama        string      `db:"nama" json:"name"`
	IdentityNO  null.String `db:"identity_no" json:"identity_no"`
	NamaJabatan string      `db:"nama_jabatan" json:"nama_jabatan"`
	NamaSatker  string      `db:"nama_satker" json:"nama_satker"`
	SignFile    null.String `db:"sign_file" json:"sign_file"`
}

//EncryptedCookies for encrypted from DPR SSO
type EncryptedCookies struct {
	Cookies1 string `db:"cookies1" json:"cookies1"`
	Cookies2 string `db:"cookies2" json:"cookies2"`
}
