package models

//Comment data struct
type Comment struct {
	ID         int    `db:"id" json:"id"`
	IDDocument int    `db:"id_document" json:"id_document" validate:"required"`
	IDUser     int    `db:"id_user" json:"id_user"`
	NameUser   string `db:"name_user" json:"name_user"`
	Comment    string `db:"comment" json:"comment" validate:"required"`
	CreatedAt  string `db:"created_at" json:"created_at"`
}
