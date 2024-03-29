package models

import (
	"image"
)

//BsreReq data struct
type BsreReq struct {
	File       string      `json:"file"`
	ImageTTD   image.Image `json:"imageTTD"`
	NIK        string      `json:"nik"`
	Passphrase string      `json:"passphrase"`
	Tampilan   string      `json:"tampilan"`
	Halaman    int         `json:"halaman"`
	Page       int         `json:"page"`
	Image      string      `json:"image"`
	LinkQR     string      `json:"linkQR"`
	XAxis      int         `json:"xAxis"`
	YAxis      int         `json:"yAxis"`
	Width      int         `json:"width"`
	Height     int         `json:"height"`
	Tag        string      `json:"tag"`
	Reason     string      `json:"reason"`
	Location   string      `json:"location"`
	Text       string      `json:"text"`
}

//EsignReq data struct
type EsignReq struct {
	DocumentID int    `json:"document_id" validate:"required"`
	FilePath   string `json:"file_path"`
	ImagePath  string `json:"image_path"`
	NIK        string `json:"nik"`
	Passphrase string `json:"passphrase" validate:"required"`
	Tampilan   string `json:"tampilan" validate:"required,oneof=invisible visible"`
	Page       int    `json:"page"`
	Image      bool   `json:"image"`
	XAxis      int    `json:"x_axis"`
	YAxis      int    `json:"y_axis"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
}

//EsignResp respon from bsre
type EsignResp struct {
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
	PathFile   string `json:"path_file"`
}

//EsignMutipleReq data struct
type EsignMutipleReq struct {
	DocumentID string `json:"document_id" validate:"required"`
	FilePath   string `json:"file_path"`
	ImagePath  string `json:"image_path"`
	NIK        string `json:"nik"`
	Passphrase string `json:"passphrase" validate:"required"`
}
