package models

type ResponseSniper struct {
	Status int    `db:"status" json:"status"`
	Error  string `db:"error" json:"error"`
	Msg    string `db:"msg" json:"msg"`
	NIP    string `db:"nip" json:"nip"`
}
