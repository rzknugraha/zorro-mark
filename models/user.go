package models

import (
	"time"
)

//User data struct
type User struct {
	ID          int
	IDDpr       int
	Nama        string
	Ktp         string
	NamaJabatan string
	NamaSatker  string
	Status      int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	NIP         string
	IDSatker    int
	IDSubSatker int
	Password    string
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
