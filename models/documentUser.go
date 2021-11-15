package models

import "gopkg.in/guregu/null.v3"

//DocumentUser data struct
type DocumentUser struct {
	ID             int         `db:"id" json:"id"`
	DocumentID     int         `db:"document_id" json:"document_id"`
	UserID         int         `db:"user_id" json:"user_id"`
	Starred        int         `db:"starred" json:"starred"`
	Shared         int         `db:"shared" json:"shared"`
	Signing        int         `db:"signing" json:"signing"`
	Labels         int         `db:"labels" json:"labels"`
	CreatedAt      string      `db:"created_at" json:"created_at"`
	UpdatedAt      string      `db:"updated_at" json:"updated_at"`
	Status         int         `db:"status" json:"status"`
	XAxis          int         `db:"x_axis" json:"x_axis"`
	YAxis          int         `db:"y_axis" json:"y_axis"`
	Width          int         `db:"width" json:"width"`
	Height         int         `db:"height" json:"height"`
	Page           int         `db:"page" json:"page"`
	Image          bool        `db:"image" json:"image"`
	Tampilan       null.String `db:"tampilan" json:"tampilan"`
	SourceUserID   null.Int    `db:"source_user_id" json:"source_user_id"`
	SourceUsername null.String `db:"source_username" json:"source_username"`
}

//DocumentUserJoinDoc data struct
type DocumentUserJoinDoc struct {
	ID             int         `db:"id" json:"id"`
	DocumentID     int         `db:"document_id" json:"document_id"`
	UserID         int         `db:"user_id" json:"user_id"`
	Starred        int         `db:"starred" json:"starred"`
	Shared         int         `db:"shared" json:"shared"`
	Signing        int         `db:"signing" json:"signing"`
	Labels         int         `db:"labels" json:"labels"`
	Signed         int         `db:"signed" json:"signed"`
	CreatedAt      string      `db:"created_at" json:"created_at"`
	UpdatedAt      string      `db:"updated_at" json:"updated_at"`
	Status         int         `db:"status" json:"status"`
	FileName       string      `db:"file_name" json:"file_name"`
	Path           string      `db:"path" json:"path"`
	XAxis          int         `db:"x_axis" json:"x_axis"`
	YAxis          int         `db:"y_axis" json:"y_axis"`
	Width          int         `db:"width" json:"width"`
	Height         int         `db:"height" json:"height"`
	Page           int         `db:"page" json:"page"`
	Image          bool        `db:"image" json:"image"`
	Tampilan       null.String `db:"tampilan" json:"tampilan"`
	SourceUserID   null.Int    `db:"source_user_id" json:"source_user_id"`
	SourceUsername null.String `db:"source_username" json:"source_username"`
}

//DocumentUserFilter data struct
type DocumentUserFilter struct {
	UserID   int    `db:"user_id" json:"user_id"`
	Starred  int    `db:"starred" json:"starred"`
	Shared   int    `db:"shared" json:"shared"`
	Signing  int    `db:"signing" json:"signing"`
	Labels   int    `db:"labels" json:"labels"`
	Signed   int    `db:"signed" json:"signed"`
	FileName string `db:"file_name" json:"file_name"`
	Sort     string `db:"sort" json:"sort"`
}

//UpdateDocReq data struct
type UpdateDocReq struct {
	FieldType  string `json:"field_type" validate:"oneof=starred signing signed labels shared status,required,alpha"`
	FieldValue int    `json:"field_value" validate:"numeric"`
	DocumentID int    `db:"document_id" json:"document_id" validate:"required"`
	UserID     int    `db:"user_id" json:"user_id"`
}

//DocumentUserMultiple data struct
type DocumentUserMultiple struct {
	ID         int         `db:"id" json:"id"`
	DocumentID string      `db:"document_id" json:"document_id"`
	UserID     int         `db:"user_id" json:"user_id"`
	Starred    int         `db:"starred" json:"starred"`
	Shared     int         `db:"shared" json:"shared"`
	Signing    int         `db:"signing" json:"signing"`
	Labels     int         `db:"labels" json:"labels"`
	CreatedAt  string      `db:"created_at" json:"created_at"`
	UpdatedAt  string      `db:"updated_at" json:"updated_at"`
	Status     int         `db:"status" json:"status"`
	XAxis      int         `db:"x_axis" json:"x_axis"`
	YAxis      int         `db:"y_axis" json:"y_axis"`
	Width      int         `db:"width" json:"width"`
	Height     int         `db:"height" json:"height"`
	Page       int         `db:"page" json:"page"`
	Image      bool        `db:"image" json:"image"`
	Tampilan   null.String `db:"tampilan" json:"tampilan"`
}

//DocumentUserSendSigning for send sign with comment
type DocumentUserSendSigning struct {
	ID         int         `db:"id" json:"id"`
	DocumentID int         `db:"document_id" json:"document_id"`
	UserID     int         `db:"user_id" json:"user_id"`
	Starred    int         `db:"starred" json:"starred"`
	Shared     int         `db:"shared" json:"shared"`
	Signing    int         `db:"signing" json:"signing"`
	Labels     int         `db:"labels" json:"labels"`
	CreatedAt  string      `db:"created_at" json:"created_at"`
	UpdatedAt  string      `db:"updated_at" json:"updated_at"`
	Status     int         `db:"status" json:"status"`
	XAxis      int         `db:"x_axis" json:"x_axis"`
	YAxis      int         `db:"y_axis" json:"y_axis"`
	Width      int         `db:"width" json:"width"`
	Height     int         `db:"height" json:"height"`
	Page       int         `db:"page" json:"page"`
	Image      bool        `db:"image" json:"image"`
	Tampilan   null.String `db:"tampilan" json:"tampilan"`
	Comment    null.String `db:"comment" json:"comment"`
}

//DocumentUserSendSigningMultiple for send multiple signing with comment
type DocumentUserSendSigningMultiple struct {
	ID         int         `db:"id" json:"id"`
	DocumentID string      `db:"document_id" json:"document_id"`
	TargetID   int         `db:"target_id" json:"target_id"`
	Starred    int         `db:"starred" json:"starred"`
	Shared     int         `db:"shared" json:"shared"`
	Signing    int         `db:"signing" json:"signing"`
	Labels     int         `db:"labels" json:"labels"`
	CreatedAt  string      `db:"created_at" json:"created_at"`
	UpdatedAt  string      `db:"updated_at" json:"updated_at"`
	Status     int         `db:"status" json:"status"`
	XAxis      int         `db:"x_axis" json:"x_axis"`
	YAxis      int         `db:"y_axis" json:"y_axis"`
	Width      int         `db:"width" json:"width"`
	Height     int         `db:"height" json:"height"`
	Page       int         `db:"page" json:"page"`
	Image      bool        `db:"image" json:"image"`
	Tampilan   null.String `db:"tampilan" json:"tampilan"`
	Comment    null.String `db:"comment" json:"comment"`
}

//CountDocUser count doc user have
type CountDocUser struct {
	UserID  int `db:"user_id" json:"user_id"`
	Starred int `db:"starred" json:"starred"`
	Shared  int `db:"shared" json:"shared"`
	Signing int `db:"signing" json:"signing"`
	Draft   int `db:"draft" json:"draft"`
	Signed  int `db:"signed" json:"signed"`
	Upload  int `db:"upload" json:"upload"`
}
