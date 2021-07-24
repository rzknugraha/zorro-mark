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
	CreatedAt  string `db:"created_at" json:"created_at"`
	UpdatedAt  string `db:"updated_at" json:"updated_at"`
	Status     int    `db:"status" json:"status"`
}

//DocumentUserJoinDoc data struct
type DocumentUserJoinDoc struct {
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
	FileName   string `db:"file_name" json:"file_name"`
	Path       string `db:"path" json:"path"`
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
