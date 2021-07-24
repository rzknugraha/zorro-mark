package models

import (
	null "gopkg.in/guregu/null.v3"
)
//DocumentActivity data struct
type DocumentActivity struct {
	ID         int    `db:"id" json:"id"`
	DocumentID int    `db:"document_id" json:"document_id"`
	UserID     int    `db:"user_id" json:"user_id"`
	Type       string `db:"type" json:"type"`
	Message    string `db:"message" json:"message"`
	CreatedAt  null.String `db:"created_at" json:"created_at"`
	UpdatedAt  null.String `db:"updated_at" json:"updated_at"`
	Status     int    `db:"status" json:"status"`
	Name       string `db:"name" json:"name"`
	NIP        string `db:"nip" json:"nip"`
}
