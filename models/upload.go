package models

//UploadResp data struct
type UploadResp struct {
	FileName string `json:"file_name" validate:"required"`
}

//FileReq request get file
type FileReq struct {
	Path string `json:"path" validate:"required"`
}
