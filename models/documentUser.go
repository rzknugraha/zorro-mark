package models

//DocumentUser data struct
type DocumentUser struct {
	ID         int    `db:"id" json:"id"`
	DocumentID int    `db:"document_id" json:"document_id"`
	UserID     int    `db:"user_id" json:"user_id"`
	Starred    int    `db:"starred" json:"starred"`
	Shared     int    `db:"shared" json:"shared"`
	Signing    int    `db:"signing" json:"signing"`
	Labels     int    `db:"labels" json:"labels"`
	Signed     int    `db:"signed" json:"signed"`
	CreatedAt  string `db:"created_at" json:"created_at"`
	UpdatedAt  string `db:"updated_at" json:"updated_at"`
	Status     int    `db:"status" json:"status"`
}
