package models

//Documents data struct
type Documents struct {
	ID        int    `db:"id" json:"id"`
	CreatedBy int    `db:"created_by" json:"created_by"`
	FileName  string `db:"file_name" json:"file_name"`
	Path      string `db:"path" json:"path"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
	Status    int    `db:"status" json:"status"`
}
