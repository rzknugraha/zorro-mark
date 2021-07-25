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
	FilePath   string `json:"file_path"`
	ImagePath  string `json:"image_path"`
	NIK        string `json:"nik"`
	Passphrase string `json:"passphrase"`
	Tampilan   string `json:"tampilan"`
	Halaman    int    `json:"halaman"`
	Page       int    `json:"page"`
	Image      string `json:"image"`
	LinkQR     string `json:"linkQR"`
	XAxis      int    `json:"x_axis"`
	YAxis      int    `json:"y_axis"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	Tag        string `json:"tag"`
	Reason     string `json:"reason"`
	Location   string `json:"location"`
	Text       string `json:"text"`
}
