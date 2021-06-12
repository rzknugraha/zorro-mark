package models

import (
	"time"
)

//User data struct
type User struct {
	ID          int       `db:"id"`
	IDDpr       int       `db:"id_dpr"`
	Nama        string    `db:"nama"`
	Ktp         string    `db:"ktp"`
	NamaJabatan string    `db:"nama_jabatan"`
	NamaSatker  string    `db:"nama_satker"`
	Status      int       `db:"status"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	NIP         string    `db:"nip"`
	IDSatker    int       `db:"id_satker"`
	IDSubSatker int       `db:"id_suB_satker"`
	Password    string    `db:"password"`
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
